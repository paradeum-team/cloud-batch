FROM golang:1.16

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

COPY ./ /cloud-batch
WORKDIR /cloud-batch

RUN git describe --tags `git rev-list --tags --max-count=1` > ./VERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-s -X gitlab.paradeum.com/pld/cloud-batch.Version=$(cat VERSION | tr -d ' \t\n\r')" -o cloud-batch

FROM alpine:3.14

RUN apk add --no-cache ca-certificates

COPY --from=0 /cloud-batch/cloud-batch /
COPY --from=0 /cloud-batch/assets /assets
COPY --from=0 /cloud-batch/VERSION /

EXPOSE 8080

CMD ["/cloud-batch"]
