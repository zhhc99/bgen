package build

import (
	"fmt"

	"github.com/zhhc99/bgen/internal/config"
	"github.com/zhhc99/bgen/internal/site"
)

func Run(projectRoot, outDir string) error {
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	s := site.New(cfg)
	if err := s.Build(projectRoot, outDir); err != nil {
		return fmt.Errorf("building site: %w", err)
	}
	fmt.Printf("build complete -> %s\n", outDir)
	return nil
}

func RunDev(projectRoot, outDir string) error {
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	cfg.BasePath = ""
	s := site.New(cfg)
	if err := s.Build(projectRoot, outDir); err != nil {
		return fmt.Errorf("building site: %w", err)
	}
	return nil
}
