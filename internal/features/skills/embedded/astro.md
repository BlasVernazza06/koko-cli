---
name: astro
description: Skill for building with the Astro web framework. Helps create Astro components and pages, configure SSR adapters, set up content collections, deploy static sites, and manage project structure and CLI commands. Use when the user needs to work with Astro, mentions .astro files, asks about static site generation (SSG), islands architecture, content collections, or deploying an Astro project.
license: MIT
metadata: 
  authors: "Astro Team"
  version: "0.0.1"
---

# Astro Usage Guide

**Always consult [docs.astro.build](https://docs.astro.build) for code examples and latest API.**

Astro is the web framework for content-driven websites.

---

## Quick Reference

### File Location
CLI looks for `astro.config.js`, `astro.config.mjs`, `astro.config.cjs`, and `astro.config.ts` in: `./`. Use `--config` for custom path.

### CLI Commands

- `npx astro dev` - Start the development server.
- `npx astro build` - Build your project and write it to disk.
- `npx astro check` - Check your project for errors.
- `npx astro add` - Add an integration.
- `npx astro sync` - Generate TypeScript types for all Astro modules.

**Re-run after adding/changing plugins.**

### Project Structure

Reference [project structure docs](https://docs.astro.build/en/basics/project-structure).

- `src/*` - Project source code (components, pages, styles, images, etc.)
- `src/pages` - **Required.** Defines all pages and routes.
- `src/components` - Components (convention, not required).
- `src/layouts` - Layout components (convention, not required).
- `src/styles` - CSS/Sass files (convention, not required).
- `public/*` - Non-code, unprocessed assets (fonts, icons, etc.); copied as-is to build output.
- `package.json` - Project manifest.
- 

