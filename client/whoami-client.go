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
	conn net.Conn
	name string
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
				}
			}
		}
	}
}

func main() {

	args := os.Args[1:]
	serverAddress := args[0]

	conn, err := net.Dial("tcp4", serverAddress)

	if err != nil {
		fmt.Println("Não foi possível conectar, tente novamente.")
		return
	}

	client := Client{conn: conn}
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
