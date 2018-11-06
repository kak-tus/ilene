FROM golang:1.10.3-alpine3.8 AS build

WORKDIR /go/src/github.com/kak-tus/ilene

COPY api ./api
COPY model ./model
COPY vendor ./vendor
COPY main.go .

RUN go install

FROM alpine:3.8

COPY --from=build /go/bin/ilene /usr/local/bin/ilene
COPY etc /etc/
COPY data /data/

RUN adduser -DH user

USER user

ENV \
  ILENE_DATA_DIR=/data \
  ILENE_REDIS_ADDRS=

EXPOSE 8080

CMD ["/usr/local/bin/ilene"]
