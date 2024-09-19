# syntax=docker/dockerfile:1
FROM golang:1.22.1  AS build-stage
 
WORKDIR /app
 
COPY ./ ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /prtf-auth-server ./cmd/sso


FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM scratch AS build-release-stage

WORKDIR /app

COPY --from=build-stage /config ./
COPY --from=build-stage /prtf-auth-server /prtf-auth-server

# Bind TCP port
EXPOSE 8085

# Run
ENTRYPOINT [ "/prtf-auth-server" ]