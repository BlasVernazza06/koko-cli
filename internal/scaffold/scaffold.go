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

type ScaffoldConfig struct {
	ProjectName    string
	Recipe         string
	InitGit        bool
	Frontend       string
	Backend        string
	PackageManager string
	Database       string
	ORM            string
	Auth           string
	Addons         string
}

func RunScaffold(targetDir string, config ScaffoldConfig) error {
	if err := CopyTemplates(targetDir, config); err != nil {
		return err
	}
	if err := GenerateDockerAndDB(targetDir, config); err != nil {
		return err
	}
	if config.InitGit {
		return InitGit(targetDir)
	}
	return nil
}

func CopyTemplates(targetDir string, config ScaffoldConfig) error {
	return walkAndCopy(targetDir, config, func(destPath string) bool {
		return !isDockerOrDBFile(destPath)
	})
}

func GenerateDockerAndDB(targetDir string, config ScaffoldConfig) error {
	return walkAndCopy(targetDir, config, func(destPath string) bool {
		return isDockerOrDBFile(destPath)
	})
}

func walkAndCopy(targetDir string, config ScaffoldConfig, filterFn func(string) bool) error {
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
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

		if filterFn != nil && !filterFn(destPath) {
			return nil
		}

		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create subdirectory for %s: %w", destPath, err)
		}

		content, err := templateFs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read virtual file %s: %w", path, err)
		}

		// If it is a binary file, copy directly without passing through template engine
		if bytes.IndexByte(content, 0) != -1 || isBinaryExtension(destPath) {
			err = os.WriteFile(destPath, content, 0644)
			if err != nil {
				return fmt.Errorf("failed to write binary file %s: %w", destPath, err)
			}
			return nil
		}

		// If it does not contain Go template tags, copy directly without parsing
		if !shouldParseAsTemplate(content) {
			err = os.WriteFile(destPath, content, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", destPath, err)
			}
			return nil
		}

		// Template Interpolation Engine:
		// Parse content as Go template and render with config
		tmpl, err := template.New(path).Delims("[[", "]]").Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to process template for %s: %w", path, err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, config)
		if err != nil {
			return fmt.Errorf("failed to interpolate variables in template %s: %w", path, err)
		}

		// Write processed file to disk
		err = os.WriteFile(destPath, buf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}
		return nil
	})

	return err
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
	case strings.HasPrefix(rel, "recipes/"):
		// Si es una receta, comprobamos si es la seleccionada
		prefix := "recipes/" + config.Recipe + "/"
		if config.Recipe != "" && strings.HasPrefix(rel, prefix) {
			dest, _ := filepath.Rel(prefix, rel)
			return dest, true
		}
		return "", false

	case rel == "docker/docker-compose.yml":
		if config.Recipe == "saas" || config.Recipe == "pern" || config.Recipe == "mern" {
			return "docker-compose.yml", true
		}
		return "", false

	case rel == "github/ci.yml":
		if config.Recipe == "saas" {
			return ".github/workflows/ci.yml", true
		}
		return "", false
	}

	return "", false
}

func isDockerOrDBFile(path string) bool {
	fullLower := strings.ToLower(path)
	return strings.Contains(fullLower, "docker") ||
		strings.Contains(fullLower, "db") ||
		strings.Contains(fullLower, "drizzle") ||
		strings.Contains(fullLower, "prisma") ||
		strings.Contains(fullLower, "schema") ||
		strings.Contains(fullLower, "database") ||
		strings.Contains(fullLower, "mongodb") ||
		strings.Contains(fullLower, "postgres") ||
		strings.Contains(fullLower, "mysql")
}

func InitGit(targetDir string) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return nil
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = targetDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	return nil
}

func isBinaryExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".ico", ".pdf", ".woff", ".woff2", ".ttf", ".eot":
		return true
	}
	return false
}

func shouldParseAsTemplate(content []byte) bool {
	s := string(content)
	return strings.Contains(s, "[[.") ||
		strings.Contains(s, "[[if ") ||
		strings.Contains(s, "[[range ") ||
		strings.Contains(s, "[[with ") ||
		strings.Contains(s, "[[define ") ||
		strings.Contains(s, "[[template ") ||
		strings.Contains(s, "[[/*")
}
