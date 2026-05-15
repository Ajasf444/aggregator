package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Ajasf444/aggregator/internal/config"
	"github.com/Ajasf444/aggregator/internal/database"
	"github.com/Ajasf444/aggregator/internal/rss"
	"github.com/google/uuid"
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
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil {
		return err
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register command expecting name argument")
	}
	ctx := context.Background()
	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]}
	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	fmt.Println("User was created.")
	fmt.Println(user)
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	s.db.DeleteUsers(ctx)
	fmt.Println("Database reset.")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	// TODO: optimization would be to build the string here rather than invoke strings.Join()
	found := false
	for i, user := range users {
		if !found && (user == s.cfg.CurrentUserName) {
			found = true
			users[i] += " (current)"
			user = users[i]
		}
		users[i] = "* " + user
	}
	fmt.Println(strings.Join(users, "\n"))
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	ctx := context.Background()
	rssFeed, err := fetchFeed(ctx, URL)
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rss.RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &rss.RSSFeed{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &rss.RSSFeed{}, err
	}
	rssFeed := &rss.RSSFeed{}
	if err := xml.Unmarshal(data, rssFeed); err != nil {
		return &rss.RSSFeed{}, err
	}
	rssFeed.UnescapeStrings()
	return rssFeed, nil
}

func handlerAddFeed(s *state, cmd command) error {
	// TODO: add logic
	return nil
}
