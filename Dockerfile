FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod download
RUN go generate ./...
RUN go build -o atomi github.com/atomi-ai/atomi

FROM golang:1.20

WORKDIR /app

COPY --from=builder /app/atomi /app/atomi

ENTRYPOINT ["/app/atomi"]
