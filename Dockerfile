FROM golang:1.21-alpine3.19 AS builder

COPY . /GRPC_SSO/

WORKDIR /GRPC_SSO/

RUN go mod download

RUN go build -o ./bin/app cmd/sso/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 GRPC_SSO/bin/app .
COPY --from=0 GRPC_SSO/config config/

EXPOSE 80

CMD ["./app"]