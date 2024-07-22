# unchain/golang-alpine-libtdl

FROM golang:alpine

RUN mkdir /.cache
RUN mkdir /.cache/go-build

RUN chmod 777 /.cache
RUN chmod 777 /.cache/go-build

RUN apk --no-cache --update add libtool ca-certificates build-base
