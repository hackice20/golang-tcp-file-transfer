package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/user/go-tcp-ftp/common"
)

type Client struct {
	serverAddr string
	key        []byte
}

func NewClient(serverAddr string, key []byte) *Client {
	return &Client{
		serverAddr: serverAddr,
		key:        key,
	}
}

func (c *Client) SendFile(filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Connect to server
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send filename
	filename := filepath.Base(filePath)
	filenameLen := uint32(len(filename))
	if err := binary.Write(conn, binary.BigEndian, filenameLen); err != nil {
		return fmt.Errorf("failed to send filename length: %v", err)
	}
	if _, err := conn.Write([]byte(filename)); err != nil {
		return fmt.Errorf("failed to send filename: %v", err)
	}

	// Read and encrypt file data
	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}

		encrypted, err := common.Encrypt(buffer[:n], c.key)
		if err != nil {
			return fmt.Errorf("error encrypting data: %v", err)
		}

		if _, err := conn.Write(encrypted); err != nil {
			return fmt.Errorf("error sending data: %v", err)
		}
	}

	fmt.Printf("File sent successfully: %s\n", filePath)
	return nil
}
