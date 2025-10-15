package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// Create a user api
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// Launch the goroutine that listens to the current user channel message
	go user.ListenMessage()

	return user
}

// Listen current user goroutine method
// Once get the message, send to the client
func (user *User) ListenMessage() {
	for {
		msg := <- user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
