/*
A very simple TCP client written in Go.

This is a toy project that I used to learn the fundamentals of writing
Go code and doing some really basic network stuff.

Maybe it will be fun for you to read. It's not meant to be
particularly idiomatic, or well-written for that matter.
*/
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
)



func main() {


	conn, err := net.Dial("tcp", ":8080")

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}

	go readConnection(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
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

func readConnection(conn net.Conn) {

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("mensagem", message)
	}
}

func handleCommands(text string) bool {
	r, err := regexp.Compile("^%.*%$")
	if err == nil {
		if r.MatchString(text) {

			switch {
			case text == "%quit%":
				fmt.Println("\b\bServer is leaving. Hanging up.")
				os.Exit(0)
			}

			return true
		}
	}

	return false
}