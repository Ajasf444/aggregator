package main

import (
	"fmt"

	"github.com/Ajasf444/aggregator/internal/config"
	"github.com/Ajasf444/aggregator/internal/database"
)

const URL = "https://wagslane.dev/index.xml"

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

func NewCommands() *commands {
	c := &commands{
		handlers: map[string]func(*state, command) error{},
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerGetUsers)
	c.register("feeds", handlerFeeds)
	c.register("agg", middlewareLoggedIn(handlerAggregate))
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	return c
}

func (c *commands) run(s *state, cmd command) error {
	callback, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unable to execute command %v", cmd.name)
	}
	return callback(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func NewState(cfg *config.Config, db *database.Queries) *state {
	return &state{
		cfg: cfg,
		db:  db,
	}
}
