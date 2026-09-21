package main

import (
	"context"
	"fmt"
	"os"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(
		context.Background(),
	)

	if err != nil {
		fmt.Errorf("could not delete all users: %w", err)
		os.Exit(1)
	}

	return err
}
