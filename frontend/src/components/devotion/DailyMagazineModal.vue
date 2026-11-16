<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';
import { useMagazineImages, type MagazineReadingKey } from '../../composables/useMagazineImages';
import { useStoryToManga } from '../../composables/useStoryToManga';
import MagazineDocument, { type MagazineMassReadings } from './MagazineDocument.vue';

const props = defineProps<{
  isOpen: boolean;
  date: string;
  formattedDate: string;
  readings: MagazineMassReadings | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const includeCommentaries = ref(true);
const includeImages = ref(true);
const includeComics = ref(true);
const isSettingsOpen = ref(false);
const settingsMenuRef = ref<HTMLElement | null>(null);

const handleDocumentClick = (event: MouseEvent) => {
  if (!isSettingsOpen.value) return;
  const target = event.target as Node | null;
  if (settingsMenuRef.value && target && !settingsMenuRef.value.contains(target)) {
    isSettingsOpen.value = false;
  }
};

const handleDocumentKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && isSettingsOpen.value) {
    isSettingsOpen.value = false;
  }
};

onMounted(() => {
  document.addEventListener('click', handleDocumentClick);
  document.addEventListener('keydown', handleDocumentKeydown);
});

onUnmounted(() => {
  document.removeEventListener('click', handleDocumentClick);
  document.removeEventListener('keydown', handleDocumentKeydown);
});

const {
  images,
  generatingKeys,
  isGeneratingAll,
  imageError,
  resolveImageSrc,
  loadImagesForDate,
  generateReadingImage,
  generateMissingImages,
} = useMagazineImages();

const {
  style: mangaStyle,
  moralLesson,
  storyAnalysis,
  storyBreakdown,
  generatedPanels,
  getProxyAssetUrl,
  loadStateFromServer,
} = useStoryToManga();

const syncMagazineData = async () => {
  if (!props.isOpen) return;
  await Promise.all([
    loadImagesForDate(props.date),
    loadStateFromServer(props.date, true),
  ]);
};

watch(
  () => [props.isOpen, props.date] as const,
  ([open]) => {
    if (open) {
      syncMagazineData();
    }
  },
  { immediate: true }
);

const handleGenerateSingleImage = async (payload: {
  key: MagazineReadingKey;
  citation: string;
  text?: string;
  context?: string;
  force?: boolean;
}) => {
  await generateReadingImage(
    props.date,
    payload.key,
    payload.citation,
    payload.text,
    payload.context,
    Boolean(payload.force)
  );
};

const handleGenerateAllMissing = async () => {
  if (!props.readings) return;
  await generateMissingImages(props.date, props.readings);
};

const buildMagazineRouteUrl = (autoPrint: boolean) => {
  const params = new URLSearchParams({
    date: props.date,
    commentaries: includeCommentaries.value ? '1' : '0',
    images: includeImages.value ? '1' : '0',
    comics: includeComics.value ? '1' : '0',
  });
  if (autoPrint) {
    params.set('print', 'true');
  }
  return `/#/magazine?${params.toString()}`;
};

const openPrintView = (autoPrint: boolean) => {
  window.open(buildMagazineRouteUrl(autoPrint), '_blank');
};
</script>

