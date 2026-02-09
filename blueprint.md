# Project Blueprint

## Overview
This is a Vue.js application using Vite, TypeScript, and Vue Router. The current goal is to resolve a reported syntax error in the router configuration.

## Current Plan
The user wants to track reading progress in the Catechism view.

1.  **Goal**: Allow users to mark paragraphs as read and persist this state.
2.  **Implementation**:
    -   Use `localStorage` to store a list of read paragraph IDs.
    -   Update `Catechism.vue` to include toggle buttons and visual indicators.
3.  **Goal**: Standardize styling across all pages for a consistent, premium feel.
    -   **Implementation**:
        -   Consolidate animations into global CSS.
        -   Update `Catechism.vue` header to match glassmorphic style.
        -   Ensure consistent border-radius and shadows.
4.  **Status**: Planning phase. `implementation_plan.md` created.

## Architecture
- **Framework**: Vue 3 (Composition API)
- **Tooling**: Vite, TypeScript
- **Routing**: Vue Router 4
