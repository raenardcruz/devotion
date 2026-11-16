<script setup lang="ts">
import { computed } from 'vue';
import { renderMarkdown } from '../../utils/markdown';
import {
  groupPanelsIntoComicRows,
  getTileAspectClass,
  getTileColSpanClass,
} from '../../utils/mangaLayout';
import type {
  ComicStyle,
  MangaGeneratedPanel,
  MangaPanelPlan,
  MangaStoryAnalysis,
  MangaStoryBreakdown,
} from '../../composables/useStoryToManga';
import type {
  MagazineReadingImage,
  MagazineReadingKey,
} from '../../composables/useMagazineImages';
import type { MagisteriumCitation } from '../common/FactCheckCitations.vue';

export interface MagazineReadingItem {
  citation: string;
  text: string;
  context?: string;
  citations?: MagisteriumCitation[];
}

export interface MagazineMassReadings {
  first_reading?: MagazineReadingItem;
  responsorial_psalm?: MagazineReadingItem;
  second_reading?: MagazineReadingItem;
  gospel?: MagazineReadingItem;
  pope_quote?: string;
}

const props = withDefaults(
  defineProps<{
    date: string;
    formattedDate: string;
    readings: MagazineMassReadings;
    images: Partial<Record<MagazineReadingKey, MagazineReadingImage>>;
    resolveImageSrc: (img?: MagazineReadingImage) => string;
    generatingKeys?: Set<MagazineReadingKey>;
    interactive?: boolean;
    includeCommentaries?: boolean;
    includeImages?: boolean;
    includeComics?: boolean;
    // Comics / Story-to-Manga data (rendered at the end)
    mangaStyle?: ComicStyle;
    storyAnalysis?: MangaStoryAnalysis | null;
    moralLesson?: string;
    storyBreakdown?: MangaStoryBreakdown | null;
    generatedPanels?: MangaGeneratedPanel[];
    getProxyAssetUrl?: (r2Key?: string) => string;
  }>(),
  {
    interactive: false,
    includeCommentaries: true,
    includeImages: true,
    includeComics: true,
    mangaStyle: 'manga',
    storyAnalysis: null,
    moralLesson: '',
    storyBreakdown: null,
    generatedPanels: () => [],
  }
);

const emit = defineEmits<{
  (
    e: 'generate-image',
    payload: {
      key: MagazineReadingKey;
      citation: string;
      text?: string;
      context?: string;
      force?: boolean;
    }
  ): void;
}>();

const cleanContext = (contextStr?: string): string => {
  if (!contextStr) return '';
  return contextStr.replace(/<details[\s\S]*?<\/details>/gi, '').trim();
};

const stripHtml = (html?: string): string => {
  if (!html) return '';
  return html
    .replace(/<[^>]+>/g, ' ')
    .replace(/\s+/g, ' ')
    .trim();
};

