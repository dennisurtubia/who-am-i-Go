package main

import (
	"bufio"
	"net"
	"strings"
	"log"
)

// Client blabla
type Client struct {
	socket net.Conn
}

// ClientManager blabla
type ClientManager struct {
	gameManager *GameManager
	clients     map[*Client]bool
}

func (clientManager *ClientManager) handleClient(client *Client) {
	defer client.socket.Close()

	log.Println("Novo cliente")

	clientManager.clients[client] = true

	scanner := bufio.NewScanner(client.socket)

	for scanner.Scan() {
		message := scanner.Text()
		
		

		commands := strings.Split(message, "::")

		log.Println("mensagem do zapp ", commands)


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
	}

	delete(clientManager.clients, client)
}

func (clientManager *ClientManager) broadcast(message string) {
	for client := range(clientManager.clients) {
		client.socket.Write([]byte(message))
	}
}

func (clientManager *ClientManager) send(client *Client, message string) {
	client.socket.Write([]byte(message))
}

func (clientManager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)

		if err != nil {
			log.Panicln("errooo ", err)
			// clientManager.clientExit(client)
			client.socket.Close()
			break
		}

		if length > 0 {

		}
	}
}

