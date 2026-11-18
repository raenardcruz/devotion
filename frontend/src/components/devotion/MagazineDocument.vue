<script setup lang="ts">
import { computed } from 'vue';
import { renderMarkdown } from '../../utils/markdown';
import {
  groupPanelsIntoComicRows,
  getTileAspectClass,
  getTileColSpanClass,
} from '../../utils/mangaLayout';
import {
  computeArticleLayoutPlan,
  computePopeQuoteLayout,
  extractCoverDeckQuote,
  type ArticleLayoutPlan,
} from '../../utils/magazineLayout';
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

const parseCitations = (reading?: MagazineReadingItem): MagisteriumCitation[] => {
  if (!reading) return [];
  if (reading.citations && reading.citations.length > 0) {
    return reading.citations.slice(0, 4);
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
    return out.slice(0, 4);
  }
  return [];
};

const commentaryMarkdownOptions = {
  paragraphClass: 'mag-commentary-p',
  listClass: 'mag-commentary-ul',
  orderedListClass: 'mag-commentary-ol',
  listItemClass: 'mag-commentary-li',
};

const issueSubtitle = computed(() => {
  const d = new Date(`${props.date}T12:00:00`);
  if (Number.isNaN(d.getTime())) return 'Liturgical Daily Edition';
  const weekday = d.toLocaleDateString('en-US', { weekday: 'long' });
  return `${weekday} of the Liturgical Year · Daily Liturgical Edition`;
});

const coverDeckQuote = computed(() =>
  extractCoverDeckQuote(props.readings.gospel?.text, props.readings.gospel?.context)
);

interface EnrichedReadingSection {
  key: Exclude<MagazineReadingKey, 'cover'>;
  kicker: string;
  label: string;
  isPsalm: boolean;
  reading: MagazineReadingItem;
  imageSrc: string;
  citations: MagisteriumCitation[];
  layout: ArticleLayoutPlan;
}

const readingSections = computed<EnrichedReadingSection[]>(() => {
  const rawList: Array<{
    key: Exclude<MagazineReadingKey, 'cover'>;
    kicker: string;
    label: string;
    isPsalm?: boolean;
    reading: MagazineReadingItem;
  }> = [];

  if (props.readings.first_reading) {
    rawList.push({
      key: 'first_reading',
      kicker: 'Part I · First Reading',
      label: 'Sacred Scripture',
      reading: props.readings.first_reading,
    });
  }
  if (props.readings.responsorial_psalm) {
    rawList.push({
      key: 'responsorial_psalm',
      kicker: 'Responsorial Psalm',
      label: 'Liturgical Song',
      isPsalm: true,
      reading: props.readings.responsorial_psalm,
    });
  }
  if (props.readings.second_reading) {
    rawList.push({
      key: 'second_reading',
      kicker: 'Second Reading',
      label: 'Apostolic Epistle',
      reading: props.readings.second_reading,
    });
  }
  if (props.readings.gospel) {
    rawList.push({
      key: 'gospel',
      kicker: 'The Holy Gospel',
      label: 'Gospel Proclamation',
      reading: props.readings.gospel,
    });
  }

  return rawList.map((item) => {
    const imageSrc = props.includeImages ? props.resolveImageSrc(props.images[item.key]) : '';
    const citations = parseCitations(item.reading);
    const layout = computeArticleLayoutPlan({
      readingText: item.reading.text,
      exegesisRaw: item.reading.context,
      isPsalm: Boolean(item.isPsalm),
      hasImage: Boolean(imageSrc),
      includeCommentaries: props.includeCommentaries,
      citationCount: citations.length,
    });

    return {
      key: item.key,
      kicker: item.kicker,
      label: item.label,
      isPsalm: Boolean(item.isPsalm),
      reading: item.reading,
      imageSrc,
      citations,
      layout,
    };
  });
});

