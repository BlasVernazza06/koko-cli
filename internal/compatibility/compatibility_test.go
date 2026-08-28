package compatibility

import (
	"testing"

	"github.com/BlasVernazza06/koko-cli/cmd/views"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

func TestGetStepOptions_Backend(t *testing.T) {
	// If Frontend == "none", Backend "none" must be disabled
	selections := []views.SelectOption{
		{Value: "none"},
	}
	opts := GetStepOptions(StepBackend, selections)
	var noneOpt *views.SelectOption
	for _, opt := range opts {
		if opt.Value == "none" {
			noneOpt = &opt
			break
		}
	}
	if noneOpt == nil || !noneOpt.Disabled {
		t.Errorf("Expected backend 'none' to be disabled when frontend is 'none'")
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

	// If Next.js project
	selectionsNext := []views.SelectOption{
		{Value: "nextjs"},
		{Value: "none"},
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

func TestGetStepOptions_ORM_DatabaseNone(t *testing.T) {
	// If Database == "none", all ORMs except none must be disabled
	selections := []views.SelectOption{
		{Value: "react"},
		{Value: "express"},
		{Value: "pnpm"},
		{Value: "none"},
	}
	opts := GetStepOptions(StepORM, selections)
	for _, opt := range opts {
		if opt.Value == "none" {
			if opt.Disabled {
				t.Errorf("Expected ORM 'none' to be enabled")
			}
		} else {
			if !opt.Disabled {
				t.Errorf("Expected ORM '%s' to be disabled when Database is 'none'", opt.Value)
			}
		}
	}
}

func TestGetStepOptions_ORM_MongoDB(t *testing.T) {
	// Node.js + MongoDB: Mongoose and Prisma enabled; Drizzle, SQLAlchemy, GORM disabled
	selections := []views.SelectOption{
		{Value: "react"},
		{Value: "express"},
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
		case "drizzle", "sqlalchemy", "gorm":
			if !opt.Disabled {
				t.Errorf("Expected ORM '%s' to be disabled for MongoDB", opt.Value)
			}
		}
	}
}

func TestGetStepOptions_ORM_PythonFastAPI(t *testing.T) {
	// FastAPI + PostgreSQL: SQLAlchemy & None enabled; Drizzle, Prisma, Mongoose, GORM disabled
	selections := []views.SelectOption{
		{Value: "none"},
		{Value: "fastapi"},
		{Value: "pip"},
		{Value: "postgres"},
	}
	opts := GetStepOptions(StepORM, selections)
	for _, opt := range opts {
		switch opt.Value {
		case "sqlalchemy", "none":
			if opt.Disabled {
				t.Errorf("Expected ORM '%s' to be enabled for FastAPI + Postgres", opt.Value)
			}
		case "drizzle", "prisma", "moongose", "gorm":
			if !opt.Disabled {
				t.Errorf("Expected ORM '%s' to be disabled for Python/FastAPI", opt.Value)
			}
		}
	}
}

func TestGetStepOptions_ORM_GoChi(t *testing.T) {
	// Go Chi + MySQL: GORM & None enabled; others disabled
	selections := []views.SelectOption{
		{Value: "none"},
		{Value: "go_chi"},
		{Value: "go_mod"},
		{Value: "mysql"},
	}
	opts := GetStepOptions(StepORM, selections)
	for _, opt := range opts {
		switch opt.Value {
		case "gorm", "none":
			if opt.Disabled {
				t.Errorf("Expected ORM '%s' to be enabled for Go Chi + MySQL", opt.Value)
			}
		case "drizzle", "prisma", "moongose", "sqlalchemy":
			if !opt.Disabled {
				t.Errorf("Expected ORM '%s' to be disabled for Go Chi", opt.Value)
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
			name: "valid fullstack nextjs + drizzle + postgres",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "nextjs",
				Backend:  "none",
				Database: "postgres",
				ORM:      "drizzle",
			},
			wantErr: false,
		},
		{
			name: "valid fullstack nuxt + prisma + mysql",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "nuxt",
				Backend:  "none",
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
			name: "invalid svelte spa without backend with mongodb",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "svelte",
				Backend:  "none",
				Database: "mongodb",
				ORM:      "moongose",
			},
			wantErr: true,
		},
		{
			name: "invalid db none with orm drizzle",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "nextjs",
				Backend:  "none",
				Database: "none",
				ORM:      "drizzle",
			},
			wantErr: true,
		},
		{
			name: "invalid go with mongoose",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "none",
				Backend:  "go_chi",
				Database: "mongodb",
				ORM:      "moongose",
			},
			wantErr: true,
		},
		{
			name: "invalid python with drizzle",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "none",
				Backend:  "fastapi",
				Database: "postgres",
				ORM:      "drizzle",
			},
			wantErr: true,
		},
		{
			name: "invalid postgres with mongoose",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "react",
				Backend:  "express",
				Database: "postgres",
				ORM:      "moongose",
			},
			wantErr: true,
		},
		{
			name: "invalid mongodb with sqlalchemy",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "none",
				Backend:  "fastapi",
				Database: "mongodb",
				ORM:      "sqlalchemy",
			},
			wantErr: true,
		},
		{
			name: "invalid pure go with pnpm package manager",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "none",
				Backend:        "go_chi",
				PackageManager: "pnpm",
			},
			wantErr: true,
		},
		{
			name: "invalid pure python with bun package manager",
			cfg: scaffold.ScaffoldConfig{
				Frontend:       "none",
				Backend:        "fastapi",
				PackageManager: "bun",
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
			name: "invalid next-auth on react spa",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "react",
				Backend:  "express",
				Auth:     "next-auth",
			},
			wantErr: true,
		},
		{
			name: "valid better-auth on nextjs",
			cfg: scaffold.ScaffoldConfig{
				Frontend: "nextjs",
				Backend:  "none",
				Auth:     "better-auth",
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
