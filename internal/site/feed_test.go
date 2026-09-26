package site

import (
	"strings"
	"testing"
)

func TestFeedCoverTitle(t *testing.T) {
	post := &Post{Title: `A "quote" & <tag>`, Cover: "/posts/test/cover.png"}
	content := buildContent(post, "https://example.com", "/posts/test/")
	if !strings.Contains(content, `alt="A &#34;quote&#34; &amp; &lt;tag&gt;"`) {
		t.Fatalf("cover title is missing or unescaped: %s", content)
	}
}
