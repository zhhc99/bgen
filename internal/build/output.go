package build

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhhc99/bgen/internal/site"
)

func buildOutput(projectRoot, outDir string, s *site.Site) error {
	outDir, err := outputPath(projectRoot, outDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outDir), 0755); err != nil {
		return err
	}
	workDir, err := os.MkdirTemp(filepath.Dir(outDir), ".bgen-build-*")
	if err != nil {
		return err
	}
	defer func() {
		if workDir != "" {
			os.RemoveAll(workDir)
		}
	}()

	next := filepath.Join(workDir, "next")
	if err := s.Build(projectRoot, next); err != nil {
		return err
	}
	previous := filepath.Join(workDir, "previous")
	hadOutput := false
	if err := os.Rename(outDir, previous); err == nil {
		hadOutput = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("saving previous output: %w", err)
	}
	if err := os.Rename(next, outDir); err != nil {
		if hadOutput {
			if restoreErr := os.Rename(previous, outDir); restoreErr != nil {
				workDir = ""
				return fmt.Errorf("replacing output: %w; previous output kept at %s: %v", err, previous, restoreErr)
			}
		}
		return fmt.Errorf("replacing output: %w", err)
	}
	return nil
}

func outputPath(projectRoot, outDir string) (string, error) {
	if outDir == "" {
		return "", fmt.Errorf("output directory must not be empty")
	}
	if info, err := os.Lstat(filepath.Clean(outDir)); err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("output must be a directory, not a file or symlink: %s", outDir)
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}
	root, err := realPath(projectRoot)
	if err != nil {
		return "", err
	}
	out, err := realPath(outDir)
	if err != nil {
		return "", err
	}
	if within(out, root) {
		return "", fmt.Errorf("output must not contain the project: %s", outDir)
	}
	for _, name := range []string{"content", "static", "layouts", "blog.yaml", ".git", ".agents", ".codex"} {
		path := filepath.Join(root, name)
		resolved, err := realPath(path)
		if err != nil {
			return "", err
		}
		for _, protected := range []string{path, resolved} {
			if within(protected, out) || within(out, protected) {
				return "", fmt.Errorf("output overlaps %s: %s", name, outDir)
			}
		}
	}
	return out, nil
}

func within(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func realPath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err == nil {
		return resolved, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	parent, err := realPath(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(path)), nil
}
