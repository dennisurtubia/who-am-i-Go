package main

import (
	"log"
	// "math/r	and"
	// "strconv"
	"strings"
	"time"
)

const MatchTime = time.Second * 120
const MasterTime = time.Second * 5
const AnswerTime = time.Second * 20

type MatchManager struct {
	gameManager *GameManager

	players     []Player
	roundPlayer *Player
	master      *Player
	response    string
	tip         string

	masterChan chan bool
	answerChan chan bool
	matchChan chan bool

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
			ticker.Stop()
			log.Println("Mestre da partida: " + matchManager.master.name)
			return
		}

	}
}

// setMasterResponse Mestre enviou resposta e dica
func (matchManager *MatchManager) setMasterResponse(client *Client, response string, tip string) {
	// todo: verificar se quem fez essa chamada é o master

	// if matchManager.gameManager.status != WaitingForMaster {
	// 	return
	// }

	log.Println("Mestre enviou resposta e dica")

	matchManager.response = response
	matchManager.tip = tip
	matchManager.gameManager.status = Game
	matchManager.masterChan <- true

}

func (matchManager *MatchManager) playerAnswer() {
	matchManager.answerChan <- true
}

func (matchManager *MatchManager) matchLoop() {

	index := 0
	log.Println("Tentando jogador ", index)

	matchManager.gameManager.clientManager.broadcast <- []byte("1234")
	
	// ticker := time.NewTicker(AnswerTime)
	index = index + 1

	// escolhe jogador e envia request

	// for {
	// 	select {
	// 	case <-matchManager.answerChan:

	// 		matchManager.gameManager.clientManager.broadcast <- []byte("round_player::")
	// 		index = index + 1

	// 		if index > len(matchManager.players) {
	// 		}

	// 		ticker = time.NewTicker(AnswerTime)
			
	// 	case <-ticker.C:
	// 		// escolhe jogador e envia request 
	// 		log.Println("Tentando jogador ", index)
	// 		matchManager.gameManager.clientManager.broadcast <- []byte("round_player::")
	// 		index = index + 1

	// 		// if index > len(matchManager.players) {
	// 		// 	log.Println("")
	// 		// }
	// 	}
	// }




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
	// matchManager.gameManager.clientManager.broadcast <- []byte("game-start\n")
	// matchManager.gameManager.clientManager.broadcast <- []byte("game-start::" + matchManager.master.name + "::" + matchManager.tip + "::" + strconv.FormatInt(matchManager.finishTime.UTC().UnixNano(), 10))
	matchManager.gameManager.clientManager.broadcast <- []byte("pqqqqqz\n")

	log.Println("aqio")
	matchManager.matchChan = make(chan bool, 1)
	matchManager.matchLoop()
	
	log.Println("aqqquququ")
	select {
	case <-matchManager.matchChan:
		break
	case <-time.After(MatchTime):
		break
	}

	log.Println("fimzera")

	// time.Sleep(MatchTime)
}

func (matchManager *MatchManager) reset() {
	matchManager.players = nil
	matchManager.roundPlayer = nil
	matchManager.master = nil
	matchManager.response = ""
	matchManager.tip = ""
}
