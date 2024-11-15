# Music Transfer API

## Project Overview

This project is a backend API built with Go, using Gorm and Gin, connected to a PostgreSQL database. The API features RabbitMQ integration for handling background tasks such as playlist and song transfers. The purpose of this project is to showcase an efficient architecture for processing large tasks asynchronously while maintaining robust API functionality.

## Features

- **RESTful API** built using Fiber for fast and lightweight web routing.
- **Database Management** with Gorm for easy handling of database operations.
- **RabbitMQ Integration** for asynchronous task processing.
- **Authentication** using JWT for secure API access.
- CRUD operations for managing **users**, **playlists**, and **songs**.
- **Concurrency Control** using RabbitMQ consumers for parallel processing.

## Technologies Used

- [Go](https://golang.org/) (Language)
- [Gin](https://gin-gonic.com/) (Web Framework)
- [Gorm](https://gorm.io/) (ORM for Database Interaction)
- [PostgreSQL](https://www.postgresql.org/) (Database)
- [RabbitMQ](https://www.rabbitmq.com/) (Message Broker)
- [AMQP Library](https://github.com/streadway/amqp) (RabbitMQ Go client)
- [JWT](https://jwt.io/) (Authentication)

## Getting Started

### Prerequisites

Before running this project, ensure that you have the following installed on your machine:

- [Go](https://golang.org/doc/install) (v1.18 or higher recommended)
- [PostgreSQL](https://www.postgresql.org/download/)
- [RabbitMQ](https://www.rabbitmq.com/download.html)
- [Docker](https://www.docker.com/products/docker-desktop) (optional, for containerized setup)
- [Make](https://www.gnu.org/software/make/manual/make.html)

### Installation

1. **Clone the Repository**
    ```bash
    git clone https://github.com/Ocheezyy/music-transfer-api.git
    cd music-transfer-api
    ```

2. **Install Dependencies**
    ```bash
    make tidy-all
    ```

3. **Set Up Environment Variables**
    For each service create a `.env` file based on the `example.env` file

4. **Run database migrations**
    ```bash
    make api-migrate
    ```

5. **Run the application**
    ```bash
    make api-run
    ```
    NOTE: Will be adding testing docker-compose functionality soon

### Project Structure
```
.
├── api/
    ├── main.go                # Entry point of the application
    ├── controllers/           # Request handlers for different endpoints
    ├── helpers/               # Helper functions
    ├── initializers/          # Initializer functions to expose singleton variables
    ├── middlewares/           # Middlewares for routes
    ├── models/                # Database models and types
    ├── producers/             # RabbitMQ producers
    ├── test/                  # Test helper functions
    ├── config/                # Configuration settings
    └── README.md              # Project documentation
├── consumer/
    ├── main.go                # Entry point of the application
    ├── initializers/          # Initializer functions to expose singleton variables
    ├── rabbitmq/              # The rabbitmq consumer client and functions related to it
    └── types                  # Types
└── Makefile                   # Makefile for bash commands to ease monorepo
```

### Future Improvements
- [ ] Finish adding unit tests for handlers and producer.
- [ ] Implement rate limiting for API endpoints.
- [ ] Enhance error handling and logging mechanisms.
- [ ] Finish setup of rabbitmq consumer client with retries, and threading
- [ ] Create docker-compose for ability to run the entire backend easily


### Contributing
Contributions are welcome! Please follow the standard GitHub workflow:
1. Fork the project.
2. Create your feature branch (git checkout -b feature/new-feature).
3. Commit your changes (git commit -m 'Add new feature').
4. Push to the branch (git push origin feature/new-feature).
5. Open a pull request.

### License
This project is licensed under the MIT License - see the LICENSE file for details.