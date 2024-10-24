package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/fummbly/gatorcli/internal/rss"
)

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error making request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error getting response from url: %v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error reading body: %v", err)
	}

	feed := rss.RSSFeed{}

	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error parsing the feed to struct: %v", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item

	}

	return &feed, nil

}
