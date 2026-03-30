package main

import (
	"errors"
	"fmt"

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
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("%v has been set as user", s.cfg.CurrentUserName)
	return nil
}
