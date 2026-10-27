# [[.ProjectName]] — SaaS Starter Monorepo

Este es un SaaS Starter profesional e integral inicializado con Koko CLI. Está estructurado como un monorrepo moderno administrado con **Turborepo** y **pnpm**.

## 🧱 Arquitectura del Proyecto

El monorrepo está compuesto por las siguientes aplicaciones y paquetes compartidos:

### Aplicaciones (`apps/`)
- **`web`**: Aplicación frontend principal construida con **Next.js** (App Router) y estilizada con **TailwindCSS v4**.

### Paquetes Compartidos (`packages/`)
- **`db`**: Persistencia de datos configurada con **Drizzle ORM** y soporte para bases de datos relacionales (PostgreSQL/Neon).
- **`auth`**: Módulo de autenticación preconfigurado con **Better Auth** que se conecta directamente al adaptador de base de datos.
- **`ui`**: Componentes de interfaz de usuario compartidos y listos para producción.
- **`validators`**: Esquemas de validación de datos basados en **Zod** compartidos entre el cliente y el servidor.
- **`typescript-config`**: Configuración compartida de TypeScript.
- **`eslint-config`**: Reglas compartidas de linter.

---

## 🚀 Inicio Rápido

### 1. Variables de Entorno
Crea un archivo `.env` en la raíz de la aplicación Next.js (`apps/web/.env`) y en el paquete de base de datos (`packages/db/.env` para las credenciales de base de datos):

```env
DATABASE_URL=tu_cadena_de_conexion_postgresql
BETTER_AUTH_SECRET=tu_secreto_para_better_auth
BETTER_AUTH_URL=http://localhost:3000
```

### 2. Instalar dependencias
Desde la raíz del monorrepo, ejecuta:
```bash
pnpm install
```

### 3. Migrar la Base de Datos
Para generar y aplicar el esquema inicial a tu base de datos:
```bash
# Dentro de packages/db
pnpm db:generate
pnpm db:push
```

### 4. Modo Desarrollo
Levanta todas las aplicaciones del monorrepo simultáneamente:
```bash
pnpm dev
```
