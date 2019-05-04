package main

import (
	"container/ring"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// const LobbyTime = time.Minute * 2
// const GameTime = time.Minute * 8
const LobbyTime = time.Second * 5
const GameTime = time.Minute * 10
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

	masterAttempt bool
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

	playerOrder ring.Ring

	masterName string
	response   string
	tip        string
}

func (gManager *GameManager) getPlayerByName(name string) *Player {
	for _, player := range append(gManager.lobbyPlayers, gManager.inGamePlayers...) {
		if player.name == name {
			return &player
		}
	}
	return nil
}

func (gManager *GameManager) getClientByName(name string) *Client {
	for _, player := range append(gManager.lobbyPlayers, gManager.inGamePlayers...) {
		if player.name == name {
			return player.client
		}
	}
	return nil
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

	if gManager.status == WaitingForClients {
		msg := "set-name::"

		isNameFree := true
		for index := 0; index < len(gManager.lobbyPlayers); index++ {
			if gManager.lobbyPlayers[index].name == name {
				isNameFree = false
				break
			}
		}

		if isNameFree {
			gManager.lobbyPlayers = append(gManager.lobbyPlayers, Player{client: client, name: name, masterAttempt: false})
			msg += "player_added"

		} else {
			msg += "already_used"
		}

		client.data <- []byte(msg)
	}
}

func (gManager *GameManager) setResponse(client *Client, response string, tip string) {
	// todo: verificar se quem fez essa chamada é o master

	if gManager.status == WaitingForMaster {

		log.Println("Rexpostaaa do mestre")

		gManager.response = response
		gManager.tip = tip
		gManager.status = InGame

		// broadcast
	}
}

func (gManager *GameManager) initLobby() {
	log.Println("Waiting players...")

	gManager.status = WaitingForClients
	gManager.waitingFinish = time.Now().Add(LobbyTime)
	gManager.lobbyPlayers = make([]Player, 0)

}

func (gManager *GameManager) sortPlayers() {

	tmp := make([]Player, 0)
	for _, player := range gManager.inGamePlayers {
		if player.name != gManager.masterName {
			tmp = append(tmp, player)
		}
	}

	rand.Shuffle(len(tmp), func(i, j int) {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})

	fmt.Println(tmp)
}

func (gManager *GameManager) initGame() {

	if len(gManager.lobbyPlayers) < 2 {
		log.Println("Jogadores insuficientes")
		return
	}

	log.Println("Initing game...")
	gManager.status = WaitingForMaster
	gManager.inGamePlayers = make([]Player, len(gManager.lobbyPlayers))
	copy(gManager.inGamePlayers, gManager.lobbyPlayers)
	gManager.lobbyPlayers = nil

	playersNames := make([]string, 0)

	for _, player := range gManager.inGamePlayers {
		playersNames = append(playersNames, player.name)
	}

	playerNamesJoin := strings.Join(playersNames[:], ",")

	gManager.cManager.broadcast <- []byte("game-init::" + playerNamesJoin)

	gManager.waitMaster()

	gManager.gameFinish = time.Now().Add(GameTime)

	log.Println("Game started...")
	gManager.cManager.broadcast <- []byte("game-start::" + gManager.masterName + "::" + gManager.tip + "::" + strconv.FormatInt(gManager.gameFinish.UTC().UnixNano(), 10))

	gManager.sortPlayers()

}

func (gManager *GameManager) waitMaster() {
	log.Println("Waiting master...")

	for gManager.status != InGame {
		//todo: ver se não vão ocorrer problemas de sincronização

		gManager.masterName = ""

		for index := 0; index < len(gManager.inGamePlayers); index++ {
			if gManager.inGamePlayers[index].masterAttempt == false {
				gManager.inGamePlayers[index].masterAttempt = true
				gManager.masterName = gManager.inGamePlayers[index].name
				break
			}
		}

		if gManager.masterName == "" {
			log.Panicln("Sem mestre irmão")
		}

		log.Println("Mestre escolhido: " + gManager.masterName)

		gManager.getClientByName(gManager.masterName).data <- []byte("game-master::" + gManager.masterName)
		// gManager.cManager.broadcast <- []byte("game-master::" + gManager.masterName)

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