const extractPullQuote = (reading?: MagazineReadingItem): string => {
  if (!reading) return '';
  const source = cleanContext(reading.context) || stripHtml(reading.text);
  if (!source) return '';
  const plain = source
    .replace(/[#*_`>]/g, '')
    .replace(/\s+/g, ' ')
    .trim();
  const sentences = plain.match(/[^.!?]+[.!?]+/g) || [];
  const candidate = sentences.find((s) => s.trim().length >= 55 && s.trim().length <= 185);
  return candidate ? candidate.trim() : '';
};

const parseCitations = (reading?: MagazineReadingItem): MagisteriumCitation[] => {
  if (!reading) return [];
  if (reading.citations && reading.citations.length > 0) {
    return reading.citations.slice(0, 3);
  }
  if (reading.context && reading.context.includes('<details')) {
    const out: MagisteriumCitation[] = [];
    const detailsMatch = reading.context.match(/<details[\s\S]*?<\/details>/i);
    if (detailsMatch) {
      const itemRegex =
        /\*\*\d+\.\s*(.*?)\*\*(?:\s*—\s*\*?(.*?)\*?)?(?:\s*\((.*?)\))?\s*\n(?:>\s*(.*?)(?=\n\n|\n\*\*|\n<\/details>|$))/gs;
      let match;
      while ((match = itemRegex.exec(detailsMatch[0])) !== null) {
        out.push({
          title: match[1]?.trim() || '',
          document_title: match[1]?.trim() || '',
          author: match[2]?.trim() || '',
          ref: match[3]?.trim() || '',
          text: match[4]?.trim() || '',
        });
      }
    }
    return out.slice(0, 3);
  }
  return [];
};

const commentaryMarkdownOptions = {
  paragraphClass: 'mag-commentary-p',
  listClass: 'mag-commentary-ul',
  listItemClass: 'mag-commentary-li',
};

const issueSubtitle = computed(() => {
  const d = new Date(`${props.date}T12:00:00`);
  if (Number.isNaN(d.getTime())) return 'Liturgical Daily Edition';
  const weekday = d.toLocaleDateString('en-US', { weekday: 'long' });
  return `${weekday} of the Liturgical Year · Daily Liturgical Edition`;
});

const readingSections = computed(() => {
  const list: Array<{
    key: Exclude<MagazineReadingKey, 'cover'>;
    kicker: string;
    label: string;
    isPsalm?: boolean;
    reading: MagazineReadingItem;
  }> = [];

  if (props.readings.first_reading) {
    list.push({
      key: 'first_reading',
      kicker: 'Part I · First Reading',
      label: 'Sacred Scripture',
      reading: props.readings.first_reading,
    });
  }
  if (props.readings.responsorial_psalm) {
    list.push({
      key: 'responsorial_psalm',
      kicker: 'Responsorial Psalm',
      label: 'Liturgical Song',
      isPsalm: true,
      reading: props.readings.responsorial_psalm,
    });
  }
  if (props.readings.second_reading) {
    list.push({
      key: 'second_reading',
      kicker: 'Second Reading',
      label: 'Apostolic Epistle',
      reading: props.readings.second_reading,
    });
  }
  if (props.readings.gospel) {
    list.push({
      key: 'gospel',
      kicker: 'The Holy Gospel',
      label: 'Gospel Proclamation',
      reading: props.readings.gospel,
    });
  }
  return list;
});

const coverImageSrc = computed(() => props.resolveImageSrc(props.images.cover));

const hasComics = computed(
  () =>
    props.includeComics &&
    Boolean(props.storyAnalysis) &&
    (props.generatedPanels?.length || 0) > 0
);

const comicLayoutRows = computed(() =>
  groupPanelsIntoComicRows(props.storyBreakdown?.panels || [])
);

const getGeneratedPanelImage = (panelNumber: number): string => {
  const found = props.generatedPanels?.find((p) => p.panelNumber === panelNumber);
  if (!found) return '';
  if (found.image && (found.image.startsWith('data:') || found.image.startsWith('http'))) {
    return found.image;
  }
  if (found.r2Key && props.getProxyAssetUrl) {
    return props.getProxyAssetUrl(found.r2Key);
  }
  return found.image || '';
};

const getPanelNarrative = (panel: MangaPanelPlan): string => {
  return (panel.narrative || panel.sceneDescription || '').trim();
};

const isGeneratingKey = (key: MagazineReadingKey): boolean => {
  return Boolean(props.generatingKeys?.has(key));
};
</script>

<template>
  <article class="magazine-root">
    <!-- ==================== PAGE 1: MAGAZINE COVER ==================== -->
    <section class="magazine-sheet magazine-cover">
      <header class="cover-masthead">
        <div class="cover-issue-bar">
          <span>VOL. MMXXVI · ISSUE {{ date }}</span>
          <span>CATHOLIC DAILY LITURGICAL MAGAZINE</span>
        </div>
        <h1 class="masthead-title">LUMEN LITURGIAE</h1>
        <p class="masthead-subtitle">{{ issueSubtitle }}</p>
      </header>

      <!-- Cover Artwork -->
      <figure v-if="includeImages" class="cover-figure">
        <div v-if="coverImageSrc" class="cover-image-frame">
          <img :src="coverImageSrc" alt="Cover Sacred Art Illustration" class="cover-image" />
        </div>
        <div v-else class="cover-placeholder">
          <div class="cover-emblem">☩</div>
          <p class="cover-placeholder-title">{{ formattedDate }}</p>
          <p class="cover-placeholder-sub">
            {{ readings.gospel?.citation || 'Daily Mass Readings & Reflections' }}
          </p>
          <button
            v-if="interactive && readings.gospel"
            type="button"
            :disabled="isGeneratingKey('cover')"
            @click="
              emit('generate-image', {
                key: 'cover',
                citation: readings.gospel.citation,
                text: readings.gospel.text,
                context: readings.gospel.context,
              })
            "
            class="mag-gen-btn print:hidden"
          >
            {{ isGeneratingKey('cover') ? 'Generating Cover Art...' : '✨ Generate Cover Illustration' }}
          </button>
        </div>

        <figcaption v-if="coverImageSrc" class="cover-caption-row">
          <span class="mag-caption">
            Cover Illumination · Inspired by {{ readings.gospel?.citation || formattedDate }}
          </span>
          <button
            v-if="interactive && readings.gospel"
            type="button"
            :disabled="isGeneratingKey('cover')"
            @click="
              emit('generate-image', {
                key: 'cover',
                citation: readings.gospel.citation,
                text: readings.gospel.text,
                context: readings.gospel.context,
                force: true,
              })
            "
            class="mag-regen-link print:hidden"
          >
            {{ isGeneratingKey('cover') ? 'Regenerating...' : '↻ Regenerate' }}
          </button>
        </figcaption>
      </figure>

      <!-- Cover Lead & Table of Contents -->
      <div class="cover-lead-grid">
        <div class="cover-headline-block">
          <span class="mag-kicker">Featured Gospel Proclamation</span>
          <h2 class="cover-headline">
            {{ readings.gospel?.citation || formattedDate }}
          </h2>
          <p v-if="extractPullQuote(readings.gospel)" class="cover-deck">
            “{{ extractPullQuote(readings.gospel) }}”
          </p>
        </div>

        <aside class="cover-toc">
          <h3 class="toc-title">In This Issue</h3>
          <ul class="toc-list">
            <li v-for="(sec, idx) in readingSections" :key="sec.key">
              <span class="toc-num">0{{ idx + 1 }}</span>
              <div>
                <span class="toc-kicker">{{ sec.kicker }}</span>
                <span class="toc-citation">{{ sec.reading.citation }}</span>
              </div>
            </li>
            <li v-if="readings.pope_quote">
              <span class="toc-num">0{{ readingSections.length + 1 }}</span>
              <div>
                <span class="toc-kicker">Papal Magisterium</span>
                <span class="toc-citation">Words of the Popes</span>
              </div>
            </li>
            <li v-if="hasComics && storyAnalysis">
              <span class="toc-num">★</span>
              <div>
                <span class="toc-kicker">Closing Parable &amp; Comics</span>
                <span class="toc-citation">{{ storyAnalysis.title }}</span>
              </div>
            </li>
          </ul>
        </aside>
      </div>

      <footer class="cover-footer">
        <span>{{ formattedDate }}</span>
        <span>Sacred Scripture · AI Theological Commentary · Illustrated Parable</span>
      </footer>
    </section>

    <!-- ==================== INTERIOR: READINGS & AI COMMENTARIES ==================== -->
    <section
      v-for="sec in readingSections"
      :key="sec.key"
      class="magazine-sheet magazine-article"
    >
      <!-- Running Header -->
      <div class="running-header">
        <span>LUMEN LITURGIAE · {{ formattedDate.toUpperCase() }}</span>
        <span>{{ sec.kicker.toUpperCase() }}</span>
      </div>

      <!-- Article Header -->
      <header class="article-header">
        <div class="article-kicker-row">
          <span class="mag-kicker">{{ sec.kicker }}</span>
          <span class="mag-badge">{{ sec.label }}</span>
        </div>
        <h2 class="article-title">{{ sec.reading.citation }}</h2>
      </header>

      <!-- Reading Editorial Illustration -->
      <figure v-if="includeImages" class="article-figure">
        <template v-if="resolveImageSrc(images[sec.key])">
          <div class="article-image-frame">
            <img
              :src="resolveImageSrc(images[sec.key])"
              :alt="`Illustration for ${sec.reading.citation}`"
              class="article-image"
            />
          </div>
          <figcaption class="cover-caption-row">
            <span class="mag-caption">
              Editorial Sacred Art · {{ sec.reading.citation }}
            </span>
            <button
              v-if="interactive"
              type="button"
              :disabled="isGeneratingKey(sec.key)"
              @click="
                emit('generate-image', {
                  key: sec.key,
                  citation: sec.reading.citation,
                  text: sec.reading.text,
                  context: sec.reading.context,
                  force: true,
                })
              "
              class="mag-regen-link print:hidden"
            >
              {{ isGeneratingKey(sec.key) ? 'Regenerating...' : '↻ Regenerate' }}
            </button>
          </figcaption>
        </template>

        <div v-else-if="interactive" class="article-image-empty print:hidden">
          <span class="text-xs text-[#5C5149]">
            No illustration generated yet for <strong>{{ sec.reading.citation }}</strong>
          </span>
          <button
            type="button"
            :disabled="isGeneratingKey(sec.key)"
            @click="
              emit('generate-image', {
                key: sec.key,
                citation: sec.reading.citation,
                text: sec.reading.text,
                context: sec.reading.context,
              })
            "
            class="mag-gen-btn-sm"
          >
            {{ isGeneratingKey(sec.key) ? 'Generating...' : '✨ Generate Image' }}
          </button>
        </div>
      </figure>

      <!-- Scripture Body -->
      <div
        v-if="sec.isPsalm"
        class="psalm-body"
        v-html="sec.reading.text"
      ></div>
      <div
        v-else
        class="magazine-columns magazine-dropcap scripture-body"
        v-html="sec.reading.text"
      ></div>

      <!-- Pull Quote -->
      <blockquote v-if="extractPullQuote(sec.reading)" class="magazine-pullquote">
        “{{ extractPullQuote(sec.reading) }}”
      </blockquote>

      <!-- AI Commentary Box -->
      <aside
        v-if="includeCommentaries && cleanContext(sec.reading.context)"
        class="ai-commentary-box"
      >
        <div class="ai-commentary-header">
          <span class="ai-commentary-kicker">Contextual Meditation</span>
          <span class="ai-commentary-tag">AI-Assisted Theological Commentary</span>
        </div>
        <div
          class="ai-commentary-body"
          v-html="renderMarkdown(cleanContext(sec.reading.context), commentaryMarkdownOptions)"
        ></div>

        <!-- Magisterium Citations -->
        <div
          v-if="parseCitations(sec.reading).length > 0"
          class="magisterium-sources"
        >
          <p class="magisterium-sources-title">Magisterium &amp; Scriptural Sources</p>
          <div
            v-for="(cit, cIdx) in parseCitations(sec.reading)"
            :key="cIdx"
            class="magisterium-source-item"
          >
            <span class="font-bold text-[#72383D]">{{ cIdx + 1 }}. {{ cit.document_title || cit.title }}</span>
            <span v-if="cit.author || cit.ref" class="text-[#5C5149]">
              — {{ [cit.author, cit.ref].filter(Boolean).join(', ') }}
            </span>
          </div>
        </div>
      </aside>
    </section>

    <!-- ==================== WORDS OF THE POPES ==================== -->
    <section v-if="readings.pope_quote" class="magazine-sheet magazine-article">
      <div class="running-header">
        <span>LUMEN LITURGIAE · {{ formattedDate.toUpperCase() }}</span>
        <span>PAPAL MAGISTERIUM</span>
      </div>
      <header class="article-header">
        <div class="article-kicker-row">
          <span class="mag-kicker">Words of the Popes</span>
          <span class="mag-badge">Holy Father's Reflection</span>
        </div>
        <h2 class="article-title">Papal Reflection on Today's Liturgy</h2>
      </header>
      <div class="magazine-columns magazine-dropcap scripture-body" v-html="readings.pope_quote"></div>
      <p class="mag-caption mt-4">Source: Vatican News Daily Liturgy</p>
    </section>

    <!-- ==================== AT THE END: DAILY PARABLE & COMICS ==================== -->
    <section v-if="hasComics && storyAnalysis" class="magazine-sheet magazine-comics-section">
      <div class="running-header">
        <span>LUMEN LITURGIAE · {{ formattedDate.toUpperCase() }}</span>
        <span>PART II · ILLUSTRATED PARABLE &amp; COMICS</span>
      </div>

      <header class="comics-masthead">
        <span class="mag-kicker">Closing Feature · Illustrated Parable</span>
        <h2 class="article-title">{{ storyAnalysis.title }}</h2>
        <div class="comics-meta-row">
          <span class="comics-meta-pill">
            📍 {{ storyAnalysis.setting.location }} · {{ storyAnalysis.setting.mood }}
          </span>
          <span class="comics-meta-pill">
            {{ mangaStyle === 'manga' ? 'Manga Edition' : 'Comic Edition' }}
          </span>
        </div>
        <div v-if="moralLesson" class="comics-lesson-box">
          <span class="comics-lesson-kicker">Lesson Applied from Today's Readings</span>
          <p class="comics-lesson-text">{{ moralLesson }}</p>
        </div>
      </header>

      <div class="comics-rows-stack">
        <div
          v-for="row in comicLayoutRows"
          :key="row.rowNumber"
          class="comic-layout-row grid grid-cols-1 md:grid-cols-6 gap-4 items-stretch"
        >
          <div
            v-for="tile in row.tiles"
            :key="tile.panel.panelNumber"
            class="comic-tile"
            :class="getTileColSpanClass(tile.colSpan)"
          >
            <div class="comic-tile-bar">
              <span>PANEL {{ tile.panel.panelNumber }}</span>
            </div>

            <div
              class="comic-tile-art"
              :class="getTileAspectClass(tile.aspectRatio, tile.colSpan)"
            >
              <img
                v-if="getGeneratedPanelImage(tile.panel.panelNumber)"
                :src="getGeneratedPanelImage(tile.panel.panelNumber)"
                :alt="`Panel ${tile.panel.panelNumber}`"
                class="w-full h-full object-cover"
              />
            </div>

            <div class="comic-tile-caption">
              <p v-if="getPanelNarrative(tile.panel)" class="comic-narrative">
                {{ getPanelNarrative(tile.panel) }}
              </p>
              <div v-if="tile.panel.dialogue" class="comic-dialogue-bubble">
                “{{ tile.panel.dialogue }}”
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== COLOPHON / BACK FOOTER ==================== -->
    <footer class="magazine-colophon">
      <div class="colophon-rule"></div>
      <p class="colophon-pastoral">
        “Digital tools serve to support human reflection, but can never replace personal prayer, attendance at the Holy Sacrifice of the Mass, or the grace of the Sacraments in your local parish community.”
      </p>
      <p class="colophon-sub">
        Lumen Liturgiae · Daily Catholic Readings Magazine · Guided by Pope Leo XIV's Encyclical <em>Magnifica Humanitas</em>
      </p>
      <p class="colophon-copyright">
        Holy Bible, New Living Translation Catholic Edition, copyright © 2016 by Tyndale House Foundation. All rights reserved. Used by permission of Tyndale House Publishers.
      </p>
    </footer>
  </article>
</template>

<style scoped>
/* Base Editorial Magazine Design Tokens */
.magazine-root {
  max-width: 210mm;
  margin: 0 auto;
  color: #1e1b18;
  font-family: 'EB Garamond', 'Cardo', 'Libre Caslon Text', Georgia, serif;
}

.magazine-sheet {
  background: #ffffff;
  border: 1px solid #e4ddd3;
  border-radius: 4px;
  padding: 16mm 16mm 14mm 16mm;
  margin-bottom: 10mm;
  box-shadow: 0 8px 28px rgba(30, 27, 24, 0.08);
}

/* Cover Masthead */
.cover-masthead {
  text-align: center;
  border-bottom: 2px solid #1e1b18;
  padding-bottom: 5mm;
  margin-bottom: 6mm;
}

.cover-issue-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-family: 'Source Sans 3', 'Inter', sans-serif;
  font-size: 8pt;
  font-weight: 700;
  letter-spacing: 0.2em;
  color: #72383d;
  border-bottom: 1px solid #d1c7bd;
  padding-bottom: 2.5mm;
  margin-bottom: 4mm;
}

.masthead-title {
  font-family: 'Playfair Display', 'Cinzel', Georgia, serif;
  font-size: 34pt;
  font-weight: 900;
  letter-spacing: -0.02em;
  line-height: 1.02;
  color: #1e1b18;
  margin: 0;
}

.masthead-subtitle {
  font-family: 'Source Sans 3', 'Inter', sans-serif;
  font-size: 9pt;
  font-weight: 600;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: #5c5149;
  margin-top: 2mm;
}

/* Cover Artwork */
.cover-figure {
  margin: 0 0 6mm 0;
}

.cover-image-frame {
  border: 1px solid #d1c7bd;
  padding: 2mm;
  background: #fdfbf7;
}

.cover-image {
  width: 100%;
  aspect-ratio: 4 / 3;
  object-fit: cover;
  display: block;
}

.cover-placeholder {
  aspect-ratio: 16 / 9;
  background: linear-gradient(135deg, #faf6ee 0%, #f1e8da 100%);
  border: 1px solid #c5a059;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 8mm;
}

.cover-emblem {
  font-size: 28pt;
  color: #72383d;
  margin-bottom: 2mm;
}

.cover-placeholder-title {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 18pt;
  font-weight: 700;
  color: #1e1b18;
}

.cover-placeholder-sub {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 9.5pt;
  color: #5c5149;
  margin-top: 1mm;
}

.cover-caption-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-top: 2mm;
}

.mag-caption {
  font-family: 'Source Sans 3', 'Inter', sans-serif;
  font-size: 8.5pt;
  font-style: italic;
  color: #5c5149;
}

/* Cover Lead & TOC */
.cover-lead-grid {
  display: grid;
  grid-template-columns: 1.35fr 1fr;
  gap: 7mm;
  border-top: 1px solid #d1c7bd;
  padding-top: 5mm;
}

.cover-headline {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 21pt;
  font-weight: 700;
  line-height: 1.15;
  letter-spacing: -0.015em;
  color: #72383d;
  margin: 1.5mm 0 2.5mm 0;
}

.cover-deck {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 12pt;
  font-style: italic;
  line-height: 1.48;
  color: #322d29;
}

.cover-toc {
  background: #faf6ee;
  border-left: 3px solid #c5a059;
  padding: 4mm 4.5mm;
}

.toc-title {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8.5pt;
  font-weight: 700;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: #72383d;
  margin: 0 0 3mm 0;
  border-bottom: 1px solid #e4ddd3;
  padding-bottom: 1.5mm;
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 2.5mm;
}

.toc-list li {
  display: flex;
  align-items: flex-start;
  gap: 2.5mm;
}

.toc-num {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 10pt;
  font-weight: 700;
  color: #c5a059;
  min-width: 16px;
}

.toc-kicker {
  display: block;
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  color: #5c5149;
}

.toc-citation {
  display: block;
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 10.5pt;
  font-weight: 600;
  color: #1e1b18;
}

.cover-footer {
  margin-top: 6mm;
  padding-top: 3mm;
  border-top: 1px solid #d1c7bd;
  display: flex;
  justify-content: space-between;
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8pt;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #5c5149;
}

/* Interior Running Header & Article Headers */
.running-header {
  display: flex;
  justify-content: space-between;
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 700;
  letter-spacing: 0.18em;
  color: #72383d;
  border-bottom: 1px solid #d1c7bd;
  padding-bottom: 2mm;
  margin-bottom: 5mm;
}

.article-header {
  margin-bottom: 4.5mm;
}

.article-kicker-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.mag-kicker {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8.5pt;
  font-weight: 700;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: #72383d;
}

.mag-badge {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  padding: 2px 8px;
  border-radius: 999px;
  background: #faf6ee;
  border: 1px solid #d1c7bd;
  color: #5c5149;
}

.article-title {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 22pt;
  font-weight: 700;
  line-height: 1.15;
  letter-spacing: -0.015em;
  color: #1e1b18;
  margin: 1.5mm 0 0 0;
}

/* Article Figure */
.article-figure {
  margin: 0 0 5mm 0;
  break-inside: avoid;
}

.article-image-frame {
  border: 1px solid #d1c7bd;
  padding: 1.5mm;
  background: #fdfbf7;
}

.article-image {
  width: 100%;
  aspect-ratio: 16 / 9;
  object-fit: cover;
  display: block;
}

.article-image-empty {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 3mm 4mm;
  background: #faf6ee;
  border: 1px dashed #c5a059;
  border-radius: 6px;
}

/* Multi-Column Editorial Scripture Body & Drop Cap */
.magazine-columns {
  column-count: 2;
  column-gap: 7mm;
  column-rule: 1px solid rgba(209, 199, 189, 0.6);
}

.scripture-body {
  font-family: 'EB Garamond', 'Cardo', Georgia, serif;
  font-size: 11pt;
  line-height: 1.58;
  text-align: justify;
  hyphens: auto;
  color: #24201c;
}

.scripture-body :deep(p) {
  margin-bottom: 2.5mm;
}

.magazine-dropcap::first-letter,
.magazine-dropcap :deep(p:first-of-type)::first-letter {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 3.2em;
  font-weight: 700;
  line-height: 0.82;
  float: left;
  padding-right: 0.09em;
  padding-top: 0.04em;
  color: #72383d;
}

.psalm-body {
  font-family: 'EB Garamond', 'Cardo', Georgia, serif;
  font-size: 11.2pt;
  line-height: 1.68;
  text-align: center;
  white-space: pre-line;
  color: #24201c;
  max-width: 140mm;
  margin: 0 auto;
}

/* Editorial Pull Quote */
.magazine-pullquote {
  margin: 5mm 0;
  padding: 3.5mm 6mm;
  border-top: 1px solid #c5a059;
  border-bottom: 1px solid #c5a059;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 12.5pt;
  font-style: italic;
  line-height: 1.45;
  text-align: center;
  color: #72383d;
  break-inside: avoid;
}

/* AI Commentary Sidebar Box */
.ai-commentary-box {
  margin-top: 5mm;
  background: #faf6ee;
  border: 1px solid #dfd6c9;
  border-left: 3.5px solid #72383d;
  padding: 4.5mm 5mm;
  break-inside: avoid;
}

.ai-commentary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  border-bottom: 1px solid #e4ddd3;
  padding-bottom: 2mm;
  margin-bottom: 2.5mm;
}

