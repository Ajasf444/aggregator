package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (r *RSSFeed) UnescapeStrings() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
	for i := range r.Channel.Item {
		r.Channel.Item[i].unescapeStrings()
	}
}

func (r *RSSItem) unescapeStrings() {
	r.Title = html.UnescapeString(r.Title)
	r.Description = html.UnescapeString(r.Description)
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")
	// TODO: Create a Client with a Timeout
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	rssFeed := &RSSFeed{}
	if err := xml.Unmarshal(data, rssFeed); err != nil {
		return &RSSFeed{}, err
	}
	rssFeed.UnescapeStrings()
	return rssFeed, nil
}
