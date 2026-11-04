---
name: postgresql
description: Production PostgreSQL engineering, schema design, B-tree and GIN indexes, connection pooling (PgBouncer), zero-downtime migrations, and EXPLAIN ANALYZE.
globs: ["**/*.sql", "src/db/**/*.ts", "docker-compose*.yml"]
---

# PostgreSQL Production Engineering & Indexing Standards

Guidelines for architecting scalable, resilient relational storage with PostgreSQL, focusing on query optimization, connection pooling, and zero-downtime schema evolution.

## Data Types & Index Strategy Reference

| Data Requirement | Recommended Type | Anti-Pattern Type | Reason / Advantage |
| :--- | :--- | :--- | :--- |
| **Timestamps** | `TIMESTAMPTZ` | `TIMESTAMP` | Preserves UTC point in time across timezones |
| **Primary Keys** | `UUIDv7` or `BIGINT` | `VARCHAR(36)` / `UUIDv4` | UUIDv7 preserves sequential index locality; avoids B-tree fragmentation |
| **Semi-Structured** | `JSONB` | `TEXT` / `JSON` | Binary format supports GIN indexing, containment operators (`@>`) |
| **Monetary Values** | `NUMERIC(12, 2)` or integer cents | `FLOAT` / `REAL` | Eliminates floating-point rounding errors |

## Index Types & When to Use

- **B-Tree (Default)**: Use for equality (`=`), range queries (`<`, `<=`, `>`, `>=`), and `ORDER BY`.
- **GIN (Generalized Inverted Index)**: Required for `JSONB`, array contains (`@>`), and full-text search vectors (`tsvector`).
- **Partial Indexes**: Index only a subset of rows meeting a condition (e.g. `CREATE INDEX idx_active_orders ON orders(user_id) WHERE status = 'pending';`).
- **Covering Indexes**: Append non-key columns using `INCLUDE (col1, col2)` to satisfy index-only scans without table page lookups.

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Unindexed Foreign Keys & Sequential Table Scans
Omitting indexes on foreign keys causes slow cascading deletes, table locks, and sequential scans during join queries.

```sql
-- BAD: Missing index on foreign key column
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id), -- Slow joins & cascading deletes!
  total_cents INTEGER NOT NULL
);
```

### ✅ Recommended: Explicit Indexing with Query Analysis

```sql
-- GOOD: Explicit foreign key index and query optimization
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  total_cents INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- B-Tree index on foreign key for O(log N) joins and deletes
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- Verify query plans in development
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM orders WHERE user_id = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
ORDER BY created_at DESC;
```

## Connection Pool Sizing (PgBouncer Formula)

Formula for sizing connection pools to PostgreSQL without CPU context switching penalties:

$$\text{Max Connections} = ((\text{CPU Cores} \times 2) + \text{Disk Spindle Count})$$

- For a 4-core server, pool size should be between 10-20 connections.
- Use **PgBouncer** in transaction pooling mode (`pool_mode = transaction`) for microservices and serverless architectures.
