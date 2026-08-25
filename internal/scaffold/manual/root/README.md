# [[.ProjectName]]

Monorepo generated with [Koko-cli](https://github.com/BlasVernazza06/koko-cli) powered by [Turborepo](https://turbo.build/repo).

## Project Structure

```text
├── apps/
│   ├── web/        # Frontend application
│   └── api/        # Backend / API service
├── packages/
│   ├── db/                 # Database schema, client & migrations
│   ├── typescript-config/  # Shared TypeScript configuration
│   └── eslint-config/      # Shared ESLint configuration
├── package.json
└── turbo.json
```

## Getting Started

1. Install dependencies:
   ```bash
   [[if eq .PackageManager "bun"]]bun install[[else if eq .PackageManager "npm"]]npm install[[else]]pnpm install[[end]]
   ```

2. Run development servers:
   ```bash
   [[if eq .PackageManager "bun"]]bun dev[[else if eq .PackageManager "npm"]]npm run dev[[else]]pnpm dev[[end]]
   ```

3. Build all applications:
   ```bash
   [[if eq .PackageManager "bun"]]bun build[[else if eq .PackageManager "npm"]]npm run build[[else]]pnpm build[[end]]
   ```
