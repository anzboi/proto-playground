FROM golang:1.14.4-buster AS builder
ENV GO111MODULE=on
WORKDIR /workspace
COPY . .
RUN go install ./chatter/chatserver

FROM debian:buster-slim
COPY --from=builder /go/bin/chatserver /bin/chatserver
EXPOSE 8080
CMD ["chatserver", "--port=:8080"]
