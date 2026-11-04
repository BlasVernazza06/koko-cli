package doctor

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/BlasVernazza06/koko-cli/internal/core/config"
)

// DiagnosticResult contiene el informe completo del diagnóstico del proyecto.
type DiagnosticResult struct {
	ConfigExists   bool
	ExistingConfig *config.KokoConfig
	DetectedConfig config.KokoConfig
	Differences    []Difference
	HealthReport   *HealthCheckReport
}

// RunDiagnostics analiza el directorio del proyecto y calcula las diferencias respecto a koko.config.json,
// ejecutando además el health check extendido de 5 módulos.
func RunDiagnostics(projectDir string) (*DiagnosticResult, error) {
	configPath := filepath.Join(projectDir, "koko.config.json")
	var existing *config.KokoConfig

	if fileExists(configPath) {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

		var cfg config.KokoConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
		existing = &cfg
	}

	detected := AnalyzeProject(projectDir)
	diffs := ComputeDifferences(existing, detected)
	healthReport := RunFullHealthCheck(context.Background(), projectDir, detected)

	return &DiagnosticResult{
		ConfigExists:   existing != nil,
		ExistingConfig: existing,
		DetectedConfig: detected,
		Differences:    diffs,
		HealthReport:   healthReport,
	}, nil
}

// ApplyFixes actualiza o crea el archivo koko.config.json con la configuración detectada
// y repara archivos faltantes como .env a partir de .env.example.
func ApplyFixes(projectDir string, detected config.KokoConfig, existing *config.KokoConfig) error {
	finalConfig := detected
	finalConfig.Schema = "https://koko-cli.dev/schema.json"
	finalConfig.Project.CLIVersion = config.CLIVersion

	if existing != nil {
		finalConfig.Project.Name = existing.Project.Name
		finalConfig.Project.CreatedAt = existing.Project.CreatedAt
	} else {
		absDir, err := filepath.Abs(projectDir)
		if err == nil {
			finalConfig.Project.Name = filepath.Base(absDir)
		} else {
			finalConfig.Project.Name = "koko-project"
		}
		finalConfig.Project.CreatedAt = time.Now().UTC()
	}

	data, err := json.MarshalIndent(finalConfig, "", "  ")
	if err != nil {
		return err
	}

	targetFile := filepath.Join(projectDir, "koko.config.json")
	if err := os.WriteFile(targetFile, data, 0644); err != nil {
		return err
	}

	// Auto-fix de archivos .env a partir de .env.example si aplica
	_, _ = FixEnvFiles(projectDir)

	return nil
}
