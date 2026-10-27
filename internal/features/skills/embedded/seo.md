---
name: seo
description: Optimize for search engine visibility and ranking. Use when asked to "improve SEO", "optimize for search", "fix meta tags", "add structured data", "sitemap optimization", or "search engine optimization".
license: MIT
metadata:
  author: web-quality-skills
  version: "2.0"
---

# SEO optimization

Search engine optimization based on Lighthouse SEO audits and Google Search guidelines. Focus on technical SEO, on-page optimization, and structured data.

## Evidence-led audit workflow

When a rendered page is available:

1. Run live Lighthouse SEO and Agentic Browsing checks when that capability is available; with Chrome DevTools MCP, use `lighthouse_audit`. Use the results to localize rendered-page failures.
2. Inspect signals Lighthouse cannot establish on its own: response headers, redirects, `robots.txt`, sitemap coverage, canonical consistency across page templates, structured-data eligibility, and Search Console evidence when the user provides access.
3. Separate technical crawl/index findings from content quality and authority. Do not invent ranking-factor weights or promise ranking changes.
4. Fix the source and re-run the same checks. For indexation or ranking outcomes, report that search-engine validation remains pending.

If live tools are unavailable, use category-specific Lighthouse CLI output plus direct source and HTTP inspection. A Lighthouse SEO score covers a useful subset of technical checks; it is not a prediction of rankings.

| Area | What this skill can verify |
|------|----------------------------|
| Crawl and index controls | Technical configuration and consistency |
| Rendered metadata and semantics | Presence, validity, and page-template issues |
| Structured data | Syntax and eligibility signals, not guaranteed rich results |
| Core Web Vitals | Link to measured field/lab evidence from the Core Web Vitals skill |
| Content usefulness and authority | Review quality, but do not assign synthetic ranking percentages |

---

## Technical SEO

### Crawlability

**robots.txt:**
```text
# /robots.txt
User-agent: *
Allow: /

# Block admin/private areas
Disallow: /admin/
Disallow: /api/
Disallow: /private/

# Don't block resources needed for rendering
# ❌ Disallow: /static/

Sitemap: https://example.com/sitemap.xml
```

**Meta robots:**
```html
<!-- Default: indexable, followable -->
<meta name="robots" content="index, follow">

<!-- Noindex specific pages -->
<meta name="robots" content="noindex, nofollow">

<!-- Indexable but don't follow links -->
<meta name="robots" content="index, nofollow">

<!-- Control snippets -->
<meta name="robots" content="max-snippet:150, max-image-preview:large">
```

**Canonical URLs:**
```html
<!-- Prevent duplicate content issues -->
<link rel="canonical" href="https://example.com/page">

<!-- Self-referencing canonical (recommended) -->
<link rel="canonical" href="https://example.com/current-page">

<!-- For paginated content -->
<link rel="canonical" href="https://example.com/products">
<!-- Or use rel="prev" / rel="next" for explicit pagination -->
```

### XML sitemap

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2024-01-15</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://example.com/products</loc>
    <lastmod>2024-01-14</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
</urlset>
```

**Sitemap best practices:**
- Maximum 50,000 URLs or 50MB per sitemap
- Use sitemap index for larger sites
- Include only canonical, indexable URLs
- Update `lastmod` when content changes
- Submit to Google Search Console

### URL structure

```
✅ Good URLs:
https://example.com/products/blue-widget
https://example.com/blog/how-to-use-widgets

❌ Poor URLs:
https://example.com/p?id=12345
https://example.com/products/item/category/subcategory/blue-widget-2024-sale-discount
```

**URL guidelines:**
- Use hyphens, not underscores
- Lowercase only
- Keep short (< 75 characters)
- Include target keywords naturally
- Avoid parameters when possible
- Use HTTPS always

### HTTPS & security

```html
<!-- Ensure all resources use HTTPS -->
<img src="https://example.com/image.jpg">

<!-- Not: -->
<img src="http://example.com/image.jpg">
```

**Security headers for SEO trust signals:**
```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
```

---

## On-page SEO

### Title tags

```html
<!-- ❌ Missing or generic -->
<title>Page</title>
<title>Home</title>

<!-- ✅ Descriptive with primary keyword -->
<title>Blue Widgets for Sale | Premium Quality | Example Store</title>
```

**Title tag guidelines:**
- Use 50–60 characters only as a rough linting proxy, not a pass/fail limit. Google truncates title links to fit the rendered device width, so preview width when tooling supports it.
- Describe the page topic naturally near the beginning
- Unique for every page
- Add the brand when it helps users distinguish the result
- Action-oriented when appropriate

Treat title-link rewriting separately from truncation. Google may build a different title link from the visible page title, headings, anchor text, and other sources even when the `<title>` is short; investigate accuracy and consistency rather than shortening it automatically. See [Google's title-link guidance](https://developers.google.com/search/docs/appearance/title-link).

### Meta descriptions

```html
<!-- ❌ Missing or duplicate -->
<meta nam

