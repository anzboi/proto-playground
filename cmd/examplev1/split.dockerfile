FROM golang:1.14.4-buster as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd/examplev1 cmd/examplev1
COPY pkg pkg

RUN go build -o bin/examplev1 ./cmd/examplev1/example-split

FROM gcr.io/distroless/base:latest
COPY --from=builder /workspace/bin/examplev1 /bin/examplev1
EXPOSE 8080 8081
CMD [ "/bin/examplev1" ]