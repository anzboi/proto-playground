package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func main() {
	target := flag.String("target", "localhost:8080", "target host:port of the service")
	rpc := flag.String("rpc", "", "fully qualified method name for the target method")
	json := flag.String("data", `{}`, "json blob holding the request data")

	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.Dial(*target, grpc.WithInsecure())
	panicIfErr(err)

	// The goal from this point is to be able to call conn.Invoke

	// first we need to discover the API specs somehow. In this example, we call the targets reflection RPCs.
	// We could also create the DescriptorSource from locally stored proto files

	descriptorSource, err := sourceFromReflection(ctx, conn)
	panicIfErr(err)

	// We can now look up a given target method
	method, ok := descriptorSource.FindMethod(protoreflect.FullName(*rpc))
	if !ok {
		panic(fmt.Errorf("rpc %s not found", *rpc))
	}

	// The method descriptor gives us all the metadata we need to invoke the method

	// dynamicpb offers a way to generate a protoreflect.Message from a descriptor. This is essentially a generic proto.Message
	// implementation coded to a specific descriptor. It is not as efficient as generated code but useable when generated code is unavailable.
	request := dynamicpb.NewMessage(method.Input())
	protojson.Unmarshal([]byte(*json), request)

	// create the response object as well
	response := dynamicpb.NewMessage(method.Output())

	// One last thing before we can invoke, we must convert the method symbol to a http path
	// the path MUST start with '/' and the method must be separated from the service by a '/'
	lastSymbolIndex := strings.LastIndex(*rpc, ".")
	methodPath := *rpc
	methodPath = "/" + methodPath[:lastSymbolIndex] + "/" + methodPath[lastSymbolIndex+1:]

	// Finally we are ready to invoke the method
	err = conn.Invoke(ctx, methodPath, request, response)
	panicIfErr(err)

	// unmarshal the response back out
	outBytes, err := protojson.Marshal(response)
	panicIfErr(err)
	fmt.Println("response:", string(outBytes))
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Gets a Descriptor source from proto reflection
//
// Send and recv messages are a little confusing because the reflection services is a streaming service.
func sourceFromReflection(ctx context.Context, conn grpc.ClientConnInterface) (DescriptorSource, error) {
	client, err := grpc_reflection_v1alpha.NewServerReflectionClient(conn).ServerReflectionInfo(ctx)
	if err != nil {
		return nil, err
	}
	defer client.CloseSend()

	// List services
	if err = client.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{},
	}); err != nil {
		return nil, err
	}
	services, err := client.Recv()
	if err != nil {
		return nil, err
	}

	// List containing files
	files := &descriptorpb.FileDescriptorSet{}
	seen := make(map[string]bool)
	for _, service := range services.GetListServicesResponse().Service {

		// List file containing a single symbol
		if err := client.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
			MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
				FileContainingSymbol: service.GetName(),
			},
		}); err != nil {
			return nil, err
		}
		file, err := client.Recv()
		if err != nil {
			return nil, err
		}

		// Load files into the file descriptor set. There should only be one by why not load them all
		for _, fileBytes := range file.GetFileDescriptorResponse().GetFileDescriptorProto() {
			fileProto := &descriptorpb.FileDescriptorProto{}
			if err := proto.Unmarshal(fileBytes, fileProto); err != nil {
				panic(err)
			}
			if seen[fileProto.GetName()] {
				continue
			}
			seen[fileProto.GetName()] = true
			files.File = append(files.File, fileProto)
		}
	}

	return descriptorSourceFromFileSet(files)
}

// converts a file descript set to a list of protoreflect.FileDescriptors
func descriptorSourceFromFileSet(fds *descriptorpb.FileDescriptorSet) (DescriptorSource, error) {
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return nil, err
	}
	descriptors := make([]protoreflect.FileDescriptor, 0, len(fds.GetFile()))
	for _, fileProto := range fds.GetFile() {
		fileDesc, err := protodesc.NewFile(fileProto, files)
		if err != nil {
			return nil, err
		}
		descriptors = append(descriptors, fileDesc)
	}
	return &fileSource{
		files: descriptors,
	}, nil
}

// DescriptorSource acts as a source for descriptors
type DescriptorSource interface {
	FindMethod(protoreflect.FullName) (protoreflect.MethodDescriptor, bool)
}

// implements protoreflect.ServiceDescriptors
type fileSource struct {
	files []protoreflect.FileDescriptor
}

var _ DescriptorSource = (*fileSource)(nil)

// Finds a method among the files in the descriptor source
func (f *fileSource) FindMethod(name protoreflect.FullName) (protoreflect.MethodDescriptor, bool) {
	for _, file := range f.files {
		// file package must match
		if !strings.HasPrefix(string(name), string(file.Package())) {
			continue
		}

		// trim the package and extract service name and method
		serviceMethod := name[len(file.Package())+1:]
		split := strings.Split(string(serviceMethod), ".")
		if len(split) != 2 {
			continue
		}
		service, method := split[0], split[1]

		// Search for the owning service
		serviceDesc := file.Services().ByName(protoreflect.Name(service))
		if serviceDesc == nil {
			continue
		}

		// Search for the method
		methodDesc := serviceDesc.Methods().ByName(protoreflect.Name(method))
		if methodDesc == nil {
			continue
		}

		// if found, return it
		return methodDesc, true
	}
	return nil, false
}
