package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func postScore(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusAccepted, nil)
}

func main() {
	router := gin.Default()
	router.POST("/score", postScore)
	router.Run("localhost:3010")
}
