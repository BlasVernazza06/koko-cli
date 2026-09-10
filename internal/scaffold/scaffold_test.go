package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManualTemplatesCompleteness(t *testing.T) {
	// 1. Check root template files
	requiredRootFiles := []string{
		"manual/root/package.json",
		"manual/root/turbo.json",
		"manual/root/pnpm-workspace.yaml",
		"manual/root/.gitignore",
		"manual/root/.npmrc",
		"manual/root/README.md",
		"manual/root/apps/.gitkeep",
		"manual/root/packages/.gitkeep",
	}
	for _, f := range requiredRootFiles {
		if _, err := templateFs.ReadFile(f); err != nil {
			t.Errorf("Missing root template file: %s (err: %v)", f, err)
		}
	}

	// 2. Check frontend templates
	frontends := []struct {
		name  string
		files []string
	}{
		{"nextjs", []string{"package.json", "tsconfig.json", "next.config.mjs", "app/layout.tsx", "app/page.tsx", "app/globals.css"}},
		{"react", []string{"package.json", "tsconfig.json", "vite.config.ts", "index.html", "src/main.tsx", "src/App.tsx"}},
		{"nuxt", []string{"package.json", "tsconfig.json", "nuxt.config.ts", "app.vue"}},
		{"svelte", []string{"package.json", "svelte.config.js", "vite.config.ts", "src/routes/+page.svelte", "src/routes/+layout.svelte"}},
	}
	for _, fe := range frontends {
		for _, f := range fe.files {
			fullPath := "manual/frontend/" + fe.name + "/" + f
			if _, err := templateFs.ReadFile(fullPath); err != nil {
				t.Errorf("Missing frontend template file for %s: %s (err: %v)", fe.name, fullPath, err)
			}
		}
	}

	// 3. Check backend templates
	backends := []struct {
		name  string
		files []string
	}{
		{"express", []string{"package.json", "tsconfig.json", "src/index.ts"}},
		{"fastapi", []string{"requirements.txt", "main.py", "app/routers/health.py", "app/routers/todos.py"}},
		{"go_chi", []string{"go.mod.tmpl", "main.go.tmpl", "handlers/health.go.tmpl", "handlers/todos.go.tmpl"}},
		{"nestjs", []string{"package.json", "nest-cli.json", "tsconfig.json", "src/main.ts", "src/app.module.ts", "src/app.controller.ts", "src/app.service.ts"}},
		{"hono", []string{"package.json", "tsconfig.json", "src/index.ts"}},
		{"spring_boot", []string{"pom.xml", "src/main/resources/application.yml", "src/main/java/com/koko/app/Application.java"}},
	}
	for _, be := range backends {
		for _, f := range be.files {
			fullPath := "manual/backend/" + be.name + "/" + f
			if _, err := templateFs.ReadFile(fullPath); err != nil {
				t.Errorf("Missing backend template file for %s: %s (err: %v)", be.name, fullPath, err)
			}
		}
	}

	// 4. Check DB and ORM templates
	dbs := []string{
		"manual/db/drizzle/drizzle.config.ts",
		"manual/db/drizzle/package.json",
		"manual/db/drizzle/postgres/schema.ts",
		"manual/db/drizzle/postgres/db.ts",
		"manual/db/drizzle/mysql/schema.ts",
		"manual/db/drizzle/mysql/db.ts",
		"manual/db/drizzle/sqlite/schema.ts",
		"manual/db/drizzle/sqlite/db.ts",
		"manual/db/prisma/package.json",
		"manual/db/prisma/postgres/schema.prisma",
		"manual/db/prisma/mysql/schema.prisma",
		"manual/db/prisma/mongodb/schema.prisma",
		"manual/db/prisma/sqlite/schema.prisma",
		"manual/db/mongoose/mongodb/package.json",
		"manual/db/mongoose/mongodb/src/db.ts",
		"manual/db/mongoose/mongodb/src/models/user.ts",
		"manual/db/sqlalchemy/postgres/database.py",
		"manual/db/sqlalchemy/postgres/models.py",
		"manual/db/sqlalchemy/mysql/database.py",
		"manual/db/sqlalchemy/sqlite/database.py",
		"manual/db/gorm/postgres/database.go.tmpl",
		"manual/db/gorm/postgres/models.go.tmpl",
		"manual/db/gorm/mysql/database.go.tmpl",
		"manual/db/gorm/sqlite/database.go.tmpl",
	}
	for _, dbFile := range dbs {
		if _, err := templateFs.ReadFile(dbFile); err != nil {
			t.Errorf("Missing DB/ORM template file: %s (err: %v)", dbFile, err)
		}
	}

	// 5. Check packages, docker, github
	shared := []string{
		"manual/root/packages/typescript-config/package.json",
		"manual/root/packages/typescript-config/base.json",
		"manual/root/packages/typescript-config/nextjs.json",
		"manual/root/packages/typescript-config/react-library.json",
		"manual/root/packages/eslint-config/package.json",
		"manual/root/packages/eslint-config/base.js",
		"manual/docker/postgres/docker-compose.yml",
		"manual/docker/mysql/docker-compose.yml",
		"manual/docker/mongoose/docker-compose.yml",
		"manual/github/ci.yml",
		"manual/api/trpc/packages/api/package.json",
		"manual/api/trpc/packages/api/src/index.ts",
		"manual/api/trpc/packages/api/src/root.ts",
		"manual/api/trpc/packages/api/src/trpc.ts",
		"manual/api/orpc/packages/api/package.json",
		"manual/api/orpc/packages/api/src/index.ts",
		"manual/api/orpc/packages/api/src/root.ts",
		"manual/api/orpc/packages/api/src/orpc.ts",
	}
	for _, sh := range shared {
		if _, err := templateFs.ReadFile(sh); err != nil {
			t.Errorf("Missing shared template file: %s (err: %v)", sh, err)
		}
	}
}


