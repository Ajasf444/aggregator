package main

import "github.com/Ajasf444/aggregator/internal/config"

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
}
