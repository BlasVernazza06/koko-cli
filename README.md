<div align="center">

# 🦾 Claw CLI

### *Grab your stack, structure your project, and start building instantly.*

[![npm version](https://img.shields.io/npm/v/cli-repo-setup?color=33cd5f&style=flat-flat)](https://www.npmjs.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-blue.svg)](file:///c:/Users/USUARIO/Desktop/Programacion/Proyects/CLI%20repo%20ai%20setup/CONTRIBUTING.md)

</div>

---

**Claw CLI** (guiado por tu fiel asistente de andamiaje, **Clawy** 🦾) es un inicializador interactivo ultra veloz para la terminal. Olvídate de configurar a mano plantillas complejas de TypeScript, linters o contenedores Docker. Claw te permite seleccionar tu stack ideal y estructurará un proyecto profesional listo para producción en segundos.

```text
     ▄███████▄
    ██  • ◡ •  ██   🦾  { CLAWY }
    █████████████
    ██  ▀███▀  ██
     ▀█████████▀
     ▄█▀     ▀█▄
```

---

## ✨ Características Principales

*   🎨 **Experiencia Visual Premium**: Flujo interactivo elegante e intuitivo diseñado con `@clack/prompts`.
*   🚀 **Configuración en Segundos**: Generación física y estructuración automática de archivos y carpetas directamente en el disco.
*   📦 **Totalmente Personalizable**: Elige exactamente el frontend, backend y herramientas secundarias que necesitas.
*   🛠️ **Listo para Producción**: Configura linters, monorrepositorios y esquemas de docker listos para levantar bases de datos.

---

## 🛠️ Stack Tecnológico Soportado

Puedes combinar de forma flexible las siguientes tecnologías al inicializar tu proyecto:

| Frontend 💻 | Backend 🔌 | Herramientas 🐳 | Librerías Útiles 🎨 |
| :--- | :--- | :--- | :--- |
| **Next.js** (⚡ React) | **NestJS** (🦁 TypeScript) | **Docker Compose** | **Zod** (🛡️ Schema Validation) |
| **React** (⚛️ Vite SPA) | **ExpressJS** (🚂 JS/TS) | **Turborepo** (Monorrepo) | **Better Auth** (🔑 Auth) |
| **Angular** (🅰️ Google) | **FastAPI** (⚡ Python) | **ESLint & Prettier** | **Motion** (🎬 Animations) |
| **Sin Frontend** (Solo API) | **Actix-web** (🦀 Rust) | **TailwindCSS** | **SVGL & Lucide** (Vector Icons) |

---

## 🚀 Instalación y Uso Rápido

### Ejecutar Directamente (Sin Instalar)
La manera más rápida de comenzar un nuevo proyecto es mediante `npx`:
```bash
npx crear-proyecto
```

### Instalación Global
Si prefieres tenerlo disponible localmente como un comando global:
```bash
npm install -g cli-repo-setup
```
Una vez instalado, ejecuta el comando interactivo:
```bash
crear-proyecto
```

---

## 💻 Ejemplo del Flujo en Terminal

Al iniciar, **Clawy** te guiará paso a paso:

```bash
🦾 Clawy: "¡Hola! Soy tu asistente de andamiaje. Vamos a estructurar tu próximo proyecto."

? ¿Cómo se llamará el proyecto? mi-super-app
? ¿Qué Frontend deseas usar? › ⚡ Next.js
? ¿Qué tecnología para el Backend? › 🦁 NestJS
? Selecciona herramientas adicionales › [x] Docker Compose, [x] ESLint + Prettier

⠋ Generando la estructura real del proyecto en disco (esto puede tardar unos segundos)...
✔ ¡Estructura de archivos y configuración creada con éxito!
```

---

## 📁 Estructura del Repositorio Generado

Dependiendo de tu elección, Claw creará una arquitectura modular y organizada:

```text
mi-super-app/
├── apps/               # Si seleccionaste Turborepo (Monorrepo)
│   ├── web/            # Aplicación Frontend (ej. Next.js)
│   └── api/            # Servidor Backend (ej. NestJS)
├── docker-compose.yml  # Si seleccionaste Docker
├── .gitignore
├── package.json
└── README.md           # Con instrucciones específicas para correr el stack seleccionado
```

---

## 🤝 Contribuir

¡Nos encantan las contribuciones! Consulta nuestra [Guía de Contribución](file:///c:/Users/USUARIO/Desktop/Programacion/Proyects/CLI%20repo%20ai%20setup/CONTRIBUTING.md) para saber cómo clonar el repositorio, configurar tu entorno local y enviar Pull Requests.

---

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Consulta el archivo [LICENSE](file:///c:/Users/USUARIO/Desktop/Programacion/Proyects/CLI%20repo%20ai%20setup/LICENSE) para obtener más detalles.
