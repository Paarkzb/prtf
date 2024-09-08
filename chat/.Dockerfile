FROM golang:1.21.4  AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /prtf-chat-server ./cmd

FROM scratch AS build-release-stage

WORKDIR /app

COPY --from=build-stage /prtf-chat-server /prtf-chat-server

EXPOSE 8071

ENTRYPOINT [ "/prtf-chat-server" ]