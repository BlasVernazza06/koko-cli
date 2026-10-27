---
name: expo-native-ui
description: Build beautiful, native-feeling Expo screens. Covers Apple HIG styling, semantic colors, native controls, SF Symbols, media, visual effects, gradients, storage, and responsive layout. For routing and navigation, use the expo-router skill; for motion and animation, use the expo-animation skill.
version: 1.1.1
license: MIT
---

# Expo Native UI Guidelines

For routes, links, stacks, tabs, modals, sheets, and headers, use the `expo-router` skill. For any motion — entering/exiting, gestures, springs, keyboard-driven UI — use the `expo-animation` skill.

> **Before picking any UI component, check `expo-ui` first.** `@expo/ui` provides native equivalents — BottomSheet, Button, Picker, Slider, Menu, Section, Switch, SegmentedControl, and more — rendered as real SwiftUI on iOS and Jetpack Compose on Android, available in Expo Go on SDK 56+ with no custom build. Load the **`expo-ui`** skill to find the right component before falling back to React Native built-ins or community libraries. This skill (`expo-native-ui`) covers the surrounding structure: Expo Router navigation, layout, styling, and visual effects.

## References

Consult these resources as needed:

```
references/
  controls.md            Native iOS: Switch, Slider, SegmentedControl, DateTimePicker, Picker
  gradients.md           CSS gradients via experimental_backgroundImage (New Arch only)
  icons.md               SF Symbols via expo-symbols SymbolView: names, weights, animations; Material icons on Android
  media.md               Camera, audio, video, and file saving
  storage.md             SQLite, AsyncStorage, SecureStore
  visual-effects.md      Blur (expo-blur) and liquid glass (expo-glass-effect)
  webgpu-three.md        3D graphics, games, GPU visualizations with WebGPU and Three.js
```

## Running the App

**CRITICAL: Always try Expo Go first before creating custom builds.**

Most Expo apps work in Expo Go without any custom native code. Before running `npx expo run:ios` or `npx expo run:android`:

1. **Start with Expo Go**: Run `npx expo start` and scan the QR code with Expo Go
2. **Check if features work**: Test your app thoroughly in Expo Go
3. **Only create custom builds when required** - see below

### When Custom Builds Are Required

You need `npx expo run:ios/android` or `eas build` ONLY when using:

- **Local Expo modules** (custom native code in `modules/`)
- **Apple targets** (widgets, app clips, extensions via `@bacons/apple-targets`)
- **Third-party native modules** not included in Expo Go
- **Custom native configuration** that can't be expressed in `app.json`

### When Expo Go Works

Expo Go supports a wide range of features out of the box:

- Most `expo-*` packages (camera, location, sensors, sqlite, etc.) — but not all: remote push notifications don't work in Expo Go on Android since SDK 53, and some packages need native capabilities Expo Go doesn't bundle (e.g. WebGPU — see `references/webgpu-three.md`)
- Expo Router navigation and deep links
- Most UI libraries (reanimated, gesture handler, etc.)

**If you're unsure, try Expo Go first.** Creating custom builds adds complexity, slower iteration, and requires Xcode/Android Studio setup.

## Code Style

- Be cautious of unterminated strings. Ensure nested backticks are escaped; never forget to escape quotes correctly.
- Always use import statements at the top of the file.
- Always use kebab-case for file names, e.g. `comment-card.tsx`
- Never use special characters in file names
- Configure tsconfig.json with path aliases, and prefer aliases over relative imports for refactors.

## Library Preferences

- **For any sheet, picker, slider, toggle, menu, or grouped-form section: use `@expo/ui` (see `expo-ui` skill) before reaching for a React Native built-in or community library** — it renders native SwiftUI/Compose and works in Expo Go on SDK 56+. For grouped/settings-style rows (short, fixed-length), use `@expo/ui`'s `List` + `ListItem`. For large or unknown-length scrolling lists (feeds, search results, catalogs), use `FlatList` or `FlashList` — `@expo/ui`'s `List` is not virtualized.
- Never use modules removed from React Native such as Picker, WebView, SafeAreaView, or AsyncStorage
- Never use legacy expo-permissions
- `expo-audio` not `expo-av`
- `expo-video` not `expo-av`
- `expo-symbols` (`SymbolView`) for SF Symbols on iOS, not `@expo/vector-icons` — see `references/icons.md`. SF Symbols are Apple-only: on Android every icon needs a Material source (`md` prop on NativeTabs triggers; in-screen options under "Android: Material Icons" in icons.md), never SF-only iconography
- `react-native-safe-area-context` not react-native SafeAreaView
- `process.env.EXPO_OS` not `Platform.OS`
- `React.use` not `React.useContext`
- `expo-image` Image component instead of intrinsic element `img`
- `expo-glass-effect` for liquid glass backdrops
- `Color` from `expo-router` for native semantic colors, not raw `PlatformColor` (type-safe, auto-adapts to light/dark)
- In SDK 56+, never import from `@react-navigation/*` directly — use `expo-router/react-navigation` instead (covers `@react-navigation/native`, `/core`, `/elements`, `/routers`)

## Responsiveness

- Wrap screens with scrollable content in a ScrollView. Screens whose

