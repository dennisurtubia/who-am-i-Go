package main

import (
	"bytes"
	"fmt"
	"net"
)

// Client blabla
type Client struct {
	socket net.Conn
	data   chan []byte
}

// ClientManager blabla
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Novo cliente")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("Cliente vazou")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)

		// buffer := bytes.NewBuffer(message)
		// decoder := gob.NewDecoder(buffer)

		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}

		if length > 0 {

			fmt.Print("Received: ", string(bytes.Trim(message, "\x00")))

			m1 := Msg1{Msg{"test_cmd"}, "777"}
			buffer := m1.encode()

			// manager.broadcast <- message
			// client.data <- []byte("galo cego")
			client.data <- buffer.Bytes()
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()

	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}
