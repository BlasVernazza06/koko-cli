# Zod Object & Schema Validation Patterns

## Schema Definitions & Invariants

```typescript
import { z } from 'zod';

export const UserSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  age: z.coerce.number().int().min(18),
  role: z.enum(['admin', 'member', 'guest']).default('member'),
});

export type User = z.infer<typeof UserSchema>;
```

## Safe Parsing & Error Formatting

```typescript
const result = UserSchema.safeParse(input);
if (!result.success) {
  // Use flatten() for clean field-level error mapping
  const fieldErrors = result.error.flatten().fieldErrors;
  return { success: false, errors: fieldErrors };
}
return { success: true, data: result.data };
```
