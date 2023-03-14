package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type submitScoreCommand struct {
	Score    uint64 `json:"score"`
	Name     string `json:"name"`
	Datetime string `json:"datetime"`
	JobID    string `json:"job_id"`
}

func handleCommands(commands <-chan *submitScoreCommand) {
	for c := range commands {
		log.Printf("received submit_score command: %s", c.Name)
	}
}

func getCommands(commands chan<- *submitScoreCommand) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("submit_scores", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	queueCommands, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for c := range queueCommands {
			log.Printf("Command: %s", c.Body)

			var command submitScoreCommand
			err := json.Unmarshal(c.Body, &command)
			if err != nil {
				log.Print(err)
				continue
			}

			commands <- &command
			// TODO: break out of loop after time elapsed or max batch size
		}
		close(commands)
	}()

	<-forever
}

func main() {
	commands := make(chan *submitScoreCommand)

	var forever chan struct{}

	go getCommands(commands)
	go handleCommands(commands)

	<-forever
}
