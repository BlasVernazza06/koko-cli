package catalog

// DependencyVersions es el catálogo maestro centralizado de versiones de dependencias
// utilizadas en todas las plantillas y generadores de Koko-cli.
var DependencyVersions = map[string]string{
	// ----------------------------------------------------
	// 🌐 Frontend Frameworks & UI
	// ----------------------------------------------------
	"next":                         "^15.1.6",
	"react":                        "^19.0.0",
	"react-dom":                    "^19.0.0",
	"vue":                          "^3.5.13",
	"vue-router":                   "^4.5.0",
	"nuxt":                         "^3.15.4",
	"@nuxtjs/tailwindcss":          "^6.13.1",
	"svelte":                       "^4.2.19",
	"@sveltejs/kit":                "^2.16.0",
	"@sveltejs/adapter-auto":       "^3.3.1",
	"@sveltejs/vite-plugin-svelte": "^3.1.2",
	"svelte-check":                 "^3.8.6",
	"vite":                         "^5.4.14",
	"@vitejs/plugin-react":         "^4.3.4",
	"tailwindcss":                  "^3.4.17",
	"postcss":                      "^8.4.49",
	"autoprefixer":                 "^10.4.20",
	"clsx":                         "^2.1.1",
	"tailwind-merge":               "^2.6.0",
	"class-variance-authority":     "^0.7.1",
	"lucide-react":                 "^0.475.0",
	"framer-motion":                "^11.18.2",
	"motion":                       "^11.18.2",
	"@ridemountainpig/svgl-react":  "^1.0.18",
	"@radix-ui/react-slot":          "^1.1.2",
	"@radix-ui/react-dialog":        "^1.1.6",
	"@radix-ui/react-dropdown-menu":  "^2.1.6",
	"@radix-ui/react-label":          "^2.1.2",
	"@radix-ui/react-select":         "^2.1.6",
	"@radix-ui/react-separator":      "^1.1.2",
	"@radix-ui/react-tooltip":        "^1.1.8",
	"@radix-ui/react-collapsible":    "^1.1.3",
	"sonner":                       "^1.7.4",
	"next-themes":                  "^0.4.4",
	"astro":                        "^4.16.18",
	"@astrojs/tailwind":            "^5.1.5",
	"expo":                         "^51.0.39",
	"react-native":                 "0.74.5",
	"lucide-react-native":          "^0.475.0",

	// ----------------------------------------------------
	// ⚙️ Backend Frameworks & Runtimes
	// ----------------------------------------------------
	"express":                  "^4.21.2",
	"hono":                     "^4.6.14",
	"@hono/node-server":        "^1.13.8",
	"@nestjs/core":             "^10.4.15",
	"@nestjs/common":           "^10.4.15",
	"@nestjs/platform-express": "^10.4.15",
	"@nestjs/cli":              "^10.4.9",
	"@nestjs/schematics":       "^10.2.3",
	"reflect-metadata":         "^0.2.2",
	"rxjs":                     "^7.8.1",
	"fastapi":                  ">=0.115.6",
	"uvicorn":                  ">=0.34.0",
	"pydantic":                 ">=2.10.4",
	"pydantic-settings":        ">=2.7.1",
	"python-dotenv":            ">=1.0.1",
	"go-chi":                   "v5.2.0",

	// ----------------------------------------------------
	// 🔌 API Layer (RPC & Data Fetching)
	// ----------------------------------------------------
	"@trpc/server":         "^10.45.2",
	"@trpc/client":         "^10.45.2",
	"@trpc/react-query":    "^10.45.2",
	"@trpc/next":           "^10.45.2",
	"@orpc/server":         "^0.80.0",
	"@orpc/client":         "^0.80.0",
	"@orpc/react-query":    "^0.80.0",
	"@orpc/openapi":        "^0.80.0",
	"@tanstack/react-query": "^5.62.11",


	// ----------------------------------------------------
	// 🗄️ Databases, ORMs & Drivers
	// ----------------------------------------------------
	"drizzle-orm":    "^0.38.3",
	"drizzle-kit":    "^0.30.1",
	"@prisma/client": "^5.22.0",
	"prisma":         "^5.22.0",
	"mongoose":       "^8.9.5",
	"postgres":       "^3.4.5",
	"mysql2":         "^3.12.0",
	"sqlite3":        "^5.1.7",
	"better-sqlite3": "^11.8.1",

	// ----------------------------------------------------
	// 🔐 Authentication, Validation & Utilities
	// ----------------------------------------------------
	"better-auth":         "^1.1.16",
	"@better-auth/cli":    "^1.1.16",
	"@clerk/nextjs":       "^6.9.9",
	"@clerk/clerk-react":  "^5.22.0",
	"@clerk/express":      "^1.3.36",
	"next-auth":           "^4.24.11",
	"zod":                 "^3.24.1",
	"cors":                "^2.8.5",
	"dotenv":              "^16.4.7",

	// ----------------------------------------------------
	// 💳 Payments & Billing
	// ----------------------------------------------------
	"stripe":        "^17.7.0",
	"@polar-sh/sdk": "^0.49.0",

	// ----------------------------------------------------
	// 📧 Email Services
	// ----------------------------------------------------
	"resend":          "^4.1.2",
	"@getbrevo/brevo": "^2.2.0",

	// ----------------------------------------------------
	// 🛠️ Tooling, TypeScript & Types
	// ----------------------------------------------------
	"typescript":         "^5.7.3",
	"tsx":                "^4.19.2",
	"ts-node":            "^10.9.2",
	"ts-node-dev":        "^2.0.0",
	"ts-loader":          "^9.5.2",
	"tsconfig-paths":     "^4.2.0",
	"source-map-support": "^0.5.21",
	"eslint":             "^9.18.0",
	"prettier":           "^3.4.2",
	"@types/node":        "^20.17.14",
	"@types/react":       "^19.0.8",
	"@types/react-dom":   "^19.0.3",
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