const popeQuoteLayout = computed(() =>
  computePopeQuoteLayout(props.readings.pope_quote)
);

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
          <span>VOL ISSUE {{ date }}</span>
          <span>CATHOLIC DAILY LITURGICAL MAGAZINE</span>
        </div>
        <h1 class="masthead-title">LIGHT OF THE LITURGY</h1>
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
          <p v-if="coverDeckQuote" class="cover-deck">
            “{{ coverDeckQuote }}”
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
        <span>Sacred Scripture · AI Assisted Theological Commentary · Illustrated Parable</span>
      </footer>
    </section>

    <!-- ==================== INTERIOR: READINGS & ADAPTIVE EXEGESIS ==================== -->
    <section
      v-for="sec in readingSections"
      :key="sec.key"
      class="magazine-sheet magazine-article"
      :class="[
        `article-scripture--${sec.layout.scriptureLayoutMode}`,
        `article-exegesis--${sec.layout.exegesisLayoutMode}`,
      ]"
    >
      <!-- Running Header -->
      <div class="running-header">
        <span>LIGHT OF THE LITURGY · {{ formattedDate.toUpperCase() }}</span>
        <span>{{ sec.kicker.toUpperCase() }}</span>
      </div>

      <!-- Article Header -->
      <header
        class="article-header"
        :class="{ 'article-header--centered': sec.layout.scriptureLayoutMode === 'psalm' }"
      >
        <div
          class="article-kicker-row"
          :class="{ 'article-kicker-row--centered': sec.layout.scriptureLayoutMode === 'psalm' }"
        >
          <span class="mag-kicker">{{ sec.kicker }}</span>
          <span class="mag-badge">{{ sec.label }}</span>
        </div>
        <h2 class="article-title">{{ sec.reading.citation }}</h2>
      </header>

      <!-- ========== SCRIPTURE OPENING: SPLIT LEAD (SHORT READING + IMAGE) ========== -->
      <div
        v-if="sec.layout.scriptureLayoutMode === 'split-lead'"
        class="scripture-split-lead"
      >
        <figure class="article-figure article-figure--split">
          <div class="article-image-frame">
            <img
              :src="sec.imageSrc"
              :alt="`Illustration for ${sec.reading.citation}`"
              class="article-image article-image--split"
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
        </figure>

        <div
          class="scripture-body scripture-body--split"
          v-html="sec.layout.formattedScriptureHtml"
        ></div>
      </div>

      <!-- ========== SCRIPTURE OPENING: STANDARD / PSALM / MEDIUM / LONG ========== -->
      <template v-else>
        <!-- Reading Editorial Illustration (or interactive prompt) -->
        <figure
          v-if="includeImages && (sec.imageSrc || interactive)"
          class="article-figure"
          :class="{
            'article-figure--psalm': sec.layout.scriptureLayoutMode === 'psalm',
            'article-figure--long': sec.layout.scriptureLayoutMode === 'two-column-long',
          }"
        >
          <template v-if="sec.imageSrc">
            <div class="article-image-frame">
              <img
                :src="sec.imageSrc"
                :alt="`Illustration for ${sec.reading.citation}`"
                class="article-image"
                :class="{
                  'article-image--panoramic': sec.layout.scriptureLayoutMode === 'two-column-long',
                }"
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

        <!-- Psalm Poetic Body -->
        <div
          v-if="sec.layout.scriptureLayoutMode === 'psalm'"
          class="psalm-body"
          v-html="sec.layout.formattedScriptureHtml"
        ></div>

        <!-- Short Reading Proclamation (single-column centered measure, no cramped 2-col stub) -->
        <div
          v-else-if="sec.layout.scriptureLayoutMode === 'single-column'"
          class="scripture-proclamation scripture-body"
          v-html="sec.layout.formattedScriptureHtml"
        ></div>

        <!-- Medium or Long Reading (balanced 2-column magazine flow) -->
        <div
          v-else
          class="magazine-columns scripture-body"
          :class="{ 'scripture-body--long': sec.layout.scriptureLayoutMode === 'two-column-long' }"
          v-html="sec.layout.formattedScriptureHtml"
        ></div>
      </template>

      <!-- Standalone Scripture Pull Quote (only when editorially warranted, never redundant) -->
      <blockquote v-if="sec.layout.scripturePullQuote" class="magazine-pullquote">
        “{{ sec.layout.scripturePullQuote }}”
      </blockquote>

      <!-- ==================== EXEGESIS MODE 1: BRIEF MEDITATION CARD ==================== -->
      <aside
        v-if="sec.layout.exegesisLayoutMode === 'brief-card'"
        class="exegesis-card"
        :class="{
          'exegesis-card--centered': !sec.layout.briefCardHasSideSources,
          'exegesis-card--with-sources': sec.layout.briefCardHasSideSources,
        }"
      >
        <div class="exegesis-card-header">
          <span class="exegesis-kicker">{{ sec.layout.exegesisKicker }}</span>
          <span class="exegesis-tag">AI-Assisted Theological Commentary</span>
        </div>

        <div :class="{ 'exegesis-card-grid': sec.layout.briefCardHasSideSources }">
          <div
            class="exegesis-body exegesis-card-body"
            v-html="renderMarkdown(sec.layout.exegesisMarkdown, commentaryMarkdownOptions)"
          ></div>

          <div
            v-if="sec.citations.length > 0"
            class="exegesis-card-sources"
          >
            <p class="magisterium-sources-title">Magisterium &amp; Sources</p>
            <div
              v-for="(cit, cIdx) in sec.citations"
              :key="cIdx"
              class="magisterium-source-item"
            >
              <span class="source-title">{{ cIdx + 1 }}. {{ cit.document_title || cit.title }}</span>
              <span v-if="cit.author || cit.ref" class="source-meta">
                — {{ [cit.author, cit.ref].filter(Boolean).join(', ') }}
              </span>
            </div>
          </div>
        </div>
      </aside>

      <!-- ==================== EXEGESIS MODE 2: EDITORIAL FEATURE ESSAY ==================== -->
      <div
        v-else-if="sec.layout.exegesisLayoutMode === 'editorial-essay'"
        class="exegesis-essay"
        :class="{ 'exegesis-essay--page-break': sec.layout.exegesisStartsNewPageInPrint }"
      >
        <!-- Print-only continuation running header when essay starts cleanly on a new page -->
        <div
          v-if="sec.layout.exegesisStartsNewPageInPrint"
          class="exegesis-print-running-header"
        >
          <span>LIGHT OF THE LITURGY · {{ formattedDate.toUpperCase() }}</span>
          <span>EXEGESIS · {{ sec.reading.citation.toUpperCase() }}</span>
        </div>

        <!-- Ornamental Editorial Rule separating Scripture from Exegesis -->
        <div
          class="exegesis-ornament-rule"
          :class="{ 'exegesis-ornament-rule--hide-print': sec.layout.exegesisStartsNewPageInPrint }"
          aria-hidden="true"
        >
          <span class="exegesis-ornament-line"></span>
          <span class="exegesis-ornament-glyph">☩</span>
          <span class="exegesis-ornament-line"></span>
        </div>

        <!-- Essay Masthead -->
        <header class="exegesis-essay-header">
          <div>
            <span class="exegesis-kicker">{{ sec.layout.exegesisKicker }}</span>
            <p class="exegesis-essay-sub">
              Theological &amp; Scriptural Reflection on {{ sec.reading.citation }}
            </p>
          </div>
          <span class="exegesis-tag">AI-Assisted Commentary</span>
        </header>

        <!-- Thematic Sections / Balanced Editorial Beats -->
        <div class="exegesis-blocks-stack">
          <template v-for="block in sec.layout.exegesisBlocks" :key="block.id">
            <!-- Mid-essay interlude pull quote for long unstructured exegesis -->
            <blockquote
              v-if="block.precedingPullQuote"
              class="magazine-pullquote magazine-pullquote--interlude"
            >
              “{{ block.precedingPullQuote }}”
            </blockquote>

            <div
              class="exegesis-section-block"
              :class="{ 'exegesis-section-block--standfirst': block.isStandfirst }"
            >
              <h3 v-if="block.heading" class="exegesis-subhead">
                <span class="exegesis-subhead-mark" aria-hidden="true">◆</span>
                <span>{{ block.heading }}</span>
              </h3>

              <div
                class="exegesis-body"
                :class="{
                  'exegesis-standfirst': block.isStandfirst,
                  'exegesis-columns': block.useTwoColumns,
                  'exegesis-single': !block.useTwoColumns && !block.isStandfirst,
                }"
                v-html="renderMarkdown(block.markdown, commentaryMarkdownOptions)"
              ></div>
            </div>
          </template>
        </div>

        <!-- Compact Magisterium & Scriptural Sources Strip -->
        <footer
          v-if="sec.citations.length > 0"
          class="exegesis-sources-footer"
        >
          <p class="magisterium-sources-title">Magisterium &amp; Scriptural Sources</p>
          <div
            class="magisterium-sources-grid"
            :class="{ 'magisterium-sources-grid--single': sec.citations.length === 1 }"
          >
            <div
              v-for="(cit, cIdx) in sec.citations"
              :key="cIdx"
              class="magisterium-source-item"
            >
              <span class="source-title">{{ cIdx + 1 }}. {{ cit.document_title || cit.title }}</span>
              <span v-if="cit.author || cit.ref" class="source-meta">
                — {{ [cit.author, cit.ref].filter(Boolean).join(', ') }}
              </span>
            </div>
          </div>
        </footer>
      </div>
    </section>

    <!-- ==================== WORDS OF THE POPES ==================== -->
    <section v-if="readings.pope_quote" class="magazine-sheet magazine-article">
      <div class="running-header">
        <span>LIGHT OF THE LITURGY · {{ formattedDate.toUpperCase() }}</span>
        <span>PAPAL MAGISTERIUM</span>
      </div>
      <header class="article-header">
        <div class="article-kicker-row">
          <span class="mag-kicker">Words of the Popes</span>
          <span class="mag-badge">Holy Father's Reflection</span>
        </div>
        <h2 class="article-title">Papal Reflection on Today's Liturgy</h2>
      </header>
      <div
        v-if="popeQuoteLayout.isShort"
        class="scripture-proclamation scripture-body"
        v-html="popeQuoteLayout.formattedHtml"
      ></div>
      <div
        v-else
        class="magazine-columns scripture-body"
        v-html="popeQuoteLayout.formattedHtml"
      ></div>
      <p class="mag-caption mt-4">Source: Vatican News Daily Liturgy</p>
    </section>

    <!-- ==================== AT THE END: DAILY PARABLE & COMICS ==================== -->
    <section v-if="hasComics && storyAnalysis" class="magazine-sheet magazine-comics-section">
      <div class="running-header">
        <span>LIGHT OF THE LITURGY · {{ formattedDate.toUpperCase() }}</span>
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
            :class="[getTileColSpanClass(tile.colSpan), `comic-span-${tile.colSpan}`]"
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
        Light of the Liturgy · Daily Catholic Readings Magazine · Guided by Pope Leo XIV's Encyclical <em>Magnifica Humanitas</em>
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

