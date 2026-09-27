FROM golang:1.23-alpine AS builder

WORKDIR /src/backend
COPY backend/go.mod ./go.mod
COPY backend/go.sum ./go.sum
RUN go mod download
COPY backend ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/aggregator ./aggregator/cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/aggregator /aggregator
EXPOSE 8081
ENTRYPOINT ["/aggregator"]
