package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Email struct {
	Address string `json:"address"`
	Body    string `json:"body"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Test Data: Mixed valid and invalid emails
	emails := []Email{
		{Address: "user1@example.com", Body: "Welcome!"},
		{Address: "invalid-email-address", Body: "Error?"}, // This should go to DLQ
		{Address: "user2@gmail.com", Body: "Invoice"},
		{Address: "spam#nomail", Body: "Spam"}, // This should go to DLQ
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, email := range emails {
		body, _ := json.Marshal(email)

		err = ch.PublishWithContext(ctx,
			"",       // exchange
			"emails", // routing key (must match consumer queue)
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent: %s", email.Address)
	}
}
