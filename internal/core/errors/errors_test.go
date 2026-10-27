package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestScaffoldErrorFormatting(t *testing.T) {
	err := NewValidationError("El nombre no es válido", "Usa letras minúsculas")
	if !strings.Contains(err.Error(), "[Validación de entrada]") {
		t.Errorf("Expected phase in error string, got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "💡 Sugerencia: Usa letras minúsculas") {
		t.Errorf("Expected suggestion in error string, got: %s", err.Error())
	}

	wrapped := errors.New("io error")
	tmplErr := NewTemplateError("apps/web/package.json", "Fallo al renderizar", wrapped, "Revisa la sintaxis")
	if !strings.Contains(tmplErr.Error(), "apps/web/package.json") {
		t.Errorf("Expected file in error string, got: %s", tmplErr.Error())
	}
	if !errors.Is(tmplErr, wrapped) {
		t.Errorf("Expected Unwrap to support errors.Is")
	}
}
