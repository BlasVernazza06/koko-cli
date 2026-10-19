# Koko CLI Stack Matrix & Options Reference

Use this reference to choose valid combinations when scaffolding with `koko init`.

## 1. Quick Recipes (`--recipe / -r`)

| Recipe | Frontend | Backend | Database & ORM | Auth | Services |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **`saas`** | Next.js (App Router) | Self (Fullstack Next.js) | PostgreSQL + Drizzle | Better-Auth | Stripe + Resend + Docker + shadcn/ui |
| **`enterprise_nestjs`** | Next.js | NestJS | PostgreSQL + Prisma | Better-Auth | Stripe + Resend + Docker + shadcn/ui |
| **`java_spring`** | React + Vite | Spring Boot 3 (Java) | PostgreSQL + JPA/Hibernate | None | Docker Compose |
| **`fastapi_react`** | React + Vite | FastAPI (Python) | PostgreSQL + SQLAlchemy | None | Docker Compose |
| **`pern`** | React + Vite | Node Express | PostgreSQL + Prisma | JWT | None |
| **`mern`** | React + Vite | Node Express | MongoDB + Mongoose | JWT | None |
| **`mobile_expo`** | Expo (React Native) | Node Express | PostgreSQL + Prisma | JWT | Docker Compose |

---

## 2. Manual Matrix Options (`koko init [flags]`)

### Frontend Framework (`--frontend / -f`)
* `nextjs`: Next.js App Router (Fullstack SSR & Server Actions).
* `react`: React + Vite (Fast Client-side SPA).
* `nuxt`: Nuxt 3 (Vue fullstack framework).
* `svelte`: SvelteKit.
* `astro`: Astro (Content & islands architecture).
* `native`: React Native with Expo.
* `none`: REST API or backend only.

### Backend Framework (`--backend / -b`)
* `express`: Node.js Express with TypeScript.
* `hono`: Ultrafast multi-runtime web framework (Cloudflare, Bun, Node).
* `nestjs`: Modular enterprise TypeScript framework.
* `fastapi`: Python async framework with Pydantic v2.
* `go_chi`: High performance Go router.
* `spring_boot`: Java Spring Boot 3 enterprise backend.
* `self`: Uses frontend server capabilities (allowed ONLY if frontend is `nextjs`, `nuxt`, `svelte`, or `astro`).
* `none`: Client-only SPA or BaaS.

### API Layer (`--api`)
* `trpc`: End-to-end type-safe RPC (requires TypeScript and React frontend: `nextjs`, `react`, `native`).
* `orpc`: OpenAPI + RPC contract generator (supports Next.js, React, Nuxt, Svelte, Astro).
* `none`: Standard REST API endpoints without RPC abstraction.

### Package Manager (`--package-manager / -p`)
* `pnpm` (Recommended), `bun`, `npm`, `yarn`.
* For pure Go: `go_mod`. Pure Python: `pip` or `uv`. Pure Java: `mvn` or `gradle`.

### Database Provider (`--database`)
* `postgres`: PostgreSQL (includes Docker Compose).
* `mysql`: MySQL / MariaDB.
* `mongodb`: Document NoSQL DB.
* `sqlite`: Lightweight local file DB.
* `none`: No persistence layer.

### ORM / ODM (`--orm`)
* `drizzle`: Lightweight type-safe SQL ORM (PostgreSQL, MySQL, SQLite).
* `prisma`: Next-gen ORM with auto type generation (PostgreSQL, MySQL, SQLite, MongoDB).
* `mongoose`: ODM exclusively for MongoDB in Node/TS.
* `sqlalchemy`: Exclusively for Python / FastAPI.
* `gorm`: Exclusively for Go / Go Chi.
* `jpa`: Exclusively for Java / Spring Boot.
* `none`: Raw queries / direct driver.

### Authentication Provider (`--auth`)
* `better-auth`: Modern TypeScript framework for Next.js, Express, Hono, Nuxt, Svelte. Auto-generates 32-byte crypto secrets.
* `clerk`: Turnkey hosted user management for React / Next.js / Expo.
* `next-auth`: Exclusively for Next.js apps.
* `jwt`: Generic token-based authentication.
* `none`: No auth layer.

### Payments (`--payments`)
* `stripe`: Global payment gateway, checkout and subscriptions.
* `polar`: Merchant of Record & digital product monetization for developers.
* `none`: No payments.

### Email (`--email`)
* `resend`: Modern email API built for developers with React Email.
* `brevo`: Transactional email and SMS platform.
* `none`: No email delivery.

### Addons (`--addons`)
Comma-separated list:
* `shadcn`: Re-usable UI component package (`packages/ui`).
* `lucide`: Lucide icons pack.
* `svgl`: SVG logos and brand library.
* `motion`: Framer Motion declarative animation library.
* `zod`: TypeScript-first schema declaration and validation.
* `docker`: Docker Compose for database and local services.
* `github_actions`: Automated CI/CD linting and test workflows.

---

## 3. Strict Compatibility Invariants

1. **Both Frontend & Backend cannot be `none`**: The project must have at least one surface.
2. **Backend `self` requires Fullstack Frontend**: Only allowed with `nextjs`, `nuxt`, `svelte`, or `astro`.
3. **SPAs cannot connect directly to DB**: If Backend is `none` and Frontend is `react` or `native`, Database must be `none`.
4. **NoSQL vs SQL ORMs**:
   - `drizzle`, `sqlalchemy`, `gorm`, and `jpa` cannot be paired with `mongodb`.
   - `mongoose` is strictly for `mongodb`.
5. **tRPC Compatibility**: Requires TypeScript ecosystem and a React-based frontend (`nextjs`, `react`, `native`). For Nuxt or Svelte, use `orpc` or `none`.
6. **Clerk Compatibility**: Requires a React frontend (`nextjs`, `react`, `native`).
