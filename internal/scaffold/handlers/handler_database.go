package handlers

import (
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/types"
)

// evaluateDatabase evalúa y mapea los esquemas, conexiones y configuración del ORM y base de datos
// hacia packages/db/ (en ecosistemas Node/TS) o apps/api/db/ (en Python/Go).
func EvaluateDatabase(rel string, config types.ScaffoldConfig) (string, bool) {
	if config.Database == "" || config.Database == "none" {
		return "", false
	}

	orm := strings.ToLower(config.ORM)
	db := strings.ToLower(config.Database)

	switch orm {
	case "drizzle":
		drizzleDbPrefix := "templates/db/drizzle/" + db + "/"
		if strings.HasPrefix(rel, drizzleDbPrefix) {
			dest := "packages/db/src/" + strings.TrimPrefix(rel, drizzleDbPrefix)
			return dest, true
		}
		if rel == "templates/db/drizzle/drizzle.config.ts" {
			return "packages/db/drizzle.config.ts", true
		}
		if rel == "templates/db/drizzle/package.json" {
			return "packages/db/package.json", true
		}

	case "prisma":
		prismaDbPrefix := "templates/db/prisma/" + db + "/"
		if strings.HasPrefix(rel, prismaDbPrefix) {
			dest := "packages/db/prisma/" + strings.TrimPrefix(rel, prismaDbPrefix)
			return dest, true
		}
		if rel == "templates/db/prisma/package.json" {
			return "packages/db/package.json", true
		}

	case "mongoose", "moongose":
		mongoosePrefix := "templates/db/mongoose/mongodb/"
		if strings.HasPrefix(rel, mongoosePrefix) {
			dest := "packages/db/" + strings.TrimPrefix(rel, mongoosePrefix)
			return dest, true
		}

	case "sqlalchemy":
		sqlAlchemyPrefix := "templates/db/sqlalchemy/" + db + "/"
		if strings.HasPrefix(rel, sqlAlchemyPrefix) {
			dest := "apps/api/app/db/" + strings.TrimPrefix(rel, sqlAlchemyPrefix)
			return dest, true
		}

	case "gorm":
		gormPrefix := "templates/db/gorm/" + db + "/"
		if strings.HasPrefix(rel, gormPrefix) {
			dest := "apps/api/db/" + strings.TrimPrefix(rel, gormPrefix)
			return dest, true
		}
	}

	return "", false
}
