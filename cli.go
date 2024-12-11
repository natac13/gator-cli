package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/natac13/gator-cli/internal/config"
	"github.com/natac13/gator-cli/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func newState() (*state, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}
	return &state{
		config: &cfg,
	}, nil
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	if user.Name != cmd.args[0] {
		return errors.New("user not found")
	}
	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Logged in as %s", cmd.args[0])
	return nil
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func newCommands() *commands {
	cmds := &commands{
		cmds: make(map[string]func(*state, command) error),
	}
	return cmds
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmds[cmd.name]
	if !ok {
		return errors.New("unknown command")
	}
	return f(s, cmd)
}

func handlerRegister(s *state, cmd command) error {
	// ensure the name was provided
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}

	// create the user
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	// set the current user
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s created and logged in\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	// ensure the name was provided
	if len(cmd.args) != 0 {
		return errors.New("no arguments expected")
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("All users deleted")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("no arguments expected")
	}

	users, err := s.db.GetUsers((context.Background()))
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.config.CurrentUserName {
			fmt.Print(" (current)\n")
		} else {
			fmt.Print("\n")
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("no arguments expected")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collection feeds every %s\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
