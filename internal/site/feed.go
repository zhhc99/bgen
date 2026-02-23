package site

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const feedMaxItems = 10

var (
	reCopyBtn = regexp.MustCompile(`<button[^>]*class="copy-btn"[^>]*>.*?</button>`)
	reAnySrc  = regexp.MustCompile(`src="([^"]+)"`)
	reAnyHref = regexp.MustCompile(`href="([^"]+)"`)
)

type rssRoot struct {
	XMLName   xml.Name   `xml:"rss"`
	Version   string     `xml:"version,attr"`
	ContentNS string     `xml:"xmlns:content,attr"`
	AtomNS    string     `xml:"xmlns:atom,attr"`
	Channel   rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	AtomLink    rssAtomLink `xml:"atom:link"`
	Description string      `xml:"description"`
	Items       []rssItem   `xml:"item"`
}

type rssAtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type rssItem struct {
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	GUID           rssGUID  `xml:"guid"`
	PubDate        string   `xml:"pubDate"`
	Description    rssCDATA `xml:"description"`
	ContentEncoded rssCDATA `xml:"content:encoded"`
}

type rssGUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type rssCDATA struct {
	Value string `xml:",cdata"`
}

func buildFeed(outDir string, s *Site) error {
	if s.Config.BaseURL == "" {
		return nil
	}

	posts := s.Posts
	if len(posts) > feedMaxItems {
		posts = posts[:feedMaxItems]
	}

	items := make([]rssItem, 0, len(posts))
	for _, p := range posts {
		permalink := s.Config.BaseURL + p.URL
		items = append(items, rssItem{
			Title:          p.Title,
			Link:           permalink,
			GUID:           rssGUID{IsPermaLink: true, Value: permalink},
			PubDate:        p.Date.UTC().Format(time.RFC1123Z),
			Description:    rssCDATA{p.Summary},
			ContentEncoded: rssCDATA{buildContent(&p, s.Config.BaseURL, p.URL)},
		})
	}

	feed := rssRoot{
		Version:   "2.0",
		ContentNS: "http://purl.org/rss/1.0/modules/content/",
		AtomNS:    "http://www.w3.org/2005/Atom",
		Channel: rssChannel{
			Title:       s.Config.Title,
			Link:        s.Config.BaseURL,
			AtomLink:    rssAtomLink{Href: s.Config.BaseURL + "/feed.xml", Rel: "self", Type: "application/rss+xml"},
			Description: s.Config.Title,
			Items:       items,
		},
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	out := append([]byte(xml.Header), data...)
	return os.WriteFile(filepath.Join(outDir, "feed.xml"), out, 0644)
}

func buildContent(p *Post, baseURL, postURL string) string {
	content := stripUIElements(string(p.Content))
	content = absolutifyURLs(content, baseURL, postURL)
	if p.Cover == "" {
		return content
	}
	coverURL := absolutifyURLs(`src="`+p.Cover+`"`, baseURL, postURL)
	cover := `<img ` + coverURL + ` alt="cover" style="max-width:100%"><br>` + "\n"
	return cover + content
}

func stripUIElements(html string) string {
	return reCopyBtn.ReplaceAllString(html, "")
}

func absolutifyURLs(html, baseURL, postURL string) string {
	postBase := baseURL + postURL
	rewrite := func(raw string) string {
		if strings.Contains(raw, "://") || strings.HasPrefix(raw, "//") {
			return raw // already absolute
		}
		if strings.HasPrefix(raw, "/") {
			return baseURL + raw
		}
		if strings.HasPrefix(raw, "./") {
			return postBase + raw[2:]
		}
		return postBase + raw // bare relative
	}
	html = reAnySrc.ReplaceAllStringFunc(html, func(m string) string {
		url := m[5 : len(m)-1]
		return `src="` + rewrite(url) + `"`
	})
	html = reAnyHref.ReplaceAllStringFunc(html, func(m string) string {
		url := m[6 : len(m)-1]
		return `href="` + rewrite(url) + `"`
	})
	return html
}
