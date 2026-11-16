<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import {
  useStoryToManga,
  type ReadingsPayload,
  type ComicStyle,
  type MangaPanelPlan,
  type MangaCharacter,
  type MangaUploadedReference,
} from '../../composables/useStoryToManga';
import { useAccessControl } from '../../composables/useAccessControl';
import { exportMangaAsA4Pdf } from '../../utils/mangaPdfExport';
import {
  groupPanelsIntoComicRows,
  getTileAspectClass,
  getTileColSpanClass,
} from '../../utils/mangaLayout';

const props = defineProps<{
  isOpen: boolean;
  date: string;
  formattedDate: string;
  readings: ReadingsPayload | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const { isAdmin, syncAdminState } = useAccessControl();

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
  r2Stored,
  currentStoryId,
  currentVersion,
  availableVersions,
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
  loadVersionFromServer,
  updateImageModel,
  updateStyle,
  getProxyAssetUrl,
  generateStoryAndManga,
  regenerateCharacter,
  regeneratePanel,
} = useStoryToManga();

const showSettingsBar = ref(false);
const showCustomStoryForm = ref(false);
const customStoryText = ref('');
const customCharacters = ref<
  Array<MangaCharacter & { refImage?: string; refFileName?: string }>
>([]);
const layoutMode = ref<'storybook' | 'grid'>('grid');
const lightboxImage = ref<{ url: string; caption: string } | null>(null);
const isExportingPdf = ref(false);

const comicLayoutRows = computed(() =>
  groupPanelsIntoComicRows(storyBreakdown.value?.panels || [])
);

const customStoryWordCount = computed(() => {
  const trimmed = customStoryText.value.trim();
  return trimmed ? trimmed.split(/\s+/).length : 0;
});

const validCustomCharacters = computed(() =>
  customCharacters.value.filter(
    (c) => c.name.trim() !== '' && c.physicalDescription.trim() !== ''
  )
);

const canSubmitCustomForm = computed(() => {
  if (customStoryWordCount.value > 500) return false;
  return customStoryWordCount.value > 0 || validCustomCharacters.value.length > 0;
});

onMounted(() => {
  syncAdminState();
  loadConfig();
});

// Automatic load or one-click generation when opened
watch(
  () => props.isOpen,
  async (open) => {
    if (open && props.readings) {
      syncAdminState();
      await loadConfig();
      if (!isAdmin.value) {
        showSettingsBar.value = false;
        showCustomStoryForm.value = false;
      }
      if (!isGenerating.value) {
        await generateStoryAndManga(props.date, props.readings, false);
      }
    }
  }
);

const handleModelChange = async (event: Event) => {
  if (!isAdmin.value) return;
  const target = event.target as HTMLSelectElement;
  if (target?.value) {
    await updateImageModel(target.value);
  }
};

const handleStyleChange = (newStyle: ComicStyle) => {
  if (!isAdmin.value) return;
  updateStyle(newStyle);
};

const handleVersionChange = async (event: Event) => {
  const target = event.target as HTMLSelectElement;
  const id = Number(target?.value);
  if (id > 0 && !isGenerating.value) {
    await loadVersionFromServer(id);
  }
};

const handleRerunAll = async () => {
  if (!isAdmin.value || !props.readings || isGenerating.value) return;
  await generateStoryAndManga(props.date, props.readings, true);
};

const prefillCurrentStory = () => {
  if (generatedStory.value) {
    customStoryText.value = generatedStory.value;
  }
};

const addCustomCharacter = () => {
  if (customCharacters.value.length >= 4) return;
  customCharacters.value.push({
    name: '',
    role: customCharacters.value.length === 0 ? 'Protagonist' : 'Supporting',
    physicalDescription: '',
    personality: '',
  });
};

const removeCustomCharacter = (idx: number) => {
  customCharacters.value.splice(idx, 1);
};

const prefillCurrentCharacters = () => {
  if (!storyAnalysis.value?.characters?.length) return;
  customCharacters.value = storyAnalysis.value.characters.slice(0, 4).map((c) => ({
    name: c.name,
    role: c.role,
    physicalDescription: c.physicalDescription,
    personality: c.personality,
  }));
};

