FROM golang:latest

WORKDIR /golang-training-university-grpc

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build ./client/cmd/main.go

CMD ["./main"]