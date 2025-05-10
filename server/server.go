package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/user/go-tcp-ftp/common"
)

type Server struct {
	port     string
	listener net.Listener
	key      []byte
}

func NewServer(port string) (*Server, error) {
	key, err := common.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %v", err)
	}

	return &Server{
		port: port,
		key:  key,
	}, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.listener = listener

	fmt.Printf("Server started on port %s\n", s.port)
	fmt.Printf("Encryption key: %x\n", s.key)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read filename length
	var filenameLen uint32
	if err := binary.Read(conn, binary.BigEndian, &filenameLen); err != nil {
		fmt.Printf("Error reading filename length: %v\n", err)
		return
	}

	// Read filename
	filename := make([]byte, filenameLen)
	if _, err := io.ReadFull(conn, filename); err != nil {
		fmt.Printf("Error reading filename: %v\n", err)
		return
	}

	// Create output file
	outputPath := filepath.Join("received", string(filename))
	if err := os.MkdirAll("received", 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Read and decrypt file data
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading data: %v\n", err)
			return
		}

		decrypted, err := common.Decrypt(buffer[:n], s.key)
		if err != nil {
			fmt.Printf("Error decrypting data: %v\n", err)
			return
		}

		if _, err := file.Write(decrypted); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}

	fmt.Printf("File received and saved: %s\n", outputPath)
}

func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
