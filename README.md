<div align="center">

# 🦾 Koko CLI

### *Grab your stack, structure your project, and start building instantly.*

[![npm version](https://img.shields.io/npm/v/koko-cli?color=33cd5f&style=flat-flat)](https://www.npmjs.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

---

**Koko CLI** (guiado por tu fiel asistente de andamiaje 🦾) es un inicializador interactivo ultra veloz para la terminal escrito en Go. Olvídate de configurar a mano plantillas complejas de TypeScript, workspaces de monorrepos, bases de datos o contenedores Docker. Koko te permite seleccionar tu receta de producción ideal y estructurará tu proyecto listo para desarrollo en segundos de forma 100% offline.

```text
     .---.   .---.
    /  _  \_/  _  \       _  ______  _  ______
   |  (o)     (o)  |     | |/ / __ \| |/ / __ \
   |     (..)      |     | ' / /  | | ' / /  | |
    \   (____)    /      | . \ \__| | . \ \__| |
     '-----------'       |_|\_\____/|_|\_\____/  v1.0.0
```

---

## ✨ Características Principales

*   🎨 **Experiencia Visual Premium**: Interfaz TUI interactiva y animada construida con Bubble Tea y Lipgloss.
*   🚀 **Instalación Instantánea**: Desempaquetado físico y estructuración en disco en milisegundos.
*   📦 **Recetas de Producción**: Bootstrapping con configuraciones listas y optimizadas (ej. SaaS Starter con autenticación y base de datos).
*   🐳 **Monorrepo e Infraestructura**: Configuración automática de workspaces de pnpm/npm, linters, TypeScript y contenedores Docker Compose locales.

---

## 🛠️ Stack Tecnológico Soportado en Recetas (v1.0.0)

| Frontend 💻 | Backend 🔌 | Herramientas 🐳 | Base de Datos & ORM 🗄️ |
| :--- | :--- | :--- | :--- |
| **Next.js** (App Router, TS) | **Node.js Express** (TS) | **Docker Compose** | **PostgreSQL** (Drizzle / Prisma) |
| **React** (Vite SPA, TS) | **Go Fiber** (REST API) | **pnpm Workspaces** | **MySQL** (Drizzle / Prisma / SQLx) |
| **Vue.js** (Vite, TS) | **Hono** (Node.js Server) | **ESLint & Prettier** | **MongoDB** (Mongoose) |

---

## 🚀 Instalación y Uso Rápido

### Ejecutar Directamente (Sin Instalar)
La manera más rápida de comenzar un nuevo proyecto es mediante `npx`:
```bash
npx koko-cli init
```

### 🛠️ Desarrollo y Compilación Local (Go)
Si estás ejecutando o modificando el CLI desde el código fuente:
```bash
# Compilar el binario
go build -o koko.exe main.go

# Ejecutar el asistente interactivo
./koko.exe init
```

---

## 💻 Flujo en Terminal

Al iniciar, **Koko** te guiará con indicadores de progreso interactivos:

```bash
  ◇ Estructurando directorios y copiando plantillas...  ✓
  ◇ Generando configuración de Docker y DB...          ✓
  ◇ Inicializando repositorio Git...                   ✓
  ◇ Creando manifiesto koko.config.json...             ✓

✨ ¡Proyecto 'mi-app' creado con éxito en 2.45s!
```

---

## 📁 Estructura del Repositorio Generado

Dependiendo de tu elección, Koko estructurará una arquitectura monorrepo modular unificada:

```text
mi-app/
├── apps/
│   ├── frontend/         # Aplicación de Frontend (ej: Next.js)
│   └── backend/          # Servidor de API (ej: Express / Go Fiber)
├── packages/
│   └── db/               # Esquemas y cliente de conexión (Prisma / Drizzle)
├── docker-compose.yml    # Contenedor de PostgreSQL / MySQL local
├── koko.config.json      # Manifiesto de configuración de Koko
├── package.json          # Root package con workspaces configurados
└── README.md
```

---

## ⚙️ Manifiesto de Configuración (koko.config.json)

Al inicializar un proyecto, se genera un archivo `koko.config.json` en la raíz que actúa como el registro central del stack y permite futuras inyecciones de código en el "Día 2" (ej. mediante `koko add`):

```json
{
  "$schema": "https://koko-cli.dev/schema.json",
  "project": {
    "name": "mi-super-app",
    "cliVersion": "v1.0.0",
    "createdAt": "2026-08-20T20:50:00Z"
  },
  "architecture": {
    "layout": "monorepo",
    "packageManager": "pnpm"
  },
  "stack": {
    "frontend": {
      "framework": "next",
      "language": "typescript",
      "styling": "tailwindcss"
    },
    "backend": {
      "framework": "express",
      "language": "typescript"
    },
    "database": {
      "provider": "postgres",
      "orm": "prisma"
    }
  },
  "features": {
    "auth": {
      "provider": "better-auth",
      "status": "installed"
    },
    "infrastructure": {
      "dockerCompose": true,
      "ciCd": "github-actions"
    }
  }
}
```

---

## 🤝 Contribuir

¡Nos encantan las contribuciones! Consulta nuestra [Guía de Contribución](CONTRIBUTING.md) para saber cómo clonar el repositorio, configurar tu entorno local y enviar Pull Requests.

---

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Consulta el archivo [LICENSE](LICENSE) para obtener más detalles.
