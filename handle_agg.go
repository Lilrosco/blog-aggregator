package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Lilrosco/blog-aggregator/internal/database"
	"github.com/lib/pq"
	"github.com/google/uuid"
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

	// fmt.Printf("Feed: %s\n", rssFeed.Channel.Title)

	for _, item := range rssFeed.Channel.Item {
		// fmt.Printf("  %d) Title: %s | Date: %s\n", index, item.Title, item.PubDate)
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)

		if err != nil {
			parsedTime, err = time.Parse(time.RFC1123, item.PubDate)

			if err != nil {
				fmt.Printf("Could not parse time for Feed ID: %s | Item Title: %s | Item PubDate: %s", feed.ID, item.Title, item.PubDate)
			}
		}

		post, err := s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID: uuid.New(),
				FeedID: feed.ID,
				Title: item.Title,
				Url: item.Link,
				Description: newNullString(item.Description),
				PublishedAt: newNullTime(parsedTime),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		)

		if pqErr, ok := errors.AsType[*pq.Error](err); ok {
			if pqErr.Code == errUniqueViolation {
				// Do nothing
				fmt.Printf("Post with url: %s already exists\n", item.Link)
			} else if err != nil {
				return fmt.Errorf("could not create post: %w", err)
			}
		} else {
			fmt.Printf("ID: %s\n", post.ID)
			fmt.Printf("Title: %s\n", post.Title)
			fmt.Println("Post has been created")
		}
	}

	return nil
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid: true,
	}
}

func newNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time: t,
		Valid: true,
	}
}
