# Project Blueprint

## Overview
This is a Vue.js application using Vite, TypeScript, and Vue Router, focused on devotion and prayer. The site is styled using a unified, premium "Sacred Parchment" design language that features warm, tactile tones, elegant serif typography, and reusable layout components.

## Current Styles & Design Systems
- **Color Palette**:
  - Primary (Gold): `#C5A059`
  - Secondary (Crimson): `#8B2635`
  - Tertiary (Slate Blue): `#8FA5D8`
  - Neutral (Charcoal): `#2D2926`
  - Background (Page): `#FDFBF7`
  - Background (Card): `#F5EFE6`
  - Border (Subtle): `#E8DFD3`
- **Typography**:
  - Headings/Quotes: `Libre Caslon Text` (Serif)
  - Body/Labels: `Source Sans 3` (Sans-serif)
- **Visual Features**: Soft shadows, custom borders, responsive grid layouts, glassmorphic accents where applicable, and bilingual toggle support.

---

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
    -   **Bilingual Display**: Automatically displays English reference text below Latin prayers when the Latin toggle is active.

3.  **Expanded Prayer Library**:
    -   Added Act of Contrition (English & Latin) to the centralized data store.

### Catechism & Study
1.  **Reading Progress Tracking**:
    -   Implemented `localStorage` persistence to track read paragraphs.
    -   Added visual indicators (grayed out) for read content.

2.  **Interactive References**:
    -   Clicking a reference opens a modal displaying the content of the referenced paragraph(s).
    -   Supports single IDs and ranges (e.g., `2095-2109`).

3.  **Catechism Quick Lookup Bubble**:
    -   Created `CatechismBubble.vue` floating widget visible on the Bible Study and Daily Mass Readings pages.
    -   Allows looking up any Catechism paragraph (1-2865).
    -   Integrated interactive nested cross-references that automatically load inside the popover on click.
    -   Implements back button tracking for reference navigation history.

### System & Infrastructure
1.  **Daily Mass Readings**:
    -   Cloudflare Functions integration with Gemini API for dynamic daily readings and context.
2.  **Global Design System**:
    -   Unified glassmorphic theme across all views.
3.  **Routing**:
    -   Hash-based routing for static hosting compatibility.
4.  **API Authentication Integration**:
    -   Secure Bearer token-based client authentication logic implemented in `useDevotionApi.ts`.
5.  **Bible Study & Context Explorer**:
    -   Implemented a page for selecting Bible book, chapter, and verse.
    -   Integrated authentication token logic to fetch scripture and contextual explanations from backend APIs.
    -   Added persistent search history with expandable/collapsible commentary cards and option to clear history.
6.  **GitHub Actions Workflow**:
    - Fixed build and deploy workflow by setting correct Node.js version to 20.x and structuring environment variables.
7.  **API Base URL Validation & Fallback**:
    - Added validation to fallback `baseUrl` to `https://devotionapi.raenardcruz.com` if `VITE_API_BASE_URL` is misconfigured (e.g., set to the token value).
8.  **Daily Readings Context Markdown Rendering**:
    - Extracted the local markdown parsing function to a reusable shared utility `src/utils/markdown.ts`.
    - Integrated Markdown parsing for the contextual meditation overlays in the Daily Mass Readings view (`DailyMassReadings.vue`), rendering headers, bold, italics, and lists beautifully.

---

## Active Plan: None (All current items completed)

## Architecture
-   **Framework**: Vue 3 (Composition API)
-   **Tooling**: Vite, TypeScript
-   **Data Storage**:
    -   `src/data/prayers.json`: Centralized prayer database.
    -   `src/components/rosary/rosaryData.ts`: Logic for generating full Rosary sequences.
    -   `src/components/divinemercy/divineMercyData.ts`: Sequence definitions for the Chaplet.
    -   `localStorage`: Bible study history persistence.

## Component Structure
-   **Common/Global**:
    -   `TopNav.vue`: Main header navigation.
    -   `BottomNav.vue`: Main bottom tab navigation.
    -   `ParchmentCard.vue`: Reusable card wrapper.
    -   `AppButton.vue`: Button variants.
    -   `AppTabs.vue`: Tab selectors.
-   **Rosary**:
    -   `Rosary.vue`: Orchestration of the 80-step prayer sequence.
    -   `RosaryCard.vue`: Displays scripture and prayer text with language toggles.
    -   `RosaryBeads.vue`: Visual progress indicator.
-   **Divine Mercy**:
    -   `DivineMercy.vue`: Main interactive view.
    -   `DivineMercyBeads.vue`: Progress indicator for the 10-bead cycle.
    -   `PhaseLabel.vue`: Dynamic header showing the current prayer category and version.
-   **Catechism**:
    -   `Catechism.vue`, `CatechismText.vue`, `CatechismReference.vue`, `CatechismBubble.vue` [NEW]: Logic for the interactive Catechism explorer and floating quick lookup.
-   **Resources**:
    -   `Resources.vue`: Interactive markdown resource reader.
-   **Views**:
    -   `BibleStudy.vue` [NEW]: Bible verse context explorer with client-side markdown conversion and local history caching.
