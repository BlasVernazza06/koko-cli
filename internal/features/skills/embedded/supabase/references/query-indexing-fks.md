# Foreign Key Indexing & Cascade Deletes

## Critical Rule: Always Index Referencing Foreign Key Columns

PostgreSQL enforces foreign key constraints but does NOT automatically index them. Unindexed foreign keys cause full sequential table scans and heavy locking during parent `DELETE` or `UPDATE` cascades.

### ❌ Anti-Pattern
```sql
CREATE TABLE order_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
  item_id UUID NOT NULL
);
-- Missing index on order_id!
```

### ✅ Recommended Pattern
```sql
CREATE TABLE order_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  item_id UUID NOT NULL
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
```
