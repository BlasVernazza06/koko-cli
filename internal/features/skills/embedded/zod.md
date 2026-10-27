---
name: zod
description: Zod schema validation best practices for type safety, parsing, and error handling. This skill should be used when defining z.object schemas, using z.string validations, safeParse, or z.infer. This skill does NOT cover React Hook Form integration patterns (use react-hook-form skill) or OpenAPI client generation (use orval skill).
---

# Zod Best Practices

Comprehensive schema validation guide for Zod in TypeScript applications. Contains 43 rules across 8 categories, prioritized by impact to guide automated refactoring and code generation.

## When to Apply

Reference these guidelines when:
- Writing new Zod schemas
- Choosing between parse() and safeParse()
- Implementing type inference with z.infer
- Handling validation errors for user feedback
- Composing complex object schemas
- Using refinements and transforms
- Optimizing bundle size and validation performance
- Reviewing Zod code for best practices

## Rule Categories by Priority

| Priority | Category | Impact | Prefix |
|----------|----------|--------|--------|
| 1 | Schema Definition | CRITICAL | `schema-` |
| 2 | Parsing & Validation | CRITICAL | `parse-` |
| 3 | Type Inference | HIGH | `type-` |
| 4 | Error Handling | HIGH | `error-` |
| 5 | Object Schemas | MEDIUM-HIGH | `object-` |
| 6 | Schema Composition | MEDIUM | `compose-` |
| 7 | Refinements & Transforms | MEDIUM | `refine-` |
| 8 | Performance & Bundle | LOW-MEDIUM | `perf-` |

## Quick Reference

### 1. Schema Definition (CRITICAL)

- `schema-use-primitives-correctly` - Use correct primitive schemas for each type
- `schema-use-unknown-not-any` - Use z.unknown() instead of z.any() for type safety
- `schema-avoid-optional-abuse` - Avoid overusing optional fields
- `schema-string-validations` - Apply string validations at schema definition
- `schema-use-enums` - Use enums for fixed string values
- `schema-coercion-for-form-data` - Use coercion for form and query data

### 2. Parsing & Validation (CRITICAL)

- `parse-use-safeparse` - Use safeParse() for user input
- `parse-async-for-async-refinements` - Use parseAsync for async refinements
- `parse-handle-all-issues` - Handle all validation issues not just first
- `parse-validate-early` - Validate at system boundaries

