# Payment Gateway Service

The Payment Gateway Service processes deposit and withdrawal transactions by dynamically selecting third-party gateways
based on country configurations. It supports callback handling for transaction status updates and publishes
transaction events to Kafka.

---

## Features

- **Dynamic Gateway Selection**: Gateways are selected based on country-specific priorities and health checks.
- **Fault Tolerance**: Retry mechanisms and circuit breakers ensure system reliability.
- **Asynchronous Processing**: Publishes transaction events to Kafka for downstream systems.
- **Support for Multiple Data Formats**: JSON and SOAP for request/response payloads.

---

## Setup Instructions

### Prerequisites

To set up the Payment Gateway Service, ensure the following dependencies are installed:

- **Go** (version 1.18+)
- **PostgreSQL** (or any supported database)
- **Kafka** (for transaction event publishing)
- **Docker** (optional, for containerized setup)

### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-repo/payment-gateway-service.git
   cd payment-gateway-service
   ```

2. **Install dependencies**:
   To install the required dependencies for the project, run:
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:
   Replace the environment variables in `.env` file accordingly.

4. **Run database migrations**:
   To set up the database schema, run:
   ```bash
   go run cmd/migrate.go
   ```

5. **Start the service**:
   To start the payment gateway service, use:
   ```bash
   go run cmd/main.go
   ```
6. **Run Unit Tests**:
   ```bash
   go test -v ./...

7. **Test the endpoints**:
   Use tools like `Postman` or `curl` to interact with the API. For example:
    - To test the `/deposit` endpoint:
      ```bash
      curl -X POST http://localhost:8090/deposit \
      -H "Content-Type: application/json" \
      -d '{"amount":100.0,"user_id":1,"country_id":1,"currency":"USD"}'
      ```
    - To test the `/withdrawal` endpoint:
      ```bash
      curl -X POST http://localhost:8090/withdrawal \
      -H "Content-Type: application/json" \
      -d '{"amount":100.0,"user_id":1,"country_id":1,"currency":"USD"}'
      ```

---

## Architectural Decisions

### Key Components

1. **Transaction Service**:
   Handles deposit/withdrawal transactions, selects gateways, and processes payments.

2. **Gateway Selector**:
   Dynamically selects gateways based on:
    - Country-based priority configuration.
    - Gateway health status.

3. **Fault Tolerance**:
    - Retry mechanisms for transient failures.
    - Circuit breakers for persistent issues.

4. **Asynchronous Processing**:
   Publishes transaction events to Kafka for downstream systems.

---

## Gateway Configurations

### Country-Based Gateways

Gateways are configured with priority settings. For example:

```plaintext
Country: US
Priority: Gateway 1 > Gateway 2 > Gateway 3
```

### Health Checks

Each gateway has a health-check endpoint to ensure availability. The system dynamically selects only healthy gateways
during transaction processing.

### Country-Based Gateway Selection

The `SelectGateway` function dynamically selects a gateway by:

1. Fetching gateways for the given country.
2. Sorting gateways by priority.
3. Returning the first healthy gateway.

---

## Folder Structure

```plaintext
payment-gateway-service/
├── cmd/               # Application entry points
├── db/                # Database operations
├── internal/          # Internal services and models
│   ├── api/           # API handlers
│   ├── kafka/         # Kafka producers
│   ├── models/        # Request/response and database models
│   ├── services/      # Core business logic
│   ├── util/          # Utility functions
├── config/            # Configuration files
├── docs/              # API documentation (OpenAPI)
└── tests/             # Unit and integration tests
```

---

## API Endpoints

### POST `/deposit`

- **Description**: Processes a deposit transaction.
- **Request Body**:
  ```json
  {
    "amount": 100.0,
    "user_id": 1,
    "country_id": 840,
    "currency": "USD"
  }
  ```
- **Response**:
  ```json
  {
    "statusCode": 200,
    "message": "Transaction processed successfully",
    "data": {
      "transaction_id": 12345,
      "status": "completed"
    }
  }
  ```

---

### POST `/withdrawal`

- **Description**: Processes a withdrawal transaction.
- **Request Body**:
  ```json
  {
    "amount": 100.0,
    "user_id": 1,
    "country_id": 840,
    "currency": "USD"
  }
  ```
- **Response**:
  ```json
  {
    "statusCode": 200,
    "message": "Transaction processed successfully",
    "data": {
      "transaction_id": 12345,
      "status": "completed"
    }
  }
  ```

---

