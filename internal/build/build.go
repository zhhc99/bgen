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
	if err := buildOutput(projectRoot, outDir, site.New(cfg)); err != nil {
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
	if err := buildOutput(projectRoot, outDir, site.New(cfg)); err != nil {
		return fmt.Errorf("building site: %w", err)
	}
	return nil
}
