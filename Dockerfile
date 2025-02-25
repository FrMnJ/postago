FROM golang:1.23-alpine 

RUN apk update && apk add --no-cache \
    bash

WORKDIR /postago
COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go mod install 

CMD ["./bin/postago"]
