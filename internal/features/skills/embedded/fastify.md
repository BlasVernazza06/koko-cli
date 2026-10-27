---
name: fastify
description: High-throughput Node.js microservice architecture, JSON Schema serialization, encapsulation, lifecycle hooks, and type-safe plugins.
globs: ["src/**/*.{ts,js}"]
---

# Fastify High-Performance Microservice Engineering

Production guidelines for building ultrafast, schema-compiled APIs using Fastify, emphasizing lifecycle hooks, request-reply encapsulation, and fast serialization.

## Fastify Architecture & Lifecycle Hooks

Fastify executes requests through an explicit, deterministic lifecycle:

$$\text{onRequest} \rightarrow \text{preParsing} \rightarrow \text{preValidation} \rightarrow \text{preHandler} \rightarrow \text{handler} \rightarrow \text{preSerialization} \rightarrow \text{onResponse}$$

| Hook | Phase | Primary Purpose | Example |
| :--- | :--- | :--- | :--- |
| `onRequest` | Earliest entrypoint | Correlation ID attachment, rate-limiting | Set `req.headers['x-request-id']` |
| `preValidation` | Before validation | Query normalization, token extraction | Extract Bearer token from header |
| `preHandler` | After validation, before handler | Business authorization, tenant context check | Verify authenticated user permissions |
| `onError` | Error caught | Centralized structured error logging | Log stack trace with context |

## Performance Engine: JSON Schema Compilation

Declaring JSON schemas for request bodies and response payloads unlocks two massive performance advantages:
1. Fast, compiled validation powered by `Ajv`.
2. High-speed JSON serialization powered by `fast-json-stringify` (up to 2x faster than `JSON.stringify`).

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Unschema'd Routes with Loose Object Decoration
Skipping JSON schemas and mutating request objects directly bypasses Fastify's compiler optimizations and breaks TypeScript inference.

```typescript
// BAD: Missing response schema, ad-hoc request mutation
fastify.post('/users', async (request, reply) => {
  // Slower JSON.stringify fallback used automatically!
  const user = await db.users.create(request.body);
  return user;
});
```

### ✅ Recommended: Strictly Typed JSON Schema with `fastify-type-provider-zod`

```typescript
import Fastify from 'fastify';
import { serializerCompiler, validatorCompiler, ZodTypeProvider } from 'fastify-type-provider-zod';
import { z } from 'zod';

const server = Fastify({ logger: true }).withTypeProvider<ZodTypeProvider>();

server.setValidatorCompiler(validatorCompiler);
server.setSerializerCompiler(serializerCompiler);

// High-performance typed route with compiled schema
server.post(
  '/users',
  {
    schema: {
      body: z.object({
        email: z.string().email(),
        name: z.string().min(2),
      }),
      response: {
        201: z.object({
          id: z.string().uuid(),
          email: z.string().email(),
          createdAt: z.string(),
        }),
      },
    },
  },
  async (request, reply) => {
    // request.body is strictly typed via ZodTypeProvider
    const { email, name } = request.body;
    const user = await userService.createUser({ email, name });
    
    // Serialized ultra-fast using compiled fast-json-stringify
    return reply.status(201).send(user);
  }
);
```

## Plugin Encapsulation Guidelines

- Use `fastify-plugin` (`fp`) only for cross-cutting singleton plugins (e.g., database connectors, global auth decorators) that must break encapsulation and be accessible to all sibling scopes.
- Keep domain route modules encapsulated; routes registered in a plugin should not leak private decorators to unrelated parent or sibling routes.
