package db

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// OSExecutor implementa Executor ejecutando procesos de forma segura en el sistema operativo.
type OSExecutor struct {
	LookPath func(file string) (string, error)
}

// NewOSExecutor crea un OSExecutor estándar.
func NewOSExecutor() *OSExecutor {
	return &OSExecutor{
		LookPath: exec.LookPath,
	}
}

// Execute corre el comando resuelto conectando los streams I/O.
func (e *OSExecutor) Execute(ctx context.Context, resCmd ResolvedCommand, stdin io.Reader, stdout, stderr io.Writer) error {
	exeName := resCmd.Executable

	// En Windows, comandos como pnpm, npm o bunx pueden requerir resolución de extensiones (.cmd, .ps1, .bat)
	// exec.LookPath resuelve automáticamente "pnpm.cmd" en Windows.
	lookPathFn := e.LookPath
	if lookPathFn == nil {
		lookPathFn = exec.LookPath
	}

	targetExe, err := lookPathFn(exeName)
	if err != nil {
		if runtime.GOOS == "windows" {
			// Intentar con sufijo .cmd
			if !strings.HasSuffix(strings.ToLower(exeName), ".cmd") &&
				!strings.HasSuffix(strings.ToLower(exeName), ".bat") &&
				!strings.HasSuffix(strings.ToLower(exeName), ".exe") {
				if candidate, errCmd := lookPathFn(exeName + ".cmd"); errCmd == nil {
					targetExe = candidate
					err = nil
				}
			}
		}
		if err != nil {
			return fmt.Errorf("executable '%s' not found in system PATH: %w", exeName, err)
		}
	}

	// Crear comando con contexto sin usar shell intermediario (evita inyección de comandos)
	cmd := exec.CommandContext(ctx, targetExe, resCmd.Args...)

	// Configurar directorio de trabajo limpio
	if resCmd.WorkDir != "" {
		cleanWorkDir := filepath.Clean(resCmd.WorkDir)
		cmd.Dir = cleanWorkDir
	}

	// Conectar flujos de I/O
	if stdin != nil {
		cmd.Stdin = stdin
	} else {
		cmd.Stdin = os.Stdin
	}

	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = os.Stdout
	}

	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	// Propagar variables de entorno actuales del proceso padre
	cmd.Env = os.Environ()

	return cmd.Run()
}

// FormatCommandSummary devuelve una representación amigable del comando para consola / dry-run.
func FormatCommandSummary(resCmd *ResolvedCommand) string {
	parts := make([]string, 0, len(resCmd.Args)+1)
	parts = append(parts, resCmd.Executable)
	for _, arg := range resCmd.Args {
		if strings.Contains(arg, " ") {
			parts = append(parts, fmt.Sprintf("%q", arg))
		} else {
			parts = append(parts, arg)
		}
	}
	return strings.Join(parts, " ")
}
