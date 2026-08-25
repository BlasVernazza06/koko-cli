package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantErr     bool
	}{
		{"valid simple", "my-app", false},
		{"valid with numbers and underscore", "koko_starter_123", false},
		{"valid with dots", "my.cool.project", false},
		{"empty string", "", true},
		{"only spaces", "   ", true},
		{"contains uppercase", "MyApp", true},
		{"contains spaces", "my app", true},
		{"starts with dot", ".my-app", true},
		{"starts with underscore", "_my-app", true},
		{"contains invalid characters", "my$app", true},
		{"contains slashes", "my/app", true},
		{"reserved node_modules", "node_modules", true},
		{"reserved windows aux", "aux", true},
		{"reserved windows con", "CON", true},
		{"reserved windows nul", "nul", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName(%q) error = %v, wantErr %v", tt.projectName, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTargetDir(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Non-existent directory should be valid
	nonExistent := filepath.Join(tmpDir, "new-project")
	if err := ValidateTargetDir(nonExistent); err != nil {
		t.Errorf("expected non-existent directory to be valid, got: %v", err)
	}

	// 2. Existing empty directory should be valid
	emptyDir := filepath.Join(tmpDir, "empty-dir")
	if err := os.Mkdir(emptyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTargetDir(emptyDir); err != nil {
		t.Errorf("expected empty directory to be valid, got: %v", err)
	}

	// 3. Existing non-empty directory should fail
	nonEmptyDir := filepath.Join(tmpDir, "non-empty-dir")
	if err := os.Mkdir(nonEmptyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nonEmptyDir, "file.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTargetDir(nonEmptyDir); err == nil {
		t.Errorf("expected non-empty directory to fail validation, got nil")
	}

	// 4. Existing file with same name should fail
	filePath := filepath.Join(tmpDir, "existing-file")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateTargetDir(filePath); err == nil {
		t.Errorf("expected existing file to fail validation, got nil")
	}
}