.ai-commentary-kicker {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8.5pt;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #72383d;
}

.ai-commentary-tag {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 600;
  color: #5c5149;
}

.ai-commentary-body {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 10.5pt;
  line-height: 1.54;
  color: #2d2926;
}

.ai-commentary-body :deep(.mag-commentary-p) {
  margin-bottom: 2mm;
}

.ai-commentary-body :deep(.mag-commentary-p:last-child) {
  margin-bottom: 0;
}

.ai-commentary-body :deep(.mag-commentary-ul) {
  list-style: disc;
  padding-left: 5mm;
  margin: 2mm 0;
}

.ai-commentary-body :deep(.mag-commentary-li) {
  margin-bottom: 1mm;
}

.magisterium-sources {
  margin-top: 3mm;
  padding-top: 2.5mm;
  border-top: 1px dashed #d1c7bd;
}

.magisterium-sources-title {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #5c5149;
  margin-bottom: 1.5mm;
}

.magisterium-source-item {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8.5pt;
  line-height: 1.4;
  margin-bottom: 1mm;
}

/* Comics Section at the End */
.comics-masthead {
  border-bottom: 2px solid #1e1b18;
  padding-bottom: 4.5mm;
  margin-bottom: 5mm;
}

.comics-meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 2.5mm;
}

