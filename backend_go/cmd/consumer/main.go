package main

import (
	"encoding/json"
	"log"
	"strings"

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

	// 1. Declare the DLQ (Dead Letter Queue) first
	_, err = ch.QueueDeclare(
		"emails_dlq", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare DLQ")

	// 2. Define arguments for the Main Queue to link it to DLQ
	args := amqp.Table{
		"x-dead-letter-exchange":    "",           // Default exchange
		"x-dead-letter-routing-key": "emails_dlq", // Send rejected msgs here
	}

	// 3. Declare the Main Queue with DLQ arguments
	q, err := ch.QueueDeclare(
		"emails", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		args,     // apply DLQ settings
	)
	failOnError(err, "Failed to declare main queue")

	// 4. Start Consuming (Auto-Ack must be FALSE)
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (FALSE: we will ack/nack manually)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var email Email
			json.Unmarshal(d.Body, &email)

			// 5. Validation Logic
			if strings.Contains(email.Address, "@") {
				// Valid Email -> Process and ACK
				log.Printf("Email Sent: %s", email.Address)
				d.Ack(false)
			} else {
				// Invalid Email -> NACK (Send to DLQ)
				// requeue: false is crucial to move it to DLQ
				log.Printf("Invalid Email: %s -> Moving to DLQ", email.Address)
				d.Nack(false, false)
			}
		}
	}()

	log.Printf(" [*] Waiting for emails. Press CTRL+C to exit")
	<-forever
}
