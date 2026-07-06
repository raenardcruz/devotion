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
    -   **Intro Bead Navigation**: Fixed bead tapping for the initial 3 Hail Marys to correctly jump to the respective intro step.
    -   **Audio Overlap Prevention**: Cleaned up pending auto-play timers on navigation and page change to avoid voice doubling/overlap.
    -   **Audio Pause & Resume**: Refactored the audio play/pause logic so that pausing retains the loaded audio instance, allowing subsequent plays to resume from the last position.
    -   **Dynamic Slide Titles**: Integrated active mystery name dynamically inside decade slide headers (e.g. "First Decade - The Annunciation - Our Father" and "First Decade - The Annunciation - Bead 1").
    -   **Sacred Orthodox & Catholic Icons**: Generated and stored 20 authentic public domain style icons locally under `/public/images/rosary/` to replace the generic Unsplash landscape images and eliminate external load issues. Updated both `rosaryData.ts` and `mysteryVisuals.ts` to reference these local assets. Prioritized the local `image` field in `Rosary.vue`.

2.  **Divine Mercy Chaplet**:
    -   **Short Version Toggle**: Added a toggle to skip to closing prayers after the first decade, ideal for quick devotion.
    -   **Centralized Data**: Migrated all prayer text to a single JSON source of truth.
    -   **Refined Navigation**: Enhanced `next`/`prev` logic to handle complex transitions between decades and closing sequences.
    -   **UI Indicators**: Added "Short Version" badges and persistent decade progress tracking.
    -   **Bilingual Display**: Automatically displays English reference text below Latin prayers when the Latin toggle is active.
    -   **Audio Navigation Stability**: Cleared playback timers on navigation changes to prevent double voice issues.
    -   **Audio Playback Resuming**: Integrated paused-audio retention to resume the chaplet's audio from where the user paused.

3.  **Expanded Prayer Library & Search**:
    -   Added Act of Contrition (English & Latin) to the centralized data store.
    -   Verified and corrected Latin character diacritics, accents, and ligatures of existing prayers based on Adoremus.
    -   Added 6 missing traditional prayers: Eternal Rest, The Angelus, Regina Caeli, Blessing Before Meals, Grace After Meals, and Nicene Creed.
    -   Implemented real-time, multilingual search functionality for the Prayers Gallery, supporting checks across both English and Latin titles and text content.
    -   Designed an elegant Empty State layout when a search returns no matching results, offering a prompt reset button.
    -   Integrated audio recordings for Eternal Rest (`eternal-rest.wav`), Nicene Creed (`nicene-creed.wav`), and Regina Caeli (`the-regina-caeli.wav`) by registering them in `AVAILABLE_AUDIO_PRAYERS`.

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
    -   Extracted the local markdown parsing function to a reusable shared utility `src/utils/markdown.ts`.
    -   Integrated Markdown parsing for the contextual meditation overlays in the Daily Mass Readings view (`DailyMassReadings.vue`), rendering headers, bold, italics, and lists beautifully.
9.  **Interactive Bible Citations & Latin Titles**:
    -   Relocated Bible chapter/verse citations in the Scriptural Rosary to directly underneath blockquote quotes as interactive `<router-link>` elements.
    -   Converted Daily Mass Readings citation headers into interactive `<router-link>` elements that navigate to the Bible Study explorer.
    -   Integrated query parameter detection in `BibleStudy.vue` to automatically trigger lookup of linked passages.
    -   Added comprehensive Latin translations for prayer, mystery, set, category, and step heading titles, rendering them dynamically in Latin mode.
    -   Enhanced the Prayers Gallery list with dual English/Latin titles for a richer, scholarly aesthetic.
10. **Rosary Completion Marker & Divine Mercy Audio Fix**:
    -   Implemented a dedicated 81st completion step in the Scriptural Rosary devotion sequence.
    -   Created custom completion slide templates inside `VerseCard.vue` with a golden-pulsing cross icon, and added interactive buttons to restart the Rosary or return home.
    -   Updated the `RosaryControls` Next button to display "Done" with a checkmark on the last step.
    -   Resolved the missing audio tracks for beads 2-10 in the Divine Mercy Chaplet by watching `[currentStepIndex, beadInDecade]` changes instead of step object reference identities.
11. **Back to Top Button**:
    -   Created a global `BackToTop.vue` component under `src/components/common/`.
    -   Integrated it in `App.vue` to dynamically show up across all scrollable views.
    -   Styled it as an elegant gold pill button (`bg-parchment-primary`) positioned at the bottom center (`left-1/2 -translate-x-1/2`) just above the navigation bar.
    -   Added a micro-animation bounce effect on the up arrow to enhance interactive aesthetics.

### Early Church Writings Redesign
- **Dynamic Asynchronous Loading**: Replaced static markdown imports with dynamic glob imports via Vite. Only the chosen markdown file is fetched and processed, preventing memory/network overflow for large files (e.g. Augustine's 21MB works).
- **Resource Indexer**: Created a script `research/generate_index.py` that extracts the main title header from all 71 Church Father `.md` files and outputs them into a static `index.json` registry.
- **Searchable Writings Selector**: Redesigned the top layout of `Resources.vue` with a premium dropdown button and searchable popover selector to navigate all 71 writings with real-time text filtering.
- **State Persistence**: Added automatic restoration of both the active writing and active chapter on page reload via `localStorage`.
- **Book-Like Page Layout & Typesetting**: Formatted the reader layout to mimic physical books, featuring centered chapter headings (numeral and title) separated by a thin gold divider, a maximum text reading width (`max-w-2xl`), justified alignment, first-line indenting (except for the first paragraph), and blockquotes rendered as centered, borderless epigraphs.
- **Inline Markdown Parser**: Added safe, lightweight inline markdown styling to convert bold (`**`, `__`), italic (`*`, `_`), and double dashes into proper HTML tags/entities.

---

## Active Plan: Add Sacred Icons to Rosary Mysteries
- **Status**: Completed.
- **Goal**: Integrate authentic Catholic and Orthodox icon/painting URLs into the Rosary's 20 mysteries to replace generic Unsplash images.



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
