package build_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zhhc99/bgen/internal/build"
)

// mustWrite 是测试辅助函数, 写文件失败直接 fatal.
func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// makeProject 在临时目录里创建一个最小博客项目.
func makeProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	mustWrite(t, filepath.Join(dir, "blog.yaml"), `
title: Test Blog
base_url: https://example.com
nav:
  - title: search
    url: /search/
  - title: tags
    url: /tags/
  - title: About
    url: /about/
`)

	mustWrite(t, filepath.Join(dir, "content/posts/hello.md"), `---
title: Hello World
date: 2024-01-01
tags: [go, test]
---

这是第一篇测试文章的正文.
`)

	mustWrite(t, filepath.Join(dir, "content/posts/math.md"), `---
title: Math Post
date: 2024-02-01
---

支持 TeX: $E = mc^2$
`)

	mustWrite(t, filepath.Join(dir, "content/about.md"), `---
title: About
---

关于页面.
`)

	return dir
}

func TestBuild_Smoke(t *testing.T) {
	if _, err := exec.LookPath("pandoc"); err != nil {
		t.Skip("pandoc not found in PATH")
	}

	dir := makeProject(t)
	outDir := filepath.Join(dir, "output")

	if err := build.Run(dir, outDir); err != nil {
		t.Fatalf("build.Run: %v", err)
	}

	// 断言关键输出文件存在
	wantFiles := []string{
		"index.html",
		"404.html",
		"posts/hello/index.html",
		"posts/math/index.html",
		"tags/index.html",
		"tags/go/index.html",
		"tags/test/index.html",
		"search/index.html",
		"search.json",
		"about/index.html",
		"style.css",
	}
	for _, rel := range wantFiles {
		path := filepath.Join(outDir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing output file: %s", rel)
		}
	}
}

func TestBuild_PostContent(t *testing.T) {
	if _, err := exec.LookPath("pandoc"); err != nil {
		t.Skip("pandoc not found in PATH")
	}

	dir := makeProject(t)
	outDir := filepath.Join(dir, "output")

	if err := build.Run(dir, outDir); err != nil {
		t.Fatalf("build.Run: %v", err)
	}

	// 文章页应包含标题
	postHTML, err := os.ReadFile(filepath.Join(outDir, "posts/hello/index.html"))
	if err != nil {
		t.Fatalf("reading post html: %v", err)
	}
	if !bytes.Contains(postHTML, []byte("Hello World")) {
		t.Error("post page missing title")
	}
	if !bytes.Contains(postHTML, []byte("这是第一篇测试文章的正文")) {
		t.Error("post page missing body content")
	}

	// 首页应包含两篇文章的标题
	indexHTML, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		t.Fatalf("reading index html: %v", err)
	}
	if !bytes.Contains(indexHTML, []byte("Hello World")) {
		t.Error("index page missing post title")
	}
	if !bytes.Contains(indexHTML, []byte("Math Post")) {
		t.Error("index page missing post title")
	}
}

func TestBuild_MissingConfig(t *testing.T) {
	dir := t.TempDir() // 空目录, 没有 blog.yaml

	err := build.Run(dir, filepath.Join(dir, "output"))
	if err == nil {
		t.Fatal("expected error for missing blog.yaml, got nil")
	}
}

func TestBuild_Idempotent(t *testing.T) {
	if _, err := exec.LookPath("pandoc"); err != nil {
		t.Skip("pandoc not found in PATH")
	}

	dir := makeProject(t)
	outDir := filepath.Join(dir, "output")

	// 连续构建两次, 都应该成功
	if err := build.Run(dir, outDir); err != nil {
		t.Fatalf("first build: %v", err)
	}
	if err := build.Run(dir, outDir); err != nil {
		t.Fatalf("second build: %v", err)
	}
}

