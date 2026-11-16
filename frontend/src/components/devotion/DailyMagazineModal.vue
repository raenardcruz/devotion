<script setup lang="ts">
import { ref, watch } from 'vue';
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
            <div class="flex flex-wrap items-center gap-2.5">
              <label class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#322D29] bg-white px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] cursor-pointer">
                <input v-model="includeCommentaries" type="checkbox" class="accent-[#72383D]" />
                <span>AI Commentaries</span>
              </label>

              <label class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#322D29] bg-white px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] cursor-pointer">
                <input v-model="includeImages" type="checkbox" class="accent-[#72383D]" />
                <span>Reading Images</span>
              </label>

              <label class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#322D29] bg-white px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] cursor-pointer">
                <input v-model="includeComics" type="checkbox" class="accent-[#72383D]" />
                <span>Comics at End</span>
              </label>

              <button
                v-if="includeImages && readings"
                type="button"
                :disabled="isGeneratingAll"
                @click="handleGenerateAllMissing"
                class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl border border-[#C5A059] bg-[#FAF6EE] hover:bg-[#F3EBD8] text-[#5C4B25] text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer disabled:opacity-50"
              >
                <div
                  v-if="isGeneratingAll"
                  class="w-3.5 h-3.5 border-2 border-[#72383D] border-t-transparent rounded-full animate-spin"
                ></div>
                <span v-else>🎨</span>
                <span>{{ isGeneratingAll ? 'Generating Art...' : 'Generate Reading Images' }}</span>
              </button>

              <button
                type="button"
                @click="openPrintView(false)"
                class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[#322D29] text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer"
              >
                <span>📄</span>
                <span>Full Page</span>
              </button>

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
