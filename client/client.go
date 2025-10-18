package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
}

func NewClient (serverIp string, serverPort int) *Client {
	// Create a new client object
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
	}

	// Connect with server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn
	
	// Return Object
	return client
}

func main () {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>> Failed to connect the server...")
	}

	fmt.Println(">>>> Successfully connected to the server....")
	select{}
}