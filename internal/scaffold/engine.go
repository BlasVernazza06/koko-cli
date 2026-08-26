package scaffold

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/BlasVernazza06/koko-cli/internal/errors"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold/processors"
	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

// GenerateVFS orquesta todo el proceso de generación del proyecto en memoria RAM (VFS).
// No toca el disco duro. Aplica templates, helpers, interpolación de variables
// y procesadores de package.json y .env.
func GenerateVFS(config ScaffoldConfig) (*vfs.VFS, error) {
	v := vfs.New()

	err := fs.WalkDir(templateFs, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return errors.NewTemplateError(path, "Fallo al acceder al sistema de plantillas embebido", walkErr, "Verifica la compilación del binario")
		}

		if d.IsDir() {
			return nil
		}

		destSubPath, shouldCopy := evaluatePath(path, config)
		if !shouldCopy {
			return nil
		}

		if strings.HasSuffix(destSubPath, ".tmpl") {
			destSubPath = strings.TrimSuffix(destSubPath, ".tmpl")
		}

		content, readErr := templateFs.ReadFile(path)
		if readErr != nil {
			return errors.NewTemplateError(path, "Fallo al leer plantilla embebida", readErr, "Verifica el archivo de plantilla")
		}

		// 1. Si es archivo binario, guardarlo directamente en VFS sin parsear
		if bytes.IndexByte(content, 0) != -1 || isBinaryExtension(destSubPath) {
			v.WriteFile(destSubPath, content)
			return nil
		}

		// 2. Si no contiene etiquetas de Go template, guardarlo tal cual
		if !shouldParseAsTemplate(content) {
			v.WriteFile(destSubPath, content)
			return nil
		}

		// 3. Renderizar template de Go con delimitadores personalizados [[ y ]]
		tmpl, parseErr := template.New(path).Delims("[[", "]]").Parse(string(content))
		if parseErr != nil {
			return errors.NewTemplateError(path, fmt.Sprintf("Error de sintaxis en plantilla %s", path), parseErr, "Verifica los bloques [[ ... ]]")
		}

		var buf bytes.Buffer
		if execErr := tmpl.Execute(&buf, config); execErr != nil {
			return errors.NewTemplateError(path, fmt.Sprintf("Fallo al interpolar variables en %s", path), execErr, "Verifica las variables del ScaffoldConfig")
		}

		v.WriteFile(destSubPath, buf.Bytes())
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Limpieza: si hay frontend o backend, remover el .gitkeep de apps/
	hasApps := (config.Frontend != "" && config.Frontend != "none") || (config.Backend != "" && config.Backend != "none")
	if hasApps {
		v.DeleteFile("apps/.gitkeep")
	}

	// Fase de Post-procesamiento en memoria (ajuste de package.json y .env)
	procCfg := processors.ProcessConfig{
		ProjectName:    config.ProjectName,
		PackageManager: config.PackageManager,
		Frontend:       config.Frontend,
		Backend:        config.Backend,
		Database:       config.Database,
		ORM:            config.ORM,
		Auth:           config.Auth,
		Addons:         config.Addons,
	}

	if err := processors.ProcessPackageJSONs(v, procCfg); err != nil {
		return nil, errors.NewPostProcessError("package.json", "Fallo al configurar scripts y dependencias", err)
	}

	if err := processors.ProcessEnvVariables(v, procCfg); err != nil {
		return nil, errors.NewPostProcessError(".env", "Fallo al generar archivo de entorno", err)
	}

	return v, nil
}
