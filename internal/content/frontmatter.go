package content

import (
	"bytes"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title   string    `yaml:"title"`
	Date    time.Time `yaml:"date"`
	Tags    []string  `yaml:"tags"`
	Slug    string    `yaml:"slug"`
	Author  string    `yaml:"author"`
	Summary string    `yaml:"summary"`
	Ignore  bool      `yaml:"ignore"`
}

type ParsedFile struct {
	Front FrontMatter
	Body  []byte
}

func Parse(data []byte) (*ParsedFile, error) {
	data = bytes.TrimSpace(data)
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return nil, fmt.Errorf("missing front matter: file must start with ---")
	}
	rest := data[4:]
	end := bytes.Index(rest, []byte("\n---"))
	if end == -1 {
		return nil, fmt.Errorf("front matter not closed")
	}
	yamlPart := rest[:end]
	body := bytes.TrimSpace(rest[end+4:])

	var fm FrontMatter
	if err := yaml.Unmarshal(yamlPart, &fm); err != nil {
		return nil, fmt.Errorf("parsing front matter yaml: %w", err)
	}
	return &ParsedFile{Front: fm, Body: body}, nil
}
