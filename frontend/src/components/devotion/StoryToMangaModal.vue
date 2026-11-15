<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import {
  useStoryToManga,
  type ReadingsPayload,
  type ComicStyle,
} from '../../composables/useStoryToManga';

const props = defineProps<{
  isOpen: boolean;
  date: string;
  formattedDate: string;
  readings: ReadingsPayload | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const {
  style,
  imageModel,
  availableImageModels,
  generatedStory,
  moralLesson,
  storyAnalysis,
  characterReferences,
  storyBreakdown,
  generatedPanels,
  isGenerating,
  isAnalyzing,
  isGeneratingCharacters,
  isChunking,
  isGeneratingPanels,
  currentStepText,
  error,
  regeneratingCharacters,
  regeneratingPanels,
  loadConfig,
  updateImageModel,
  updateStyle,
  generateStoryAndManga,
  regenerateCharacter,
  regeneratePanel,
} = useStoryToManga();

const showSettingsBar = ref(false);
const lightboxImage = ref<{ url: string; caption: string } | null>(null);
const isDownloadingComposite = ref(false);

onMounted(() => {
  loadConfig();
});

// One-click automatic trigger when opened
watch(
  () => props.isOpen,
  async (open) => {
    if (open && props.readings) {
      await loadConfig();
      if (!isGenerating.value) {
        await generateStoryAndManga(props.date, props.readings, false);
      }
    }
  }
);

const handleModelChange = async (event: Event) => {
  const target = event.target as HTMLSelectElement;
  if (target?.value) {
    await updateImageModel(target.value);
  }
};

const handleStyleChange = (newStyle: ComicStyle) => {
  updateStyle(newStyle);
};

const handleRerunAll = async () => {
  if (!props.readings || isGenerating.value) return;
  await generateStoryAndManga(props.date, props.readings, true);
};

const getGeneratedPanelImage = (panelNumber: number): string | undefined => {
  return generatedPanels.value.find((p) => p.panelNumber === panelNumber)?.image;
};

const openLightbox = (url: string, caption: string) => {
  lightboxImage.value = { url, caption };
};

const closeLightbox = () => {
  lightboxImage.value = null;
};

const downloadDataUrl = (dataUrl: string, filename: string) => {
  const link = document.createElement('a');
  link.href = dataUrl;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

const loadImageElement = (src: string): Promise<HTMLImageElement> => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = 'anonymous';
    img.onload = () => resolve(img);
    img.onerror = (err) => reject(err);
    img.src = src;
  });
};

