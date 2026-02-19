package site

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type searchItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Date  string `json:"date"`
}

func (s *Site) render(projectRoot, outPath string) error {
	searchEnabled := s.Config.Nav["search"] != ""
	tagsEnabled := s.Config.Nav["tags"] != ""

	// 预加载模板
	coreTemplates := []string{"index", "404", "single", "page"}
	for _, name := range coreTemplates {
		if _, err := s.getTemplate(projectRoot, name); err != nil {
			return fmt.Errorf("preloading template %s: %w", name, err)
		}
	}
	if searchEnabled {
		if _, err := s.getTemplate(projectRoot, "search"); err != nil {
			return fmt.Errorf("preloading template search: %w", err)
		}
	}
	if tagsEnabled {
		for _, name := range []string{"tags", "tag"} {
			if _, err := s.getTemplate(projectRoot, name); err != nil {
				return fmt.Errorf("preloading template %s: %w", name, err)
			}
		}
	}

	if err := s.copyStaticFiles(projectRoot, outPath); err != nil {
		return err
	}
	if err := s.copyCoverImages(outPath); err != nil {
		return err
	}
	if err := s.copyBundleImages(outPath); err != nil {
		return err
	}
	if searchEnabled {
		if err := s.writeSearchJSON(outPath); err != nil {
			return err
		}
	}

	type renderJob struct {
		path string
		name string
		data any
	}

	base := struct{ Site *Site }{Site: s}

	jobs := []renderJob{
		{"index.html", "index", base},
		{"404.html", "404", base},
	}
	if searchEnabled {
		jobs = append(jobs, renderJob{"search/index.html", "search", base})
	}
	if tagsEnabled {
		jobs = append(jobs, renderJob{"tags/index.html", "tags", base})
		for tag, posts := range s.Tags {
			tag, posts := tag, posts
			jobs = append(jobs, renderJob{
				"tags/" + tag + "/index.html",
				"tag",
				struct {
					Site  *Site
					Tag   string
					Posts []Post
				}{s, tag, posts},
			})
		}
	}
	for _, p := range s.Posts {
		p := p
		jobs = append(jobs, renderJob{
			"posts/" + p.Slug + "/index.html",
			"single",
			struct {
				Site *Site
				Post Post
			}{s, p},
		})
	}
	for _, pg := range s.Pages {
		pg := pg
		jobs = append(jobs, renderJob{
			pg.Slug + "/index.html",
			"page",
			struct {
				Site *Site
				Page Page
			}{s, pg},
		})
	}

	for _, job := range jobs {
		if err := s.renderPage(projectRoot, outPath, job.path, job.name, job.data); err != nil {
			return fmt.Errorf("rendering %s: %w", job.path, err)
		}
	}
	return nil
}

func (s *Site) renderPage(projectRoot, outPath, relPath, name string, data any) error {
	tmpl, err := s.getTemplate(projectRoot, name)
	if err != nil {
		return err
	}
	dest := filepath.Join(outPath, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.ExecuteTemplate(f, "base.html", data)
}

func (s *Site) getTemplate(projectRoot, name string) (*template.Template, error) {
	if tmpl, ok := s.templateCache[name]; ok {
		return tmpl, nil
	}
	tmpl, err := s.loadTemplate(projectRoot, name)
	if err != nil {
		return nil, err
	}
	s.templateCache[name] = tmpl
	return tmpl, nil
}

func (s *Site) loadTemplate(projectRoot, name string) (*template.Template, error) {
	read := func(embPath, userPath string) ([]byte, error) {
		if projectRoot != "" {
			if data, err := os.ReadFile(filepath.Join(projectRoot, userPath)); err == nil {
				return data, nil
			}
		}
		return embeddedFS.ReadFile(embPath)
	}

	baseData, err := read("templates/base.html", "layouts/base.html")
	if err != nil {
		return nil, fmt.Errorf("loading base.html: %w", err)
	}
	pageData, err := read("templates/"+name+".html", "layouts/"+name+".html")
	if err != nil {
		return nil, fmt.Errorf("loading %s.html: %w", name, err)
	}

	tmpl := template.New("base.html")
	if _, err := tmpl.Parse(string(baseData)); err != nil {
		return nil, err
	}
	if _, err := tmpl.New(name + ".html").Parse(string(pageData)); err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (s *Site) copyStaticFiles(projectRoot, outPath string) error {
	err := fs.WalkDir(embeddedFS, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel("static", path)
		dest := filepath.Join(outPath, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		src, err := embeddedFS.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		return writeFile(dest, src)
	})
	if err != nil {
		return err
	}
	return copyDir(filepath.Join(projectRoot, "static"), outPath)
}

func (s *Site) copyCoverImages(outPath string) error {
	for _, p := range s.Posts {
		if p.CoverSrc == "" {
			continue
		}
		dest := filepath.Join(outPath, "posts", p.Slug, "cover"+filepath.Ext(p.CoverSrc))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		f, err := os.Open(p.CoverSrc)
		if err != nil {
			return err
		}
		if err := writeFile(dest, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

func (s *Site) copyBundleImages(outPath string) error {
	for _, p := range s.Posts {
		if len(p.BundleImages) == 0 {
			continue
		}
		postDir := filepath.Join(outPath, "posts", p.Slug)
		if err := os.MkdirAll(postDir, 0755); err != nil {
			return err
		}
		for relPath, absPath := range p.BundleImages {
			dest := filepath.Join(postDir, relPath)
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return err
			}
			f, err := os.Open(absPath)
			if err != nil {
				return fmt.Errorf("opening %s: %w", absPath, err)
			}
			if err := writeFile(dest, f); err != nil {
				f.Close()
				return fmt.Errorf("copying %s: %w", relPath, err)
			}
			f.Close()
		}
	}
	return nil
}

func (s *Site) writeSearchJSON(outPath string) error {
	items := make([]searchItem, len(s.Posts))
	for i, p := range s.Posts {
		items[i] = searchItem{
			Title: p.Title,
			URL:   p.URL,
			Date:  p.Date.Format("2006-01-02"),
		}
	}
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outPath, "search.json"), data, 0644)
}

func writeFile(dest string, src io.Reader) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, src)
	return err
}

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dest, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		return writeFile(target, f)
	})
}
