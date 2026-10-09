package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Lilrosco/blog-aggregator/internal/database"
)

func handleBrowsePosts(s *state, cmd command, user database.User) error {
	// limit expected as int32
	limit := int32(2)

	if len(cmd.Args) == 1 {
		// Ensure value is in bounds for 32 Bit Int
		conv, err := strconv.ParseInt(cmd.Args[0], 10, 32)

		if err != nil {
			fmt.Printf("Could not parse %s as LIMIT defaulting to 2\n", cmd.Args[0])
		} else if conv > 0 {
			limit = int32(conv)
		}
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit: limit,
		},
	)

	if err != nil {
		return fmt.Errorf("could fetch posts for username: %s | %w", user.Name, err)
	}

	for _, post := range posts {
		prettyPrintPost(post)
	}

	return nil
}

func prettyPrintPost(post database.GetPostsForUserRow) {
	fmt.Printf("Feed Name: %s\n", post.FeedName)
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("Description: %s\n", post.Description.String)
	fmt.Printf("URL: %s\n", post.Url)
	fmt.Printf("Published At: %s\n", post.PublishedAt.Time)
	fmt.Println("--------------------------------")
}
