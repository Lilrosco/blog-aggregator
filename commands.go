package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.Name]

	if !exists {
		return errors.New("Command does not exist")
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func registerCommands() *commands {
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("register", handleRegister)
	cmds.register("login", handleLogin)
	cmds.register("reset", handleReset)
	cmds.register("users", handleGetUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	return &cmds
}
