FROM golang:1.22.1-alpine  AS build-stage
 
 
COPY ./ /app/
WORKDIR /app/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/sso


FROM build-stage AS run-test-stage
RUN go test ./tests -count=1 -v

FROM alpine:latest AS build-release-stage

WORKDIR /app/

COPY --from=0 /app/server ./
COPY --from=0 /app/config config/
COPY --from=0 /app/.env ./

# Bind TCP port
# EXPOSE 8085

# Run
ENTRYPOINT [ "./server" ]

