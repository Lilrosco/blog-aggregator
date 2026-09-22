package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Lilrosco/blog-aggregator/internal/database"
	"github.com/lib/pq"
	"github.com/google/uuid"
)

const errUniqueViolation = "23505"

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]
	user, err := s.db.GetUser(
		context.Background(),
		username,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("User with name: %s does not exists\n", username)
		}

		return fmt.Errorf("could not fetch user: %s - %w", username, err)
	}

	err = s.cfg.SetUser(user.Name)

	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: %s\n", username)

	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: username,
		},
	)

	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		if pqErr.Code == errUniqueViolation {
			return fmt.Errorf("User with name: %s already exists\n", username)
		} else if err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}
	}

	err = s.cfg.SetUser(username)

	if err != nil {
		return fmt.Errorf("could not set current user in config: %w", err)
	}

	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Println("Username has been created")

	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(
		context.Background(),
	)

	if err != nil {
		return fmt.Errorf("could not delete all users: %w", err)
	}

	fmt.Println("Users table has been reset")

	return nil
}

func handleGetUsers(s *state, cmd command) error {
	currentUserName := s.cfg.CurrentUserName

	users, err := s.db.GetUsers(
		context.Background(),
	)

	if err != nil {
		return fmt.Errorf("could fetch all users: %w", err)
	}

	for _, user := range users {
		if user.Name == currentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
