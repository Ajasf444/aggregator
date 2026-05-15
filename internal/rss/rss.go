package rss

import "html"

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
