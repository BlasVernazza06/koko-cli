package db

import (
	"context"
	"io"
)

// Action representa la operación solicitada sobre la base de datos.
type Action string

const (
	ActionMigrate  Action = "migrate"
	ActionPush     Action = "push"
	ActionGenerate Action = "generate"
	ActionStudio   Action = "studio"
	ActionSeed     Action = "seed"
	ActionReset    Action = "reset"
)

// SupportedActions contiene las acciones válidas en orden de visualización.
var SupportedActions = []Action{
	ActionMigrate,
	ActionPush,
	ActionGenerate,
	ActionStudio,
	ActionSeed,
	ActionReset,
}

// ProjectDBContext describe el stack de persistencia y runtime detectado en el proyecto.
type ProjectDBContext struct {
	ProjectDir     string // Ruta raíz del proyecto
	ConfigExists   bool   // Si se encontró koko.config.json
	Layout         string // "monorepo" o "standalone"
	PackageManager string // "pnpm", "npm", "bun", "yarn"
	Provider       string // "postgres", "mysql", "mongodb", "sqlite", etc.
	ORM            string // "drizzle", "prisma", "mongoose", "sqlalchemy", etc.
	RootScripts    map[string]string // scripts definidos en root package.json
	DBScripts      map[string]string // scripts definidos en packages/db/package.json
	HasDBPackage   bool   // Si existe packages/db
}

// ResolvedCommand representa el comando resuelto que debe ejecutarse.
type ResolvedCommand struct {
	Executable string   // Binario a ejecutar (ej: "pnpm", "npm", "npx", "bun")
	Args       []string // Argumentos completos sin concatenar
	WorkDir    string   // Directorio de trabajo donde debe ejecutarse
	SourceDesc string   // Descripción de la regla aplicada (para logging/dry-run)
}

// Options define las opciones de ejecución para el comando DB.
type Options struct {
	Dir    string    // Directorio objetivo (-d, --dir)
	Action Action    // Acción solicitada
	Args   []string  // Argumentos adicionales tras "--" o flags desconocidos
	DryRun bool      // Si es true, imprime el comando sin ejecutarlo
	Stdin  io.Reader // Inyección opcional para testing/streams
	Stdout io.Writer // Inyección opcional para testing/streams
	Stderr io.Writer // Inyección opcional para testing/streams
}

// Executor describe la interfaz encargada de correr procesos en el sistema operativo.
type Executor interface {
	Execute(ctx context.Context, cmd ResolvedCommand, stdin io.Reader, stdout, stderr io.Writer) error
}
