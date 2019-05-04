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
	log.Println("running on " + listener.Addr().String())

	cManager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	gManager := GameManager{}

	cManager.gManager = &gManager
	gManager.cManager = &cManager

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
