FROM golang:1.16

ENV GOPROXY=https://goproxy.io

COPY ./ /go/src/cloud-batch
WORKDIR /go/src/cloud-batch

RUN  git describe --tags `git rev-list --tags --max-count=1`|tr -d ' \t\n\r' > ./VERSION && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o cloud-batch

FROM alpine:3.14

RUN apk add --no-cache ca-certificates

COPY --from=0 /go/src/cloud-batch/cloud-batch /
COPY --from=0 /go/src/cloud-batch/assets /assets
COPY --from=0 /go/src/cloud-batch/VERSION /

EXPOSE 8080

CMD ["/cloud-batch"]