// Shareable composite comic page download (ported from ShareableComicLayout in victorhuangwq/story-to-manga)
const downloadCompositeComicPage = async () => {
  if (generatedPanels.value.length === 0) return;
  isDownloadingComposite.value = true;
  try {
    const panels = [...generatedPanels.value].sort((a, b) => a.panelNumber - b.panelNumber);
    const cols = panels.length <= 2 ? 1 : 2;
    const rows = Math.ceil(panels.length / cols);
    const cellW = cols === 1 ? 1040 : 510;
    const cellH = 440;
    const pad = 40;
    const headerH = 150;
    const footerH = 60;
    const canvasW = 1120;
    const canvasH = headerH + rows * (cellH + 24) + footerH;

    const canvas = document.createElement('canvas');
    canvas.width = canvasW;
    canvas.height = canvasH;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Background
    ctx.fillStyle = '#FBF8F3';
    ctx.fillRect(0, 0, canvasW, canvasH);

    // Header Title
    ctx.fillStyle = '#72383D';
    ctx.font = 'bold 30px Georgia, serif';
    ctx.textAlign = 'center';
    const title = storyAnalysis.value?.title || `Daily Readings ${style.value === 'manga' ? 'Manga' : 'Comic'}`;
    ctx.fillText(title, canvasW / 2, 56);

    ctx.fillStyle = '#5C5149';
    ctx.font = '16px Georgia, serif';
    ctx.fillText(`${props.formattedDate} · Generated with ${imageModel.value}`, canvasW / 2, 88);

    if (moralLesson.value) {
      ctx.fillStyle = '#322D29';
      ctx.font = 'italic 15px Georgia, serif';
      const maxLesson = moralLesson.value.length > 125 ? moralLesson.value.slice(0, 122) + '...' : moralLesson.value;
      ctx.fillText(`Lesson: "${maxLesson}"`, canvasW / 2, 118);
    }

    // Draw panels
    for (let i = 0; i < panels.length; i++) {
      const p = panels[i];
      if (!p || !p.image) continue;
      const col = i % cols;
      const row = Math.floor(i / cols);
      const x = pad + col * (cellW + 20);
      const y = headerH + row * (cellH + 24);

      ctx.fillStyle = '#FFFFFF';
      ctx.strokeStyle = '#322D29';
      ctx.lineWidth = 3;
      ctx.fillRect(x, y, cellW, cellH);
      ctx.strokeRect(x, y, cellW, cellH);

      try {
        const img = await loadImageElement(p.image);
        const scale = Math.min((cellW - 12) / img.width, (cellH - 12) / img.height);
        const drawW = img.width * scale;
        const drawH = img.height * scale;
        const drawX = x + (cellW - drawW) / 2;
        const drawY = y + (cellH - drawH) / 2;
        ctx.drawImage(img, drawX, drawY, drawW, drawH);
      } catch {
        // Skip broken image
      }

      // Panel badge
      ctx.fillStyle = '#322D29';
      ctx.fillRect(x + 10, y + 10, 82, 26);
      ctx.fillStyle = '#FFFFFF';
      ctx.font = 'bold 13px sans-serif';
      ctx.textAlign = 'center';
      ctx.fillText(`Panel ${p.panelNumber}`, x + 51, y + 28);
    }

    // Footer
    ctx.fillStyle = '#72383D';
    ctx.font = '13px sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText('Catholic Digital Devotional · Story-to-Manga Edition', canvasW / 2, canvasH - 22);

    const compositeUrl = canvas.toDataURL('image/png');
    downloadDataUrl(compositeUrl, `Devotion_Story_${style.value}_${props.date}.png`);
  } finally {
    isDownloadingComposite.value = false;
  }
};

const stepStates = computed(() => {
  const s1Done = !!storyAnalysis.value;
  const s2Done =
    !!storyAnalysis.value &&
    characterReferences.value.length >= (storyAnalysis.value.characters?.length || 1);
  const s3Done = !!storyBreakdown.value && (storyBreakdown.value.panels?.length || 0) > 0;
  const s4Done =
    s3Done &&
    generatedPanels.value.length >= (storyBreakdown.value?.panels?.length || 1);

  return [
    {
      id: 'story',
      label: '1. Story & Lesson',
      active: isAnalyzing.value,
      done: s1Done,
    },
    {
      id: 'characters',
      label: '2. Character Sheets',
      active: isGeneratingCharacters.value,
      done: s2Done,
    },
    {
      id: 'layout',
      label: '3. Panel Plan',
      active: isChunking.value,
      done: s3Done,
    },
    {
      id: 'panels',
      label: '4. Manga Panels',
      active: isGeneratingPanels.value,
      done: s4Done,
    },
  ];
});
</script>

