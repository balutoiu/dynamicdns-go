FROM golang:1.13-alpine as build

WORKDIR /go/src/app
ADD . /go/src/app/

RUN CGO_ENABLED=0 go build -o /go/bin/dynamicdns-go

FROM alpine:latest
COPY --from=build /go/bin/dynamicdns-go /

CMD ["/dynamicdns-go"]
