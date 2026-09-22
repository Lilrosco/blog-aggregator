package main

import (
	"context"
	"fmt"
	"time"
)

func handleAgg(s *state, cmd command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rssFeed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")

	if err != nil {
		return fmt.Errorf("error fetching feed: %w\n", err)
	}

	fmt.Println(rssFeed)

	return nil
}
