package main

import (
	"fmt"
	"strconv"
	"time"
)

type GameStatus int

const (
	WaitingForClients GameStatus = 0
	InGame            GameStatus = 1
)

type Player struct {
	client *Client
	name   string
}

// GameManager a
type GameManager struct {
	cManager ClientManager

	status        GameStatus
	lobbyPlayers  []Player
	inGamePlayers []Player

	waitingFinish time.Time
	gameFinish    time.Time
}

func (gManager *GameManager) getGameInfo(client *Client) {
	msg := "get-game-info::"

	fmt.Println(gManager.lobbyPlayers, len(gManager.lobbyPlayers))

	if gManager.status == WaitingForClients {
		msg += "waiting::" + strconv.Itoa(len(gManager.lobbyPlayers)) + "::" + strconv.FormatInt(gManager.waitingFinish.UTC().UnixNano(), 10)
	} else {
		msg += "ingame::" + strconv.FormatInt(gManager.gameFinish.UTC().UnixNano(), 10)
	}

	client.data <- []byte(msg)
}

func (gManager *GameManager) setName(client *Client, name string) {
	msg := "set-name::"

	isNameFree := true
	for index := 0; index < len(gManager.lobbyPlayers); index++ {
		if gManager.lobbyPlayers[index].name == name {
			isNameFree = false
			break
		}
	}

	if isNameFree {
		gManager.lobbyPlayers = append(gManager.lobbyPlayers, Player{client: client, name: name})
		msg += "player_added"

		gManager.broadcastNewPlayer(name)
	} else {
		msg += "already_used"
	}

	client.data <- []byte(msg)
}

func (gManager *GameManager) broadcastNewPlayer(name string) {
	gManager.cManager.broadcast <- []byte("new-player::" + name)
}

func (gManager *GameManager) start() {
	fmt.Println("[GameManager] start")

	for {
		gManager.status = WaitingForClients
		gManager.waitingFinish = time.Now().Add(time.Minute * 2)
		gManager.lobbyPlayers = make([]Player, 0)

		time.Sleep(time.Minute * 2)

		gManager.status = InGame
		gManager.gameFinish = time.Now().Add(time.Minute * 8)
		gManager.inGamePlayers = make([]Player, len(gManager.lobbyPlayers))
		copy(gManager.inGamePlayers, gManager.lobbyPlayers)
		gManager.lobbyPlayers = nil

		gManager.cManager.broadcast <- []byte("game-start::" + strconv.FormatInt(gManager.gameFinish.UTC().UnixNano(), 10))

	}

}
