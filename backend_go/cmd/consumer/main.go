package main

import (
	"log"
	"os"
	"time"

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
	// 1. Connect to RabbitMQ
	conn, err := amqp.Dial(getRabbitMQURL())
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// 2. Declare Dead Letter Exchange (DLX)
	dlxName := "dlx_exchange"
	err = ch.ExchangeDeclare(dlxName, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare DLX:", err)
	}

	// 3. Declare DLQ (Dead Letter Queue)
	dlqName := "emails_dlq"
	_, err = ch.QueueDeclare(dlqName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare DLQ:", err)
	}

	// 4. Bind DLQ to DLX
	err = ch.QueueBind(dlqName, "emails", dlxName, false, nil)
	if err != nil {
		log.Fatal("Failed to bind DLQ:", err)
	}

	// 5. Declare Main Queue with DLX configuration
	qName := "emails"
	args := amqp.Table{
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": "emails",
	}
	q, err := ch.QueueDeclare(qName, true, false, false, false, args)
	if err != nil {
		log.Fatal("Failed to declare main queue:", err)
	}

	// 6. Start Consuming
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil) // Auto-Ack: false
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	log.Println("Waiting for messages...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			body := string(d.Body)
			log.Printf("Received: %s", body)

			// Simulate processing time
			time.Sleep(1 * time.Second)

			// Logic: Reject messages containing "error" to test DLQ
			if body == "error_mail" {
				log.Println("Error detected! Moving to DLQ...")
				d.Nack(false, false) // Requeue: false (Send to DLQ)
			} else {
				log.Println("Processed successfully")
				d.Ack(false)
			}
		}
	}()

	<-forever
}
