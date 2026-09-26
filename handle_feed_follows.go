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

func createFeedFollow(s *state, feed database.Feed, user database.User) (database.CreateFeedFollowRow, error) {
	feed_follow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID: uuid.New(),
			FeedID: feed.ID,
			UserID: user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)

	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		if pqErr.Code == errUniqueViolation {
			return database.CreateFeedFollowRow{}, fmt.Errorf("FeedFollow with user|feed: %s already exists\n", user.Name, feed.Name)
		} else if err != nil {
			return database.CreateFeedFollowRow{}, fmt.Errorf("could not create feed_follow: %w", err)
		}
	}

	fmt.Printf("Feed Name: %s\n", feed_follow.FeedName)
	fmt.Printf("User Name: %s\n", feed_follow.UserName)
	fmt.Println("FeedFollow has been created")

	return feed_follow, nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeed(
		context.Background(),
		url,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Feed with url: %s does not exists\n", url)
		}

		return fmt.Errorf("could not fetch feed: %s - %w", url, err)
	}

	// Also create entry for feed follow
	_, err = createFeedFollow(s, feed, user)

	if err != nil {
		return err
	}

	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	feed_follows, err := s.db.GetFeedFollowForUser(
		context.Background(),
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("could fetch all feed_follows for user - %s | %w", user.Name, err)
	}

	fmt.Printf("* Feeds for User: %s\n", user.Name)

	for _, feed_follow := range feed_follows {
		fmt.Printf("* %s\n", feed_follow.FeedName)
	}

	return nil
}
