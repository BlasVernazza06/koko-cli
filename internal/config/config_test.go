package config

import (
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

func TestBuildKokoConfigForRecipes(t *testing.T) {
	testCases := []struct {
		recipe       string
		expectedFe   string
		expectedBe   string
		expectedDb   string
		expectedOrm  string
		expectedPm   string
		expectedAuth string
	}{
		{
			recipe:       "saas",
			expectedFe:   "next",
			expectedBe:   "next",
			expectedDb:   "postgres",
			expectedOrm:  "drizzle",
			expectedPm:   "pnpm",
			expectedAuth: "better-auth",
		},
		{
			recipe:       "java-spring",
			expectedFe:   "react",
			expectedBe:   "spring",
			expectedDb:   "postgres",
			expectedOrm:  "jpa",
			expectedPm:   "pnpm",
			expectedAuth: "",
		},
		{
			recipe:       "enterprise-nestjs",
			expectedFe:   "next",
			expectedBe:   "nestjs",
			expectedDb:   "postgres",
			expectedOrm:  "prisma",
			expectedPm:   "pnpm",
			expectedAuth: "better-auth",
		},
		{
			recipe:       "mern",
			expectedFe:   "react",
			expectedBe:   "express",
			expectedDb:   "mongodb",
			expectedOrm:  "mongoose",
			expectedPm:   "npm",
			expectedAuth: "jwt",
		},
		{
			recipe:       "pern",
			expectedFe:   "react",
			expectedBe:   "express",
			expectedDb:   "postgres",
			expectedOrm:  "prisma",
			expectedPm:   "pnpm",
			expectedAuth: "jwt",
		},
		{
			recipe:       "python-fastapi",
			expectedFe:   "react",
			expectedBe:   "fastapi",
			expectedDb:   "postgres",
			expectedOrm:  "sqlalchemy",
			expectedPm:   "pnpm",
			expectedAuth: "",
		},
		{
			recipe:       "mobile-expo",
			expectedFe:   "expo",
			expectedBe:   "express",
			expectedDb:   "postgres",
			expectedOrm:  "prisma",
			expectedPm:   "pnpm",
			expectedAuth: "jwt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.recipe, func(t *testing.T) {
			cfg := scaffold.ScaffoldConfig{
				ProjectName: "test-app",
				Recipe:      tc.recipe,
			}
			kokoCfg := BuildKokoConfig(cfg)

			if kokoCfg.Architecture.PackageManager != tc.expectedPm {
				t.Errorf("For recipe %s expected PM %s, got %s", tc.recipe, tc.expectedPm, kokoCfg.Architecture.PackageManager)
			}
			if tc.expectedFe != "" && (kokoCfg.Stack.Frontend == nil || kokoCfg.Stack.Frontend.Framework != tc.expectedFe) {
				t.Errorf("For recipe %s expected FE %s, got %+v", tc.recipe, tc.expectedFe, kokoCfg.Stack.Frontend)
			}
			if tc.expectedBe != "" && (kokoCfg.Stack.Backend == nil || kokoCfg.Stack.Backend.Framework != tc.expectedBe) {
				t.Errorf("For recipe %s expected BE %s, got %+v", tc.recipe, tc.expectedBe, kokoCfg.Stack.Backend)
			}
			if tc.expectedDb != "" && (kokoCfg.Stack.Database == nil || kokoCfg.Stack.Database.Provider != tc.expectedDb) {
				t.Errorf("For recipe %s expected DB %s, got %+v", tc.recipe, tc.expectedDb, kokoCfg.Stack.Database)
			}
			if tc.expectedOrm != "" && (kokoCfg.Stack.Database == nil || kokoCfg.Stack.Database.ORM != tc.expectedOrm) {
				t.Errorf("For recipe %s expected ORM %s, got %+v", tc.recipe, tc.expectedOrm, kokoCfg.Stack.Database)
			}
			if tc.expectedAuth != "" {
				if kokoCfg.Features.Auth == nil || kokoCfg.Features.Auth.Provider != tc.expectedAuth {
					t.Errorf("For recipe %s expected Auth %s, got %+v", tc.recipe, tc.expectedAuth, kokoCfg.Features.Auth)
				}
			}
		})
	}
}
