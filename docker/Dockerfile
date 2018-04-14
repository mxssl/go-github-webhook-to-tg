FROM golang:1.10.1-alpine as builder

# install dep package manager
RUN apk add --no-cache ca-certificates curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# build app
WORKDIR /go/src/go-github-webhook-to-tg
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o gobot

# copy app binary file to new clear image
FROM alpine:latest
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/go-github-webhook-to-tg/gobot .
RUN chmod +x gobot
CMD ["./gobot"]
