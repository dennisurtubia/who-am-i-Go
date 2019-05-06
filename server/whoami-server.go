package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	listener, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Println(err.Error())
	}

	log.SetPrefix("[whoami-server] ")
	log.Println("Executando em " + listener.Addr().String())

	clientManager := ClientManager{
		clients:    make(map[*Client]bool),
	}

	gameManager := GameManager{}

	clientManager.gameManager = &gameManager
	gameManager.clientManager = &clientManager

	// go clientManager.start()
	go gameManager.start()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		client := &Client{socket: connection}
		go clientManager.handleClient(client)
		// clientManager.register <- client
		// clientManager.clientEnter(client)

		// go clientManager.receive(client)
		// go clientManager.send(client)
	}
}
