package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Package manager conventions regex:
	// - Must only contain lowercase alphanumeric characters, hyphens (-), underscores (_), or dots (.)
	// - Must not start with . or _
	validNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

	// Reserved package names or problematic names
	reservedNames = map[string]bool{
		"node_modules": true,
		"favicon.ico":  true,
		"package.json": true,
		"tsconfig":     true,
		// Windows device reserved names
		"con":  true,
		"prn":  true,
		"aux":  true,
		"nul":  true,
		"com1": true,
		"com2": true,
		"com3": true,
		"com4": true,
		"com5": true,
		"com6": true,
		"com7": true,
		"com8": true,
		"com9": true,
		"lpt1": true,
		"lpt2": true,
		"lpt3": true,
		"lpt4": true,
		"lpt5": true,
		"lpt6": true,
		"lpt7": true,
		"lpt8": true,
		"lpt9": true,
	}
)

// ValidateProjectName checks that the project name complies with package manager conventions.
func ValidateProjectName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(trimmed) > 214 {
		return fmt.Errorf("project name must be 214 characters or less")
	}

	if strings.ContainsAny(trimmed, " \t\n\r") {
		return fmt.Errorf("project name cannot contain spaces")
	}

	if strings.ToLower(trimmed) != trimmed {
		return fmt.Errorf("project name must be lowercase (e.g. '%s')", strings.ToLower(trimmed))
	}

	if strings.HasPrefix(trimmed, ".") || strings.HasPrefix(trimmed, "_") {
		return fmt.Errorf("project name cannot start with a dot '.' or underscore '_'")
	}

	if !validNameRegex.MatchString(trimmed) {
		return fmt.Errorf("project name can only contain lowercase letters, numbers, hyphens (-), underscores (_), and dots (.)")
	}

	if reservedNames[strings.ToLower(trimmed)] {
		return fmt.Errorf("'%s' is a reserved name and cannot be used", trimmed)
	}

	return nil
}

// ValidateTargetDir checks if the target directory already exists and is non-empty.
func ValidateTargetDir(dirPath string) error {
	trimmed := strings.TrimSpace(dirPath)
	if trimmed == "" {
		return fmt.Errorf("target directory path cannot be empty")
	}

	info, err := os.Stat(trimmed)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist, which is ideal
		}
		return fmt.Errorf("unable to inspect directory '%s': %w", trimmed, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("a file named '%s' already exists", trimmed)
	}

	// Check if directory is empty
	entries, err := os.ReadDir(trimmed)
	if err != nil {
		return fmt.Errorf("unable to read directory '%s': %w", trimmed, err)
	}

	if len(entries) > 0 {
		return fmt.Errorf("directory '%s' already exists and is not empty", filepath.Base(trimmed))
	}

	return nil
}

// Validate performs both package manager name validation and directory existence check.
func Validate(name string) error {
	if err := ValidateProjectName(name); err != nil {
		return err
	}
	if err := ValidateTargetDir(name); err != nil {
		return err
	}
	return nil
}
