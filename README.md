#  RabbitMQ Queue Manager

This project is a **Full-Stack** application designed to monitor RabbitMQ queues, message counts, and consumer statuses in real-time.

> **Status:** This project is currently **under active development (WIP)**. Architecture and features are subject to change as the learning process continues.

##  Architecture & Tech Stack

The project follows a **Monorepo** structure and utilizes the following technologies:

* **Message Broker:** RabbitMQ (Running on Docker)
* **Backend (BFF):** Go (Golang) & Gin Framework
    * *Role:* Consumes the RabbitMQ Management API, filters raw data, and serves optimized JSON to the mobile client.
* **Mobile:** React Native (Expo)
    * *Role:* Displays queue lists and visualizes critical states (e.g., DLQ alerts).
