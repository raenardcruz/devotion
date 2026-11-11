<template>
  <Teleport to="body">
    <transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 z-[150] flex items-center justify-center p-4 sm:p-6 overflow-y-auto"
        @keydown.esc="handleClose"
      >
        <!-- Backdrop -->
        <div
          class="fixed inset-0 bg-[#322D29]/50 backdrop-blur-sm transition-opacity"
          @click="handleClose"
        ></div>

        <!-- Dialog Panel -->
        <div
          class="relative w-full max-w-2xl bg-[#FCFAF6] border border-[#D1C7BD] rounded-3xl shadow-2xl overflow-hidden flex flex-col z-10 my-auto animate-modal-enter text-[#322D29]"
          @click.stop
        >
          <!-- Modal Header -->
          <div class="px-6 py-5 border-b border-[#D1C7BD]/60 flex items-center justify-between bg-white/70 backdrop-blur-md">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-2xl bg-gradient-to-tr from-[#72383D]/15 to-[#AC9C8D]/20 border border-[#72383D]/20 flex items-center justify-center text-[#72383D]">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"/>
                  <path d="M6 6h10M6 10h10M6 14h10"/>
                </svg>
              </div>
              <div>
                <h3 class="font-serif text-lg md:text-xl font-bold text-[#72383D]">Download EPUB Edition</h3>
                <p class="text-[11px] text-[#322D29]/60 font-medium">Ready for Apple Books, Kindle, Kobo, and e-readers</p>
              </div>
            </div>

            <button
              @click="handleClose"
              :disabled="isGenerating"
              class="w-8 h-8 rounded-full flex items-center justify-center text-[#322D29]/50 hover:text-[#72383D] hover:bg-[#D1C7BD]/30 transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
              aria-label="Close modal"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="18" x2="6" y2="6"></line>
              </svg>
            </button>
          </div>

          <!-- Modal Body -->
          <div class="p-6 md:p-8 space-y-6 max-h-[75vh] overflow-y-auto">
            <!-- Book Metadata Summary -->
            <div class="bg-white/80 border border-[#D1C7BD]/70 rounded-2xl p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-xs">
              <div class="space-y-0.5">
                <span class="text-[10px] tracking-widest uppercase font-bold text-[#72383D]/80">Church Father Resource</span>
                <h4 class="font-serif text-base font-bold text-[#322D29]">{{ bookTitle }}</h4>
                <p v-if="author || era" class="text-xs text-[#322D29]/70 italic">
                  {{ author }} <span v-if="era" class="font-normal font-sans text-[#322D29]/50">({{ era }})</span>
                </p>
              </div>
              <div class="flex items-center gap-2 self-start sm:self-auto text-[11px] bg-[#EFE9E1]/80 text-[#72383D] font-medium px-3 py-1.5 rounded-xl border border-[#D1C7BD]/50 whitespace-nowrap">
                <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
                </svg>
                <span>{{ statsText }}</span>
              </div>
            </div>

            <!-- Download Scope Selection (if multiple sections exist) -->
            <div v-if="hasMultipleSections" class="space-y-2">
              <label class="block text-xs font-serif font-bold text-[#72383D] uppercase tracking-wider">
                Download Scope
              </label>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                <button
                  type="button"
                  @click="selectedScope = 'all'"
                  :disabled="isGenerating"
                  class="p-3 text-left rounded-xl border transition-all flex items-start gap-2.5 cursor-pointer"
                  :class="selectedScope === 'all'
                    ? 'border-[#72383D] bg-[#72383D]/5 shadow-xs'
                    : 'border-[#D1C7BD]/70 bg-white/70 hover:border-[#72383D]/50'"
                >
                  <div
                    class="w-4 h-4 rounded-full border mt-0.5 flex items-center justify-center shrink-0 transition-colors"
                    :class="selectedScope === 'all' ? 'border-[#72383D] bg-[#72383D]' : 'border-[#D1C7BD]'"
                  >
                    <div v-if="selectedScope === 'all'" class="w-1.5 h-1.5 rounded-full bg-white"></div>
                  </div>
                  <div>
                    <div class="text-xs font-semibold text-[#322D29]">Complete Work</div>
                    <div class="text-[11px] text-[#322D29]/60">All {{ sections.length }} sections ({{ totalChaptersCount }} chapters)</div>
                  </div>
                </button>

                <button
                  type="button"
                  @click="selectedScope = 'current'"
                  :disabled="isGenerating"
                  class="p-3 text-left rounded-xl border transition-all flex items-start gap-2.5 cursor-pointer"
                  :class="selectedScope === 'current'
                    ? 'border-[#72383D] bg-[#72383D]/5 shadow-xs'
                    : 'border-[#D1C7BD]/70 bg-white/70 hover:border-[#72383D]/50'"
                >
                  <div
                    class="w-4 h-4 rounded-full border mt-0.5 flex items-center justify-center shrink-0 transition-colors"
                    :class="selectedScope === 'current' ? 'border-[#72383D] bg-[#72383D]' : 'border-[#D1C7BD]'"
                  >
                    <div v-if="selectedScope === 'current'" class="w-1.5 h-1.5 rounded-full bg-white"></div>
                  </div>
                  <div class="overflow-hidden">
                    <div class="text-xs font-semibold text-[#322D29] truncate">Current Section</div>
                    <div class="text-[11px] text-[#322D29]/60 truncate">{{ activeSectionTitle || 'Selected Section' }}</div>
                  </div>
                </button>
              </div>
            </div>

            <!-- Cover Image Customization Section -->
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-xs font-serif font-bold text-[#72383D] uppercase tracking-wider">
                    Cover Artwork <span class="text-[10px] font-sans font-normal text-[#322D29]/50 normal-case">(Optional)</span>
                  </h4>
                  <p class="text-[11px] text-[#322D29]/60 mt-0.5">
                    Upload a custom image (JPEG, PNG, WebP) or keep our generated classic parchment cover.
                  </p>
                </div>
                <span
                  class="text-[10px] font-medium px-2.5 py-1 rounded-full uppercase tracking-wider border shrink-0"
                  :class="customCoverUrl 
                    ? 'bg-emerald-50 text-emerald-700 border-emerald-200' 
                    : 'bg-[#EFE9E1] text-[#72383D] border-[#D1C7BD]/60'"
                >
                  {{ customCoverUrl ? 'Custom Image' : 'Classic Default' }}
                </span>
              </div>

              <!-- Cover Preview and Upload Zone Container -->
              <div class="grid grid-cols-1 sm:grid-cols-12 gap-5 items-center bg-white/60 border border-[#D1C7BD]/70 rounded-2xl p-4 md:p-5">
                <!-- Book Mockup Preview -->
                <div class="sm:col-span-5 flex flex-col items-center justify-center">
                  <div
                    class="relative w-32 sm:w-36 aspect-[1/1.5] rounded-lg shadow-md border border-[#D1C7BD] overflow-hidden group bg-white flex items-center justify-center transition-transform hover:scale-[1.02] duration-300"
                  >
                    <!-- Spine depth shadow -->
                    <div class="absolute left-0 top-0 bottom-0 w-3 bg-gradient-to-r from-black/25 via-black/10 to-transparent z-10 pointer-events-none"></div>

                    <!-- Cover Image -->
                    <img
                      v-if="previewCoverUrl"
                      :src="previewCoverUrl"
                      alt="EPUB Cover Preview"
                      class="w-full h-full object-cover select-none"
                    />
                    <div v-else class="w-full h-full flex flex-col items-center justify-center p-3 text-center bg-[#FAF7F2]">
                      <div class="w-6 h-6 border-2 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
                      <span class="text-[9px] text-[#72383D] mt-2 font-serif">Rendering cover...</span>
                    </div>

                    <!-- Hover Quick Preview Tooltip -->
                    <div class="absolute inset-0 bg-[#322D29]/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center text-white text-[10px] font-medium p-2 text-center pointer-events-none">
                      Cover preview
                    </div>
                  </div>
                  <span class="text-[10px] text-[#322D29]/50 mt-2 font-mono">1200 × 1800 px</span>
                </div>

                <!-- Upload Drag-and-Drop Area -->
                <div class="sm:col-span-7 flex flex-col gap-3">
                  <div
                    class="border-2 border-dashed rounded-xl p-4 text-center transition-colors cursor-pointer relative"
                    :class="[
                      isDragging 
                        ? 'border-[#72383D] bg-[#72383D]/5' 
                        : 'border-[#D1C7BD] hover:border-[#72383D]/60 hover:bg-white/90',
                      isGenerating ? 'pointer-events-none opacity-50' : ''
                    ]"
                    @dragover.prevent="isDragging = true"
                    @dragleave.prevent="isDragging = false"
                    @drop.prevent="handleFileDrop"
                    @click="triggerFileInput"
                  >
                    <input
                      ref="fileInputRef"
                      type="file"
                      accept="image/png,image/jpeg,image/webp,image/jpg"
                      class="hidden"
                      @change="handleFileSelect"
                      :disabled="isGenerating"
                    />

                    <div class="flex flex-col items-center gap-1.5 py-1">
                      <div class="w-9 h-9 rounded-full bg-[#EFE9E1] text-[#72383D] flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <rect width="18" height="18" x="3" y="3" rx="2" ry="2"/>
                          <circle cx="9" cy="9" r="2"/>
                          <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"/>
                        </svg>
                      </div>
                      <div class="text-xs font-semibold text-[#72383D]">
                        {{ customCoverFile ? 'Replace cover image' : 'Upload custom cover image' }}
                      </div>
                      <p class="text-[11px] text-[#322D29]/50 leading-tight">
                        Drop image here or <span class="text-[#72383D] underline">browse</span>
                      </p>
                      <p class="text-[10px] text-[#322D29]/40 mt-1">PNG, JPG, or WebP</p>
                    </div>
                  </div>

                  <!-- File Information & Reset Controls -->
                  <div v-if="customCoverFile" class="flex items-center justify-between bg-white px-3 py-2 rounded-xl border border-[#D1C7BD]/60 text-xs">
                    <div class="flex items-center gap-2 overflow-hidden pr-2">
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-emerald-600 shrink-0">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                      <span class="truncate font-medium text-[#322D29]">{{ customCoverFile.name }}</span>
                      <span class="text-[10px] text-[#322D29]/40 shrink-0">({{ formattedFileSize }})</span>
                    </div>
                    <button
                      type="button"
                      @click="removeCustomCover"
                      :disabled="isGenerating"
                      class="text-[11px] text-red-600 hover:text-red-700 hover:underline shrink-0 font-medium ml-2 cursor-pointer"
                    >
                      Reset to Default
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Generation Progress / Status Bar -->
            <div v-if="isGenerating" class="space-y-2 bg-[#FAF7F2] border border-[#D1C7BD] rounded-2xl p-4 animate-fade-in">
              <div class="flex items-center justify-between text-xs font-medium">
                <span class="text-[#72383D] flex items-center gap-2">
                  <div class="w-3.5 h-3.5 border-2 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
                  {{ progressStatus || 'Generating EPUB...' }}
                </span>
                <span class="font-mono text-[#322D29]/70 font-semibold">{{ progressPercent }}%</span>
              </div>
              <div class="w-full bg-[#E5DFD6] rounded-full h-2 overflow-hidden">
                <div
                  class="bg-gradient-to-r from-[#72383D] to-[#AC9C8D] h-2 rounded-full transition-all duration-300"
                  :style="{ width: `${progressPercent}%` }"
                ></div>
              </div>
            </div>

            <!-- Error Banner -->
            <div v-if="errorMessage" class="bg-red-50 border border-red-200 rounded-xl p-3.5 text-xs text-red-700 flex items-start gap-2.5">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="shrink-0 mt-0.5">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              <div>
                <strong class="font-semibold">Generation failed:</strong> {{ errorMessage }}
              </div>
            </div>

            <!-- Success Notification -->
            <div v-if="downloadSuccess" class="bg-emerald-50 border border-emerald-200 rounded-xl p-3.5 text-xs text-emerald-800 flex items-center gap-2.5">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-emerald-600 shrink-0">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
              <div>
                <strong>EPUB Ready!</strong> Your download has started.
              </div>
            </div>
          </div>

          <!-- Modal Footer Actions -->
          <div class="px-6 py-4 border-t border-[#D1C7BD]/60 bg-white/70 backdrop-blur-md flex items-center justify-between gap-3">
            <button
              type="button"
              @click="handleClose"
              :disabled="isGenerating"
              class="px-5 py-2.5 rounded-full text-xs font-semibold text-[#322D29]/70 hover:text-[#322D29] hover:bg-[#D1C7BD]/30 transition-all cursor-pointer disabled:opacity-40"
            >
              Close
            </button>

            <button
              type="button"
              @click="startDownload"
              :disabled="isGenerating"
              class="px-6 py-2.5 rounded-full bg-gradient-to-r from-[#72383D] to-[#322D29] text-white text-xs font-semibold tracking-wide hover:from-[#8B464C] hover:to-[#453E38] transition-all shadow-sm hover:shadow active:scale-95 flex items-center gap-2 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg v-if="!isGenerating" xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              <div v-else class="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              <span>{{ isGenerating ? 'Generating EPUB...' : 'Download (.epub)' }}</span>
            </button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import {
  generateDefaultCover,
  processUploadedCover,
  generateEpub,
  downloadFile,
  type SectionData
} from '../../utils/epubGenerator';

const props = defineProps<{
  isOpen: boolean;
  bookTitle: string;
  author: string;
  era?: string;
  sections: SectionData[];
  currentSectionId?: string;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'downloaded'): void;
}>();

// State
const selectedScope = ref<'all' | 'current'>('all');
const isDragging = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

const defaultCoverBlob = ref<Blob | null>(null);
const defaultCoverUrl = ref<string>('');

const customCoverFile = ref<File | null>(null);
const customCoverBlob = ref<Blob | null>(null);
const customCoverUrl = ref<string>('');

const isGenerating = ref(false);
const progressPercent = ref(0);
const progressStatus = ref('');
const errorMessage = ref('');
const downloadSuccess = ref(false);

// Active section title if downloading single section
const activeSectionTitle = computed(() => {
  if (!props.currentSectionId || !props.sections) return '';
  const found = props.sections.find(s => s.id === props.currentSectionId);
  return found ? found.title : '';
});

const hasMultipleSections = computed(() => {
  return props.sections && props.sections.length > 1;
});

const totalChaptersCount = computed(() => {
  if (!props.sections) return 0;
  return props.sections.reduce((acc, s) => acc + s.chapters.length, 0);
});

const statsText = computed(() => {
  if (selectedScope.value === 'current') {
    const sec = props.sections.find(s => s.id === props.currentSectionId);
    const count = sec ? sec.chapters.length : 0;
    return `1 Section • ${count} ${count === 1 ? 'Chapter' : 'Chapters'}`;
  }
  const secCount = props.sections.length;
  const chapCount = totalChaptersCount.value;
  return `${secCount} ${secCount === 1 ? 'Section' : 'Sections'} • ${chapCount} Chapters`;
});

const previewCoverUrl = computed(() => {
  return customCoverUrl.value || defaultCoverUrl.value;
});

const formattedFileSize = computed(() => {
  if (!customCoverFile.value) return '';
  const bytes = customCoverFile.value.size;
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
});

// Render or regenerate default cover when metadata changes
const refreshDefaultCover = async () => {
  if (!props.bookTitle) return;
  try {
    const result = await generateDefaultCover(props.bookTitle, props.author, props.era);
    defaultCoverBlob.value = result.blob;
    defaultCoverUrl.value = result.dataUrl;
  } catch (err) {
    console.error('Failed to generate default cover:', err);
  }
};

watch(() => props.isOpen, (open) => {
  if (open) {
    errorMessage.value = '';
    downloadSuccess.value = false;
    refreshDefaultCover();
  }
});

watch(() => props.bookTitle, () => {
  if (props.isOpen) {
    refreshDefaultCover();
  }
});

onMounted(() => {
  if (props.isOpen) {
    refreshDefaultCover();
  }
});

// File input handlers
const triggerFileInput = () => {
  fileInputRef.value?.click();
};

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file) {
    applyUploadedFile(file);
  }
};

