---
name: clerk-react-patterns
description: 'React SPA auth patterns with @clerk/react for Vite/CRA - ClerkProvider
  setup, useAuth/useUser/useClerk hooks, React Router protected routes, custom sign-in
  flows. Triggers on: Vite Clerk setup, React Router auth, useAuth hook, protected
  route, custom sign-in form React.'
license: MIT
allowed-tools: WebFetch
metadata:
  author: clerk
  version: 1.0.0
---

# React SPA Patterns

> This skill covers `@clerk/react` for Vite/CRA SPAs. For Next.js use `clerk-nextjs-patterns`. For TanStack Start use `clerk-tanstack-patterns`.

## What Do You Need?

| Task | Reference |
|------|-----------|
| useAuth / useUser / useClerk hooks | references/hooks.md |
| Protected routes with React Router | references/protected-routes.md |
| Custom sign-in / sign-up forms | references/custom-flows.md |
| React Router v6/v7 integration | references/router-integration.md |

## References

| Reference | Description |
|-----------|-------------|
| `references/hooks.md` | useAuth, isLoaded guard |
| `references/protected-routes.md` | ProtectedRoute pattern |
| `references/custom-flows.md` | useSignIn, useSignUp flows |
| `references/router-integration.md` | React Router v6/v7 setup |

## Setup

```
npm install @clerk/react
```

`.env`:
```
VITE_CLERK_PUBLISHABLE_KEY=pk_...
```

`src/main.tsx`:
```tsx
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { ClerkProvider } from '@clerk/react'
import App from './App.tsx'

const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ClerkProvider publishableKey={PUBLISHABLE_KEY}>
      <App />
    </ClerkProvider>
  </StrictMode>,
)

