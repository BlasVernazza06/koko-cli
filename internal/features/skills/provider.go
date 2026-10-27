package skills

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const autoskillsRawBaseURL = "https://raw.githubusercontent.com/midudev/autoskills/main/packages/autoskills/skills-registry"

// Provider orquesta el suministro híbrido (Red Github -> Cache Local -> Fallback Embebido).
type Provider struct {
	CacheDir   string
	HTTPClient *http.Client
}

// NewProvider instancia el proveedor configurando el directorio de caché de usuario.
func NewProvider() (*Provider, error) {
	userCache, err := os.UserCacheDir()
	if err != nil {
		userCache = os.TempDir()
	}
	cacheDir := filepath.Join(userCache, "koko", "skills-cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		cacheDir = filepath.Join(os.TempDir(), "koko-skills-cache")
		_ = os.MkdirAll(cacheDir, 0755)
	}

	return &Provider{
		CacheDir: cacheDir,
		HTTPClient: &http.Client{
			Timeout: 4 * time.Second,
		},
	}, nil
}

// FetchSkill resuelve el contenido de una skill siguiendo la estrategia híbrida.
func (p *Provider) FetchSkill(ctx context.Context, id string) (*SkillDefinition, error) {
	normID := strings.ToLower(strings.TrimSpace(id))

	// 1. Intentar descargar desde GitHub raw de autoskills si hay conexión
	if netSkill, err := p.fetchFromNetwork(ctx, normID); err == nil && netSkill != nil {
		_ = p.saveToCache(normID, netSkill.Content)
		return netSkill, nil
	}

	// 2. Intentar leer desde la caché local del SO
	if cachedContent, err := p.loadFromCache(normID); err == nil && cachedContent != "" {
		return &SkillDefinition{
			ID:          normID,
			Name:        FormatSkillTitle(normID),
			Category:    categorizeSkill(normID),
			Description: extractDescriptionFromFrontmatter(cachedContent),
			Content:     cachedContent,
			Source:      "cache",
		}, nil
	}

	// 3. Fallback embebido (go:embed)
	if emb, ok := GetEmbeddedSkill(normID); ok {
		return emb, nil
	}

	return nil, fmt.Errorf("skill '%s' not found across network, cache or embedded fallback", normID)
}

func (p *Provider) fetchFromNetwork(ctx context.Context, id string) (*SkillDefinition, error) {
	url := fmt.Sprintf("%s/%s/SKILL.md", autoskillsRawBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // Límite seguro de 1MB
	if err != nil {
		return nil, err
	}

	content := string(data)
	return &SkillDefinition{
		ID:          id,
		Name:        FormatSkillTitle(id),
		Category:    categorizeSkill(id),
		Description: extractDescriptionFromFrontmatter(content),
		Content:     content,
		Source:      "network",
	}, nil
}

func (p *Provider) loadFromCache(id string) (string, error) {
	cachePath := filepath.Join(p.CacheDir, id+".md")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (p *Provider) saveToCache(id string, content string) error {
	cachePath := filepath.Join(p.CacheDir, id+".md")
	return os.WriteFile(cachePath, []byte(content), 0644)
}
