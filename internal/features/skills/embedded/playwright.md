---
name: playwright-best-practices
description: Use when writing Playwright tests, fixing flaky tests, debugging failures, implementing Page Object Model, configuring CI/CD, optimizing performance, mocking APIs, handling authentication or OAuth, testing accessibility (axe-core), file uploads/downloads, date/time mocking, WebSockets, geolocation, permissions, multi-tab/popup flows, mobile/responsive layouts, touch gestures, GraphQL, error handling, offline mode, multi-user collaboration, third-party services (payments, email verification), console error monitoring, global setup/teardown, test annotations (skip, fixme, slow), test tags (@smoke, @fast, @critical, filtering with --grep), project dependencies, security testing (XSS, CSRF, auth), performance budgets (Web Vitals, Lighthouse), iframes, component testing, canvas/WebGL, service workers/PWA, test coverage, i18n/localization, Electron apps, or browser extension testing. Covers E2E, component, API, visual, accessibility, security, Electron, and extension testing.
license: MIT
metadata:
  author: currents.dev
  version: "1.2"
---

# Playwright Best Practices

This skill provides comprehensive guidance for all aspects of Playwright test development, from writing new tests to debugging and maintaining existing test suites.

## Activity-Based Reference Guide

Consult these references based on what you're doing:

### Writing New Tests

**When to use**: Creating new test files, writing test cases, implementing test scenarios

| Activity                            | Reference Files                                                                                                                               |
| ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| **Writing E2E tests**               | [test-suite-structure.md](core/test-suite-structure.md), [locators.md](core/locators.md), [assertions-waiting.md](core/assertions-waiting.md) |
| **Writing component tests**         | [component-testing.md](testing-patterns/component-testing.md), [test-suite-structure.md](core/test-suite-structure.md)                        |
| **Writing API tests**               | [api-testing.md](testing-patterns/api-testing.md), [test-suite-structure.md](core/test-suite-structure.md)                                    |
| **Writing GraphQL tests**           | [graphql-testing.md](testing-patterns/graphql-testing.md), [api-testing.md](testing-patterns/api-testing.md)                                  |
| **Writing visual regression tests** | [visual-regression.md](testing-patterns/visual-regression.md), [canvas-webgl.md](testing-patterns/canvas-webgl.md)                            |
| **Structuring test code with POM**  | [page-object-model.md](core/page-object-model.md), [test-suite-structure.md](core/test-suite-structure.md)                                    |
| **Setting up test data/fixtures**   | [fixtures-hooks.md](core/fixtures-hooks.md), [test-data.md](core/test-data.md)                                                                |
| **Handling authentication**         | [authentication.md](advanced/authentication.md), [authentication-flows.md](advanced/authentication-flows.md)                                  |
| **Testing date/time features**      | [clock-mocking.md](advanced/clock-mocking.md)                                                                                                 |
| **Testing file upload/download**    | [file-operations.md](testing-patterns/file-operations.md), [file-upload-download.md](testing-patterns/file-upload-download.md)                |
| **Testing forms/validation**        | [forms-validation.md](testing-patterns/forms-validation.md)                                                                                   |
| **Testing drag and drop**           | [drag-drop.md](testing-patterns/drag-drop.md)                                                                                                 |
| **Testing accessibility**           | [accessibility.md](testing-patterns/accessibility.md)                                                                                         |
| **Testing security (XSS, CSRF)**    | [security-testing.md](testing-patterns/security-testing.md)                                                                                   |
| **Using test annotations**          | [annotations.md](core/annotations.md)                                                                                                         |
| **Using test tags**                 | [test-tags.md](core/test-tags.md)                                    