func TestScaffoldManualNextjsExpressPostgresDrizzle(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-next-express"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName:    projectName,
		Frontend:       "nextjs",
		Backend:        "express",
		PackageManager: "pnpm",
		Database:       "postgres",
		ORM:            "drizzle",
		Addons:         "docker_cicd",
		InitGit:        false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold failed: %v", err)
	}

	// Verify apps/.gitkeep and packages/.gitkeep are not present
	for _, keep := range []string{filepath.Join(targetDir, "apps", ".gitkeep"), filepath.Join(targetDir, "packages", ".gitkeep")} {
		if _, err := os.Stat(keep); !os.IsNotExist(err) {
			t.Errorf("Expected .gitkeep to not exist, but it exists: %s", keep)
		}
	}


	// Verify Root Monorepo files
	expectedRootFiles := []string{
		"package.json",
		"turbo.json",
		"pnpm-workspace.yaml",
		".gitignore",
		".npmrc",
		"README.md",
		"docker-compose.yml",
		filepath.Join(".github", "workflows", "ci.yml"),
	}
	for _, f := range expectedRootFiles {
		p := filepath.Join(targetDir, f)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("Expected root file %s not found", f)
		}
	}

	// Verify Frontend in apps/web
	expectedFrontendFiles := []string{
		filepath.Join("apps", "web", "package.json"),
		filepath.Join("apps", "web", "app", "layout.tsx"),
		filepath.Join("apps", "web", "app", "page.tsx"),
		filepath.Join("apps", "web", "next.config.mjs"),
	}
	for _, f := range expectedFrontendFiles {
		p := filepath.Join(targetDir, f)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("Expected frontend file %s not found", f)
		}
	}

	// Verify Backend in apps/api
	expectedBackendFiles := []string{
		filepath.Join("apps", "api", "package.json"),
		filepath.Join("apps", "api", "src", "index.ts"),
		filepath.Join("apps", "api", "tsconfig.json"),
	}
	for _, f := range expectedBackendFiles {
		p := filepath.Join(targetDir, f)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("Expected backend file %s not found", f)
		}
	}

	// Verify Drizzle in packages/db
	expectedDbFiles := []string{
		filepath.Join("packages", "db", "drizzle.config.ts"),
		filepath.Join("packages", "db", "package.json"),
		filepath.Join("packages", "db", "src", "schema.ts"),
		filepath.Join("packages", "db", "src", "db.ts"),
	}
	for _, f := range expectedDbFiles {
		p := filepath.Join(targetDir, f)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("Expected db file %s not found", f)
		}
	}

	// Verify variable interpolation in root package.json
	pkgData, err := os.ReadFile(filepath.Join(targetDir, "package.json"))
	if err != nil {
		t.Fatalf("Failed to read root package.json: %v", err)
	}
	if !strings.Contains(string(pkgData), `"name": "test-next-express"`) {
		t.Errorf("Expected interpolated project name in package.json, got: %s", string(pkgData))
	}
}

