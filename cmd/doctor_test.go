package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestDoctorCommand_Comprehensive5Modules(t *testing.T) {
	tempDir := t.TempDir()

	// Crear proyecto simulado con .env.example
	exampleContent := "DATABASE_URL=postgres://localhost:5432/db\nAUTH_SECRET=your_secret_here\n"
	_ = os.WriteFile(filepath.Join(tempDir, ".env.example"), []byte(exampleContent), 0644)

	pkgJSON := `{
		"name": "sample-app",
		"dependencies": {
			"next": "^14.2.0",
			"react": "^18.3.0",
			"tailwindcss": "^3.4.0",
			"drizzle-orm": "^0.30.0",
			"stripe": "^17.0.0"
		}
	}`
	_ = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(pkgJSON), 0644)

	// 1. Ejecutar doctor sin --fix
	out := captureOutput(func() {
		rootCmd.SetArgs([]string{"doctor", "--dir", tempDir})
		_ = rootCmd.Execute()
	})

	// Verificar la presencia de los 5 bloques
	if !strings.Contains(out, "[1/5] Entorno & Herramientas") {
		t.Errorf("expected module 1 in output: %s", out)
	}
	if !strings.Contains(out, "[2/5] Stack & Arquitectura") {
		t.Errorf("expected module 2 in output: %s", out)
	}
	if !strings.Contains(out, "[3/5] Configuración & Variables (.env)") {
		t.Errorf("expected module 3 in output: %s", out)
	}
	if !strings.Contains(out, "[4/5] Base de Datos & Puertos Locales") {
		t.Errorf("expected module 4 in output: %s", out)
	}
	if !strings.Contains(out, "[5/5] Sincronización de Manifiesto") {
		t.Errorf("expected module 5 in output: %s", out)
	}

	// 2. Ejecutar doctor con --fix
	outFix := captureOutput(func() {
		rootCmd.SetArgs([]string{"doctor", "--dir", tempDir, "--fix"})
		_ = rootCmd.Execute()
	})

	if !strings.Contains(outFix, "successfully updated") {
		t.Errorf("expected success message with --fix, got: %s", outFix)
	}

	// Verificar que .env y koko.config.json fueron creados
	if _, err := os.Stat(filepath.Join(tempDir, ".env")); err != nil {
		t.Errorf(".env should be created by doctor --fix")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "koko.config.json")); err != nil {
		t.Errorf("koko.config.json should be created by doctor --fix")
	}
}
