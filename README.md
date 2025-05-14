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
- `GET /example/get/users/all` - Returns a message indicating where all users would be returned
- `GET /example/get/user?email=user@example.com` - Returns user details for the specified email
- `POST /example/create/user` - Creates a new user
- `PUT /example/update/user` - Updates an existing user

### POST /example/create/user
Request body:
```json
{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com"
}
```

### PUT /example/update/user
Request body (at least firstName OR lastName must be provided):
```json
{
    "email": "john.doe@example.com",
    "firstName": "Johnny",
    "lastName": "Doe"
}
```
