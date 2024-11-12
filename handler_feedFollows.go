package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fummbly/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd input, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to get feed: %v", err)
	}

	follows, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to created feed follow: %v", err)
	}

	fmt.Printf("%s is now following %s\n", follows.User, follows.Feed)
	return nil
}

func handlerFollowing(s *state, cmd input, user database.User) error {

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get following feeds: %v", err)
	}

	fmt.Println("Following =============")
	for _, feed := range following {

		fmt.Printf(" * %s\n", feed.Feed)

	}
	return nil
}

func handlerUnfollow(s *state, cmd input, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	err := s.db.DeleteFollow(context.Background(), database.DeleteFollowParams{
		Name: user.Name,
		Url:  url,
	})
	if err != nil {
		return fmt.Errorf("Failed to delete follow: %s", err)
	}

	fmt.Printf("%s has stopped following %s\n", user.Name, url)

	return nil
}
