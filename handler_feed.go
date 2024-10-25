package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fummbly/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_reps>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to parse time: %v", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("Failed scraping: %v", err)
		}

	}

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feed: %v", err)
	}

	fmt.Printf("Fetching feed: %s\n", feed.Name)

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now().UTC(),
		LastFetchedAt: sql.NullTime{
			Time: time.Now().UTC(),
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to mark feed as fetched: %v", err)

	}

	feedRSS, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Failed to fetch feed rss: %v", err)
	}

	fmt.Println()
	fmt.Printf("Title: %s\n", feedRSS.Channel.Title)
	fmt.Println("=================================")
	fmt.Println("Items")

	for _, item := range feedRSS.Channel.Item {
		fmt.Printf(" * %s\n", item.Title)
	}

	return nil

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url", cmd.Name)
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create feed: %v", err)
	}

	fmt.Println("Feed Successfully created")
	fmt.Println("===============================")
	printFeed(feed, user)
	fmt.Println("===============================")

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
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