func TestBuild_Navigation(t *testing.T) {
	if _, err := exec.LookPath("pandoc"); err != nil {
		t.Skip("pandoc not found in PATH")
	}
	dir := makeProject(t)
	out := filepath.Join(dir, "output")
	mustWrite(t, filepath.Join(dir, "blog.yaml"), `
title: Test Blog
base_url: https://example.com/~alice
nav:
  - title: GitHub
    url: https://github.com/zhhc99
  - title: Find
    url: /search/
  - title: Topics
    url: /tags/
`)
	if err := build.Run(dir, out); err != nil {
		t.Fatal(err)
	}
	html, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	nav := string(html)
	last := -1
	for _, link := range []string{
		`href="https://github.com/zhhc99">GitHub</a>`,
		`href="/~alice/search/">Find</a>`,
		`href="/~alice/tags/">Topics</a>`,
	} {
		pos := strings.Index(nav, link)
		if pos <= last {
			t.Fatalf("missing or unordered navigation link: %s", link)
		}
		last = pos
	}
	if strings.Contains(nav, `href="/~alice/about/"`) {
		t.Fatal("unconfigured page added to navigation")
	}
	if _, err := os.Stat(filepath.Join(out, "about/index.html")); err != nil {
		t.Fatal("unlisted page was not generated:", err)
	}
	for page, title := range map[string]string{"search": "Find", "tags": "Topics"} {
		html, err := os.ReadFile(filepath.Join(out, page, "index.html"))
		if err != nil || !strings.Contains(string(html), "<h1>"+title+"</h1>") {
			t.Fatalf("%s page does not use navigation title: %v", page, err)
		}
	}
	if err := build.RunDev(dir, out); err != nil {
		t.Fatal(err)
	}
	html, err = os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil || !strings.Contains(string(html), `href="/search/">Find</a>`) {
		t.Fatalf("preview navigation uses deployment prefix: %v", err)
	}
	for _, tc := range []struct {
		nav    string
		search bool
		tags   bool
	}{
		{"nav:\n  - title: Search\n    url: /search/\n", true, false},
		{"nav:\n  - title: Tags\n    url: /tags/\n", false, true},
		{"nav: []\n", false, false},
		{"", false, false},
	} {
		mustWrite(t, filepath.Join(dir, "blog.yaml"), "title: Test\n"+tc.nav)
		if err := build.Run(dir, out); err != nil {
			t.Fatal(err)
		}
		for rel, exists := range map[string]bool{
			"search/index.html":  tc.search,
			"search.json":        tc.search,
			"tags/index.html":    tc.tags,
			"tags/go/index.html": tc.tags,
			"about/index.html":   true,
			"feed.xml":           false,
		} {
			_, err := os.Stat(filepath.Join(out, rel))
			if exists && err != nil || !exists && !os.IsNotExist(err) {
				t.Errorf("%s: exists = %v, stat = %v", rel, exists, err)
			}
		}
	}
}

func TestBuild_RemovesUnpublishedContent(t *testing.T) {
	if _, err := exec.LookPath("pandoc"); err != nil {
		t.Skip("pandoc not found in PATH")
	}
	dir := makeProject(t)
	out := filepath.Join(dir, "output")
	mustWrite(t, filepath.Join(dir, "content/posts/bundle/index.md"), "---\ntitle: Bundle\ntags: [bundle]\n---\n![asset](asset.png)\n")
	mustWrite(t, filepath.Join(dir, "content/posts/bundle/asset.png"), "asset")
	mustWrite(t, filepath.Join(dir, "content/posts/hello.png"), "cover")
	mustWrite(t, filepath.Join(dir, "static/removed.txt"), "old asset")
	mustWrite(t, filepath.Join(dir, "static/kept.txt"), "current asset")
	if err := build.Run(dir, out); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"posts/hello/cover.png", "posts/bundle/asset.png", "about/index.html"} {
		if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
			t.Fatal(err)
		}
	}
	for _, rel := range []string{"posts/hello.md", "posts/bundle/index.md", "about.md"} {
		path := filepath.Join(dir, "content", rel)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		mustWrite(t, path, strings.Replace(string(data), "---\n", "---\nignore: true\n", 1))
	}
	for _, rel := range []string{"content/posts/math.md", "static/removed.txt"} {
		if err := os.Remove(filepath.Join(dir, rel)); err != nil {
			t.Fatal(err)
		}
	}
	if err := build.Run(dir, out); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"posts", "about", "tags/go", "tags/test", "tags/bundle", "removed.txt"} {
		if _, err := os.Stat(filepath.Join(out, rel)); !os.IsNotExist(err) {
			t.Errorf("stale output remains: %s (%v)", rel, err)
		}
	}
	for _, rel := range []string{"index.html", "search.json", "feed.xml", "tags/index.html"} {
		data, err := os.ReadFile(filepath.Join(out, rel))
		if err != nil {
			t.Fatal(err)
		}
		for _, title := range []string{"Hello World", "Math Post", "Bundle"} {
			if bytes.Contains(data, []byte(title)) {
				t.Errorf("%s still contains %s", rel, title)
			}
		}
	}
	if data, err := os.ReadFile(filepath.Join(out, "kept.txt")); err != nil || string(data) != "current asset" {
		t.Fatalf("static asset lost: %v", err)
	}
}

func TestBuild_FailurePreservesOutput(t *testing.T) {
	if _, err := exec.LookPath("pandoc"); err != nil {
		t.Skip("pandoc not found in PATH")
	}
	dir := makeProject(t)
	out := filepath.Join(dir, "output")
	if err := build.Run(dir, out); err != nil {
		t.Fatal(err)
	}
	before := make(map[string]string)
	err := filepath.WalkDir(out, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		before[path] = string(data)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, "layouts/index.html"), `{{define "content"}}{{.MissingField}}{{end}}`)
	if err := build.Run(dir, out); err == nil {
		t.Fatal("expected template rendering failure")
	}
	for path, want := range before {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Errorf("previous output changed: %s (%v)", path, err)
		}
	}
	leftovers, err := filepath.Glob(filepath.Join(dir, ".bgen-build-*"))
	if err != nil || len(leftovers) != 0 {
		t.Fatalf("temporary files remain: %v (%v)", leftovers, err)
	}
}
