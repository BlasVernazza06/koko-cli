package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

func TestEngineGenerateVFSSaaS(t *testing.T) {
	cfg := ScaffoldConfig{
		ProjectName:    "my-saas-app",
		Recipe:         "saas",
		PackageManager: "pnpm",
	}

	virtualFS, err := GenerateVFS(cfg)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	if virtualFS.FileCount() == 0 {
		t.Errorf("Expected files in VFS, got 0")
	}

	// Verify package.json updated in memory
	if !virtualFS.Exists("package.json") {
		t.Errorf("Expected package.json to exist in VFS")
	}

	var rootPkg map[string]interface{}
	if err := virtualFS.ReadJSON("package.json", &rootPkg); err != nil {
		t.Fatalf("Failed to read package.json from VFS: %v", err)
	}

	if rootPkg["name"] != "my-saas-app" {
		t.Errorf("Expected root package name 'my-saas-app', got '%v'", rootPkg["name"])
	}

	// Verify .env was generated in VFS
	if !virtualFS.Exists(".env") {
		t.Errorf("Expected .env in VFS")
	}
}

func TestWriteTreeSafePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "koko-test-writer-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	v := vfs.New()
	v.WriteString("package.json", `{"name": "test-app"}`)
	v.WriteString("src/index.ts", `console.log("hello");`)

	targetDir := filepath.Join(tmpDir, "test-app")
	if err := WriteTree(v, targetDir); err != nil {
		t.Fatalf("WriteTree failed: %v", err)
	}

	// Verify files written to disk
	if _, err := os.Stat(filepath.Join(targetDir, "package.json")); os.IsNotExist(err) {
		t.Errorf("Expected package.json on disk")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "src", "index.ts")); os.IsNotExist(err) {
		t.Errorf("Expected src/index.ts on disk")
	}
}
