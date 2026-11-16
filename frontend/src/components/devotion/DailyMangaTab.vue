<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import {
  useStoryToManga,
  type MangaPanelPlan,
} from '../../composables/useStoryToManga';
import { exportMangaAsA4Pdf } from '../../utils/mangaPdfExport';
import {
  groupPanelsIntoComicRows,
  getTileAspectClass,
  getTileColSpanClass,
} from '../../utils/mangaLayout';

const props = defineProps<{
  date: string;
  formattedDate: string;
  isAdmin: boolean;
}>();

const emit = defineEmits<{
  (e: 'open-generator'): void;
}>();

const {
  style,
  moralLesson,
  storyAnalysis,
  storyBreakdown,
  generatedPanels,
  currentStoryId,
  currentVersion,
  availableVersions,
  getProxyAssetUrl,
  loadStateFromServer,
  loadVersionFromServer,
} = useStoryToManga();

const isLoading = ref(false);
const layoutMode = ref<'storybook' | 'grid'>('grid');
const isExportingPdf = ref(false);
const lightboxImage = ref<{ url: string; caption: string } | null>(null);

const comicLayoutRows = computed(() =>
  groupPanelsIntoComicRows(storyBreakdown.value?.panels || [])
);

const loadManga = async () => {
  isLoading.value = true;
  try {
    await loadStateFromServer(props.date, true);
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadManga();
});

watch(
  () => props.date,
  () => {
    loadManga();
  }
);

const handleVersionChange = async (event: Event) => {
  const target = event.target as HTMLSelectElement;
  const id = Number(target?.value);
  if (id > 0) {
    isLoading.value = true;
    try {
      await loadVersionFromServer(id);
    } finally {
      isLoading.value = false;
    }
  }
};

const getGeneratedPanel = (panelNumber: number) => {
  return generatedPanels.value.find((p) => p.panelNumber === panelNumber);
};

const getGeneratedPanelImage = (panelNumber: number): string | undefined => {
  return getGeneratedPanel(panelNumber)?.image;
};

const getPanelPlan = (panelNumber: number): MangaPanelPlan | undefined => {
  return storyBreakdown.value?.panels.find((p) => p.panelNumber === panelNumber);
};

const getPanelNarrative = (panel: MangaPanelPlan): string => {
  return (panel.narrative || panel.sceneDescription || '').trim();
};

const getPanelLightboxCaption = (panel: MangaPanelPlan): string => {
  const narrative = getPanelNarrative(panel);
  return narrative ? `Panel ${panel.panelNumber} — ${narrative}` : `Panel ${panel.panelNumber}`;
};

const openLightbox = (url: string, caption: string) => {
  lightboxImage.value = { url, caption };
};

const closeLightbox = () => {
  lightboxImage.value = null;
};

const loadImageElement = (src: string, r2Key?: string): Promise<HTMLImageElement> => {
  const attemptLoad = (url: string) =>
    new Promise<HTMLImageElement>((resolve, reject) => {
      const img = new Image();
      img.crossOrigin = 'anonymous';
      img.onload = () => resolve(img);
      img.onerror = (err) => reject(err);
      img.src = url;
    });

  return attemptLoad(src).catch((err) => {
    const proxyUrl = getProxyAssetUrl(r2Key);
    if (proxyUrl && proxyUrl !== src) {
      return attemptLoad(proxyUrl);
    }
    throw err;
  });
};

const downloadA4Pdf = async () => {
  if (generatedPanels.value.length === 0) return;
  isExportingPdf.value = true;
  try {
    await exportMangaAsA4Pdf({
      date: props.date,
      formattedDate: props.formattedDate,
      style: style.value,
      storyAnalysis: storyAnalysis.value,
      moralLesson: moralLesson.value,
      panels: generatedPanels.value,
      getPanelPlan,
      getPanelNarrative,
      loadImageElement,
    });
  } finally {
    isExportingPdf.value = false;
  }
};
</script>

