# [[.ProjectName]] — Enterprise NestJS + Next.js Monorepo

Arquitectura monorrepo escalable y tipada para aplicaciones empresariales de gran escala, construida con **Next.js** (App Router), **NestJS**, **PostgreSQL** (Prisma ORM), **Better-Auth / JWT**, **Stripe**, **Resend** y **Turborepo**.

## 🧱 Estructura del Monorrepo

### Aplicaciones (`apps/`)
- **`web`**: Frontend empresarial en **Next.js** (App Router, Tailwind CSS, Lucide icons, Stripe y Resend helpers).
- **`api`**: Backend modular en **NestJS** (Swagger en `/api/docs`, PrismaService, validación Zod/DTOs).

### Paquetes Compartidos (`packages/`)
- **`db`**: Esquema de base de datos relacional y cliente **Prisma** tipado.
- **`validators`**: Esquemas de validación **Zod** compartidos entre frontend y backend.
- **`typescript-config`**: Configuración de TypeScript unificada.
- **`eslint-config`**: Reglas de linting unificadas.

## 🚀 Inicio Rápido

### 1. Variables de Entorno
Copia y configura los archivos `.env`:
```bash
cp .env.example .env
```
Contenido básico para `.env`:
```env
DATABASE_URL="postgresql://postgres:password123@localhost:5432/[[.ProjectName]]"
PORT=4000
STRIPE_API_KEY="sk_test_..."
RESEND_API_KEY="re_..."
```

### 2. Levantar Servicios con Docker
```bash
docker compose up -d
```

### 3. Instalar y Migrar
```bash
pnpm install
pnpm db:generate
pnpm db:push
```

### 4. Modo Desarrollo
```bash
pnpm dev
```
- **Frontend (Next.js)**: `http://localhost:3000`
- **Backend (NestJS API)**: `http://localhost:4000`
- **Swagger Docs**: `http://localhost:4000/api/docs`
