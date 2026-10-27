# Matriz de Compatibilidad y Restricciones Técnicas (Core Stack)

Documento depurado de restricciones esenciales entre tecnologías base (Frontend, Backend, Runtime, Database, ORM, Auth, Payments y API).

---

## 1. Frontend

### Reglas Generales de Selección
- **Exclusividad Web:** Solo se permite seleccionar como máximo **1 framework Web** a la vez (`tanstack-router`, `tanstack-start`, `react-router`, `next`, `nuxt`, `svelte`, `astro`).
- **Exclusividad Nativa:** Solo se permite seleccionar como máximo **1 framework Nativo** a la vez (`native-bare`, `native-uniwind`, `native-unistyles`).

### Restricciones Específicas por Framework

| Frontend | Incompatibilidades | Motivo Técnico |
| :--- | :--- | :--- |
| **`astro`** | • API `trpc`<br>• Auth `clerk` | Falta de soporte/bindings oficiales para tRPC y Clerk en este stack (enfoque MPA). |
| **`nuxt`** | • API `trpc`<br>• Auth `clerk` | tRPC y Clerk requieren el ecosistema React en esta integración. |
| **`svelte`** | • API `trpc`<br>• Auth `clerk` | tRPC y Clerk requieren el ecosistema React en esta integración. |
| **`tanstack-router`** | • Backend `self` | Es una SPA estática y no produce rutas de servidor para el modo monolítico `self`. |
| **`react-router`** | • Backend `self` | No está soportado dentro de `FULLSTACK_FRONTENDS` para modo `self`. |
| **`next`** | Ninguna | Soporte total con Auth (Clerk, Better-Auth), APIs (tRPC, oRPC) y Backend `self`. |
| **`tanstack-start`** | Ninguna | Soporte total con Auth (Clerk, Better-Auth), APIs (tRPC, oRPC) y Backend `self`. |

---

## 2. Backend

### Modos y Servidores

#### A. `self` (Fullstack Monolítico)
- **Runtime forzado:** `runtime: none`.
- **Frontends soportados (`FULLSTACK_FRONTENDS`):** `next`, `tanstack-start`, `nuxt`, `svelte`, `astro`.
- **Incompatibilidades:**
  - Incompatible con frontends puramente clientes (`tanstack-router`, `react-router`).
  - Con `--auth clerk`: Solo compatible si el frontend es `next` o `tanstack-start`.

#### B. `hono`, `express`, `fastify` (Servidores Dedicados)
- Requieren `--runtime bun` o `--runtime node`.
- Compatibles con todos los frontends cliente a través de una API (`trpc` u `orpc`).
- Compatibles con `auth: clerk` y `auth: better-auth`.

#### C. `none` (Sin Backend / Solo Frontend)
- **Opciones forzadas a `none`:** `runtime: none`, `database: none`, `orm: none`, `api: none`, `auth: none`, `payments: none`.

---

## 3. Runtime

| Runtime | Backends Permitidos |
| :--- | :--- |
| **`bun`** | `hono`, `express`, `fastify` |
| **`node`** | `hono`, `express`, `fastify` |
| **`none`** | `self` o `none` |

---

## 4. Database & ORM

### Dependencia Cruzada
- La selección de una base de datos requiere obligatoriamente un ORM (y viceversa).
- Si `database === 'none'`, entonces `orm` debe ser `none`.

### Matriz ORM vs Database

| ORM | Bases de Datos Soportadas | Bases de Datos Incompatibles |
| :--- | :--- | :--- |
| **`drizzle`** | `sqlite`, `postgres`, `mysql` | `mongodb` |
| **`prisma`** | `sqlite`, `postgres`, `mysql`, `mongodb` | Ninguna |
| **`mongoose`** | `mongodb` | `sqlite`, `postgres`, `mysql` |
| **`none`** | `none` | Cualquier base de datos |

---

## 5. Autenticación (`auth`)

### `clerk`
- **Frontends compatibles:** Solo ecosistema React Web (`next`, `tanstack-start`, `react-router`, `tanstack-router`) o React Native (`native-bare`, `native-uniwind`, `native-unistyles`).
- **Frontends incompatibles:** `nuxt`, `svelte`, `astro`.
- **Backends compatibles:** Servidores dedicados (`hono`, `express`, `fastify`), o backend `self` (exclusivamente con `next` o `tanstack-start`).

### `better-auth`
- Compatible con todos los frontends (`next`, `tanstack-start`, `react-router`, `tanstack-router`, `nuxt`, `svelte`, `astro`, `native-*`).
- Compatible con todos los backends (`self`, `hono`, `express`, `fastify`).

### `none`
- Inhabilita la configuración de pagos (`payments`).

---

## 6. Pagos (`payments`)

### `polar`
- Requiere obligatoriamente `--auth better-auth`.
- Incompatible con `auth: clerk` o `auth: none`.

---

## 7. API (`api`)

### `trpc`
- **Frontends incompatibles:** `nuxt`, `svelte`, `astro`. (Deben usar `orpc` o `none`).
- **Frontends compatibles:** `next`, `tanstack-start`, `react-router`, `tanstack-router`, `native-*`.

### `orpc`
- Compatible con todos los frontends.

### `none`
- Sin capa de API RPC tipada entre cliente y servidor.