const handleCharacterRefUpload = (idx: number, event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input?.files?.[0];
  if (!file || !customCharacters.value[idx]) return;
  const reader = new FileReader();
  reader.onload = () => {
    const result = typeof reader.result === 'string' ? reader.result : '';
    if (result && customCharacters.value[idx]) {
      customCharacters.value[idx].refImage = result;
      customCharacters.value[idx].refFileName = file.name;
    }
  };
  reader.readAsDataURL(file);
};

const clearCharacterRefUpload = (idx: number) => {
  if (!customCharacters.value[idx]) return;
  customCharacters.value[idx].refImage = undefined;
  customCharacters.value[idx].refFileName = undefined;
};

const handleGenerateFromCustomStory = async () => {
  if (!isAdmin.value || isGenerating.value || !canSubmitCustomForm.value) return;
  const storyToUse = customStoryText.value.trim();
  const charsToUse: MangaCharacter[] = validCustomCharacters.value.map((c) => ({
    name: c.name.trim(),
    role: c.role.trim() || 'Character',
    physicalDescription: c.physicalDescription.trim(),
    personality: c.personality.trim() || 'Kind and relatable',
  }));
  const uploadedRefs: MangaUploadedReference[] = validCustomCharacters.value
    .filter((c) => c.refImage)
    .map((c, i) => ({
      id: `custom_char_${i}`,
      name: c.name.trim(),
      image: c.refImage!,
      fileName: c.refFileName,
    }));

  showCustomStoryForm.value = false;
  await generateStoryAndManga(
    props.date,
    props.readings || {},
    true,
    storyToUse || undefined,
    charsToUse.length > 0 ? charsToUse : undefined,
    uploadedRefs.length > 0 ? uploadedRefs : undefined
  );
};

const getGeneratedPanel = (panelNumber: number) => {
  return generatedPanels.value.find((p) => p.panelNumber === panelNumber);
};

const getGeneratedPanelImage = (panelNumber: number): string | undefined => {
  return getGeneratedPanel(panelNumber)?.image;
};

const openLightbox = (url: string, caption: string) => {
  lightboxImage.value = { url, caption };
};

const closeLightbox = () => {
  lightboxImage.value = null;
};

