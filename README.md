# QR Payment Service – DANA Integration

This service facilitates QRIS CPM payments via DANA, encompassing functionalities for payment initiation, status querying, and order cancellation. It adheres to the specifications outlined in the DANA Merchant Portal API documentation.
![image](https://github.com/user-attachments/assets/c0c3b383-c78d-4948-be8e-a7a4343961bc)
![image](https://github.com/user-attachments/assets/969cbc28-03d8-495e-b4ae-9269eff59bd7)

## Features

* **Payment Initiation**: Processes QRIS CPM payments through DANA.
* **Status Querying**: Checks the status of pending transactions.
* **Order Cancellation**: Cancels orders in cases of payment failures or timeouts.

## Architecture Overview

The service is structured into distinct components:

* **Services**: Business logic handling payment, query, and cancellation operations.
* **Repositories**: Interfaces for database interactions.
* **Models**: Structs representing data models.
* **Requests/Responses**: Structs for API request and response payloads.
* **Utilities**: Helper functions for tasks like signature generation.
* **Handlers**: Handling the Http Request from client.
* **Routes**: Managing api routes.

## API Endpoints

### 1. Initiate Payment

* **Endpoint**: `/api/payment`

* **Method**: `POST`

* **Request Body**:

  ```json
  {
    "trxId": "string",
    "qrContent": "string",
    "amount": 10000,
    "productCode": "string",
    "branchId": "string",
    "counterId": "string",
    "caCode": "string"
  }
  ```

* **Response**:

  ```json
  {
    "trxId": "string",
    "trxConfirm": "string",
    "responseCode": "00",
    "responseMessage": "SUCCESS",
    "paidAt": "2025-06-01T17:16:02+07:00"
  }
  ```

### 2. Query Payment Status

* **Endpoint**: `/api/payment/status`

* **Method**: `POST`

* **Request Body**:

  ```json
  {
    "trxId": "string",
    "productCode": "string"
  }
  ```

* **Response**:

  ```json
  {
    "trxId": "TP8120250530102450001",
    "responseCode": "00",
    "responseMessage": "SUCCESS",
    "latestTransactionStatus": "00",
    "paidTime": "2025-06-01T17:16:02+07:00",
    "amount":"100000"
  }
  ```

### 3. Cancel Order

* **Endpoint**: `/api/payment/cancel`

* **Method**: `POST`

* **Request Body**:

  ```json
  {
    "trxId": "string",
    "productCode": "string"
  }
  ```

* **Response**:

  ```json
  {
    "responseCode": "00",
    "responseMessage": "Order canceled successfully",
    "trxId":"TP8120250530102450001",
    "cancelTime":"2020-12-21T17:07:25+07:00",
    "transactionDate":"2020-12-21T17:55:11+07:00"
  }
  ```

## Configuration

Ensure the following configurations are set, typically via environment variables or configuration files:

* **Client Secret**: Used for generating HMAC signatures.
* **Partner ID**: Provided by DANA.
* **External ID**: Unique identifier for external systems.
* **Channel ID**: Identifier for the communication channel.
* **Authorization Token**: Bearer token for API authentication.
* **API URLs**: Endpoints for payment, status, and cancellation APIs.

## Logging

The service logs each step of the transaction process, including:

* Incoming requests.
* Outgoing requests to DANA.
* Responses from DANA.
* Errors and exceptions.

Logs are stored in the `tracelog` repository for auditing and debugging purposes.

## Error Handling

The service returns appropriate HTTP status codes and error messages for various failure scenarios, such as:

* Invalid request payloads.
* Failed communication with DANA APIs.
* Internal server errors.

Ensure to handle these errors gracefully in your client applications.

## Dependencies

* **Go Modules**: Ensure all required Go modules are installed.
* **External Libraries**: Utilize standard libraries for HTTP requests, JSON handling, and cryptographic operations.

## Running the Service

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/adrianuscharlie/Golang-Payment-QR.git
   ```

2. **Navigate to the Project Directory**:

   ```bash
   cd Golang-Payment-QR
   ```

3. **Set Up Configuration**:

   Configure environment variables or configuration files with the necessary credentials and API endpoints.

4. **Run the Service**:

   ```bash
   go run ./cmd
   ```

## Testing

Implement unit and integration tests to ensure the reliability of each component. Use Go's testing framework for writing and executing tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to customize this README further to match your project's specifics and organizational standards.
