package examples

import (
	"encoding/json"
	"fmt"

	"github.com/anzboi/proto-playground/pkg/types"
	"google.golang.org/protobuf/proto"
)

func ExampleMarshal() {
	m := &types.Message{
		IntField:    1,
		StringField: "hello",
	}
	raw, _ := proto.Marshal(m)

	fmt.Println("Message:", m)
	fmt.Println("bytes:", raw)

	// Output: int_field:1 string_field:"hello"
	// [8 1 18 5 104 101 108 108 111]
}

func ExampleProtoVSJson() {
	m := &types.Message{
		IntField:    1,
		StringField: "hello",
	}
	protoBytes, _ := proto.Marshal(m)
	jsonBytes, _ := json.Marshal(m)

	fmt.Println("proto:", protoBytes)
	fmt.Println("json: ", jsonBytes)
	fmt.Println("reduction: ", float64(len(protoBytes))/float64(len(jsonBytes)))

	// Output: proto: [8 1 18 5 104 101 108 108 111]
	// json:  [123 34 105 110 116 95 102 105 101 108 100 34 58 49 44 34 115 116 114 105 110 103 95 102 105 101 108 100 34 58 34 104 101 108 108 111 34 125]
	// reduction:  0.23684210526315788
}

func ExampleProtoInterchangeability() {
	m1 := types.Message{
		IntField:    1,
		StringField: "hello",
	}

	m2 := types.MessageV2{}

	raw, _ := proto.Marshal(&m1)
	_ = proto.Unmarshal(raw, &m2)

	fmt.Println("MessageV1:", &m1)
	fmt.Println("MessageV2:", &m2)
	fmt.Println("bytes:", raw)

	// Output: MessageV1: int_field:1 string_field:"hello"
	// MessageV2: int_field:1 string_field:"hello"
	// bytes: [8 1 18 5 104 101 108 108 111]
}

func ExampleProtoWireCompatibility() {
	m1 := types.Message{
		IntField:    1,
		StringField: "hello",
	}

	m3 := types.MessageV3{}

	raw, _ := proto.Marshal(&m1)
	_ = proto.Unmarshal(raw, &m3)

	fmt.Println("MessageV1:", &m1)
	fmt.Println("MessageV3:", &m3)
	fmt.Println("bytes:", raw)

	// Output: MessageV1: int_field:1 string_field:"hello"
	// MessageV3: int_field:1 string_field:"hello"
	// bytes: [8 1 18 5 104 101 108 108 111]
}

func ExampleProtoWireCompatibilityBackwards() {
	m3 := types.MessageV3{
		IntField:    1,
		StringField: "hello",
		StringMap:   map[string]string{"foo": "bar"},
	}

	m1 := types.Message{}

	raw, _ := proto.Marshal(&m3)
	_ = proto.Unmarshal(raw, &m1)

	fmt.Println("MessageV3:", &m3)
	fmt.Println("MessageV1:", &m1)
	fmt.Println("bytes:", raw)

	// Output: MessageV3: int_field:1  string_field:"hello"  string_map:{key:"foo"  value:"bar"}
	// MessageV1: int_field:1  string_field:"hello"  3:"\n\x03foo\x12\x03bar"
	// bytes: [8 1 18 5 104 101 108 108 111 26 10 10 3 102 111 111 18 3 98 97 114]
}

func ExampleProtoJSONBreak() {
	m1 := types.Message{
		IntField:    1,
		StringField: "hello",
	}

	m4 := types.MessageV4{}

	raw, _ := proto.Marshal(&m1)
	_ = proto.Unmarshal(raw, &m4)

	json1, _ := json.Marshal(m1)
	json4, _ := json.Marshal(m4)

	fmt.Println("MessageV1:", &m1)
	fmt.Println("Json1:", string(json1))
	fmt.Println()
	fmt.Println("MessageV4:", &m4)
	fmt.Println("Json4:", string(json4))

	// Output: MessageV1: int_field:1 string_field:"hello"
	// Json1: {"int_field":1,"string_field":"hello"}

	// MessageV4: foo:1 string_field:"hello"
	// Json4: {"foo":1,"string_field":"hello"}
}

func ExampleEnum() {
	e := types.Status{
		Code: types.Code_Code_ONE,
	}

	raw, _ := proto.Marshal(&e)
	fmt.Println("status:", &e)
	fmt.Println("bytes:", raw)

	// Output: status: code:Code_ONE
	// bytes: [8 1]
}

func ExampleEnumNotExists() {
	// Set marshalled proto message manually
	raw := []byte{8, 2}
	e := types.Status{}
	_ = proto.Unmarshal(raw, &e)

	fmt.Println("status:", &e)

	// set enum field to an unknown value
	raw = []byte{8, 3}
	_ = proto.Unmarshal(raw, &e)

	fmt.Println("status:", &e)

	// Output: status: code:Code_TWO
	// status: code:3
}
