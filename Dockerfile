FROM ankrnetwork/alpine:v1.0.0 AS builder
LABEL stage=builder
RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY ./ ./
WORKDIR /go/src/app/chat_server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o test_circleci .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/app/chat_server/test_circleci .
CMD ["./test_circleci"]
