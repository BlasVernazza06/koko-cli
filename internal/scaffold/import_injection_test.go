package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

func TestAutoWiringNextJSTRPCClerkShadcn(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := types.ScaffoldConfig{
		ProjectName:    "test-nextjs-clerk-trpc",
		Frontend:       "nextjs",
		Backend:        "none",
		API:            "trpc",
		Auth:           "clerk",
		Database:       "postgres",
		ORM:            "prisma",
		Addons:         "shadcn,motion,lucide",
		PackageManager: "pnpm",
	}

	if err := RunScaffold(tmpDir, cfg); err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}

	// 1. Validate layout.tsx
	layoutPath := filepath.Join(tmpDir, "apps", "web", "app", "layout.tsx")
	layoutContent, err := os.ReadFile(layoutPath)
	if err != nil {
		t.Fatalf("Failed to read layout.tsx: %v", err)
	}
	layoutStr := string(layoutContent)

	if !strings.Contains(layoutStr, "<ClerkProvider>") {
		t.Errorf("layout.tsx missing <ClerkProvider>")
	}
	if !strings.Contains(layoutStr, "<TRPCReactProvider>") {
		t.Errorf("layout.tsx missing <TRPCReactProvider>")
	}
	if !strings.Contains(layoutStr, "<ThemeProvider") {
		t.Errorf("layout.tsx missing <ThemeProvider")
	}

	// 2. Validate page.tsx
	pagePath := filepath.Join(tmpDir, "apps", "web", "app", "page.tsx")
	pageContent, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("Failed to read page.tsx: %v", err)
	}
	pageStr := string(pageContent)

	if !strings.Contains(pageStr, "SignedIn") || !strings.Contains(pageStr, "UserButton") {
		t.Errorf("page.tsx missing Clerk SignedIn/UserButton components")
	}
	if !strings.Contains(pageStr, "framer-motion") {
		t.Errorf("page.tsx missing framer-motion import")
	}

	// 3. Validate packages/api/src/context.ts
	ctxPath := filepath.Join(tmpDir, "packages", "api", "src", "context.ts")
	if ctxContent, err := os.ReadFile(ctxPath); err == nil {
		ctxStr := string(ctxContent)
		if !strings.Contains(ctxStr, "@test-nextjs-clerk-trpc/db") {
			t.Errorf("context.ts missing db import")
		}
		if !strings.Contains(ctxStr, "db,") {
			t.Errorf("context.ts missing db injection in context return")
		}
	} else {
		t.Fatalf("Failed to read context.ts: %v", err)
	}

	validateProjectTree(t, tmpDir, cfg.ProjectName)
}

func TestAutoWiringReactORPCBetterAuth(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := types.ScaffoldConfig{
		ProjectName:    "test-react-orpc-auth",
		Frontend:       "react",
		Backend:        "hono",
		API:            "orpc",
		Auth:           "better-auth",
		Database:       "postgres",
		ORM:            "drizzle",
		Addons:         "lucide,motion",
		PackageManager: "pnpm",
	}

	if err := RunScaffold(tmpDir, cfg); err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}

	// 1. Validate main.tsx
	mainPath := filepath.Join(tmpDir, "apps", "web", "src", "main.tsx")
	mainContent, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("Failed to read main.tsx: %v", err)
	}
	mainStr := string(mainContent)

	if !strings.Contains(mainStr, "<ORPCReactProvider>") {
		t.Errorf("main.tsx missing <ORPCReactProvider>")
	}

	// 2. Validate App.tsx
	appPath := filepath.Join(tmpDir, "apps", "web", "src", "App.tsx")
	appContent, err := os.ReadFile(appPath)
	if err != nil {
		t.Fatalf("Failed to read App.tsx: %v", err)
	}
	appStr := string(appContent)

	if !strings.Contains(appStr, "authClient") {
		t.Errorf("App.tsx missing authClient import")
	}
	if !strings.Contains(appStr, "Sparkles") {
		t.Errorf("App.tsx missing Sparkles lucide icon")
	}

	// 3. Validate Hono index.ts
	honoPath := filepath.Join(tmpDir, "apps", "api", "src", "index.ts")
	honoContent, err := os.ReadFile(honoPath)
	if err != nil {
		t.Fatalf("Failed to read Hono index.ts: %v", err)
	}
	honoStr := string(honoContent)

	if !strings.Contains(honoStr, "handleORPC") {
		t.Errorf("Hono index.ts missing handleORPC")
	}
	if !strings.Contains(honoStr, "auth.handler(c.req.raw)") {
		t.Errorf("Hono index.ts missing Better-Auth handler mounting")
	}

	validateProjectTree(t, tmpDir, cfg.ProjectName)
}

func TestAutoWiringExpoClerkLucide(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := types.ScaffoldConfig{
		ProjectName:    "test-expo-clerk",
		Frontend:       "native",
		Backend:        "none",
		Auth:           "clerk",
		Addons:         "lucide",
		PackageManager: "pnpm",
	}

	if err := RunScaffold(tmpDir, cfg); err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}

	appPath := filepath.Join(tmpDir, "apps", "web", "App.tsx")
	appContent, err := os.ReadFile(appPath)
	if err != nil {
		t.Fatalf("Failed to read Expo App.tsx: %v", err)
	}
	appStr := string(appContent)

	if !strings.Contains(appStr, "lucide-react-native") {
		t.Errorf("Expo App.tsx missing lucide-react-native import")
	}
	if strings.Contains(appStr, "from 'lucide-react'") {
		t.Errorf("Expo App.tsx should not import web lucide-react")
	}
	if !strings.Contains(appStr, "ClerkProvider") {
		t.Errorf("Expo App.tsx missing ClerkProvider")
	}

	// Validate package.json
	pkgPath := filepath.Join(tmpDir, "apps", "web", "package.json")
	pkgContent, err := os.ReadFile(pkgPath)
	if err != nil {
		t.Fatalf("Failed to read Expo package.json: %v", err)
	}
	pkgStr := string(pkgContent)

	if !strings.Contains(pkgStr, "lucide-react-native") {
		t.Errorf("Expo package.json missing lucide-react-native dependency")
	}
	if strings.Contains(pkgStr, "\"lucide-react\":") {
		t.Errorf("Expo package.json should not contain web lucide-react")
	}
	if !strings.Contains(pkgStr, "@clerk/clerk-expo") {
		t.Errorf("Expo package.json missing @clerk/clerk-expo")
	}

	validateProjectTree(t, tmpDir, cfg.ProjectName)
}

