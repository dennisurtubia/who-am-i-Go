package main

import "time"
import "fmt"





func main()  {

	time1 := time.Now()
	time.Sleep(4 * time.Second)
	time2 := time.Now()

	fmt.Println(int(100 * (1/time2.Sub(time1).Seconds())))
}

// precisa travar no loop
// enquanto 