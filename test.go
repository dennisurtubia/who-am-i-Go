package main

import (
	"log"
	"strconv"
	"time"
)

func main() {

	now := time.Now().Add(time.Hour * 2)
	tstmp := strconv.FormatInt(now.Unix(), 10)

	str := tstmp
	i, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		log.Print(err)

	} else {
		t := time.Unix(i, 0)
		log.Print(t)
	}
}

// precisa travar no loop
// enquanto
