<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useDate } from '../composables/useDate';
import { useDevotionApi } from '../composables/useDevotionApi';
import { renderMarkdown } from '../utils/markdown';
import type { MagisteriumCitation } from '../components/common/FactCheckCitations.vue';

interface Reading {
  citation: string;
  text: string;
  context?: string;
  citations?: MagisteriumCitation[];
}

interface MassReadings {
  first_reading?: Reading;
  responsorial_psalm?: Reading;
  second_reading?: Reading;
  gospel?: Reading;
  pope_quote?: string;
}

const route = useRoute();
const router = useRouter();
const { getLocalISOString } = useDate();
const { getDevotion, loading, error } = useDevotionApi();

const initialDate = typeof route.query.date === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(route.query.date)
  ? route.query.date
  : getLocalISOString();

const selectedDate = ref(initialDate);
const dateInput = ref(initialDate);
const readings = ref<MassReadings | null>(null);

// Font size state (persisted in localStorage for Kindle convenience)
const FONT_SIZE_KEY = 'kindle_readings_font_size';
const fontSize = ref<number>(20);

const loadSavedFontSize = () => {
  try {
    const saved = localStorage.getItem(FONT_SIZE_KEY);
    if (saved) {
      const parsed = parseInt(saved, 10);
      if (!isNaN(parsed) && parsed >= 14 && parsed <= 34) {
        fontSize.value = parsed;
      }
    }
  } catch {
    // Ignore storage errors on restricted browsers
  }
};

const saveFontSize = (size: number) => {
  try {
    localStorage.setItem(FONT_SIZE_KEY, String(size));
  } catch {
    // Ignore storage errors
  }
};

const decreaseFont = () => {
  if (fontSize.value > 14) {
    fontSize.value -= 2;
    saveFontSize(fontSize.value);
  }
};

const increaseFont = () => {
  if (fontSize.value < 34) {
    fontSize.value += 2;
    saveFontSize(fontSize.value);
  }
};

const resetFont = () => {
  fontSize.value = 20;
  saveFontSize(20);
};

// Track which citation sections are expanded (expanded by default or toggleable without animations)
const expandedCitations = ref<Record<string, boolean>>({
  first_reading: false,
  responsorial_psalm: false,
  second_reading: false,
  gospel: false,
});

const toggleCitations = (key: string) => {
  expandedCitations.value[key] = !expandedCitations.value[key];
};

const cleanContext = (contextStr?: string) => {
  if (!contextStr) return '';
  return contextStr.replace(/<details[\s\S]*?<\/details>/gi, '').trim();
};

const getCitations = (reading?: Reading): MagisteriumCitation[] => {
  if (!reading) return [];
  if (reading.citations && reading.citations.length > 0) {
    return reading.citations;
  }
  if (reading.context && reading.context.includes('<details')) {
    const citations: MagisteriumCitation[] = [];
    const detailsMatch = reading.context.match(/<details[\s\S]*?<\/details>/i);
    if (detailsMatch) {
      const detailsContent = detailsMatch[0];
      const itemRegex = /\*\*\d+\.\s*(.*?)\*\*(?:\s*—\s*\*?(.*?)\*?)?(?:\s*\((.*?)\))?\s*\n(?:>\s*(.*?)(?=\n\n|\n\*\*|\n<\/details>|$))/gs;
      let match;
      while ((match = itemRegex.exec(detailsContent)) !== null) {
        citations.push({
          title: match[1]?.trim() || '',
          document_title: match[1]?.trim() || '',
          author: match[2]?.trim() || '',
          ref: match[3]?.trim() || '',
          text: match[4]?.trim() || '',
        });
      }
    }
    return citations;
  }
  return [];
};

const kindleMarkdownOptions = {
  paragraphClass: 'kindle-md-p',
  listClass: 'kindle-md-ul',
  listItemClass: 'kindle-md-li',
};

const fetchReadings = async () => {
  try {
    const data = await getDevotion<MassReadings>(selectedDate.value);
    readings.value = data;
  } catch (err) {
    console.error('Error fetching readings for Kindle view:', err);
  }
};

const applyDate = (newDate: string) => {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(newDate)) return;
  selectedDate.value = newDate;
  dateInput.value = newDate;
  router.replace({ query: { ...route.query, date: newDate } });
  fetchReadings();
};

const handleDateSubmit = () => {
  applyDate(dateInput.value.trim());
};

