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
	status        GameStatus
	lobbyPlayers  []Player
	inGameClients map[*Player]bool

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

		//todo: broadcast
	} else {
		msg += "already_used"
	}

	client.data <- []byte(msg)
}

func (gManager *GameManager) start() {
	fmt.Println("[GameManager] start")

	gManager.status = WaitingForClients
	gManager.waitingFinish = time.Now().Add(time.Minute * 2)

}
