package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/errors"
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
		return errors.NewValidationError("El nombre del proyecto no puede estar vacío", "Escribe un nombre, por ejemplo: 'mi-app' o 'my-saas'")
	}

	if len(trimmed) > 214 {
		return errors.NewValidationError("El nombre del proyecto no puede superar los 214 caracteres", "Elige un nombre más corto y descriptivo")
	}

	if strings.ContainsAny(trimmed, " \t\n\r") {
		return errors.NewValidationError("El nombre del proyecto no puede contener espacios", fmt.Sprintf("Usa guiones medios en su lugar: '%s'", strings.ReplaceAll(trimmed, " ", "-")))
	}

	if strings.ToLower(trimmed) != trimmed {
		return errors.NewValidationError("El nombre del proyecto debe estar en minúsculas", fmt.Sprintf("Usa: '%s'", strings.ToLower(trimmed)))
	}

	if strings.HasPrefix(trimmed, ".") || strings.HasPrefix(trimmed, "_") {
		return errors.NewValidationError("El nombre no puede comenzar con un punto '.' o guion bajo '_'", "Comienza con una letra o número (ej: 'app-web')")
	}

	if !validNameRegex.MatchString(trimmed) {
		return errors.NewValidationError("El nombre contiene caracteres especiales no permitidos", "Usa únicamente letras minúsculas, números, guiones (-) y puntos (.)")
	}

	if reservedNames[strings.ToLower(trimmed)] {
		return errors.NewValidationError(fmt.Sprintf("'%s' es un nombre reservado por el sistema", trimmed), "Elige otro nombre para tu proyecto (ej: 'mi-"+trimmed+"')")
	}

	return nil
}

// ValidateTargetDir checks if the target directory already exists and is non-empty.
func ValidateTargetDir(dirPath string) error {
	trimmed := strings.TrimSpace(dirPath)
	if trimmed == "" {
		return errors.NewValidationError("La ruta del directorio de destino no puede estar vacía", "Especifica un directorio válido")
	}

	info, err := os.Stat(trimmed)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist, which is ideal
		}
		return errors.NewValidationError(fmt.Sprintf("No se pudo inspeccionar el directorio '%s'", trimmed), "Verifica los permisos de lectura en la carpeta actual")
	}

	if !info.IsDir() {
		return errors.NewValidationError(fmt.Sprintf("Ya existe un archivo llamado '%s'", trimmed), "Elige otro nombre o elimina el archivo existente")
	}

	// Check if directory is empty
	entries, err := os.ReadDir(trimmed)
	if err != nil {
		return errors.NewValidationError(fmt.Sprintf("No se pudo leer el contenido de '%s'", trimmed), "Verifica los permisos de lectura")
	}

	if len(entries) > 0 {
		return errors.NewValidationError(fmt.Sprintf("El directorio '%s' ya existe y no está vacío", filepath.Base(trimmed)), "Elige un directorio nuevo o vacía la carpeta antes de continuar")
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
