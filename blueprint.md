# Project Blueprint

## Overview
This is a Vue.js application using Vite, TypeScript, and Vue Router, focused on devotion and prayer.

## Recent Features

### Devotional Experiences
1.  **Scriptural Rosary**:
    -   **Full Prayer Sequence**: Implemented an 80-step sequence covering Intro, 5 Mysteries (10 beads each), and Closing.
    -   **Interactive Beads**: Dynamic visualization of progress based on the current step type (Intro vs. Decade).
    -   **Latin/English Toggle**: Real-time language switching for all traditional prayers while maintaining English scripture meditation.
    -   **Premium Styling**: Glassmorphic dark theme with amber/rose accents and smooth transitions.

2.  **Divine Mercy Chaplet**:
    -   **Short Version Toggle**: Added a toggle to skip to closing prayers after the first decade, ideal for quick devotion.
    -   **Centralized Data**: Migrated all prayer text to a single JSON source of truth.
    -   **Refined Navigation**: Enhanced `next`/`prev` logic to handle complex transitions between decades and closing sequences.
    -   **UI Indicators**: Added "Short Version" badges and persistent decade progress tracking.

### Catechism & Study
1.  **Reading Progress Tracking**:
    -   Implemented `localStorage` persistence to track read paragraphs.
    -   Added visual indicators (grayed out) for read content.

2.  **Interactive References**:
    -   Clicking a reference opens a modal displaying the content of the referenced paragraph(s).
    -   Supports single IDs and ranges (e.g., `2095-2109`).

### System & Infrastructure
1.  **Daily Mass Readings**:
    -   Cloudflare Functions integration with Gemini API for dynamic daily readings and context.
2.  **Global Design System**:
    -   Unified glassmorphic theme across all views.
    -   Responsive navigation and interactive controls.
3.  **Routing**:
    -   Hash-based routing for static hosting compatibility.

## Architecture
-   **Framework**: Vue 3 (Composition API)
-   **Tooling**: Vite, TypeScript
-   **Data Storage**:
    -   `src/data/prayers.json`: Centralized prayer database.
    -   `src/components/rosary/rosaryData.ts`: Logic for generating full Rosary sequences.
    -   `src/components/divinemercy/divineMercyData.ts`: Sequence definitions for the Chaplet.

## Component Structure
-   **Rosary**:
    -   `Rosary.vue`: Orchestration of the 80-step prayer sequence.
    -   `RosaryCard.vue`: Displays scripture and prayer text with language toggles.
    -   `RosaryBeads.vue`: Visual progress indicator using canvas-like or SVG beads.
-   **Divine Mercy**:
    -   `DivineMercy.vue`: Main interactive view with version and language controls.
    -   `DivineMercyBeads.vue`: Progress indicator for the 10-bead cycle.
    -   `PhaseLabel.vue`: Dynamic header showing the current prayer category and version.
-   **Catechism**:
    -   `Catechism.vue`, `CatechismText.vue`, `CatechismReference.vue`: Logic for the interactive Catechism explorer.
