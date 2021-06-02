FROM golang:1.16 AS builder

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/deliveryIndex

FROM alpine:3.10

COPY --from=builder /app/bin/deliveryIndex /app

RUN apk --no-cache add ca-certificates

RUN chmod 777 ./app

ENTRYPOINT ["./app"]