/* Harmonize API.Bible verse numbers with warm editorial palette */
.magazine-root :deep(.v) {
  font-family: 'Source Sans 3', 'Inter', sans-serif !important;
  font-size: 7.2pt !important;
  font-weight: 700 !important;
  color: #8b5a5e !important;
  vertical-align: super !important;
  line-height: 0 !important;
  margin-right: 2px !important;
  margin-left: 1px !important;
}

.magazine-root :deep(.v.mag-v-lead) {
  display: none !important;
}

.magazine-sheet {
  background: #ffffff;
  border: 1px solid #e4ddd3;
  border-radius: 4px;
  padding: 15mm 16mm 14mm 16mm;
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
  font-size: 8.2pt;
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
  margin-bottom: 4.5mm;
}

.article-header {
  margin-bottom: 4.5mm;
  padding-bottom: 3mm;
  border-bottom: 1px solid #ece5dc;
}

.article-header--centered {
  text-align: center;
}

.article-kicker-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.article-kicker-row--centered {
  justify-content: center;
  gap: 10px;
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

/* Article Figure & Split Lead */
.article-figure {
  margin: 0 0 5mm 0;
  break-inside: avoid;
}

.article-figure--psalm {
  max-width: 132mm;
  margin: 0 auto 5.5mm auto;
}

.article-figure--split {
  margin: 0;
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

.article-image--split {
  aspect-ratio: 4 / 3;
}

.article-image--panoramic {
  aspect-ratio: 2 / 1;
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

/* Split Lead: Short Reading paired cleanly beside its Sacred Art */
.scripture-split-lead {
  display: grid;
  grid-template-columns: 0.92fr 1.08fr;
  gap: 6.5mm;
  align-items: center;
  margin-bottom: 5mm;
}

/* Short Reading Proclamation Band (no image) */
.scripture-proclamation {
  max-width: 140mm;
  margin: 0 auto 4mm auto;
  padding: 3.5mm 5.5mm;
  border-left: 2.5px solid #c5a059;
  background: linear-gradient(90deg, rgba(250, 246, 238, 0.65) 0%, rgba(255, 255, 255, 0) 100%);
  font-size: 11.4pt !important;
  line-height: 1.66 !important;
  text-align: left !important;
}

/* Multi-Column Editorial Scripture Body */
.magazine-columns {
  column-count: 2;
  column-gap: 7mm;
  column-rule: 1px solid rgba(209, 199, 189, 0.6);
}

.scripture-body {
  font-family: 'EB Garamond', 'Cardo', Georgia, serif;
  font-size: 10.8pt;
  line-height: 1.56;
  text-align: justify;
  hyphens: auto;
  color: #24201c;
}

.scripture-body--split {
  font-size: 10.9pt;
  line-height: 1.62;
  padding-left: 4mm;
  border-left: 1.5px solid rgba(197, 160, 89, 0.55);
}

.scripture-body--long {
  font-size: 10.3pt;
  line-height: 1.52;
}

.scripture-body :deep(p) {
  margin-bottom: 2.4mm;
}

.scripture-body :deep(p:last-child) {
  margin-bottom: 0;
}

/* API.Bible internal subheadings and poetic lines inside prose readings */
.scripture-body :deep(.s),
.scripture-body :deep(.s1),
.scripture-body :deep(.s2),
.scripture-body :deep(.ms),
.scripture-body :deep(.ms1) {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 10.4pt;
  font-weight: 700;
  font-style: italic;
  color: #72383d;
  text-align: left;
  margin: 3mm 0 1.4mm 0;
  break-after: avoid;
}

.scripture-body :deep(.q1) {
  padding-left: 4mm;
  text-indent: -2mm;
  margin-bottom: 1mm;
  text-align: left;
}

.scripture-body :deep(.q2),
.scripture-body :deep(.q3) {
  padding-left: 7mm;
  text-indent: -2mm;
  margin-bottom: 1mm;
  text-align: left;
}

/* Illuminated Magazine Drop Cap (explicitly wrapped around opening letter) */
.scripture-body :deep(.mag-dropcap-letter) {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 3.1em;
  font-weight: 700;
  line-height: 0.8;
  float: left;
  padding-right: 0.09em;
  padding-top: 0.05em;
  color: #72383d;
}

.scripture-body :deep(.mag-dropcap-quote) {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1.45em;
  font-weight: 700;
  line-height: 0.9;
  float: left;
  padding-right: 0.03em;
  color: #c5a059;
}

/* Psalm Poetic Layout */
.psalm-body {
  font-family: 'EB Garamond', 'Cardo', Georgia, serif;
  font-size: 11.2pt;
  line-height: 1.68;
  text-align: center;
  white-space: pre-line;
  color: #24201c;
  max-width: 132mm;
  margin: 0 auto;
}

.psalm-body :deep(p) {
  margin-bottom: 2.2mm;
}

/* Editorial Pull Quote */
.magazine-pullquote {
  margin: 5mm auto;
  padding: 3.2mm 6.5mm;
  max-width: 152mm;
  border-top: 1px solid #c5a059;
  border-bottom: 1px solid #c5a059;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 12pt;
  font-style: italic;
  line-height: 1.45;
  text-align: center;
  color: #72383d;
  break-inside: avoid;
}

.magazine-pullquote--interlude {
  margin: 4.5mm auto;
  font-size: 11.5pt;
  background: linear-gradient(90deg, transparent, rgba(250, 246, 238, 0.7), transparent);
}

/* ==================== EXEGESIS MODE 1: BRIEF MEDITATION CARD ==================== */
.exegesis-card {
  margin-top: 5.5mm;
  background: #faf6ee;
  border: 1px solid #dfd6c9;
  border-left: 3.5px solid #72383d;
  padding: 4.5mm 5.5mm;
}

.exegesis-card--centered {
  max-width: 152mm;
  margin-left: auto;
  margin-right: auto;
}

.exegesis-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  border-bottom: 1px solid #e4ddd3;
  padding-bottom: 2mm;
  margin-bottom: 2.8mm;
}

.exegesis-card-grid {
  display: grid;
  grid-template-columns: 1.55fr 1fr;
  gap: 5.5mm;
  align-items: start;
}

.exegesis-card-sources {
  border-left: 1px solid #dfd6c9;
  padding-left: 4.5mm;
}

/* ==================== EXEGESIS MODE 2: EDITORIAL FEATURE ESSAY ==================== */
.exegesis-essay {
  margin-top: 5.5mm;
}

.exegesis-print-running-header {
  display: none;
}

.exegesis-ornament-rule {
  display: flex;
  align-items: center;
  gap: 4mm;
  margin: 5mm 0 4mm 0;
}

.exegesis-ornament-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg, rgba(197, 160, 89, 0.25), #c5a059, rgba(197, 160, 89, 0.25));
}

.exegesis-ornament-glyph {
  font-size: 11pt;
  color: #72383d;
  line-height: 1;
}

.exegesis-essay-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 10px;
  border-bottom: 1.5px solid #1e1b18;
  padding-bottom: 2.5mm;
  margin-bottom: 4.5mm;
  break-after: avoid;
}

