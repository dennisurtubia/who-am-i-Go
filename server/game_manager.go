package main

import (
	"log"
	"strconv"
	"time"
)

// const LobbyTime = time.Minute * 2
// const GameTime = time.Minute * 8
// const LobbyTime = time.Second * 5
const MasterTime = time.Second * 10

type GameStatus int

const (
	WaitingForClients GameStatus = 0
	WaitingForMaster  GameStatus = 1
	InGame            GameStatus = 2
	WaitingAnswer     GameStatus = 3
	Lobby                        = 4
	Game              GameStatus = 5
)

type Player struct {
	client *Client
	name   string

	masterAttempt bool
}

// GameManager a
type GameManager struct {
	clientManager *ClientManager
	lobbyManager  *LobbyManager
	matchManager  *MatchManager

	status        GameStatus
	lobbyPlayers  []Player
	inGamePlayers []Player

	waitingFinish time.Time
	gameFinish    time.Time
	playerTimeout time.Time

	roundPlayer Player

	masterName string
	response   string
	tip        string
}

func (gameManager *GameManager) getPlayerByName(name string) *Player {
	for _, player := range append(gameManager.lobbyPlayers, gameManager.inGamePlayers...) {
		if player.name == name {
			return &player
		}
	}
	return nil
}

func (gameManager *GameManager) getClientByName(name string) *Client {
	for _, player := range append(gameManager.lobbyPlayers, gameManager.inGamePlayers...) {
		if player.name == name {
			return player.client
		}
	}
	return nil
}

func (gameManager *GameManager) getGameInfo(client *Client) {
	msg := "get-game-info::"

	if gameManager.status == Lobby {
		msg += "waiting::" + strconv.Itoa(len(gameManager.lobbyManager.players)) + "::" + strconv.FormatInt(gameManager.lobbyManager.waitingTime.UTC().UnixNano(), 10)
	} else {
		msg += "ingame::" + strconv.FormatInt(gameManager.matchManager.finishTime.UTC().UnixNano(), 10)
	}

	client.data <- []byte(msg)
}

func (gameManager *GameManager) waitPlayerAnswer() {
	gameManager.clientManager.broadcast <- []byte("round_player::" + gameManager.roundPlayer.name)
}

func (gameManager *GameManager) gameLoop() {
	gameManager.sortPlayers()
	for index := 0; index < len(gameManager.inGamePlayers); index++ {

		gameManager.roundPlayer = gameManager.inGamePlayers[index]
		gameManager.waitPlayerAnswer()
	}
}

func (gameManager *GameManager) start() {

	log.SetPrefix("GameManager ")
	log.Println("Start")

	lobbyManager := LobbyManager{gameManager: gameManager}
	matchManager := MatchManager{gameManager: gameManager}

	for {

		lobbyManager.start()
		copy(matchManager.players, lobbyManager.players)
		matchManager.start()

		log.Println("Jogo terminou. Reiniciando...")
		lobbyManager.reset()
		matchManager.reset()
	}
}
