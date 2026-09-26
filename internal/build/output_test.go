package build

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOutputPath(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "project")
	if err := os.MkdirAll(filepath.Join(root, "content"), 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "blog.yaml")
	if err := os.WriteFile(file, []byte("title: Test"), 0644); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(parent, "alias")
	if err := os.Symlink(root, alias); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"", root, parent, string(filepath.Separator), file,
		filepath.Join(root, "content"), filepath.Join(root, "content/output"),
		filepath.Join(root, "static/output"), filepath.Join(root, "layouts"),
		filepath.Join(root, ".git/objects"), filepath.Join(root, ".agents"), filepath.Join(root, ".codex"),
		alias, filepath.Join(alias, "content/new/output"),
	} {
		if _, err := outputPath(root, path); err == nil {
			t.Errorf("accepted unsafe output: %q", path)
		}
	}
	for _, path := range []string{filepath.Join(root, "output"), filepath.Join(root, "dist/site"), filepath.Join(parent, "public")} {
		if got, err := outputPath(root, path); err != nil || got != path {
			t.Errorf("rejected output %s: %s (%v)", path, got, err)
		}
	}
	external := filepath.Join(parent, "assets")
	if err := os.Mkdir(external, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(root, "static")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{external, filepath.Join(external, "output")} {
		if _, err := outputPath(root, path); err == nil {
			t.Errorf("accepted output inside linked source: %s", path)
		}
	}
}
