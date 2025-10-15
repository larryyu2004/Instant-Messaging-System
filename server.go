package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// Online User List
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	// Message Boardcast Channel
	Message chan string
}

// Create a server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

// Goroutine for listening the boardcast from the message
// Once get new message, send to all user
func (server *Server) listenMessager () {
	for {
		msg := <- server.Message

		// Send the message to all the online users
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) Boardcast (user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// Connect with current business
	// fmt.Println("Connect Successfully")

	user := NewUser(conn)

	// User login, add user into the OnlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// Boardcast user login message
	server.Boardcast(user, "login")

	// Current handler block
	select {

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
