package doctor

// CheckStatus define el resultado del check de salud.
type CheckStatus string

const (
	StatusOk   CheckStatus = "OK"
	StatusWarn CheckStatus = "WARN"
	StatusFail CheckStatus = "FAIL"
	StatusInfo CheckStatus = "INFO"
)

// HealthItem representa un ítem individual de diagnóstico dentro de un módulo.
type HealthItem struct {
	Name        string      `json:"name"`        // Nombre descriptivo (ej: "Node.js runtime")
	Status      CheckStatus `json:"status"`      // OK, WARN, FAIL, INFO
	Message     string      `json:"message"`     // Mensaje descriptivo (ej: "v20.10.0 (>= 18.0.0)")
	Detail      string      `json:"detail"`      // Detalle adicional u observación
	FixCommand  string      `json:"fixCommand"`  // Comando sugerido para resolver (ej: "nvm install 20")
}

// HealthCheckReport consolida los 5 módulos de diagnóstico para un proyecto.
type HealthCheckReport struct {
	// [1/5] Entorno & Herramientas
	Environment []HealthItem `json:"environment"`

	// [2/5] Stack & Arquitectura (manejado junto a DiagnosticResult)
	Stack []HealthItem `json:"stack"`

	// [3/5] Configuración & Variables (.env)
	EnvVars []HealthItem `json:"envVars"`

	// [4/5] Base de Datos & Puertos Locales
	DatabaseAndPorts []HealthItem `json:"databaseAndPorts"`

	// [5/5] Estado de AI Skills y Contexto del Editor
	Skills []HealthItem `json:"skills"`

	// Acciones recomendadas consolidadas
	RecommendedActions []string `json:"recommendedActions"`
}
