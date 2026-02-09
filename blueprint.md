# Project Blueprint

## Overview
This is a Vue.js application using Vite, TypeScript, and Vue Router, focused on devotion and prayer.

## Recent Features

### System
1.  **Routing**:
    -   Switched from HTML5 History Mode to Hash Mode for better compatibility with static hosting environments.

### Catechism
1.  **Reading Progress Tracking**:
    -   Implemented `localStorage` persistence to track read paragraphs.
    -   Added toggle buttons (checkmarks) to paragraphs for easy tracking.
    -   Added visual indicators (grayed out) for read content.

2.  **Interactive References**:
    -   Replaced plain text references `(843)` with interactive components.
    -   Clicking a reference opens a modal displaying the content of the referenced paragraph(s).
    -   Supports single IDs and ranges (e.g., `2095-2109`).

3.  **Styling**:
    -   Consolidated animations into global CSS.
    -   Implemented glassmorphic header and modal styles.
    -   Added custom scrollbars for modal content.

## Architecture
-   **Framework**: Vue 3 (Composition API)
-   **Tooling**: Vite, TypeScript
-   **Routing**: Vue Router 4 (Hash Mode)
-   **State Management**: Composition API + LocalStorage (for simple persistence)

## Component Structure
-   `Catechism.vue`: Main view, handles state and modal logic.
-   `CatechismText.vue`: Parsers paragraph text to inject interactive references.
-   `CatechismReference.vue`: Handles individual reference logic and click events.
