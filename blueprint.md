# Project Blueprint

## Overview
This is a Vue.js application using Vite, TypeScript, and Vue Router. The current goal is to resolve a reported syntax error in the router configuration.

## Current Plan
The user reported "Cannot find module '../views/Home.vue'" in `src/router/index.ts`.
1.  Verify file existence (Done).
2.  Update `src/vite-env.d.ts` to include `declare module '*.vue'` shim.
3.  Verify the fix.

## Architecture
- **Framework**: Vue 3 (Composition API)
- **Tooling**: Vite, TypeScript
- **Routing**: Vue Router 4