.exegesis-kicker {
  display: block;
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8.5pt;
  font-weight: 700;
  letter-spacing: 0.17em;
  text-transform: uppercase;
  color: #72383d;
}

.exegesis-essay-sub {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 11pt;
  font-style: italic;
  color: #4b433c;
  margin: 0.8mm 0 0 0;
  line-height: 1.25;
}

.exegesis-tag {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #5c5149;
  white-space: nowrap;
}

.exegesis-blocks-stack {
  display: flex;
  flex-direction: column;
  gap: 4.5mm;
}

.exegesis-section-block--standfirst {
  margin-bottom: 0.5mm;
}

.exegesis-subhead {
  display: flex;
  align-items: baseline;
  gap: 2mm;
  font-family: 'Playfair Display', Georgia, serif !important;
  font-size: 11.4pt !important;
  font-weight: 700;
  letter-spacing: -0.005em;
  line-height: 1.25;
  color: #72383d;
  margin: 0 0 2.2mm 0;
  padding-bottom: 1.2mm;
  border-bottom: 1px solid #eae2d7;
  break-after: avoid;
  page-break-after: avoid;
}

.exegesis-subhead-mark {
  font-size: 7pt;
  color: #c5a059;
  transform: translateY(-1px);
  flex-shrink: 0;
}

