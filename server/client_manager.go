package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// Client blabla
type Client struct {
	socket net.Conn
	data   chan []byte
}

// ClientManager blabla
type ClientManager struct {
	gameManager *GameManager
	clients     map[*Client]bool
	broadcast   chan []byte
	register    chan *Client
	unregister  chan *Client
}

func (clientManager *ClientManager) start() {
	for {
		select {
		case connection := <-clientManager.register:
			clientManager.clients[connection] = true
			fmt.Println("Novo cliente")
		case connection := <-clientManager.unregister:
			if _, ok := clientManager.clients[connection]; ok {
				close(connection.data)
				delete(clientManager.clients, connection)
				fmt.Println("Cliente vazou")
			}
		case message := <-clientManager.broadcast:
			for connection := range clientManager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(clientManager.clients, connection)
				}
			}
		}
	}
}

func (clientManager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)

		if err != nil {
			clientManager.unregister <- client
			client.socket.Close()
			break
		}

		if length > 0 {

			commands := strings.Split(string(bytes.Trim(message, "\x00")), "::")

			for index := 0; index < len(commands); index++ {
				commands[index] = strings.TrimSpace(commands[index])
			}

			if len(commands) > 0 {
				switch commands[0] {
				case "get-game-info":
					{

						clientManager.gameManager.getGameInfo(client)
					}

				case "set-name":
					{
						clientManager.gameManager.lobbyManager.setName(client, commands[1])
					}

				case "set-response":
					{
						response := commands[1]
						tip := commands[2]

						clientManager.gameManager.matchManager.setMasterResponse(client, response, tip)
					}
				}

			}

			// get-game-status

			// m1 := Msg1{Msg{"test_cmd"}, "777"}
			// buffer := m1.encode()

			// clientManager.broadcast <- message
			// client.data <- []byte("galo cego")
			// client.data <- buffer.Bytes()
		}
	}
}

func (clientManager *ClientManager) send(client *Client) {
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