func TestScaffoldManualReactGoChiMysqlGorm(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-react-gochi"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName:    projectName,
		Frontend:       "react",
		Backend:        "go_chi",
		PackageManager: "pnpm",
		Database:       "mysql",
		ORM:            "gorm",
		Addons:         "docker",
		InitGit:        false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold failed: %v", err)
	}

	// Verify React in apps/web
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "web", "vite.config.ts")); os.IsNotExist(err) {
		t.Errorf("apps/web/vite.config.ts not found")
	}

	// Verify Go Chi in apps/api
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "api", "main.go")); os.IsNotExist(err) {
		t.Errorf("apps/api/main.go not found")
	}

	// Verify GORM in apps/api/db
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "api", "db", "database.go")); os.IsNotExist(err) {
		t.Errorf("apps/api/db/database.go not found")
	}

	// Verify Docker Compose with MySQL
	dcData, err := os.ReadFile(filepath.Join(targetDir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("docker-compose.yml not found: %v", err)
	}
	if !strings.Contains(string(dcData), "mysql:") {
		t.Errorf("Expected mysql in docker-compose.yml, got: %s", string(dcData))
	}
}

func TestScaffoldManualNuxtFastapiSqliteSqlalchemy(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-nuxt-fastapi"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName:    projectName,
		Frontend:       "nuxt",
		Backend:        "fastapi",
		PackageManager: "npm",
		Database:       "sqlite",
		ORM:            "sqlalchemy",
		Addons:         "none",
		InitGit:        false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold failed: %v", err)
	}

	// Verify Nuxt in apps/web
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "web", "nuxt.config.ts")); os.IsNotExist(err) {
		t.Errorf("apps/web/nuxt.config.ts not found")
	}

	// Verify FastAPI in apps/api
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "api", "main.py")); os.IsNotExist(err) {
		t.Errorf("apps/api/main.py not found")
	}

	// Verify SQLAlchemy in apps/api/app/db
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "api", "app", "db", "database.py")); os.IsNotExist(err) {
		t.Errorf("apps/api/app/db/database.py not found")
	}
}

func TestScaffoldManualSvelteNestjsMongoMongoose(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-svelte-nestjs"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName:    projectName,
		Frontend:       "svelte",
		Backend:        "nestjs",
		PackageManager: "bun",
		Database:       "mongodb",
		ORM:            "mongoose",
		Addons:         "docker_cicd",
		InitGit:        false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold failed: %v", err)
	}

	// Verify Svelte in apps/web
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "web", "svelte.config.js")); os.IsNotExist(err) {
		t.Errorf("apps/web/svelte.config.js not found")
	}

	// Verify NestJS in apps/api
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "api", "src", "main.ts")); os.IsNotExist(err) {
		t.Errorf("apps/api/src/main.ts not found")
	}

	// Verify Mongoose in packages/db
	if _, err := os.Stat(filepath.Join(targetDir, "packages", "db", "src", "db.ts")); os.IsNotExist(err) {
		t.Errorf("packages/db/src/db.ts not found")
	}

	// Verify Docker Compose with MongoDB
	dcData, err := os.ReadFile(filepath.Join(targetDir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("docker-compose.yml not found: %v", err)
	}
	if !strings.Contains(string(dcData), "mongodb:") {
		t.Errorf("Expected mongodb in docker-compose.yml, got: %s", string(dcData))
	}
}

