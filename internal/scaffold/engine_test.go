package scaffold

import (
	"os"
	"path/filepath"
	"strings"
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

func TestEngineGenerateVFSFullstackSelf(t *testing.T) {
	cfg := ScaffoldConfig{
		ProjectName:    "my-next-self-app",
		Frontend:       "nextjs",
		Backend:        "self",
		PackageManager: "pnpm",
		Database:       "postgres",
		ORM:            "drizzle",
		Auth:           "better-auth",
	}

	virtualFS, err := GenerateVFS(cfg)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	// Should generate apps/web
	if !virtualFS.Exists("apps/web/package.json") {
		t.Errorf("Expected apps/web/package.json to exist in VFS")
	}

	// Should NOT generate apps/api because backend is self
	if virtualFS.Exists("apps/api/package.json") {
		t.Errorf("Expected apps/api to not exist when backend is 'self'")
	}

	// Should generate packages/db
	if !virtualFS.Exists("packages/db/drizzle.config.ts") {
		t.Errorf("Expected packages/db/drizzle.config.ts in VFS")
	}

	// Should generate shared packages/auth
	if !virtualFS.Exists("packages/auth/src/index.ts") {
		t.Errorf("Expected packages/auth/src/index.ts in VFS")
	}

	// Should generate Next.js auth endpoint route and auth-client
	if !virtualFS.Exists("apps/web/app/api/auth/[...all]/route.ts") {
		t.Errorf("Expected apps/web/app/api/auth/[...all]/route.ts in VFS")
	}
	if !virtualFS.Exists("apps/web/src/lib/auth-client.ts") {
		t.Errorf("Expected apps/web/src/lib/auth-client.ts in VFS")
	}
}

func TestEngineGenerateVFSAstro(t *testing.T) {
	cfg := ScaffoldConfig{
		ProjectName:    "my-astro-app",
		Frontend:       "astro",
		Backend:        "self",
		PackageManager: "pnpm",
		Auth:           "better-auth",
	}

	virtualFS, err := GenerateVFS(cfg)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	if !virtualFS.Exists("apps/web/astro.config.mjs") {
		t.Errorf("Expected apps/web/astro.config.mjs in VFS")
	}

	if !virtualFS.Exists("apps/web/src/pages/api/auth/[...all].ts") {
		t.Errorf("Expected apps/web/src/pages/api/auth/[...all].ts in VFS")
	}

	if !virtualFS.Exists("apps/web/src/lib/auth-client.ts") {
		t.Errorf("Expected apps/web/src/lib/auth-client.ts in VFS")
	}
}

func TestEngineGenerateVFSSvelte(t *testing.T) {
	cfg := ScaffoldConfig{
		ProjectName:    "my-svelte-app",
		Frontend:       "svelte",
		Backend:        "self",
		PackageManager: "pnpm",
		Auth:           "better-auth",
	}

	virtualFS, err := GenerateVFS(cfg)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	if !virtualFS.Exists("apps/web/src/hooks.server.ts") {
		t.Errorf("Expected apps/web/src/hooks.server.ts in VFS")
	}

	if !virtualFS.Exists("apps/web/src/lib/auth-client.ts") {
		t.Errorf("Expected apps/web/src/lib/auth-client.ts in VFS")
	}
}

func TestEngineGenerateVFSNuxt(t *testing.T) {
	cfg := ScaffoldConfig{
		ProjectName:    "my-nuxt-app",
		Frontend:       "nuxt",
		Backend:        "self",
		PackageManager: "pnpm",
		Auth:           "better-auth",
	}

	virtualFS, err := GenerateVFS(cfg)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	if !virtualFS.Exists("apps/web/server/api/auth/[...all].ts") {
		t.Errorf("Expected apps/web/server/api/auth/[...all].ts in VFS")
	}

	if !virtualFS.Exists("apps/web/utils/auth-client.ts") {
		t.Errorf("Expected apps/web/utils/auth-client.ts in VFS")
	}
}

func TestEngineGenerateVFSNative(t *testing.T) {
	cfg := ScaffoldConfig{
		ProjectName:    "my-native-app",
		Frontend:       "native",
		Backend:        "express",
		PackageManager: "pnpm",
	}

	virtualFS, err := GenerateVFS(cfg)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	if !virtualFS.Exists("apps/web/app.json") {
		t.Errorf("Expected apps/web/app.json in VFS for React Native")
	}
}

func TestEngineGenerateVFSReactWithAuth(t *testing.T) {
	cfg := ScaffoldConfig{
		ProjectName:    "my-react-app",
		Frontend:       "react",
		Backend:        "express",
		PackageManager: "pnpm",
		Auth:           "better-auth",
	}

	virtualFS, err := GenerateVFS(cfg)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	if !virtualFS.Exists("apps/web/src/lib/auth-client.ts") {
		t.Errorf("Expected apps/web/src/lib/auth-client.ts in VFS for React SPA")
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

func TestEngineGenerateVFSGitignoreTemplating(t *testing.T) {
	// Case 1: Next.js + FastAPI + PostgreSQL
	cfgPython := ScaffoldConfig{
		ProjectName: "py-app",
		Frontend:    "nextjs",
		Backend:     "fastapi",
		Database:    "postgres",
	}

	vfsPy, err := GenerateVFS(cfgPython)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	giPy, ok := vfsPy.ReadFile(".gitignore")
	if !ok {
		t.Fatalf(".gitignore was not generated in VFS")
	}
	giPyStr := string(giPy)

	if !strings.Contains(giPyStr, ".next/") {
		t.Errorf("Expected .next/ in gitignore for Next.js")
	}
	if !strings.Contains(giPyStr, "__pycache__/") || !strings.Contains(giPyStr, ".venv/") {
		t.Errorf("Expected Python entries in gitignore for FastAPI")
	}
	if strings.Contains(giPyStr, "*.exe") || strings.Contains(giPyStr, ".svelte-kit/") {
		t.Errorf("Did not expect Go or Svelte entries in Next.js + FastAPI gitignore")
	}

	// Case 2: Svelte + Go Chi + SQLite
	cfgGo := ScaffoldConfig{
		ProjectName: "go-app",
		Frontend:    "svelte",
		Backend:     "go_chi",
		Database:    "sqlite",
	}

	vfsGo, err := GenerateVFS(cfgGo)
	if err != nil {
		t.Fatalf("GenerateVFS failed: %v", err)
	}

	giGo, ok := vfsGo.ReadFile(".gitignore")
	if !ok {
		t.Fatalf(".gitignore was not generated in VFS")
	}
	giGoStr := string(giGo)

	if !strings.Contains(giGoStr, ".svelte-kit/") {
		t.Errorf("Expected .svelte-kit/ in gitignore for Svelte")
	}
	if !strings.Contains(giGoStr, "bin/") || !strings.Contains(giGoStr, "*.exe") {
		t.Errorf("Expected Go entries in gitignore for Go Chi")
	}
	if !strings.Contains(giGoStr, "*.sqlite") {
		t.Errorf("Expected SQLite entries in gitignore")
	}
	if strings.Contains(giGoStr, ".next/") || strings.Contains(giGoStr, "__pycache__/") {
		t.Errorf("Did not expect Next.js or Python entries in Svelte + Go Chi gitignore")
	}
}
