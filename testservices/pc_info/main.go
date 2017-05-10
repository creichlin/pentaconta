package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	cwd, _ := os.Getwd()
	fmt.Println("arguments: " + strings.Join(os.Args[1:], " "))
	fmt.Println("cwd: " + cwd)
	time.Sleep(time.Second)
}
