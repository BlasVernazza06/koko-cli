package scaffold

// ¡Este archivo es tuyo para programar!
// Aquí debes implementar la lógica que lea las plantillas de la carpeta 'templates'
// y las escriba en el directorio destino del usuario.
//
// Tips para empezar:
// 1. Define tu estructura de configuración (por ejemplo, 'ScaffoldConfig').
// 2. Investiga cómo usar la directiva '//go:embed' para importar la carpeta de plantillas.
// 3. Escribe una función 'RunScaffold(targetDir string, config ScaffoldConfig) error'.
import (
	"embed"
	"fmt"
	"os"
)

//go:embed all:templates
var templateFs embed.FS

type ScaffoldConfig struct {
	ProjectName   string
	Frontend      string
	Backend       string
	Database      string
	Docker        bool
	GithubActions bool
}

func RunScaffold(targetDir string, config ScaffoldConfig) error {
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("error al crear el directorio de destino: %w", err)
	}

	return nil
}