const downloadDataUrl = async (imageUrl: string, filename: string, r2Key?: string) => {
  if (!imageUrl) return;
  if (imageUrl.startsWith('data:') || imageUrl.startsWith('blob:')) {
    const link = document.createElement('a');
    link.href = imageUrl;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    return;
  }

  const tryFetchBlob = async (url: string) => {
    const resp = await fetch(url);
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
    return await resp.blob();
  };

  try {
    let blob: Blob;
    try {
      blob = await tryFetchBlob(imageUrl);
    } catch {
      const proxyUrl = getProxyAssetUrl(r2Key);
      if (!proxyUrl || proxyUrl === imageUrl) throw new Error('Download failed');
      blob = await tryFetchBlob(proxyUrl);
    }
    const objUrl = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = objUrl;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    setTimeout(() => URL.revokeObjectURL(objUrl), 1000);
  } catch {
    window.open(imageUrl, '_blank', 'noopener,noreferrer');
  }
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

const downloadMangaA4Pdf = async () => {
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
        class="manga-print-modal fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-black/65 backdrop-blur-xs overflow-y-auto"
        @click.self="emit('close')"
      >
        <div
          class="manga-print-dialog bg-[#FBF8F3] border-2 border-[#D1C7BD] rounded-2xl max-w-5xl w-full max-h-[92vh] flex flex-col shadow-2xl overflow-hidden text-[#322D29] my-auto"
        >
          <!-- Modal Header -->
          <div class="manga-no-print px-5 py-4 bg-gradient-to-r from-[#72383D] to-[#322D29] text-white flex flex-wrap items-center justify-between gap-3 shrink-0">
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
                  <span
                    v-if="storyAnalysis"
                    class="text-[10px] font-mono font-bold uppercase px-2 py-0.5 rounded-full bg-amber-300/20 text-amber-100 border border-amber-200/35"
                  >
                    v{{ currentVersion }}
                  </span>
                  <span
                    v-if="r2Stored"
                    class="text-[10px] font-bold uppercase px-2 py-0.5 rounded-full bg-emerald-400/20 text-emerald-100 border border-emerald-300/35"
                    title="Stored in Cloudflare R2 & PostgreSQL"
                  >
                    R2 Stored
                  </span>
                </div>
                <p class="text-xs text-white/80">
                  {{ formattedDate }} · Problem → Listening to Readings → Applying Lesson → Resolution
                </p>
              </div>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
              <!-- Toggle Model & Style Settings (Admin Only) -->
              <button
                v-if="isAdmin"
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

              <!-- Write Custom Story & Characters (Admin Only) -->
              <button
                v-if="isAdmin"
                type="button"
                @click="showCustomStoryForm = !showCustomStoryForm"
                :disabled="isGenerating"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-white/25 text-xs font-bold uppercase tracking-wider transition-colors disabled:opacity-50 cursor-pointer"
                :class="showCustomStoryForm ? 'bg-white text-[#72383D]' : 'bg-white/15 hover:bg-white/25 text-white'"
                title="Add your own story and/or characters"
              >
                <span>✍️</span>
                <span>Story &amp; Characters</span>
              </button>

              <!-- One-Click Regenerate (Admin Only) -->
              <button
                v-if="isAdmin"
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

          <!-- Collapsible Image Generation Model & Style Settings Bar (Admin Only) -->
          <div
            v-if="isAdmin && showSettingsBar"
            class="manga-no-print px-5 py-3.5 bg-[#EFE9E1] border-b border-[#D1C7BD] flex flex-col sm:flex-row sm:items-center justify-between gap-4 shrink-0"
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
          <div
            v-if="isAdmin || isGenerating || storyAnalysis"
            class="manga-no-print px-5 py-3 bg-[#FAF6EE] border-b border-[#D1C7BD]/70 shrink-0"
          >
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
          <div class="manga-print-body p-5 sm:p-6 overflow-y-auto space-y-6 flex-1">
            <!-- Custom Story & Characters Composer (Admin Only) -->
            <section
              v-if="isAdmin && showCustomStoryForm"
              class="manga-no-print bg-white border-2 border-[#72383D] rounded-2xl p-5 shadow-[3px_3px_0_#72383D] space-y-5"
            >
              <!-- Custom Story Textarea -->
              <div class="space-y-2.5">
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <div>
                    <h3 class="font-serif font-bold text-base sm:text-lg text-[#72383D]">
                      ✍️ Custom Story &amp; Characters
                    </h3>
                    <p class="text-xs text-[#322D29]/70">
                      Write or paste your own story (optional if you only want to customize characters for today’s readings) and/or add up to 4 characters.
                    </p>
                  </div>
                  <button
                    v-if="generatedStory"
                    type="button"
                    @click="prefillCurrentStory"
                    class="px-3 py-1 rounded-lg border border-[#D1C7BD] bg-[#FAF6EE] hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors cursor-pointer"
                  >
                    Use Current Story
                  </button>
                </div>

                <textarea
                  v-model="customStoryText"
                  rows="5"
                  placeholder="Write or paste your story here (up to 500 words)... Leave blank to generate a story from today's Mass readings using your custom characters below."
                  class="w-full rounded-xl border-2 border-[#D1C7BD] bg-[#FBF8F3] p-3.5 text-sm sm:text-base font-comic text-[#1E1B18] focus:outline-none focus:border-[#72383D] leading-relaxed"
                ></textarea>

                <div class="text-xs font-bold" :class="customStoryWordCount > 500 ? 'text-red-700' : 'text-[#322D29]/70'">
                  {{ customStoryWordCount }} / 500 words
                  <span v-if="customStoryWordCount > 500"> — Please shorten to 500 words or fewer.</span>
                </div>
              </div>

              <!-- Custom Characters Builder -->
              <div class="pt-4 border-t border-[#D1C7BD]/70 space-y-3">
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <div>
                    <h4 class="font-serif font-bold text-sm sm:text-base text-[#322D29]">
                      👤 Custom Characters ({{ customCharacters.length }}/4)
                    </h4>
                    <p class="text-xs text-[#322D29]/65">
                      Define character names, appearance, and optional reference photos. Leave empty to auto-extract characters from the story.
                    </p>
                  </div>
                  <div class="flex items-center gap-2">
                    <button
                      v-if="storyAnalysis?.characters?.length"
                      type="button"
                      @click="prefillCurrentCharacters"
                      class="px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FAF6EE] hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors cursor-pointer"
                    >
                      Load Current Characters
                    </button>
                    <button
                      type="button"
                      @click="addCustomCharacter"
                      :disabled="customCharacters.length >= 4"
                      class="px-3 py-1.5 rounded-lg bg-[#72383D]/10 hover:bg-[#72383D]/20 text-xs font-bold text-[#72383D] transition-colors disabled:opacity-40 cursor-pointer"
                    >
                      + Add Character
                    </button>
                  </div>
                </div>

                <div v-if="customCharacters.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-3.5">
                  <div
                    v-for="(char, idx) in customCharacters"
                    :key="idx"
                    class="p-3.5 rounded-xl border border-[#D1C7BD] bg-[#FAF6EE]/70 space-y-2.5"
                  >
                    <div class="flex items-center justify-between gap-2">
                      <span class="text-xs font-bold uppercase tracking-wider text-[#72383D]">
                        Character {{ idx + 1 }}
                      </span>
                      <button
                        type="button"
                        @click="removeCustomCharacter(idx)"
                        class="text-xs font-bold text-red-700 hover:text-red-900 cursor-pointer"
                      >
                        Remove
                      </button>
                    </div>

                    <div class="grid grid-cols-2 gap-2">
                      <input
                        v-model="char.name"
                        type="text"
                        placeholder="Name (e.g. Leo) *"
                        class="px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs sm:text-sm font-semibold text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                      />
                      <input
                        v-model="char.role"
                        type="text"
                        placeholder="Role (e.g. Protagonist)"
                        class="px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs sm:text-sm text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                      />
                    </div>

                    <textarea
                      v-model="char.physicalDescription"
                      rows="2"
                      placeholder="Physical appearance (age, hair, clothes, distinctive features) *"
                      class="w-full px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs sm:text-sm text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                    ></textarea>

                    <input
                      v-model="char.personality"
                      type="text"
                      placeholder="Personality traits (optional, e.g. Curious, gentle)"
                      class="w-full px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs sm:text-sm text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                    />

                    <!-- Optional Reference Image Upload -->
                    <div class="flex items-center justify-between gap-2 pt-1">
                      <label class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] cursor-pointer">
                        <span>📷 {{ char.refImage ? 'Change Photo' : 'Reference Photo (Optional)' }}</span>
                        <input
                          type="file"
                          accept="image/*"
                          class="hidden"
                          @change="handleCharacterRefUpload(idx, $event)"
                        />
                      </label>

                      <div v-if="char.refImage" class="flex items-center gap-2">
                        <img
                          :src="char.refImage"
                          :alt="char.name || 'Character reference'"
                          class="w-8 h-8 rounded-md object-cover border border-[#D1C7BD]"
                        />
                        <button
                          type="button"
                          @click="clearCharacterRefUpload(idx)"
                          class="text-[11px] font-bold text-red-700 hover:underline cursor-pointer"
                        >
                          Clear
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Submit Bar -->
              <div class="flex flex-wrap items-center justify-end gap-2 pt-2 border-t border-[#D1C7BD]/70">
                <button
                  type="button"
                  @click="showCustomStoryForm = false"
                  class="px-3.5 py-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors cursor-pointer"
                >
                  Cancel
                </button>
                <button
                  type="button"
                  @click="handleGenerateFromCustomStory"
                  :disabled="!canSubmitCustomForm || isGenerating"
                  class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs transition-colors disabled:opacity-50 cursor-pointer"
                >
                  <span>🎨</span>
                  <span>{{ isGenerating ? 'Generating...' : 'Generate Manga' }}</span>
                </button>
              </div>
            </section>

            <!-- Error Alert -->
            <div
              v-if="error"
              class="manga-no-print p-4 rounded-2xl bg-red-50 border border-red-200 text-red-900 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3"
            >
              <div>
                <h4 class="font-bold text-sm">Story-to-Manga Generation Notice</h4>
                <p class="text-xs text-red-700 mt-0.5">{{ error }}</p>
              </div>
              <button
                v-if="isAdmin"
                type="button"
                @click="handleRerunAll"
                class="px-4 py-2 rounded-xl bg-red-800 hover:bg-red-900 text-white text-xs font-bold uppercase tracking-wider shrink-0 cursor-pointer"
              >
                Retry Generation
              </button>
              <router-link
                v-else
                to="/admin"
                @click="emit('close')"
                class="px-4 py-2 rounded-xl bg-[#72383D] hover:bg-[#542529] text-white text-xs font-bold uppercase tracking-wider shrink-0"
              >
                Admin Login
              </router-link>
            </div>

            <!-- Initial Loading State before Step 1 completes -->
            <div
              v-if="isAnalyzing && !storyAnalysis"
              class="manga-no-print py-16 flex flex-col items-center justify-center text-center space-y-3"
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
              class="bg-white/95 border-2 border-[#1E1B18] rounded-2xl p-5 shadow-[3px_3px_0_#1E1B18] space-y-4"
            >
              <div class="flex flex-wrap items-start justify-between gap-2 border-b border-[#D1C7BD]/60 pb-3">
                <div>
                  <span class="font-comic font-bold text-xs uppercase tracking-[0.16em] text-[#72383D] block">
                    Custom Parable · Inspired by {{ formattedDate }} Readings
                  </span>
                  <h3 class="manga-print-title text-2xl sm:text-3xl font-comic-display text-[#1E1B18] mt-1">
                    {{ storyAnalysis.title }}
                  </h3>
                </div>
                <div class="flex items-center gap-2 text-xs font-comic font-bold text-[#322D29]/80 bg-[#FAF6EE] px-3 py-1 rounded-full border border-[#1E1B18]/40">
                  <span>📍 {{ storyAnalysis.setting.location }}</span>
                  <span>·</span>
                  <span>{{ storyAnalysis.setting.mood }}</span>
                </div>
              </div>

              <!-- Moral Lesson Highlight -->
              <div
                v-if="moralLesson"
                class="p-3.5 rounded-xl bg-[#FFFDF5] border-2 border-[#1E1B18] border-l-4 border-l-[#72383D] text-sm sm:text-base text-[#1E1B18]"
              >
                <span class="font-comic font-bold uppercase tracking-wider text-xs text-[#72383D] block mb-0.5">
                  Lesson Applied from Today's Readings
                </span>
                <p class="font-comic font-bold italic leading-relaxed">{{ moralLesson }}</p>
              </div>

              <!-- Unabridged Narrative Text (shown directly while planning panels, collapsible once Story + Panels flow is ready) -->
              <div
                v-if="!storyBreakdown?.panels?.length"
                class="prose max-w-none text-base sm:text-lg font-comic font-bold leading-relaxed text-[#1E1B18] whitespace-pre-line"
              >
                {{ generatedStory }}
              </div>
              <details
                v-else-if="generatedStory"
                class="manga-no-print rounded-xl border border-[#D1C7BD]/80 bg-[#FAF6EE]/70 px-4 py-2.5"
              >
                <summary class="cursor-pointer text-xs font-bold uppercase tracking-wider text-[#72383D] select-none">
                  Read Full Unabridged Story Text
                </summary>
                <div class="prose max-w-none text-base font-comic leading-relaxed text-[#1E1B18] whitespace-pre-line mt-3 pt-3 border-t border-[#D1C7BD]/60">
                  {{ generatedStory }}
                </div>
              </details>
            </section>

            <!-- Stage 2: Character Reference Sheets -->
            <section v-if="storyAnalysis && storyAnalysis.characters?.length" class="manga-no-print space-y-3">
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
                      @click="downloadDataUrl(characterReferences[idx].image, `${char.name}_reference.jpg`, characterReferences[idx].r2Key)"
                      class="flex-1 px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FAF6EE] hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] transition-colors cursor-pointer"
                    >
                      Download
                    </button>
                    <button
                      v-if="isAdmin"
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

            <!-- Stage 3 & 4: Separate Artwork & Narrative Comic Panels -->
            <section v-if="storyBreakdown?.panels?.length" class="space-y-5">
              <div class="manga-no-print flex flex-wrap items-center justify-between gap-3">
                <div>
                  <h3 class="text-base sm:text-lg font-serif font-bold text-[#72383D]">
                    Illustrated Story &amp; {{ style === 'manga' ? 'Manga' : 'Comic' }} Panels ({{ generatedPanels.length }}/{{ storyBreakdown.panels.length }})
                  </h3>
                  <p class="text-xs text-[#322D29]/65">
                    Each story beat pairs an artwork panel with a separate comic-lettered narrative panel. Click any panel to enlarge{{ isAdmin ? ', download, or regenerate' : ' or download' }}.
                  </p>
                </div>

                <div class="flex items-center gap-2 flex-wrap">
                  <!-- Version Selector -->
                  <div v-if="availableVersions.length > 1" class="inline-flex items-center gap-1.5">
                    <label for="manga-modal-version-select" class="text-[11px] font-bold uppercase tracking-wider text-[#72383D]">
                      Version
                    </label>
                    <select
                      id="manga-modal-version-select"
                      :value="currentStoryId || ''"
                      @change="handleVersionChange"
                      :disabled="isGenerating"
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

                  <button
                    v-if="generatedPanels.length > 0"
                    type="button"
                    @click="downloadMangaA4Pdf"
                    :disabled="isExportingPdf"
                    class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs cursor-pointer disabled:opacity-50"
                    title="Download A4 PDF"
                  >
                    <span>📄</span>
                    <span>{{ isExportingPdf ? 'Creating A4 PDF...' : 'Download A4 PDF' }}</span>
                  </button>
                </div>
              </div>

              <!-- Separate Artwork & Narrative Panels (Storybook Layout) -->
              <div v-if="layoutMode === 'storybook'" class="space-y-6">
                <article
                  v-for="(panel, idx) in storyBreakdown.panels"
                  :key="panel.panelNumber"
                  class="manga-print-row grid grid-cols-1 md:grid-cols-12 gap-4 sm:gap-5 items-stretch"
                >
                  <!-- 1) Standalone Artwork Panel -->
                  <div
                    class="md:col-span-6 manga-print-col-6 comic-panel-frame bg-white rounded-2xl overflow-hidden flex flex-col"
                    :class="idx % 2 === 1 ? 'md:order-2' : 'md:order-1'"
                  >
                    <div class="px-4 py-2 bg-[#1E1B18] text-white flex items-center justify-between text-xs">
                      <span class="font-comic-display text-sm sm:text-base tracking-wider uppercase">
                        Panel {{ panel.panelNumber }}
                      </span>
                      <span class="manga-no-print text-white/75 text-[11px]">{{ panel.cameraAngle }} · {{ panel.visualMood }}</span>
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
                  </div>

                  <!-- 2) Separate Narrative Text Panel -->
                  <div
                    class="md:col-span-6 manga-print-col-6 comic-narrative-panel rounded-2xl overflow-hidden flex flex-col justify-between"
                    :class="idx % 2 === 1 ? 'md:order-1' : 'md:order-2'"
                  >
                    <div>
                      <div class="px-4 py-2 bg-[#72383D] text-white flex items-center justify-between gap-2">
                        <span class="font-comic-display text-sm sm:text-base tracking-wider uppercase">
                          Story · Part {{ panel.panelNumber }} of {{ storyBreakdown.panels.length }}
                        </span>
                        <span
                          v-if="panel.characters?.length"
                          class="text-xs font-comic font-bold text-white/90 truncate max-w-[200px]"
                          :title="panel.characters.join(', ')"
                        >
                          {{ panel.characters.join(', ') }}
                        </span>
                      </div>

                      <div class="p-5 sm:p-6 space-y-4">
                        <p class="manga-print-narrative-text font-comic font-bold text-lg sm:text-xl leading-relaxed text-[#1E1B18]">
                          {{ getPanelNarrative(panel) }}
                        </p>

                        <div
                          v-if="panel.dialogue"
                          class="manga-print-dialogue-text p-3.5 rounded-2xl bg-white border-2 border-[#1E1B18] text-base sm:text-lg font-comic font-bold italic text-[#1E1B18] shadow-[2px_2px_0_#1E1B18]"
                        >
                          💬 “{{ panel.dialogue }}”
                        </div>

                        <details
                          v-if="panel.narrative && panel.sceneDescription"
                          class="manga-no-print text-xs text-[#322D29]/65"
                        >
                          <summary class="cursor-pointer font-semibold hover:text-[#72383D] select-none">
                            Visual scene direction
                          </summary>
                          <p class="mt-1.5 leading-relaxed pl-2 border-l border-[#D1C7BD]">
                            {{ panel.sceneDescription }}
                          </p>
                        </details>
                      </div>
                    </div>

                    <!-- Panel Actions -->
                    <div
                      v-if="getGeneratedPanelImage(panel.panelNumber)"
                      class="manga-no-print px-5 py-3 border-t border-[#D1C7BD]/70 bg-[#FAF6EE]/80 flex items-center justify-end gap-2"
                    >
                      <button
                        type="button"
                        @click="downloadDataUrl(getGeneratedPanelImage(panel.panelNumber)!, `panel_${panel.panelNumber}_${date}.jpg`, getGeneratedPanel(panel.panelNumber)?.r2Key)"
                        class="px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors cursor-pointer"
                      >
                        Download Panel
                      </button>
                      <button
                        v-if="isAdmin"
                        type="button"
                        @click="regeneratePanel(panel.panelNumber, date)"
                        :disabled="regeneratingPanels.has(panel.panelNumber) || isGenerating"
                        class="px-3 py-1.5 rounded-lg bg-[#72383D] hover:bg-[#542529] text-white text-xs font-bold transition-colors disabled:opacity-50 cursor-pointer"
                      >
                        Regenerate
                      </button>
                    </div>
                  </div>
                </article>
              </div>

              <!-- Common Comic Page Layout (6-column variable-tile rows) -->
              <div v-else class="space-y-5">
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
                        <span class="manga-no-print text-white/75 text-[10px] font-mono">
                          {{ tile.aspectRatio }} · {{ tile.panel.cameraAngle }}
                        </span>
                      </div>

                      <!-- Tile Artwork matching exact tile aspect ratio -->
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
                        <div
                          v-else
                          class="p-5 text-center space-y-2 max-w-sm"
                        >
                          <div
                            v-if="isGeneratingPanels"
                            class="w-9 h-9 mx-auto border-3 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"
                          ></div>
                          <p class="text-xs font-semibold text-[#72383D] uppercase tracking-wider">
                            {{ isGeneratingPanels ? `Illustrating Panel ${tile.panel.panelNumber} (${tile.aspectRatio})...` : `Tile ${tile.aspectRatio}` }}
                          </p>
                          <p class="text-xs text-[#322D29]/80 leading-relaxed line-clamp-4">
                            {{ tile.panel.sceneDescription }}
                          </p>
                        </div>

                        <!-- Regenerating Overlay -->
                        <div
                          v-if="regeneratingPanels.has(tile.panel.panelNumber)"
                          class="absolute inset-0 bg-black/55 flex flex-col items-center justify-center text-white space-y-2"
                        >
                          <div class="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                          <span class="text-xs font-semibold">Regenerating Panel {{ tile.panel.panelNumber }}...</span>
                        </div>
                      </div>

                      <!-- Integrated Comic Caption & Dialogue Strip -->
                      <div class="p-3.5 sm:p-4 bg-[#FFFDF5] border-t-2 border-[#1E1B18] flex-1 flex flex-col justify-between gap-2.5">
                        <div class="space-y-2">
                          <p class="manga-print-narrative-text font-comic font-bold text-sm sm:text-base text-[#1E1B18] leading-snug">
                            {{ getPanelNarrative(tile.panel) }}
                          </p>
                          <div
                            v-if="tile.panel.dialogue"
                            class="manga-print-dialogue-text px-3 py-2 rounded-xl bg-white border-2 border-[#1E1B18] text-xs sm:text-sm font-comic font-bold italic text-[#1E1B18] shadow-[2px_2px_0_#1E1B18]"
                          >
                            💬 “{{ tile.panel.dialogue }}”
                          </div>
                        </div>

                        <!-- Panel Actions -->
                        <div
                          v-if="getGeneratedPanelImage(tile.panel.panelNumber)"
                          class="manga-no-print pt-2 border-t border-[#D1C7BD]/60 flex items-center justify-end gap-2"
                        >
                          <button
                            type="button"
                            @click="downloadDataUrl(getGeneratedPanelImage(tile.panel.panelNumber)!, `panel_${tile.panel.panelNumber}_${date}.jpg`, getGeneratedPanel(tile.panel.panelNumber)?.r2Key)"
                            class="px-2.5 py-1 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] transition-colors cursor-pointer"
                          >
                            Download
                          </button>
                          <button
                            v-if="isAdmin"
                            type="button"
                            @click="regeneratePanel(tile.panel.panelNumber, date)"
                            :disabled="regeneratingPanels.has(tile.panel.panelNumber) || isGenerating"
                            class="px-2.5 py-1 rounded-lg bg-[#72383D] hover:bg-[#542529] text-white text-[11px] font-bold transition-colors disabled:opacity-50 cursor-pointer"
                          >
                            Regenerate
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </section>
          </div>
        </div>

        <!-- Fullscreen Image Lightbox -->
        <div
          v-if="lightboxImage"
          class="manga-no-print fixed inset-0 z-60 bg-black/85 flex flex-col items-center justify-center p-4"
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
      </div>
    </transition>
  </Teleport>
</template>
