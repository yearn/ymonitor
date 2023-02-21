FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/yearn/ymonitor
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/ymonitor

FROM scratch
COPY --from=builder /go/bin/ymonitor /go/bin/ymonitor
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
ENTRYPOINT ["/go/bin/ymonitor"]