func TestScaffoldRecipeSaaS(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-recipe-saas"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "saas",
		InitGit:     false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for recipe saas failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "package.json")); os.IsNotExist(err) {
		t.Errorf("Root package.json not found for recipe saas")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "web", "package.json")); os.IsNotExist(err) {
		t.Errorf("apps/web/package.json not found for recipe saas")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "web", "lib", "email.ts")); os.IsNotExist(err) {
		t.Errorf("apps/web/lib/email.ts not found for recipe saas")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "web", "components", "billing", "pricing-cards.tsx")); os.IsNotExist(err) {
		t.Errorf("pricing-cards.tsx not found for recipe saas")
	}
}

func TestScaffoldRecipeJavaSpring(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-recipe-java-spring"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "java-spring",
		InitGit:     false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for recipe java-spring failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "package.json")); os.IsNotExist(err) {
		t.Errorf("Root package.json not found for recipe java-spring")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "backend", "pom.xml")); os.IsNotExist(err) {
		t.Errorf("apps/backend/pom.xml not found for recipe java-spring")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "frontend", "src", "App.tsx")); os.IsNotExist(err) {
		t.Errorf("apps/frontend/src/App.tsx not found for recipe java-spring")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "docker-compose.yml")); os.IsNotExist(err) {
		t.Errorf("docker-compose.yml not found for recipe java-spring")
	}
}

func TestScaffoldRecipeEnterpriseNestJS(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-recipe-nestjs"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "enterprise-nestjs",
		InitGit:     false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for recipe enterprise-nestjs failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "apps", "web", "package.json")); os.IsNotExist(err) {
		t.Errorf("apps/web/package.json not found for recipe enterprise-nestjs")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "api", "src", "main.ts")); os.IsNotExist(err) {
		t.Errorf("apps/api/src/main.ts not found for recipe enterprise-nestjs")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "packages", "db", "prisma", "schema.prisma")); os.IsNotExist(err) {
		t.Errorf("schema.prisma not found for recipe enterprise-nestjs")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "docker-compose.yml")); os.IsNotExist(err) {
		t.Errorf("docker-compose.yml not found for recipe enterprise-nestjs")
	}
}

func TestScaffoldRecipeMERN(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-recipe-mern"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "mern",
		InitGit:     false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for recipe mern failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "apps", "backend", "src", "models", "item.ts")); os.IsNotExist(err) {
		t.Errorf("models/item.ts not found for recipe mern")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "frontend", "src", "App.tsx")); os.IsNotExist(err) {
		t.Errorf("apps/frontend/src/App.tsx not found for recipe mern")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "docker-compose.yml")); os.IsNotExist(err) {
		t.Errorf("docker-compose.yml not found for recipe mern")
	}
}

func TestScaffoldRecipePERN(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-recipe-pern"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "pern",
		InitGit:     false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for recipe pern failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "apps", "backend", "src", "index.ts")); os.IsNotExist(err) {
		t.Errorf("apps/backend/src/index.ts not found for recipe pern")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "packages", "db", "schema.prisma")); os.IsNotExist(err) {
		t.Errorf("packages/db/schema.prisma not found for recipe pern")
	}
}

