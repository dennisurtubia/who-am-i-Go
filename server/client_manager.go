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

func (cManager *ClientManager) start() {
	for {
		select {
		case connection := <-cManager.register:
			cManager.clients[connection] = true
			fmt.Println("Novo cliente")
		case connection := <-cManager.unregister:
			if _, ok := cManager.clients[connection]; ok {
				close(connection.data)
				delete(cManager.clients, connection)
				fmt.Println("Cliente vazou")
			}
		case message := <-cManager.broadcast:
			for connection := range cManager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(cManager.clients, connection)
				}
			}
		}
	}
}

func (cManager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)

		if err != nil {
			cManager.unregister <- client
			client.socket.Close()
			break
		}

		if length > 0 {

			// log.p("Received: ", string(bytes.Trim(message, "\x00")))

			commands := strings.Split(string(bytes.Trim(message, "\x00")), "::")

			for index := 0; index < len(commands); index++ {
				commands[index] = strings.TrimSpace(commands[index])
			}

			if len(commands) > 0 {
				switch commands[0] {
				case "get-game-info":
					{

						cManager.gameManager.getGameInfo(client)
					}

				case "set-name":
					{
						cManager.gameManager.setName(client, commands[1])
					}

				case "set-response":
					{
						response := commands[1]
						tip := commands[2]

						cManager.gameManager.setResponse(client, response, tip)
					}
				}

			}

			// get-game-status

			// m1 := Msg1{Msg{"test_cmd"}, "777"}
			// buffer := m1.encode()

			// cManager.broadcast <- message
			// client.data <- []byte("galo cego")
			// client.data <- buffer.Bytes()
		}
	}
}

func (cManager *ClientManager) send(client *Client) {
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
