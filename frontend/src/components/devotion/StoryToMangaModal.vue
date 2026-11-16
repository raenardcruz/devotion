<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import {
  useStoryToManga,
  type ReadingsPayload,
  type ComicStyle,
  type MangaPanelPlan,
  type MangaSetting,
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
  loadStateFromServer,
  loadStateForDate,
  loadVersionFromServer,
  saveStateToServer,
  resetWorkflowState,
  updateImageModel,
  updateStyle,
  getProxyAssetUrl,
  generateStoryWithAI,
  analyzeUserStory,
  generateAllCharacterRefs,
  planStoryPanels,
  generateAllPanels,
  regenerateCharacter,
  regeneratePanel,
} = useStoryToManga();

type WizardStep = 1 | 2 | 3;

const activeStep = ref<WizardStep>(1);
const showSettingsBar = ref(false);
const isLoadingInitial = ref(false);
const storyDraft = ref('');
const titleDraft = ref('');
const moralLessonDraft = ref('');
const settingDraft = ref<MangaSetting>({
  timePeriod: 'Modern day',
  location: '',
  mood: '',
});
const lastAnalyzedStoryText = ref('');
const createNewVersionOnSave = ref(false);

const characterUploads = ref<Record<number, { image: string; fileName?: string }>>({});
const refineFromExisting = ref<Record<number, boolean>>({});
const editingPanelNumber = ref<number | null>(null);
const isSavingChanges = ref(false);
const saveFeedback = ref<string | null>(null);

const layoutMode = ref<'storybook' | 'grid'>('grid');
const lightboxImage = ref<{ url: string; caption: string } | null>(null);
const isExportingPdf = ref(false);

const comicLayoutRows = computed(() =>
  groupPanelsIntoComicRows(storyBreakdown.value?.panels || [])
);

const storyWordCount = computed(() => {
  const trimmed = storyDraft.value.trim();
  return trimmed ? trimmed.split(/\s+/).length : 0;
});

const canProceedStep1 = computed(() => {
  return (
    storyWordCount.value > 0 &&
    storyWordCount.value <= 500 &&
    !isAnalyzing.value &&
    !isGenerating.value
  );
});

const charactersList = computed(() => storyAnalysis.value?.characters || []);

const generatedCharacterCount = computed(() => {
  return charactersList.value.filter((_, idx) => Boolean(characterReferences.value[idx]?.image)).length;
});

const allCharactersHaveImages = computed(() => {
  return (
    charactersList.value.length > 0 &&
    generatedCharacterCount.value === charactersList.value.length
  );
});

const canProceedStep2 = computed(() => {
  return (
    charactersList.value.length > 0 &&
    charactersList.value.every(
      (c) => c.name.trim() !== '' && c.physicalDescription.trim() !== ''
    ) &&
    !isGenerating.value
  );
});

const plannedPanelCount = computed(() => storyBreakdown.value?.panels?.length || 0);
const generatedPanelCount = computed(() => generatedPanels.value.filter((p) => Boolean(p.image)).length);
const allPanelsGenerated = computed(
  () => plannedPanelCount.value > 0 && generatedPanelCount.value >= plannedPanelCount.value
);

const syncDraftsFromAnalysis = () => {
  storyDraft.value = generatedStory.value || storyAnalysis.value?.story || '';
  titleDraft.value = storyAnalysis.value?.title || '';
  moralLessonDraft.value = moralLesson.value || storyAnalysis.value?.moralLesson || '';
  settingDraft.value = {
    timePeriod: storyAnalysis.value?.setting?.timePeriod || 'Modern day',
    location: storyAnalysis.value?.setting?.location || '',
    mood: storyAnalysis.value?.setting?.mood || '',
  };
  lastAnalyzedStoryText.value = storyDraft.value.trim();
};

const clearTransientUiState = () => {
  characterUploads.value = {};
  refineFromExisting.value = {};
  editingPanelNumber.value = null;
  saveFeedback.value = null;
};

onMounted(() => {
  syncAdminState();
  loadConfig();
});

