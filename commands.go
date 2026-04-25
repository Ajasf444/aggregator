package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ajasf444/aggregator/internal/config"
	"github.com/Ajasf444/aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	callback, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unable to execute command %v", cmd.name)
	}
	return callback(s, cmd)
}

func NewCommands() commands {
	return commands{
		handlers: map[string]func(*state, command) error{},
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login command expecting username argument")
	}
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("%v has been set as user.\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register command expecting name argument")
	}
	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]}
	ctx := context.Background()
	// TODO: invoke queries from state
	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}
