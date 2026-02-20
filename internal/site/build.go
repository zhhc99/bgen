package site

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/zhhc99/bgen/internal/content"
	"github.com/zhhc99/bgen/internal/pandoc"
)

const (
	contentDir = "content"
	outputDir  = "output"
	postsDir   = "content/posts"
)

var coverExts = []string{"jpg", "jpeg", "png", "webp", "gif"}

func (s *Site) Build(projectRoot string) error {
	outPath := filepath.Join(projectRoot, outputDir)
	if err := os.MkdirAll(outPath, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	if err := s.loadPosts(filepath.Join(projectRoot, postsDir)); err != nil {
		return fmt.Errorf("loading posts: %w", err)
	}
	if err := s.loadPages(filepath.Join(projectRoot, contentDir)); err != nil {
		return fmt.Errorf("loading pages: %w", err)
	}
	if err := s.render(projectRoot, outPath); err != nil {
		return fmt.Errorf("rendering: %w", err)
	}
	if err := buildFeed(projectRoot, s); err != nil {
		return fmt.Errorf("building feed: %w", err)
	}
	return nil
}

func (s *Site) loadPosts(postsPath string) error {
	entries, err := os.ReadDir(postsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading posts dir: %w", err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(postsPath, entry.Name())

		var post *Post
		if entry.IsDir() {
			post, err = s.loadBundlePost(fullPath)
		} else if filepath.Ext(entry.Name()) == ".md" {
			post, err = s.loadFlatPost(fullPath)
		} else {
			continue
		}

		if err != nil {
			return fmt.Errorf("%s: %w", entry.Name(), err)
		}
		if post != nil {
			s.Posts = append(s.Posts, *post)
		}
	}

	sort.Slice(s.Posts, func(i, j int) bool {
		return s.Posts[i].Date.After(s.Posts[j].Date)
	})
	for _, p := range s.Posts {
		for _, tag := range p.Tags {
			s.Tags[tag] = append(s.Tags[tag], p)
		}
	}
	return nil
}

func (s *Site) loadFlatPost(mdPath string) (*Post, error) {
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return nil, err
	}
	pf, err := content.Parse(data)
	if err != nil {
		return nil, err
	}
	if pf.Front.Ignore {
		return nil, nil
	}

	base := strings.TrimSuffix(filepath.Base(mdPath), ".md")
	slug := pf.Front.Slug
	if slug == "" {
		slug = base
	}
	return s.buildPost(pf, slug, findCover(filepath.Dir(mdPath), base))
}

func (s *Site) loadBundlePost(bundleDir string) (*Post, error) {
	indexPath := filepath.Join(bundleDir, "index.md")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	pf, err := content.Parse(data)
	if err != nil {
		return nil, err
	}
	if pf.Front.Ignore {
		return nil, nil
	}

	slug := pf.Front.Slug
	if slug == "" {
		slug = filepath.Base(bundleDir)
	}

	// bundle 封面优先找 cover.*, 退回到 index.*
	coverSrc := findCover(bundleDir, "cover")
	if coverSrc == "" {
		coverSrc = findCover(bundleDir, "index")
	}

	post, err := s.buildPost(pf, slug, coverSrc)
	if err != nil {
		return nil, err
	}
	post.BundleImages = extractImageRefs(pf.Body, bundleDir)
	return post, nil
}

func (s *Site) buildPost(pf *content.ParsedFile, slug, coverSrc string) (*Post, error) {
	result, err := pandoc.Convert(pf.Body)
	if err != nil {
		return nil, err
	}

	summary := pf.Front.Summary
	if summary == "" {
		summary = content.ExtractSummary(pf.Body)
	}

	author := pf.Front.Author
	if author == "" {
		author = s.Config.FrontMatterDefaults.Author
	}

	var coverURL string
	if coverSrc != "" {
		coverURL = "/posts/" + slug + "/cover" + filepath.Ext(coverSrc)
	}

	return &Post{
		Title:    pf.Front.Title,
		Date:     pf.Front.Date,
		Tags:     pf.Front.Tags,
		Slug:     slug,
		URL:      "/posts/" + slug + "/",
		Summary:  summary,
		Author:   author,
		Cover:    coverURL,
		CoverSrc: coverSrc,
		Content:  template.HTML(result.Body),
		TOC:      template.HTML(result.TOC),
	}, nil
}

func (s *Site) loadPages(contentPath string) error {
	entries, err := os.ReadDir(contentPath)
	if err != nil {
		return fmt.Errorf("reading content dir: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := filepath.Join(contentPath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		pf, err := content.Parse(data)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		result, err := pandoc.Convert(pf.Body)
		if err != nil {
			return err
		}
		slug := strings.TrimSuffix(entry.Name(), ".md")
		s.Pages[slug] = Page{
			Title:   pf.Front.Title,
			Slug:    slug,
			URL:     "/" + slug + "/",
			Content: template.HTML(result.Body),
		}
	}
	return nil
}

func findCover(dir, base string) string {
	for _, ext := range coverExts {
		p := filepath.Join(dir, base+"."+ext)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// extractImageRefs 从 markdown 正文中提取本地图片引用.
// 返回 map: 相对路径 -> 绝对路径.
func extractImageRefs(body []byte, bundleDir string) map[string]string {
	re := regexp.MustCompile(`!\[[^\]]*\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
	images := make(map[string]string)
	for _, m := range re.FindAllSubmatch(body, -1) {
		if len(m) < 2 {
			continue
		}
		imgPath := string(m[1])
		if strings.HasPrefix(imgPath, "http://") || strings.HasPrefix(imgPath, "https://") {
			continue
		}
		absPath := filepath.Join(bundleDir, imgPath)
		if _, err := os.Stat(absPath); err == nil {
			images[imgPath] = absPath
		}
	}
	return images
}
