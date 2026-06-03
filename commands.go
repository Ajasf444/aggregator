package main

import (
	"context"
	"errors"
	"fmt"
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

func NewCommands() *commands {
	c := &commands{
		handlers: map[string]func(*state, command) error{},
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerGetUsers)
	c.register("agg", handlerAggregate)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	return c
}

func (c *commands) run(s *state, cmd command) error {
	callback, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unable to execute command %v", cmd.name)
	}
	return callback(s, cmd)
}

func NewState(cfg *config.Config, db *database.Queries) *state {
	return &state{
		cfg: cfg,
		db:  db,
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
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
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
	s.db.DeleteFeeds(ctx)
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
	rssFeed, err := rss.FetchFeed(ctx, URL)
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("addfeed command expecting a name and url")
	}
	ctx := context.Background()
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return err
	}
	fmt.Println("Feed was created.")
	fmt.Println(feed)

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", follow)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%v %v %v\n", feed.Feedname, feed.Url, feed.Username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	URL := cmd.args[0]
	ctx := context.Background()
	feed, err := s.db.GetFeedFromURL(ctx, URL)
	if err != nil {
		return err
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	follow, err := s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", follow)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for i := range feedFollows {
		fmt.Println(feedFollows[i])
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return nil
		}
		return handler(s, cmd, user)
	}
}
