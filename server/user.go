package main

import (
	"net"
	"strings"
)

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
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// Message Format: rename|xxxxx
		newName := strings.Split(msg, "|")[1]

		// Check if the newName exists
		_, ok := user.server.OnlineMap[newName]
		if ok {
			// exist
			user.SendMsg("New user name has been used")
		} else {
			// inexist
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMsg("You have updated user name" + user.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// Message Format: to|<username>|<message>

		// Get username
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMsg("Incorrect message format, please use format: to|<username>|<message>")
			return
		}

		// Based on the username, get user object
		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok {
			user.SendMsg("The user not exists")
			return
		}

		// Get message, send the message based on the user object
		content := strings.Split(msg, "|")[2]

		if content == "" {
			user.SendMsg("No content, please try again!")
			return
		}
		remoteUser.SendMsg(user.Name + " says to you: " + content)

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
