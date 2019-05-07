package main

import (
	"log"
	"strconv"
)

// const LobbyTime = time.Minute * 2
// const GameTime = time.Minute * 8
// const LobbyTime = time.Second * 5

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
	score         int
}

// GameManager a
type GameManager struct {
	clientManager *ClientManager
	lobbyManager  *LobbyManager
	matchManager  *MatchManager

	status GameStatus
}

func (gameManager *GameManager) getPlayerByName(name string) *Player {
	for _, player := range append(gameManager.lobbyManager.players, gameManager.matchManager.players...) {
		if player.name == name {
			return &player
		}
	}
	return nil
}

func (gameManager *GameManager) getClientByName(name string) *Client {
	for _, player := range append(gameManager.lobbyManager.players, gameManager.matchManager.players...) {
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

	gameManager.clientManager.send(client, msg)
	// client.data <- []byte(msg)
	// client.socket.Write([]byte(msg + "\n"))
}

// func (gameManager *GameManager) waitPlayerAnswer() {
// 	gameManager.clientManager.broadcast <- []byte("round_player::" + gameManager.roundPlayer.name)
// }

// func (gameManager *GameManager) gameLoop() {
// 	gameManager.sortPlayers()
// 	for index := 0; index < len(gameManager.inGamePlayers); index++ {

// 		gameManager.roundPlayer = gameManager.inGamePlayers[index]
// 		gameManager.waitPlayerAnswer()
// 	}
// }

func (gameManager *GameManager) start() {

	log.SetPrefix("GameManager ")
	log.Println("Start")

	gameManager.lobbyManager = &LobbyManager{gameManager: gameManager, players: make([]Player, 0)}
	gameManager.matchManager = &MatchManager{gameManager: gameManager, players: make([]Player, 0)}

	for {

		gameManager.lobbyManager.start()
		// copy(gameManager.matchManager.players, gameManager.lobbyManager.players[:])
		gameManager.matchManager.players = gameManager.lobbyManager.players
		gameManager.matchManager.start()

		log.Println("Jogo terminou. Reiniciando...")
		// gameManager.lobbyManager.reset()
		// gameManager.matchManager.reset()
	}
}
