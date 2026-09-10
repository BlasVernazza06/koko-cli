package compatibility

import (
	"testing"

	"github.com/BlasVernazza06/koko-cli/cmd/views"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

func TestGetStepOptions_Backend(t *testing.T) {
	// If Frontend == "none", Backend "none" and "self" must be disabled
	selections := []views.SelectOption{
		{Value: "none"},
	}
	opts := GetStepOptions(StepBackend, selections)
	for _, opt := range opts {
		if opt.Value == "none" || opt.Value == "self" {
			if !opt.Disabled {
				t.Errorf("Expected backend '%s' to be disabled when frontend is 'none'", opt.Value)
			}
		}
	}
}

func TestGetStepOptions_Backend_SelfCondition(t *testing.T) {
	// Fullstack frontends: "self" must be enabled
	fullstacks := []string{"nextjs", "nuxt", "svelte", "astro"}
	for _, f := range fullstacks {
		selections := []views.SelectOption{
			{Value: f},
		}
		opts := GetStepOptions(StepBackend, selections)
		var selfOpt *views.SelectOption
		for _, opt := range opts {
			if opt.Value == "self" {
				selfOpt = &opt
				break
			}
		}
		if selfOpt == nil || selfOpt.Disabled {
			t.Errorf("Expected backend 'self' to be enabled for fullstack frontend '%s'", f)
		}
	}

	// Client-only / mobile frontends: "self" must be disabled
	clientOnly := []string{"react", "native", "none"}
	for _, f := range clientOnly {
		selections := []views.SelectOption{
			{Value: f},
		}
		opts := GetStepOptions(StepBackend, selections)
		var selfOpt *views.SelectOption
		for _, opt := range opts {
			if opt.Value == "self" {
				selfOpt = &opt
				break
			}
		}
		if selfOpt == nil || !selfOpt.Disabled {
			t.Errorf("Expected backend 'self' to be disabled for non-fullstack frontend '%s'", f)
		}
	}
}

func TestGetStepOptions_API(t *testing.T) {
	// If Node/TS ecosystem (Next.js + self), tRPC and oRPC enabled
	selectionsNode := []views.SelectOption{
		{Value: "nextjs"},
		{Value: "self"},
	}
	optsNode := GetStepOptions(StepAPI, selectionsNode)
	for _, opt := range optsNode {
		if opt.Disabled {
			t.Errorf("Expected API option '%s' to be enabled for Next.js", opt.Value)
		}
	}

	// If pure Go backend, tRPC and oRPC must be disabled
	selectionsGo := []views.SelectOption{
		{Value: "none"},
		{Value: "go_chi"},
	}
	optsGo := GetStepOptions(StepAPI, selectionsGo)
	for _, opt := range optsGo {
		if opt.Value == "trpc" || opt.Value == "orpc" {
			if !opt.Disabled {
				t.Errorf("Expected API option '%s' to be disabled for pure Go backend", opt.Value)
			}
		}
	}
}

func TestGetStepOptions_PackageManager(t *testing.T) {
	// If pure Go project
	selectionsGo := []views.SelectOption{
		{Value: "none"},
		{Value: "go_chi"},
	}
	optsGo := GetStepOptions(StepPackageManager, selectionsGo)
	if len(optsGo) != 1 || optsGo[0].Value != "go_mod" {
		t.Errorf("Expected go_mod for pure Go project, got: %+v", optsGo)
	}

	// If pure Python project
	selectionsPy := []views.SelectOption{
		{Value: "none"},
		{Value: "fastapi"},
	}
	optsPy := GetStepOptions(StepPackageManager, selectionsPy)
	if len(optsPy) == 0 || (optsPy[0].Value != "pip" && optsPy[0].Value != "uv") {
		t.Errorf("Expected pip/uv for pure Python project, got: %+v", optsPy)
	}

	// If pure Java Spring Boot project
	selectionsJava := []views.SelectOption{
		{Value: "none"},
		{Value: "spring_boot"},
	}
	optsJava := GetStepOptions(StepPackageManager, selectionsJava)
	if len(optsJava) != 2 || optsJava[0].Value != "mvn" {
		t.Errorf("Expected mvn/gradle for pure Java project, got: %+v", optsJava)
	}

	// If Next.js project
	selectionsNext := []views.SelectOption{
		{Value: "nextjs"},
		{Value: "self"},
	}
	optsNext := GetStepOptions(StepPackageManager, selectionsNext)
	if len(optsNext) != 3 || optsNext[0].Value != "pnpm" {
		t.Errorf("Expected pnpm/npm/bun for Next.js project, got: %+v", optsNext)
	}
}

func TestGetStepOptions_Database_ClientSPA(t *testing.T) {
	// React SPA without backend: All databases except none must be disabled
	selections := []views.SelectOption{
		{Value: "react"},
		{Value: "none"},
		{Value: "none"},
		{Value: "pnpm"},
	}
	opts := GetStepOptions(StepDatabase, selections)
	for _, opt := range opts {
		if opt.Value == "none" {
			if opt.Disabled {
				t.Errorf("Expected database 'none' to be enabled for SPA without backend")
			}
		} else {
			if !opt.Disabled {
				t.Errorf("Expected database '%s' to be disabled for SPA without backend", opt.Value)
			}
		}
	}
}

func TestGetStepOptions_Database_SelfBackend(t *testing.T) {
	// Next.js with "self" backend: Databases must be enabled
	selections := []views.SelectOption{
		{Value: "nextjs"},
		{Value: "self"},
		{Value: "none"},
		{Value: "pnpm"},
	}
	opts := GetStepOptions(StepDatabase, selections)
	var pgOpt *views.SelectOption
	for _, opt := range opts {
		if opt.Value == "postgres" {
			pgOpt = &opt
			break
		}
	}
	if pgOpt == nil || pgOpt.Disabled {
		t.Errorf("Expected database 'postgres' to be enabled for Next.js + self backend")
	}
}

func TestGetStepOptions_Auth(t *testing.T) {
	// Next.js project: all auth options enabled
	selectionsNext := []views.SelectOption{
		{Value: "nextjs"},
		{Value: "self"},
		{Value: "none"},
		{Value: "pnpm"},
		{Value: "postgres"},
		{Value: "drizzle"},
	}
	optsNext := GetStepOptions(StepAuth, selectionsNext)
	for _, opt := range optsNext {
		if opt.Disabled {
			t.Errorf("Expected auth '%s' to be enabled for Next.js", opt.Value)
		}
	}

	// React + Express: NextAuth disabled, Better Auth / Clerk enabled
	selectionsReact := []views.SelectOption{
		{Value: "react"},
		{Value: "express"},
		{Value: "none"},
		{Value: "pnpm"},
		{Value: "postgres"},
		{Value: "drizzle"},
	}
	optsReact := GetStepOptions(StepAuth, selectionsReact)
	for _, opt := range optsReact {
		if opt.Value == "next-auth" && !opt.Disabled {
			t.Errorf("Expected NextAuth to be disabled for React SPA")
		}
		if (opt.Value == "better-auth" || opt.Value == "clerk") && opt.Disabled {
			t.Errorf("Expected '%s' to be enabled for React + Express", opt.Value)
		}
	}

	// Pure Go backend: TS auth disabled
	selectionsGo := []views.SelectOption{
		{Value: "none"},
		{Value: "go_chi"},
		{Value: "none"},
		{Value: "go_mod"},
		{Value: "postgres"},
		{Value: "gorm"},
	}
	optsGo := GetStepOptions(StepAuth, selectionsGo)
	for _, opt := range optsGo {
		if opt.Value != "none" && !opt.Disabled {
			t.Errorf("Expected auth '%s' to be disabled for pure Go backend", opt.Value)
		}
	}
}

func TestGetStepOptions_ORM_MongoDB(t *testing.T) {
	// Node.js + MongoDB: Mongoose and Prisma enabled; Drizzle, SQLAlchemy, GORM, JPA disabled
	selections := []views.SelectOption{
		{Value: "react"},
		{Value: "express"},
		{Value: "none"},
		{Value: "pnpm"},
		{Value: "mongodb"},
	}
	opts := GetStepOptions(StepORM, selections)
	for _, opt := range opts {
		switch opt.Value {
		case "moongose", "prisma", "none":
			if opt.Disabled {
				t.Errorf("Expected ORM '%s' to be enabled for Node.js + MongoDB", opt.Value)
			}
		case "drizzle", "sqlalchemy", "gorm", "jpa":
			if !opt.Disabled {
				t.Errorf("Expected ORM '%s' to be disabled for MongoDB", opt.Value)
			}
		}
	}
}

func TestGetStepOptions_ORM_Java(t *testing.T) {
	// Java backend: JPA and none enabled, other ORMs disabled
	selections := []views.SelectOption{
		{Value: "react"},
		{Value: "spring_boot"},
		{Value: "none"},
		{Value: "pnpm"},
		{Value: "postgres"},
	}
	opts := GetStepOptions(StepORM, selections)
	for _, opt := range opts {
		switch opt.Value {
		case "jpa", "none":
			if opt.Disabled {
				t.Errorf("Expected ORM '%s' to be enabled for Java Spring Boot", opt.Value)
			}
		default:
			if !opt.Disabled {
				t.Errorf("Expected ORM '%s' to be disabled for Java Spring Boot", opt.Value)
			}
		}
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     scaffold.ScaffoldConfig
		wantErr bool
	}{
		{
			name: "valid fullstack nextjs + self backend + drizzle + postgres + better-auth + trpc",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "nextjs",
				Backend:        "self",
				API:            "trpc",
				PackageManager: "pnpm",
				Database:       "postgres",
				ORM:            "drizzle",
				Auth:           "better-auth",
			},
			wantErr: false,
		},
		{
			name: "valid fullstack astro + self backend + postgres + drizzle",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "astro",
				Backend:        "self",
				PackageManager: "pnpm",
				Database:       "postgres",
				ORM:            "drizzle",
				Auth:           "none",
			},
			wantErr: false,
		},
		{
			name: "valid react + java spring boot + postgres + jpa",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "react",
				Backend:        "spring_boot",
				PackageManager: "pnpm",
				Database:       "postgres",
				ORM:            "jpa",
				Auth:           "none",
			},
			wantErr: false,
		},
		{
			name: "invalid trpc on pure java backend",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "none",
				Backend:        "spring_boot",
				API:            "trpc",
				PackageManager: "mvn",
			},
			wantErr: true,
		},
		{
			name: "invalid react spa with self backend",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "react",
				Backend:        "self",
				PackageManager: "pnpm",
			},
			wantErr: true,
		},
		{
			name: "invalid native expo with self backend",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "native",
				Backend:        "self",
				PackageManager: "pnpm",
			},
			wantErr: true,
		},
		{
			name: "valid fullstack nuxt + prisma + mysql",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "nuxt",
				Backend:  "self",
				Database: "mysql",
				ORM:      "prisma",
			},
			wantErr: false,
		},
		{
			name: "valid react + express + postgres + drizzle",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "react",
				Backend:  "express",
				Database: "postgres",
				ORM:      "drizzle",
			},
			wantErr: false,
		},
		{
			name: "valid pure python fastapi + postgres + sqlalchemy",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "none",
				Backend:        "fastapi",
				PackageManager: "pip",
				Database:       "postgres",
				ORM:            "sqlalchemy",
			},
			wantErr: false,
		},
		{
			name: "valid pure go go_chi + sqlite + gorm",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "none",
				Backend:        "go_chi",
				PackageManager: "go_mod",
				Database:       "sqlite",
				ORM:            "gorm",
			},
			wantErr: false,
		},
		{
			name: "invalid none + none",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "none",
				Backend:  "none",
			},
			wantErr: true,
		},
		{
			name: "invalid react spa without backend with postgres db",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "react",
				Backend:  "none",
				Database: "postgres",
				ORM:      "drizzle",
			},
			wantErr: true,
		},
		{
			name: "invalid db none with orm drizzle",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "nextjs",
				Backend:  "self",
				Database: "none",
				ORM:      "drizzle",
			},
			wantErr: true,
		},
		{
			name: "invalid better-auth on pure go backend",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "none",
				Backend:        "go_chi",
				PackageManager: "go_mod",
				Auth:           "better-auth",
			},
			wantErr: true,
		},
		{
			name: "invalid clerk on pure python backend",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "none",
				Backend:        "fastapi",
				PackageManager: "pip",
				Auth:           "clerk",
			},
			wantErr: true,
		},
		{
			name: "invalid next-auth on react spa",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "react",
				Backend:  "express",
				Auth:     "next-auth",
			},
			wantErr: true,
		},
		{
			name: "valid clerk on nextjs",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "nextjs",
				Backend:  "self",
				Auth:     "clerk",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

