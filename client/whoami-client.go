package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	conn     net.Conn
	name     string
	isMaster bool
}

func (client *Client) send(message string) {
	client.conn.Write([]byte(message + "\n"))
}

func (client *Client) receiveMessages() {
	scanner := bufio.NewScanner(client.conn)
	reader := bufio.NewReader(os.Stdin)

	for scanner.Scan() {
		message := scanner.Text()
		commands := strings.Split(message, "::")

		for index := 0; index < len(commands); index++ {
			commands[index] = strings.TrimSpace(commands[index])
		}

		if len(commands) > 0 {
			switch commands[0] {
			case "get-game-info":
				{
					if commands[1] == "waiting" {

						i, err := strconv.ParseInt(commands[3], 10, 64)

						if err == nil {
							t := time.Unix(i, 0)
							fmt.Println("Atualmente há:", commands[2], "jogadores conectados.")
							fmt.Println("Próxima partida inicia em:", int(t.Sub(time.Now()).Seconds()), "segundos.")

							fmt.Println("Para começar, nos diga: quem é você?")
							fmt.Print("\\> ")
							text, _ := reader.ReadString('\n')
							client.send("set-name::" + text)

						}

					} else if commands[1] == "ingame" {
						i, err := strconv.ParseInt(commands[2], 10, 64)

						if err == nil {
							t := time.Unix(i, 0)
							fmt.Println("Partida em andamento.")
							fmt.Println("Próxima partida inicia em:", int(t.Sub(time.Now()).Seconds()), "segundos.")
						}
					}

				}
			case "set-name":
				{
					if commands[1] == "player_added" {
						fmt.Println("Aguardando partida...")
						client.name = commands[2]

					} else if commands[1] == "already_used" {

						fmt.Println("ERRO: nome atualmente em uso.")
						fmt.Print("Novo nome: ")

						text, _ := reader.ReadString('\n')
						client.send("set-name::" + text)

					}
				}

			case "game-init":
				{
					fmt.Println("----------------")
					fmt.Println("INICIANDO PARTIDA")
					fmt.Println("----------------")
					names := strings.Split(commands[1], ",")

					fmt.Println("Jogadores conectados: " + strconv.Itoa(len(names)))

					for _, name := range names {
						fmt.Printf("\t[%s]\n", name)
					}

					fmt.Println("Aguardando definição do mestre...")
				}

			case "game-master":
				{
					if commands[1] == client.name {
						fmt.Println("------------------------------------------------")
						fmt.Println("VOCÊ FOI ESCOLHIDO COMO MESTRE DA PARTIDA")
						fmt.Println("------------------------------------------------")
						fmt.Print("Informe a dica: ")

						tip, _ := reader.ReadString('\n')

						fmt.Print("Informe a resposta: ")

						response, _ := reader.ReadString('\n')

						client.send("set-response::" + strings.TrimSpace(tip) + "::" + response)

					}
				}

			case "game-start":
				{
					if commands[1] == client.name {
						fmt.Println("----------------------")
						fmt.Println("VOCÊ É O MESTRE")
						fmt.Println("----------------------")

						client.isMaster = true

						i, err := strconv.ParseInt(commands[3], 10, 64)

						if err == nil {
							t := time.Unix(i, 0)
							fmt.Println("Partida termina daqui:", int(t.Sub(time.Now()).Seconds()), "segundos.")
						}
						fmt.Println("Dica: " + commands[2])
						fmt.Println("Aguardando respostas...")
					} else {
						fmt.Println("----------------------")
						fmt.Println("PARTIDA INICIADA")
						fmt.Println("----------------------")
						fmt.Println("Mestre: " + commands[1])

						i, err := strconv.ParseInt(commands[3], 10, 64)

						if err == nil {
							t := time.Unix(i, 0)
							fmt.Println("Partida termina daqui:", int(t.Sub(time.Now()).Seconds()), "segundos.")
						}
						fmt.Println("Dica: " + commands[2])
						fmt.Println("Aguardando definição do jogador da vez...")

					}
				}

			case "round-player":
				{
					// jogador da vez

					//mostrar
					if commands[1] == client.name {
						fmt.Println("É a sua vez de perguntar !")
						fmt.Print("Digite a pergunta: ")

						question, _ := reader.ReadString('\n')

						client.send("player-question::" + question)

						// mestre
					} else {
						fmt.Println("Jogador da vez: " + commands[1])
					}
				}

			case "player-question":
				{
					fmt.Println("Pergunta do jogador: " + commands[1])
					fmt.Println("--> " + commands[2])

					validResponse := false

					for !validResponse {
						fmt.Println("Responda 's' para SIM e 'n' para não.")
						fmt.Print("\\> ")
						response, _ := reader.ReadString('\n')

						if response == "s\n" {
							validResponse = true
							client.send("master-response::true")

						} else if response == "n\n" {
							validResponse = true
							client.send("master-response::false")
						}
					}

				}

			case "master-response":
				{
					if !client.isMaster {

						if commands[1] != client.name {

							fmt.Printf("%s perguntou: %s\n", commands[1], commands[2])
						}

						if commands[3] == "true" {

							fmt.Println("A resposta para a pergunta é: SIM! ")
						} else if commands[3] == "false" {
							fmt.Println("A resposta para a pergunta é: NÃO ):")
						}

						if commands[1] == client.name {

							fmt.Println("Você tem direito a tentar uma resposta.")
							fmt.Print("\\> ")

							response, _ := reader.ReadString('\n')

							client.send("player-response::" + response)
						}

					}
				}

			case "player-response":
				{
					if !client.isMaster {
						if commands[1] == client.name {
							if commands[2] == "true" {
								fmt.Println("Parabéns! Você acertou !!! Sua pontuação: " + commands[3])
							} else {
								fmt.Println("Não foi dessa vez ): . Aguarde o fim da partida...")
							}
						} else {
							if commands[2] == "true" {
								fmt.Printf("%s acertou e fez %s pontos.\n", commands[1], commands[3])
							} else {
								fmt.Printf("%s error.\n", commands[1])
							}
						}
					}
				}

			case "game-end":
				{

				}

			default:
				{
					fmt.Println("Comando desconhecido.")
					fmt.Println(commands)

				}
			}

		}
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("[Ajuda] go run client/ 127.0.0.1:8080")
		return
	}

	args := os.Args[1:]
	serverAddress := args[0]

	conn, err := net.Dial("tcp4", serverAddress)

	if err != nil {
		fmt.Println("Não foi possível conectar, tente novamente.")
		return
	}

	client := Client{conn: conn, name: "", isMaster: false}
	go client.receiveMessages()

	fmt.Println("-------------------------------")
	fmt.Println("Bem vindo ao quem sou eu.")
	fmt.Println("-------------------------------")

	client.send("get-game-info")

	for {
	}

	// for {
	// 	reader := bufio.NewReader(os.Stdin)

	// 	select {
	// 	case msg := <-client.gameStatusChan:
	// 		{
	// 			if msg == "waiting" {
	// 				fmt.Println("Para começar, nos diga: quem é você?")
	// 				fmt.Print("\\> ")
	// 				text, _ := reader.ReadString('\n')

	// 				client.setName(text)
	// 			}
	// 		}
	// 	}

	// 	fmt.Print("> ")
	// 	text, _ := reader.ReadString('\n')

	// 	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	// 	_, err := conn.Write([]byte(text))
	// 	if err != nil {
	// 		fmt.Println("Error writing to stream.")
	// 		break
	// 	}
	// }
}
