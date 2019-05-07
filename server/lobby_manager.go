package main

import (
	"log"
	"time"
)

// LobbyTime Tempo de espera no lobby
const LobbyTime = time.Second * 20

// LobbyManager Gerenciador do lobby
type LobbyManager struct {
	gameManager *GameManager

	players     []Player
	waitingTime time.Time
}

// setName Cliente envia o nome enquanto está no lobby
func (lobbyManager *LobbyManager) setName(client *Client, name string) {

	if lobbyManager.gameManager.status != Lobby {
		return
	}

	msg := "set-name::"

	isNameFree := true
	for index := 0; index < len(lobbyManager.players); index++ {
		if lobbyManager.players[index].name == name {
			isNameFree = false
			break
		}
	}

	if isNameFree {
		lobbyManager.players = append(lobbyManager.players, Player{client: client, name: name, masterAttempt: false})
		msg += "player_added"

	} else {
		msg += "already_used"
	}

	// client.data <- []byte(msg)
	lobbyManager.gameManager.clientManager.send(client, msg)

}

func (lobbyManager *LobbyManager) start() {

	log.SetPrefix("LobbyManager")

	log.Println("Esperando jogadores")

	// lobbyManager.players = make([]Player, 0)
	lobbyManager.waitingTime = time.Now().Add(LobbyTime)
	lobbyManager.gameManager.status = Lobby

	time.Sleep(LobbyTime)
}

func (lobbyManager *LobbyManager) reset() {
	lobbyManager.players = nil
}
