---
name: clerk-nextjs-patterns
description: Advanced Next.js patterns - middleware, Server Actions, caching with
  Clerk.
license: MIT
allowed-tools: WebFetch
compatibility: Requires NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY and CLERK_SECRET_KEY. For manual JWT verification (standalone API servers without Clerk middleware), additionally requires CLERK_JWT_KEY or CLERK_PEM_PUBLIC_KEY.
metadata:
  author: clerk
  version: 2.2.0
---

# Next.js Patterns

> **Version**: Check `package.json` for the SDK version — see `clerk` skill for the version table. Core 2 differences are noted inline with `> **Core 2 ONLY (skip if current SDK):**` callouts.

For basic setup, see `clerk-setup` skill.

> **Route protection:** call `await auth.protect()` in each page, Route Handler, and Server Function. Keep `clerkMiddleware()` bare. `createRouteMatcher` is deprecated in `@clerk/nextjs` 7 and will be removed in the next major version. Don't add it. If a project already uses it, leave it working and point the user to the [migration guide](https://clerk.com/docs/guides/development/upgrading/upgrade-guides/migrate-from-create-route-matcher).

## What Do You Need?

| Task | Reference |
|------|-----------|
| Server vs client auth (`auth()` vs hooks) | references/server-vs-client.md |
| Protect routes (check in each page, handler, and Server Function) | references/middleware-strategies.md |
| Protect Server Actions | references/server-actions.md |
| API route auth (401 vs 403) | references/api-routes.md |
| Cache auth data (user-scoped caching) | references/caching-auth.md |

## References

| Reference | Description |
|-----------|-------------|
| `references/server-vs-client.md` | `await auth()` vs hooks |
| `references/middleware-strategies.md` | `auth.protect()` in each resource, bare `clerkMiddleware()` in `proxy.ts` (Next.js <=15: `middleware.ts`) |
| `references/server-actions.md` | Protect mutations |
| `references/api-routes.md` | 401 vs 403 |
| `references/caching-auth.md` | User-scoped caching |

## Mental Model

Server vs Client = different auth APIs:
- **Server**: `await auth()` from `@clerk/nextjs/server` (async!)
- **Client**: `useAuth()` hook from `@clerk/nextjs` (sync)

Never mix them. Server Components use server imports, Client Components use hooks.

Key properties from `auth()`:
- `isAuthenticated` — boolean, replaces the `!!userId` pattern
- `sessionStatus` — `'active'` | `'pending'`, for detecting incomplete session tasks
- `userId`, `orgId`, `orgSlug`, `has()`, `protect()` — unchanged

> **Core 2 ONLY (skip if current SDK):** `isAuthenticated` and `sessionStatus` are not available. Check `!!userId` instead.

## Minimal Pattern

```typescript
// Server Component
import { auth } from '@clerk/nextjs/server'

export default async function Page() {
  const { isAuthenticated, userId } = await auth()  // MUST await!
  if (!isAuthenticated) return <p>Not signed in</p>
  return <p>Hello {userId}</p>
}
```

> **Core 2 ONLY (skip if current SDK):** `isAuthenticated` is not available. Use `if (!userId)` instead.

### Conditional Rendering with `<Show>`

For client-side conditional rendering based on auth state. `<Show>` covers both authentication checks and authorization (feature, plan, role, permission) in one component.

**Authentication check:**

```tsx
import { Show } from '@clerk/nextjs'

<Show when="signed-in" fallback={<p>Please sign in</p>}>
  <Dashboard />
</Show>
```

**Authorization checks (B2B):**

```tsx
// Feature-based (preferred — features can move between plans without redeploy)
<Show when={{ f

