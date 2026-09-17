package scaffold

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

type pkgJSON struct {
	Name                 string            `json:"name"`
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	TranspilePackages    []string          `json:"transpilePackages"`
}

type tsConfigJSON struct {
	Extends string `json:"extends"`
}

// validateProjectTree inspects a generated scaffolded project on disk and validates:
// 1. All package.jsons are valid JSON.
// 2. All workspace:* dependencies point to an existing package in the monorepo.
// 3. No unresolved template markers ([[, ]]) exist in files.
// 4. TSConfig extends paths resolve to existing packages.
func validateProjectTree(t *testing.T, projectDir string, projectName string) {
	t.Helper()

	workspacePackages := make(map[string]string) // name -> relative path

	// Pass 1: Collect all workspace packages
	err := filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "package.json" {
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				t.Errorf("[%s] Failed to read %s: %v", projectName, path, readErr)
				return nil
			}
			var p pkgJSON
			if jsonErr := json.Unmarshal(data, &p); jsonErr != nil {
				t.Errorf("[%s] Invalid JSON in %s: %v", projectName, path, jsonErr)
				return nil
			}
			if p.Name != "" {
				rel, _ := filepath.Rel(projectDir, path)
				workspacePackages[p.Name] = rel
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("[%s] WalkDir failed: %v", projectName, err)
	}

	// Pass 2: Validate workspace dependency links and template rendering
	_ = filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(projectDir, path)

		// 1. Verify no unresolved Go template tags
		content, readErr := os.ReadFile(path)
		if readErr == nil {
			if shouldParseAsTemplate(content) {
				t.Errorf("[%s] Found unrendered template tag in %s:\n%s", projectName, rel, string(content))
			}
		}

		// 2. Verify package.json workspace references
		if d.Name() == "package.json" {
			var p pkgJSON
			_ = json.Unmarshal(content, &p)

			allDeps := make(map[string]string)
			for k, v := range p.Dependencies {
				allDeps[k] = v
			}
			for k, v := range p.DevDependencies {
				allDeps[k] = v
			}
			for k, v := range p.PeerDependencies {
				allDeps[k] = v
			}

			for depName, version := range allDeps {
				if strings.HasPrefix(version, "workspace:") {
					if _, exists := workspacePackages[depName]; !exists {
						t.Errorf("[%s in %s] Broken workspace dependency: %q not found in workspace (available packages: %v)",
							projectName, rel, depName, getKeys(workspacePackages))
					}
				}
			}
		}

		// 3. Verify tsconfig.json extends references
		if d.Name() == "tsconfig.json" {
			var tsCfg tsConfigJSON
			_ = json.Unmarshal(content, &tsCfg)
			if tsCfg.Extends != "" {
				// E.g. "@memo123/typescript-config/base.json" or "@repo/typescript-config/nextjs.json"
				parts := strings.Split(tsCfg.Extends, "/")
				if len(parts) >= 2 && strings.HasPrefix(tsCfg.Extends, "@") {
					basePkg := parts[0] + "/" + parts[1]
					if _, exists := workspacePackages[basePkg]; !exists {
						t.Errorf("[%s in %s] tsconfig extends non-existent package %q (available: %v)",
							projectName, rel, basePkg, getKeys(workspacePackages))
					}
				}
			}
		}

		return nil
	})
}

func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// TestRecipesMatrixValidation executes scaffolding on every predefined recipe
// and validates all workspace connections, template interpolation, and package relationships.
func TestRecipesMatrixValidation(t *testing.T) {
	recipes := []string{
		"saas",
		"pern",
		"mern",
		"fastapi_react",
		"enterprise_nestjs",
		"mobile_expo",
		"java_spring",
	}

	for _, recipe := range recipes {
		t.Run("Recipe_"+recipe, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "koko-test-recipe-"+recipe+"-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			projectName := "test-" + recipe + "-app"
			cfg := types.ScaffoldConfig{
				ProjectName: projectName,
				Recipe:      recipe,
			}

			if err := RunScaffold(tmpDir, cfg); err != nil {
				t.Fatalf("Scaffold failed for recipe %s: %v", recipe, err)
			}

			validateProjectTree(t, tmpDir, projectName)
		})
	}
}

// TestManualCombinationsMatrixValidation runs scaffolding across a broad permutation of manual configurations
// testing frontends, backends, databases, ORMs, and addons.
func TestManualCombinationsMatrixValidation(t *testing.T) {
	testCases := []struct {
		name string
		cfg  types.ScaffoldConfig
	}{
		{
			name: "NextJS_Express_Postgres_Drizzle_BetterAuth_Shadcn",
			cfg: types.ScaffoldConfig{
				ProjectName: "koko-fullstack-1",
				Frontend:    "nextjs",
				Backend:     "express",
				Database:    "postgres",
				ORM:         "drizzle",
				Auth:        "better-auth",
				Addons:      "docker,shadcn,zod,lucide",
				API:         "trpc",
			},
		},
		{
			name: "React_Hono_MySQL_Prisma_Clerk",
			cfg: types.ScaffoldConfig{
				ProjectName: "koko-fullstack-2",
				Frontend:    "react",
				Backend:     "hono",
				Database:    "mysql",
				ORM:         "prisma",
				Auth:        "clerk",
				Addons:      "docker,github_actions",
				API:         "orpc",
			},
		},
		{
			name: "Astro_NestJS_MongoDB_Mongoose_Docker",
			cfg: types.ScaffoldConfig{
				ProjectName: "koko-fullstack-3",
				Frontend:    "astro",
				Backend:     "nestjs",
				Database:    "mongodb",
				ORM:         "mongoose",
				Addons:      "docker",
			},
		},
		{
			name: "Nuxt_GoChi_SQLite_Drizzle",
			cfg: types.ScaffoldConfig{
				ProjectName: "koko-fullstack-4",
				Frontend:    "nuxt",
				Backend:     "go_chi",
				Database:    "sqlite",
				ORM:         "drizzle",
			},
		},
		{
			name: "Svelte_FastAPI_Postgres_SQLAlchemy",
			cfg: types.ScaffoldConfig{
				ProjectName: "koko-fullstack-5",
				Frontend:    "svelte",
				Backend:     "fastapi",
				Database:    "postgres",
				ORM:         "sqlalchemy",
				Addons:      "docker,github_actions",
			},
		},
		{
			name: "NextJS_FastAPI_Postgres_Prisma_BetterAuth",
			cfg: types.ScaffoldConfig{
				ProjectName: "koko-fullstack-6",
				Frontend:    "nextjs",
				Backend:     "fastapi",
				Database:    "postgres",
				ORM:         "prisma",
				Auth:        "better-auth",
				Addons:      "docker,shadcn,motion",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "koko-test-manual-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			if err := RunScaffold(tmpDir, tc.cfg); err != nil {
				t.Fatalf("Scaffold failed for manual case %s: %v", tc.name, err)
			}

			validateProjectTree(t, tmpDir, tc.cfg.ProjectName)
		})
	}
}
