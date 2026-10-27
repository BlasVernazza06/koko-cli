---
name: clerk-expo
description: Add Clerk authentication to Expo and React Native apps using @clerk/expo.
  Use for Expo setup, prebuilt native components (AuthView, UserButton), custom sign-in/sign-up
  flows (email, password, SMS/phone OTP, MFA), OAuth/SSO, native Google/Apple sign-in,
  Expo Router protected routes, biometrics, and push notifications. Do not use for
  native Swift/iOS, native Android/Kotlin, or web-only framework projects.
license: MIT
allowed-tools: WebFetch
metadata:
  author: clerk
  version: 2.0.0
compatibility: Requires @clerk/expo v3.4+ (written against v3.6.x, July 2026). Expo SDK 53-56, React Native 0.75+.
---

# Clerk Expo (React Native)

Implement Clerk in Expo / React Native projects. This skill inlines verified patterns for the stable surface (provider, token cache, flows) and requires source inspection of the installed `@clerk/expo` package for anything volatile (component props, hook signatures).

## Activation Rules

Activate when either is true:
- The user asks for auth in an Expo or React Native app, or mentions `@clerk/expo`, `ClerkProvider`, Expo Router auth, or Clerk hooks in a native app.
- The project is Expo/React Native (`app.json` / `app.config.js`, `expo` in `package.json`, `metro.config.js`, `@clerk/expo` dependency).

Route away when:
- Native iOS/Swift project (`.xcodeproj`, `Package.swift`) → `clerk-swift`
- Native Android/Kotlin project (`build.gradle` without React Native) → `clerk-android`
- Web-only framework (Next.js, Remix, plain React, etc.) → the matching framework skill

## Intent Map

Match what the user asked for, then load the reference(s) listed. Load only what the task needs.

| User intent (examples) | Path | Reference |
|------------------------|------|-----------|
| "Add auth to my app" / "add sign-in with Clerk" | Prebuilt native components (default) | references/setup.md + references/prebuilt-components.md |
| "Add auth" but Expo Go / web / custom UI required | Custom flows | references/setup.md + references/custom-flows.md |
| "Add phone / SMS auth", "email OTP", "passwordless" | Custom flow, `phoneCode` / `emailCode` | references/custom-flows.md |
| "Sign in with Google/Apple/GitHub", "social login", "SSO" | Browser SSO or native buttons | references/sso-and-native-auth.md |
| "MFA / 2FA / TOTP", "forgot password", "email link" | Custom flow additions | references/custom-flows.md |
| "Protect routes/screens", "redirect if signed out" | Expo Router guards | references/protected-routes.md |
| "Show user profile", "org switching", "push notifications", "sign out", "call my backend" | App recipes | references/recipes.md |
| "Biometric login", "Face ID", "passkeys" | Device features | references/recipes.md |

## Default Path Decision

When the user says "add auth" without specifying UI:

1. **Default to prebuilt native components** (`AuthView` + `UserButton` from `@clerk/expo/native`). Fastest to working auth; UI is maintained by Clerk. Tell the developer they are in beta and require a development build.
2. **Fall back to custom flows** when any of these hold — say why when you switch:
   - The project must run in Expo Go (no dev build).
   - The app targets web (native components don't render on web).
   - The developer wants their own UI or a specific brand experience beyond theming.
3. If the developer has an existing auth UI, extend what's there — don't rip out custom flows to insert `AuthView` (or vice versa) without being asked.

Do not blend prebuilt components and custom flows for the same auth step (e.g. `AuthView` plus a custom password form). Blending is allowed only when the developer explicitly asks.

## Quick Workflow

1. Confirm project type (Expo/RN) and pick the path per the Intent Map / Default Path rules.
2. Follow references/setup.md: install, env key, provider, token cache, config plugin, build type.
3. Verify dashboard prerequisites (Gate 2 and Gate 3 below).
4. Implement from the selected reference only.
5. Verify by building, not just by writing:
   - Run the project's typecheck (`npx tsc --noEmit` or equivalent).
   - Build and launch: `npx expo run:ios` / `run:android` for native features, `npx expo start` for Expo Go flows. If the build fails, fix and rebuild iteratively — build errors against the installed SDK are the ground truth when this skill and the SDK disagree. After ~5 failed fix attempts, stop and ask the developer how to proceed instead of thrashing.
   - Walk the developer through one real sign-in, then confirm the session survives an app restart (token cache working).

## Execution Gates (Do Not Skip)

1. **Publishable key** — Read from `process.env.EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY` (`.env` file). Never `NEXT_P

