# TransactionApp

This project consists of implementing an application for storing and querying purchase transactions


## 🚀 Starting

These instructions will allow you to get a copy of the project running on your local machine for development and testing purposes.

### 📋 Requirements

Tools: 

- [Docker](https://docs.docker.com/desktop/)
- [Golang](https://golang.org/doc/install)


## ⚙️ Running the tests

To run the API unit tests, use the command `make mock` and then the `make test`.

- `make mock`: creates interface implementations, with the aim of performing dependency injection to execute unit tests.
- `make test`: run the unit tests.


## 📦 Development

There are two basic commands for running the project:

- `make run-services`: starts the containers for the resources used by the service.
- `make run`: wrapper to command `go run main.go`.

## 📋 Documentation

To acess the aplication documentation run the service and then access the following url:
- `http://0.0.0.0:5055/swagger/`

