# Tailwind CSS v4 Theme Directives & Arbitrary Values Prevention

## Using `@theme` Variables over Hardcoded Hex Values

```css
@import "tailwindcss";

@theme {
  --color-brand-primary: oklch(0.65 0.24 280);
  --color-brand-secondary: oklch(0.85 0.15 150);
}
```

## Class Ordering Hierarchy

1. Display & Layout (`flex`, `grid`, `hidden`, `absolute`)
2. Box Model (`p-*`, `m-*`, `w-*`, `h-*`)
3. Typography (`text-*`, `font-*`, `leading-*`)
4. Visuals (`bg-*`, `border-*`, `rounded-*`, `shadow-*`)
5. Interactive states (`hover:*`, `focus-visible:*`, `disabled:*`)
