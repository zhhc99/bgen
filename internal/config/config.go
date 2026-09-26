package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type HeroConfig struct {
	Header  string `yaml:"header"`
	Content string `yaml:"content"`
}

type FrontMatterDefaults struct {
	Author string `yaml:"author"`
}

type NavItem struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

type Config struct {
	Title               string              `yaml:"title"`
	BaseURL             string              `yaml:"base_url"`
	BasePath            string              `yaml:"-"` // derived from BaseURL, e.g. "/~john"
	Hero                HeroConfig          `yaml:"hero"`
	Nav                 []NavItem           `yaml:"nav"`
	L10n                map[string]string   `yaml:"l10n"`
	FrontMatterDefaults FrontMatterDefaults `yaml:"front-matter-defaults"`
}

func Load(projectRoot string) (*Config, error) {
	path := filepath.Join(projectRoot, "blog.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading blog.yaml: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing blog.yaml: %w", err)
	}
	for _, item := range cfg.Nav {
		if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.URL) == "" {
			return nil, fmt.Errorf("nav entries require title and url")
		}
	}

	if cfg.BaseURL != "" {
		cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
		if u, err := url.Parse(cfg.BaseURL); err == nil {
			cfg.BasePath = strings.TrimRight(u.Path, "/")
		}
	}

	return &cfg, nil
}

func (c *Config) NavTitle(path string) string {
	for _, item := range c.Nav {
		if item.URL == path {
			return item.Title
		}
	}
	return ""
}

func (c *Config) NavURL(path string) string {
	if strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "//") {
		return c.BasePath + path
	}
	return path
}
