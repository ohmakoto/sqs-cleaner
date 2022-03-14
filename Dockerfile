FROM golang:1.17.8 as builder

WORKDIR /go/build
COPY . .
RUN go mod download \
    && go build -o sqs-cleaner


FROM gcr.io/distroless/base-debian11

COPY --from=builder /go/build/sqs-cleaner /sqs-cleaner

ENTRYPOINT ["/sqs-cleaner"]