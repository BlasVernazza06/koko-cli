package processors

import (
	"strings"
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

func TestProcessPackageJSONs(t *testing.T) {
	v := vfs.New()
	v.WriteString("package.json", `{"name": "placeholder", "scripts": {}}`)
	v.WriteString("apps/web/package.json", `{"name": "placeholder-web"}`)
	v.WriteString("packages/db/package.json", `{"name": "placeholder-db"}`)

	cfg := ProcessConfig{
		ProjectName:    "acme-saas",
		PackageManager: "pnpm",
		Database:       "postgres",
		ORM:            "drizzle",
	}

	if err := ProcessPackageJSONs(v, cfg); err != nil {
		t.Fatalf("ProcessPackageJSONs failed: %v", err)
	}

	var rootPkg map[string]interface{}
	if err := v.ReadJSON("package.json", &rootPkg); err != nil {
		t.Fatalf("Failed to read root package.json: %v", err)
	}

	if rootPkg["name"] != "acme-saas" {
		t.Errorf("Expected root name 'acme-saas', got '%v'", rootPkg["name"])
	}

	scripts := rootPkg["scripts"].(map[string]interface{})
	if scripts["dev"] != "pnpm -r dev" {
		t.Errorf("Expected dev script 'pnpm -r dev', got '%v'", scripts["dev"])
	}
	if scripts["db:push"] != "pnpm --filter @acme-saas/db db:push" {
		t.Errorf("Expected db:push script with filter, got '%v'", scripts["db:push"])
	}

	var webPkg map[string]interface{}
	_ = v.ReadJSON("apps/web/package.json", &webPkg)
	if webPkg["name"] != "@acme-saas/web" {
		t.Errorf("Expected web package name '@acme-saas/web', got '%v'", webPkg["name"])
	}
}

func TestProcessEnvVariables(t *testing.T) {
	v := vfs.New()
	cfg := ProcessConfig{
		ProjectName: "shop-app",
		Backend:     "express",
		Frontend:    "nextjs",
		Database:    "postgres",
		Auth:        "better-auth",
	}

	if err := ProcessEnvVariables(v, cfg); err != nil {
		t.Fatalf("ProcessEnvVariables failed: %v", err)
	}

	envContent, ok := v.ReadString(".env")
	if !ok {
		t.Fatalf("Expected .env file to be created")
	}

	if !strings.Contains(envContent, "DATABASE_URL=\"postgresql://postgres:password@localhost:5432/shop-app?schema=public\"") {
		t.Errorf("Expected Postgres connection string in .env, got: %s", envContent)
	}
	if !strings.Contains(envContent, "BETTER_AUTH_SECRET") {
		t.Errorf("Expected Better Auth secret in .env, got: %s", envContent)
	}
	if !strings.Contains(envContent, "NEXT_PUBLIC_API_URL=http://localhost:4000") {
		t.Errorf("Expected Next.js API url in .env, got: %s", envContent)
	}
}
