---
name: express
description: Production-grade Express REST architecture, TypeScript typing, async error handling middleware, validation schemas, and security hardening.
globs: ["src/**/*.{ts,js}"]
---

# Express & TypeScript REST API Engineering

Production standards for scalable Express services, emphasizing modular router layering, asynchronous error propagation, input validation schemas, and security headers.

## Architecture Layers Reference

| Layer | Responsibility | Allowed Dependencies |
| :--- | :--- | :--- |
| **Routes** (`routes/`) | Endpoint URLs, HTTP methods, middleware attachments | Controller functions, auth/validation middlewares |
| **Middlewares** (`middlewares/`) | Authentication, schema validation, rate-limiting, CORS | Token decoders, Zod schemas, logger |
| **Controllers** (`controllers/`) | Request unpacking, response formatting, status codes | Business Service layer |
| **Services** (`services/`) | Domain logic, data mutations, external API calls | Database models, repositories, third-party clients |

## Security & Reliability Checklist

- **Helmet**: Always apply `helmet()` to secure HTTP response headers against clickjacking and MIME sniffing.
- **Async Handling**: Wrap asynchronous handler routes or adopt Express 5 to prevent unhandled rejected promises from crashing the process.
- **Centralized Error Handler**: Implement a 4-argument error middleware `(err, req, res, next)` at the very bottom of the middleware chain.
- **Rate Limiting**: Protect authentication and public mutation routes with `express-rate-limit`.

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Unhandled Async Errors & Missing Validation
Failing to catch errors in asynchronous Express handlers causes process unhandled rejections and leaks stack traces to clients.

```typescript
// BAD: Unhandled rejected promises, no schema validation, raw 500 error leak
app.post('/api/users', async (req, res) => {
  // If db.createUser throws, the request hangs or crashes the Node.js process!
  const user = await db.createUser(req.body);
  res.json(user);
});
```

### ✅ Recommended: Schema Validation Middleware & Structured Error Handler

```typescript
import express, { Request, Response, NextFunction } from 'express';
import helmet from 'helmet';
import { z } from 'zod';

const app = express();
app.use(helmet());
app.use(express.json());

// Input validation schema
const CreateUserSchema = z.object({
  email: z.string().email(),
  name: z.string().min(2),
});

// Reusable validation middleware
export const validate = (schema: z.ZodSchema) => 
  (req: Request, res: Response, next: NextFunction) => {
    const result = schema.safeParse(req.body);
    if (!result.success) {
      return res.status(400).json({ errors: result.error.flatten().fieldErrors });
    }
    req.body = result.data;
    next();
  };

// Async route wrapper (for Express 4)
export const asyncHandler = (fn: (req: Request, res: Response, next: NextFunction) => Promise<any>) =>
  (req: Request, res: Response, next: NextFunction) => {
    Promise.resolve(fn(req, res, next)).catch(next);
  };

// Controlled route
app.post('/api/users', validate(CreateUserSchema), asyncHandler(async (req: Request, res: Response) => {
  const user = await userService.createUser(req.body);
  return res.status(201).json(user);
}));

// Centralized error handler
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  console.error('[API Error]:', err.message);
  res.status(500).json({ error: 'Internal Server Error' });
});
```
