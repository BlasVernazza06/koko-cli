package scaffold

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:templates
var templateFs embed.FS

type FrontendType string

const (
	FrontendNext  FrontendType = "next"
	FrontendReact FrontendType = "react"
	FrontendVue   FrontendType = "vue"
	FrontendNone  FrontendType = "none"
)

type BackendType string

const (
	BackendFiber   BackendType = "fiber"
	BackendExpress BackendType = "express"
	BackendHono    BackendType = "hono"
	BackendNone    BackendType = "none"
)

type DatabaseType string

const (
	DatabasePostgres DatabaseType = "postgres"
	DatabaseMySQL    DatabaseType = "mysql"
	DatabaseNone     DatabaseType = "none"
)

type ScaffoldConfig struct {
	ProjectName   string
	Frontend      FrontendType
	Backend       BackendType
	Database      DatabaseType
	Docker        bool
	GithubActions bool
}

func RunScaffold(targetDir string, config ScaffoldConfig) error {
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("error al crear el directorio de destino: %w", err)
	}

	err = fs.WalkDir(templateFs, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		destSubPath, shouldCopy := evaluatePath(path, config)
		if !shouldCopy {
			return nil
		}

		destPath := filepath.Join(targetDir, destSubPath)

		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			return fmt.Errorf("error al crear subdirectorio para %s: %w", destPath, err)
		}

		content, err := templateFs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error al leer archivo virtual %s: %w", path, err)
		}

		// Motor de Interpolación (Task 1.4):
		// Parseamos el contenido como una plantilla de Go y la renderizamos con config
		tmpl, err := template.New(path).Delims("[[", "]]").Parse(string(content))
		if err != nil {
			return fmt.Errorf("error al procesar plantilla para %s: %w", path, err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, config)
		if err != nil {
			return fmt.Errorf("error al interpolar variables en plantilla %s: %w", path, err)
		}

		// Escribimos el archivo procesado en el disco real
		err = os.WriteFile(destPath, buf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("error al escribir archivo %s: %w", destPath, err)
		}
		fmt.Printf("✓ Generado: %s\n", destSubPath)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error recorriendo las plantillas: %w", err)
	}

	// Inicializar repositorio Git aquí (fuera del bucle walk)
	if err := initGit(targetDir); err != nil {
		fmt.Printf("⚠️  Advertencia: %v\n", err)
	}

	return nil
}

// evaluatePath decide si un archivo debe copiarse según la configuración
// y retorna la ruta de destino limpia y un booleano (si debe copiarse o no)
func evaluatePath(path string, config ScaffoldConfig) (string, bool) {
	// 1. Quitamos el prefijo "templates/" de la ruta.
	rel, err := filepath.Rel("templates", path)
	if err != nil {
		return "", false
	}

	// 2. Normalizamos las barras diagonales a formato '/' (estilo Unix)
	rel = filepath.ToSlash(rel)

	// 3. Filtramos los archivos según la configuración del usuario y ruteamos al monorepo:
	switch {
	case strings.HasPrefix(rel, "frontend/react-vite-ts/"):
		if config.Frontend != FrontendReact {
			return "", false
		}
		dest, _ := filepath.Rel("frontend/react-vite-ts", rel)
		return filepath.Join("apps", "frontend", dest), true

	case strings.HasPrefix(rel, "frontend/next-app-router/"):
		if config.Frontend != FrontendNext {
			return "", false
		}
		dest, _ := filepath.Rel("frontend/next-app-router", rel)
		return filepath.Join("apps", "frontend", dest), true

	case strings.HasPrefix(rel, "frontend/vue-vite-ts/"):
		if config.Frontend != FrontendVue {
			return "", false
		}
		dest, _ := filepath.Rel("frontend/vue-vite-ts", rel)
		return filepath.Join("apps", "frontend", dest), true

	case strings.HasPrefix(rel, "backend/go-fiber/"):
		if config.Backend != BackendFiber {
			return "", false
		}
		dest, _ := filepath.Rel("backend/go-fiber", rel)
		return filepath.Join("apps", "backend", dest), true

	case strings.HasPrefix(rel, "backend/node-express/"):
		if config.Backend != BackendExpress {
			return "", false
		}
		dest, _ := filepath.Rel("backend/node-express", rel)
		return filepath.Join("apps", "backend", dest), true

	case strings.HasPrefix(rel, "backend/hono-node/"):
		if config.Backend != BackendHono {
			return "", false
		}
		dest, _ := filepath.Rel("backend/hono-node", rel)
		return filepath.Join("apps", "backend", dest), true

	case strings.HasPrefix(rel, "db/postgres/"):
		if config.Database != DatabasePostgres {
			return "", false
		}
		if strings.HasSuffix(rel, ".prisma") && config.Backend != BackendExpress {
			return "", false
		}
		if strings.HasSuffix(rel, ".go") && config.Backend != BackendFiber {
			return "", false
		}
		dest, _ := filepath.Rel("db/postgres", rel)
		return filepath.Join("packages", "db", dest), true

	case strings.HasPrefix(rel, "db/mysql/"):
		if config.Database != DatabaseMySQL {
			return "", false
		}
		if strings.HasSuffix(rel, ".prisma") && config.Backend != BackendExpress {
			return "", false
		}
		if strings.HasSuffix(rel, ".go") && config.Backend != BackendFiber {
			return "", false
		}
		dest, _ := filepath.Rel("db/mysql", rel)
		return filepath.Join("packages", "db", dest), true

	case rel == "docker/docker-compose.yml":
		if !config.Docker {
			return "", false
		}
		return "docker-compose.yml", true

	case rel == "github/ci.yml":
		if !config.GithubActions {
			return "", false
		}
		return ".github/workflows/ci.yml", true
	}

	return "", false
}

func initGit(targetDir string) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return nil
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = targetDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("no se pudo inicializar git; %w", err)
	}

	return nil
}
