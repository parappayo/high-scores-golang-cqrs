
install:
	go get github.com/google/uuid
	go get github.com/gin-gonic/gin
	go get github.com/rabbitmq/amqp091-go
	go get github.com/lib/pq

build:
	go build -o bin/get-scores ./cmd/get-scores
	go build -o bin/post-score ./cmd/post-score
	go build -o bin/scores-worker ./cmd/scores-worker

format:
	gofmt -w .
