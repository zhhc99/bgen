package content_test

import (
	"testing"
	"time"

	"github.com/zhhc99/bgen/internal/content"
)

func TestParse(t *testing.T) {
	t.Run("基本字段", func(t *testing.T) {
		input := []byte(`---
title: Hello
date: 2024-01-15
tags: [go, blog]
slug: my-slug
author: alice
---

正文内容.
`)
		pf, err := content.Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if pf.Front.Title != "Hello" {
			t.Errorf("title: got %q, want %q", pf.Front.Title, "Hello")
		}
		wantDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		if !pf.Front.Date.Equal(wantDate) {
			t.Errorf("date: got %v, want %v", pf.Front.Date, wantDate)
		}
		if len(pf.Front.Tags) != 2 || pf.Front.Tags[0] != "go" || pf.Front.Tags[1] != "blog" {
			t.Errorf("tags: got %v, want [go blog]", pf.Front.Tags)
		}
		if pf.Front.Slug != "my-slug" {
			t.Errorf("slug: got %q, want %q", pf.Front.Slug, "my-slug")
		}
		if pf.Front.Author != "alice" {
			t.Errorf("author: got %q, want %q", pf.Front.Author, "alice")
		}
		if string(pf.Body) != "正文内容." {
			t.Errorf("body: got %q", string(pf.Body))
		}
	})

	t.Run("ignore 字段", func(t *testing.T) {
		input := []byte("---\ntitle: Draft\nignore: true\n---\n\n内容\n")
		pf, err := content.Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !pf.Front.Ignore {
			t.Error("ignore should be true")
		}
	})

	t.Run("缺少 front matter", func(t *testing.T) {
		input := []byte("没有 front matter 的文件\n")
		_, err := content.Parse(input)
		if err == nil {
			t.Fatal("expected error for missing front matter")
		}
	})

	t.Run("front matter 未关闭", func(t *testing.T) {
		input := []byte("---\ntitle: Oops\n\n正文\n")
		_, err := content.Parse(input)
		if err == nil {
			t.Fatal("expected error for unclosed front matter")
		}
	})

	t.Run("空 body", func(t *testing.T) {
		input := []byte("---\ntitle: Empty\n---\n")
		pf, err := content.Parse(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(pf.Body) != 0 {
			t.Errorf("expected empty body, got %q", pf.Body)
		}
	})
}
