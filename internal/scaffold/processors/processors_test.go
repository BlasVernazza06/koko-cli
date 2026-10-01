package processors

import (
	"strings"
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/catalog"
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

	var dbPkg map[string]interface{}
	_ = v.ReadJSON("packages/db/package.json", &dbPkg)
	deps := dbPkg["dependencies"].(map[string]interface{})
	if deps["postgres"] != catalog.GetVersion("postgres") {
		t.Errorf("Expected dynamically added postgres dependency '%s', got '%v'", catalog.GetVersion("postgres"), deps["postgres"])
	}
}

func TestProcessPackageJSONs_PaymentsAndEmailAddons(t *testing.T) {
	v := vfs.New()
	v.WriteString("package.json", `{"name": "placeholder"}`)
	v.WriteString("apps/web/package.json", `{"name": "placeholder-web", "dependencies": {}}`)
	v.WriteString("apps/api/package.json", `{"name": "placeholder-api", "dependencies": {}}`)

	cfg := ProcessConfig{
		ProjectName:    "acme-commerce",
		PackageManager: "pnpm",
		Frontend:       "nextjs",
		Backend:        "express",
		Addons:         "stripe,polar,resend,brevo",
	}

	if err := ProcessPackageJSONs(v, cfg); err != nil {
		t.Fatalf("ProcessPackageJSONs failed: %v", err)
	}

	// Verify web package.json has payments & email deps
	var webPkg map[string]interface{}
	if err := v.ReadJSON("apps/web/package.json", &webPkg); err != nil {
		t.Fatalf("Failed to read web package.json: %v", err)
	}
	webDeps := webPkg["dependencies"].(map[string]interface{})
	for _, pkgName := range []string{"stripe", "@polar-sh/sdk", "resend", "@getbrevo/brevo"} {
		if _, ok := webDeps[pkgName]; !ok {
			t.Errorf("Expected %s in apps/web dependencies, got %+v", pkgName, webDeps)
		}
	}

	// Verify api package.json has payments & email deps
	var apiPkg map[string]interface{}
	if err := v.ReadJSON("apps/api/package.json", &apiPkg); err != nil {
		t.Fatalf("Failed to read api package.json: %v", err)
	}
	apiDeps := apiPkg["dependencies"].(map[string]interface{})
	for _, pkgName := range []string{"stripe", "@polar-sh/sdk", "resend", "@getbrevo/brevo"} {
		if _, ok := apiDeps[pkgName]; !ok {
			t.Errorf("Expected %s in apps/api dependencies, got %+v", pkgName, apiDeps)
		}
	}
}

func TestAddDependency(t *testing.T) {
	pkg := make(map[string]interface{})

	AddDependency(pkg, "drizzle-orm", false)
	AddDependency(pkg, "typescript", true)

	deps := pkg["dependencies"].(map[string]interface{})
	if deps["drizzle-orm"] != catalog.GetVersion("drizzle-orm") {
		t.Errorf("Expected drizzle-orm version '%s', got '%v'", catalog.GetVersion("drizzle-orm"), deps["drizzle-orm"])
	}

	devDeps := pkg["devDependencies"].(map[string]interface{})
	if devDeps["typescript"] != catalog.GetVersion("typescript") {
		t.Errorf("Expected typescript version '%s', got '%v'", catalog.GetVersion("typescript"), devDeps["typescript"])
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
		Addons:      "stripe,polar,resend",
	}
	if err := ProcessEnvVariables(v, cfg); err != nil {
		t.Fatalf("ProcessEnvVariables failed: %v", err)
	}
	// 1. Validaciones en .env (Archivo Local)
	envContent, ok := v.ReadString(".env")
	if !ok {
		t.Fatalf("Expected .env file to be created")
	}
	// Database URL
	if !strings.Contains(envContent, "DATABASE_URL=\"postgresql://postgres:password@localhost:5432/shop-app?schema=public\"") {
		t.Errorf("Expected Postgres connection string in .env, got: %s", envContent)
	}
	// URLs de API
	if !strings.Contains(envContent, "NEXT_PUBLIC_API_URL=http://localhost:4000") {
		t.Errorf("Expected Next.js API url in .env, got: %s", envContent)
	}
	// Better-Auth Secret: Verificar que no sea el placeholder estático y tenga formato base64 seguro
	if strings.Contains(envContent, "supersecret-auth-key") {
		t.Errorf("Expected generated random secret, but found static placeholder in .env")
	}
	if !strings.Contains(envContent, "BETTER_AUTH_SECRET=\"") {
		t.Errorf("Expected BETTER_AUTH_SECRET in .env, got: %s", envContent)
	}
	// Addons de Pagos y Email en .env
	if !strings.Contains(envContent, "STRIPE_SECRET_KEY=\"sk_test_...\"") {
		t.Errorf("Expected Stripe secret key in .env, got: %s", envContent)
	}
	if !strings.Contains(envContent, "POLAR_ACCESS_TOKEN=\"polar_atfs_...\"") {
		t.Errorf("Expected Polar access token in .env, got: %s", envContent)
	}
	if !strings.Contains(envContent, "RESEND_API_KEY=\"re_...\"") {
		t.Errorf("Expected Resend API key in .env, got: %s", envContent)
	}
	// 2. Validaciones en .env.example (Archivo para Git)
	exampleContent, ok := v.ReadString(".env.example")
	if !ok {
		t.Fatalf("Expected .env.example file to be created")
	}
	// En .env.example debe mantenerse el placeholder seguro para Git
	if !strings.Contains(exampleContent, "BETTER_AUTH_SECRET=\"your-better-auth-secret\"") {
		t.Errorf("Expected generic placeholder in .env.example, got: %s", exampleContent)
	}
	if !strings.Contains(exampleContent, "STRIPE_SECRET_KEY=\"sk_test_...\"") {
		t.Errorf("Expected Stripe placeholder in .env.example, got: %s", exampleContent)
	}
}
