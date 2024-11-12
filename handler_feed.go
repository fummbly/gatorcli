package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fummbly/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd input) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_reps>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to parse time: %v", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Failed to get next feed to fectch")
		return
	}

	log.Printf("Found %s feed to fetch", feed.Name)
	scrapeFeed(s.db, feed)

}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark the feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Failed to fetch data from feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        int32(uuid.New().ID()),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt.Time,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}

	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))

}

func handlerRemoveFeed(s *state, cmd input) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <feed-url>", cmd.Name)
	}

	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Failed to find feed: %v", err)
	}

	err = s.db.RemoveFeed(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("Failed to remove feed: %v", err)
	}

	fmt.Printf("Successfully deleted: %s\n", feed.Name)
	return nil

}

func handlerAddFeed(s *state, cmd input, user database.User) error {
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

func handlerGetFeeds(s *state, cmd input) error {

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
