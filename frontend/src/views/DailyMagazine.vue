<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useDate } from '../composables/useDate';
import { useDevotionApi } from '../composables/useDevotionApi';
import { useMagazineImages, type MagazineReadingKey } from '../composables/useMagazineImages';
import { useStoryToManga } from '../composables/useStoryToManga';
import MagazineDocument, { type MagazineMassReadings } from '../components/devotion/MagazineDocument.vue';

const route = useRoute();
const router = useRouter();
const { getLocalISOString } = useDate();
const { getDevotion, loading, error } = useDevotionApi();

const selectedDate = computed(() => {
  const q = route.query.date;
  return typeof q === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(q) ? q : getLocalISOString();
});

const includeCommentaries = ref(route.query.commentaries !== '0');
const includeImages = ref(route.query.images !== '0');
const includeComics = ref(route.query.comics !== '0');

const readings = ref<MagazineMassReadings | null>(null);

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

const formattedDate = computed(() => {
  const d = new Date(`${selectedDate.value}T12:00:00`);
  return d.toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  });
});

const waitForImagesAndFonts = async () => {
  if (typeof document !== 'undefined' && document.fonts?.ready) {
    await document.fonts.ready;
  }
  await nextTick();
  const imgs = Array.from(document.querySelectorAll('img'));
  await Promise.all(
    imgs.map((img) => {
      if (img.complete) return Promise.resolve();
      return new Promise<void>((resolve) => {
        img.addEventListener('load', () => resolve(), { once: true });
        img.addEventListener('error', () => resolve(), { once: true });
      });
    })
  );
};

const triggerPrint = async () => {
  await waitForImagesAndFonts();
  window.print();
};

const loadMagazine = async () => {
  const date = selectedDate.value;
  try {
    const [data] = await Promise.all([
      getDevotion<MagazineMassReadings>(date),
      loadImagesForDate(date),
      loadStateFromServer(date, true),
    ]);
    readings.value = data;

    if (route.query.print === 'true') {
      setTimeout(() => {
        triggerPrint();
      }, 350);
    }
  } catch (err) {
    console.error('Error loading daily magazine:', err);
  }
};

const handleGenerateSingleImage = async (payload: {
  key: MagazineReadingKey;
  citation: string;
  text?: string;
  context?: string;
  force?: boolean;
}) => {
  await generateReadingImage(
    selectedDate.value,
    payload.key,
    payload.citation,
    payload.text,
    payload.context,
    Boolean(payload.force)
  );
};

const handleGenerateAllMissing = async () => {
  if (!readings.value) return;
  await generateMissingImages(selectedDate.value, readings.value);
};

onMounted(() => {
  loadMagazine();
});
</script>

<template>
  <div class="min-h-screen bg-[#EFEAE2] print:bg-white py-6 print:py-0 px-3 sm:px-6 print:px-0">
    <!-- Non-printable Editorial Toolbar -->
    <div class="print:hidden max-w-[210mm] mx-auto mb-5 bg-[#FDFBF7] border border-[#D1C7BD] rounded-2xl p-3.5 shadow-sm flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center gap-2.5">
        <button
          type="button"
          @click="router.push('/daily-readings')"
          class="px-3 py-1.5 rounded-xl border border-[#D1C7BD] bg-white hover:bg-[#EFE9E1] text-[#322D29] text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer"
        >
          ← Daily Readings
        </button>
        <span class="text-xs font-serif font-bold text-[#72383D]">
          Lumen Liturgiae · {{ formattedDate }}
        </span>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <label class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#322D29] bg-white px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] cursor-pointer">
          <input v-model="includeCommentaries" type="checkbox" class="accent-[#72383D]" />
          <span>AI Commentaries</span>
        </label>
        <label class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#322D29] bg-white px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] cursor-pointer">
          <input v-model="includeImages" type="checkbox" class="accent-[#72383D]" />
          <span>Images</span>
        </label>
        <label class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#322D29] bg-white px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] cursor-pointer">
          <input v-model="includeComics" type="checkbox" class="accent-[#72383D]" />
          <span>Comics</span>
        </label>

        <button
          v-if="includeImages && readings"
          type="button"
          :disabled="isGeneratingAll"
          @click="handleGenerateAllMissing"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-[#C5A059] bg-[#FAF6EE] hover:bg-[#F3EBD8] text-[#5C4B25] text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer disabled:opacity-50"
        >
          <span>{{ isGeneratingAll ? 'Generating Art...' : '🎨 Generate Images' }}</span>
        </button>

        <button
          type="button"
          @click="triggerPrint"
          class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-xl bg-gradient-to-r from-[#72383D] to-[#542529] text-white text-xs font-bold uppercase tracking-wider shadow-xs cursor-pointer"
        >
          <span>🖨️</span>
          <span>Print / Save PDF</span>
        </button>
      </div>

      <div
        v-if="imageError"
        class="w-full px-3 py-2 rounded-lg bg-red-50 border border-red-200 text-xs text-red-800"
      >
        {{ imageError }}
      </div>
    </div>

    <!-- Loading / Error States -->
    <div v-if="loading && !readings" class="max-w-[210mm] mx-auto py-24 text-center">
      <div class="w-10 h-10 border-4 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin mx-auto mb-3"></div>
      <p class="text-sm font-serif text-[#5C5149]">Typesetting daily magazine edition...</p>
    </div>

    <div v-else-if="error" class="max-w-[210mm] mx-auto bg-red-50 border border-red-200 rounded-2xl p-6 text-center text-red-800">
      <p class="font-bold mb-1">Could not load magazine</p>
      <p class="text-xs">{{ error }}</p>
    </div>

    <!-- Magazine Document -->
    <MagazineDocument
      v-else-if="readings"
      :date="selectedDate"
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
</template>

<style>
@page {
  size: A4 portrait;
  margin: 14mm 15mm 16mm 15mm;
}

@media print {
  html,
  body {
    background: #ffffff !important;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
}
</style>
