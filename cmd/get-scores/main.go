package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type topScore struct {
	Rank     uint64 `json:"rank"`
	Score    uint64 `json:"score"`
	Name     string `json:"name"`
	Datetime string `json:"datetime"`
}

// test data
var scores = []topScore{
	{Rank: 1, Score: 1234560, Name: "Jason", Datetime: "2023-03-10 1500"},
	{Rank: 2, Score: 970000, Name: "Ulysses", Datetime: "2023-02-10 1500"},
	{Rank: 3, Score: 870000, Name: "Steinbeck", Datetime: "2023-01-10 1500"},
}

func getScores(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, scores)
}

func main() {
	router := gin.Default()
	router.GET("/scores", getScores)
	router.Run("localhost:3000")
}
