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

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	user, err := s.db.GetUser(
		context.Background(),
		s.cfg.CurrentUserName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("User with name: %s does not exists\n", user.Name)
		}

		return fmt.Errorf("could not fetch user: %s - %w", user.Name, err)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID: uuid.New(),
			Name: name,
			Url: url,
			UserID: user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)

	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		if pqErr.Code == errUniqueViolation {
			return fmt.Errorf("Feed with url: %s already exists\n", url)
		} else if err != nil {
			return fmt.Errorf("could not create feed: %w", err)
		}
	}

	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("Url: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)
	fmt.Printf("Username: %s\n", user.Name)
	fmt.Println("Feed has been created")

	_, err = createFeedFollow(s, feed, user)

	if err != nil {
		return err
	}

	return nil
}

func handleGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUserName(
		context.Background(),
	)

	if err != nil {
		return fmt.Errorf("could fetch all feeds: %w", err)
	}

	for _, feed := range feeds {
		if feed.Name_2.Valid {
			fmt.Printf("* %s\n", feed.Name)
			fmt.Printf("* %s\n", feed.Url)
			fmt.Printf("* %s\n", feed.Name_2.String)
		} else {
			return fmt.Errorf("could not properly fetch user attached to feed - Name: %s | Url: %s", feed.Name, feed.Url)
		}
	}

	return nil
}
