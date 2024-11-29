FROM golang:1.21.4-alpine  AS build-stage
 
WORKDIR /app/

COPY go.mod go.sum /app/
RUN go mod download

COPY ./ /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /app/

COPY --from=0 /app/server ./
COPY --from=0 /app/configs/ configs/
COPY --from=0 /app/.env ./

# Bind TCP port
# EXPOSE 8086

# Run
ENTRYPOINT [ "./server" ]

