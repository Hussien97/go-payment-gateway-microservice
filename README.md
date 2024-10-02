
# Payment Gateway Microservice

## Documentation

### Summary
This project implements a payment gateway microservice that supports deposit and withdrawal transactions. It uses a microservices architecture with Docker for containerization and includes integrations with Kafka, Redis, and PostgreSQL. The solution efficiently handles transaction requests through a microservices architecture that emphasizes resilience, scalability, and security. The design allows for easy integration of new gateways and data formats, ensuring a future-proof system that can adapt to evolving business needs.

### Components
- **API Gateway(assumption)**: Acts as the entry point for client requests, directing them to the Payment Gateway Microservice.
- **Payment Gateway Microservice**: Responsible for processing transaction requests. It validates incoming data, saves transaction details concurrently to PostgreSQL and Redis for quick access, and publishes messages to Kafka for asynchronous processing by relevant gateways.
- **Kafka**: Serves as the message broker, enabling decoupled communication between the microservice and the gateways. Each transaction is published to specific topics based on its data format (e.g., JSON, SOAP) ensuring that only the related gateways will consume it.
- **Redis Cache**: Provides a caching layer for frequently accessed data, improving performance by reducing the number of direct database queries.
- **PostgreSQL**: A relational database used to persist transaction details and statuses, ensuring data integrity and reliability.
- **Gateways**: External services that listen to Kafka topics corresponding to their respective data formats. They process incoming transactions and send callbacks to update transaction statuses.

### Request Flow
1. A client submits a transaction request via the API Gateway.
2. The API Gateway forwards the request to the Payment Gateway Microservice.
3. The microservice validates the request and saves the transaction to PostgreSQL and Redis.
4. The transaction data is masked and published to the appropriate Kafka topic.
5. Relevant gateways consume the transaction messages from Kafka, process them, and send callbacks to update the transaction status.

### Implementation Logic
- **Fault Tolerance**: The microservice uses circuit breakers to handle failures when communicating with Kafka, ensuring that failed requests are marked appropriately in PostgreSQL to prevent duplication of transactions.

### Kafka Topic Strategy
- The system uses distinct Kafka topics for different data formats, ensuring relevant gateways manage specific messages and enabling the addition of new topics as required.

### Security Measures
- **Data Masking**: Sensitive information is masked before transmission to Kafka, ensuring transaction details remain protected.
- **Digital Signatures**: Digital signatures can be applied to callback requests using HMAC, ensuring that only authorized gateways can invoke callback endpoints.

### Ease of Adding New Gateways
- The middleware responsible for validating data formats can be extended to include new formats, facilitating rapid integration of new gateway types.


## How to Run Locally
1. Ensure you have Docker installed.
2. pull the code from the repository
2. Run `docker-compose up --build` to start all services.
3. For running tests only, use `docker-compose up --build test` to log test results (it will take 10 seconds to run the tests as it will wait for kafka, redis, db to run first).

## Documentation Links
- [Technical Documentation](https://drive.google.com/file/d/1tUuOjMrFeTRT5lhQ62b3KWtuNuwYpj1t/view?usp=sharing)
- [API Documentation](https://app.swaggerhub.com/apis/HUSSIENCIS/Payment-Gateway-Microservice/1.0.0)
- [Repository](https://github.com/Hussien97/go-payment-gateway-microservice)

### Conclusion
The Payment Gateway Microservice is designed to provide robust transaction processing while ensuring high availability, fault tolerance, and security. The design strategies employed facilitate easy scalability and flexibility for future enhancements, making the system well-equipped to handle evolving business requirements.
