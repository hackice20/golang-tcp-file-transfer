package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/user/go-tcp-ftp/client"
	"github.com/user/go-tcp-ftp/server"
)

func main() {
	// Define command-line flags
	mode := flag.String("mode", "", "Mode: 'server' or 'client'")
	port := flag.String("port", "8080", "Port to listen on (server) or connect to (client)")
	file := flag.String("file", "", "File to send (client mode only)")
	key := flag.String("key", "", "Encryption key in hex format (client mode only)")
	flag.Parse()

	switch *mode {
	case "server":
		runServer(*port)
	case "client":
		if *file == "" {
			fmt.Println("Error: -file flag is required in client mode")
			os.Exit(1)
		}
		if *key == "" {
			fmt.Println("Error: -key flag is required in client mode")
			os.Exit(1)
		}
		runClient(*port, *file, *key)
	default:
		fmt.Println("Error: -mode flag must be either 'server' or 'client'")
		os.Exit(1)
	}
}

func runServer(port string) {
	srv, err := server.NewServer(port)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting server...")
	if err := srv.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func runClient(port, filePath, keyHex string) {
	// Decode the encryption key
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		fmt.Printf("Error decoding key: %v\n", err)
		os.Exit(1)
	}

	// Create and run client
	cl := client.NewClient("localhost:"+port, key)
	if err := cl.SendFile(filePath); err != nil {
		fmt.Printf("Error sending file: %v\n", err)
		os.Exit(1)
	}
}
