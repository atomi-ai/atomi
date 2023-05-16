FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod download
RUN go generate ./...
RUN go build -o atomi github.com/atomi-ai/atomi
RUN go build -o init-db github.com/atomi-ai/atomi/exp/init-db

#===============================#
FROM golang:1.20

WORKDIR /app

COPY --from=builder /app/atomi /app/atomi
COPY --from=builder /app/init-db /app/init-db

ENTRYPOINT ["/app/atomi"]
