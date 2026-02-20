package build

import (
	"fmt"

	"github.com/zhhc99/bgen/internal/config"
	"github.com/zhhc99/bgen/internal/site"
)

func Run(projectRoot string) error {
	return runWithConfig(projectRoot, func(cfg *config.Config) {})
}

func RunDev(projectRoot string) error {
	return runWithConfig(projectRoot, func(cfg *config.Config) {
		cfg.BasePath = ""
	})
}

func runWithConfig(projectRoot string, patch func(*config.Config)) error {
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	patch(cfg)
	s := site.New(cfg)
	if err := s.Build(projectRoot); err != nil {
		return fmt.Errorf("building site: %w", err)
	}
	fmt.Println("build complete -> output/")
	return nil
}
