package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	conn           net.Conn
	gameStatusChan chan string
}

func (client *Client) getGameStatus() {
	client.send("get-game-info")
}

func (client *Client) setName(name string) {
	client.send("set-name::" + name)
}

func (client *Client) send(message string) {
	client.conn.Write([]byte(message + "\n"))
}

func (client *Client) receiveMessages() {
	scanner := bufio.NewScanner(client.conn)

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
							log.Println("Atualmente há:", commands[2], "jogadores conectados.")
							log.Println("Próxima partida inicia em:", int(t.Sub(time.Now()).Seconds()), "segundos.")
						}

					} else if commands[1] == "ingame" {
						i, err := strconv.ParseInt(commands[2], 10, 64)

						if err == nil {
							t := time.Unix(i, 0)
							log.Println("Partida em andamento.")
							log.Println("Próxima partida inicia em:", int(t.Sub(time.Now()).Seconds()), "segundos.")
						}
					}
					client.gameStatusChan <- commands[1]
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
		log.Panicln("Não foi possível conectar, tente novamente.")
	}

	client := Client{conn: conn, gameStatusChan: make(chan string, 1)}
	go client.receiveMessages()

	log.Println("-------------------------------")
	log.Println("Bem vindo ao quem sou eu.")
	log.Println("-------------------------------")

	client.getGameStatus()

	for {
		reader := bufio.NewReader(os.Stdin)

		select {
		case msg := <-client.gameStatusChan:
			{
				if msg == "waiting" {
					log.Println("Para começar, nos diga: quem é você?")
					log.Print("\\> ")
					text, _ := reader.ReadString('\n')

					client.setName(text)
				}
			}
		}

		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}
