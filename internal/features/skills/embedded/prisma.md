---
name: prisma-client-api
description: Prisma Client API reference covering model queries, filters, operators, and client methods. Use when writing database queries, using CRUD operations, filtering data, or configuring Prisma Client. Triggers on "prisma query", "findMany", "create", "update", "delete", "$transaction".
license: MIT
metadata:
  author: prisma
  version: "7.9.1"
---

# Prisma Client API Reference

Complete API reference for Prisma Client. This skill provides guidance on model queries, filtering, relations, and client methods for current Prisma projects.

## When to Apply

Reference this skill when:
- Writing database queries with Prisma Client
- Performing CRUD operations (create, read, update, delete)
- Filtering and sorting data
- Working with relations
- Using transactions
- Configuring client options

## Rule Categories by Priority

| Priority | Category | Impact | Prefix |
|----------|----------|--------|--------|
| 1 | Client Construction | HIGH | `constructor` |
| 2 | Model Queries | CRITICAL | `model-queries` |
| 3 | Query Shape | HIGH | `query-options` |
| 4 | Filtering | HIGH | `filters` |
| 5 | Relations | HIGH | `relations` |
| 6 | Transactions | CRITICAL | `transactions` |
| 7 | Raw SQL | CRITICAL | `raw-queries` |
| 8 | Client Methods | MEDIUM | `client-methods` |

## Quick Reference

- `constructor` - `PrismaClient` setup, adapter wiring, logging, and SQL commenter plugins
- `model-queries` - CRUD operations and bulk operations
- `query-options` - `select`, `include`, `omit`, sort, pagination
- `filters` - scalar and logical filter operators
- `relations` - relation reads and nested writes
- `transactions` - array and interactive transaction patterns
- `raw-queries` - `$queryRaw` and `$executeRaw` safety
- `client-methods` - lifecycle methods, extensions, and `satisfies` patterns for `prisma-client`

## Client Instantiation

```typescript
import { PrismaClient } from '../generated/client'
import { PrismaPg } from '@prisma/adapter-pg'

const adapter = new PrismaPg({
  connectionString: process.env.DATABASE_URL
})

const prisma = new PrismaClient({ adapter })
```

## Model Query Methods

| Method | Description |
|--------|-------------|
| `findUnique()` | Find one record 

