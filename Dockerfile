FROM golang:1.20 as build
WORKDIR /go/src/github.com/kubernetescode-aaserver
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux GOPROXY="https://goproxy.cn" go build .

FROM alpine:3.14
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/kubernetescode-aaserver/kubernetescode-aaserver /
ENTRYPOINT ["/kubernetescode-aaserver"]