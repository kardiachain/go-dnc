FROM golang:1.19.4 as builder
RUN mkdir -p "$GOPATH/src/github.com/kardiachain/go-kardia"
WORKDIR /go/src/github.com/kardiachain/go-kardia
RUN apt-get update && apt-get install -y libzmq3-dev
ADD . .
WORKDIR /go/src/github.com/kardiachain/go-kardia/cmd
RUN go install

FROM alpine:3.18
RUN apk add ca-certificates
ENV PATH="${PATH}:/go/bin"
WORKDIR /go/bin
COPY --from=builder /go/bin/* .
COPY --from=builder /go/src/github.com/kardiachain/go-kardia/cmd/cfg .

ENTRYPOINT ["./cmd"]
