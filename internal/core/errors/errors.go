package errors

import (
	"fmt"
	"strings"
)

// PhaseType representa la etapa exacta del ciclo de vida de scaffolding donde ocurrió el error.
type PhaseType string

const (
	PhaseValidation     PhaseType = "Validación de entrada"
	PhaseTemplateRender PhaseType = "Renderizado de plantillas"
	PhasePostProcess    PhaseType = "Ajuste de configuración"
	PhaseDiskWrite      PhaseType = "Escritura en disco"
	PhaseUnknown        PhaseType = "General"
)

// ScaffoldError es un error estructurado y tipado que contiene contexto de la fase,
// el archivo involucrado y sugerencias claras de remediación para el usuario.
type ScaffoldError struct {
	Phase      PhaseType
	File       string
	Message    string
	Err        error
	Suggestion string
}

func (e *ScaffoldError) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s] %s", e.Phase, e.Message))
	if e.File != "" {
		b.WriteString(fmt.Sprintf(" (archivo: %s)", e.File))
	}
	if e.Err != nil {
		b.WriteString(fmt.Sprintf(": %v", e.Err))
	}
	if e.Suggestion != "" {
		b.WriteString(fmt.Sprintf("\n💡 Sugerencia: %s", e.Suggestion))
	}
	return b.String()
}

func (e *ScaffoldError) Unwrap() error {
	return e.Err
}

// NewValidationError crea un error de validación de configuración o nombres.
func NewValidationError(message string, suggestion string) *ScaffoldError {
	return &ScaffoldError{
		Phase:      PhaseValidation,
		Message:    message,
		Suggestion: suggestion,
	}
}

// NewTemplateError crea un error durante el procesamiento o renderizado de un template.
func NewTemplateError(file string, message string, err error, suggestion string) *ScaffoldError {
	return &ScaffoldError{
		Phase:      PhaseTemplateRender,
		File:       file,
		Message:    message,
		Err:        err,
		Suggestion: suggestion,
	}
}

// NewPostProcessError crea un error durante la mutación de package.json o .env.
func NewPostProcessError(file string, message string, err error) *ScaffoldError {
	return &ScaffoldError{
		Phase:   PhasePostProcess,
		File:    file,
		Message: message,
		Err:     err,
	}
}

// NewDiskWriteError crea un error durante la escritura física a disco.
func NewDiskWriteError(file string, message string, err error, suggestion string) *ScaffoldError {
	return &ScaffoldError{
		Phase:      PhaseDiskWrite,
		File:       file,
		Message:    message,
		Err:        err,
		Suggestion: suggestion,
	}
}
