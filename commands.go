package main

import (
	"errors"

	"github.com/Ajasf444/aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("expecting username argument")
	}
	return nil
}
