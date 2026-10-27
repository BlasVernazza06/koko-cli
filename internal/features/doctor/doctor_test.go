package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDoctorDiagnosticsAndFix(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Crear estructura simulada de un proyecto Monorepo con Next.js + Tailwind + Better-Auth + Stripe + Resend + Drizzle
	webDir := filepath.Join(tempDir, "apps", "web")
	if err := os.MkdirAll(webDir, 0755); err != nil {
		t.Fatalf("failed to create webDir: %v", err)
	}

	pkgJSON := `{
		"name": "my-web-app",
		"dependencies": {
			"next": "^14.2.13",
			"react": "^18.3.1",
			"better-auth": "^1.1.0",
			"tailwindcss": "^3.4.12",
			"stripe": "^17.7.0",
			"resend": "^4.1.2",
			"drizzle-orm": "^0.30.10",
			"postgres": "^3.4.4",
			"lucide-react": "^0.446.0"
		},
		"devDependencies": {
			"typescript": "^5.6.2"
		}
	}`
	if err := os.WriteFile(filepath.Join(webDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("failed to write web package.json: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tempDir, "pnpm-lock.yaml"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to write pnpm-lock.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "docker-compose.yml"), []byte(""), 0644); err != nil {
		t.Fatalf("failed to write docker-compose.yml: %v", err)
	}

	// 2. Correr diagnóstico sin archivo previo (koko.config.json no existe)
	diag, err := RunDiagnostics(tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diag.ConfigExists {
		t.Errorf("expected ConfigExists to be false, got true")
	}
	if diag.DetectedConfig.Architecture.Layout != "monorepo" {
		t.Errorf("expected layout 'monorepo', got %s", diag.DetectedConfig.Architecture.Layout)
	}
	if diag.DetectedConfig.Architecture.PackageManager != "pnpm" {
		t.Errorf("expected pm 'pnpm', got %s", diag.DetectedConfig.Architecture.PackageManager)
	}
	if diag.DetectedConfig.Stack.Frontend == nil || diag.DetectedConfig.Stack.Frontend.Framework != "next" {
		t.Errorf("expected frontend framework 'next', got %+v", diag.DetectedConfig.Stack.Frontend)
	}
	if diag.DetectedConfig.Stack.Frontend.Styling != "tailwindcss" {
		t.Errorf("expected styling 'tailwindcss', got %s", diag.DetectedConfig.Stack.Frontend.Styling)
	}
	if diag.DetectedConfig.Stack.Frontend.Icons != "lucide" {
		t.Errorf("expected icons 'lucide', got %s", diag.DetectedConfig.Stack.Frontend.Icons)
	}
	if diag.DetectedConfig.Stack.Database == nil || diag.DetectedConfig.Stack.Database.ORM != "drizzle" {
		t.Errorf("expected orm 'drizzle', got %+v", diag.DetectedConfig.Stack.Database)
	}
	if diag.DetectedConfig.Features.Auth == nil || diag.DetectedConfig.Features.Auth.Provider != "better-auth" {
		t.Errorf("expected auth 'better-auth', got %+v", diag.DetectedConfig.Features.Auth)
	}
	if diag.DetectedConfig.Features.Payments == nil || diag.DetectedConfig.Features.Payments.Provider != "stripe" {
		t.Errorf("expected payments 'stripe', got %+v", diag.DetectedConfig.Features.Payments)
	}
	if diag.DetectedConfig.Features.Email == nil || diag.DetectedConfig.Features.Email.Provider != "resend" {
		t.Errorf("expected email 'resend', got %+v", diag.DetectedConfig.Features.Email)
	}
	if diag.DetectedConfig.Features.Infrastructure == nil || !diag.DetectedConfig.Features.Infrastructure.DockerCompose {
		t.Errorf("expected dockerCompose true, got %+v", diag.DetectedConfig.Features.Infrastructure)
	}

	// 3. Aplicar Fixes y comprobar que el archivo se guardó
	err = ApplyFixes(tempDir, diag.DetectedConfig, diag.ExistingConfig)
	if err != nil {
		t.Fatalf("failed to apply fixes: %v", err)
	}

	// 4. Volver a correr diagnóstico: debe detectar que el archivo existe y no hay diferencias
	diagAfter, err := RunDiagnostics(tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !diagAfter.ConfigExists {
		t.Errorf("expected ConfigExists to be true after apply fixes")
	}
	if len(diagAfter.Differences) != 0 {
		t.Errorf("expected 0 differences after fix, got %d: %+v", len(diagAfter.Differences), diagAfter.Differences)
	}
}
