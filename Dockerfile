# ---------- BUILD STAGE ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build only the server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o grpc-server ./server

# ---------- RUNTIME STAGE ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/grpc-server .

EXPOSE 50051

USER nonroot:nonroot

CMD ["/app/grpc-server"]
