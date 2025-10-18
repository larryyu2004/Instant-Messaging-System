package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // Current Mode
}

func NewClient(serverIp string, serverPort int) *Client {
	// Create a new client object
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
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

// Deal with responses from the server
func (client *Client) DealResponse() {
	// Once client.conn has data, it will be copied directly to stdout standard output, permanently blocking the listener
	io.Copy(os.Stdout, client.conn)
	/*
		for {
			buf := make()
			client.conn.Read(buf)
			fmt.Println(buf)
		}
	*/
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. Public Chat Mode")
	fmt.Println("2. Private Chat Mode")
	fmt.Println("3. Rename")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Invalid")
		return false
	}
}

// Query online users
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// Private chat mode
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	// Show list of online users first
	client.SelectUsers()

	fmt.Println(">>>> Please enter the username to chat with (type 'exit' to quit):")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>> Please enter your message (type 'exit' to quit):")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// If the message is not empty, send it
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err:", err)
					break
				}
			}

			// Prompt again for the next message
			fmt.Println(">>>> Please enter your message (type 'exit' to quit):")
			fmt.Scanln(&chatMsg)
		}

		// After exiting chat with this user, prompt for another
		client.SelectUsers()
		fmt.Println(">>>> Please enter the username to chat with (type 'exit' to quit):")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	// Prompt the user to enter a message
	var chatMsg string

	fmt.Println(">>>> Please enter your message (type 'exit' to quit):")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// If the message is not empty, send it to the server
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}

		// Continue prompting the user for new input
		fmt.Println(">>>> Please enter your message (type 'exit' to quit):")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>> Please enter new username")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		// Different Cases
		switch client.flag {
		case 1:
			// Public Chat Mode
			client.PublicChat()
			break
		case 2:
			// Private Chat Mode
			client.PrivateChat()
			break
		case 3:
			// Rename
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set server ip address")
	flag.IntVar(&serverPort, "port", 8888, "Set server port")
}

func main() {
	// Parse Commends
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>> Failed to connect the server...")
	}

	// Start a goroutine to deal with server message
	go client.DealResponse()

	fmt.Println(">>>> Successfully connected to the server....")
	client.Run()
}
