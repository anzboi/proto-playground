EXAMPLE1_MUXED_IMAGE?=example-muxed
EXAMPLE1_SPLIT_IMAGE?=example-split

examplev1: bin/examplev1-muxed.image bin/examplev1-split.image

bin/examplev1-muxed.image:
	docker build -t ${EXAMPLE1_MUXED_IMAGE} -f ./cmd/examplev1/muxed.dockerfile .
	echo ${EXAMPLE1_MUXED_IMAGE} >bin/examplev1-muxed.image

bin/examplev1-split.image:
	docker build -t ${EXAMPLE1_SPLIT_IMAGE} -f ./cmd/examplev1/split.dockerfile .
	echo ${EXAMPLE1_SPLIT_IMAGE} >bin/examplev1-split.image

clean:
	rm -rf bin