package config

import "testing"

func TestText(t *testing.T) {
	cfg := &Config{TextOverride: map[string]string{"toc": "目录", "go-home": ""}}
	for key, want := range map[string]string{
		"toc":                "目录",
		"search-placeholder": "Search posts...",
		"not-found":          "Page not found.",
		"go-home":            "",
		"copy":               "Copy",
	} {
		if got := cfg.Text(key); got != want {
			t.Errorf("Text(%q) = %q, want %q", key, got, want)
		}
	}
	if got := (&Config{}).Text("toc"); got != "Table of Contents" {
		t.Errorf("default toc = %q", got)
	}
}
