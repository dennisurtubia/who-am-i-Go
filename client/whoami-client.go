package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp4", ":8080")

	if err != nil {
		fmt.Println("Não foi possível conectar, por favor, tente novamente")
		os.Exit(1)
	}
		
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		fmt.Fprintf(conn, "get-game-info")

		message := scanner.Text()
		commands := strings.Split(message, "::")
		fmt.Println(commands)
		fmt.Println("Bem vindo ao quem sou eu")

		conn.Write([]byte("get-game-info"))

		switch commands[0] {
		case "get-game-info":
			{
				fmt.Println("chegou aq")
				// stateUser, numberOfPlayres := commands[1], commands[2]
				// if stateUser == "waiting" {
				// 	fmt.Println("Aguardando o inicio da partida..."+"\n"+"Número de jogadores no momento: ", numberOfPlayres)
				// 	fmt.Print("\n" + "Para começar, nos forneça seu nickgame: ")
				// }
			}
		case "set-name":
			{

			}
		case "game-init":
			{

			}
		}
	}
}
