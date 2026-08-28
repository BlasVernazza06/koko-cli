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
	"github.com/BlasVernazza06/koko-cli/internal/types"
)

//go:embed all:templates all:manual
var templateFs embed.FS

type ScaffoldConfig = types.ScaffoldConfig

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
