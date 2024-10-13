FROM golang:1.22.1-alpine  AS build-stage
 
WORKDIR /app
 
COPY ./ ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/sso


FROM build-stage AS run-test-stage
RUN go test ./tests -count=1 -v

FROM scratch AS build-release-stage

WORKDIR /app

COPY --from=build-stage /server /server
COPY config/config_local.yaml /config/config_local.yaml

# Bind TCP port
EXPOSE 8085

# Run
ENTRYPOINT [ "/server" ]

