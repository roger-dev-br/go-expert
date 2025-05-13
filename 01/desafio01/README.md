# Go Currency Quote

This project consists of a simple client-server application written in Go that retrieves the current dollar exchange rate and stores it in a text file.

## Project Structure

```
go-currency-quote
├── client.go        # Client implementation that requests the dollar exchange rate
├── server.go        # Server implementation that provides the exchange rate
├── cotacao.txt      # File where the current dollar exchange rate is stored
├── go.mod           # Module definition for the Go project
└── README.md        # Documentation for the project
```

## How to Run

### Prerequisites

- Go installed on your machine
- SQLite database (for server-side data persistence)

### Running the Server

1. Navigate to the project directory:
   ```
   cd go-currency-quote
   ```

2. Run the server:
   ```
   go run server.go
   ```

The server will start listening on port 8080.

### Running the Client

1. In a new terminal window, navigate to the project directory:
   ```
   cd go-currency-quote
   ```

2. Run the client:
   ```
   go run client.go
   ```

The client will make a request to the server to get the current dollar exchange rate and save it to `cotacao.txt` in the format: `Dólar: {valor}`.

## Error Handling

Both the client and server implement context timeouts to handle potential delays:

- The client has a timeout of 300ms for receiving the response from the server.
- The server has a timeout of 200ms for fetching the exchange rate from the external API and 10ms for persisting the data into the SQLite database.

Errors will be logged if any of these timeouts are exceeded.