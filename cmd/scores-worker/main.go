package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func getCommands() error {
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

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("hello"),
		})
	if err != nil {
		return err
	}

	commands, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for c := range commands {
		log.Printf("Command: %s", c.Body)
	}
	return nil
}

func main() {
	err := getCommands()
	if err != nil {
		log.Print(err)
	}
}
