package main

import "time"
import "fmt"
import "strconv"





func main()  {

	time1 := time.Now()
	time.Sleep(1 * time.Second)
	time2 := time.Now()

	fmt.Println(strconv.Itoa(int(100 * (1/time2.Sub(time1).Seconds()))))
}

// precisa travar no loop
// enquanto 