FROM golang:1.18 as builder
WORKDIR /investment
COPY . .

ENV CGO_ENABLED=0

RUN go build -o ./bin/investment ./cmd/investment/main.go

FROM alpine

RUN apk add --no-cache tzdata

COPY --from=builder /investment/bin/investment /investment

EXPOSE 5000

CMD ["/investment"]
