package types

// ScaffoldConfig representa las opciones seleccionadas por el usuario
// para la generación del proyecto.
type ScaffoldConfig struct {
	ProjectName    string
	Recipe         string
	InitGit        bool
	Frontend       string
	Backend        string
	PackageManager string
	Database       string
	ORM            string
	Auth           string
	Addons         string
}
