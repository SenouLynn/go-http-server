# Go HTTP Server

A simple HTTP server implementation using Go's standard library.

## Project Setup

To initialize this project from scratch, the following commands were used:

```bash
# Initialize a new git repository
git init

# Create a new Go module
go mod init go-http-server
```

## Running the Server

To run the server:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Project Structure

- `main.go` - Contains the HTTP server implementation
- `go.mod` - Go module definition
- `.gitignore` - Specifies which files Git should ignore

## Available Endpoints

- `GET /` - Returns a welcome message
