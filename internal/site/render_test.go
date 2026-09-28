package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zhhc99/bgen/internal/config"
)

func TestRenderStylesheets(t *testing.T) {
	defaultStyle, err := embeddedFS.ReadFile("static/style.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name     string
		custom   bool
		css      string
		style    string
		basePath string
	}{
		{name: "default theme"},
		{name: "incremental CSS", custom: true, css: ":root { --accent: red; }"},
		{name: "empty custom CSS", custom: true},
		{name: "replacement theme", style: "body { color: blue; }"},
		{name: "replacement theme with custom CSS", custom: true, css: "body { color: red; }", style: "body { color: blue; }"},
		{name: "subpath", custom: true, css: ":root { --accent: red; }", basePath: "/blog"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root, out := t.TempDir(), t.TempDir()
			if err := os.Mkdir(filepath.Join(root, "static"), 0755); err != nil {
				t.Fatal(err)
			}
			if tt.custom {
				if err := os.WriteFile(filepath.Join(root, "static/custom.css"), []byte(tt.css), 0644); err != nil {
					t.Fatal(err)
				}
			}
			wantStyle := string(defaultStyle)
			if tt.style != "" {
				wantStyle = tt.style
				if err := os.WriteFile(filepath.Join(root, "static/style.css"), []byte(tt.style), 0644); err != nil {
					t.Fatal(err)
				}
			}
			s := New(&config.Config{
				Title: "Test Blog", BasePath: tt.basePath,
				Nav: []config.NavItem{{Title: "Search", URL: "/search/"}, {Title: "Tags", URL: "/tags/"}},
			})
			post := Post{Title: "Hello", Slug: "hello", URL: "/posts/hello/"}
			s.Posts = []Post{post}
			s.Tags["test"] = []Post{post}
			s.Pages["about"] = Page{Title: "About", Slug: "about", URL: "/about/"}
			if err := s.render(root, out); err != nil {
				t.Fatal(err)
			}
			style, err := os.ReadFile(filepath.Join(out, "style.css"))
			if err != nil || string(style) != wantStyle {
				t.Fatalf("base stylesheet changed unexpectedly: %v", err)
			}
			custom, err := os.ReadFile(filepath.Join(out, "custom.css"))
			if tt.custom {
				if err != nil || string(custom) != tt.css {
					t.Fatalf("custom stylesheet was not copied unchanged: %v", err)
				}
			} else if !os.IsNotExist(err) {
				t.Fatalf("unexpected custom stylesheet: %v", err)
			}
			for _, page := range []string{
				"index.html", "404.html", "search/index.html", "tags/index.html",
				"tags/test/index.html", "posts/hello/index.html", "about/index.html",
			} {
				data, err := os.ReadFile(filepath.Join(out, page))
				if err != nil {
					t.Fatal(err)
				}
				html := string(data)
				baseLink := `<link rel="stylesheet" href="` + tt.basePath + `/style.css">`
				customLink := `<link rel="stylesheet" href="` + tt.basePath + `/custom.css">`
				if !strings.Contains(html, baseLink) {
					t.Errorf("%s: missing base stylesheet", page)
				}
				if tt.custom {
					if strings.Count(html, customLink) != 1 || strings.Index(html, customLink) < strings.Index(html, baseLink) {
						t.Errorf("%s: custom stylesheet must appear once after the base stylesheet", page)
					}
					if strings.Index(html, customLink) < strings.Index(html, "github-dark.min.css") {
						t.Errorf("%s: custom stylesheet must follow syntax highlighting styles", page)
					}
				} else if strings.Contains(html, "/custom.css") {
					t.Errorf("%s: references nonexistent custom stylesheet", page)
				}
			}
		})
	}
}
