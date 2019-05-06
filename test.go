package main

import "time"
import "fmt"

func responseSent(c chan bool) {
	time.Sleep(time.Second * 7)
	c <- true
}

func waitingLoop(c chan bool)  {

	uptimeTicker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-uptimeTicker.C:
			fmt.Println("trocando manager")
		case <-c:
			fmt.Println("pronto")
			uptimeTicker.Stop()
			return
		}

	}

}

func main()  {

	c := make(chan bool, 1)
	go responseSent(c)
	waitingLoop(c)
}

// precisa travar no loop
// enquanto 