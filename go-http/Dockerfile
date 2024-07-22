FROM golang:1.15.2-alpine as GO_BUILD
COPY . /go-http
WORKDIR /go-http

RUN go build -o /go/bin/go-http

FROM alpine:3.10
WORKDIR app
COPY --from=GO_BUILD /go/bin/go-http ./
EXPOSE 8080
CMD ./go-http