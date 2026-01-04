package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type MessageStats struct {
	Ack     int `json:"ack"`     // Total messages successfully processed
	Publish int `json:"publish"` // Total messages sent to this queue
}

type Queue struct {
	Name         string       `json:"name"`
	Messages     int          `json:"messages"` // Messages waiting to be processed
	Consumers    int          `json:"consumers"`
	MessageStats MessageStats `json:"message_stats"` // This field contains the history
}

func main() {
	r := gin.Default()

	// Enable CORS for mobile access
	r.Use(cors.Default())

	r.GET("/queues", func(c *gin.Context) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:15672/api/queues", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req.SetBasicAuth("guest", "guest")

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "RabbitMQ unreachable"})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var allQueues []Queue
		if err := json.Unmarshal(body, &allQueues); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON parse error"})
			return
		}

		// Directly return all queues found in RabbitMQ
		if allQueues == nil {
			allQueues = []Queue{}
		}
		c.JSON(http.StatusOK, allQueues)
	})

	log.Println("🔌 API Server running on :8080")
	r.Run(":8080")
}