const shiftDate = (days: number) => {
  const d = new Date(`${selectedDate.value}T12:00:00`);
  if (isNaN(d.getTime())) return;
  d.setDate(d.getDate() + days);
  applyDate(getLocalISOString(d));
};

const goToToday = () => {
  applyDate(getLocalISOString());
};

const selectedDateValue = computed(() => new Date(`${selectedDate.value}T12:00:00`));

const formattedDate = computed(() => {
  if (isNaN(selectedDateValue.value.getTime())) return selectedDate.value;
  return selectedDateValue.value.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
});

const readingSections = computed(() => {
  if (!readings.value) return [];
  const sections: Array<{
    key: string;
    id: string;
    label: string;
    shortLabel: string;
    badge: string;
    meditationTitle: string;
    reading: Reading;
    isPsalm?: boolean;
  }> = [];

  if (readings.value.first_reading) {
    sections.push({
      key: 'first_reading',
      id: 'section-first-reading',
      label: 'First Reading',
      shortLabel: 'First Reading',
      badge: 'Sacred Scripture',
      meditationTitle: 'Contextual Meditation',
      reading: readings.value.first_reading,
    });
  }
  if (readings.value.responsorial_psalm) {
    sections.push({
      key: 'responsorial_psalm',
      id: 'section-psalm',
      label: 'Responsorial Psalm',
      shortLabel: 'Psalm',
      badge: 'Sacred Scripture',
      meditationTitle: 'Psalm Context',
      reading: readings.value.responsorial_psalm,
      isPsalm: true,
    });
  }
  if (readings.value.second_reading) {
    sections.push({
      key: 'second_reading',
      id: 'section-second-reading',
      label: 'Second Reading',
      shortLabel: 'Second Reading',
      badge: 'Sacred Scripture',
      meditationTitle: 'Contextual Meditation',
      reading: readings.value.second_reading,
    });
  }
  if (readings.value.gospel) {
    sections.push({
      key: 'gospel',
      id: 'section-gospel',
      label: 'Gospel Reading',
      shortLabel: 'Gospel',
      badge: 'Sacred Gospel',
      meditationTitle: 'Gospel Reflection',
      reading: readings.value.gospel,
    });
  }
  return sections;
});

watch(
  () => route.query.date,
  (newDate) => {
    if (typeof newDate === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(newDate) && newDate !== selectedDate.value) {
      selectedDate.value = newDate;
      dateInput.value = newDate;
      fetchReadings();
    }
  }
);

onMounted(() => {
  loadSavedFontSize();
  fetchReadings();
});
</script>

