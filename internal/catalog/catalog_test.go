package catalog

import (
	"strings"
	"testing"
)

func TestDependencyCatalogIntegrity(t *testing.T) {
	if len(DependencyVersions) == 0 {
		t.Fatal("Dependency catalog is empty")
	}

	for pkg, ver := range DependencyVersions {
		if strings.TrimSpace(pkg) == "" {
			t.Errorf("Catalog contains empty package name")
		}
		if strings.TrimSpace(ver) == "" {
			t.Errorf("Catalog package %s has empty version", pkg)
		}

		// Ensure no accidental brackets or template tags in catalog keys or values
		if strings.Contains(pkg, "[") || strings.Contains(pkg, "]") {
			t.Errorf("Catalog key %q contains invalid characters", pkg)
		}
		if strings.Contains(ver, "[") || strings.Contains(ver, "]") {
			t.Errorf("Catalog version for %s (%q) contains invalid characters", pkg, ver)
		}
	}
}

func TestGetVersion(t *testing.T) {
	nextVer := GetVersion("next")
	if nextVer == "latest" || nextVer == "" {
		t.Errorf("Expected defined version for next, got: %s", nextVer)
	}

	nonExistent := GetVersion("non-existent-pkg-12345", "^1.0.0")
	if nonExistent != "^1.0.0" {
		t.Errorf("Expected fallback version ^1.0.0, got: %s", nonExistent)
	}

	defaultFallback := GetVersion("non-existent-pkg-67890")
	if defaultFallback != "latest" {
		t.Errorf("Expected default fallback 'latest', got: %s", defaultFallback)
	}
}
