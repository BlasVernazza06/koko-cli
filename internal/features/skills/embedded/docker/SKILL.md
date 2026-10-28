---
name: docker
description: Multi-stage Docker builds, Alpine/Distroless minimal images, layer caching, non-root users, security hardening, and Docker Compose orchestration.
globs: ["Dockerfile*", "docker-compose*.yml", ".dockerignore"]
---

# Production Docker Containerization & Compose Standards

Containerization guidelines for packaging high-performance Go, Node.js, and Python microservices with minimal attack surface and zero-overhead layer caching.

## Container Hardening & Image Sizing Reference

| Language / Runtime | Builder Base Image | Production Base Image | Approximate Final Size |
| :--- | :--- | :--- | :--- |
| **Go** | `golang:1.22-alpine` | `gcr.io/distroless/static-debian12:nonroot` | ~5-15 MB |
| **Node.js (Next/Hono)** | `node:22-alpine` | `node:22-alpine` (with non-root `node` user) | ~80-140 MB |
| **Python (FastAPI)** | `python:3.12-slim` | `python:3.12-slim` (with custom appuser) | ~90-160 MB |

## Core Principles Checklist

1. **Multi-Stage Builds**: Keep compiler toolchains, dev dependencies, and temporary source files out of the production runner image.
2. **Order of Layers**: Place rarely changing layers (system packages, lockfiles, dependency installations) before frequently changing source code.
3. **Never Run as Root**: Always define and switch to a non-privileged system user (`USER nonroot` or `USER node`).
4. **Graceful Shutdown**: Configure `STOPSIGNAL SIGTERM` and ensure PID 1 forwards signals to worker threads.
5. **Strict `.dockerignore`**: Exclude `node_modules`, `.git`, `.next`, `dist`, `coverage`, and `.env*` from build contexts.

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Monolithic Single-Stage Dockerfile as Root
Bundling compiler tools and running as root leads to massive image sizes (~1.5GB) and severe container breakout vulnerabilities.

```dockerfile
# BAD: Single-stage build, running as root with build dependencies included
FROM node:22
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

### ✅ Recommended: Multi-Stage Production Build with Non-Root User

```dockerfile
# Stage 1: Dependency resolution & build
FROM node:22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache libc6-compat
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

# Stage 2: Minimal hardened runner
FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Security: Non-root user execution
USER node

COPY --from=builder --chown=node:node /app/package.json ./
COPY --from=builder --chown=node:node /app/node_modules ./node_modules
COPY --from=builder --chown=node:node /app/dist ./dist

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

EXPOSE 3000
CMD ["node", "dist/index.js"]
```

## Go Multi-Stage Distroless Example

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot AS runner
WORKDIR /
COPY --from=builder /app/server /server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
```

## Production Docker Compose Orchestration

```yaml
version: "3.8"

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - DATABASE_URL=postgres://app:secret@db:5432/app_db
    depends_on:
      db:
        condition: service_healthy
    networks:
      - internal_net

  db:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: app_db
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app -d app_db"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - internal_net

volumes:
  pgdata:

networks:
  internal_net:
    driver: bridge
```
