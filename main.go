package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Lilrosco/blog-aggregator/internal/config"
	"github.com/Lilrosco/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func registerCommands() *commands {
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	return &cmds
}

func main() {
	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	dbURL := cfg.DBURL
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}

	dbQueries := database.New(db)

	programState := &state{
		db: dbQueries,
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