<template>
  <Teleport to="body">
    <transition name="fade">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-black/65 backdrop-blur-xs overflow-y-auto"
        @click.self="emit('close')"
      >
        <div
          class="bg-[#FBF8F3] border-2 border-[#D1C7BD] rounded-2xl max-w-5xl w-full max-h-[92vh] flex flex-col shadow-2xl overflow-hidden text-[#322D29] my-auto"
        >
          <!-- Modal Header -->
          <div class="px-5 py-4 bg-gradient-to-r from-[#72383D] to-[#322D29] text-white flex flex-wrap items-center justify-between gap-3 shrink-0">
            <div class="flex items-center gap-3">
              <span class="text-2xl">🎨</span>
              <div>
                <div class="flex items-center gap-2 flex-wrap">
                  <h2 class="text-lg sm:text-xl font-serif font-bold tracking-wide">
                    Daily Readings Story-to-Manga
                  </h2>
                  <span class="text-[10px] font-mono uppercase px-2 py-0.5 rounded-full bg-white/15 border border-white/25">
                    {{ imageModel }}
                  </span>
                </div>
                <p class="text-xs text-white/80">
                  {{ formattedDate }} · Problem → Listening to Readings → Applying Lesson → Resolution
                </p>
              </div>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
              <!-- Toggle Model & Style Settings -->
              <button
                type="button"
                @click="showSettingsBar = !showSettingsBar"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-white/15 hover:bg-white/25 text-white text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer"
                title="Choose Image Generation Model & Style"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="3"></circle>
                  <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                </svg>
                <span>Model Settings</span>
              </button>

              <!-- One-Click Regenerate -->
              <button
                type="button"
                @click="handleRerunAll"
                :disabled="isGenerating"
                class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-xl bg-[#EFE9E1] hover:bg-white text-[#72383D] text-xs font-bold uppercase tracking-wider transition-colors disabled:opacity-50 cursor-pointer"
              >
                <span>✨</span>
                <span>{{ isGenerating ? 'Generating...' : 'Regenerate Story' }}</span>
              </button>

              <!-- Close Button -->
              <button
                type="button"
                @click="emit('close')"
                class="p-1.5 rounded-lg text-white/80 hover:text-white hover:bg-white/10 transition-colors cursor-pointer"
                aria-label="Close modal"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>

          <!-- Collapsible Image Generation Model & Style Settings Bar -->
          <div
            v-if="showSettingsBar"
            class="px-5 py-3.5 bg-[#EFE9E1] border-b border-[#D1C7BD] flex flex-col sm:flex-row sm:items-center justify-between gap-4 shrink-0"
          >
            <div class="flex flex-col sm:flex-row sm:items-center gap-3 flex-1">
              <div class="flex-1 max-w-xs">
                <label for="manga-image-model-select" class="block text-[10px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                  Image Generation Model
                </label>
                <select
                  id="manga-image-model-select"
                  :value="imageModel"
                  @change="handleModelChange"
                  :disabled="isGenerating"
                  class="w-full px-3 py-1.5 rounded-xl border border-[#D1C7BD] bg-white text-xs font-mono font-semibold text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/30"
                >
                  <option v-for="m in availableImageModels" :key="m" :value="m">
                    {{ m }}
                  </option>
                </select>
              </div>

              <div>
                <span class="block text-[10px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                  Visual Art Style
                </span>
                <div class="inline-flex rounded-xl border border-[#D1C7BD] bg-white p-0.5">
                  <button
                    type="button"
                    @click="handleStyleChange('manga')"
                    :disabled="isGenerating"
                    class="px-3 py-1 rounded-lg text-xs font-bold transition-colors cursor-pointer"
                    :class="style === 'manga' ? 'bg-[#72383D] text-white' : 'text-[#322D29]/75 hover:text-[#72383D]'"
                  >
                    🇯🇵 Manga (B&amp;W)
                  </button>
                  <button
                    type="button"
                    @click="handleStyleChange('comic')"
                    :disabled="isGenerating"
                    class="px-3 py-1 rounded-lg text-xs font-bold transition-colors cursor-pointer"
                    :class="style === 'comic' ? 'bg-[#72383D] text-white' : 'text-[#322D29]/75 hover:text-[#72383D]'"
                  >
                    🦸 Comic (Color)
                  </button>
                </div>
              </div>
            </div>

            <div class="text-[11px] text-[#322D29]/70">
              Default model: <code class="font-mono font-semibold text-[#72383D]">gemini-nano-banana-2.1</code>
            </div>
          </div>

          <!-- 4-Step Pipeline Progress Bar -->
          <div class="px-5 py-3 bg-[#FAF6EE] border-b border-[#D1C7BD]/70 shrink-0">
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
              <div
                v-for="step in stepStates"
                :key="step.id"
                class="flex items-center gap-2 px-3 py-1.5 rounded-xl border text-xs font-semibold transition-all"
                :class="
                  step.active
                    ? 'bg-[#72383D]/10 border-[#72383D] text-[#72383D]'
                    : step.done
                      ? 'bg-emerald-50 border-emerald-200 text-emerald-800'
                      : 'bg-white/60 border-[#D1C7BD]/60 text-[#322D29]/50'
                "
              >
                <span
                  v-if="step.active"
                  class="w-3.5 h-3.5 border-2 border-[#72383D] border-t-transparent rounded-full animate-spin shrink-0"
                ></span>
                <span v-else-if="step.done" class="text-emerald-600 font-bold shrink-0">✓</span>
                <span v-else class="w-2 h-2 rounded-full bg-[#D1C7BD] shrink-0"></span>
                <span class="truncate">{{ step.label }}</span>
              </div>
            </div>
            <p v-if="currentStepText && isGenerating" class="text-xs text-[#72383D] font-medium mt-2 animate-pulse">
              {{ currentStepText }}
            </p>
          </div>

          <!-- Modal Body (Scrollable) -->
          <div class="p-5 sm:p-6 overflow-y-auto space-y-6 flex-1">
            <!-- Error Alert -->
            <div
              v-if="error"
              class="p-4 rounded-2xl bg-red-50 border border-red-200 text-red-900 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3"
            >
              <div>
                <h4 class="font-bold text-sm">Story-to-Manga Generation Notice</h4>
                <p class="text-xs text-red-700 mt-0.5">{{ error }}</p>
              </div>
              <button
                type="button"
                @click="handleRerunAll"
                class="px-4 py-2 rounded-xl bg-red-800 hover:bg-red-900 text-white text-xs font-bold uppercase tracking-wider shrink-0 cursor-pointer"
              >
                Retry Generation
              </button>
            </div>

            <!-- Initial Loading State before Step 1 completes -->
            <div
              v-if="isAnalyzing && !storyAnalysis"
              class="py-16 flex flex-col items-center justify-center text-center space-y-3"
            >
              <div class="w-12 h-12 border-4 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
              <h3 class="font-serif text-lg font-bold text-[#72383D]">
                Writing Custom Story from Today's Readings...
              </h3>
              <p class="text-xs text-[#322D29]/70 max-w-md">
                Crafting a story where a protagonist faces a daily problem, listens to today's Mass readings, applies the lesson, and finds resolution.
              </p>
            </div>

            <!-- Stage 1: Generated Custom Story & Lesson Card -->
            <section
              v-if="storyAnalysis"
              class="bg-white/85 border border-[#D1C7BD] rounded-2xl p-5 shadow-xs space-y-4"
            >
              <div class="flex flex-wrap items-start justify-between gap-2 border-b border-[#D1C7BD]/60 pb-3">
                <div>
                  <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D] block">
                    Custom Parable · Inspired by {{ formattedDate }} Readings
                  </span>
                  <h3 class="text-xl sm:text-2xl font-serif font-bold text-[#322D29] mt-0.5">
                    {{ storyAnalysis.title }}
                  </h3>
                </div>
                <div class="flex items-center gap-2 text-[11px] text-[#322D29]/70 bg-[#FAF6EE] px-3 py-1 rounded-full border border-[#D1C7BD]">
                  <span>📍 {{ storyAnalysis.setting.location }}</span>
                  <span>·</span>
                  <span>{{ storyAnalysis.setting.mood }}</span>
                </div>
              </div>

              <!-- Moral Lesson Highlight -->
              <div
                v-if="moralLesson"
                class="p-3.5 rounded-xl bg-[#72383D]/8 border-l-4 border-[#72383D] text-xs sm:text-sm text-[#322D29]"
              >
                <span class="font-bold uppercase tracking-wider text-[10px] text-[#72383D] block mb-0.5">
                  Lesson Applied from Today's Readings
                </span>
                <p class="font-serif italic">{{ moralLesson }}</p>
              </div>

              <!-- Narrative Text -->
              <div class="prose max-w-none text-sm sm:text-base font-serif leading-relaxed text-[#322D29] whitespace-pre-line">
                {{ generatedStory }}
              </div>
            </section>

            <!-- Stage 2: Character Reference Sheets -->
            <section v-if="storyAnalysis && storyAnalysis.characters?.length" class="space-y-3">
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="text-base sm:text-lg font-serif font-bold text-[#72383D]">
                    Character Reference Sheets
                  </h3>
                  <p class="text-xs text-[#322D29]/65">
                    Generated with <code class="font-mono text-[#72383D]">{{ imageModel }}</code> to keep characters consistent across every panel
                  </p>
                </div>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
                <div
                  v-for="(char, idx) in storyAnalysis.characters"
                  :key="char.name"
                  class="bg-white/90 border border-[#D1C7BD] rounded-2xl p-4 flex flex-col justify-between shadow-xs"
                >
                  <div>
                    <!-- Reference Image or Loading Placeholder -->
                    <div class="relative aspect-square w-full rounded-xl bg-[#EFE9E1] border border-[#D1C7BD] overflow-hidden mb-3 flex items-center justify-center">
                      <img
                        v-if="characterReferences[idx]?.image"
                        :src="characterReferences[idx].image"
                        :alt="char.name"
                        class="w-full h-full object-cover cursor-pointer hover:scale-105 transition-transform"
                        @click="openLightbox(characterReferences[idx].image, `${char.name} — Character Reference`)"
                      />
                      <div
                        v-else-if="isGeneratingCharacters"
                        class="flex flex-col items-center justify-center p-4 text-center space-y-2"
                      >
                        <div class="w-8 h-8 border-3 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
                        <span class="text-[11px] text-[#72383D] font-semibold">Designing {{ char.name }}...</span>
                      </div>
                      <div v-else class="text-xs text-[#322D29]/50 p-4 text-center">
                        Awaiting character sheet
                      </div>

                      <!-- Regenerating Overlay -->
                      <div
                        v-if="regeneratingCharacters.has(idx)"
                        class="absolute inset-0 bg-black/55 flex flex-col items-center justify-center text-white space-y-2"
                      >
                        <div class="w-7 h-7 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                        <span class="text-xs font-semibold">Regenerating...</span>
                      </div>
                    </div>

                    <div class="flex items-center justify-between gap-2 mb-1">
                      <h4 class="font-serif font-bold text-base text-[#322D29]">{{ char.name }}</h4>
                      <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-[#EFE9E1] text-[#72383D]">
                        {{ char.role }}
                      </span>
                    </div>
                    <p class="text-xs text-[#322D29]/75 line-clamp-3">
                      {{ char.physicalDescription }}
                    </p>
                  </div>

                  <div v-if="characterReferences[idx]?.image" class="flex items-center gap-2 mt-3 pt-2.5 border-t border-[#D1C7BD]/50">
                    <button
                      type="button"
                      @click="downloadDataUrl(characterReferences[idx].image, `${char.name}_reference.jpg`)"
                      class="flex-1 px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FAF6EE] hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] transition-colors cursor-pointer"
                    >
                      Download
                    </button>
                    <button
                      type="button"
                      @click="regenerateCharacter(idx, date)"
                      :disabled="regeneratingCharacters.has(idx) || isGenerating"
                      class="flex-1 px-2.5 py-1.5 rounded-lg bg-[#72383D]/10 hover:bg-[#72383D]/20 text-[11px] font-bold text-[#72383D] transition-colors disabled:opacity-50 cursor-pointer"
                    >
                      Regenerate
                    </button>
                  </div>
                </div>
              </div>
            </section>

            <!-- Stage 3 & 4: Sequential Manga / Comic Panels -->
            <section v-if="storyBreakdown?.panels?.length" class="space-y-4">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <h3 class="text-base sm:text-lg font-serif font-bold text-[#72383D]">
                    Illustrated {{ style === 'manga' ? 'Manga' : 'Comic' }} Panels ({{ generatedPanels.length }}/{{ storyBreakdown.panels.length }})
                  </h3>
                  <p class="text-xs text-[#322D29]/65">
                    Click any panel image to view full size, download, or regenerate
                  </p>
                </div>

                <button
                  v-if="generatedPanels.length > 0"
                  type="button"
                  @click="downloadCompositeComicPage"
                  :disabled="isDownloadingComposite"
                  class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs cursor-pointer disabled:opacity-50"
                >
                  <span>📥</span>
                  <span>{{ isDownloadingComposite ? 'Creating Page...' : 'Download Full Comic Page' }}</span>
                </button>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                <div
                  v-for="panel in storyBreakdown.panels"
                  :key="panel.panelNumber"
                  class="bg-white border-2 border-[#322D29] rounded-2xl overflow-hidden shadow-sm flex flex-col justify-between"
                >
                  <div>
                    <!-- Panel Header Bar -->
                    <div class="px-4 py-2 bg-[#322D29] text-white flex items-center justify-between text-xs">
                      <span class="font-bold uppercase tracking-wider">Panel {{ panel.panelNumber }}</span>
                      <span class="text-white/75 text-[11px]">{{ panel.cameraAngle }} · {{ panel.visualMood }}</span>
                    </div>

                    <!-- Generated Panel Image or Storyboard Plan Preview -->
                    <div class="relative aspect-square w-full bg-[#F5F0E6] flex items-center justify-center overflow-hidden">
                      <img
                        v-if="getGeneratedPanelImage(panel.panelNumber)"
                        :src="getGeneratedPanelImage(panel.panelNumber)"
                        :alt="`Panel ${panel.panelNumber}`"
                        class="w-full h-full object-contain cursor-pointer hover:scale-[1.01] transition-transform"
                        @click="openLightbox(getGeneratedPanelImage(panel.panelNumber)!, `Panel ${panel.panelNumber}: ${panel.sceneDescription}`)"
                      />
                      <div
                        v-else
                        class="p-6 text-center space-y-3 max-w-sm"
                      >
                        <div
                          v-if="isGeneratingPanels"
                          class="w-10 h-10 mx-auto border-3 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"
                        ></div>
                        <p class="text-xs font-semibold text-[#72383D] uppercase tracking-wider">
                          {{ isGeneratingPanels ? `Illustrating Panel ${panel.panelNumber}...` : 'Storyboard Layout' }}
                        </p>
                        <p class="text-xs text-[#322D29]/80 leading-relaxed">
                          {{ panel.sceneDescription }}
                        </p>
                      </div>

                      <!-- Regenerating Overlay -->
                      <div
                        v-if="regeneratingPanels.has(panel.panelNumber)"
                        class="absolute inset-0 bg-black/55 flex flex-col items-center justify-center text-white space-y-2"
                      >
                        <div class="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                        <span class="text-xs font-semibold">Regenerating Panel {{ panel.panelNumber }}...</span>
                      </div>
                    </div>

                    <!-- Scene & Speech Bubble Caption -->
                    <div class="p-4 space-y-2.5 bg-[#FBF8F3] border-t border-[#D1C7BD]">
                      <div
                        v-if="panel.dialogue"
                        class="p-3 rounded-xl bg-white border-2 border-[#322D29] text-xs sm:text-sm font-serif text-[#322D29] shadow-2xs"
                      >
                        💬 "{{ panel.dialogue }}"
                      </div>
                      <p class="text-xs text-[#322D29]/75 leading-relaxed">
                        {{ panel.sceneDescription }}
                      </p>
                    </div>
                  </div>

                  <!-- Panel Actions -->
                  <div
                    v-if="getGeneratedPanelImage(panel.panelNumber)"
                    class="px-4 py-2.5 bg-[#FAF6EE] border-t border-[#D1C7BD]/70 flex items-center justify-end gap-2"
                  >
                    <button
                      type="button"
                      @click="downloadDataUrl(getGeneratedPanelImage(panel.panelNumber)!, `panel_${panel.panelNumber}_${date}.jpg`)"
                      class="px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors cursor-pointer"
                    >
                      Download Panel
                    </button>
                    <button
                      type="button"
                      @click="regeneratePanel(panel.panelNumber, date)"
                      :disabled="regeneratingPanels.has(panel.panelNumber) || isGenerating"
                      class="px-3 py-1.5 rounded-lg bg-[#72383D] hover:bg-[#542529] text-white text-xs font-bold transition-colors disabled:opacity-50 cursor-pointer"
                    >
                      Regenerate
                    </button>
                  </div>
                </div>
              </div>
            </section>
          </div>
        </div>

        <!-- Fullscreen Image Lightbox -->
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
          <p class="text-white text-sm font-serif mt-3 text-center max-w-xl">
            {{ lightboxImage.caption }}
          </p>
        </div>
      </div>
    </transition>
  </Teleport>
</template>