func TestScaffoldRecipeFastAPIReact(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-recipe-fastapi"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "python-fastapi",
		InitGit:     false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for recipe python-fastapi failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "apps", "backend", "main.py")); os.IsNotExist(err) {
		t.Errorf("apps/backend/main.py not found for recipe python-fastapi")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "backend", "app", "schemas", "item.py")); os.IsNotExist(err) {
		t.Errorf("app/schemas/item.py not found for recipe python-fastapi")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "frontend", "src", "App.tsx")); os.IsNotExist(err) {
		t.Errorf("apps/frontend/src/App.tsx not found for recipe python-fastapi")
	}
}

func TestScaffoldRecipeMobileExpo(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-recipe-expo"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "mobile-expo",
		InitGit:     false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for recipe mobile-expo failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "apps", "mobile", "App.tsx")); os.IsNotExist(err) {
		t.Errorf("apps/mobile/App.tsx not found for recipe mobile-expo")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "apps", "api", "src", "index.ts")); os.IsNotExist(err) {
		t.Errorf("apps/api/src/index.ts not found for recipe mobile-expo")
	}
	if _, err := os.Stat(filepath.Join(targetDir, "packages", "db", "prisma", "schema.prisma")); os.IsNotExist(err) {
		t.Errorf("schema.prisma not found for recipe mobile-expo")
	}
}

func TestScaffoldManualAddonShadcnUI(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-manual-shadcn"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName:    projectName,
		Frontend:       "nextjs",
		Backend:        "self",
		API:            "trpc",
		PackageManager: "pnpm",
		Database:       "postgres",
		ORM:            "drizzle",
		Auth:           "better-auth",
		Addons:         "shadcn,lucide,svgl,motion,zod,docker",
		InitGit:        false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for manual stack with shadcn failed: %v", err)
	}

	// 1. Verify packages/ui structure
	requiredUIFiles := []string{
		filepath.Join(targetDir, "packages", "ui", "package.json"),
		filepath.Join(targetDir, "packages", "ui", "components.json"),
		filepath.Join(targetDir, "packages", "ui", "tsconfig.json"),
		filepath.Join(targetDir, "packages", "ui", "src", "index.ts"),
		filepath.Join(targetDir, "packages", "ui", "src", "lib", "utils.ts"),
		filepath.Join(targetDir, "packages", "ui", "src", "components", "ui", "button.tsx"),
		filepath.Join(targetDir, "packages", "ui", "src", "hooks", "use-mobile.tsx"),
	}

	for _, f := range requiredUIFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("Expected shadcn UI file not found: %s", f)
		}
	}

	// 2. Verify apps/web package.json has dependencies
	webPkgBytes, err := os.ReadFile(filepath.Join(targetDir, "apps", "web", "package.json"))
	if err != nil {
		t.Fatalf("Could not read apps/web/package.json: %v", err)
	}
	webPkgStr := string(webPkgBytes)

	if !strings.Contains(webPkgStr, "@test-manual-shadcn/ui") {
		t.Errorf("Expected apps/web/package.json to depend on @test-manual-shadcn/ui, got: %s", webPkgStr)
	}
	if !strings.Contains(webPkgStr, "lucide-react") {
		t.Errorf("Expected apps/web/package.json to have lucide-react, got: %s", webPkgStr)
	}
	if !strings.Contains(webPkgStr, "@svgl/react") {
		t.Errorf("Expected apps/web/package.json to have @svgl/react, got: %s", webPkgStr)
	}
	if !strings.Contains(webPkgStr, "framer-motion") {
		t.Errorf("Expected apps/web/package.json to have framer-motion, got: %s", webPkgStr)
	}
	if !strings.Contains(webPkgStr, "zod") {
		t.Errorf("Expected apps/web/package.json to have zod, got: %s", webPkgStr)
	}
}

