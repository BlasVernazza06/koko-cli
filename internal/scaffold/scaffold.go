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

	"github.com/BlasVernazza06/koko-cli/internal/catalog"
)

//go:embed all:templates all:manual
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

// RunScaffold genera el proyecto en memoria (VFS), escribe los archivos
// atómicamente en disco e inicializa Git si fue solicitado.
func RunScaffold(targetDir string, config ScaffoldConfig) error {
	virtualFS, err := GenerateVFS(config)
	if err != nil {
		return err
	}

	if err := WriteTree(virtualFS, targetDir); err != nil {
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

	err = fs.WalkDir(templateFs, ".", func(path string, d fs.DirEntry, err error) error {
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
		funcMap := template.FuncMap{
			"version": func(pkg string) string {
				return catalog.GetVersion(pkg)
			},
		}
		tmpl, err := template.New(path).Delims("[[", "]]").Funcs(funcMap).Parse(string(content))
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

	if err == nil {
		hasApps := (config.Frontend != "" && config.Frontend != "none") || (config.Backend != "" && config.Backend != "none")
		if hasApps {
			_ = os.Remove(filepath.Join(targetDir, "apps", ".gitkeep"))
		}
	}

	return err
}

// evaluatePath decides if a file should be copied according to config
// and returns the clean destination path and a boolean (shouldCopy)
func evaluatePath(path string, config ScaffoldConfig) (string, bool) {
	rel := filepath.ToSlash(path)

	// Mode 1: Recipe Templates
	if config.Recipe != "" {
		recipePrefix := "templates/recipes/" + config.Recipe + "/"
		if strings.HasPrefix(rel, recipePrefix) {
			dest := strings.TrimPrefix(rel, recipePrefix)
			return dest, true
		}

		if rel == "manual/docker/docker-compose.yml" || rel == "templates/docker/docker-compose.yml" {
			if config.Recipe == "saas" || config.Recipe == "pern" || config.Recipe == "mern" {
				return "docker-compose.yml", true
			}
		}

		if rel == "manual/github/ci.yml" || rel == "templates/github/ci.yml" {
			if config.Recipe == "saas" {
				return ".github/workflows/ci.yml", true
			}
		}

		return "", false
	}

	// Mode 2: Manual Configuration Templates

	// Monorepo root files -> root
	if strings.HasPrefix(rel, "manual/root/") {
		if rel == "manual/root/apps/.gitkeep" {
			if (config.Frontend != "" && config.Frontend != "none") || (config.Backend != "" && config.Backend != "none") {
				return "", false
			}
		}
		dest := strings.TrimPrefix(rel, "manual/root/")
		return dest, true
	}

	// Shared packages -> packages/
	if strings.HasPrefix(rel, "manual/packages/") {
		dest := strings.TrimPrefix(rel, "manual/")
		return dest, true
	}

	// Frontend selection -> apps/web
	if config.Frontend != "" && config.Frontend != "none" {
		frontendPrefix := "manual/frontend/" + config.Frontend + "/"
		if strings.HasPrefix(rel, frontendPrefix) {
			dest := "apps/web/" + strings.TrimPrefix(rel, frontendPrefix)
			return dest, true
		}
	}

	// Backend selection -> apps/api
	if config.Backend != "" && config.Backend != "none" {
		backendPrefix := "manual/backend/" + config.Backend + "/"
		if strings.HasPrefix(rel, backendPrefix) {
			dest := "apps/api/" + strings.TrimPrefix(rel, backendPrefix)
			return dest, true
		}
	}

	// Database and ORM routing
	if config.Database != "" && config.Database != "none" {
		orm := strings.ToLower(config.ORM)
		db := strings.ToLower(config.Database)

		if orm == "drizzle" {
			// Specific db files
			drizzleDbPrefix := "manual/db/drizzle/" + db + "/"
			if strings.HasPrefix(rel, drizzleDbPrefix) {
				dest := "packages/db/src/" + strings.TrimPrefix(rel, drizzleDbPrefix)
				return dest, true
			}
			// Root config files for drizzle package
			if rel == "manual/db/drizzle/drizzle.config.ts" {
				return "packages/db/drizzle.config.ts", true
			}
			if rel == "manual/db/drizzle/package.json" {
				return "packages/db/package.json", true
			}
		} else if orm == "prisma" {
			// Specific prisma schemas
			prismaDbPrefix := "manual/db/prisma/" + db + "/"
			if strings.HasPrefix(rel, prismaDbPrefix) {
				dest := "packages/db/prisma/" + strings.TrimPrefix(rel, prismaDbPrefix)
				return dest, true
			}
			// Root package.json for prisma package
			if rel == "manual/db/prisma/package.json" {
				return "packages/db/package.json", true
			}
		} else if orm == "mongoose" || orm == "moongose" {
			mongoosePrefix := "manual/db/mongoose/mongodb/"
			if strings.HasPrefix(rel, mongoosePrefix) {
				dest := "packages/db/" + strings.TrimPrefix(rel, mongoosePrefix)
				return dest, true
			}
		} else if orm == "sqlalchemy" {
			sqlAlchemyPrefix := "manual/db/sqlalchemy/" + db + "/"
			if strings.HasPrefix(rel, sqlAlchemyPrefix) {
				dest := "apps/api/app/db/" + strings.TrimPrefix(rel, sqlAlchemyPrefix)
				return dest, true
			}
		} else if orm == "gorm" {
			gormPrefix := "manual/db/gorm/" + db + "/"
			if strings.HasPrefix(rel, gormPrefix) {
				dest := "apps/api/db/" + strings.TrimPrefix(rel, gormPrefix)
				return dest, true
			}
		}
	}

	// Addons: Docker Compose
	if rel == "manual/docker/docker-compose.yml" {
		if config.Addons == "docker" || config.Addons == "docker_cicd" || (config.Database != "" && config.Database != "none" && config.Database != "sqlite") {
			return "docker-compose.yml", true
		}
	}

	// Addons: GitHub Actions CI
	if rel == "manual/github/ci.yml" {
		if config.Addons == "github_actions" || config.Addons == "docker_cicd" {
			return ".github/workflows/ci.yml", true
		}
	}

	return "", false
}

func isDockerOrDBFile(path string) bool {
	fullLower := strings.ToLower(filepath.ToSlash(path))
	return strings.HasSuffix(fullLower, "docker-compose.yml") ||
		strings.Contains(fullLower, "/packages/db/") ||
		strings.Contains(fullLower, "/apps/api/db/") ||
		strings.Contains(fullLower, "/apps/api/app/db/") ||
		strings.Contains(fullLower, "schema.prisma") ||
		strings.Contains(fullLower, "drizzle.config.ts")
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
		strings.Contains(s, "[[ version ") ||
		strings.Contains(s, "[[version ") ||
		strings.Contains(s, "[[if ") ||
		strings.Contains(s, "[[range ") ||
		strings.Contains(s, "[[with ") ||
		strings.Contains(s, "[[define ") ||
		strings.Contains(s, "[[template ") ||
		strings.Contains(s, "[[/*")
}
