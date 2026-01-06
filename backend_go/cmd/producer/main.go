package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Helper: Get RabbitMQ URL from env or default to localhost
func getRabbitMQURL() string {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return "amqp://guest:guest@localhost:5672/"
	}
	return url
}

func main() {
	// 1. Connect
	conn, err := amqp.Dial(getRabbitMQURL())
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	qName := "emails"

	// 2. Send 5 test messages
	for i := 1; i <= 5; i++ {
		body := fmt.Sprintf("Email #%d", i)

		// Every 3rd message is "bad" to test DLQ
		if i%3 == 0 {
			body = "error_mail"
		}

		err = ch.Publish(
			"",    // Exchange
			qName, // Routing Key
			false, // Mandatory
			false, // Immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		if err != nil {
			log.Println("Publish error:", err)
		} else {
			log.Printf("Sent: %s", body)
		}
	}
}