/* Exegesis Body Typography (shared between brief card and essay) */
.exegesis-body {
  font-family: 'EB Garamond', 'Cardo', Georgia, serif;
  font-size: 10.2pt;
  line-height: 1.54;
  color: #26221e;
}

.exegesis-card-body {
  font-size: 10.4pt;
  line-height: 1.56;
}

.exegesis-standfirst {
  font-size: 10.9pt;
  font-style: italic;
  line-height: 1.56;
  color: #36302b;
  padding-left: 3.5mm;
  border-left: 2px solid #c5a059;
  max-width: 158mm;
}

.exegesis-columns {
  column-count: 2;
  column-gap: 7mm;
  column-rule: 1px solid rgba(209, 199, 189, 0.55);
  text-align: justify;
  hyphens: auto;
}

.exegesis-single {
  max-width: 148mm;
  text-align: left;
}

.exegesis-body :deep(.mag-commentary-p) {
  margin-bottom: 2.2mm;
}

.exegesis-body :deep(.mag-commentary-p:last-child) {
  margin-bottom: 0;
}

.exegesis-body :deep(.mag-commentary-ul),
.exegesis-body :deep(.mag-commentary-ol) {
  padding-left: 4.5mm;
  margin: 1.5mm 0 2.2mm 0;
  text-align: left;
}

