package catalog

// DependencyVersions es el catálogo maestro centralizado de versiones de dependencias
// utilizadas en todas las plantillas y generadores de Koko-cli.
var DependencyVersions = map[string]string{
	// ----------------------------------------------------
	// 🌐 Frontend Frameworks & UI
	// ----------------------------------------------------
	"next":                         "^14.2.13",
	"react":                        "^18.3.1",
	"react-dom":                    "^18.3.1",
	"vue":                          "^3.5.8",
	"vue-router":                   "^4.4.5",
	"nuxt":                         "^3.13.1",
	"@nuxtjs/tailwindcss":          "^6.12.1",
	"svelte":                       "^4.2.19",
	"@sveltejs/kit":                "^2.5.27",
	"@sveltejs/adapter-auto":       "^3.2.5",
	"@sveltejs/vite-plugin-svelte": "^3.1.2",
	"svelte-check":                 "^3.8.6",
	"vite":                         "^5.4.6",
	"@vitejs/plugin-react":         "^4.3.1",
	"tailwindcss":                  "^3.4.12",
	"postcss":                      "^8.4.47",
	"autoprefixer":                 "^10.4.20",
	"clsx":                         "^2.1.1",
	"tailwind-merge":               "^2.5.2",
	"lucide-react":                 "^0.446.0",

	// ----------------------------------------------------
	// ⚙️ Backend Frameworks & Runtimes
	// ----------------------------------------------------
	"express":                  "^4.19.2",
	"hono":                     "^4.6.0",
	"@hono/node-server":        "^1.13.0",
	"@nestjs/core":             "^10.4.1",
	"@nestjs/common":           "^10.4.1",
	"@nestjs/platform-express": "^10.4.1",
	"@nestjs/cli":              "^10.4.5",
	"@nestjs/schematics":       "^10.1.4",
	"reflect-metadata":         "^0.2.2",
	"rxjs":                     "^7.8.1",
	"fastapi":                  ">=0.115.0",
	"uvicorn":                  ">=0.30.6",
	"pydantic":                 ">=2.9.1",
	"pydantic-settings":        ">=2.5.2",
	"python-dotenv":            ">=1.0.1",
	"go-chi":                   "v5.0.12",

	// ----------------------------------------------------
	// 🗄️ Databases, ORMs & Drivers
	// ----------------------------------------------------
	"drizzle-orm":    "^0.30.10",
	"drizzle-kit":    "^0.21.4",
	"@prisma/client": "^5.19.1",
	"prisma":         "^5.19.1",
	"mongoose":       "^8.6.1",
	"postgres":       "^3.4.4",
	"mysql2":         "^3.9.7",
	"sqlite3":        "^5.1.7",
	"better-sqlite3": "^11.3.0",

	// ----------------------------------------------------
	// 🔐 Authentication, Validation & Utilities
	// ----------------------------------------------------
	"better-auth":   "^1.1.0",
	"@clerk/nextjs": "^5.7.0",
	"next-auth":     "^4.24.8",
	"zod":           "^3.23.8",
	"cors":          "^2.8.5",
	"dotenv":        "^16.4.5",

	// ----------------------------------------------------
	// 🛠️ Tooling, TypeScript & Types
	// ----------------------------------------------------
	"typescript":         "^5.6.2",
	"tsx":                "^4.19.1",
	"ts-node":            "^10.9.2",
	"ts-node-dev":        "^2.0.0",
	"ts-loader":          "^9.5.1",
	"tsconfig-paths":     "^4.2.0",
	"source-map-support": "^0.5.21",
	"eslint":             "^8.57.1",
	"prettier":           "^3.3.3",
	"@types/node":        "^20.16.5",
	"@types/react":       "^18.3.8",
	"@types/react-dom":   "^18.3.0",
	"@types/express":     "^4.17.21",
	"@types/cors":        "^2.8.17",
}

// GetVersion busca una dependencia en el catálogo maestro y devuelve su versión.
// Si el paquete no se encuentra registrado, devuelve fallback (o "latest" por defecto).
func GetVersion(pkg string, fallback ...string) string {
	if ver, exists := DependencyVersions[pkg]; exists {
		return ver
	}
	if len(fallback) > 0 && fallback[0] != "" {
		return fallback[0]
	}
	return "latest"
}

// Has comprueba si un paquete está registrado en el catálogo maestro.
func Has(pkg string) bool {
	_, exists := DependencyVersions[pkg]
	return exists
}

// GetAllVersions devuelve una copia segura de todo el mapa de versiones.
func GetAllVersions() map[string]string {
	cp := make(map[string]string, len(DependencyVersions))
	for k, v := range DependencyVersions {
		cp[k] = v
	}
	return cp
}