const handleFileDrop = (event: DragEvent) => {
  isDragging.value = false;
  const file = event.dataTransfer?.files?.[0];
  if (file) {
    applyUploadedFile(file);
  }
};

const applyUploadedFile = async (file: File) => {
  errorMessage.value = '';
  if (!file.type.startsWith('image/')) {
    errorMessage.value = 'Please select a valid image file (JPEG, PNG, or WebP).';
    return;
  }

  try {
    const processed = await processUploadedCover(file);
    customCoverFile.value = file;
    customCoverBlob.value = processed.blob;
    customCoverUrl.value = processed.dataUrl;
  } catch (err: any) {
    errorMessage.value = `Failed to process image: ${err.message || 'Unknown error'}`;
  }
};

const removeCustomCover = () => {
  customCoverFile.value = null;
  customCoverBlob.value = null;
  customCoverUrl.value = '';
  if (fileInputRef.value) {
    fileInputRef.value.value = '';
  }
};

const handleClose = () => {
  if (isGenerating.value) return;
  emit('close');
};

const startDownload = async () => {
  if (isGenerating.value) return;

  isGenerating.value = true;
  progressPercent.value = 0;
  progressStatus.value = 'Initializing...';
  errorMessage.value = '';
  downloadSuccess.value = false;

  try {
    let sectionsToExport = props.sections;
    let exportTitle = props.bookTitle;

    if (selectedScope.value === 'current' && props.currentSectionId) {
      const match = props.sections.find(s => s.id === props.currentSectionId);
      if (match) {
        sectionsToExport = [match];
        if (match.title && match.title !== 'default') {
          exportTitle = `${props.bookTitle} - ${match.title}`;
        }
      }
    }

    const coverToUse = customCoverBlob.value || defaultCoverBlob.value;

    const epubBlob = await generateEpub({
      title: exportTitle,
      author: props.author,
      era: props.era,
      sections: sectionsToExport,
      coverBlob: coverToUse,
      onProgress: (percent, status) => {
        progressPercent.value = percent;
        progressStatus.value = status;
      }
    });

    // Create safe filename
    const safeTitle = exportTitle
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '_')
      .replace(/^_+|_+$/g, '');
    const filename = `${safeTitle || 'church_resource'}.epub`;

    downloadFile(epubBlob, filename);
    downloadSuccess.value = true;
    emit('downloaded');

    // Auto dismiss after brief success display
    setTimeout(() => {
      if (downloadSuccess.value && !isGenerating.value) {
        handleClose();
      }
    }, 2200);
  } catch (err: any) {
    console.error('EPUB Generation Error:', err);
    errorMessage.value = err.message || 'An error occurred while generating the EPUB file.';
  } finally {
    isGenerating.value = false;
  }
};
</script>

<style scoped>
@keyframes modalEnter {
  from {
    opacity: 0;
    transform: scale(0.96) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.animate-modal-enter {
  animation: modalEnter 0.25s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-fade-in {
  animation: fadeIn 0.2s ease-out forwards;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
