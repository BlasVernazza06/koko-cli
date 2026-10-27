---
name: sleek-design-mobile-apps
description: Design mobile app screens with Sleek, edit Sleek projects, and implement their designs in React Native or HTML.
metadata:
  requires-env: SLEEK_API_KEY
---

# Designing with Sleek

[![Design mobile apps in minutes](https://raw.githubusercontent.com/sleekdotdesign/agent-skills/main/assets/hero.png)](https://sleek.design)

## Overview

[sleek.design](https://sleek.design) provides a REST API at `/api/v1/*` to create projects, design screens from a brief, and fetch their HTML and renders.

**Base URL**: `https://sleek.design`
**Auth**: `Authorization: Bearer $SLEEK_API_KEY`, except the public device start/poll endpoints
**Content-Type**: send `application/json` for JSON request bodies.
**CORS**: Enabled on all `/api/v1/*` endpoints
**Envelope**: check HTTP status first. Successful JSON wraps its payload in `data`; lists also include sibling `pagination`. Errors use top-level `code`, `status` and `message`, with optional `data.url`. Successful screenshots and share cards return image bytes; `204` has no body.
**Ids**: opaque 11-character tokens like `V1StGXR8Z5j`, with no type prefix. Nothing about an id says whether it names a project, screen, component, version or theme, so track each one by the response field it arrived in.
**Parsing responses**: write the body to a file (`curl -o run.json`) and parse the file. Piping JSON through `echo` corrupts it: zsh expands the escaped `\n` inside string values into real newlines, which makes the body invalid JSON.
**API docs**: OpenAPI spec at `https://sleek.design/api/v1/spec.json`; browsable docs at `https://sleek.design/api/v1/docs`. Fetch the spec for any contract detail not covered here.

---

## Designing or implementing

A single run can take both directions — "design me a fitness app and build it in Expo" is a design pass followed by an implementation pass. Gate each pointer on the artefact in front of you rather than on a one-time reading of the request:

**Designing** — you are about to send a chat message to Sleek. Read [designing.md](designing.md) first; it carries the create → chat → screenshot loop, how to author a style direction, and when to seed one from a reference instead.

**Implementing.** Before writing code from a Sleek design, read [implementing.md](implementing.md) for version selection, native conversion, working features, and acceptance checks.

## User shot and review shot

Screenshots come in two framings, and each file below asks for one by name:

- A **user shot** is the default render: viewport height, phone-shaped. This is the one you show the user.
- A **review shot** adds `fullHeight: true`, one component per request. This is the one you judge your own work against.

---

## Prerequisites: API Key

If `SLEEK_API_KEY` is not set, use the device flow so the user never handles the raw key:

1. `POST https://sleek.design/api/v1/device/start` (no auth) with body `{"source": "your-tool-slug"}`. The response contains a `verificationUrl`, a human-checkable `userCode`, a secret `deviceCode`, and a poll `interval` in seconds.
2. Show the user the `verificationUrl` and the `userCode`, and tell them to confirm the code matches before approving.
3. Poll `POST https://sleek.design/api/v1/device/poll` with `{"deviceCode": "..."}` every `interval` seconds. On `pending`, keep waiting. On `approved`, store `data.key` pr

