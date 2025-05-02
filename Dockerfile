FROM golang:1.24 as builder

WORKDIR /mikrotik-exporter

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mikrotik-exporter cmd/mikrotik-exporter/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /mikrotik-exporter/mikrotik-exporter .

EXPOSE 9111

CMD ["./mikrotik-exporter"]
