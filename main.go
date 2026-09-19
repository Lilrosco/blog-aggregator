package main

import (
	"log"
	"os"

	"github.com/Lilrosco/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func registerCommands() *commands {
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	return &cmds
}

func main() {
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := registerCommands()
	args := os.Args

	if len(args) < 2 {
		log.Fatalf("no arguments were passed into program")
		os.Exit(1)
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err = cmds.run(programState, cmd); err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
