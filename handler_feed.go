package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fummbly/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {

	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %v", err)
	}

	fmt.Println(feed.Channel.Item[0])

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url", cmd.Name)
	}

	currUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("Failed to get current user: %v", err)
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    currUser.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create feed: %v", err)
	}

	fmt.Println("Feed Successfully created")
	fmt.Println("===============================")
	printFeed(feed, currUser)
	fmt.Println("===============================")

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create follow: %v", err)
	}

	fmt.Printf("%s is now following %s\n", follow.User, follow.Feed)

	return nil

}

func handlerGetFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds: %v", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Failed to get user: %v", err)
		}
		fmt.Println("===================")
		printFeed(feed, user)
		fmt.Println("===================")
		fmt.Println()

	}

	return nil

}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:        %d\n", feed.ID)
	fmt.Printf("* CreatedAt:        %v\n", feed.CreatedAt)
	fmt.Printf("* UpdatedAt:        %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:        %s\n", feed.Name)
	fmt.Printf("* URL:        %s\n", feed.Url)
	fmt.Printf("* User:        %s\n", user.Name)

}
