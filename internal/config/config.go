package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type HeroConfig struct {
	Header  string `yaml:"header"`
	Content string `yaml:"content"`
}

type FrontMatterDefaults struct {
	Author string `yaml:"author"`
}

type Config struct {
	Title               string              `yaml:"title"`
	BaseURL             string              `yaml:"base_url"`
	Hero                HeroConfig          `yaml:"hero"`
	Nav                 map[string]string   `yaml:"nav"`
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
	return &cfg, nil
}
