FROM golang:1.13.4-alpine as builder

WORKDIR /build

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ci-league cmd/*.go

WORKDIR /dist
RUN cp /build/ci-league ./ci-league
RUN cp /build/template.html template.html

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist .
EXPOSE 8000
CMD ["./ci-league"]