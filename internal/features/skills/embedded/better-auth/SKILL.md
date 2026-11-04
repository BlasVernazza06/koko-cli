---
id: better-auth
title: Better Auth Modern Framework
category: Web Frameworks
description: Modern comprehensive TypeScript authentication framework, sessions, multi-factor auth, social providers, and database adapters.
tags: [auth, better-auth, typescript, session, security]
globs: ["src/lib/auth*.ts", "src/auth*.ts"]
---

# Better Auth Modern Authentication

Better Auth is a TypeScript-first, framework-agnostic authentication framework built for modern fullstack architectures (Next.js, Nuxt, SvelteKit, Express, Hono).

## Key Features & Protocol Support

- **Database Adapters**: Native drivers for Drizzle, Prisma, Kysely, and MongoDB.
- **Session Strategies**: Database sessions or stateless JWTs with rotation.
- **Pluggable Architecture**: 2FA (TOTP), Organizations/Teams, Magic Links, Passkeys (WebAuthn).
- **TypeScript First**: Full client-server RPC schema inference (`authClient.useSession()`).

## Comparison Table: Auth Alternatives

| Feature | Better Auth | NextAuth / Auth.js | Clerk |
| :--- | :--- | :--- | :--- |
| **Data Ownership** | 100% Self-Hosted in your DB | 100% Self-Hosted in your DB | Hosted SaaS (Vendor Lock-in) |
| **Framework Agnostic** | Next, Nuxt, Astro, Svelte, Hono | Primarily Next.js | Next, React, Expo |
| **Database Sync** | Auto-generates schemas (Drizzle/Prisma)| Manual schema migrations | Hosted Cloud DB |
| **Pricing** | Free & Open Source | Free & Open Source | Tiered active user fees |

## Server Configuration Example (`src/lib/auth.ts`)

```typescript
import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "@/db"; // Drizzle instance
import * as schema from "@/db/schema";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    schema: schema,
  }),
  emailAndPassword: {
    enabled: true,
    autoSignIn: true,
  },
  socialProviders: {
    github: {
      clientId: process.env.GITHUB_CLIENT_ID!,
      clientSecret: process.env.GITHUB_CLIENT_SECRET!,
    },
    google: {
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    },
  },
  session: {
    expiresIn: 60 * 60 * 24 * 7, // 7 days
    updateAge: 60 * 60 * 24, // 1 day
  },
});
```

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Unverified Client-Side Session Checks
Checking sessions exclusively on the client exposes protected mutations to forged requests.

```typescript
// BAD: Client-side check only; API handler executes unconditionally
export async function POST(req: Request) {
  // Missing server-side session verification!
  const data = await req.json();
  await deleteCustomerAccount(data.id);
  return Response.json({ success: true });
}
```

### ✅ Recommended: Strict Server Context Verification

```typescript
import { auth } from "@/lib/auth";
import { headers } from "next/headers";

export async function POST(req: Request) {
  const session = await auth.api.getSession({
    headers: await headers(),
  });

  if (!session || !session.user) {
    return new Response(JSON.stringify({ error: "Unauthorized" }), { status: 401 });
  }

  const data = await req.json();
  await deleteCustomerAccount(session.user.id, data.id);
  return Response.json({ success: true });
}
```
