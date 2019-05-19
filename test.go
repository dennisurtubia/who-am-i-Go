package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("> ")
	a, _ := reader.ReadString('\n')
	log.Println(a)
}

// precisa travar no loop
// enquanto
