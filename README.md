# RabbitMQ Queue Manager

This project is a **Full-Stack** application designed to monitor RabbitMQ queues, message counts, and consumer statuses in real-time. It demonstrates a microservices architecture running entirely on Docker.

> **Status:** This project is currently **under active development (WIP)**. Architecture and features are subject to change as the learning process continues.

## Architecture & Tech Stack

The project follows a **Monorepo** structure and utilizes the following technologies, orchestrated via **Docker Compose**:

* **Infrastructure:** Docker & Docker Compose (Orchestration)
* **Message Broker:** RabbitMQ 4 (Management Plugin enabled)
* **Backend (BFF - `s_api`):** Go (Golang) & Gin Framework
    * *Role:* Consumes the RabbitMQ Management API, filters raw data, and serves optimized JSON to the mobile client.
* **Worker (`s_consumer`):** Go (Golang)
    * *Role:* Background service that continuously processes messages from the queue.
* **Mobile:** React Native (Expo)
    * *Role:* Displays queue lists and visualizes critical states (e.g., DLQ alerts).

## Simulation Scenario: Email Service

To demonstrate real-world usage, this project simulates an **Email Dispatch System**:

1.  **Producer (Traffic Generator):** A script that simulates sending bulk emails. It sends 5 messages at a time.
    * **Success Scenario:** Standard messages like "Email #1" are sent.
    * **Failure Scenario:** To test error handling, every 3rd message is intentionally sent as content "error_mail".
2.  **Consumer (Worker):** Reads messages from the `emails` queue.
    * If the message is valid, it processes it (simulated with a delay) and sends an **ACK** (Acknowledge).
    * If the message is "error_mail", it sends a **NACK** (Negative Acknowledge), causing the message to be routed to the **DLQ (Dead Letter Queue)**.

This setup allows us to monitor both successful processing and error rates via the mobile app.

## Getting Started

### Prerequisites

* Docker Engine (Running)
* Node.js & npm (For the mobile app)

### 1. Run the Backend Infrastructure

The entire backend system (RabbitMQ, API, and Consumer) is Dockerized. To start all services, run the following command in the root directory:
```
docker-compose up --build
```
This command will:

-   Start **RabbitMQ** on ports `5672` & `15672`.
    
-   Start the **Go API** on port `8080`.
    
-   Start the **Consumer** worker (which will immediately start listening to queues).

### 2. Triggering Test Data (Producer)

Since the system runs inside a Docker network, you need to run the producer script inside that network to send messages. Open a new terminal and run:


```
docker-compose run s_api go run cmd/producer/main.go

```

This command spins up a temporary container, executes the Go producer code to send test emails, and then exits. You will see the queues update in the mobile app immediately.