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

func handleCommand(command *submitScoreCommand) {
	log.Printf("received submit_score command")
}

func getCommands(commandHandler func(*submitScoreCommand)) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("submit_scores", false, false, false, false, nil)
	if err != nil {
		return err
	}

	commands, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	var forever chan struct{}

	go func() {
		for c := range commands {
			log.Printf("Command: %s", c.Body)

			var command submitScoreCommand
			err := json.Unmarshal(c.Body, &command)
			if err != nil {
				log.Print(err)
				continue
			}

			commandHandler(&command)
		}
	}()

	<-forever
	return nil
}

func main() {
	err := getCommands(handleCommand)
	if err != nil {
		log.Print(err)
	}
}
