EXAMPLE1_MUXED_IMAGE?=example-muxed
EXAMPLE1_CMUXED_IMAGE?=example-cmuxed
EXAMPLE1_SPLIT_IMAGE?=example-split

examplev1: bin/examplev1-muxed.image bin/examplev1-cmuxed.image bin/examplev1-split.image

bin/examplev1-muxed.image:
	mkdir -p bin
	docker build -t ${EXAMPLE1_MUXED_IMAGE} -f ./cmd/examplev1/muxed.dockerfile .
	echo ${EXAMPLE1_MUXED_IMAGE} >bin/examplev1-muxed.image

bin/examplev1-cmuxed.image:
	mkdir -p bin
	docker build -t ${EXAMPLE1_CMUXED_IMAGE} -f ./cmd/examplev1/cmuxed.dockerfile .
	echo ${EXAMPLE1_CMUXED_IMAGE} >bin/examplev1-cmuxed.image

bin/examplev1-split.image:
	mkdir -p bin
	docker build -t ${EXAMPLE1_SPLIT_IMAGE} -f ./cmd/examplev1/split.dockerfile .
	echo ${EXAMPLE1_SPLIT_IMAGE} >bin/examplev1-split.image

clean:
	rm -rf bin