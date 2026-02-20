package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

var conflictPaths = []string{
	"blog.yaml",
	"content/about.md",
	"content/posts/my-post.md",
	"content/posts/my-bundle-post",
}

func Run(projectRoot string) error {
	for _, rel := range conflictPaths {
		if _, err := os.Stat(filepath.Join(projectRoot, rel)); err == nil {
			return fmt.Errorf("%q already exists; run bgen init in an empty directory", rel)
		}
	}

	files := map[string]string{
		"blog.yaml":                             blogYAML,
		"content/about.md":                      aboutMD,
		"content/posts/my-post.md":              myPostMD,
		"content/posts/my-bundle-post/index.md": myBundleMD,
	}

	for rel, content := range files {
		path := filepath.Join(projectRoot, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", rel, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", rel, err)
		}
	}
	assets := map[string][]byte{
		"content/posts/my-bundle-post/laifu-the-cat.webp": laifuTheCat,
	}
	for rel, data := range assets {
		path := filepath.Join(projectRoot, rel)
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", rel, err)
		}
	}

	fmt.Println("bgen: initialized! Run `bgen serve` to preview your blog.")
	fmt.Println("bgen: see https://github.com/zhhc99/bgen for usage, cover images, and theme customization.")
	return nil
}

const blogYAML = `title: Alice's Blog
base_url: https://example.com   # your deployed site URL
hero:
  header: Alice
  content: Welcome to my blog!
nav:
  search: search
  tags: tags
l10n:
  toc: Table of Contents
front-matter-defaults:
  author: Alice
`

const aboutMD = `---
title: About
---

Hi! This blog is powered by [bgen](https://github.com/zhhc99/bgen),
a static blog generator built with minimal cognitive load in mind.

No complex configuration, no theme jungle — just write Markdown and ship.
`

const myPostMD = `---
title: My First Post
date: 2024-01-01
tags: [hello, markdown]
---

Welcome to bgen! Write your posts in Markdown right here.

## Code

Syntax highlighting:

` + "```" + `python
def greet(name: str) -> str:
    return f"Hello, {name}!"
` + "```" + `

## Math

Inline math: $E = mc^2$

Display math:

$$
\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
$$

## Images with captions

bgen uses Pandoc, so image captions just work:

![This text becomes a figure caption](https://picsum.photos/500)

## What's next?

- To add a cover image, place a same-named image file next to this post (e.g. ` + "`my-post.jpg`" + `).
- To customize the theme, edit ` + "`static/style.css`" + ` or templates under ` + "`layouts/`" + `.
- You can download templates of style sheets / layout files from the GitHub repo below.
- See https://github.com/zhhc99/bgen for full documentation.
`

const myBundleMD = `---
title: A Bundle Post
date: 2024-01-02
tags: [bundle]
---

When a post has multiple images or assets, organize it as a bundle:
place ` + "`index.md`" + ` and all related files in a folder together.

` + "```" + `
content/posts/my-bundle-post/
├── index.md           ← this file
├── cover.jpg          ← cover image
└── laifu-the-cat.webp ← referenced in the post body
` + "```" + `

Then reference images with relative paths:

![A cute cat "laifu"](./laifu-the-cat.webp)
`
