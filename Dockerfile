FROM golang:1.20-alpine AS build

WORKDIR /app
COPY ./ /app
RUN go mod download && go build -o bin/alice_kuzya_proxy main.go

FROM alpine:3.16
WORKDIR /app/
COPY --from=build /app/bin /app/

CMD ["/app/alice_kuzya_proxy"]
