package main

import (
	"fmt"
	"net"
)

func main() {

	fmt.Println("Starting server...")

	listener, err := net.Listen("tcp4", ":8080")
	if err != nil {
		fmt.Println(err.Error())
	}

	manager := ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.start()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		client := &Client{socket: connection, data: make(chan []byte)}
		manager.register <- client

		go manager.receive(client)
		go manager.send(client)

	}

}
