package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// Online User List
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// Message Boardcast Channel
	Message chan string
}

// Create a server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// Goroutine for listening the boardcast from the message
// Once get new message, send to all user
func (server *Server) listenMessager() {
	for {
		msg := <-server.Message

		// Send the message to all the online users
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) Boardcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// Connect with current business
	// fmt.Println("Connect Successfully")

	user := NewUser(conn, server)

	user.Online()
	// Listen if user is active
	isLive := make(chan bool)

	// Accept users' messages from client
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// Extract user message, remove '\n' in the end
			msg := string(buf[:n-1])

			// User handle the message
			user.ProcessingMessage(msg)

			// Any message from user make user's state to be active 
			isLive <- true
		}
	}()

	// Current handler block
	for {
		select {
		case <- isLive:
			// Current user is active, reset the timer
			// No need to do anything, just for activating the select statement, reset the following timer

			// 10 seconds later
		case <- time.After(time.Second * 60):
			// Timed Out
			// Kick current user out

			user.SendMsg("You are kicked out")
			// Destory user resource
			close(user.C)

			// Destory connection
			conn.Close()

			return
		}
	}
	

}

// Launch the server interface
func (server *Server) LaunchServer() {
	// Socket Listen
	// Sprintf = String Print format
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	defer listener.Close()

	// Launch goroutine for listening the message
	go server.listenMessager()

	for {
		// Accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// Do Handler
		go server.Handler(conn)
	}
}
