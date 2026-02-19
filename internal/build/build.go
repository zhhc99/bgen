package build

import (
	"fmt"

	"github.com/zhhc99/bgen/internal/config"
	"github.com/zhhc99/bgen/internal/site"
)

func Run(projectRoot string) error {
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	s := site.New(cfg)
	if err := s.Build(projectRoot); err != nil {
		return fmt.Errorf("building site: %w", err)
	}

	fmt.Println("build complete -> output/")
	return nil
}
