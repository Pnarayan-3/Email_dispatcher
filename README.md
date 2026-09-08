# Email Dispatcher

A lightweight concurrent email dispatcher built with Go. It reads recipients from a CSV file, renders a reusable HTML email template for each recipient, and dispatches emails concurrently using multiple worker goroutines.

The project uses **Go channels**, **goroutines**, **WaitGroup**, **CSV processing**, **HTML templates**, and **SMTP**. For local development, it is designed to work with Mailpit as a local SMTP server and email inbox.

## Features

- 📧 Reads recipient information from a CSV file
- ⚡ Sends emails concurrently using multiple worker goroutines
- 🔄 Uses a Go channel to distribute recipients to workers
- 🧵 Uses `sync.WaitGroup` to wait for all workers to finish
- 📝 Generates email content from an HTML template
- 🧪 Supports local email testing with Mailpit
- 🛠️ Simple project structure with no external Go dependencies

## Project Structure

```text
email-dispatcher/
├── consumer.go      # Worker implementation and SMTP email sending
├── producer.go      # Reads recipients from the CSV file
├── main.go          # Application entry point, workers, channel and template execution
├── email.tmpl       # HTML email template
├── emails.csv       # Recipient input data
├── info.md          # Local Mailpit setup command
├── go.mod           # Go module definition
└── .github/
    └── workflows/
        └── go-ci.yml
```

## Prerequisites

- Go 1.27+
- Docker (recommended for local SMTP testing)

Verify Go:

```bash
go version
```

## Run Mailpit Locally

The application currently sends SMTP traffic to:

```text
localhost:1025
```

Start Mailpit:

```bash
docker run -d \
  --restart unless-stopped \
  --name=mailpit \
  -p 8025:8025 \
  -p 1025:1025 \
  axllent/mailpit
```

Mailpit provides:

- SMTP server: `localhost:1025`
- Web UI: `http://localhost:8025`

Open the Mailpit UI after sending emails to inspect the captured messages.

## Input CSV

The application expects `emails.csv` to contain a header followed by recipient records:

```csv
Name,Email
Alice,alice@example.com
Bob,bob@example.com
Charlie,charlie@example.com
```

The first row is skipped because it is treated as the CSV header.

> Use test addresses when working with Mailpit. The messages are captured locally and are not delivered to real recipients.

## Email Template

The email body is loaded from `email.tmpl`.

The template can use fields from the `Recipient` struct:

```html
<h1>Hello {{.Name}}</h1>
<p>Welcome to our email campaign.</p>
```

Available fields:

```text
.Name
.Email
```

## Run the Application

From the project root:

```bash
go run .
```

You should see worker output similar to:

```text
worker 1: Sending email to alice@example.com
worker 2: Sending email to bob@example.com
worker 1: Sent email to alice@example.com
worker 2: Sent email to bob@example.com
```

Then open Mailpit:

```text
http://localhost:8025
```

## How Concurrency Works

The application creates a channel:

```go
recipientChannel := make(chan Recipient)
```

The producer reads recipients from the CSV file and sends them to the channel:

```text
CSV → Producer → Channel
```

Five workers consume recipients concurrently:

```go
workerCount := 5
```

Each worker waits for a recipient and sends an email through SMTP.

```text
                 Channel
              /    |    |    |    \
             W1    W2   W3   W4    W5
              \    |    |    |    /
                    SMTP
```

A `sync.WaitGroup` ensures the main goroutine waits until all workers finish processing.

## Configuration

The current implementation has the SMTP server and port defined in `consumer.go`:

```go
smtpHost := "localhost"
smtpPort := "1025"
```

For production use, these values should be moved to environment variables or another configuration mechanism.

For example:

```text
SMTP_HOST
SMTP_PORT
SMTP_USERNAME
SMTP_PASSWORD
```

Avoid committing real SMTP credentials to GitHub.

## Error Handling

The application currently handles:

- CSV file opening errors
- CSV parsing errors
- Email template parsing/execution errors
- SMTP sending errors

The worker currently terminates the process when `smtp.SendMail` fails because it uses `log.Fatal`. For a production-grade dispatcher, consider returning the error to a centralized error handler instead.

## CI Pipeline

The repository includes a GitHub Actions workflow that can:

1. Check out the repository
2. Set up Go
3. Verify formatting
4. Run `go vet`
5. Run tests
6. Build the application

Workflow location:

```text
.github/workflows/go-ci.yml
```

## Learning Goals

This project is useful for practicing:

- Go concurrency
- Goroutines
- Channels
- Worker-pool pattern
- `sync.WaitGroup`
- CSV parsing
- HTML templating
- SMTP communication
- Error handling
- GitHub Actions CI/CD

