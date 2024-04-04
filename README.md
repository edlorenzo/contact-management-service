# Enhanced Contact Management Service

> Contact management system that allows users to add, view, edit, and delete contact information, enriched with data from an external API.

## Prerequisites

Download and install the following:

- Go 1.21 - https://golang.org/
- GORM - https://gorm.io/
- go-migrate - [go-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## Setup

Go recommends having a single workspace directory for all your Go projects. Before cloning, [set up your Go workspace first](https://golang.org/doc/code.html).

Create the project directory inside the workspace:

```zsh
$ git clone git@github.com:go-ed/contact-management-service.git
```

Declare environment variables

- Refer to the [.env.sample](https://github.com/edlorenzo/contact-management-service/blob/main/.env.example).
- Change the values accordingly
- Rename the .env.sample to .env

```
APP_HOST=127.0.0.1
APP_PORT=8000
APP_NAME=contact-management-service
APP_ENVIRONMENT=dev
APP_READ_TIMEOUT=60
APP_WRITE_TIMEOUT=60

PROCESS=SERVER

DB_HOST=localhost
DB_DATABASE=assessment
DB_USERNAME=assessment
DB_PASSWORD=P@ssw0rd
DB_PORT=5433

MIGRATION_PATH=D:/_0.0.1_Source_Code/_Golang/@my-repo/0403/contact-management-service/cmd/migrations

USER_API_URL=https://jsonplaceholder.typicode.com/users
```

- Note: make sure the `MIGRATION_PATH` is updated.

To start the server locally:

```zsh
$ cd assessments
$ go mod tidy
$ go run cmd/main.go or make run 
```

Start the server via docker-compose:

```
docker-compose up -d --force-recreate
```

## Healthcheck

- Access the healthcheck endpoint (http://localhost:8003/health/liveness | http://localhost:8003/health/readiness) to see if the service is up.

## Postman Collection

- Import the [collection](https://www.postman.com/dark-desert-384866/workspace/my-public-workspace/collection/2409862-b7599ead-f5d7-4367-a18c-d86d6983d2f5) to Postman to see the sample payload and response.
