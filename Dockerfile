FROM golang:1.11-alpine3.8

RUN mkdir /go/src/webhooker

COPY main.go /go/src/webhooker
COPY config.json /go/src/webhooker

WORKDIR /go/src/webhooker

RUN go build -o webhooker main.go

ENTRYPOINT ["/go/src/webhooker/webhooker"]
