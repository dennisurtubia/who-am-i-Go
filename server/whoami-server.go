package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Println(err.Error())
	}

	log.SetPrefix("[whoami-server] ")
	log.Println("running on " + listener.Addr().String())

	gManager := GameManager{lobbyPlayers: make([]Player, 0)}

	cManager := ClientManager{
		gManager:   &gManager,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go cManager.start()
	go gManager.start()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		client := &Client{socket: connection, data: make(chan []byte)}
		cManager.register <- client

		go cManager.receive(client)
		go cManager.send(client)

	}

}
