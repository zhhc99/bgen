package pandoc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Result struct {
	Body string
	TOC  string
}

func Convert(markdown []byte) (*Result, error) {
	args := []string{
		"-f", "markdown+tex_math_dollars+pipe_tables+fenced_code_blocks+implicit_figures",
		"-t", "html",
		"--standalone", // required: pandoc only emits <nav id="TOC"> in standalone mode
		"--mathjax",
		"--no-highlight", // disable pandoc's built-in highlighting; we use highlight.js in base.html
		"--toc",
		"--toc-depth=3",
	}

	cmd := exec.Command("pandoc", args...)
	cmd.Stdin = bytes.NewReader(markdown)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pandoc: %w\n%s", err, stderr.String())
	}

	toc, body := splitTOC(extractBody(stdout.String()))
	return &Result{Body: injectCopyButtons(body), TOC: toc}, nil
}

// injectCopyButtons inserts a copy button inside every <pre> block so that
// the JS in copy.js can bind to it without touching the DOM structure.
func injectCopyButtons(html string) string {
	const btn = `<button class="copy-btn" aria-label="Copy">` +
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" ` +
		`stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` +
		`<rect x="9" y="9" width="13" height="13" rx="2"/>` +
		`<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>` +
		`</svg></button>`
	return strings.ReplaceAll(html, `</pre>`, btn+`</pre>`)
}

// extractBody pulls the content between <body> and </body> from a standalone
// HTML document, falling back to the full string if the tags are absent.
func extractBody(html string) string {
	start := strings.Index(html, "<body>")
	end := strings.LastIndex(html, "</body>")
	if start == -1 || end == -1 || end <= start {
		return html
	}
	return strings.TrimSpace(html[start+len("<body>") : end])
}

// splitTOC extracts the <nav id="TOC">…</nav> block pandoc prepends in --toc
// mode, handling nested <nav> tags to locate the correct closing tag.
func splitTOC(html string) (toc, body string) {
	const marker = `<nav id="TOC"`
	start := strings.Index(html, marker)
	if start == -1 {
		return "", html
	}

	depth, i := 0, start
	for i < len(html) {
		switch {
		case strings.HasPrefix(html[i:], "<nav"):
			depth++
			i += 4
		case strings.HasPrefix(html[i:], "</nav>"):
			depth--
			i += 6
			if depth == 0 {
				return strings.TrimSpace(html[start:i]), strings.TrimSpace(html[i:])
			}
		default:
			i++
		}
	}
	return "", html
}
