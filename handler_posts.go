package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/fummbly/gatorcli/internal/database"
	"github.com/fummbly/gatorcli/internal/rss"
	"github.com/google/uuid"
)

func createPost(db *database.Queries, feedID int32, rssItem rss.RSSItem) (database.Post, error) {

	timeLayot := "Mon, 2 Jan 2006 15:04:04 -0700"
	parsedTime, err := time.Parse(timeLayot, rssItem.PubDate)
	if err != nil {
		return database.Post{}, fmt.Errorf("Failed to parse time: %v", err)
	}

	post, err := db.CreatePost(context.Background(), database.CreatePostParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Title:     rssItem.Title,
		Url:       rssItem.Link,
		Description: sql.NullString{
			String: rssItem.Description, Valid: true,
		},
		PublishedAt: parsedTime,
		FeedID:      feedID,
	})
	if err != nil {
		return database.Post{}, fmt.Errorf("Failed to create post %s: %v", rssItem.Title, err)
	}

	return post, nil

}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("Usage: %s <limit> (optional)", cmd.Name)
	}
	fmt.Println("Starting browsing")
	var limit int
	var err error

	if len(cmd.Args) == 1 {
		limit, err = strconv.Atoi(cmd.Args[0])

		if err != nil {
			return fmt.Errorf("Failed to convert limit value: %v", err)
		}
	} else {
		limit = 2
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Failed to get posts for user: %v", err)
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil

}

func printPost(post database.GetPostsForUserRow) {
	fmt.Println()
	fmt.Printf("Post: %s\n", post.Title)
	fmt.Printf("Description: %s\n", post.Description.String)
	fmt.Printf("Publish Date: %v\n", post.PublishedAt)

}
