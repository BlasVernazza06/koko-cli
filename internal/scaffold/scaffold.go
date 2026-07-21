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
	DatabaseMongoDB  DatabaseType = "mongodb"
	DatabaseNone     DatabaseType = "none"
)

type ORMType string

const (
	ORMPrisma  ORMType = "prisma"
	ORMSqlx    ORMType = "sqlx"
	ORMDrizzle ORMType = "drizzle"
	ORMNone    ORMType = "none"
)

type ScaffoldConfig struct {
	ProjectName   string
	Frontend      FrontendType
	Backend       BackendType
	Database      DatabaseType
	ORM           ORMType
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

		if strings.HasSuffix(destSubPath, ".tmpl") {
			destSubPath = strings.TrimSuffix(destSubPath, ".tmpl")
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
	case strings.HasPrefix(rel, "frontend/react-vite/"):
		if config.Frontend != FrontendReact {
			return "", false
		}
		dest, _ := filepath.Rel("frontend/react-vite", rel)
		return filepath.Join("apps", "frontend", dest), true

	case strings.HasPrefix(rel, "frontend/next/"):
		if config.Frontend != FrontendNext {
			return "", false
		}
		dest, _ := filepath.Rel("frontend/next", rel)
		return filepath.Join("apps", "frontend", dest), true

	case strings.HasPrefix(rel, "frontend/vue/"):
		if config.Frontend != FrontendVue {
			return "", false
		}
		dest, _ := filepath.Rel("frontend/vue", rel)
		return filepath.Join("apps", "frontend", dest), true

	case strings.HasPrefix(rel, "backend/fiber/"):
		if config.Backend != BackendFiber {
			return "", false
		}
		dest, _ := filepath.Rel("backend/fiber", rel)
		return filepath.Join("apps", "backend", dest), true

	case strings.HasPrefix(rel, "backend/express/"):
		if config.Backend != BackendExpress {
			return "", false
		}
		dest, _ := filepath.Rel("backend/express", rel)
		return filepath.Join("apps", "backend", dest), true

	case strings.HasPrefix(rel, "backend/hono/"):
		if config.Backend != BackendHono {
			return "", false
		}
		dest, _ := filepath.Rel("backend/hono", rel)
		return filepath.Join("apps", "backend", dest), true

	case strings.HasPrefix(rel, "db/prisma/postgres/"):
		if config.Database != DatabasePostgres || config.ORM != ORMPrisma {
			return "", false
		}
		dest, _ := filepath.Rel("db/prisma/postgres", rel)
		return filepath.Join("packages", "db", dest), true

	case strings.HasPrefix(rel, "db/prisma/mysql/"):
		if config.Database != DatabaseMySQL || config.ORM != ORMPrisma {
			return "", false
		}
		dest, _ := filepath.Rel("db/prisma/mysql", rel)
		return filepath.Join("packages", "db", dest), true

	case strings.HasPrefix(rel, "db/prisma/mongodb/"):
		if config.Database != DatabaseMongoDB || config.ORM != ORMPrisma {
			return "", false
		}
		dest, _ := filepath.Rel("db/prisma/mongodb", rel)
		return filepath.Join("packages", "db", dest), true

	case strings.HasPrefix(rel, "db/sqlx/postgres/"):
		if config.Database != DatabasePostgres || config.ORM != ORMSqlx {
			return "", false
		}
		dest, _ := filepath.Rel("db/sqlx/postgres", rel)
		return filepath.Join("packages", "db", dest), true

	case strings.HasPrefix(rel, "db/sqlx/mysql/"):
		if config.Database != DatabaseMySQL || config.ORM != ORMSqlx {
			return "", false
		}
		dest, _ := filepath.Rel("db/sqlx/mysql", rel)
		return filepath.Join("packages", "db", dest), true

	case rel == "db/drizzle/package.json":
		if config.ORM != ORMDrizzle {
			return "", false
		}
		return filepath.Join("packages", "db", "package.json"), true

	case rel == "db/drizzle/drizzle.config.ts":
		if config.ORM != ORMDrizzle {
			return "", false
		}
		return filepath.Join("packages", "db", "drizzle.config.ts"), true

	case strings.HasPrefix(rel, "db/drizzle/postgres/"):
		if config.Database != DatabasePostgres || config.ORM != ORMDrizzle {
			return "", false
		}
		dest, _ := filepath.Rel("db/drizzle/postgres", rel)
		return filepath.Join("packages", "db", dest), true

	case strings.HasPrefix(rel, "db/drizzle/mysql/"):
		if config.Database != DatabaseMySQL || config.ORM != ORMDrizzle {
			return "", false
		}
		dest, _ := filepath.Rel("db/drizzle/mysql", rel)
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
