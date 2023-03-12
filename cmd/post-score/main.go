package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type submitScoreRequest struct {
	Score    uint64 `json:"score"`
	Name     string `json:"name"`
	Datetime string `json:"datetime"`
}

type submitScoreResponse struct {
	JobID string `json:"job_id"`
}

type submitScoreCommand struct {
	Score    uint64 `json:"score"`
	Name     string `json:"name"`
	Datetime string `json:"datetime"`
	JobID    string `json:"job_id"`
}

func createCommand(request *submitScoreRequest, jobID uuid.UUID) *submitScoreCommand {
	// TODO: validation here
	return &submitScoreCommand{
		Score:    request.Score,
		Name:     request.Name,
		Datetime: request.Datetime,
		JobID:    jobID.String(),
	}
}

func sendSubmitScore(command *submitScoreCommand) error {
	jsonBody, err := json.Marshal(command)
	if err != nil {
		return err
	}

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
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        jsonBody,
		})
	return err
}

func postScore(ctx *gin.Context) {
	var request submitScoreRequest
	if err := ctx.BindJSON(&request); err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jobID := uuid.New()

	if err := sendSubmitScore(createCommand(&request, jobID)); err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := submitScoreResponse{JobID: jobID.String()}
	ctx.JSON(http.StatusAccepted, response)
}

func main() {
	router := gin.Default()
	router.POST("/", postScore)
	router.Run("localhost:3010")
}
