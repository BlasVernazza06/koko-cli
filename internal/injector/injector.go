package injector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold/processors"
	"github.com/BlasVernazza06/koko-cli/internal/vfs"
)

// AddAddon inyecta un nuevo servicio o herramienta en un proyecto existente.
func AddAddon(projectDir string, addonArg string) error {
	addon := strings.ToLower(strings.TrimSpace(addonArg))
	configPath := filepath.Join(projectDir, "koko.config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read koko.config.json: %w", err)
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var kokoCfg config.KokoConfig
	if err := json.Unmarshal(data, &kokoCfg); err != nil {
		return fmt.Errorf("invalid koko.config.json format: %w", err)
	}

	// 1. Verificar idempotencia / ya instalado
	switch addon {
	case "stripe":
		if kokoCfg.Features.Payments != nil && kokoCfg.Features.Payments.Provider == "stripe" {
			return fmt.Errorf("stripe payment provider is already configured in this project")
		}
		kokoCfg.Features.Payments = &config.PaymentInfo{Provider: "stripe"}

	case "polar":
		if kokoCfg.Features.Payments != nil && kokoCfg.Features.Payments.Provider == "polar" {
			return fmt.Errorf("polar payment provider is already configured in this project")
		}
		kokoCfg.Features.Payments = &config.PaymentInfo{Provider: "polar"}

	case "resend":
		if kokoCfg.Features.Email != nil && kokoCfg.Features.Email.Provider == "resend" {
			return fmt.Errorf("resend email provider is already configured in this project")
		}
		kokoCfg.Features.Email = &config.EmailInfo{Provider: "resend"}

	case "brevo":
		if kokoCfg.Features.Email != nil && kokoCfg.Features.Email.Provider == "brevo" {
			return fmt.Errorf("brevo email provider is already configured in this project")
		}
		kokoCfg.Features.Email = &config.EmailInfo{Provider: "brevo"}

	case "docker":
		if kokoCfg.Features.Infrastructure != nil && kokoCfg.Features.Infrastructure.DockerCompose {
			return fmt.Errorf("docker Compose is already configured in this project")
		}
		if kokoCfg.Features.Infrastructure == nil {
			kokoCfg.Features.Infrastructure = &config.Infrastructure{}
		}
		kokoCfg.Features.Infrastructure.DockerCompose = true

	case "github_actions":
		if kokoCfg.Features.Infrastructure != nil && kokoCfg.Features.Infrastructure.CICD == "github-actions" {
			return fmt.Errorf("gitHub Actions CI is already configured in this project")
		}
		if kokoCfg.Features.Infrastructure == nil {
			kokoCfg.Features.Infrastructure = &config.Infrastructure{}
		}
		kokoCfg.Features.Infrastructure.CICD = "github-actions"

	case "shadcn":
		if kokoCfg.Stack.Frontend != nil && kokoCfg.Stack.Frontend.UILibrary == "shadcn" {
			return fmt.Errorf("shadcn/ui is already configured in this project")
		}
		if kokoCfg.Stack.Frontend != nil {
			kokoCfg.Stack.Frontend.UILibrary = "shadcn"
		}

	case "lucide":
		if kokoCfg.Stack.Frontend != nil && kokoCfg.Stack.Frontend.Icons == "lucide" {
			return fmt.Errorf("lucide Icons is already configured in this project")
		}
		if kokoCfg.Stack.Frontend != nil {
			kokoCfg.Stack.Frontend.Icons = "lucide"
		}

	case "svgl", "motion", "zod":
		// Addons de UI / Validaciones
	default:
		return fmt.Errorf("unsupported addon: %s", addon)
	}

	// 2. Guardar manifiesto koko.config.json actualizado
	updatedData, err := json.MarshalIndent(kokoCfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode updated config: %w", err)
	}
	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to save koko.config.json: %w", err)
	}

	// 3. Cargar archivos del proyecto en VFS para inyección segura
	v := vfs.New()
	targetFiles := []string{
		"package.json",
		"apps/web/package.json",
		"apps/api/package.json",
		"packages/db/package.json",
		"packages/auth/package.json",
		"packages/ui/package.json",
		".env",
		".env.example",
	}

	for _, file := range targetFiles {
		fullPath := filepath.Join(projectDir, file)
		if content, err := os.ReadFile(fullPath); err == nil {
			v.WriteFile(file, bytes.TrimPrefix(content, []byte("\xef\xbb\xbf")))
		}
	}

	// 4. Mapear KokoConfig a ProcessConfig
	procCfg := processors.ProcessConfig{
		ProjectName:    kokoCfg.Project.Name,
		PackageManager: kokoCfg.Architecture.PackageManager,
		Addons:         addon,
	}

	if kokoCfg.Stack.Frontend != nil {
		procCfg.Frontend = kokoCfg.Stack.Frontend.Framework
	}
	if kokoCfg.Stack.Backend != nil {
		procCfg.Backend = kokoCfg.Stack.Backend.Framework
	}
	if kokoCfg.Stack.Database != nil {
		procCfg.Database = kokoCfg.Stack.Database.Provider
		procCfg.ORM = kokoCfg.Stack.Database.ORM
	}
	if kokoCfg.Features.Auth != nil {
		procCfg.Auth = kokoCfg.Features.Auth.Provider
	}
	if kokoCfg.Features.Payments != nil {
		procCfg.Payments = kokoCfg.Features.Payments.Provider
	}
	if kokoCfg.Features.Email != nil {
		procCfg.Email = kokoCfg.Features.Email.Provider
	}

	// 5. Inyectar dependencias y variables .env
	if err := processors.ProcessPackageJSONs(v, procCfg); err != nil {
		return fmt.Errorf("failed to update package.json files: %w", err)
	}

	if err := processors.ProcessEnvVariables(v, procCfg); err != nil {
		return fmt.Errorf("failed to update .env files: %w", err)
	}

	// 6. Escribir árbol VFS de vuelta al disco duro
	return scaffold.WriteTree(v, projectDir)
}
