FROM golang:1-alpine

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod ./

COPY go.sum* ./

RUN go mod download || true

COPY . .

RUN go build -o main .

EXPOSE 8081

CMD ["./main"]
