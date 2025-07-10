package main

// This is a simple client for the go-chat-server.

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn     net.Conn
	name     string
	messages chan string
}

type Server struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) Start() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()

			fmt.Printf("Client %s connected\n", client.name)
			go func() {
				s.broadcast <- fmt.Sprintf("%s has joined the chat", client.name)
			}()

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.messages)
				s.mutex.Unlock()

				fmt.Printf("Client %s disconnected\n", client.name)
				go func() {
					s.broadcast <- fmt.Sprintf("%s has left the chat", client.name)
				}()
			} else {
				s.mutex.Unlock()
			}
		case message := <-s.broadcast:
			s.mutex.RLock()
			for client := range s.clients {
				select {
				case client.messages <- message:
				default:
					s.mutex.RUnlock()
					s.unregister <- client
					s.mutex.RLock()
				}
			}
			s.mutex.RUnlock()
		}
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Enter your name: "))
	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" {
		name = "Anonymous"
	}

	client := &Client{
		conn:     conn,
		name:     name,
		messages: make(chan string, 100),
	}

	s.register <- client

	go s.writeMessages(client)

	client.messages <- "Welcome to the chat!"
	client.messages <- "Type '/help' for commands or just start chatting!"

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		conn.Write([]byte("\033[1A\033[2K"))

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		if strings.HasPrefix(message, "/") {
			s.handleCommand(client, message)
			continue
		}
		formattedMessage := fmt.Sprintf("%s: %s", client.name, message)
		s.broadcast <- formattedMessage

	}
	s.unregister <- client
}

func (s *Server) handleCommand(client *Client, message string) {
	switch message {
	case "/help":
		client.messages <- "Available commands:"
		client.messages <- "/help - Show this help message"
		client.messages <- "/users - List online users"
		client.messages <- "/quit - Leave the chat"
	case "/users":
		s.mutex.RLock()
		userList := "Online users: "
		users := make([]string, 0, len(s.clients))
		for client := range s.clients {
			users = append(users, client.name)
		}
		s.mutex.RUnlock()

		userList += strings.Join(users, ", ")
		client.messages <- userList
	case "/quit":
		client.messages <- "Bye!"
		client.conn.Close()

	default:
		client.messages <- "Unknown command try /help"
	}
}

func (s *Server) writeMessages(client *Client) {
	defer func() {
		client.conn.Close()
		s.unregister <- client
	}()
	for message := range client.messages {
		_, err := client.conn.Write([]byte(message + "\n"))
		if err != nil {
			break
		}
	}
}

func main() {
	server := NewServer()
	go server.Start()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server started on port 8080\n")
	fmt.Println("Clients can connect using: telnet localhost 8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go server.handleClient(conn)
	}
}