func TestAutoWiringExpressBetterAuthResend(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := types.ScaffoldConfig{
		ProjectName:    "test-express-auth-resend",
		Frontend:       "react",
		Backend:        "express",
		API:            "trpc",
		Auth:           "better-auth",
		Database:       "postgres",
		ORM:            "drizzle",
		Addons:         "resend",
		PackageManager: "pnpm",
	}

	if err := RunScaffold(tmpDir, cfg); err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}

	// 1. Validate Express index.ts
	expressPath := filepath.Join(tmpDir, "apps", "api", "src", "index.ts")
	expressContent, err := os.ReadFile(expressPath)
	if err != nil {
		t.Fatalf("Failed to read Express index.ts: %v", err)
	}
	expressStr := string(expressContent)

	if !strings.Contains(expressStr, "toNodeHandler(auth)") {
		t.Errorf("Express index.ts missing toNodeHandler(auth)")
	}
	if !strings.Contains(expressStr, "trpcMiddleware") {
		t.Errorf("Express index.ts missing trpcMiddleware")
	}

	// 2. Validate packages/auth/src/index.ts
	authPath := filepath.Join(tmpDir, "packages", "auth", "src", "index.ts")
	authContent, err := os.ReadFile(authPath)
	if err != nil {
		t.Fatalf("Failed to read auth index.ts: %v", err)
	}
	authStr := string(authContent)

	if !strings.Contains(authStr, "Resend") || !strings.Contains(authStr, "sendResetPassword") {
		t.Errorf("auth index.ts missing Resend integration")
	}

	// 3. Validate packages/auth/package.json has resend and @test-express-auth-resend/db
	authPkgPath := filepath.Join(tmpDir, "packages", "auth", "package.json")
	authPkgContent, err := os.ReadFile(authPkgPath)
	if err != nil {
		t.Fatalf("Failed to read auth package.json: %v", err)
	}
	authPkgStr := string(authPkgContent)

	if !strings.Contains(authPkgStr, "\"resend\":") {
		t.Errorf("packages/auth package.json missing resend dependency")
	}
	if !strings.Contains(authPkgStr, "@test-express-auth-resend/db") {
		t.Errorf("packages/auth package.json missing @project/db workspace dependency")
	}

	validateProjectTree(t, tmpDir, cfg.ProjectName)
}

func TestAutoWiringFastAPISQLAlchemyLifespan(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := types.ScaffoldConfig{
		ProjectName:    "test-fastapi-db",
		Frontend:       "svelte",
		Backend:        "fastapi",
		Database:       "postgres",
		ORM:            "sqlalchemy",
		PackageManager: "pnpm",
	}

	if err := RunScaffold(tmpDir, cfg); err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}

	mainPath := filepath.Join(tmpDir, "apps", "api", "main.py")
	mainContent, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("Failed to read FastAPI main.py: %v", err)
	}
	mainStr := string(mainContent)

	if !strings.Contains(mainStr, "asynccontextmanager") || !strings.Contains(mainStr, "lifespan=lifespan") {
		t.Errorf("FastAPI main.py missing lifespan context manager")
	}

	reqPath := filepath.Join(tmpDir, "apps", "api", "requirements.txt")
	reqContent, err := os.ReadFile(reqPath)
	if err != nil {
		t.Fatalf("Failed to read FastAPI requirements.txt: %v", err)
	}
	reqStr := string(reqContent)

	if !strings.Contains(reqStr, "sqlalchemy") || !strings.Contains(reqStr, "psycopg2-binary") {
		t.Errorf("FastAPI requirements.txt missing sqlalchemy or postgres driver")
	}

	validateProjectTree(t, tmpDir, cfg.ProjectName)
}

func TestAutoWiringGoChiGORMHealth(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := types.ScaffoldConfig{
		ProjectName:    "test-gochi-gorm",
		Frontend:       "nuxt",
		Backend:        "go_chi",
		Database:       "postgres",
		ORM:            "gorm",
		PackageManager: "pnpm",
	}

	if err := RunScaffold(tmpDir, cfg); err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}

	// 1. Validate main.go
	mainPath := filepath.Join(tmpDir, "apps", "api", "main.go")
	mainContent, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("Failed to read Go Chi main.go: %v", err)
	}
	mainStr := string(mainContent)

	if !strings.Contains(mainStr, "db.Connect()") {
		t.Errorf("Go Chi main.go missing db.Connect()")
	}

	// 2. Validate handlers/health.go
	healthPath := filepath.Join(tmpDir, "apps", "api", "handlers", "health.go")
	healthContent, err := os.ReadFile(healthPath)
	if err != nil {
		t.Fatalf("Failed to read Go Chi health.go: %v", err)
	}
	healthStr := string(healthContent)

	if !strings.Contains(healthStr, "sqlDB.Ping()") {
		t.Errorf("Go Chi health.go missing database ping check")
	}

	validateProjectTree(t, tmpDir, cfg.ProjectName)
}
