package doctor

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/BlasVernazza06/koko-cli/internal/core/config"
)

type dummyConn struct {
	net.Conn
}

func (d *dummyConn) Close() error {
	return nil
}

func TestHealthChecker_Environment(t *testing.T) {
	tempDir := t.TempDir()

	checker := &HealthChecker{
		lookPath: func(file string) (string, error) {
			if file == "node" || file == "pnpm" || file == "docker" {
				return "/usr/bin/" + file, nil
			}
			return "", os.ErrNotExist
		},
		runCommand: func(ctx context.Context, name string, args ...string) (string, error) {
			if name == "node" {
				return "v20.10.0", nil
			}
			if name == "pnpm" {
				return "9.1.0", nil
			}
			if name == "docker" {
				return "Docker version 27.0", nil
			}
			return "", nil
		},
		tcpDial: func(network, address string, timeout time.Duration) (net.Conn, error) {
			return nil, os.ErrNotExist
		},
	}

	detected := config.KokoConfig{
		Architecture: config.ArchitectureInfo{PackageManager: "pnpm"},
		Features:     config.FeaturesInfo{Infrastructure: &config.Infrastructure{DockerCompose: true}},
	}

	items := checker.CheckEnvironment(context.Background(), tempDir, detected)
	if len(items) != 4 {
		t.Fatalf("expected 4 items (node, pnpm, git, docker), got %d", len(items))
	}

	// Node check
	if items[0].Status != StatusOk || !strings.Contains(items[0].Message, "v20.10.0") {
		t.Errorf("expected node OK v20.10.0, got %+v", items[0])
	}

	// PM check
	if items[1].Status != StatusOk || !strings.Contains(items[1].Message, "v9.1.0") {
		t.Errorf("expected pnpm OK, got %+v", items[1])
	}

	// Git uninitialized
	if items[2].Status != StatusWarn || items[2].FixCommand != "git init" {
		t.Errorf("expected git warn, got %+v", items[2])
	}

	// Docker running
	if items[3].Status != StatusOk {
		t.Errorf("expected docker daemon OK, got %+v", items[3])
	}
}

func TestHealthChecker_EnvFilesAndAutoFix(t *testing.T) {
	tempDir := t.TempDir()

	exampleContent := `
# Database
DATABASE_URL="postgres://user:password@localhost:5432/mydb"
AUTH_SECRET="your_secret_here"
PORT=3000
`
	if err := os.WriteFile(filepath.Join(tempDir, ".env.example"), []byte(exampleContent), 0644); err != nil {
		t.Fatalf("failed to write .env.example: %v", err)
	}

	checker := NewHealthChecker()

	// 1. Sin .env presente: debe alertar StatusFail
	items := checker.CheckEnvFiles(tempDir)
	if len(items) == 0 || items[0].Status != StatusFail {
		t.Fatalf("expected missing .env fail, got %+v", items)
	}

	// 2. Ejecutar AutoFix
	copied, err := FixEnvFiles(tempDir)
	if err != nil {
		t.Fatalf("unexpected fix error: %v", err)
	}
	if copied != 1 {
		t.Errorf("expected 1 file copied, got %d", copied)
	}
	if !fileExists(filepath.Join(tempDir, ".env")) {
		t.Errorf(".env should exist after fix")
	}

	// 3. Volver a chequear: ahora tiene placeholder en AUTH_SECRET
	itemsAfterFix := checker.CheckEnvFiles(tempDir)
	foundPlaceholder := false
	for _, it := range itemsAfterFix {
		if it.Status == StatusWarn && strings.Contains(it.Detail, "AUTH_SECRET") {
			foundPlaceholder = true
		}
	}
	if !foundPlaceholder {
		t.Errorf("expected placeholder warning for AUTH_SECRET, got %+v", itemsAfterFix)
	}
}

func TestHealthChecker_DatabaseAndPorts(t *testing.T) {
	tempDir := t.TempDir()

	// Crear schema.prisma simulado
	prismaDir := filepath.Join(tempDir, "prisma")
	_ = os.MkdirAll(prismaDir, 0755)
	_ = os.WriteFile(filepath.Join(prismaDir, "schema.prisma"), []byte("datasource db {}"), 0644)

	// Iniciar un listener TCP local simulado para puerto 3000 (puerto en conflicto)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	listeningPort := listener.Addr().(*net.TCPAddr).Port

	checker := &HealthChecker{
		lookPath: exec.LookPath,
		runCommand: func(ctx context.Context, name string, args ...string) (string, error) {
			return "", nil
		},
		tcpDial: func(network, address string, timeout time.Duration) (net.Conn, error) {
			if strings.HasSuffix(address, fmt.Sprintf(":%d", listeningPort)) {
				// Simular conexión exitosa (puerto ocupado)
				return &dummyConn{}, nil
			}
			return nil, os.ErrNotExist
		},
	}

	detected := config.KokoConfig{
		Stack: config.StackInfo{
			Frontend: &config.FrontendInfo{Framework: "next"},
			Database: &config.DatabaseInfo{Provider: "postgres", ORM: "prisma"},
		},
	}

	items := checker.CheckDatabaseAndPorts(tempDir, detected)
	foundSchema := false
	for _, it := range items {
		if strings.Contains(it.Name, "Database Schema") && it.Status == StatusOk {
			foundSchema = true
		}
	}
	if !foundSchema {
		t.Errorf("expected schema OK, got %+v", items)
	}
}

func TestHealthChecker_SkillsContext(t *testing.T) {
	tempDir := t.TempDir()

	detected := config.KokoConfig{
		Stack: config.StackInfo{
			Frontend: &config.FrontendInfo{Framework: "next", Styling: "tailwindcss", UILibrary: "shadcn"},
			Database: &config.DatabaseInfo{ORM: "drizzle"},
		},
		Features: config.FeaturesInfo{
			Payments: &config.PaymentInfo{Provider: "stripe"},
		},
	}

	checker := NewHealthChecker()

	// 1. Sin reglas instaladas: debe advertir
	items := checker.CheckSkillsContext(tempDir, detected)
	if len(items) == 0 || items[0].Status != StatusWarn {
		t.Errorf("expected warn when skills missing, got %+v", items)
	}

	// 2. Instalar reglas simuladas en .cursor/rules/
	rulesDir := filepath.Join(tempDir, ".cursor", "rules")
	_ = os.MkdirAll(rulesDir, 0755)
	for _, s := range []string{"nextjs", "tailwindcss", "shadcn-ui", "drizzle", "stripe-best-practices"} {
		_ = os.WriteFile(filepath.Join(rulesDir, s+".mdc"), []byte("---"), 0644)
	}

	itemsAfter := checker.CheckSkillsContext(tempDir, detected)
	if len(itemsAfter) == 0 || itemsAfter[0].Status != StatusOk {
		t.Errorf("expected OK when skills installed, got %+v", itemsAfter)
	}
}
