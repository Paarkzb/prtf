# syntax=docker/dockerfile:1
FROM golang:1.21.4  AS build-stage
 
# Set destination for COPY
WORKDIR /app
 
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY ./ ./
 
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /prtf-quiz-server ./cmd

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage

WORKDIR /app

COPY /configs/config.yaml /app/configs/config.yaml
COPY --from=build-stage /prtf-quiz-server /prtf-quiz-server

# Bind TCP port
EXPOSE 8086

# Run
ENTRYPOINT [ "/prtf-quiz-server" ]