<template>
  <div id="kindle-top" class="kindle-page" :style="{ fontSize: fontSize + 'px' }">
    <!-- Top Utility Bar -->
    <header class="kindle-header">
      <div class="kindle-topbar">
        <div class="kindle-brand">
          <strong>☦ Faith &amp; Devotion</strong>
          <span class="kindle-tag">Kindle / E-Ink Edition</span>
        </div>
        <div class="kindle-font-controls" aria-label="Font size controls">
          <button type="button" class="kindle-btn" @click="decreaseFont" title="Decrease font size">A&minus;</button>
          <button type="button" class="kindle-btn" @click="resetFont" title="Reset font size">{{ fontSize }}px</button>
          <button type="button" class="kindle-btn" @click="increaseFont" title="Increase font size">A+</button>
          <router-link to="/daily-readings" class="kindle-btn kindle-link-btn">Standard View</router-link>
        </div>
      </div>

      <!-- Liturgical Date Title -->
      <div class="kindle-title-block">
        <div class="kindle-kicker">Daily Mass Readings &amp; Meditations</div>
        <h1 class="kindle-h1">{{ formattedDate }}</h1>
      </div>

      <!-- E-Ink Friendly Date Navigation -->
      <div class="kindle-date-bar">
        <div class="kindle-date-nav">
          <button type="button" class="kindle-btn" @click="shiftDate(-1)">&larr; Prev Day</button>
          <button type="button" class="kindle-btn" @click="goToToday">Today</button>
          <button type="button" class="kindle-btn" @click="shiftDate(1)">Next Day &rarr;</button>
        </div>
        <form class="kindle-date-form" @submit.prevent="handleDateSubmit">
          <label for="kindle-date-input" class="kindle-label">Date (YYYY-MM-DD):</label>
          <input
            id="kindle-date-input"
            v-model="dateInput"
            type="date"
            placeholder="YYYY-MM-DD"
            class="kindle-input"
          />
          <button type="submit" class="kindle-btn kindle-btn-primary">Go</button>
        </form>
      </div>

      <!-- Quick Jump Links -->
      <nav v-if="readings && !loading" class="kindle-jump-nav" aria-label="Jump to reading">
        <span class="kindle-jump-label">Jump to:</span>
        <a
          v-for="sec in readingSections"
          :key="sec.key"
          :href="'#' + sec.id"
          class="kindle-jump-link"
        >
          {{ sec.shortLabel }}
        </a>
        <a
          v-if="readings.pope_quote"
          href="#section-pope-quote"
          class="kindle-jump-link"
        >
          Words of the Popes
        </a>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="kindle-main">
      <!-- Loading State (No CSS Spinners for E-Ink) -->
      <div v-if="loading" class="kindle-status-box">
        <p><strong>Loading the liturgy for {{ selectedDate }}...</strong></p>
        <p class="kindle-small">Please wait while readings and meditations are retrieved.</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="kindle-status-box kindle-error-box">
        <p><strong>Liturgy Fetch Error</strong></p>
        <p>{{ error }}</p>
        <button type="button" class="kindle-btn kindle-btn-primary" @click="fetchReadings">
          Try Again
        </button>
      </div>

      <!-- Readings Content (Single Continuous Page) -->
      <div v-else-if="readings" class="kindle-content">
        <article
          v-for="sec in readingSections"
          :key="sec.key"
          :id="sec.id"
          class="kindle-section"
        >
          <div class="kindle-section-header">
            <div class="kindle-section-meta">
              <span class="kindle-section-type">{{ sec.label }}</span>
              <span class="kindle-badge">{{ sec.badge }}</span>
            </div>
            <div class="kindle-citation-row">
              <h2 class="kindle-h2">{{ sec.reading.citation }}</h2>
              <a href="#kindle-top" class="kindle-top-link">&uarr; Top</a>
            </div>
          </div>

          <!-- Scripture Passage -->
          <div
            class="kindle-scripture-text"
            :class="{ 'kindle-psalm-text': sec.isPsalm }"
            v-html="sec.reading.text"
          ></div>

          <!-- Contextual Meditation -->
          <div v-if="sec.reading.context" class="kindle-meditation-box">
            <div class="kindle-meditation-header">
              <h3 class="kindle-h3">{{ sec.meditationTitle }}</h3>
              <span class="kindle-badge">AI-Assisted Reflection</span>
            </div>
            <div
              class="kindle-meditation-body"
              v-html="renderMarkdown(cleanContext(sec.reading.context), kindleMarkdownOptions)"
            ></div>

            <!-- Fact-Checked Magisterium Citations -->
            <div v-if="getCitations(sec.reading).length > 0" class="kindle-citations-wrapper">
              <button
                type="button"
                class="kindle-btn kindle-citations-toggle"
                @click="toggleCitations(sec.key)"
              >
                <span>
                  Fact-Checked Magisterium Citations ({{ getCitations(sec.reading).length }})
                </span>
                <span>{{ expandedCitations[sec.key] ? '[Hide −]' : '[Show +]' }}</span>
              </button>

              <div v-if="expandedCitations[sec.key]" class="kindle-citations-list">
                <div
                  v-for="(item, idx) in getCitations(sec.reading)"
                  :key="idx"
                  class="kindle-citation-item"
                >
                  <div class="kindle-citation-item-header">
                    <strong>{{ idx + 1 }}. {{ item.document_title || item.title || 'Magisterial Document' }}</strong>
                    <span v-if="item.ref" class="kindle-ref-tag">[{{ item.ref }}]</span>
                  </div>
                  <div v-if="item.author" class="kindle-citation-author">
                    &mdash; {{ item.author }}
                  </div>
                  <blockquote
                    v-if="item.text"
                    class="kindle-citation-quote"
                    v-html="renderMarkdown(item.text, kindleMarkdownOptions)"
                  ></blockquote>
                  <div v-if="item.url" class="kindle-citation-link-wrap">
                    <a :href="item.url" target="_blank" rel="noopener noreferrer" class="kindle-link">
                      View Source Document &rarr;
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </article>

        <!-- Words of the Popes -->
        <article
          v-if="readings.pope_quote"
          id="section-pope-quote"
          class="kindle-section"
        >
          <div class="kindle-section-header">
            <div class="kindle-section-meta">
              <span class="kindle-section-type">Words of the Popes</span>
              <span class="kindle-badge">Papal Magisterium</span>
            </div>
            <div class="kindle-citation-row">
              <h2 class="kindle-h2">Holy Father's Reflection</h2>
              <a href="#kindle-top" class="kindle-top-link">&uarr; Top</a>
            </div>
          </div>

          <div class="kindle-scripture-text" v-html="readings.pope_quote"></div>

          <div class="kindle-source-note">
            Source: Vatican News Daily Liturgy
          </div>
        </article>

        <!-- Pastoral Callout -->
        <aside class="kindle-pastoral-box">
          <div class="kindle-pastoral-title">Pastoral Call to Community &amp; Prayer</div>
          <p class="kindle-pastoral-text">
            "Digital tools serve to support human reflection, but can never replace personal prayer, attendance at the Holy Sacrifice of the Mass, or the grace of the Sacraments in your local parish community."
          </p>
          <div class="kindle-pastoral-sub">
            Guided by Pope Leo XIV's Encyclical <em>Magnifica Humanitas</em>
          </div>
        </aside>

        <!-- Copyright -->
        <footer class="kindle-footer">
          <p>
            Holy Bible, New Living Translation Catholic Edition, copyright &copy; 2016 by Tyndale House Foundation. All rights reserved. Used by permission of Tyndale House Publishers.
          </p>
          <p class="kindle-footer-links">
            <a href="#kindle-top" class="kindle-link">&uarr; Back to Top</a>
            <span> &middot; </span>
            <a :href="'/kindle.html?date=' + selectedDate" class="kindle-link">Standalone HTML Mode</a>
          </p>
        </footer>
      </div>
    </main>
  </div>
