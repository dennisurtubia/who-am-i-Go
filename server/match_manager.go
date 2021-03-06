package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	// "math/rand"
	"strconv"
	"strings"
	"time"
)

const MatchTime = time.Second * 120
const MasterTime = time.Second * 20
const QuestionTime = time.Second * 20
const MasterAnswerTime = time.Second * 20

type Highscore struct {
	playerName string
	score      int
}
type MatchManager struct {
	gameManager *GameManager

	players     []Player
	roundPlayer *Player
	master      *Player
	response    string
	tip         string

	masterChan chan bool
	// playerQuestionChan chan bool
	// masterResponseChan chan bool
	// playerResponseChan chan bool
	matchChan chan bool

	playerQuestionChan chan string
	masterResponseChan chan string
	playerResponseChan chan string

	finishTime time.Time

	responseStart time.Time
	responseEnd   time.Time
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
	// ticker := time.NewTicker(MasterTime)

	// for index := 0; index < len(matchManager.players); index++ {
	// 	if matchManager.players[index].masterAttempt == false {
	// 		matchManager.players[index].masterAttempt = true
	// 		matchManager.master = &matchManager.players[index]
	// 		break
	// 	}
	// }

	// if matchManager.master == nil {
	// 	log.Println("Não foi possível escolher um mestre")
	// 	return
	// }

	matchManager.master = &matchManager.players[0]

	log.Println("Mestre da partida: " + matchManager.master.name)

	// matchManager.gameManager.getClientByName(matchManager.master.name).data <- []byte("game-master::" + matchManager.master.name)

	matchManager.gameManager.clientManager.send(matchManager.gameManager.getClientByName(matchManager.master.name), "game-master::"+matchManager.master.name)
	<-matchManager.masterChan
	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		for index := 0; index < len(matchManager.players); index++ {
	// 			if matchManager.players[index].masterAttempt == false {
	// 				matchManager.players[index].masterAttempt = true
	// 				matchManager.master = &matchManager.players[index]
	// 				break
	// 			}
	// 		}

	// 		if matchManager.master == nil {
	// 			log.Println("Não foi possível escolher um mestre")
	// 			return
	// 		}

	// 		log.Println("Tentando mestre: " + matchManager.master.name)

	// 		matchManager.gameManager.getClientByName(matchManager.master.name).socket.Write([]byte("game-master::" + matchManager.master.name))
	// 	case <-matchManager.masterChan:
	// 		ticker.Stop()
	// 		log.Println("Mestre da partida: " + matchManager.master.name)
	// 		return
	// 	}

	// }
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

func (matchManager *MatchManager) playerQuestion(question string) {
	log.Println("jogador perguntou")
	matchManager.playerQuestionChan <- question
}

func (matchManager *MatchManager) masterResponse(response string) {
	log.Println("master respondeu ", response)
	matchManager.masterResponseChan <- response
}

func (matchManager *MatchManager) playerResponse(response string) {
	log.Println("jogador respondeu")
	matchManager.playerResponseChan <- response
}

func (matchManager *MatchManager) selectPlayer(index *int) bool {

	if matchManager.master == &matchManager.players[*index] {
		*index++
	}

	if *index > len(matchManager.players) {
		log.Println("acabou")
		return false
	}

	matchManager.gameManager.clientManager.broadcast("round-player::" + matchManager.players[*index].name)
	(*index)++

	return true
}

func (matchManager *MatchManager) matchLoop() {

	for index := 0; index < len(matchManager.players); index++ {

		matchManager.playerQuestionChan = make(chan string, 1)
		matchManager.masterResponseChan = make(chan string, 1)
		matchManager.playerResponseChan = make(chan string, 1)

		if matchManager.master == &matchManager.players[index] {
			index++
		}

		log.Println("Jogador da vez: ", matchManager.players[index])

		if index >= len(matchManager.players) {
			log.Println("acabou")
			return
		}

		matchManager.gameManager.clientManager.broadcast("round-player::" + matchManager.players[index].name)

		playerQuestion := <-matchManager.playerQuestionChan // timeout

		matchManager.gameManager.clientManager.send(matchManager.gameManager.getClientByName(matchManager.master.name), "player-question::"+matchManager.players[index].name+"::"+playerQuestion)

		masterResponse := <-matchManager.masterResponseChan //timeout

		matchManager.gameManager.clientManager.broadcast("master-response::" + matchManager.players[index].name + "::" + playerQuestion + "::" + masterResponse)

		matchManager.responseStart = time.Now()
		playerResponse := <-matchManager.playerResponseChan
		matchManager.responseEnd = time.Now()

		if playerResponse == matchManager.response {
			score := int(100 * (1 / (matchManager.responseEnd.Sub(matchManager.responseStart).Seconds())))
			// log.Println("Resposta correta. Score: ", strconv.Itoa(score))
			matchManager.players[index].score = score
			// player := matchManager.gameManager.getPlayerByName(matchManager.players[index].name)
			// player.score = score
			matchManager.gameManager.clientManager.broadcast("player-response::" + matchManager.players[index].name + "::true::" + strconv.Itoa(matchManager.players[index].score))
		} else {
			matchManager.gameManager.clientManager.broadcast("player-response::" + matchManager.players[index].name + "::false")
		}
	}
	matchManager.matchChan <- true

}

func (matchManager *MatchManager) start() {
	log.SetPrefix("MatchManager")

	log.Println("Configurando partida")

	// matchManager.players = make([]Player, 0)

	for index := 0; index < len(matchManager.players); index++ {
		matchManager.players[index].score = 0
	}

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
	matchManager.gameManager.clientManager.broadcast("game-init::" + playerNamesJoin)

	// Espera definição do mestre
	matchManager.masterChan = make(chan bool, 1)
	matchManager.waitMasterResponse()

	if matchManager.master == nil {
		return
	}

	matchManager.finishTime = time.Now().Add(MatchTime)

	log.Println("Iniciando partida")
	// matchManager.gameManager.clientManager.broadcast <- []byte("game-start\n")
	matchManager.gameManager.clientManager.broadcast("game-start::" + matchManager.master.name + "::" + matchManager.tip + "::" + strconv.FormatInt(matchManager.finishTime.Unix(), 10))

	matchManager.matchChan = make(chan bool, 1)
	go matchManager.matchLoop()

	select {
	case <-matchManager.matchChan:
		break
	case <-time.After(MatchTime):
		break
	}

	maxScorePlayer := matchManager.players[1]
	for index := 2; index < len(matchManager.players); index++ {
		if matchManager.players[index].score > maxScorePlayer.score {
			maxScorePlayer = matchManager.players[index]
		}
	}

	/*
		retorna nesse formato:

		game-end::jorge::jorge:100,dennis:80,carlos sumare:20
	*/

	if maxScorePlayer.score == 0 {
		matchManager.gameManager.clientManager.broadcast("game-end")

		b, _ := ioutil.ReadFile("./highscores.txt")
		s := strings.TrimSpace(string(b))
		matchManager.gameManager.clientManager.broadcast("highscore::" + s)

	} else {
		winStr := "game-end::" + maxScorePlayer.name + "::"

		tmp := make([]string, 0)
		for index := 1; index < len(matchManager.players); index++ {
			tmp = append(tmp, matchManager.players[index].name+":"+strconv.Itoa((matchManager.players[index].score)))
		}
		winStr = winStr + strings.Join(tmp[:], ",")
		log.Println("Fim da partida: ", winStr)

		b, _ := ioutil.ReadFile("./highscores.txt")
		s := strings.TrimSpace(string(b))

		lines := strings.Split(s, ",")

		hscores := make([]Highscore, 0)

		for _, line := range lines {
			if len(line) > 0 {
				tmp2 := strings.Split(line, ":")
				name := tmp2[0]
				score, _ := strconv.Atoi(tmp2[1])
				hscores = append(hscores, Highscore{playerName: name, score: score})
			}
		}
		for _, player := range matchManager.players[1:] {
			if player.score != 0 {
				hscores = append(hscores, Highscore{playerName: player.name, score: player.score})
			}
		}

		sort.Slice(hscores, func(i, j int) bool {
			return hscores[i].score > hscores[j].score
		})

		if len(hscores) > 10 {
			hscores = hscores[:10]
		}

		out := make([]string, 0)
		for _, hscore := range hscores {
			out = append(out, hscore.playerName+":"+strconv.Itoa(hscore.score))
		}

		fmt.Println(out)

		ioutil.WriteFile("./highscores.txt", []byte(strings.Join(out, ",")), 0777)

		matchManager.gameManager.clientManager.broadcast(winStr)
		matchManager.gameManager.clientManager.broadcast("highscore::" + strings.Join(out, ","))
	}

}

func (matchManager *MatchManager) reset() {
	matchManager.players = nil
	matchManager.roundPlayer = nil
	matchManager.master = nil
	matchManager.response = ""
	matchManager.tip = ""
}