.comics-meta-pill {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8pt;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 999px;
  background: #faf6ee;
  border: 1px solid #d1c7bd;
  color: #322d29;
}

.comics-lesson-box {
  margin-top: 3.5mm;
  padding: 3mm 4mm;
  background: #fffdf5;
  border: 1.5px solid #1e1b18;
  border-left: 4px solid #72383d;
  border-radius: 6px;
}

.comics-lesson-kicker {
  display: block;
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  color: #72383d;
}

.comics-lesson-text {
  font-family: 'Comic Neue', 'Comic Sans MS', sans-serif;
  font-size: 10.5pt;
  font-weight: 700;
  font-style: italic;
  color: #1e1b18;
  margin-top: 1mm;
}

.comics-rows-stack {
  display: flex;
  flex-direction: column;
  gap: 5mm;
}

.comic-layout-row {
  break-inside: avoid;
}

.comic-tile {
  border: 2px solid #1e1b18;
  border-radius: 8px;
  overflow: hidden;
  background: #fffdf5;
  display: flex;
  flex-direction: column;
}

.comic-tile-bar {
  background: #1e1b18;
  color: #ffffff;
  padding: 4px 10px;
  font-family: 'Bangers', 'Comic Neue', Impact, sans-serif;
  font-size: 10pt;
  letter-spacing: 0.06em;
}

