---
name: hono
description: Use when building Hono web applications or when the user asks about Hono APIs, routing, middleware, JSX, validation, testing, or streaming. TRIGGER when code imports from 'hono' or 'hono/*', or user mentions Hono. Use Hono CLI to inspect and test the app.
---

# Hono Skill

Build Hono web applications. This skill provides inline API knowledge for AI. Use Hono CLI to inspect and test the app.

## Latest Documentation

For details beyond this inline reference, fetch the latest documentation from https://hono.dev. Get the index of doc pages from `https://hono.dev/llms.txt`, then fetch a page with the `Accept: text/markdown` header to receive it as Markdown:

```bash
curl -H "Accept: text/markdown" https://hono.dev/docs/helpers/cookie
```

## Hono CLI

Verify with [Hono CLI](https://github.com/honojs/cli), not with throwaway scripts or a dev server. It is a dev dependency (projects from `create-hono` already have it), loads the app in-process with `app.request()`, and prints JSON. Install it if missing, then read its help:

<!-- TODO at the 1.0 release: change @hono/cli@next to @hono/cli -->

```bash
npm install -D @hono/cli@next
npx hono --help
```

Follow the note at the top. It explains the JSON output contract and the workflow. `npx hono <command> --help` has examples and notes for each command (`routes`, `request`, `batch`, `snapshot`, `benchmark`, `ssg`, `init`).

<!-- TODO at the Hono v5 GA: change `npm create hono@next` to `npm create hono@latest` here and in hono-jsx -->

Notes:

- `npx hono routes` lists the routes; `npx hono request /path` sends one request — no server needed. If the app is not at the default path (`src/index.ts`), pass the file to every command: `npx hono request / src/app.ts`. Do not pass credentials directly in CLI arguments; use environment variables for sensitive values.
- Cloudflare bindings (KV, D1, R2, vars) reach `c.env` automatically — no server, no flag. In a project with `cloudflare.config.ts` (the `cf` CLI) and a Vite config, `request`, `batch`, and `snapshot` run through the Vite dev server by default, and `c.env` has the local bindings; `request` and `batch` default to Vite as well in a Vite project with no `src/index.*` (an app built by a Vite plugin; `snapshot` still needs the file), and `--runtime vite` selects it explicitly. In a project with a wrangler config, the app runs on Node.js with the real local bindings loaded; for the full Workers runtime there, add `--runtime workerd` (pass no file argument; the entry is `main` in the wrangler config). `--trace`, `--watch`, and a file argument run the app on Node.js, where `c.env` is empty in a `cf` project.
- `routes` and `snapshot` also load apps that import `cloudflare:workers` (Durable Objects); the import is stubbed while the routes are read.
- To add a Hono app to an existing directory, run `npx hono init --template cloudflare-workers` (or `nodejs`, `bun`, `deno`, ...; without `--template` it is detected from the files in the directory), then install the dependencies. It copies the `create-hono` template and never overwrites an existing file. For an empty directory, `npm create hono@next` does the same with prompts.
- For several requests, or a flow that keeps state (POST, then use the returned id), run them in one `hono batch -` call. One JSON object per line; `save` a value and use it as `{{id}}` in later steps. The steps share one app instance. Declare the expected status/body per step with `expect` (body is a deep partial match) and iterate until the summary shows `"failed": 0`:

  ```bash
  npx hono batch - <<'EOF'
  {"method":"POST","path":"/users","body":{"name":"Alice"},"save":{"id":".id"},"expect":{"status":201}}
  {"path":"/users/{{id}}","expect":{"status":200,"body":{"name":"Alice"}}}
  EOF
  ```

- Before changing existing routes, capture the current behavior: `npx hono snapshot` prints it as batch JSONL lines (real responses become the `expect`). Keep the lines, make the change, then rerun them with `npx hono batch -` until `"failed": 0`.
- On a large API, use `npx hono snapshot --status-only` and `npx hono batch - --compact` — same loop, much smaller output. A failed step carries `diff`: fix what it names.

---

## Hono API Reference

### App Constructor

```ts
import { Hono } from 'hono'

const app = new Hono()

// With TypeScript generics
type AppEnv = {
  Bindings: Env // generated from the bindings config — see "Environment (Cloudflare Workers)"
  Variables: { user: User }
}
const app = new Hono<AppEnv>()
```

### Routing Methods

```ts
app.get('/path', handler)
app.post('/path', handler)
app.put('/path', handler)
app.delete('/path', handler)
app.patch('/path', handler)
app.options('/path', handler)
app.all('/path', handler) // all HTTP methods
app.on('PURGE', '/path', handler) // custom method
app.on(['PUT', 'DELETE'], '/path', handler) // multiple methods
```

### Routing Patterns

```ts
// Path parameters
app.get('/user/:name', (c) => {
  const name = c.req.param('name')
  return c.json({ name })
})

// Multiple params
app.get('/posts/:id/comments/:commentId', (c) => {
  const { id, commentId } = c.req.param()
})

// Optional parameters
app.get('/api/animal/:type?', (c) => c.text('Animal!'))

// Wildcards
app.get('/wild/*/card', (c) => c.text('Wildcard'))

// Regexp constraints
app.get('/post/:date{[0-9]+}/:title{[a-z]+}', (c) => {
  const { date, title } = c.req.param()
})

// Chained routes
app
  .get('/endpoint', (c) => c.text('GET'))
  .post((c) => c.text('POST'))
  .delete((c) => c.text('DELETE'))
```

### Route Grouping

```ts
// Using route()
const api = new Hono()
api.get('/users', (c) => c.json([]))

const app = new Hono()
app.route('/api', api) // mounts at /api/users

// Using basePath()
const app = new Hono().basePath('/api')
app.get('/users', (c) => c.json([])) // GET /api/users
```

### Error Handling

```ts
app.notFound((c) => c.json({ message: 'Not Found' }, 404))

app.onError((err, c) => {
  console.error(err)
  return c.json({ message: 'Internal Server Error' }, 500)
})
```

---

## Context (c)

### Response Methods

```ts
c.text('Hello') // text/plain
c.json({ message: 'Hello' }) // application/json
c.html('<h1>Hello</h1>') // text/html
c.redirect('/new-path') // 302 redirect
c.redirect('/new-path', 301) // 301 redirect
c.body('raw body', 200, headers) // raw response
c.notFound() // 404 response
```

### Headers & Status

```ts
c.status(201)
c.header('X-Custom', 'value')
c.header('Cache-Control', 'no-store')
```

### Variables (request-scoped data)

```ts
// In middleware
c.set('user', { id: 1, name: 'Alice' })

// In handler
const user = c.get('user')
// or
const user = c.var.user
```

### Environment (Cloudflare Workers)

Never write the `Bindings` type by hand: generate it, and rerun after every change to the bindings. Which tool a project uses is visible from its files.

**`cf` CLI projects** (`cloudflare.config.ts`; this is what `create-hono` generates for `cloudflare-workers` now): bindings are declared in `cloudflare.config.ts` with the `bindings` helpers from `cf/config`. Local bindings need no IDs; do not create Cloudflare resources (KV, D1, R2) or set IDs unless the user asks. Then run `npm run typecheck` (= `cf workers types && tsc`; it writes `.cloudflare/types`) and use the generated `Env`:

```ts
// cloudflare.config.ts
import { bindings, defineConfig } from 'cf/config

