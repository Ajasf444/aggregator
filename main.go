package main

import (
	"fmt"
	"os"

	"github.com/Ajasf444/aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	s := state{cfg: &cfg}
	commands := NewCommands()
	commands.register("login", handlerLogin)

	allArgs := os.Args
	if len(allArgs) == 1 {
		fmt.Println("Not enough arguments provided.")
		os.Exit(1)
	}
	args := allArgs[1:]
	cmdString := args[1]
	cmd, ok := commands.handlers[cmdString]
	if !ok {
		fmt.Printf("Command %v not found", cmd)
		os.Exit(1)
	}
	// TODO: write switch statement to handle cmd and arguments
}
