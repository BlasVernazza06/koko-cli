package processors

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

// generateRandomSecret genera una clave criptográfica de N bytes en Base64 URL-safe
func generateRandomSecret(byteLength int) string {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "koko_secret_fallback_32_chars_random_key"
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// ProcessEnvVariables genera automáticamente archivos .env y .env.example
// con las configuraciones y credenciales listas según el stack seleccionado.
func ProcessEnvVariables(v *vfs.VFS, cfg ProcessConfig) error {
	var envLines []string
	var exampleLines []string

	header := fmt.Sprintf("# ----------------------------------------------------\n# Environment variables for %s\n# ----------------------------------------------------\n", cfg.ProjectName)
	envLines = append(envLines, header)
	exampleLines = append(exampleLines, header)

	// 1. URLs de Frontend y Backend
	if cfg.Backend != "" && cfg.Backend != "none" {
		serverVars := []string{
			"PORT=4000",
			"NODE_ENV=development",
		}
		envLines = append(envLines, serverVars...)
		exampleLines = append(exampleLines, serverVars...)

		if cfg.Frontend != "" && cfg.Frontend != "none" {
			var apiUrl string
			if cfg.Frontend == "nextjs" {
				apiUrl = "NEXT_PUBLIC_API_URL=http://localhost:4000"
			} else {
				apiUrl = "VITE_API_URL=http://localhost:4000"
			}
			envLines = append(envLines, apiUrl)
			exampleLines = append(exampleLines, apiUrl)
		}
		envLines = append(envLines, "")
		exampleLines = append(exampleLines, "")
	}

	// 2. Base de datos
	db := strings.ToLower(cfg.Database)
	if db != "" && db != "none" {
		envLines = append(envLines, "# Database Connection")
		exampleLines = append(exampleLines, "# Database Connection")

		var dbUrl string
		switch db {
		case "postgres":
			dbUrl = fmt.Sprintf("DATABASE_URL=\"postgresql://postgres:password@localhost:5432/%s?schema=public\"", cfg.ProjectName)
		case "mysql":
			dbUrl = fmt.Sprintf("DATABASE_URL=\"mysql://root:password@localhost:3306/%s\"", cfg.ProjectName)
		case "mongodb":
			dbUrl = fmt.Sprintf("DATABASE_URL=\"mongodb://localhost:27017/%s\"", cfg.ProjectName)
		case "sqlite":
			dbUrl = "DATABASE_URL=\"file:./local.db\""
		}
		envLines = append(envLines, dbUrl, "")
		exampleLines = append(exampleLines, dbUrl, "")
	}

	// 3. Autenticación (Librerías Locales vs SaaS Dashboards)
	auth := strings.ToLower(cfg.Auth)
	if auth != "" && auth != "none" {
		envLines = append(envLines, "# Authentication")
		exampleLines = append(exampleLines, "# Authentication")

		if strings.Contains(auth, "better") {
			secret := generateRandomSecret(32)
			envLines = append(envLines, fmt.Sprintf("BETTER_AUTH_SECRET=\"%s\"", secret))
			envLines = append(envLines, "BETTER_AUTH_URL=\"http://localhost:4000\"")

			exampleLines = append(exampleLines, "BETTER_AUTH_SECRET=\"your-better-auth-secret\"")
			exampleLines = append(exampleLines, "BETTER_AUTH_URL=\"http://localhost:4000\"")

		} else if strings.Contains(auth, "nextauth") {
			secret := generateRandomSecret(32)
			envLines = append(envLines, fmt.Sprintf("NEXTAUTH_SECRET=\"%s\"", secret))
			envLines = append(envLines, "NEXTAUTH_URL=\"http://localhost:3000\"")

			exampleLines = append(exampleLines, "NEXTAUTH_SECRET=\"your-nextauth-secret\"")
			exampleLines = append(exampleLines, "NEXTAUTH_URL=\"http://localhost:3000\"")

		} else if strings.Contains(auth, "jwt") {
			secret := generateRandomSecret(32)
			envLines = append(envLines, fmt.Sprintf("JWT_SECRET=\"%s\"", secret))
			exampleLines = append(exampleLines, "JWT_SECRET=\"your-jwt-secret\"")

		} else if strings.Contains(auth, "clerk") {
			envLines = append(envLines, "NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=\"pk_test_...\"", "CLERK_SECRET_KEY=\"sk_test_...\"")
			exampleLines = append(exampleLines, "NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=\"pk_test_...\"", "CLERK_SECRET_KEY=\"sk_test_...\"")
		}

		envLines = append(envLines, "")
		exampleLines = append(exampleLines, "")
	}

	// 4. Addons de Pagos (Stripe, Polar) y Email (Resend, Brevo)
	payments := strings.ToLower(cfg.Payments)
	if strings.Contains(payments, "stripe") {
		envLines = append(envLines, "# Stripe Payments", "STRIPE_SECRET_KEY=\"sk_test_...\"", "NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=\"pk_test_...\"", "")
		exampleLines = append(exampleLines, "# Stripe Payments", "STRIPE_SECRET_KEY=\"sk_test_...\"", "NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=\"pk_test_...\"", "")
	}
	if strings.Contains(payments, "polar") {
		envLines = append(envLines, "# Polar Payments", "POLAR_ACCESS_TOKEN=\"polar_atfs_...\"", "POLAR_ORGANIZATION_ID=\"...\"", "")
		exampleLines = append(exampleLines, "# Polar Payments", "POLAR_ACCESS_TOKEN=\"polar_atfs_...\"", "POLAR_ORGANIZATION_ID=\"...\"", "")
	}

	email := strings.ToLower(cfg.Email)
	if strings.Contains(email, "resend") {
		envLines = append(envLines, "# Resend Email", "RESEND_API_KEY=\"re_...\"", "")
		exampleLines = append(exampleLines, "# Resend Email", "RESEND_API_KEY=\"re_...\"", "")
	}
	if strings.Contains(email, "brevo") {
		envLines = append(envLines, "# Brevo Email", "BREVO_API_KEY=\"xkeysib-...\"", "")
		exampleLines = append(exampleLines, "# Brevo Email", "BREVO_API_KEY=\"xkeysib-...\"", "")
	}

	// 5. Escritura e inyección inteligente en VFS
	ensureEnvContent(v, ".env", envLines)
	ensureEnvContent(v, ".env.example", exampleLines)

	return nil
}

func ensureEnvContent(v *vfs.VFS, fileName string, lines []string) {
	if !v.Exists(fileName) {
		v.WriteString(fileName, strings.Join(lines, "\n"))
		return
	}

	existing, ok := v.ReadString(fileName)
	if !ok || strings.TrimSpace(existing) == "" {
		v.WriteString(fileName, strings.Join(lines, "\n"))
		return
	}

	var toAdd []string
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) > 0 {
			key := parts[0]
			if key != "" && !strings.Contains(existing, key+"=") {
				toAdd = append(toAdd, line)
			}
		}
	}

	if len(toAdd) > 0 {
		if !strings.HasSuffix(existing, "\n") {
			existing += "\n"
		}
		newContent := existing + strings.Join(toAdd, "\n") + "\n"
		v.WriteString(fileName, newContent)
	}
}
