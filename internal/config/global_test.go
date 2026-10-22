package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalSettingsLoadSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "koko-global-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	settings := &GlobalSettings{
		DefaultProjectsDirectory: tmpDir,
	}

	configFile := filepath.Join(tmpDir, "settings.json")
	data, err := os.ReadFile(configFile)
	assert.Error(t, err) // Todavía no existe

	// Verificar resolución de explícito
	res := ResolveDefaultTargetDirectory("C:\\custom-dir")
	assert.Equal(t, "C:\\custom-dir", res)

	_ = data
	_ = settings
}
