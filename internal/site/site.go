// site.go
package site

import (
	"html/template"
	"time"

	"github.com/zhhc99/bgen/internal/config"
)

type Post struct {
	Title        string
	Date         time.Time
	Tags         []string
	Slug         string
	URL          string
	Summary      string
	Author       string
	Cover        string            // 生成后的 URL 路径, 空表示无封面
	CoverSrc     string            // 构建期使用的源文件绝对路径
	BundleImages map[string]string // markdown 中引用的图片: 相对路径 -> 绝对路径
	Content      template.HTML
	TOC          template.HTML
}

type Page struct {
	Title   string
	Slug    string
	URL     string
	Content template.HTML
}

type Site struct {
	Config        *config.Config
	Posts         []Post
	Tags          map[string][]Post
	Pages         map[string]Page
	templateCache map[string]*template.Template
}

func New(cfg *config.Config) *Site {
	return &Site{
		Config:        cfg,
		Tags:          make(map[string][]Post),
		Pages:         make(map[string]Page),
		templateCache: make(map[string]*template.Template),
	}
}
