FROM golang:1.21-alpine as builder
ADD . /go/src/apps/signer/
WORKDIR /go/src/apps/signer/
RUN go build -o /go/bin/signer ./cmd/signer

FROM alpine:latest
COPY --from=builder /go/bin/signer /usr/local/bin/signer
CMD [ "/usr/local/bin/signer" ]
EXPOSE 8000