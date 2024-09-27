# REST API

This is a restful server written in go, connected with CockroachDB and dockerized. 
This README will guide you through the steps to set up and test the API.

## Getting Started

### Prerequisites

- Go 1.22.3
- Docker
- Docker Compose

### Starting the API

To start the API, you can use Docker Compose. Run the following command in the terminal:

```bash
docker-compose up --build
```
As soon as the server starts running, dummy data will be added to the database through POST requests.

### Running Tests

To run the tests, execute the following command the terminal:

```bash
go test -v
```

**Note**: The tests will modify the data in the database.

### Manual Testing

The server is running on :8888, so you can manually send http requests **or** 
you can send requests running the client.go. 

**Warning**: To run client.go you need to pass an argument
- c to create a product (POST)
- r to read products from the database (GET)
- u to update a product (UPDATE)
- d to delete a product (DELETE)
- db to add data.json to the database

```bash
cd client
go run client.go r
```

**Note**: Make sure to modify the arguments passed to the functions inside client.go
