package doctor

// SectionType clasifica la sección a la que pertenece una dependencia.
type SectionType string

const (
	SectionFrontend       SectionType = "frontend"
	SectionFrontendStyle  SectionType = "frontend_styling"
	SectionFrontendUI     SectionType = "frontend_ui"
	SectionFrontendIcons  SectionType = "frontend_icons"
	SectionBackend        SectionType = "backend"
	SectionBackendDep     SectionType = "backend_dep"
	SectionDatabaseORM    SectionType = "db_orm"
	SectionDatabaseDriver SectionType = "db_driver"
	SectionAPI            SectionType = "api"
	SectionAuth           SectionType = "auth"
	SectionPayments       SectionType = "payments"
	SectionEmail          SectionType = "email"
	SectionLanguage       SectionType = "language"
)

// TechRule mapea una dependencia de package.json a una propiedad de KokoConfig.
type TechRule struct {
	Section SectionType
	Value   string
	Lang    string // Opcional, para asociar lenguaje (ej: "typescript", "python", "go", "java")
}

// TechCatalog mapea todos los paquetes de internal/catalog/versions.go a reglas de Koko.
var TechCatalog = map[string]TechRule{
	// 🌐 Frontend Frameworks
	"next":          {Section: SectionFrontend, Value: "next", Lang: "typescript"},
	"react":         {Section: SectionFrontend, Value: "react"},
	"vue":           {Section: SectionFrontend, Value: "vue"},
	"nuxt":          {Section: SectionFrontend, Value: "nuxt", Lang: "typescript"},
	"svelte":        {Section: SectionFrontend, Value: "svelte"},
	"@sveltejs/kit": {Section: SectionFrontend, Value: "svelte", Lang: "typescript"},
	"astro":         {Section: SectionFrontend, Value: "astro"},
	"expo":          {Section: SectionFrontend, Value: "expo", Lang: "typescript"},
	"react-native":  {Section: SectionFrontend, Value: "expo"},

	// 🎨 Styling
	"tailwindcss":         {Section: SectionFrontendStyle, Value: "tailwindcss"},
	"@nuxtjs/tailwindcss": {Section: SectionFrontendStyle, Value: "tailwindcss"},
	"@astrojs/tailwind":   {Section: SectionFrontendStyle, Value: "tailwindcss"},

	// 🧩 UI & Components
	"@radix-ui/react-slot":          {Section: SectionFrontendUI, Value: "shadcn"},
	"@radix-ui/react-dialog":        {Section: SectionFrontendUI, Value: "shadcn"},
	"@radix-ui/react-dropdown-menu":  {Section: SectionFrontendUI, Value: "shadcn"},

	// 🖼️ Icons
	"lucide-react":        {Section: SectionFrontendIcons, Value: "lucide"},
	"lucide-react-native": {Section: SectionFrontendIcons, Value: "lucide-react-native"},

	// ⚙️ Backend Frameworks
	"express":                  {Section: SectionBackend, Value: "express", Lang: "typescript"},
	"hono":                     {Section: SectionBackend, Value: "hono", Lang: "typescript"},
	"@hono/node-server":        {Section: SectionBackend, Value: "hono", Lang: "typescript"},
	"@nestjs/core":             {Section: SectionBackend, Value: "nestjs", Lang: "typescript"},
	"@nestjs/common":           {Section: SectionBackend, Value: "nestjs", Lang: "typescript"},
	"@nestjs/platform-express": {Section: SectionBackend, Value: "nestjs", Lang: "typescript"},
	"fastapi":                  {Section: SectionBackend, Value: "fastapi", Lang: "python"},
	"uvicorn":                  {Section: SectionBackend, Value: "fastapi", Lang: "python"},
	"go-chi":                   {Section: SectionBackend, Value: "go_chi", Lang: "go"},

	// ⚙️ Backend Dependencies
	"zod":              {Section: SectionBackendDep, Value: "zod"},
	"cors":             {Section: SectionBackendDep, Value: "cors"},
	"dotenv":           {Section: SectionBackendDep, Value: "dotenv"},
	"reflect-metadata": {Section: SectionBackendDep, Value: "reflect-metadata"},
	"rxjs":             {Section: SectionBackendDep, Value: "rxjs"},

	// 🔌 API Layer
	"@trpc/server":      {Section: SectionAPI, Value: "trpc"},
	"@trpc/client":      {Section: SectionAPI, Value: "trpc"},
	"@trpc/react-query": {Section: SectionAPI, Value: "trpc"},
	"@trpc/next":        {Section: SectionAPI, Value: "trpc"},
	"@orpc/server":      {Section: SectionAPI, Value: "orpc"},
	"@orpc/client":      {Section: SectionAPI, Value: "orpc"},
	"@orpc/react-query": {Section: SectionAPI, Value: "orpc"},
	"@orpc/openapi":     {Section: SectionAPI, Value: "orpc"},

	// 🗄️ Database ORMs
	"drizzle-orm":    {Section: SectionDatabaseORM, Value: "drizzle"},
	"drizzle-kit":    {Section: SectionDatabaseORM, Value: "drizzle"},
	"@prisma/client": {Section: SectionDatabaseORM, Value: "prisma"},
	"prisma":         {Section: SectionDatabaseORM, Value: "prisma"},
	"mongoose":       {Section: SectionDatabaseORM, Value: "mongoose"},

	// 🗄️ Database Drivers / Providers
	"postgres":       {Section: SectionDatabaseDriver, Value: "postgres"},
	"mysql2":         {Section: SectionDatabaseDriver, Value: "mysql"},
	"sqlite3":        {Section: SectionDatabaseDriver, Value: "sqlite"},
	"better-sqlite3": {Section: SectionDatabaseDriver, Value: "sqlite"},

	// 🔐 Authentication
	"better-auth":        {Section: SectionAuth, Value: "better-auth"},
	"@better-auth/cli":   {Section: SectionAuth, Value: "better-auth"},
	"@clerk/nextjs":      {Section: SectionAuth, Value: "clerk"},
	"@clerk/clerk-react": {Section: SectionAuth, Value: "clerk"},
	"@clerk/express":     {Section: SectionAuth, Value: "clerk"},
	"next-auth":          {Section: SectionAuth, Value: "next-auth"},

	// 💳 Payments
	"stripe":        {Section: SectionPayments, Value: "stripe"},
	"@polar-sh/sdk": {Section: SectionPayments, Value: "polar"},

	// 📧 Email
	"resend":          {Section: SectionEmail, Value: "resend"},
	"@getbrevo/brevo": {Section: SectionEmail, Value: "brevo"},

	// 🛠️ Language
	"typescript": {Section: SectionLanguage, Value: "typescript"},
}
