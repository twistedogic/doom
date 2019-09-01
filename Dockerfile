FROM golang:1.12.1-alpine3.9 AS build-env
RUN apk add --no-cache git mercurial musl-dev gcc
WORKDIR /workspace
COPY ./ ./
RUN go build -o goapp doom.go

# final stage
FROM alpine
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=build-env /workspace/goapp /app/
ENTRYPOINT ./goapp
