FROM golang:1.22.5-alpine  AS build-stage
 

WORKDIR /app/

COPY go.mod go.sum /app/
RUN go mod download

COPY ./ /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd


FROM build-stage AS run-test-stage
RUN go test ./tests -count=1 -v

FROM alpine:latest AS build-release-stage

WORKDIR /app/

COPY --from=0 /app/server ./
# COPY --from=0 /app/config config/
# COPY --from=0 /app/.env ./

ENTRYPOINT [ "./server" ]

