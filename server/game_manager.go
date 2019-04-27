package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// const LobbyTime = time.Minute * 2
// const GameTime = time.Minute * 8
const LobbyTime = time.Second * 10
const GameTime = time.Second * 15
const MasterTime = time.Second * 10

type GameStatus int

const (
	WaitingForClients GameStatus = 0
	WaitingForMaster  GameStatus = 1
	InGame            GameStatus = 2
)

type Player struct {
	client *Client
	name   string
}

// GameManager a
type GameManager struct {
	cManager *ClientManager

	status        GameStatus
	lobbyPlayers  []Player
	inGamePlayers []Player

	waitingFinish time.Time
	gameFinish    time.Time
	playerTimeout time.Time

	masterName string
	response   string
	tip        string
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

func (gManager *GameManager) setResponse(client *Client, response string, tip string) {
	// todo: verificar se quem fez essa chamada é o master

	if gManager.status == WaitingForMaster {

		gManager.response = response
		gManager.tip = tip
		gManager.status = InGame

		// broadcast
	}
}

func (gManager *GameManager) broadcastNewPlayer(name string) {
	gManager.cManager.broadcast <- []byte("new-player::" + name)
}

func (gManager *GameManager) initLobby() {
	log.Println("Waiting players...")

	gManager.status = WaitingForClients
	gManager.waitingFinish = time.Now().Add(LobbyTime)
	gManager.lobbyPlayers = make([]Player, 0)

}

func (gManager *GameManager) initGame() {

	log.Println("Initing game...")
	gManager.status = WaitingForMaster
	gManager.inGamePlayers = make([]Player, len(gManager.lobbyPlayers))
	copy(gManager.inGamePlayers, gManager.lobbyPlayers)
	gManager.lobbyPlayers = nil

	gManager.waitMaster()

	gManager.gameFinish = time.Now().Add(GameTime)

	log.Println("Game started...")
	gManager.cManager.broadcast <- []byte("game-start::" + strconv.FormatInt(gManager.gameFinish.UTC().UnixNano(), 10))

}

func (gManager *GameManager) waitMaster() {
	log.Println("Waiting master...")
	for gManager.status != InGame {
		//todo: novo master deve ser diferente do anterior no caso do timeout estourar
		//todo: ver se não vão ocorrer problemas de sincronização
		masterIndex := rand.Intn(len(gManager.inGamePlayers))
		gManager.masterName = gManager.inGamePlayers[masterIndex].name

		gManager.cManager.broadcast <- []byte("game-master::" + gManager.masterName)

		time.Sleep(MasterTime)
	}
	log.Println("Master: " + gManager.masterName)
}

func (gManager *GameManager) start() {

	log.SetPrefix("GameManager ")
	log.Println("Start")

	for {

		gManager.initLobby()
		time.Sleep(LobbyTime)

		gManager.initGame()
		time.Sleep(GameTime)

		log.Println("Game finished...")
	}
}
