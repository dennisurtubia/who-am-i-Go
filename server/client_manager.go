package main

import (
	"bufio"
	"log"
	"net"
	"strings"
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
		log.Println("[Mensagem] ", message)

		commands := strings.Split(message, "::")

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
					response := commands[2]
					tip := commands[1]

					clientManager.gameManager.matchManager.setMasterResponse(client, response, tip)
				}
			case "player-question":
				{
					log.Println("playerQuestion")
					question := commands[1]
					clientManager.gameManager.matchManager.playerQuestion(question)
				}

			case "master-response":
				{
					response := commands[1]
					clientManager.gameManager.matchManager.masterResponse(response)
				}

			case "player-response":
				{
					response := commands[1]
					clientManager.gameManager.matchManager.playerResponse(response)
				}
			}

		}
	}

	log.Println("Cliente saiu")

	delete(clientManager.clients, client)
}

func (clientManager *ClientManager) broadcast(message string) {
	for client := range clientManager.clients {
		clientManager.send(client, message)
	}
}

func (clientManager *ClientManager) send(client *Client, message string) {
	client.socket.Write([]byte(message + "\n"))
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
