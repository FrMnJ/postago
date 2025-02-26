FROM golang:1.23-alpine 

RUN apk update && apk add --no-cache \
    bash

# Create a new user and group
RUN addgroup -S postago && adduser -S postago -G postago

WORKDIR /postago
COPY . .

# Change ownership of the working directory
RUN chown -R postago:postago /postago

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go mod tidy 
RUN task build

# Switch to the new user
USER postago

CMD ["./bin/postago"]
