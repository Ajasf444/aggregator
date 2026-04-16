package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Ajasf444/aggregator/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	dbURL := cfg.DBURL
	db, err := sql.Open("postgres", dbURL)
	s := state{cfg: &cfg}
	commands := NewCommands()
	commands.register("login", handlerLogin)

	allArgs := os.Args
	if len(allArgs) == 1 {
		fmt.Println("Not enough arguments provided.")
		os.Exit(1)
	}
	args := allArgs[1:]
	cmdName := args[0]
	callback, ok := commands.handlers[cmdName]
	if !ok {
		fmt.Printf("Command %v not found.\n", cmdName)
		os.Exit(1)
	}
	// TODO: handle this better
	err = callback(&s, command{name: cmdName, args: args[1:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
