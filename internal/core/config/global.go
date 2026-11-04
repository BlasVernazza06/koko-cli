package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GlobalSettings almacena las preferencias globales del usuario para Koko CLI.
type GlobalSettings struct {
	DefaultProjectsDirectory string `json:"defaultProjectsDirectory,omitempty"`
}

func getGlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".koko")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

// LoadGlobalSettings lee las preferencias de ~/.koko/settings.json.
func LoadGlobalSettings() (*GlobalSettings, error) {
	path, err := getGlobalConfigPath()
	if err != nil {
		return &GlobalSettings{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &GlobalSettings{}, nil
		}
		return &GlobalSettings{}, err
	}

	var settings GlobalSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return &GlobalSettings{}, nil
	}

	return &settings, nil
}

// SaveGlobalSettings guarda las preferencias en ~/.koko/settings.json.
func SaveGlobalSettings(settings *GlobalSettings) error {
	path, err := getGlobalConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ResolveDefaultTargetDirectory determina la ruta de destino según la jerarquía:
// 1. Si targetDir explícito fue proporcionado, usarlo.
// 2. Si DefaultProjectsDirectory está configurado en settings.json, usarlo.
// 3. Fallback: Directorio Desktop del usuario (~/Desktop). Si no existe, UserHomeDir.
func ResolveDefaultTargetDirectory(explicitTargetDir string) string {
	if explicitTargetDir != "" {
		return explicitTargetDir
	}

	settings, err := LoadGlobalSettings()
	if err == nil && settings.DefaultProjectsDirectory != "" {
		if _, err := os.Stat(settings.DefaultProjectsDirectory); err == nil {
			return settings.DefaultProjectsDirectory
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		desktop := filepath.Join(home, "Desktop")
		if info, err := os.Stat(desktop); err == nil && info.IsDir() {
			return desktop
		}
		return home
	}

	return "."
}
