package vfs

import (
	"bytes"
	"encoding/json"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// VFS representa un sistema de archivos virtual en memoria.
// Todas las operaciones de creación, edición y orquestación de templates
// ocurren sobre esta estructura antes de ser volcadas a disco.
type VFS struct {
	mu    sync.RWMutex
	files map[string][]byte
}

// New crea una nueva instancia limpia de Virtual File System.
func New() *VFS {
	return &VFS{
		files: make(map[string][]byte),
	}
}

// normalizePath limpia y estandariza la ruta a formato UNIX ("/") relativo.
func normalizePath(p string) string {
	clean := filepath.ToSlash(filepath.Clean(p))
	clean = strings.TrimPrefix(clean, "./")
	clean = strings.TrimPrefix(clean, "/")
	return clean
}

// WriteFile escribe o sobreescribe un archivo en memoria.
func (v *VFS) WriteFile(filePath string, content []byte) {
	v.mu.Lock()
	defer v.mu.Unlock()

	norm := normalizePath(filePath)
	if norm == "" || norm == "." {
		return
	}
	cp := make([]byte, len(content))
	copy(cp, content)
	v.files[norm] = cp
}

// WriteString escribe un string como archivo en memoria.
func (v *VFS) WriteString(filePath string, content string) {
	v.WriteFile(filePath, []byte(content))
}

// ReadFile obtiene los bytes de un archivo en memoria.
func (v *VFS) ReadFile(filePath string) ([]byte, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	data, exists := v.files[normalizePath(filePath)]
	if !exists {
		return nil, false
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return cp, true
}

// ReadString obtiene el contenido en string de un archivo en memoria.
func (v *VFS) ReadString(filePath string) (string, bool) {
	data, exists := v.ReadFile(filePath)
	if !exists {
		return "", false
	}
	return string(data), true
}

// ReadJSON deserializa un archivo JSON en memoria hacia una estructura destino.
func (v *VFS) ReadJSON(filePath string, target any) error {
	data, exists := v.ReadFile(filePath)
	if !exists {
		return path.ErrBadPattern
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(target)
}

// WriteJSON serializa una estructura como JSON formateado y la guarda en el VFS.
func (v *VFS) WriteJSON(filePath string, data any, indent string) error {
	if indent == "" {
		indent = "  "
	}
	encoded, err := json.MarshalIndent(data, "", indent)
	if err != nil {
		return err
	}
	// Agregar salto de línea al final por convención POSIX
	encoded = append(encoded, '\n')
	v.WriteFile(filePath, encoded)
	return nil
}

// DeleteFile elimina un archivo del VFS.
func (v *VFS) DeleteFile(filePath string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	norm := normalizePath(filePath)
	if _, exists := v.files[norm]; exists {
		delete(v.files, norm)
		return true
	}
	return false
}

// Exists comprueba si un archivo existe en el VFS.
func (v *VFS) Exists(filePath string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	_, exists := v.files[normalizePath(filePath)]
	return exists
}

// ListFiles devuelve una lista ordenada alfabéticamente de todas las rutas de archivos en el VFS.
func (v *VFS) ListFiles() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	res := make([]string, 0, len(v.files))
	for k := range v.files {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

// FileCount devuelve la cantidad total de archivos en memoria.
func (v *VFS) FileCount() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.files)
}

// Clear vacía todos los archivos del VFS.
func (v *VFS) Clear() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.files = make(map[string][]byte)
}
