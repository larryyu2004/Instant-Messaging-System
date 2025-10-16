package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// Create a user api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// Launch the goroutine that listens to the current user channel message
	go user.ListenMessage()

	return user
}

// User Online
func (user *User) Online() {
	// User login, add user into the OnlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// Boardcast user login message
	user.server.Boardcast(user, "login")
}

// User Offline
func (user *User) Offline() {
	// User logout, remove user from the OnlineMap
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// Boardcast user logout message
	user.server.Boardcast(user, "logout")
}

// Send message to the client of current user
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// Processing Message for users
func (user *User) ProcessingMessage(msg string) {
	if msg == "who" {
		// Search current online users
		user.server.mapLock.Lock()
		for _, userOnline := range user.server.OnlineMap {
			onlineMsg := "[" + userOnline.Addr + "]" + userOnline.Name + ":" + "Online...\n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else {
		user.server.Boardcast(user, msg)
	}

}

// Listen current user goroutine method
// Once get the message, send to the client
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
