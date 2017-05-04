package main

import (
	"fmt"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Unstable main started")
	time.Sleep(time.Millisecond * 100)
	syscall.Exit(2)
}
