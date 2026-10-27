---
name: tailwind-css-patterns
description: Provides comprehensive Tailwind CSS utility-first styling patterns including responsive design, layout utilities, flexbox, grid, spacing, typography, colors, and modern CSS best practices. Use when styling React/Vue/Svelte components, building responsive layouts, implementing design systems, or optimizing CSS workflow.
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

# Tailwind CSS Development Patterns

Expert guide for building modern, responsive user interfaces with Tailwind CSS utility-first framework. Covers v4.1+ features including CSS-first configuration, custom utilities, and enhanced developer experience.

## Overview

Provides actionable patterns for responsive, accessible UIs with Tailwind CSS v4.1+. Covers utility composition, dark mode, component patterns, and performance optimization.

## When to Use

- Styling React/Vue/Svelte components
- Building responsive layouts and grids
- Implementing design systems
- Adding dark mode support
- Optimizing CSS workflow

## Quick Reference

### Responsive Breakpoints

| Prefix | Min Width | Description |
|--------|-----------|-------------|
| `sm:` | 640px | Small screens |
| `md:` | 768px | Tablets |
| `lg:` | 1024px | Desktops |
| `xl:` | 1280px | Large screens |
| `2xl:` | 1536px | Extra large |

### Common Patterns

```html
<!-- Center content -->
<div class="flex items-center justify-center min-h-screen">
  Content
</div>

<!-- Responsive grid -->
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
  <!-- Items -->
</div>

<!-- Card component -->
<div class="bg-white rounded-lg shadow-lg p-6">
  <h3 class="text-xl font-bold">Title</h3>
  <p class="text-gray-600">Description</p>
</div>
```

## Instructions

1. **Start Mobile-First**: Write base styles for mobile, add responsive prefixes (`sm:`, `md:`, `lg:`) for larger screens
2. **Use Design Tokens**: Leverage Tailwind's spacing, color, and typography scales
3. **Compose Utilities**: Combine multiple utilities for complex styles
4. **Extract Components**: Create reusable component classes for repeated patterns
5. **Configure Theme**: Customize design tokens in `tailwind.config.js` or using `@theme`
6. **Verify Changes**: Test at each breakpoint using DevTools responsive mode. Check for visual regressions and accessibility issues before committing.

## Examples

### Responsive Card Component

```tsx
function ProductCard({ product }: { product: Product }) {
  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden sm:flex">
      <img className="h-48 w-full object-cover sm:h-auto sm:w-48" src={product.image} />
      <div className

