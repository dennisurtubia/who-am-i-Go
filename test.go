package main

import (
	"fmt"
	"time"
)

func main() {
	status := 0
	// statusChange := make(chan bool)
	// status <- 0

	go func() {
		time.Sleep(6 * time.Second)
		status = 1
		// statusChange <- true
		// status <- 1
	}()

	for status == 0 {
		fmt.Println("new loop")
		time.Sleep(2 * time.Second)
	}

	fmt.Println("exit")

}
