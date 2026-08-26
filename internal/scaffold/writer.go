package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/errors"
	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

// assertSafeWritePath previene ataques de path traversal y escrituras accidentales
// fuera de la carpeta designada para el proyecto.
func assertSafeWritePath(targetDir, destPath string) error {
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return errors.NewDiskWriteError(destPath, "No se pudo resolver la ruta base del proyecto", err, "Verifica los permisos del sistema")
	}

	absDest, err := filepath.Abs(destPath)
	if err != nil {
		return errors.NewDiskWriteError(destPath, "No se pudo resolver la ruta destino del archivo", err, "Verifica los caracteres en el nombre")
	}

	rel, err := filepath.Rel(absTarget, absDest)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return errors.NewDiskWriteError(destPath, fmt.Sprintf("Ruta insegura fuera del proyecto (%s)", destPath), nil, "No se permite escribir fuera del directorio del proyecto")
	}

	return nil
}

// WriteTree vuelca todos los archivos del VFS a disco de manera segura y atómica.
// Si ocurre un error a mitad de camino, ejecuta un rollback automático para no dejar archivos huérfanos.
func WriteTree(virtualFS *vfs.VFS, targetDir string) error {
	if virtualFS.FileCount() == 0 {
		return errors.NewDiskWriteError(targetDir, "No hay archivos en memoria para escribir", nil, "Verifica la selección de plantillas")
	}

	targetDirAbs, err := filepath.Abs(targetDir)
	if err != nil {
		return errors.NewDiskWriteError(targetDir, "Error al resolver directorio de destino", err, "Verifica la ruta ingresada")
	}

	// Verificar si el directorio destino ya existía antes de empezar
	dirExistedPrior := false
	if _, err := os.Stat(targetDirAbs); err == nil {
		dirExistedPrior = true
	}

	// Registrar archivos escritos para permitir rollback en caso de fallo
	writtenFiles := make([]string, 0, virtualFS.FileCount())

	rollback := func() {
		// Eliminar archivos escritos
		for _, f := range writtenFiles {
			_ = os.Remove(f)
		}
		// Si el directorio principal fue creado por este proceso y falló, eliminarlo completamente
		if !dirExistedPrior {
			_ = os.RemoveAll(targetDirAbs)
		}
	}

	// Crear carpeta raíz
	if err := os.MkdirAll(targetDirAbs, 0755); err != nil {
		return errors.NewDiskWriteError(targetDirAbs, "No se pudo crear el directorio del proyecto", err, "Revisa los permisos de escritura en la carpeta actual")
	}

	for _, relPath := range virtualFS.ListFiles() {
		destPath := filepath.Join(targetDirAbs, filepath.FromSlash(relPath))

		if err := assertSafeWritePath(targetDirAbs, destPath); err != nil {
			rollback()
			return err
		}

		// Asegurar que exista el subdirectorio
		parentDir := filepath.Dir(destPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			rollback()
			return errors.NewDiskWriteError(destPath, "No se pudo crear la estructura de carpetas", err, "Revisa los permisos de disco")
		}

		content, _ := virtualFS.ReadFile(relPath)
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			rollback()
			return errors.NewDiskWriteError(destPath, "Fallo al escribir archivo en disco", err, "Verifica espacio en disco o permisos")
		}

		writtenFiles = append(writtenFiles, destPath)
	}

	return nil
}
