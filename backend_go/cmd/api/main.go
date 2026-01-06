package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type MessageStats struct {
	Ack     int `json:"ack"`
	Publish int `json:"publish"`
}

type Queue struct {
	Name         string       `json:"name"`
	Messages     int          `json:"messages"`
	Consumers    int          `json:"consumers"`
	MessageStats MessageStats `json:"message_stats"`
}

// Determine RabbitMQ Management API URL based on environment
func getManagementAPI() string {
	amqpURL := os.Getenv("RABBITMQ_URL")

	// Check if running inside Docker container
	if strings.Contains(amqpURL, "rabbitmq") {
		return "http://s_rabbitmq:15672/api/queues"
	}

	// Default to localhost for local development
	return "http://localhost:15672/api/queues"
}

func main() {
	r := gin.Default()

	// Enable CORS for mobile access
	r.Use(cors.Default())

	r.GET("/queues", func(c *gin.Context) {
		url := getManagementAPI()

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req.SetBasicAuth("guest", "guest")

		resp, err := client.Do(req)
		if err != nil {
			log.Println("RabbitMQ Connection Error:", err)
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

		// Handle empty response safely
		if allQueues == nil {
			allQueues = []Queue{}
		}
		c.JSON(http.StatusOK, allQueues)
	})

	log.Println("API Server running on :8080")
	r.Run(":8080")
}
