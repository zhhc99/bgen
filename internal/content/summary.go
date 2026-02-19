package content

import (
	"bytes"
	"strings"
)

const maxSummaryRunes = 150

func ExtractSummary(body []byte) string {
	for _, para := range bytes.Split(body, []byte("\n\n")) {
		p := strings.TrimSpace(string(para))
		if p == "" {
			continue
		}
		// 跳过标题、代码块、图片、表格
		if strings.HasPrefix(p, "#") || strings.HasPrefix(p, "```") ||
			strings.HasPrefix(p, "![") || strings.HasPrefix(p, "|") {
			continue
		}
		p = strings.ReplaceAll(p, "\n", " ")
		return truncateRunes(p, maxSummaryRunes)
	}
	return ""
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
