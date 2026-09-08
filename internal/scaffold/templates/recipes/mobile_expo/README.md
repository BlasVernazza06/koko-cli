# [[.ProjectName]] — Mobile Expo + Express API Monorepo

Monorrepo multiplataforma moderno para aplicaciones móviles con **Expo (React Native)**, **Node.js Express (TypeScript)**, **PostgreSQL (Prisma ORM)** y **Turborepo**.

## 🧱 Estructura del Monorrepo

### Aplicaciones (`apps/`)
- **`mobile`**: Aplicación móvil multiplataforma desarrollada con **Expo** (React Native, TypeScript, Lucide icons).
- **`api`**: Servidor backend en **Express (TypeScript)** con conexión a **Prisma ORM** y validación **Zod**.

### Paquetes Compartidos (`packages/`)
- **`db`**: Esquema de datos **Prisma** con PostgreSQL.
- **`validators`**: Esquemas de validación **Zod** tipados compartidos.

## 🚀 Inicio Rápido

### 1. Iniciar Base de Datos con Docker
```bash
docker compose up -d
```

### 2. Instalar Dependencias y Generar Esquema
```bash
pnpm install
pnpm db:generate
pnpm db:push
```

### 3. Iniciar el Backend
```bash
pnpm --filter api dev
```
La API REST estará corriendo en `http://localhost:8080` (con endpoints `/api/health` y `/api/tasks`).

### 4. Iniciar la App Móvil con Expo
```bash
pnpm --filter mobile dev
```
Abre **Expo Go** en tu dispositivo físico escaneando el código QR en la terminal, o presiona `a` para Android Emulator / `i` para iOS Simulator.
