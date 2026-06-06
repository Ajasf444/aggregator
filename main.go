package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Ajasf444/aggregator/internal/config"
	"github.com/Ajasf444/aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	allArgs := os.Args
	if len(allArgs) == 1 {
		fmt.Println("Not enough arguments provided.")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println("unable to connect to database")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s := NewState(&cfg, dbQueries)
	commands := NewCommands()

	args := allArgs[1:]
	cmdName, cmdArgs := args[0], args[1:]
	err = commands.run(s, command{name: cmdName, args: cmdArgs})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
