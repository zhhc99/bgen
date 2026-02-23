package build_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
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
  search: search
  tags: tags
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
