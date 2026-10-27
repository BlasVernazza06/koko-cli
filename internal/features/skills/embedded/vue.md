---
name: vue-best-practices
description: MUST be used for Vue.js tasks. Strongly recommends Composition API with `<script setup>` and TypeScript as the standard approach. Covers Vue 3, SSR, Volar, vue-tsc. Load for any Vue, .vue files, Vue Router, Pinia, or Vite with Vue work. ALWAYS use Composition API unless the project explicitly requires Options API.
license: MIT
metadata:
  author: github.com/vuejs-ai
  version: "18.0.0"
---

# Vue Best Practices Workflow

Use this skill as an instruction set. Follow the workflow in order unless the user explicitly asks for a different order.

## Core Principles
- **Keep state predictable:** one source of truth, derive everything else.
- **Make data flow explicit:** Props down, Events up for most cases.
- **Favor small, focused components:** easier to test, reuse, and maintain.
- **Avoid unnecessary re-renders:** use computed properties and watchers wisely.
- **Readability counts:** write clear, self-documenting code.

## 1) Confirm architecture before coding (required)

- Default stack: Vue 3 + Composition API + `<script setup lang="ts">`.
- If the project explicitly uses Options API, load `vue-options-api-best-practices` skill if available.
- If the project explicitly uses JSX, load `vue-jsx-best-practices` skill if available.

### 1.1 Must-read core references (required)

- Before implementing any Vue task, make sure to read and apply these core references:
  - `references/reactivity.md`
  - `references/sfc.md`
  - `references/component-data-flow.md`
  - `references/composables.md`
- Keep these references in active working context for the entire task, not only when a specific issue appears.

### 1.2 Plan component boundaries before coding (required)

Create a brief component map before implementation for any non-trivial feature.

- Define each component's single responsibility in one sentence.
- Keep entry/root and route-level view components as composition surfaces by default.
- Move feature UI and feature logic out of entry/root/view components unless the task is intentionally a tiny single-file demo.
- Define props/emits contracts for each child component in the map.
- Prefer a feature folder layout (`components/<feature>/...`, `composables/use<Feature>.ts`) when adding more than one component.

## 2) Apply essential Vue foundations (required)

These are essential, must-know foundations. Apply all of them in every Vue task using the core references already loaded in section `1.1`.

### Reactivity

- Must-read reference from `1.1`: [reactivity](references/reactivity.md)
- Keep source state minimal (`ref`/`reactive`), derive everything possible with `computed`.
- Use watchers for side effects if needed.
- Avoid recomputing expensive logic in templates.

### SFC structure and template safety

- Must-read reference from `1.1`: [sfc](references/sfc.md)
- Keep SFC sections in this order: `<script>` → `<template>` → `<style>`.
- Keep SFC responsibilities focused; split large components.
- Keep templates declarative; move branching/derivation to script.
- Apply Vue template safety rules (`v-html`, list rendering, conditional rendering choices).

### Keep components focused

Split a component when it has **more than one clear responsibility** (e.g. data orchestration + UI, or multiple independent UI sections).

- Prefer **smaller components + composables** over one “mega component”
- Move **UI sections** into child componen

