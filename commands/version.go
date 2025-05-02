package commands

import "fmt"

var Version = "dev"

func RunVersion() {
	fmt.Println("pushy version", Version)
}