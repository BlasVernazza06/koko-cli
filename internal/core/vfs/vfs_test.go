package vfs

import (
	"reflect"
	"testing"
)

type samplePkg struct {
	Name    string            `json:"name"`
	Scripts map[string]string `json:"scripts"`
}

func TestVFSBasicOperations(t *testing.T) {
	fs := New()

	fs.WriteString("package.json", `{"name": "test-app"}`)
	fs.WriteFile("src/index.ts", []byte("console.log('hello');"))

	if fs.FileCount() != 2 {
		t.Fatalf("Expected 2 files, got %d", fs.FileCount())
	}

	if !fs.Exists("package.json") {
		t.Errorf("Expected package.json to exist")
	}

	if !fs.Exists("./src/index.ts") {
		t.Errorf("Expected normalized path to exist")
	}

	content, ok := fs.ReadString("package.json")
	if !ok || content != `{"name": "test-app"}` {
		t.Errorf("Unexpected content: %s", content)
	}

	files := fs.ListFiles()
	expected := []string{"package.json", "src/index.ts"}
	if !reflect.DeepEqual(files, expected) {
		t.Errorf("ListFiles mismatch. Expected %v, got %v", expected, files)
	}

	fs.DeleteFile("src/index.ts")
	if fs.Exists("src/index.ts") {
		t.Errorf("Expected src/index.ts to be deleted")
	}
}

func TestVFSJSONOperations(t *testing.T) {
	fs := New()

	pkg := samplePkg{
		Name: "@koko/my-app",
		Scripts: map[string]string{
			"dev":   "pnpm dev",
			"build": "pnpm build",
		},
	}

	if err := fs.WriteJSON("package.json", pkg, "  "); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var readPkg samplePkg
	if err := fs.ReadJSON("package.json", &readPkg); err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if readPkg.Name != pkg.Name {
		t.Errorf("Expected name %s, got %s", pkg.Name, readPkg.Name)
	}
	if readPkg.Scripts["dev"] != "pnpm dev" {
		t.Errorf("Expected dev script 'pnpm dev', got %s", readPkg.Scripts["dev"])
	}
}
