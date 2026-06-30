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

### System & Infrastructure
1.  **Daily Mass Readings**:
    -   Cloudflare Functions integration with Gemini API for dynamic daily readings and context.
2.  **Global Design System**:
    -   Unified glassmorphic theme across all views.
    -   Responsive navigation and interactive controls.
3.  **Routing**:
    -   Hash-based routing for static hosting compatibility.

---

## Active Plan: Centralized "Sacred Parchment" Theme & Reusable Components

The objective of this plan is to unify the design of all pages under the "Sacred Parchment" style guide.

### Step-by-Step Action Plan:
1.  **Configure index.html & style.css**:
    -   Import `Libre Caslon Text` and `Source Sans 3` in `index.html`.
    -   Add custom color extensions (`parchment-*`) and font configurations to the Tailwind CDN setup inside `index.html`.
    -   Set page background and default styling in `src/style.css`.
2.  **Implement Central Reusable Components**:
    -   Create `src/components/common/TopNav.vue` (Global Header).
    -   Create `src/components/common/BottomNav.vue` (Global Footer / Tabs).
    -   Create `src/components/common/ParchmentCard.vue` (Base Card container).
    -   Create `src/components/common/AppButton.vue` (Base Button).
    -   Create `src/components/common/AppTabs.vue` (Segmented Tabs control).
3.  **Apply Central Layout and Components to Views**:
    -   **Home.vue**: Align layout with Image 2 (Word of the Day, structured list of devotions/readings, floating button).
    -   **Rosary.vue**: Align with Image 3 (2-column layout, left card with Luke ref, Annunciation title, blockquote, 10 bead tracker; right column with classical image and meditation quote card).
    -   **DailyMassReadings.vue**: Align with Image 4 (Ordinary Time heading, large date, switchable readings tabs, floating Meditate button).
    -   **DivineMercy.vue**: Clean up layout and apply warm parchment colors, custom buttons, and bottom nav.
    -   **Catechism.vue**: Apply warm styling to paragraph items and sidebar navigation, utilizing Libre Caslon Text and gold checkmarks.
    -   **PrayersGallery.vue**: Restyle list of cards and prayer detail modal.
    -   **Resources.vue**: Restyle documents reader with sidebar navigation and search bar.
4.  **Verification**:
    -   Run local development server.
    -   Test all pages on different screens for visual correctness and alignment.
    -   Verify production build output via `npm run build`.

---

## Architecture
-   **Framework**: Vue 3 (Composition API)
-   **Tooling**: Vite, TypeScript
-   **Data Storage**:
    -   `src/data/prayers.json`: Centralized prayer database.
    -   `src/components/rosary/rosaryData.ts`: Logic for generating full Rosary sequences.
    -   `src/components/divinemercy/divineMercyData.ts`: Sequence definitions for the Chaplet.

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
    -   `Catechism.vue`, `CatechismText.vue`, `CatechismReference.vue`: Logic for the interactive Catechism explorer.
-   **Resources**:
    -   `Resources.vue`: Interactive markdown resource reader.
