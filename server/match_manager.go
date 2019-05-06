package main

import (
	"log"
	// "math/r	and"
	"strconv"
	"strings"
	"time"
)

const MatchTime = time.Minute * 10
const MasterTime = time.Second * 10


type MatchManager struct {
	gameManager *GameManager

	players     []Player
	roundPlayer *Player
	master      *Player
	response    string
	tip         string

	masterChan chan bool

	finishTime time.Time
}

func (gameManager *GameManager) initGame() {

	// log.Println("Initing game...")
	// gameManager.status = WaitingForMaster
	// gameManager.inGamePlayers = make([]Player, len(gameManager.lobbyPlayers))
	// copy(gameManager.inGamePlayers, gameManager.lobbyPlayers)
	// gameManager.lobbyPlayers = nil

}

// func (gameManager *GameManager) sortPlayers() {

// 	tmp := make([]Player, 0)
// 	for _, player := range gameManager.inGamePlayers {
// 		if player.name != gameManager.masterName {
// 			tmp = append(tmp, player)
// 		}
// 	}

// 	rand.Shuffle(len(tmp), func(i, j int) {
// 		tmp[i], tmp[j] = tmp[j], tmp[i]
// 	})
// }

// waitMasterResponse  Escolhe mestre e aguarda resposta e dica.
func (matchManager *MatchManager) waitMasterResponse() {

	
	ticker := time.NewTicker(MasterTime)

	
	for index := 0; index < len(matchManager.players); index++ {
		if matchManager.players[index].masterAttempt == false {
			matchManager.players[index].masterAttempt = true
			matchManager.master = &matchManager.players[index]
			break
		}
	}

	if matchManager.master == nil {
		log.Println("Não foi possível escolher um mestre")
		return
	}

	log.Println("Tentando mestre: " + matchManager.master.name)

	matchManager.gameManager.getClientByName(matchManager.master.name).data <- []byte("game-master::" + matchManager.master.name)

	for {
		select {
		case <-ticker.C:
			for index := 0; index < len(matchManager.players); index++ {
				if matchManager.players[index].masterAttempt == false {
					matchManager.players[index].masterAttempt = true
					matchManager.master = &matchManager.players[index]
					break
				}
			}
	
			if matchManager.master == nil {
				log.Println("Não foi possível escolher um mestre")
				return
			}
	
			log.Println("Tentando mestre: " + matchManager.master.name)
	
			matchManager.gameManager.getClientByName(matchManager.master.name).data <- []byte("game-master::" + matchManager.master.name)
		case <-matchManager.masterChan:
			log.Println("pronto")
			ticker.Stop()
			log.Println("Mestre da partida: " + matchManager.master.name)
			return
		}

	}
}

// setMasterResponse Mestre enviou resposta e dica
func (matchManager *MatchManager) setMasterResponse(client *Client, response string, tip string) {
	// todo: verificar se quem fez essa chamada é o master

	if matchManager.gameManager.status != WaitingForMaster {
		return
	}

	log.Println("Mestre enviou resposta e dica")

	matchManager.response = response
	matchManager.tip = tip
	matchManager.gameManager.status = Game
	matchManager.masterChan <- true

}

func (matchManager *MatchManager) start() {
	log.SetPrefix("MatchManager")

	log.Println("Configurando partida")

	// matchManager.players = make([]Player, 0)
	matchManager.finishTime = time.Now().Add(MatchTime)
	matchManager.gameManager.status = Game

	if len(matchManager.players) < 2 {
		log.Println("Jogadores insuficientes")
		return
	}

	// Envia nome dos jogadores separados por vírgula [game-init]
	matchManager.gameManager.status = WaitingForMaster
	playersNames := make([]string, 0)

	for _, player := range matchManager.players {
		playersNames = append(playersNames, player.name)
	}

	playerNamesJoin := strings.Join(playersNames[:], ",")
	matchManager.gameManager.clientManager.broadcast <- []byte("game-init::" + playerNamesJoin)

	// Espera definição do mestre
	matchManager.masterChan = make(chan bool, 1)
	matchManager.waitMasterResponse()

	if matchManager.master == nil {
		return
	}

	matchManager.finishTime = time.Now().Add(MatchTime)

	log.Println("Iniciando partida")
	matchManager.gameManager.clientManager.broadcast <- []byte("game-start::" + matchManager.master.name + "::" + matchManager.tip + "::" + strconv.FormatInt(matchManager.finishTime.UTC().UnixNano(), 10))

	time.Sleep(MatchTime)
}

func (matchManager *MatchManager) reset() {
	matchManager.players = nil
	matchManager.roundPlayer = nil
	matchManager.master = nil
	matchManager.response = ""
	matchManager.tip = ""
}
