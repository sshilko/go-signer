FROM golang:1.21-alpine as builder
ADD . /go/src/apps/worker/
WORKDIR /go/src/apps/worker/
RUN go build -o /go/bin/worker ./cmd/worker

FROM alpine:latest
COPY --from=builder /go/bin/worker /usr/local/bin/worker
CMD [ "/usr/local/bin/worker" ]