.comic-tile-art {
  width: 100%;
  background: #f5f0e6;
  overflow: hidden;
}

.comic-tile-caption {
  padding: 2.5mm 3mm;
  border-top: 1.5px solid #1e1b18;
  display: flex;
  flex-direction: column;
  gap: 2mm;
}

.comic-narrative {
  font-family: 'Comic Neue', 'Comic Sans MS', sans-serif;
  font-size: 9.5pt;
  font-weight: 700;
  line-height: 1.38;
  color: #1e1b18;
}

.comic-dialogue-bubble {
  background: #ffffff;
  border: 1.5px solid #1e1b18;
  border-radius: 8px;
  padding: 1.5mm 2.5mm;
  font-family: 'Comic Neue', 'Comic Sans MS', sans-serif;
  font-size: 9pt;
  font-weight: 700;
  font-style: italic;
  color: #1e1b18;
}

/* Colophon */
.magazine-colophon {
  text-align: center;
  padding: 6mm 10mm 10mm 10mm;
  break-inside: avoid;
}

.colophon-rule {
  height: 1px;
  background: #c5a059;
  width: 50mm;
  margin: 0 auto 4mm auto;
}

.colophon-pastoral {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 10pt;
  font-style: italic;
  color: #322d29;
  line-height: 1.5;
  margin-bottom: 2mm;
}

