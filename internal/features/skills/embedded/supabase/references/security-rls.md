# Row Level Security (RLS) Performance & Subquery Caching

## Critical Rule: Wrap `auth.uid()` in Scalar Subqueries

When writing RLS policies, never call `auth.uid()` directly in the `USING` or `WITH CHECK` expression.

### ❌ Anti-Pattern (Uncached O(N) evaluation)
```sql
CREATE POLICY "Users can select own orders"
ON orders FOR SELECT
USING (auth.uid() = user_id);
```

### ✅ Recommended Pattern (Cached scalar subquery)
```sql
CREATE POLICY "Users can select own orders"
ON orders FOR SELECT
USING ((SELECT auth.uid()) = user_id);
```

## JWT Metadata Claims Access

Avoid joining the `users` table to check tenant IDs or roles. Read directly from the signed JWT claims:

```sql
CREATE POLICY "Tenant isolation"
ON projects FOR ALL
USING (
  tenant_id = ((SELECT auth.jwt()) -> 'app_metadata' ->> 'tenant_id')::uuid
);
```