.exegesis-body :deep(.mag-commentary-ul) {
  list-style: disc;
}

.exegesis-body :deep(.mag-commentary-ol) {
  list-style: decimal;
}

.exegesis-body :deep(.mag-commentary-li) {
  margin-bottom: 1.4mm;
  break-inside: avoid;
  page-break-inside: avoid;
}

.exegesis-body :deep(.mag-commentary-li:last-child) {
  margin-bottom: 0;
}

/* Override global Cinzel Decorative h3/h4/h5 if any nested markdown headers appear */
.exegesis-body :deep(h3),
.exegesis-body :deep(h4),
.exegesis-body :deep(h5) {
  font-family: 'Source Sans 3', 'Inter', sans-serif !important;
  font-size: 8.8pt !important;
  font-weight: 700 !important;
  letter-spacing: 0.1em !important;
  text-transform: uppercase !important;
  color: #72383d !important;
  margin: 3mm 0 1.4mm 0 !important;
  line-height: 1.3 !important;
  text-align: left;
  break-after: avoid;
  page-break-after: avoid;
}

.exegesis-body :deep(strong) {
  font-weight: 700;
  color: #1e1b18;
}

.exegesis-body :deep(blockquote) {
  margin: 2.2mm 0;
  padding: 1.5mm 3mm;
  border-left: 2px solid #72383d;
  background: #faf6ee;
  font-style: italic;
  color: #36302b;
  text-align: left;
  break-inside: avoid;
  page-break-inside: avoid;
}

