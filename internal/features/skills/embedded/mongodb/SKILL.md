---
name: mongodb
description: Production MongoDB document modeling, embedding vs referencing, compound indexes (ESR rule), aggregation pipelines, and transaction sessions.
globs: ["src/models/**/*.{ts,js}", "src/db/mongo*.ts"]
---

# Production MongoDB Architecture & Document Modeling

Guidelines for high-throughput, schema-disciplined MongoDB applications, prioritizing data access patterns, index efficiency, and bounded arrays.

## Data Modeling Strategy: Embed vs Reference

| Criterion | Choose Embedding (1-to-few) | Choose Referencing (1-to-many / 1-to-squillions) |
| :--- | :--- | :--- |
| **Cardinality** | Bounded items (e.g. up to 10 addresses) | Growing indefinitely (e.g. comments, logs) |
| **Access Pattern** | Data always retrieved together in one query | Data accessed and paginated independently |
| **Atomicity** | Requires atomic writes to parent and child | Updates happen across independent lifecycles |
| **Document Size Limit**| Well within the 16 MB BSON document limit | Risk of approaching the 16 MB limit over time |

## The Equality, Sort, Range (ESR) Index Rule

When creating compound indexes to support queries with filtering and sorting, order index keys following the ESR rule:

1. **E (Equality)**: Exact match fields first (`{ status: 1 }`).
2. **S (Sort)**: Sort order fields second (`{ createdAt: -1 }`).
3. **R (Range)**: Range condition fields last (`{ price: 1 }`).

Compound Index Example:
```javascript
// Matches query: { status: "active", price: { $gte: 100 } }, sort: { createdAt: -1 }
db.orders.createIndex({ status: 1, createdAt: -1, price: 1 });
```

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Unbounded Array Anti-Pattern
Pushing sub-documents directly into an unbounded array eventually hits the 16MB document limit and causes heavy document relocation on disk.

```typescript
// BAD: Unbounded sub-document growth inside user document
const UserSchema = new Schema({
  name: String,
  activityLogs: [
    { action: String, timestamp: Date } // Can grow to millions of rows, breaking 16MB limit!
  ]
});
```

### ✅ Recommended: Bucketing / Separate Collection with Foreign Key Index

```typescript
// GOOD: Separate collection with compound index on foreign key and timestamp
const ActivityLogSchema = new Schema({
  userId: { type: Schema.Types.ObjectId, ref: 'User', required: true, index: true },
  action: { type: String, required: true },
  timestamp: { type: Date, default: Date.now },
  metadata: { type: Map, of: Schema.Types.Mixed }
});

// Compound index supporting paginated user activity queries
ActivityLogSchema.index({ userId: 1, timestamp: -1 });

// Query with streaming or pagination
export async function getUserActivity(userId: string, page = 1, limit = 20) {
  return ActivityLog.find({ userId })
    .sort({ timestamp: -1 })
    .skip((page - 1) * limit)
    .limit(limit)
    .lean();
}
```

## Multi-Document Transactions

```typescript
import mongoose from 'mongoose';

export async function transferBalance(senderId: string, recipientId: string, amount: number) {
  const session = await mongoose.startSession();
  session.startTransaction();

  try {
    await Account.findByIdAndUpdate(senderId, { $inc: { balance: -amount } }, { session });
    await Account.findByIdAndUpdate(recipientId, { $inc: { balance: amount } }, { session });

    await session.commitTransaction();
  } catch (error) {
    await session.abortTransaction();
    throw error;
  } finally {
    session.endSession();
  }
}
```