</template>

<style scoped>
/* Strict E-Ink optimization: zero animations, zero transitions, solid black on white */
.kindle-page,
.kindle-page *,
.kindle-page *::before,
.kindle-page *::after {
  animation: none !important;
  transition: none !important;
  box-shadow: none !important;
  text-shadow: none !important;
  backdrop-filter: none !important;
}

.kindle-page {
  background-color: #ffffff;
  color: #000000;
  font-family: Bookerly, Georgia, 'Cardo', 'EB Garamond', 'Times New Roman', serif;
  line-height: 1.6;
  min-height: 100vh;
  max-width: 760px;
  margin: 0 auto;
  padding: 12px 16px 48px 16px;
}

.kindle-header {
  border-bottom: 2px solid #000000;
  padding-bottom: 14px;
  margin-bottom: 20px;
}

.kindle-topbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding-bottom: 10px;
  border-bottom: 1px solid #000000;
}

.kindle-brand {
  font-size: 0.85em;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.kindle-tag {
  font-family: sans-serif;
  font-size: 0.75em;
  border: 1px solid #000000;
  padding: 1px 6px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.kindle-font-controls {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.kindle-btn {
  background-color: #ffffff;
  color: #000000;
  border: 1.5px solid #000000;
  padding: 5px 10px;
  font-family: sans-serif;
  font-size: 13px;
  font-weight: bold;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  line-height: 1.2;
}

.kindle-btn-primary {
  background-color: #000000;
  color: #ffffff;
}

.kindle-title-block {
  padding: 14px 0 10px 0;
  text-align: center;
}

.kindle-kicker {
  font-family: sans-serif;
  font-size: 0.65em;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  margin-bottom: 4px;
}

.kindle-h1 {
  font-size: 1.5em;
  font-weight: bold;
  margin: 0;
  line-height: 1.25;
  color: #000000;
}

.kindle-date-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 10px 0;
  border-top: 1px solid #000000;
}

.kindle-date-nav {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.kindle-date-form {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.kindle-label {
  font-family: sans-serif;
  font-size: 12px;
  font-weight: bold;
}

.kindle-input {
  border: 1.5px solid #000000;
  background-color: #ffffff;
  color: #000000;
  padding: 4px 8px;
  font-family: sans-serif;
  font-size: 13px;
}

.kindle-jump-nav {
  border-top: 1px solid #000000;
  padding-top: 10px;
  margin-top: 4px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  font-family: sans-serif;
  font-size: 13px;
}

.kindle-jump-label {
  font-weight: bold;
  text-transform: uppercase;
  font-size: 11px;
  letter-spacing: 0.05em;
}

.kindle-jump-link {
  color: #000000;
  text-decoration: underline;
  font-weight: bold;
  padding: 3px 8px;
  border: 1px solid #000000;
  background-color: #f5f5f5;
}

.kindle-status-box {
  border: 2px solid #000000;
  padding: 20px;
  text-align: center;
  margin: 20px 0;
}

.kindle-small {
  font-size: 0.85em;
}

.kindle-section {
  border-bottom: 2px solid #000000;
  padding-bottom: 24px;
  margin-bottom: 28px;
}

.kindle-section-header {
  border-bottom: 1px solid #000000;
  padding-bottom: 8px;
  margin-bottom: 16px;
}

.kindle-section-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 4px;
}

.kindle-section-type {
  font-family: sans-serif;
  font-size: 0.65em;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.kindle-badge {
  font-family: sans-serif;
  font-size: 11px;
  font-weight: bold;
  border: 1px solid #000000;
  padding: 1px 6px;
  background-color: #f5f5f5;
}

.kindle-citation-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
}

.kindle-h2 {
  font-size: 1.25em;
  font-weight: bold;
  margin: 0;
  color: #000000;
}

.kindle-top-link {
  font-family: sans-serif;
  font-size: 12px;
  font-weight: bold;
  color: #000000;
  text-decoration: underline;
  white-space: nowrap;
}

.kindle-scripture-text {
  white-space: pre-line;
  font-size: 1em;
  line-height: 1.65;
  margin-bottom: 18px;
}

.kindle-psalm-text {
  line-height: 1.75;
}

.kindle-meditation-box {
  border: 1.5px solid #000000;
  border-left-width: 5px;
  padding: 14px 16px;
  margin-top: 16px;
  background-color: #fafafa;
}

.kindle-meditation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  border-bottom: 1px solid #000000;
  padding-bottom: 6px;
  margin-bottom: 10px;
}

.kindle-h3 {
  font-family: sans-serif;
  font-size: 0.75em;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin: 0;
}

.kindle-meditation-body {
  font-size: 0.92em;
  line-height: 1.55;
}

.kindle-meditation-body :deep(.kindle-md-p) {
  margin: 0 0 10px 0;
}

.kindle-meditation-body :deep(.kindle-md-p:last-child) {
  margin-bottom: 0;
}

.kindle-meditation-body :deep(.kindle-md-ul) {
  margin: 8px 0;
  padding-left: 20px;
  list-style-type: disc;
}

.kindle-meditation-body :deep(.kindle-md-li) {
  margin-bottom: 6px;
}

.kindle-citations-wrapper {
  margin-top: 14px;
  padding-top: 10px;
  border-top: 1px solid #000000;
}

.kindle-citations-toggle {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  text-align: left;
  padding: 8px 10px;
  background-color: #ffffff;
}

.kindle-citations-list {
  margin-top: 10px;
}

.kindle-citation-item {
  border: 1px solid #000000;
  background-color: #ffffff;
  padding: 10px 12px;
  margin-bottom: 10px;
  font-size: 0.85em;
}

.kindle-citation-item:last-child {
  margin-bottom: 0;
}

.kindle-citation-item-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 4px;
}

.kindle-ref-tag {
  font-family: monospace;
  font-size: 0.9em;
  font-weight: bold;
  white-space: nowrap;
}

.kindle-citation-author {
  font-style: italic;
  margin-bottom: 6px;
}

.kindle-citation-quote {
  margin: 6px 0;
  padding-left: 10px;
  border-left: 3px solid #000000;
}

.kindle-citation-link-wrap {
  margin-top: 6px;
}

.kindle-link {
  color: #000000;
  font-weight: bold;
  text-decoration: underline;
}

.kindle-source-note {
  font-family: sans-serif;
  font-size: 0.7em;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  border-top: 1px solid #000000;
  padding-top: 8px;
  margin-top: 12px;
}

.kindle-pastoral-box {
  border: 1.5px solid #000000;
  padding: 14px 16px;
  text-align: center;
  margin-top: 24px;
}

.kindle-pastoral-title {
  font-family: sans-serif;
  font-size: 0.7em;
  font-weight: bold;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin-bottom: 6px;
}

.kindle-pastoral-text {
  font-size: 0.88em;
  margin: 0 0 6px 0;
}

.kindle-pastoral-sub {
  font-family: sans-serif;
  font-size: 0.65em;
  font-weight: bold;
  text-transform: uppercase;
}

.kindle-footer {
  margin-top: 20px;
  text-align: center;
  font-family: sans-serif;
  font-size: 12px;
  line-height: 1.4;
}

.kindle-footer-links {
  margin-top: 8px;
}
</style>