/* Compact Magisterium & Scriptural Sources Strip */
.exegesis-sources-footer {
  margin-top: 5mm;
  padding: 3.5mm 4.5mm;
  background: #faf6ee;
  border-top: 2px solid #c5a059;
  border-bottom: 1px solid #dfd6c9;
  break-inside: avoid;
  page-break-inside: avoid;
}

.magisterium-sources-title {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 7.5pt;
  font-weight: 700;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  color: #72383d;
  margin: 0 0 2mm 0;
  break-after: avoid;
  page-break-after: avoid;
}

.magisterium-sources-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 2mm 5mm;
}

.magisterium-sources-grid--single {
  grid-template-columns: 1fr;
}

.magisterium-source-item {
  font-family: 'Source Sans 3', sans-serif;
  font-size: 8.3pt;
  line-height: 1.38;
  margin-bottom: 1mm;
  break-inside: avoid;
  page-break-inside: avoid;
}

.source-title {
  font-weight: 700;
  color: #72383d;
}

.source-meta {
  color: #5c5149;
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
  background: #ffffff;
  border: 1px solid #e4ddd3;
  border-radius: 4px;
  padding: 8mm 12mm 10mm 12mm;
  box-shadow: 0 8px 28px rgba(30, 27, 24, 0.08);
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

@media screen and (max-width: 768px) {
  .magazine-sheet {
    padding: 6mm;
  }
  .cover-lead-grid,
  .scripture-split-lead,
  .exegesis-card-grid,
  .magisterium-sources-grid {
    grid-template-columns: 1fr;
  }
  .exegesis-card-sources {
    border-left: none;
    border-top: 1px dashed #d1c7bd;
    padding-left: 0;
    padding-top: 3mm;
  }
  .scripture-body--split {
    border-left: none;
    padding-left: 0;
  }
  .magazine-columns,
  .exegesis-columns {
    column-count: 1;
    text-align: left;
  }
  .exegesis-essay-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  .masthead-title {
    font-size: 24pt;
  }
}

@media print {
  .magazine-root {
    max-width: none;
    width: 100%;
    margin: 0;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
  .magazine-sheet {
    border: none;
    border-radius: 0;
    box-shadow: none;
    padding: 0;
    margin: 0;
    background: #ffffff !important;
  }
  .magazine-colophon {
    border: none;
    border-radius: 0;
    box-shadow: none;
    padding: 6mm 10mm 10mm 10mm;
    background: #ffffff !important;
  }
  .magazine-sheet + .magazine-sheet {
    break-before: page;
    page-break-before: always;
  }
  .masthead-title {
    font-size: 34pt;
  }
  .cover-lead-grid {
    grid-template-columns: 1.35fr 1fr;
    gap: 6mm;
    padding-top: 4mm;
  }
  .magazine-cover .cover-masthead {
    padding-bottom: 4mm;
    margin-bottom: 4.5mm;
  }
  .magazine-cover .cover-figure {
    margin-bottom: 4.5mm;
  }
  .magazine-cover .cover-image {
    max-height: 110mm;
  }
  .magazine-cover .cover-footer {
    margin-top: 4.5mm;
  }
  .article-image {
    max-height: 76mm;
  }
  .article-image--panoramic {
    max-height: 66mm;
  }
  .article-image--split {
    max-height: 68mm;
  }
  .article-figure--psalm .article-image {
    max-height: 66mm;
  }
  .scripture-split-lead {
    display: grid !important;
    grid-template-columns: 0.92fr 1.08fr !important;
    gap: 6.5mm !important;
    break-inside: avoid;
    page-break-inside: avoid;
  }
  .magazine-columns {
    column-count: 2 !important;
    column-gap: 7mm;
    column-rule: 1px solid rgba(209, 199, 189, 0.6);
  }
  .exegesis-columns {
    column-count: 2 !important;
    column-gap: 7mm;
    column-rule: 1px solid rgba(209, 199, 189, 0.55);
  }
  .exegesis-card {
    break-inside: avoid;
    page-break-inside: avoid;
  }
  .exegesis-card-grid {
    display: grid !important;
    grid-template-columns: 1.55fr 1fr !important;
    gap: 5.5mm !important;
  }
  .exegesis-essay--page-break {
    break-before: page;
    page-break-before: always;
    margin-top: 0;
  }
  .exegesis-print-running-header {
    display: flex !important;
    justify-content: space-between;
    font-family: 'Source Sans 3', sans-serif;
    font-size: 7.5pt;
    font-weight: 700;
    letter-spacing: 0.18em;
    color: #72383d;
    border-bottom: 1px solid #d1c7bd;
    padding-bottom: 2mm;
    margin-bottom: 4.5mm;
  }
  .exegesis-ornament-rule--hide-print {
    display: none !important;
  }
  .exegesis-essay-header {
    break-after: avoid;
    page-break-after: avoid;
  }
  .exegesis-subhead {
    break-after: avoid;
    page-break-after: avoid;
  }
  .article-header {
    break-after: avoid;
    page-break-after: avoid;
  }
  .scripture-body :deep(p),
  .exegesis-body :deep(.mag-commentary-p),
  .exegesis-body :deep(.mag-commentary-li) {
    orphans: 3;
    widows: 3;
  }
  .magisterium-sources-grid {
    display: grid !important;
    grid-template-columns: repeat(2, minmax(0, 1fr)) !important;
  }
  .magisterium-sources-grid--single {
    grid-template-columns: 1fr !important;
  }
  .article-figure:has(> .article-image-empty):not(:has(> .article-image-frame)),
  .article-image-empty,
  .mag-gen-btn,
  .mag-gen-btn-sm,
  .mag-regen-link {
    display: none !important;
  }
  .comic-layout-row {
    display: grid !important;
    grid-template-columns: repeat(6, minmax(0, 1fr)) !important;
    gap: 4mm !important;
    align-items: stretch !important;
    break-inside: avoid;
    page-break-inside: avoid;
  }
  .comic-span-6 {
    grid-column: span 6 / span 6 !important;
  }
  .comic-span-4 {
    grid-column: span 4 / span 4 !important;
  }
  .comic-span-3 {
    grid-column: span 3 / span 3 !important;
  }
  .comic-span-2 {
    grid-column: span 2 / span 2 !important;
  }
}
</style>
