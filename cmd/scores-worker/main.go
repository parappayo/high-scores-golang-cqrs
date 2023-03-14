package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type submitScoreCommand struct {
	Score    uint64 `json:"score"`
	EventID  string `json:"event_id"`
	Name     string `json:"name"`
	Datetime string `json:"datetime"`
	Region   string `json:"region"`
	JobID    string `json:"job_id"`
}

func getCommands(commands chan<- *submitScoreCommand) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("submit_scores", false, false, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}

	queueCommands, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for c := range queueCommands {
		var command submitScoreCommand
		err := json.Unmarshal(c.Body, &command)
		if err != nil {
			log.Println(err)
			continue
		}

		commands <- &command
		// TODO: break out of loop after time elapsed or max batch size
	}
	close(commands)
}

func handleCommands(commands <-chan *submitScoreCommand) {
	db, err := sql.Open("postgres", "postgresql://postgres@postgres")
	if err != nil {
		log.Fatal(err)
	}

	for c := range commands {
		result, err := db.Exec(
			"INSERT INTO scores (event_id, username, score, submitted_on, region) VALUES ($1, $2, $3, $4, $5)",
			c.EventID,
			c.Name,
			c.Score,
			c.Datetime,
			c.Region)
		if err != nil {
			log.Printf("error %s on job %s\n", err, c.JobID)
		} else {
			rowCount, err := result.RowsAffected()
			if err != nil {
				log.Printf("completed job %s\n", c.JobID)
				log.Println(err)
			} else {
				log.Printf("completed job %s rows affected %d\n", c.JobID, rowCount)
			}
		}

		// TODO: update the job status in Redis
	}
}

func main() {
	commands := make(chan *submitScoreCommand)

	var forever chan struct{}

	go getCommands(commands)
	go handleCommands(commands)

	<-forever
}
