package main

import (
	"context"
	"errors"
	"fmt"
	"database/sql"
	"log"
	"os"

	"github.com/Lilrosco/blog-aggregator/internal/config"
	"github.com/Lilrosco/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

// For pq database error
const errUniqueViolation = "23505"

type state struct {
	db *database.Queries
	cfg *config.Config
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
		os.Exit(1)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if s.cfg.CurrentUserName == "" {
			return errors.New("No user is logged in!")
		}

		user, err := s.db.GetUser(
			context.Background(),
			s.cfg.CurrentUserName,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("User with name: %s does not exists\n", s.cfg.CurrentUserName)
			}

			return fmt.Errorf("could not fetch user: %s - %w", s.cfg.CurrentUserName, err)
		}

		return handler(s, cmd, user)
	}
}