.colophon-sub {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8pt;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #72383d;
  margin-bottom: 2mm;
}

.colophon-copyright {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  color: #5c5149;
}

/* Interactive Buttons (hidden in print) */
.mag-gen-btn {
  margin-top: 3mm;
  padding: 7px 14px;
  border-radius: 999px;
  background: #72383d;
  color: #ffffff;
  font-family: 'Source Sans 3', sans-serif;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  cursor: pointer;
}

.mag-gen-btn-sm {
  padding: 5px 11px;
  border-radius: 8px;
  background: #72383d;
  color: #ffffff;
  font-family: 'Source Sans 3', sans-serif;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
}

.mag-regen-link {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 11px;
  font-weight: 700;
  color: #72383d;
  cursor: pointer;
}

@media (max-width: 768px) {
  .magazine-sheet {
    padding: 6mm;
  }
  .cover-lead-grid {
    grid-template-columns: 1fr;
  }
  .magazine-columns {
    column-count: 1;
  }
  .masthead-title {
    font-size: 24pt;
  }
}

@media print {
  .magazine-root {
    max-width: none;
    margin: 0;
  }
  .magazine-sheet {
    border: none;
    border-radius: 0;
    box-shadow: none;
    padding: 0;
    margin: 0 0 8mm 0;
    background: #ffffff !important;
  }
  .magazine-cover {
    break-after: page;
  }
  .magazine-article {
    break-after: page;
  }
  .magazine-comics-section {
    break-before: page;
  }
}
</style>
