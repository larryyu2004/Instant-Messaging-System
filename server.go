package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// Create a server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (server *Server) Handler(conn net.Conn) {
	// Connect with current business
	fmt.Println("Connect Successfully")
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