func TestScaffoldManualTRPC(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-trpc-app"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName:    projectName,
		Frontend:       "nextjs",
		Backend:        "self",
		API:            "trpc",
		PackageManager: "pnpm",
		Database:       "postgres",
		ORM:            "drizzle",
		Auth:           "better-auth",
		Addons:         "docker",
		InitGit:        false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for tRPC stack failed: %v", err)
	}

	// 1. Verify packages/api files
	requiredAPIFiles := []string{
		filepath.Join(targetDir, "packages", "api", "package.json"),
		filepath.Join(targetDir, "packages", "api", "tsconfig.json"),
		filepath.Join(targetDir, "packages", "api", "src", "index.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "trpc.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "root.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "context.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "routers", "health.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "routers", "user.ts"),
	}
	for _, f := range requiredAPIFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("Expected tRPC file not found: %s", f)
		}
	}

	// 2. Verify Next.js integration
	routeFile := filepath.Join(targetDir, "apps", "web", "app", "api", "trpc", "[trpc]", "route.ts")
	if _, err := os.Stat(routeFile); os.IsNotExist(err) {
		t.Errorf("Expected Next.js tRPC route handler not found: %s", routeFile)
	}

	// 3. Verify apps/web package.json has @repo/api and @trpc/* dependencies
	webPkgBytes, err := os.ReadFile(filepath.Join(targetDir, "apps", "web", "package.json"))
	if err != nil {
		t.Fatalf("Could not read apps/web/package.json: %v", err)
	}
	webPkgStr := string(webPkgBytes)
	if !strings.Contains(webPkgStr, "@test-trpc-app/api") {
		t.Errorf("Expected apps/web/package.json to depend on @test-trpc-app/api, got: %s", webPkgStr)
	}
	if !strings.Contains(webPkgStr, "@trpc/client") {
		t.Errorf("Expected apps/web/package.json to have @trpc/client, got: %s", webPkgStr)
	}
	if !strings.Contains(webPkgStr, "@tanstack/react-query") {
		t.Errorf("Expected apps/web/package.json to have @tanstack/react-query, got: %s", webPkgStr)
	}
}

func TestScaffoldManualORPC(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-orpc-app"
	targetDir := filepath.Join(tmpDir, projectName)

	cfg := ScaffoldConfig{
		ProjectName:    projectName,
		Frontend:       "react",
		Backend:        "express",
		API:            "orpc",
		PackageManager: "pnpm",
		Database:       "postgres",
		ORM:            "drizzle",
		Auth:           "none",
		Addons:         "docker",
		InitGit:        false,
	}

	err := RunScaffold(targetDir, cfg)
	if err != nil {
		t.Fatalf("RunScaffold for oRPC stack failed: %v", err)
	}

	// 1. Verify packages/api files
	requiredAPIFiles := []string{
		filepath.Join(targetDir, "packages", "api", "package.json"),
		filepath.Join(targetDir, "packages", "api", "tsconfig.json"),
		filepath.Join(targetDir, "packages", "api", "src", "index.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "orpc.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "root.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "openapi.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "context.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "routers", "health.ts"),
		filepath.Join(targetDir, "packages", "api", "src", "routers", "user.ts"),
	}
	for _, f := range requiredAPIFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("Expected oRPC file not found: %s", f)
		}
	}

	// 2. Verify Express server integration
	expressHandler := filepath.Join(targetDir, "apps", "api", "src", "orpc.ts")
	if _, err := os.Stat(expressHandler); os.IsNotExist(err) {
		t.Errorf("Expected Express oRPC handler not found: %s", expressHandler)
	}

	// 3. Verify apps/api package.json has dependencies
	apiPkgBytes, err := os.ReadFile(filepath.Join(targetDir, "apps", "api", "package.json"))
	if err != nil {
		t.Fatalf("Could not read apps/api/package.json: %v", err)
	}
	apiPkgStr := string(apiPkgBytes)
	if !strings.Contains(apiPkgStr, "@test-orpc-app/api") {
		t.Errorf("Expected apps/api/package.json to depend on @test-orpc-app/api, got: %s", apiPkgStr)
	}
	if !strings.Contains(apiPkgStr, "@orpc/server") {
		t.Errorf("Expected apps/api/package.json to have @orpc/server, got: %s", apiPkgStr)
	}
}

