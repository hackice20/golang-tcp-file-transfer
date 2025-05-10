# Go TCP File Transfer

A secure file transfer tool built in Go that uses TCP with AES-GCM encryption.

## Features

- Secure file transfer using TCP
- AES-GCM encryption for data security
- Simple command-line interface
- Support for large files
- Automatic key generation for server

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone this repository
3. Run `go mod download` to install dependencies

## Usage

### Starting the Server

```bash
go run main.go -mode server -port 8080
```

The server will start and display an encryption key. You'll need this key for the client.

### Sending a File

```bash
go run main.go -mode client -port 8080 -file path/to/your/file -key <encryption-key>
```

Replace `<encryption-key>` with the key displayed by the server.

## Example

1. Start the server:
```bash
go run main.go -mode server -port 8080
```

2. In another terminal, send a file:
```bash
go run main.go -mode client -port 8080 -file test.txt -key <key-from-server>
```

The file will be received in the `received` directory on the server side.

## Security

- Uses AES-GCM encryption for secure file transfer
- Each server instance generates a unique encryption key
- The key must be shared securely between server and client

## License

MIT 