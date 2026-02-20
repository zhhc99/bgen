package site

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"time"
)

type rssRoot struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func buildFeed(projectRoot string, s *Site) error {
	items := make([]rssItem, 0, len(s.Posts))
	for _, p := range s.Posts {
		items = append(items, rssItem{
			Title:       p.Title,
			Link:        s.Config.BaseURL + p.URL,
			PubDate:     p.Date.UTC().Format(time.RFC1123Z),
			Description: p.Summary,
		})
	}

	feed := rssRoot{
		Version: "2.0",
		Channel: rssChannel{
			Title:       s.Config.Title,
			Link:        s.Config.BaseURL,
			Description: s.Config.Title,
			Items:       items,
		},
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	out := append([]byte(xml.Header), data...)
	return os.WriteFile(filepath.Join(projectRoot, "output", "feed.xml"), out, 0644)
}
