FROM golang:1.23.9-alpine

WORKDIR /app/server/subscription

COPY . .

RUN apk update && apk add --no-cache make bash

CMD ["sh", "-c", "make docker-run"]
