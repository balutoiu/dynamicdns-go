FROM alpine:3.18.0 as alpine
RUN apk add -U --no-cache ca-certificates

FROM golang:latest as gobuilder
ADD . /dynamicdns-go
WORKDIR /dynamicdns-go
RUN make build

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /dynamicdns-go/_output/dynamicdns-go /dynamicdns-go
CMD ["/dynamicdns-go"]
