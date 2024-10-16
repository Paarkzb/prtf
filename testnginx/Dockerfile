# Temporary image to coordinate the build...
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY ./ ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /service .

# The actual output image
FROM alpine:latest AS releaser
WORKDIR /app
COPY --from=builder /service /service
ENV PORT 8080
ENTRYPOINT ["/service"]