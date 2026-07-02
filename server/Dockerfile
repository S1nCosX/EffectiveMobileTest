FROM golang:1.26.4-alpine3.24 AS builder
WORKDIR /src
COPY ./src .
RUN GOOS=linux go build -a -o server .

FROM alpine:3.24
WORKDIR /root/
COPY --from=builder /src/server .
CMD ["./server"]