<template>
  <Teleport to="body">
    <transition name="fade">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex flex-col bg-[#24201C]/80 backdrop-blur-xs"
        @click.self="emit('close')"
      >
        <!-- Studio Toolbar -->
        <header class="bg-[#FDFBF7] border-b border-[#D1C7BD] px-4 py-3 shadow-sm shrink-0">
          <div class="max-w-6xl mx-auto flex flex-wrap items-center justify-between gap-3">
            <div class="flex items-center gap-3">
              <button
                type="button"
                @click="emit('close')"
                class="p-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[#322D29] transition-colors cursor-pointer"
                aria-label="Close Magazine Preview"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
              <div>
                <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D] block">
                  Lumen Liturgiae · Magazine Edition
                </span>
                <h2 class="text-base sm:text-lg font-serif font-bold text-[#1E1B18]">
                  {{ formattedDate }}
                </h2>
              </div>
            </div>

            <!-- Options & Actions -->
            <div class="flex items-center gap-2">
              <div ref="settingsMenuRef" class="relative">
                <button
                  type="button"
                  @click="isSettingsOpen = !isSettingsOpen"
                  aria-haspopup="menu"
                  :aria-expanded="isSettingsOpen"
                  title="Magazine sections, image generation & full-page view"
                  class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[#322D29] text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="3"></circle>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                  </svg>
                  <span class="hidden sm:inline">Options</span>
                </button>

                <div
                  v-if="isSettingsOpen"
                  role="menu"
                  class="absolute right-0 mt-2 w-64 rounded-2xl bg-[#FDFBF7] border border-[#D1C7BD] shadow-xl p-2 z-50 space-y-1"
                >
                  <p class="px-2.5 py-1 text-[10px] font-bold uppercase tracking-widest text-[#72383D]/70">
                    Include Sections
                  </p>
                  <label class="flex items-center justify-between px-2.5 py-2 rounded-xl hover:bg-[#FAF6EE] text-xs font-semibold text-[#322D29] cursor-pointer">
                    <span>AI Commentaries</span>
                    <input v-model="includeCommentaries" type="checkbox" class="accent-[#72383D]" />
                  </label>
                  <label class="flex items-center justify-between px-2.5 py-2 rounded-xl hover:bg-[#FAF6EE] text-xs font-semibold text-[#322D29] cursor-pointer">
                    <span>Reading Images</span>
                    <input v-model="includeImages" type="checkbox" class="accent-[#72383D]" />
                  </label>
                  <label class="flex items-center justify-between px-2.5 py-2 rounded-xl hover:bg-[#FAF6EE] text-xs font-semibold text-[#322D29] cursor-pointer">
                    <span>Comics at End</span>
                    <input v-model="includeComics" type="checkbox" class="accent-[#72383D]" />
                  </label>

                  <div class="my-1 border-t border-[#D1C7BD]/60"></div>

                  <button
                    v-if="includeImages && readings"
                    type="button"
                    role="menuitem"
                    :disabled="isGeneratingAll"
                    @click="handleGenerateAllMissing(); isSettingsOpen = false"
                    class="w-full flex items-center gap-2 px-2.5 py-2 rounded-xl hover:bg-[#FAF6EE] text-xs font-bold text-[#72383D] cursor-pointer disabled:opacity-50 text-left"
                  >
                    <span>🎨</span>
                    <span>{{ isGeneratingAll ? 'Generating Art...' : 'Generate Reading Images' }}</span>
                  </button>

                  <button
                    type="button"
                    role="menuitem"
                    @click="openPrintView(false); isSettingsOpen = false"
                    class="w-full flex items-center gap-2 px-2.5 py-2 rounded-xl hover:bg-[#FAF6EE] text-xs font-bold text-[#322D29] cursor-pointer text-left"
                  >
                    <span>📄</span>
                    <span>Open Full-Page View</span>
                  </button>
                </div>
              </div>

              <button
                type="button"
                @click="openPrintView(true)"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] hover:from-[#8B444A] hover:to-[#72383D] text-white text-xs font-bold uppercase tracking-wider shadow-xs transition-all cursor-pointer active:scale-95"
              >
                <span>🖨️</span>
                <span>Print / Save PDF</span>
              </button>
            </div>
          </div>

          <div
            v-if="imageError"
            class="max-w-6xl mx-auto mt-2 px-3 py-2 rounded-lg bg-red-50 border border-red-200 text-xs text-red-800"
          >
            {{ imageError }}
          </div>
        </header>

        <!-- Scrollable Magazine Preview -->
        <div class="flex-grow overflow-y-auto p-4 sm:p-8">
          <MagazineDocument
            v-if="readings"
            :date="date"
            :formattedDate="formattedDate"
            :readings="readings"
            :images="images"
            :resolveImageSrc="resolveImageSrc"
            :generatingKeys="generatingKeys"
            :interactive="true"
            :includeCommentaries="includeCommentaries"
            :includeImages="includeImages"
            :includeComics="includeComics"
            :mangaStyle="mangaStyle"
            :storyAnalysis="storyAnalysis"
            :moralLesson="moralLesson"
            :storyBreakdown="storyBreakdown"
            :generatedPanels="generatedPanels"
            :getProxyAssetUrl="getProxyAssetUrl"
            @generate-image="handleGenerateSingleImage"
          />
        </div>
      </div>
    </transition>
  </Teleport>
</template>
