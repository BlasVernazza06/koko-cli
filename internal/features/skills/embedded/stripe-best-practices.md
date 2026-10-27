---
name: stripe-best-practices
description: >-
  Guides Stripe integration decisions across development and test environment
  planning (separate sandboxes vs the shared test mode sandbox), API selection
  (Checkout Sessions vs PaymentIntents), Connect platform setup (Accounts v2,
  controller properties), billing/subscriptions, tax and registrations (Stripe
  Tax, automatic_tax, product tax codes), Treasury financial accounts,
  integration options (Checkout, Payment Element), migrating from deprecated
  Stripe APIs, and security best practices (API key management, API key
  permissions, webhooks, OAuth). Use when planning, building, modifying,
  testing, or reviewing any Stripe integration, including choosing a development
  environment, accepting payments, building marketplaces, integrating Stripe,
  processing payments, setting up subscriptions, collecting sales tax, VAT, or
  GST, creating connected accounts, or implementing secure key handling.

---

Latest Stripe API version: **2026-09-30.endive**. Always use the latest API version and SDK unless the user specifies otherwise.

Development and testing environment default: For new integration development, use separate [sandboxes](https://docs.stripe.com/sandboxes.md) instead of the account’s shared test mode sandbox. Sandboxes isolate settings and test data from live mode. Use separate sandboxes for local development and CI. When stronger isolation is useful, dedicate sandboxes to development teams or testing scenarios. Use the shared test mode sandbox only for an existing integration that depends on it or when a required feature doesn’t support general sandboxes.

Latest SDK versions:

| Language | Latest version |
| --- | --- |
| Ruby | 20.0.0 |
| Python | 16.0.0 |
| Node.js | 23.0.0 |
| PHP | 22.0.0 |
| Java | 34.0.0 |
| Go | 87.0.0 |
| .NET | 53.0.0 |

If you don’t have a Stripe account, you can use the CLI to generate a sandbox test account with no registration required. Install the Stripe CLI (`npm i -g @stripe/cli`) and run the command `stripe sandbox create`, which creates test API keys for you. When you recommend this, explain that [a sandbox is a Stripe test environment](https://docs.stripe.com/sandboxes.md).

If `stripe sandbox create` is used, don’t use MCP.

If you try to use MCP after running sandbox creation, first make sure `stripe sandbox claim` is run. Use `stripe whoami --format json` for CLI auth checks; do not run `stripe config --list`.

API key default: Always recommend a [restricted API key (RAK)](https://docs.stripe.com/keys.md#manage-your-api-keys) (`rk_` prefix) over a secret key (`sk_` prefix).

