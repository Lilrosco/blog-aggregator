package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func handleAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeDuration := cmd.Args[0]
	timeBetweenRequests, err := time.ParseDuration(timeDuration)

	if err != nil {
		return fmt.Errorf("error parsing duration for %s : %w", timeDuration, err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	log.Printf("Collecting feeds every %s...\n", timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		return fmt.Errorf("could not update feed's last_fetched: %s - %w", feed.Name, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rssFeed, err := fetchFeed(ctx, feed.Url)

	if err != nil {
		return fmt.Errorf("error fetching feed: %s - %w", feed.Url, err)
	}

	fmt.Printf("Feed: %s\n", rssFeed.Channel.Title)

	for index, item := range rssFeed.Channel.Item {
		fmt.Printf("  %d) Title: %s\n", index, item.Title)
	}

	return nil
}
