# Agni - Self-Hostable Notification Engine

Agni is a notification engine designed to be self-hostable, lightweight, and scalable. It provides a robust solution for managing notifications across multiple channels, built with modern technologies like Go, Fiber, Redis, and Docker.

## 🚀 Project Overview
Agni aims to simplify notification management by offering a modular and extensible architecture. It is designed to handle high-throughput notification delivery with features like retry queues, rate limiting, and message deduplication.

## 🔧 Tech Stack
- **Language**: Go
- **Web Framework**: Fiber
- **Database**: Redis (optimized for queues, ephemeral data, and TTLs)
- **Containerization**: Docker
- **Communication**: HTTP/gRPC SDK
- **Optional Features**: Webhooks, Rate Limiting, Retry Queues, Message Deduplication

## ✅ Planned Features
### 1. API Server
- REST endpoints for:
    - Sending notification requests
    - Querying notification status
    - Health checks
    - Authentication (API keys, JWT)

### 2. Message Processor
- Processes notification tasks from Redis queues
- Sends messages via multiple channels (e.g., email, SMS, push notifications)
- Implements:
    - Retry logic
    - Rate limiting
    - Failure logging

### 3. Redis Integration
- **List/Stream**: For managing message queues
- **Hash**: For storing metadata
- **Sorted Set**: For scheduling future notifications
- **TTL**: For expiring processed notifications

### 4. Agni SDK
- Client SDK for application servers to:
    - Send notifications to the Agni server
    - Check delivery status
- Includes built-in support for retries, timeouts, and fallbacks

### 5. Notification Sender Modules
- Modular design with pluggable adapters for:
    - Email (SMTP/Sendgrid)
    - SMS (Twilio/etc.)
    - Web Push
    - Slack, Discord, etc.

## 🐳 Docker Support
- The Docker container will include:
    - Agni server binary
    - Configuration files (`config.yaml` or `.env`)
    - Health check scripts

## 🌟 Why Agni?
Agni is designed to be:
- **Fast**: Built with Fiber for high performance
- **Scalable**: Redis-backed architecture for efficient data handling
- **Extensible**: Modular design for easy integration with new channels
- **Self-Hostable**: Full control over your notification infrastructure

Stay tuned for updates as we build Agni into a powerful notification engine for your applications!