<template>
  <div class="space-y-6">
    <!-- Loading State -->
    <div
      v-if="isLoading && !storyAnalysis"
      class="bg-white/80 border border-[#D1C7BD] rounded-2xl py-16 flex flex-col items-center justify-center space-y-3"
    >
      <div class="w-10 h-10 border-4 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
      <p class="text-xs font-bold uppercase tracking-wider text-[#72383D]">
        Loading illustrated story...
      </p>
    </div>

    <template v-else-if="storyAnalysis">
      <!-- Story Header & Lesson Card -->
      <section class="bg-white/95 border-2 border-[#1E1B18] rounded-2xl p-5 shadow-[3px_3px_0_#1E1B18] space-y-4">
        <div class="flex flex-wrap items-start justify-between gap-2 border-b border-[#D1C7BD]/60 pb-3">
          <div>
            <div class="flex items-center gap-2 flex-wrap">
              <span class="font-comic font-bold text-xs uppercase tracking-[0.16em] text-[#72383D]">
                Custom Parable · Inspired by {{ formattedDate }} Readings
              </span>
              <span class="text-[10px] font-mono font-bold uppercase px-2 py-0.5 rounded-full bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">
                v{{ currentVersion }} · {{ style === 'manga' ? 'Manga' : 'Comic' }}
              </span>
            </div>
            <h3 class="text-2xl sm:text-3xl font-comic-display text-[#1E1B18] mt-1">
              {{ storyAnalysis.title }}
            </h3>
          </div>
          <div class="flex items-center gap-2 text-xs font-comic font-bold text-[#322D29]/80 bg-[#FAF6EE] px-3 py-1 rounded-full border border-[#1E1B18]/40">
            <span>📍 {{ storyAnalysis.setting.location }}</span>
            <span>·</span>
            <span>{{ storyAnalysis.setting.mood }}</span>
          </div>
        </div>

        <div
          v-if="moralLesson"
          class="p-3.5 rounded-xl bg-[#FFFDF5] border-2 border-[#1E1B18] border-l-4 border-l-[#72383D] text-sm sm:text-base text-[#1E1B18]"
        >
          <span class="font-comic font-bold uppercase tracking-wider text-xs text-[#72383D] block mb-0.5">
            Lesson Applied from Today's Readings
          </span>
          <p class="font-comic font-bold italic leading-relaxed">{{ moralLesson }}</p>
        </div>
      </section>

      <!-- Controls Bar -->
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center gap-2 flex-wrap">
          <!-- Version Switcher -->
          <div v-if="availableVersions.length > 1" class="inline-flex items-center gap-1.5">
            <label for="daily-tab-version-select" class="text-[11px] font-bold uppercase tracking-wider text-[#72383D]">
              Version
            </label>
            <select
              id="daily-tab-version-select"
              :value="currentStoryId || ''"
              @change="handleVersionChange"
              class="px-2.5 py-1.5 rounded-xl border border-[#D1C7BD] bg-white text-xs font-bold text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/30 cursor-pointer"
            >
              <option v-for="(ver, vIdx) in availableVersions" :key="ver.id" :value="ver.id">
                v{{ ver.version }}{{ vIdx === 0 ? ' (Latest)' : '' }} — {{ ver.title }}
              </option>
            </select>
          </div>

          <!-- Layout Mode Switcher -->
          <div class="inline-flex rounded-xl border border-[#D1C7BD] bg-white p-0.5">
            <button
              type="button"
              @click="layoutMode = 'grid'"
              class="px-3 py-1.5 rounded-lg text-xs font-bold transition-colors cursor-pointer"
              :class="layoutMode === 'grid' ? 'bg-[#72383D] text-white' : 'text-[#322D29]/75 hover:text-[#72383D]'"
            >
              🎞️ Comic Layout
            </button>
            <button
              type="button"
              @click="layoutMode = 'storybook'"
              class="px-3 py-1.5 rounded-lg text-xs font-bold transition-colors cursor-pointer"
              :class="layoutMode === 'storybook' ? 'bg-[#72383D] text-white' : 'text-[#322D29]/75 hover:text-[#72383D]'"
            >
              📖 Story + Panels
            </button>
          </div>
        </div>

        <div class="flex items-center gap-2 flex-wrap">
          <button
            v-if="isAdmin"
            type="button"
            @click="emit('open-generator')"
            class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[#72383D] text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer"
          >
            <span>✨</span>
            <span>Edit / Regenerate</span>
          </button>

          <button
            v-if="generatedPanels.length > 0"
            type="button"
            @click="downloadA4Pdf"
            :disabled="isExportingPdf"
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs cursor-pointer disabled:opacity-50"
          >
            <span>📄</span>
            <span>{{ isExportingPdf ? 'Creating A4 PDF...' : 'Download A4 PDF' }}</span>
          </button>
        </div>
      </div>

      <!-- Panels: Storybook Layout -->
      <div v-if="storyBreakdown?.panels?.length && layoutMode === 'storybook'" class="space-y-6">
        <article
          v-for="(panel, idx) in storyBreakdown.panels"
          :key="panel.panelNumber"
          class="grid grid-cols-1 md:grid-cols-12 gap-4 sm:gap-5 items-stretch"
        >
          <!-- Artwork Panel -->
          <div
            class="md:col-span-6 comic-panel-frame bg-white rounded-2xl overflow-hidden flex flex-col"
            :class="idx % 2 === 1 ? 'md:order-2' : 'md:order-1'"
          >
            <div class="px-4 py-2 bg-[#1E1B18] text-white flex items-center justify-between text-xs">
              <span class="font-comic-display text-sm sm:text-base tracking-wider uppercase">
                Panel {{ panel.panelNumber }}
              </span>
              <span class="text-white/75 text-[11px]">{{ panel.cameraAngle }} · {{ panel.visualMood }}</span>
            </div>
            <div
              class="relative w-full bg-[#F5F0E6] flex items-center justify-center overflow-hidden flex-1"
              :class="getTileAspectClass(panel.aspectRatio, panel.colSpan)"
            >
              <img
                v-if="getGeneratedPanelImage(panel.panelNumber)"
                :src="getGeneratedPanelImage(panel.panelNumber)"
                :alt="`Panel ${panel.panelNumber}`"
                class="w-full h-full object-contain cursor-pointer hover:scale-[1.01] transition-transform"
                @click="openLightbox(getGeneratedPanelImage(panel.panelNumber)!, getPanelLightboxCaption(panel))"
              />
            </div>
          </div>

          <!-- Narrative Panel -->
          <div
            class="md:col-span-6 comic-narrative-panel rounded-2xl overflow-hidden flex flex-col justify-between"
            :class="idx % 2 === 1 ? 'md:order-1' : 'md:order-2'"
          >
            <div>
              <div class="px-4 py-2 bg-[#72383D] text-white flex items-center justify-between gap-2">
                <span class="font-comic-display text-sm sm:text-base tracking-wider uppercase">
                  Story · Part {{ panel.panelNumber }} of {{ storyBreakdown.panels.length }}
                </span>
                <span
                  v-if="panel.characters?.length"
                  class="text-xs font-comic font-bold text-white/90 truncate max-w-[180px]"
                  :title="panel.characters.join(', ')"
                >
                  {{ panel.characters.join(', ') }}
                </span>
              </div>

              <div class="p-5 space-y-4">
                <p class="font-comic font-bold text-lg sm:text-xl leading-relaxed text-[#1E1B18]">
                  {{ getPanelNarrative(panel) }}
                </p>

                <div
                  v-if="panel.dialogue"
                  class="p-3.5 rounded-2xl bg-white border-2 border-[#1E1B18] text-base sm:text-lg font-comic font-bold italic text-[#1E1B18] shadow-[2px_2px_0_#1E1B18]"
                >
                  💬 “{{ panel.dialogue }}”
                </div>
              </div>
            </div>
          </div>
        </article>
      </div>

      <!-- Panels: Common Comic Page Layout (6-column variable-tile rows) -->
      <div v-else-if="storyBreakdown?.panels?.length" class="space-y-5">
        <div
          v-for="row in comicLayoutRows"
          :key="row.rowNumber"
          class="grid grid-cols-1 md:grid-cols-6 gap-4 items-stretch"
        >
          <div
            v-for="tile in row.tiles"
            :key="tile.panel.panelNumber"
            :class="[getTileColSpanClass(tile.colSpan), 'flex flex-col']"
          >
            <div class="comic-panel-frame bg-white rounded-2xl overflow-hidden flex flex-col h-full">
              <div class="px-3.5 py-2 bg-[#1E1B18] text-white flex items-center justify-between gap-2 text-xs">
                <span class="font-comic-display text-sm tracking-wider uppercase">
                  Panel {{ tile.panel.panelNumber }}
                </span>
                <span class="text-white/75 text-[10px] font-mono">
                  {{ tile.aspectRatio }} · {{ tile.panel.cameraAngle }}
                </span>
              </div>

              <div
                class="relative w-full bg-[#F5F0E6] flex items-center justify-center overflow-hidden"
                :class="getTileAspectClass(tile.aspectRatio, tile.colSpan)"
              >
                <img
                  v-if="getGeneratedPanelImage(tile.panel.panelNumber)"
                  :src="getGeneratedPanelImage(tile.panel.panelNumber)"
                  :alt="`Panel ${tile.panel.panelNumber}`"
                  class="w-full h-full object-cover cursor-pointer hover:scale-[1.01] transition-transform"
                  @click="openLightbox(getGeneratedPanelImage(tile.panel.panelNumber)!, getPanelLightboxCaption(tile.panel))"
                />
              </div>

              <div class="p-3.5 sm:p-4 bg-[#FFFDF5] border-t-2 border-[#1E1B18] flex-1 flex flex-col gap-2">
                <p class="font-comic font-bold text-sm sm:text-base text-[#1E1B18] leading-snug">
                  {{ getPanelNarrative(tile.panel) }}
                </p>
                <div
                  v-if="tile.panel.dialogue"
                  class="px-3 py-2 rounded-xl bg-white border-2 border-[#1E1B18] text-xs sm:text-sm font-comic font-bold italic text-[#1E1B18] shadow-[2px_2px_0_#1E1B18]"
                >
                  💬 “{{ tile.panel.dialogue }}”
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Fullscreen Image Lightbox -->
    <Teleport to="body">
      <div
        v-if="lightboxImage"
        class="fixed inset-0 z-60 bg-black/85 flex flex-col items-center justify-center p-4"
        @click="closeLightbox"
      >
        <img
          :src="lightboxImage.url"
          :alt="lightboxImage.caption"
          class="max-w-full max-h-[82vh] rounded-xl border-2 border-white/20 shadow-2xl object-contain"
          @click.stop
        />
        <p class="text-white text-base sm:text-lg font-comic font-bold mt-3 text-center max-w-2xl leading-relaxed">
          {{ lightboxImage.caption }}
        </p>
      </div>
    </Teleport>
  </div>
</template>
