FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /synology-dsm-notification-lark

FROM alpine:latest
COPY --from=builder /synology-dsm-notification-lark /synology-dsm-notification-lark

EXPOSE 8080

CMD [ "./synology-dsm-notification-lark" ]