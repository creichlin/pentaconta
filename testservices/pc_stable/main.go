package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Stable main started")
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("I'm doing fine")
	}
}
