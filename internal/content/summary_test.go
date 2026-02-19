package content_test

import (
	"strings"
	"testing"

	"github.com/zhhc99/bgen/internal/content"
)

func TestExtractSummary(t *testing.T) {
	t.Run("取第一段", func(t *testing.T) {
		body := []byte("第一段内容.\n\n第二段内容.")
		got := content.ExtractSummary(body)
		if got != "第一段内容." {
			t.Errorf("got %q", got)
		}
	})

	t.Run("跳过标题取正文", func(t *testing.T) {
		body := []byte("# 标题\n\n正文段落.")
		got := content.ExtractSummary(body)
		if got != "正文段落." {
			t.Errorf("got %q", got)
		}
	})

	t.Run("跳过代码块", func(t *testing.T) {
		body := []byte("```go\nfmt.Println()\n```\n\n正文段落.")
		got := content.ExtractSummary(body)
		if got != "正文段落." {
			t.Errorf("got %q", got)
		}
	})

	t.Run("跳过图片", func(t *testing.T) {
		body := []byte("![图注](image.png)\n\n正文段落.")
		got := content.ExtractSummary(body)
		if got != "正文段落." {
			t.Errorf("got %q", got)
		}
	})

	t.Run("跳过表格", func(t *testing.T) {
		body := []byte("| A | B |\n|---|---|\n| 1 | 2 |\n\n正文段落.")
		got := content.ExtractSummary(body)
		if got != "正文段落." {
			t.Errorf("got %q", got)
		}
	})

	t.Run("CJK 150 字截断", func(t *testing.T) {
		// 200 个中文字符
		long := strings.Repeat("字", 200)
		body := []byte(long)
		got := content.ExtractSummary(body)
		runes := []rune(got)
		// 应截断到 150 rune + "..."
		if len(runes) != 153 {
			t.Errorf("got %d runes, want 153 (150 + ...)", len(runes))
		}
		if !strings.HasSuffix(got, "...") {
			t.Error("truncated string should end with ...")
		}
	})

	t.Run("短文本不截断", func(t *testing.T) {
		body := []byte("短文本.")
		got := content.ExtractSummary(body)
		if got != "短文本." {
			t.Errorf("got %q", got)
		}
	})

	t.Run("空 body", func(t *testing.T) {
		got := content.ExtractSummary([]byte{})
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("全是标题无正文", func(t *testing.T) {
		body := []byte("# H1\n\n## H2\n\n### H3")
		got := content.ExtractSummary(body)
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})
}
