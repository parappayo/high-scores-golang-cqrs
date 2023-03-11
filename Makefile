
install:
	go get github.com/google/uuid
	go get github.com/gin-gonic/gin
	go get github.com/rabbitmq/amqp091-go

build:
	go build -o bin/get-scores ./cmd/get-scores
	go build -o bin/post-score ./cmd/post-score

format:
	gofmt -w .
