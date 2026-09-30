# [[.ProjectName]] — Enterprise Java Spring Boot + React SPA

Monorrepo empresarial de alto rendimiento con **Java Spring Boot 3**, **React 19 SPA**, **PostgreSQL** y **Docker**.

## 🧱 Estructura del Monorrepo

- **`apps/frontend`**: React 19 SPA con Vite, TypeScript y Tailwind CSS.
- **`apps/backend`**: API REST en Java Spring Boot 3 con Spring Data JPA y validación estricta.

## 🚀 Inicio Rápido

### 1. Levantar Base de Datos
```bash
docker compose up -d
```

### 2. Ejecutar Backend (Spring Boot)
```bash
cd apps/backend
./mvnw spring-boot:run # O mvn spring-boot:run
```
La API estará disponible en `http://localhost:8080` (con endpoints `/api/health` y `/api/items`).

### 3. Ejecutar Frontend (React)
```bash
pnpm install
pnpm dev
```
La aplicación web estará disponible en `http://localhost:5173`.