// Load saved story on open without auto-triggering generation
watch(
  () => props.isOpen,
  async (open) => {
    if (!open) return;
    syncAdminState();
    isLoadingInitial.value = true;
    clearTransientUiState();

    try {
      await loadConfig();
      if (!isAdmin.value) {
        showSettingsBar.value = false;
      }

      let hasExisting = await loadStateFromServer(props.date, true);
      if (!hasExisting && isAdmin.value) {
        hasExisting = await loadStateForDate(props.date);
      }

      if (hasExisting && storyAnalysis.value) {
        syncDraftsFromAnalysis();
        createNewVersionOnSave.value = false;
        activeStep.value = 3;
      } else {
        resetWorkflowState(false);
        storyDraft.value = '';
        titleDraft.value = '';
        moralLessonDraft.value = '';
        settingDraft.value = { timePeriod: 'Modern day', location: '', mood: '' };
        lastAnalyzedStoryText.value = '';
        createNewVersionOnSave.value = false;
        activeStep.value = 1;
      }
    } finally {
      isLoadingInitial.value = false;
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
    const ok = await loadVersionFromServer(id);
    if (ok) {
      syncDraftsFromAnalysis();
      clearTransientUiState();
      createNewVersionOnSave.value = false;
    }
  }
};

const handleStartNewStory = () => {
  if (!isAdmin.value || isGenerating.value) return;
  const hadVersions = availableVersions.value.length > 0;
  resetWorkflowState(true);
  storyDraft.value = '';
  titleDraft.value = '';
  moralLessonDraft.value = '';
  settingDraft.value = { timePeriod: 'Modern day', location: '', mood: '' };
  lastAnalyzedStoryText.value = '';
  clearTransientUiState();
  createNewVersionOnSave.value = hadVersions;
  activeStep.value = 1;
};

// Step 1: AI Pencil icon auto-generates a story from today's readings (stays on Step 1 so user can review/edit)
const handleAiGenerateStory = async () => {
  if (!isAdmin.value || isGenerating.value) return;
  const hadVersions = availableVersions.value.length > 0;
  const analysis = await generateStoryWithAI(props.date, props.readings);
  if (analysis) {
    if (hadVersions) {
      createNewVersionOnSave.value = true;
      currentStoryId.value = null;
    }
    syncDraftsFromAnalysis();
    clearTransientUiState();
  }
};

// Step 1 -> Step 2
const handleStep1Next = async () => {
  if (!isAdmin.value || !canProceedStep1.value) return;
  const trimmed = storyDraft.value.trim();

  if (storyAnalysis.value && trimmed === lastAnalyzedStoryText.value) {
    if (titleDraft.value.trim()) {
      storyAnalysis.value.title = titleDraft.value.trim();
    }
    if (moralLessonDraft.value.trim()) {
      moralLesson.value = moralLessonDraft.value.trim();
      storyAnalysis.value.moralLesson = moralLessonDraft.value.trim();
    }
    if (settingDraft.value.location.trim()) {
      storyAnalysis.value.setting = { ...settingDraft.value };
    }
    activeStep.value = 2;
    return;
  }

  const hadVersions = availableVersions.value.length > 0;
  if (hadVersions && trimmed !== generatedStory.value.trim()) {
    createNewVersionOnSave.value = true;
    currentStoryId.value = null;
  }

  const analysis = await analyzeUserStory(props.date, props.readings, trimmed, {
    title: titleDraft.value,
    moralLesson: moralLessonDraft.value,
  });
  if (analysis) {
    syncDraftsFromAnalysis();
    clearTransientUiState();
    activeStep.value = 2;
  }
};

// Step 2: Character Design Adjustments
const addCharacter = () => {
  if (!storyAnalysis.value || storyAnalysis.value.characters.length >= 4) return;
  storyAnalysis.value.characters.push({
    name: '',
    role: storyAnalysis.value.characters.length === 0 ? 'Protagonist' : 'Supporting',
    physicalDescription: '',
    personality: '',
  });
  characterReferences.value = [
    ...characterReferences.value,
    { name: '', image: '', description: '' },
  ];
};

const removeCharacter = (idx: number) => {
  if (!storyAnalysis.value || storyAnalysis.value.characters.length <= 1) return;
  storyAnalysis.value.characters.splice(idx, 1);

  const nextRefs = [...characterReferences.value];
  if (idx < nextRefs.length) {
    nextRefs.splice(idx, 1);
  }
  characterReferences.value = nextRefs;

  const nextUploads: Record<number, { image: string; fileName?: string }> = {};
  Object.entries(characterUploads.value).forEach(([k, v]) => {
    const keyIdx = Number(k);
    if (keyIdx < idx) nextUploads[keyIdx] = v;
    else if (keyIdx > idx) nextUploads[keyIdx - 1] = v;
  });
  characterUploads.value = nextUploads;

  const nextRefine: Record<number, boolean> = {};
  Object.entries(refineFromExisting.value).forEach(([k, v]) => {
    const keyIdx = Number(k);
    if (keyIdx < idx) nextRefine[keyIdx] = v;
    else if (keyIdx > idx) nextRefine[keyIdx - 1] = v;
  });
  refineFromExisting.value = nextRefine;
};

const handleCharacterRefUpload = (idx: number, event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input?.files?.[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = () => {
    const result = typeof reader.result === 'string' ? reader.result : '';
    if (result) {
      characterUploads.value = {
        ...characterUploads.value,
        [idx]: { image: result, fileName: file.name },
      };
    }
  };
  reader.readAsDataURL(file);
};

const clearCharacterRefUpload = (idx: number) => {
  const next = { ...characterUploads.value };
  delete next[idx];
  characterUploads.value = next;
};

const buildUploadsForCharacter = (idx: number): MangaUploadedReference[] => {
  const char = storyAnalysis.value?.characters[idx];
  if (!char || !char.name.trim()) return [];
  const uploaded = characterUploads.value[idx];
  if (uploaded?.image) {
    return [
      {
        id: `upload_${idx}`,
        name: char.name.trim(),
        image: uploaded.image,
        fileName: uploaded.fileName,
      },
    ];
  }
  const existingImg = characterReferences.value[idx]?.image;
  if (refineFromExisting.value[idx] && existingImg) {
    return [
      {
        id: `refine_${idx}`,
        name: char.name.trim(),
        image: existingImg,
      },
    ];
  }
  return [];
};

const handleGenerateSingleCharacter = async (idx: number) => {
  if (!isAdmin.value || isGenerating.value) return;
  const uploads = buildUploadsForCharacter(idx);
  await regenerateCharacter(
    idx,
    generatedPanels.value.length > 0 ? props.date : undefined,
    uploads.length > 0 ? uploads : undefined
  );
};

const handleGenerateAllCharacters = async (onlyMissing = false) => {
  if (!isAdmin.value || !canProceedStep2.value) return;
  const allUploads: MangaUploadedReference[] = [];
  charactersList.value.forEach((_, idx) => {
    allUploads.push(...buildUploadsForCharacter(idx));
  });
  await generateAllCharacterRefs(
    allUploads.length > 0 ? allUploads : undefined,
    onlyMissing
  );
};

// Step 2 -> Step 3
const handleStep2Next = async () => {
  if (!isAdmin.value || !canProceedStep2.value) return;
  activeStep.value = 3;
  if (!storyBreakdown.value?.panels?.length) {
    await planStoryPanels();
  }
};

// Step 3: Panels
const handleReplanStoryboard = async () => {
  if (!isAdmin.value || isGenerating.value) return;
  editingPanelNumber.value = null;
  await planStoryPanels();
};

const handleGenerateAllPanels = async (onlyMissing = false) => {
  if (!isAdmin.value || isGenerating.value) return;
  editingPanelNumber.value = null;
  const ok = await generateAllPanels(props.date, createNewVersionOnSave.value, onlyMissing);
  if (ok) {
    createNewVersionOnSave.value = false;
    saveFeedback.value = 'Saved to story archive';
    setTimeout(() => {
      saveFeedback.value = null;
    }, 3000);
  }
};

const handleSaveEdits = async () => {
  if (!isAdmin.value || isGenerating.value || isSavingChanges.value) return;
  isSavingChanges.value = true;
  saveFeedback.value = null;
  try {
    const saved = await saveStateToServer(props.date, createNewVersionOnSave.value);
    if (saved) {
      createNewVersionOnSave.value = false;
      editingPanelNumber.value = null;
      saveFeedback.value = 'Changes saved';
      setTimeout(() => {
        saveFeedback.value = null;
      }, 3000);
    }
  } finally {
    isSavingChanges.value = false;
  }
};

const toggleEditPanel = (panelNumber: number) => {
  editingPanelNumber.value = editingPanelNumber.value === panelNumber ? null : panelNumber;
};

// Navigation via Stepper Tabs
const canNavigateToStep = (target: WizardStep): boolean => {
  if (!isAdmin.value || isGenerating.value) return false;
  if (target === 1) return true;
  if (target === 2) return Boolean(storyAnalysis.value);
  if (target === 3) return Boolean(storyAnalysis.value && charactersList.value.length > 0);
  return false;
};

const navigateToStep = async (target: WizardStep) => {
  if (!canNavigateToStep(target)) return;
  activeStep.value = target;
  if (target === 3 && !storyBreakdown.value?.panels?.length) {
    await planStoryPanels();
  }
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

const wizardSteps = computed(() => {
  const step1Done = Boolean(storyAnalysis.value && storyDraft.value.trim());
  const step2Done = allCharactersHaveImages.value;
  const step3Done = allPanelsGenerated.value;

  return [
    {
      step: 1 as WizardStep,
      title: '1. Story',
      subtitle: 'Write or AI-generate',
      active: activeStep.value === 1,
      busy: isAnalyzing.value,
      done: step1Done && activeStep.value !== 1,
      enabled: canNavigateToStep(1),
    },
    {
      step: 2 as WizardStep,
      title: '2. Characters',
      subtitle: 'Design & adjust',
      active: activeStep.value === 2,
      busy: isGeneratingCharacters.value,
      done: step2Done && activeStep.value !== 2,
      enabled: canNavigateToStep(2),
    },
    {
      step: 3 as WizardStep,
      title: '3. Panels',
      subtitle: 'Plan & illustrate',
      active: activeStep.value === 3,
      busy: isChunking.value || isGeneratingPanels.value,
      done: step3Done,
      enabled: canNavigateToStep(3),
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
                    v-if="currentStoryId"
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
                  {{ formattedDate }} · Step-by-Step Studio: Story → Characters → Panels
                </p>
              </div>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
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
                <span>Model &amp; Style</span>
              </button>

              <button
                v-if="isAdmin && (storyAnalysis || storyDraft.trim())"
                type="button"
                @click="handleStartNewStory"
                :disabled="isGenerating"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-white/15 hover:bg-white/25 text-white text-xs font-bold uppercase tracking-wider transition-colors disabled:opacity-50 cursor-pointer"
                title="Start a brand-new story from Step 1"
              >
                <span>+ New Story</span>
              </button>

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

          <!-- 3-Step Interactive Guide Bar (Admin Only) -->
          <div
            v-if="isAdmin && !isLoadingInitial"
            class="manga-no-print px-5 py-3 bg-[#FAF6EE] border-b border-[#D1C7BD]/70 shrink-0"
          >
            <div class="grid grid-cols-3 gap-2 sm:gap-3">
              <button
                v-for="item in wizardSteps"
                :key="item.step"
                type="button"
                @click="navigateToStep(item.step)"
                :disabled="!item.enabled"
                class="flex items-center gap-2.5 px-3 py-2 rounded-xl border text-left transition-all"
                :class="[
                  item.active
                    ? 'bg-[#72383D] border-[#72383D] text-white shadow-xs'
                    : item.done
                      ? 'bg-emerald-50 border-emerald-200 text-emerald-900 hover:bg-emerald-100/70'
                      : 'bg-white/75 border-[#D1C7BD]/70 text-[#322D29]/65',
                  item.enabled ? 'cursor-pointer' : 'opacity-55 cursor-not-allowed',
                ]"
              >
                <span
                  v-if="item.busy"
                  class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin shrink-0"
                ></span>
                <span
                  v-else-if="item.done"
                  class="w-5 h-5 rounded-full bg-emerald-600 text-white text-[11px] font-bold flex items-center justify-center shrink-0"
                >
                  ✓
                </span>
                <span
                  v-else
                  class="w-5 h-5 rounded-full text-[11px] font-bold flex items-center justify-center shrink-0"
                  :class="item.active ? 'bg-white/20 text-white' : 'bg-[#EFE9E1] text-[#72383D]'"
                >
                  {{ item.step }}
                </span>

                <div class="min-w-0">
                  <div class="text-xs sm:text-sm font-bold truncate">{{ item.title }}</div>
                  <div
                    class="text-[10px] truncate hidden sm:block"
                    :class="item.active ? 'text-white/80' : 'opacity-70'"
                  >
                    {{ item.subtitle }}
                  </div>
                </div>
              </button>
            </div>

            <p v-if="currentStepText && isGenerating" class="text-xs text-[#72383D] font-semibold mt-2 animate-pulse">
              {{ currentStepText }}
            </p>
          </div>

          <!-- Modal Body (Scrollable) -->
          <div class="manga-print-body p-5 sm:p-6 overflow-y-auto space-y-6 flex-1">
            <!-- Initial Load Spinner -->
            <div
              v-if="isLoadingInitial"
              class="py-16 flex flex-col items-center justify-center text-center space-y-3"
            >
              <div class="w-10 h-10 border-4 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
              <p class="text-xs font-bold uppercase tracking-wider text-[#72383D]">
                Loading Story Studio...
              </p>
            </div>

            <!-- Error Alert -->
            <div
              v-if="error && !isLoadingInitial"
              class="manga-no-print p-4 rounded-2xl bg-red-50 border border-red-200 text-red-900 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3"
            >
              <div>
                <h4 class="font-bold text-sm">Story-to-Manga Notice</h4>
                <p class="text-xs text-red-700 mt-0.5">{{ error }}</p>
              </div>
              <router-link
                v-if="!isAdmin"
                to="/admin"
                @click="emit('close')"
                class="px-4 py-2 rounded-xl bg-[#72383D] hover:bg-[#542529] text-white text-xs font-bold uppercase tracking-wider shrink-0"
              >
                Admin Login
              </router-link>
            </div>

            <!-- ============================================================ -->
            <!-- ADMIN STEP-BY-STEP WIZARD                                    -->
            <!-- ============================================================ -->
            <template v-if="isAdmin && !isLoadingInitial">
              <!-- STEP 1: WRITE OR AI-GENERATE STORY -->
              <section v-if="activeStep === 1" class="space-y-5">
                <div class="rounded-2xl bg-[#FAF6EE] border border-[#D1C7BD] p-4 flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D]">
                      Step 1 of 3 · Story
                    </span>
                    <h3 class="text-base sm:text-lg font-serif font-bold text-[#322D29]">
                      Write Your Story or Auto-Generate with AI
                    </h3>
                    <p class="text-xs text-[#322D29]/75 mt-0.5">
                      Enter your own story below (up to 500 words), or click the <strong>AI Pencil</strong> button to automatically craft a story inspired by today’s Mass readings.
                    </p>
                  </div>

                  <!-- Version Selector if prior versions exist -->
                  <div v-if="availableVersions.length > 0" class="flex items-center gap-2">
                    <label for="step1-version-select" class="text-[11px] font-bold uppercase tracking-wider text-[#72383D]">
                      Saved
                    </label>
                    <select
                      id="step1-version-select"
                      :value="currentStoryId || ''"
                      @change="handleVersionChange"
                      :disabled="isGenerating"
                      class="px-2.5 py-1.5 rounded-xl border border-[#D1C7BD] bg-white text-xs font-bold text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/30 cursor-pointer"
                    >
                      <option v-if="!currentStoryId" value="">New Draft</option>
                      <option v-for="(ver, vIdx) in availableVersions" :key="ver.id" :value="ver.id">
                        v{{ ver.version }}{{ vIdx === 0 ? ' (Latest)' : '' }} — {{ ver.title }}
                      </option>
                    </select>
                  </div>
                </div>

                <!-- Story Composer Card -->
                <div class="bg-white border-2 border-[#D1C7BD] rounded-2xl p-5 shadow-xs space-y-4">
                  <!-- Title & AI Pencil Toolbar -->
                  <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-3">
                    <div class="flex-1">
                      <label for="manga-story-title" class="block text-xs font-bold uppercase tracking-wider text-[#72383D] mb-1">
                        Story Title <span class="font-normal text-[#322D29]/55">(optional — auto-titled if blank)</span>
                      </label>
                      <input
                        id="manga-story-title"
                        v-model="titleDraft"
                        type="text"
                        :disabled="isAnalyzing"
                        placeholder="e.g., Leo's Courage at School"
                        class="w-full px-3.5 py-2 rounded-xl border border-[#D1C7BD] bg-[#FBF8F3] text-sm font-serif font-bold text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                      />
                    </div>

                    <!-- AI Pencil Icon Button -->
                    <button
                      type="button"
                      @click="handleAiGenerateStory"
                      :disabled="isGenerating"
                      class="inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs transition-all active:scale-95 disabled:opacity-50 cursor-pointer shrink-0"
                      title="Auto-generate a story from today's Mass readings using AI"
                    >
                      <span
                        v-if="isAnalyzing"
                        class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"
                      ></span>
                      <!-- AI Pencil + Sparkles SVG Icon -->
                      <svg
                        v-else
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        aria-hidden="true"
                      >
                        <path d="M12 20h9" />
                        <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4Z" />
                        <path d="M5 3v4" />
                        <path d="M3 5h4" />
                        <path d="M19 14v3" />
                        <path d="M17.5 15.5h3" />
                      </svg>
                      <span>{{ isAnalyzing ? 'Writing Story...' : 'AI Auto-Write Story' }}</span>
                    </button>
                  </div>

                  <!-- Story Textarea -->
                  <div class="space-y-1.5">
                    <div class="flex items-center justify-between gap-2">
                      <label for="manga-story-textarea" class="text-xs font-bold uppercase tracking-wider text-[#72383D]">
                        Story Narrative *
                      </label>
                      <div class="flex items-center gap-3">
                        <button
                          v-if="storyDraft.trim()"
                          type="button"
                          @click="storyDraft = ''"
                          :disabled="isAnalyzing"
                          class="text-[11px] font-bold text-[#72383D] hover:underline cursor-pointer"
                        >
                          Clear
                        </button>
                        <span
                          class="text-xs font-bold"
                          :class="storyWordCount > 500 ? 'text-red-700' : 'text-[#322D29]/65'"
                        >
                          {{ storyWordCount }} / 500 words
                        </span>
                      </div>
                    </div>

                    <div class="relative">
                      <textarea
                        id="manga-story-textarea"
                        v-model="storyDraft"
                        rows="8"
                        :disabled="isAnalyzing"
                        placeholder="Write or paste your story here (up to 500 words)...&#10;&#10;Or click the AI Auto-Write Story pencil button above to automatically craft a story from today's Mass readings that you can review and edit!"
                        class="w-full rounded-xl border-2 border-[#D1C7BD] bg-[#FBF8F3] p-4 text-sm sm:text-base font-comic text-[#1E1B18] focus:outline-none focus:border-[#72383D] leading-relaxed disabled:opacity-60"
                      ></textarea>
                    </div>
                  </div>

                  <!-- Moral Lesson & Setting (Optional / Editable) -->
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-3 pt-2 border-t border-[#D1C7BD]/60">
                    <div>
                      <label for="manga-moral-lesson" class="block text-[11px] font-bold uppercase tracking-wider text-[#322D29]/75 mb-1">
                        Lesson / Takeaway <span class="font-normal text-[#322D29]/50">(optional)</span>
                      </label>
                      <input
                        id="manga-moral-lesson"
                        v-model="moralLessonDraft"
                        type="text"
                        :disabled="isAnalyzing"
                        placeholder="e.g., Forgiving from the heart brings peace."
                        class="w-full px-3 py-2 rounded-xl border border-[#D1C7BD] bg-[#FBF8F3] text-xs sm:text-sm text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                      />
                    </div>
                    <div>
                      <label for="manga-setting-location" class="block text-[11px] font-bold uppercase tracking-wider text-[#322D29]/75 mb-1">
                        Setting / Location <span class="font-normal text-[#322D29]/50">(optional)</span>
                      </label>
                      <input
                        id="manga-setting-location"
                        v-model="settingDraft.location"
                        type="text"
                        :disabled="isAnalyzing"
                        placeholder="e.g., Parish church and neighborhood park"
                        class="w-full px-3 py-2 rounded-xl border border-[#D1C7BD] bg-[#FBF8F3] text-xs sm:text-sm text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                      />
                    </div>
                  </div>
                </div>

                <!-- Step 1 Footer -->
                <div class="flex flex-wrap items-center justify-between gap-3 pt-2">
                  <p class="text-xs text-[#322D29]/65">
                    Next, we’ll extract the characters from your story so you can customize their look.
                  </p>
                  <button
                    type="button"
                    @click="handleStep1Next"
                    :disabled="!canProceedStep1"
                    class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs transition-all disabled:opacity-50 cursor-pointer"
                  >
                    <span
                      v-if="isAnalyzing"
                      class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"
                    ></span>
                    <span>{{ isAnalyzing ? 'Analyzing Story...' : 'Next: Characters →' }}</span>
                  </button>
                </div>
              </section>

              <!-- STEP 2: GENERATE & ADJUST CHARACTER DESIGN -->
              <section v-else-if="activeStep === 2 && storyAnalysis" class="space-y-5">
                <div class="rounded-2xl bg-[#FAF6EE] border border-[#D1C7BD] p-4 flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D]">
                      Step 2 of 3 · Character Design
                    </span>
                    <h3 class="text-base sm:text-lg font-serif font-bold text-[#322D29]">
                      Generate &amp; Adjust Characters for “{{ storyAnalysis.title }}”
                    </h3>
                    <p class="text-xs text-[#322D29]/75 mt-0.5">
                      Customize each character’s appearance, role, or reference photo, then generate or refine their character design sheets.
                    </p>
                  </div>

                  <div class="flex items-center gap-2 flex-wrap">
                    <!-- Quick Art Style Toggle -->
                    <div class="inline-flex rounded-xl border border-[#D1C7BD] bg-white p-0.5">
                      <button
                        type="button"
                        @click="handleStyleChange('manga')"
                        :disabled="isGenerating"
                        class="px-2.5 py-1 rounded-lg text-[11px] font-bold transition-colors cursor-pointer"
                        :class="style === 'manga' ? 'bg-[#72383D] text-white' : 'text-[#322D29]/75 hover:text-[#72383D]'"
                      >
                        🇯🇵 Manga
                      </button>
                      <button
                        type="button"
                        @click="handleStyleChange('comic')"
                        :disabled="isGenerating"
                        class="px-2.5 py-1 rounded-lg text-[11px] font-bold transition-colors cursor-pointer"
                        :class="style === 'comic' ? 'bg-[#72383D] text-white' : 'text-[#322D29]/75 hover:text-[#72383D]'"
                      >
                        🦸 Comic
                      </button>
                    </div>

                    <button
                      type="button"
                      @click="addCharacter"
                      :disabled="charactersList.length >= 4 || isGenerating"
                      class="px-3 py-1.5 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors disabled:opacity-40 cursor-pointer"
                    >
                      + Add Character ({{ charactersList.length }}/4)
                    </button>

                    <button
                      type="button"
                      @click="handleGenerateAllCharacters(false)"
                      :disabled="!canProceedStep2 || isGenerating"
                      class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs transition-colors disabled:opacity-50 cursor-pointer"
                    >
                      <span>🎨</span>
                      <span>{{ allCharactersHaveImages ? 'Regenerate All Designs' : 'Generate All Designs' }}</span>
                    </button>
                  </div>
                </div>

                <!-- Character Design Studio Cards -->
                <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                  <div
                    v-for="(char, idx) in storyAnalysis.characters"
                    :key="idx"
                    class="bg-white border-2 border-[#D1C7BD] rounded-2xl p-4 shadow-xs flex flex-col sm:flex-row gap-4"
                  >
                    <!-- Left: Character Design Sheet Preview -->
                    <div class="sm:w-44 shrink-0 flex flex-col gap-2">
                      <div class="relative aspect-square w-full rounded-xl bg-[#EFE9E1] border border-[#D1C7BD] overflow-hidden flex items-center justify-center">
                        <img
                          v-if="characterReferences[idx]?.image"
                          :src="characterReferences[idx].image"
                          :alt="char.name || `Character ${idx + 1}`"
                          class="w-full h-full object-cover cursor-pointer hover:scale-105 transition-transform"
                          @click="openLightbox(characterReferences[idx].image, `${char.name} — Character Design`)"
                        />
                        <div
                          v-else-if="regeneratingCharacters.has(idx)"
                          class="flex flex-col items-center justify-center p-3 text-center space-y-2"
                        >
                          <div class="w-8 h-8 border-3 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
                          <span class="text-[11px] text-[#72383D] font-semibold">Designing...</span>
                        </div>
                        <div v-else class="p-3 text-center space-y-1.5">
                          <span class="text-2xl block opacity-60">👤</span>
                          <p class="text-[11px] text-[#322D29]/60 font-medium">
                            Not generated yet
                          </p>
                        </div>

                        <!-- Regenerating Overlay -->
                        <div
                          v-if="regeneratingCharacters.has(idx) && characterReferences[idx]?.image"
                          class="absolute inset-0 bg-black/55 flex flex-col items-center justify-center text-white space-y-2"
                        >
                          <div class="w-7 h-7 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                          <span class="text-xs font-semibold">Updating...</span>
                        </div>
                      </div>

                      <div class="flex items-center gap-1.5">
                        <button
                          type="button"
                          @click="handleGenerateSingleCharacter(idx)"
                          :disabled="!char.name.trim() || !char.physicalDescription.trim() || regeneratingCharacters.has(idx) || isGenerating"
                          class="flex-1 px-2.5 py-1.5 rounded-lg bg-[#72383D] hover:bg-[#542529] text-white text-[11px] font-bold transition-colors disabled:opacity-50 cursor-pointer"
                        >
                          {{ characterReferences[idx]?.image ? 'Regenerate' : 'Generate' }}
                        </button>
                        <button
                          v-if="characterReferences[idx]?.image"
                          type="button"
                          @click="downloadDataUrl(characterReferences[idx].image, `${char.name || 'character'}_design.jpg`, characterReferences[idx].r2Key)"
                          class="px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FAF6EE] hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] transition-colors cursor-pointer"
                          title="Download character sheet"
                        >
                          ⬇
                        </button>
                      </div>
                    </div>

                    <!-- Right: Editable Design Controls -->
                    <div class="flex-1 space-y-2.5 min-w-0">
                      <div class="flex items-center justify-between gap-2">
                        <span class="text-[11px] font-bold uppercase tracking-wider text-[#72383D]">
                          Character {{ idx + 1 }}
                        </span>
                        <button
                          v-if="storyAnalysis.characters.length > 1"
                          type="button"
                          @click="removeCharacter(idx)"
                          :disabled="isGenerating"
                          class="text-[11px] font-bold text-red-700 hover:underline disabled:opacity-40 cursor-pointer"
                        >
                          Remove
                        </button>
                      </div>

                      <div class="grid grid-cols-2 gap-2">
                        <div>
                          <label class="block text-[10px] font-bold uppercase text-[#322D29]/65 mb-0.5">Name *</label>
                          <input
                            v-model="char.name"
                            type="text"
                            placeholder="Character name"
                            class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FBF8F3] text-xs font-bold text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                          />
                        </div>
                        <div>
                          <label class="block text-[10px] font-bold uppercase text-[#322D29]/65 mb-0.5">Role</label>
                          <input
                            v-model="char.role"
                            type="text"
                            placeholder="Protagonist / Mentor"
                            class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FBF8F3] text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                          />
                        </div>
                      </div>

                      <div>
                        <label class="block text-[10px] font-bold uppercase text-[#322D29]/65 mb-0.5">
                          Visual Design / Appearance *
                        </label>
                        <textarea
                          v-model="char.physicalDescription"
                          rows="3"
                          placeholder="Age, hair style, outfit, colors, accessories, facial features..."
                          class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FBF8F3] text-xs text-[#1E1B18] leading-relaxed focus:outline-none focus:border-[#72383D]"
                        ></textarea>
                      </div>

                      <div>
                        <label class="block text-[10px] font-bold uppercase text-[#322D29]/65 mb-0.5">Personality</label>
                        <input
                          v-model="char.personality"
                          type="text"
                          placeholder="e.g., Curious, gentle, hopeful"
                          class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-[#FBF8F3] text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                        />
                      </div>

                      <!-- Optional Reference Photo & Refine Existing Toggle -->
                      <div class="pt-1 flex flex-wrap items-center justify-between gap-2">
                        <label class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg border border-[#D1C7BD] bg-[#FAF6EE] hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] cursor-pointer">
                          <span>📷 {{ characterUploads[idx]?.image ? 'Change Photo' : 'Reference Photo' }}</span>
                          <input
                            type="file"
                            accept="image/*"
                            class="hidden"
                            @change="handleCharacterRefUpload(idx, $event)"
                          />
                        </label>

                        <div v-if="characterUploads[idx]?.image" class="flex items-center gap-1.5">
                          <img
                            :src="characterUploads[idx].image"
                            :alt="char.name || 'Reference'"
                            class="w-7 h-7 rounded-md object-cover border border-[#D1C7BD]"
                          />
                          <button
                            type="button"
                            @click="clearCharacterRefUpload(idx)"
                            class="text-[11px] font-bold text-red-700 hover:underline cursor-pointer"
                          >
                            Clear
                          </button>
                        </div>

                        <label
                          v-else-if="characterReferences[idx]?.image"
                          class="inline-flex items-center gap-1.5 text-[11px] font-semibold text-[#322D29]/80 cursor-pointer select-none"
                          title="Pass the current character sheet as a reference when regenerating so tweaks keep the same face/identity"
                        >
                          <input
                            v-model="refineFromExisting[idx]"
                            type="checkbox"
                            class="rounded border-[#D1C7BD] text-[#72383D] focus:ring-[#72383D]"
                          />
                          <span>Refine current look</span>
                        </label>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Step 2 Footer -->
                <div class="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-[#D1C7BD]/70">
                  <button
                    type="button"
                    @click="activeStep = 1"
                    :disabled="isGenerating"
                    class="px-4 py-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold uppercase tracking-wider text-[#322D29] transition-colors disabled:opacity-50 cursor-pointer"
                  >
                    ← Back to Story
                  </button>

                  <div class="flex items-center gap-2 flex-wrap">
                    <span v-if="!allCharactersHaveImages" class="text-xs text-[#322D29]/70">
                      {{ generatedCharacterCount }}/{{ charactersList.length }} character sheets ready
                    </span>
                    <button
                      v-if="!allCharactersHaveImages"
                      type="button"
                      @click="handleGenerateAllCharacters(true)"
                      :disabled="!canProceedStep2 || isGenerating"
                      class="px-4 py-2 rounded-xl border border-[#72383D] bg-[#FAF6EE] hover:bg-[#EFE9E1] text-[#72383D] text-xs font-bold uppercase tracking-wider transition-colors disabled:opacity-50 cursor-pointer"
                    >
                      Generate Missing Designs
                    </button>
                    <button
                      type="button"
                      @click="handleStep2Next"
                      :disabled="!canProceedStep2"
                      class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs transition-all disabled:opacity-50 cursor-pointer"
                    >
                      <span>{{ allCharactersHaveImages ? 'Next: Generate Panels →' : 'Continue to Panels →' }}</span>
                    </button>
                  </div>
                </div>
              </section>

              <!-- STEP 3: PLAN & GENERATE PANELS -->
              <section v-else-if="activeStep === 3 && storyAnalysis" class="space-y-5">
                <!-- Story Header Summary -->
                <div class="bg-white/95 border-2 border-[#1E1B18] rounded-2xl p-5 shadow-[3px_3px_0_#1E1B18] space-y-3">
                  <div class="flex flex-wrap items-start justify-between gap-2 border-b border-[#D1C7BD]/60 pb-3">
                    <div>
                      <span class="font-comic font-bold text-xs uppercase tracking-[0.16em] text-[#72383D] block">
                        Step 3 of 3 · Panels · {{ formattedDate }}
                      </span>
                      <h3 class="manga-print-title text-2xl sm:text-3xl font-comic-display text-[#1E1B18] mt-1">
                        {{ storyAnalysis.title }}
                      </h3>
                    </div>
                    <div class="flex items-center gap-2 text-xs font-comic font-bold text-[#322D29]/80 bg-[#FAF6EE] px-3 py-1 rounded-full border border-[#1E1B18]/40">
                      <span>📍 {{ storyAnalysis.setting.location || 'Daily Life' }}</span>
                      <span v-if="storyAnalysis.setting.mood">· {{ storyAnalysis.setting.mood }}</span>
                    </div>
                  </div>

                  <div
                    v-if="moralLesson"
                    class="p-3 rounded-xl bg-[#FFFDF5] border-2 border-[#1E1B18] border-l-4 border-l-[#72383D] text-sm text-[#1E1B18]"
                  >
                    <span class="font-comic font-bold uppercase tracking-wider text-xs text-[#72383D] block mb-0.5">
                      Lesson Applied from Today's Readings
                    </span>
                    <p class="font-comic font-bold italic leading-relaxed">{{ moralLesson }}</p>
                  </div>

                  <details v-if="generatedStory" class="manga-no-print rounded-xl border border-[#D1C7BD]/80 bg-[#FAF6EE]/70 px-4 py-2">
                    <summary class="cursor-pointer text-xs font-bold uppercase tracking-wider text-[#72383D] select-none">
                      Read Full Story Text
                    </summary>
                    <div class="prose max-w-none text-sm sm:text-base font-comic leading-relaxed text-[#1E1B18] whitespace-pre-line mt-2.5 pt-2.5 border-t border-[#D1C7BD]/60">
                      {{ generatedStory }}
                    </div>
                  </details>
                </div>

                <!-- Planning Storyboard Spinner -->
                <div
                  v-if="isChunking && !storyBreakdown?.panels?.length"
                  class="manga-no-print py-14 flex flex-col items-center justify-center text-center space-y-3 bg-white/80 border border-[#D1C7BD] rounded-2xl"
                >
                  <div class="w-10 h-10 border-4 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
                  <h4 class="font-serif text-base font-bold text-[#72383D]">
                    Planning Storyboard &amp; Comic Layout...
                  </h4>
                  <p class="text-xs text-[#322D29]/70 max-w-md">
                    Breaking your story into panels with camera angles, dialogue, and narrative captions.
                  </p>
                </div>

                <!-- Panel Controls Bar -->
                <div v-else-if="storyBreakdown?.panels?.length" class="space-y-5">
                  <div class="manga-no-print flex flex-wrap items-center justify-between gap-3 bg-[#FAF6EE] border border-[#D1C7BD] rounded-2xl p-4">
                    <div>
                      <h4 class="text-base font-serif font-bold text-[#72383D]">
                        {{ style === 'manga' ? 'Manga' : 'Comic' }} Panels ({{ generatedPanelCount }}/{{ plannedPanelCount }})
                      </h4>
                      <p class="text-xs text-[#322D29]/70">
                        Review or edit any panel’s scene &amp; dialogue below, then generate all panels or one panel at a time.
                      </p>
                    </div>

                    <div class="flex items-center gap-2 flex-wrap">
                      <span
                        v-if="saveFeedback"
                        class="text-xs font-bold text-emerald-700 bg-emerald-50 border border-emerald-200 px-2.5 py-1 rounded-lg"
                      >
                        ✓ {{ saveFeedback }}
                      </span>

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
                        type="button"
                        @click="handleReplanStoryboard"
                        :disabled="isGenerating"
                        class="px-3 py-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors disabled:opacity-50 cursor-pointer"
                        title="Re-generate panel breakdown from the story"
                      >
                        Re-plan Layout
                      </button>

                      <button
                        v-if="!allPanelsGenerated"
                        type="button"
                        @click="handleGenerateAllPanels(generatedPanelCount > 0)"
                        :disabled="isGenerating"
                        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs transition-colors disabled:opacity-50 cursor-pointer"
                      >
                        <span>🎨</span>
                        <span>{{ generatedPanelCount > 0 ? 'Generate Remaining Panels' : 'Generate All Panels' }}</span>
                      </button>

                      <button
                        v-else
                        type="button"
                        @click="handleGenerateAllPanels(false)"
                        :disabled="isGenerating"
                        class="px-3.5 py-2 rounded-xl border border-[#72383D] bg-white hover:bg-[#FAF6EE] text-[#72383D] text-xs font-bold uppercase tracking-wider transition-colors disabled:opacity-50 cursor-pointer"
                      >
                        Regenerate All
                      </button>

                      <button
                        v-if="generatedPanelCount > 0"
                        type="button"
                        @click="handleSaveEdits"
                        :disabled="isGenerating || isSavingChanges"
                        class="px-3.5 py-2 rounded-xl border border-emerald-600 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold uppercase tracking-wider transition-colors disabled:opacity-50 cursor-pointer"
                      >
                        {{ isSavingChanges ? 'Saving...' : 'Save' }}
                      </button>

                      <button
                        v-if="generatedPanelCount > 0"
                        type="button"
                        @click="downloadMangaA4Pdf"
                        :disabled="isExportingPdf"
                        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs cursor-pointer disabled:opacity-50"
                      >
                        <span>📄</span>
                        <span>{{ isExportingPdf ? 'Creating PDF...' : 'A4 PDF' }}</span>
                      </button>
                    </div>
                  </div>

                  <!-- Storybook Layout -->
                  <div v-if="layoutMode === 'storybook'" class="space-y-6">
                    <article
                      v-for="(panel, idx) in storyBreakdown.panels"
                      :key="panel.panelNumber"
                      class="manga-print-row grid grid-cols-1 md:grid-cols-12 gap-4 sm:gap-5 items-stretch"
                    >
                      <!-- Artwork Panel -->
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
                          <div v-else class="p-6 text-center space-y-3 max-w-sm">
                            <div
                              v-if="regeneratingPanels.has(panel.panelNumber)"
                              class="w-10 h-10 mx-auto border-3 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"
                            ></div>
                            <p class="text-xs font-semibold text-[#72383D] uppercase tracking-wider">
                              {{ regeneratingPanels.has(panel.panelNumber) ? `Illustrating Panel ${panel.panelNumber}...` : 'Ready to Illustrate' }}
                            </p>
                            <p class="text-xs text-[#322D29]/80 leading-relaxed">
                              {{ panel.sceneDescription }}
                            </p>
                          </div>

                          <div
                            v-if="regeneratingPanels.has(panel.panelNumber) && getGeneratedPanelImage(panel.panelNumber)"
                            class="absolute inset-0 bg-black/55 flex flex-col items-center justify-center text-white space-y-2"
                          >
                            <div class="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                            <span class="text-xs font-semibold">Regenerating Panel {{ panel.panelNumber }}...</span>
                          </div>
                        </div>
                      </div>

                      <!-- Narrative Panel + Admin Edit -->
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
                            >
                              {{ panel.characters.join(', ') }}
                            </span>
                          </div>

                          <div v-if="editingPanelNumber === panel.panelNumber" class="manga-no-print p-4 space-y-2.5 bg-[#FAF6EE]">
                            <div>
                              <label class="block text-[10px] font-bold uppercase text-[#72383D] mb-0.5">Narrative Caption</label>
                              <textarea
                                v-model="panel.narrative"
                                rows="3"
                                class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                              ></textarea>
                            </div>
                            <div>
                              <label class="block text-[10px] font-bold uppercase text-[#72383D] mb-0.5">Speech Bubble Dialogue</label>
                              <input
                                v-model="panel.dialogue"
                                type="text"
                                class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                              />
                            </div>
                            <div>
                              <label class="block text-[10px] font-bold uppercase text-[#72383D] mb-0.5">Visual Scene Direction</label>
                              <textarea
                                v-model="panel.sceneDescription"
                                rows="2"
                                class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                              ></textarea>
                            </div>
                          </div>

                          <div v-else class="p-5 sm:p-6 space-y-4">
                            <p class="manga-print-narrative-text font-comic font-bold text-lg sm:text-xl leading-relaxed text-[#1E1B18]">
                              {{ getPanelNarrative(panel) }}
                            </p>
                            <div
                              v-if="panel.dialogue"
                              class="manga-print-dialogue-text p-3.5 rounded-2xl bg-white border-2 border-[#1E1B18] text-base sm:text-lg font-comic font-bold italic text-[#1E1B18] shadow-[2px_2px_0_#1E1B18]"
                            >
                              💬 “{{ panel.dialogue }}”
                            </div>
                          </div>
                        </div>

                        <div class="manga-no-print px-5 py-3 border-t border-[#D1C7BD]/70 bg-[#FAF6EE]/80 flex items-center justify-between gap-2 flex-wrap">
                          <button
                            type="button"
                            @click="toggleEditPanel(panel.panelNumber)"
                            class="px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors cursor-pointer"
                          >
                            {{ editingPanelNumber === panel.panelNumber ? 'Done Editing' : 'Edit Scene / Text' }}
                          </button>

                          <div class="flex items-center gap-2">
                            <button
                              v-if="getGeneratedPanelImage(panel.panelNumber)"
                              type="button"
                              @click="downloadDataUrl(getGeneratedPanelImage(panel.panelNumber)!, `panel_${panel.panelNumber}_${date}.jpg`, getGeneratedPanel(panel.panelNumber)?.r2Key)"
                              class="px-3 py-1.5 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold text-[#322D29] transition-colors cursor-pointer"
                            >
                              Download
                            </button>
                            <button
                              type="button"
                              @click="regeneratePanel(panel.panelNumber, props.date)"
                              :disabled="regeneratingPanels.has(panel.panelNumber) || isGenerating"
                              class="px-3 py-1.5 rounded-lg bg-[#72383D] hover:bg-[#542529] text-white text-xs font-bold transition-colors disabled:opacity-50 cursor-pointer"
                            >
                              {{ getGeneratedPanelImage(panel.panelNumber) ? 'Regenerate' : 'Generate Panel' }}
                            </button>
                          </div>
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
                            <div v-else class="p-5 text-center space-y-2 max-w-sm">
                              <div
                                v-if="regeneratingPanels.has(tile.panel.panelNumber)"
                                class="w-9 h-9 mx-auto border-3 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"
                              ></div>
                              <p class="text-xs font-semibold text-[#72383D] uppercase tracking-wider">
                                {{ regeneratingPanels.has(tile.panel.panelNumber) ? `Illustrating Panel ${tile.panel.panelNumber}...` : `Ready (${tile.aspectRatio})` }}
                              </p>
                              <p class="text-xs text-[#322D29]/80 leading-relaxed line-clamp-4">
                                {{ tile.panel.sceneDescription }}
                              </p>
                            </div>

                            <div
                              v-if="regeneratingPanels.has(tile.panel.panelNumber) && getGeneratedPanelImage(tile.panel.panelNumber)"
                              class="absolute inset-0 bg-black/55 flex flex-col items-center justify-center text-white space-y-2"
                            >
                              <div class="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                              <span class="text-xs font-semibold">Regenerating Panel {{ tile.panel.panelNumber }}...</span>
                            </div>
                          </div>

                          <div class="p-3.5 sm:p-4 bg-[#FFFDF5] border-t-2 border-[#1E1B18] flex-1 flex flex-col justify-between gap-2.5">
                            <div v-if="editingPanelNumber === tile.panel.panelNumber" class="manga-no-print space-y-2">
                              <div>
                                <label class="block text-[10px] font-bold uppercase text-[#72383D]">Narrative</label>
                                <textarea
                                  v-model="tile.panel.narrative"
                                  rows="2"
                                  class="w-full px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                                ></textarea>
                              </div>
                              <div>
                                <label class="block text-[10px] font-bold uppercase text-[#72383D]">Dialogue</label>
                                <input
                                  v-model="tile.panel.dialogue"
                                  type="text"
                                  class="w-full px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                                />
                              </div>
                              <div>
                                <label class="block text-[10px] font-bold uppercase text-[#72383D]">Scene Prompt</label>
                                <textarea
                                  v-model="tile.panel.sceneDescription"
                                  rows="2"
                                  class="w-full px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#1E1B18] focus:outline-none focus:border-[#72383D]"
                                ></textarea>
                              </div>
                            </div>

                            <div v-else class="space-y-2">
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

                            <div class="manga-no-print pt-2 border-t border-[#D1C7BD]/60 flex items-center justify-between gap-1.5 flex-wrap">
                              <button
                                type="button"
                                @click="toggleEditPanel(tile.panel.panelNumber)"
                                class="px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] transition-colors cursor-pointer"
                              >
                                {{ editingPanelNumber === tile.panel.panelNumber ? 'Done' : 'Edit' }}
                              </button>

                              <div class="flex items-center gap-1.5">
                                <button
                                  v-if="getGeneratedPanelImage(tile.panel.panelNumber)"
                                  type="button"
                                  @click="downloadDataUrl(getGeneratedPanelImage(tile.panel.panelNumber)!, `panel_${tile.panel.panelNumber}_${date}.jpg`, getGeneratedPanel(tile.panel.panelNumber)?.r2Key)"
                                  class="px-2.5 py-1 rounded-lg border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[11px] font-bold text-[#322D29] transition-colors cursor-pointer"
                                >
                                  Download
                                </button>
                                <button
                                  type="button"
                                  @click="regeneratePanel(tile.panel.panelNumber, props.date)"
                                  :disabled="regeneratingPanels.has(tile.panel.panelNumber) || isGenerating"
                                  class="px-2.5 py-1 rounded-lg bg-[#72383D] hover:bg-[#542529] text-white text-[11px] font-bold transition-colors disabled:opacity-50 cursor-pointer"
                                >
                                  {{ getGeneratedPanelImage(tile.panel.panelNumber) ? 'Regenerate' : 'Generate' }}
                                </button>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Fallback if chunking failed -->
                <div
                  v-else
                  class="manga-no-print py-12 flex flex-col items-center justify-center text-center space-y-3 bg-white/80 border border-[#D1C7BD] rounded-2xl"
                >
                  <p class="text-sm font-serif font-bold text-[#322D29]">No storyboard panels planned yet.</p>
                  <button
                    type="button"
                    @click="handleReplanStoryboard"
                    :disabled="isGenerating"
                    class="px-4 py-2 rounded-xl bg-[#72383D] hover:bg-[#542529] text-white text-xs font-bold uppercase tracking-wider cursor-pointer"
                  >
                    Plan Storyboard Panels
                  </button>
                </div>

                <!-- Step 3 Footer -->
                <div class="manga-no-print flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-[#D1C7BD]/70">
                  <button
                    type="button"
                    @click="activeStep = 2"
                    :disabled="isGenerating"
                    class="px-4 py-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-xs font-bold uppercase tracking-wider text-[#322D29] transition-colors disabled:opacity-50 cursor-pointer"
                  >
                    ← Back to Characters
                  </button>

                  <button
                    type="button"
                    @click="emit('close')"
                    class="px-5 py-2 rounded-xl bg-[#322D29] hover:bg-[#1E1B18] text-white text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer"
                  >
                    Done
                  </button>
                </div>
              </section>
            </template>

            <!-- ============================================================ -->
            <!-- NON-ADMIN READ-ONLY VIEW                                     -->
            <!-- ============================================================ -->
            <template v-else-if="!isAdmin && !isLoadingInitial && storyAnalysis">
              <section class="bg-white/95 border-2 border-[#1E1B18] rounded-2xl p-5 shadow-[3px_3px_0_#1E1B18] space-y-4">
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

                <div
                  v-if="moralLesson"
                  class="p-3.5 rounded-xl bg-[#FFFDF5] border-2 border-[#1E1B18] border-l-4 border-l-[#72383D] text-sm sm:text-base text-[#1E1B18]"
                >
                  <span class="font-comic font-bold uppercase tracking-wider text-xs text-[#72383D] block mb-0.5">
                    Lesson Applied from Today's Readings
                  </span>
                  <p class="font-comic font-bold italic leading-relaxed">{{ moralLesson }}</p>
                </div>

                <details v-if="generatedStory" class="manga-no-print rounded-xl border border-[#D1C7BD]/80 bg-[#FAF6EE]/70 px-4 py-2.5">
                  <summary class="cursor-pointer text-xs font-bold uppercase tracking-wider text-[#72383D] select-none">
                    Read Full Unabridged Story Text
                  </summary>
                  <div class="prose max-w-none text-base font-comic leading-relaxed text-[#1E1B18] whitespace-pre-line mt-3 pt-3 border-t border-[#D1C7BD]/60">
                    {{ generatedStory }}
                  </div>
                </details>
              </section>

              <section v-if="storyBreakdown?.panels?.length" class="space-y-5">
                <div class="manga-no-print flex flex-wrap items-center justify-between gap-3">
                  <h3 class="text-base sm:text-lg font-serif font-bold text-[#72383D]">
                    Illustrated Story &amp; {{ style === 'manga' ? 'Manga' : 'Comic' }} Panels
                  </h3>

                  <div class="flex items-center gap-2 flex-wrap">
                    <div v-if="availableVersions.length > 1" class="inline-flex items-center gap-1.5">
                      <label for="reader-version-select" class="text-[11px] font-bold uppercase tracking-wider text-[#72383D]">
                        Version
                      </label>
                      <select
                        id="reader-version-select"
                        :value="currentStoryId || ''"
                        @change="handleVersionChange"
                        class="px-2.5 py-1.5 rounded-xl border border-[#D1C7BD] bg-white text-xs font-bold text-[#322D29] focus:outline-none cursor-pointer"
                      >
                        <option v-for="(ver, vIdx) in availableVersions" :key="ver.id" :value="ver.id">
                          v{{ ver.version }}{{ vIdx === 0 ? ' (Latest)' : '' }} — {{ ver.title }}
                        </option>
                      </select>
                    </div>

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
                    >
                      <span>📄</span>
                      <span>{{ isExportingPdf ? 'Creating A4 PDF...' : 'Download A4 PDF' }}</span>
                    </button>
                  </div>
                </div>

                <!-- Storybook Layout -->
                <div v-if="layoutMode === 'storybook'" class="space-y-6">
                  <article
                    v-for="(panel, idx) in storyBreakdown.panels"
                    :key="panel.panelNumber"
                    class="manga-print-row grid grid-cols-1 md:grid-cols-12 gap-4 sm:gap-5 items-stretch"
                  >
                    <div
                      class="md:col-span-6 manga-print-col-6 comic-panel-frame bg-white rounded-2xl overflow-hidden flex flex-col"
                      :class="idx % 2 === 1 ? 'md:order-2' : 'md:order-1'"
                    >
                      <div class="px-4 py-2 bg-[#1E1B18] text-white flex items-center justify-between text-xs">
                        <span class="font-comic-display text-sm sm:text-base tracking-wider uppercase">Panel {{ panel.panelNumber }}</span>
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
                          class="w-full h-full object-contain cursor-pointer"
                          @click="openLightbox(getGeneratedPanelImage(panel.panelNumber)!, getPanelLightboxCaption(panel))"
                        />
                      </div>
                    </div>

                    <div
                      class="md:col-span-6 manga-print-col-6 comic-narrative-panel rounded-2xl overflow-hidden flex flex-col justify-between"
                      :class="idx % 2 === 1 ? 'md:order-1' : 'md:order-2'"
                    >
                      <div>
                        <div class="px-4 py-2 bg-[#72383D] text-white flex items-center justify-between gap-2">
                          <span class="font-comic-display text-sm sm:text-base tracking-wider uppercase">
                            Story · Part {{ panel.panelNumber }} of {{ storyBreakdown.panels.length }}
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
                        </div>
                      </div>
                    </div>
                  </article>
                </div>

                <!-- Comic Grid Layout -->
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
                          <span class="font-comic-display text-sm tracking-wider uppercase">Panel {{ tile.panel.panelNumber }}</span>
                          <span class="manga-no-print text-white/75 text-[10px] font-mono">{{ tile.aspectRatio }} · {{ tile.panel.cameraAngle }}</span>
                        </div>
                        <div
                          class="relative w-full bg-[#F5F0E6] flex items-center justify-center overflow-hidden"
                          :class="getTileAspectClass(tile.aspectRatio, tile.colSpan)"
                        >
                          <img
                            v-if="getGeneratedPanelImage(tile.panel.panelNumber)"
                            :src="getGeneratedPanelImage(tile.panel.panelNumber)"
                            :alt="`Panel ${tile.panel.panelNumber}`"
                            class="w-full h-full object-cover cursor-pointer"
                            @click="openLightbox(getGeneratedPanelImage(tile.panel.panelNumber)!, getPanelLightboxCaption(tile.panel))"
                          />
                        </div>
                        <div class="p-3.5 sm:p-4 bg-[#FFFDF5] border-t-2 border-[#1E1B18] flex-1 flex flex-col gap-2">
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
                      </div>
                    </div>
                  </div>
                </div>
              </section>
            </template>
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
