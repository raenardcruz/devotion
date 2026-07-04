<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import AppButton from '../components/common/AppButton.vue';
import writingsIndex from '../components/resources/index.json';

interface Chapter {
  id: string;
  numeral: string;
  title: string;
  content: string[];
}

interface Section {
  id: string;
  title: string;
  chapters: Chapter[];
}

interface Document {
  id: string;
  title: string;
  sections: Section[];
}

// Vite dynamic imports map
const writingModules = import.meta.glob('../components/resources/*.md', { query: '?raw', import: 'default' });

// State
const documents = ref<Document[]>([]);

// Get initial writing from localStorage or default
const getInitialWritingId = (): string => {
  const savedDocId = localStorage.getItem('selected_writing_id');
  const defaultDocId = writingsIndex[0]?.id || 'alexander_of_alexandria';
  return savedDocId && writingsIndex.some(w => w.id === savedDocId)
    ? savedDocId
    : defaultDocId;
};

const selectedDocId = ref<string>(getInitialWritingId());

const getInitialChapterId = (docId: string): string => {
  const savedChapterId = localStorage.getItem('selected_chapter_id');
  if (savedChapterId && savedChapterId.startsWith(docId)) {
    return savedChapterId;
  }
  return '';
};

const selectedChapterId = ref<string>(getInitialChapterId(selectedDocId.value));
const readChapters = ref<Record<string, boolean>>({});
const searchQuery = ref<string>('');
const writingSearchQuery = ref<string>('');
const isSelectorOpen = ref<boolean>(false);
const isLoading = ref<boolean>(false);
const errorMsg = ref<string>('');

// Parse resource markdown
const parseMarkdown = (docId: string, rawText: string): Document => {
  const lines = rawText.split('\n');
  let title = '';
  const sections: Section[] = [];
  let currentSection: Section | null = null;
  let currentChapter: Chapter | null = null;

  for (let i = 0; i < lines.length; i++) {
    const rawLine = lines[i];
    if (rawLine === undefined) continue;
    const line = rawLine.trim();
    if (!line) continue;

    if (line.startsWith('# ')) {
      title = line.substring(2).trim();
      continue;
    }

    if (line.startsWith('## ')) {
      if (currentChapter && currentSection) {
        currentSection.chapters.push(currentChapter);
        currentChapter = null;
      }
      
      const secTitle = line.substring(3).trim();
      const secId = secTitle.toLowerCase().replace(/[^a-z0-9]+/g, '_');
      
      currentSection = {
        id: secId,
        title: secTitle,
        chapters: []
      };
      sections.push(currentSection);
      continue;
    }

    const chapMatch = line.match(/^#### Chapter\s+([IVXLCDM\d]+):\s*(.*)/i);
    if (chapMatch) {
      if (currentChapter) {
        if (!currentSection) {
          currentSection = {
            id: 'default',
            title: title || 'Content',
            chapters: []
          };
          sections.push(currentSection);
        }
        currentSection.chapters.push(currentChapter);
      }
      
      const num = chapMatch[1] || '';
      const chTitle = (chapMatch[2] || '').trim();
      const chId = `${docId}_${currentSection ? currentSection.id : 'default'}_ch_${num.toLowerCase()}`;
      
      currentChapter = {
        id: chId,
        numeral: num,
        title: chTitle,
        content: []
      };
      continue;
    }

    if (currentChapter) {
      currentChapter.content.push(line);
    }
  }

  if (currentChapter) {
    if (!currentSection) {
      currentSection = {
        id: 'default',
        title: title || 'Content',
        chapters: []
      };
      sections.push(currentSection);
    }
    currentSection.chapters.push(currentChapter);
  }

  return {
    id: docId,
    title: title || 'Resource',
    sections
  };
};

// Dynamic loading function
const loadWriting = async (docId: string) => {
  // Check if we already have it loaded
  const existing = documents.value.find(d => d.id === docId);
  if (existing) {
    return existing;
  }

  isLoading.value = true;
  errorMsg.value = '';
  try {
    const path = `../components/resources/${docId}.md`;
    const loadFn = writingModules[path];
    if (!loadFn) {
      throw new Error(`Writing file for ID '${docId}' not found.`);
    }
    const rawContent = await loadFn() as string;
    const doc = parseMarkdown(docId, rawContent);
    documents.value.push(doc);
    return doc;
  } catch (err: any) {
    console.error(err);
    errorMsg.value = `Failed to load document: ${err.message}`;
    return null;
  } finally {
    isLoading.value = false;
  }
};

const selectWriting = async (docId: string) => {
  selectedDocId.value = docId;
  isSelectorOpen.value = false;
  writingSearchQuery.value = '';
  const doc = await loadWriting(docId);
  if (doc && doc.sections.length > 0) {
    const firstSection = doc.sections[0];
    if (firstSection && firstSection.chapters.length > 0) {
      const firstChapter = firstSection.chapters[0];
      if (firstChapter) {
        selectedChapterId.value = firstChapter.id;
      }
    }
  }
};

// Initialize
onMounted(async () => {
  // Restore read progress
  const stored = localStorage.getItem('resource_read_progress');
  if (stored) {
    try {
      readChapters.value = JSON.parse(stored);
    } catch (e) {
      console.error('Error parsing reading progress', e);
    }
  }

  // Load the initial writing
  const doc = await loadWriting(selectedDocId.value);
  
  // Set default chapter if none was restored
  if (!selectedChapterId.value && doc && doc.sections.length > 0) {
    const firstSection = doc.sections[0];
    if (firstSection && firstSection.chapters.length > 0) {
      const firstChapter = firstSection.chapters[0];
      if (firstChapter) {
        selectedChapterId.value = firstChapter.id;
      }
    }
  }
});

// Watch writing and chapter changes to persist progress
watch(selectedDocId, (newDocId) => {
  localStorage.setItem('selected_writing_id', newDocId);
});

watch(selectedChapterId, (newChapterId) => {
  if (newChapterId) {
    localStorage.setItem('selected_chapter_id', newChapterId);
  }
});

const activeDoc = computed(() => {
  return documents.value.find(d => d.id === selectedDocId.value) || null;
});

const selectedWritingTitle = computed(() => {
  const writing = writingsIndex.find(w => w.id === selectedDocId.value);
  if (writing) {
    return writing.years ? `${writing.title} (${writing.years})` : writing.title;
  }
  return 'Select Writing';
});

const filteredWritings = computed(() => {
  if (!writingSearchQuery.value.trim()) return writingsIndex;
  const query = writingSearchQuery.value.toLowerCase();
  return writingsIndex.filter(w => 
    w.title.toLowerCase().includes(query) || 
    (w.years && w.years.toLowerCase().includes(query))
  );
});

const allDocChapters = computed(() => {
  if (!activeDoc.value) return [];
  const list: Chapter[] = [];
  activeDoc.value.sections.forEach(s => {
    list.push(...s.chapters);
  });
  return list;
});

const activeChapter = computed(() => {
  if (!activeDoc.value) return null;
  for (const s of activeDoc.value.sections) {
    const found = s.chapters.find(c => c.id === selectedChapterId.value);
    if (found) return found;
  }
  return null;
});

const flatChapters = computed(() => {
  return allDocChapters.value;
});

const currentChapterIndex = computed(() => {
  return flatChapters.value.findIndex(c => c.id === selectedChapterId.value);
});

const filteredSections = computed(() => {
  if (!activeDoc.value) return [];
  if (!searchQuery.value.trim()) return activeDoc.value.sections;

  const query = searchQuery.value.toLowerCase();
  return activeDoc.value.sections.map(s => {
    const filteredChapters = s.chapters.filter(c => 
      c.title.toLowerCase().includes(query) || 
      c.numeral.toLowerCase().includes(query)
    );
    return {
      ...s,
      chapters: filteredChapters
    };
  }).filter(s => s.chapters.length > 0);
});

const docProgress = computed(() => {
  const chapters = allDocChapters.value;
  if (chapters.length === 0) return 0;
  const readCount = chapters.filter(c => readChapters.value[c.id]).length;
  return Math.round((readCount / chapters.length) * 100);
});

const readCountText = computed(() => {
  const chapters = allDocChapters.value;
  const readCount = chapters.filter(c => readChapters.value[c.id]).length;
  return `${readCount} of ${chapters.length} chapters`;
});

// Methods
const toggleReadStatus = (chapterId: string) => {
  readChapters.value[chapterId] = !readChapters.value[chapterId];
  localStorage.setItem('resource_read_progress', JSON.stringify(readChapters.value));
};

const selectChapter = (chapterId: string) => {
  selectedChapterId.value = chapterId;
  const reader = document.getElementById('reader-content');
  if (reader) {
    if (window.innerWidth < 1024) {
      setTimeout(() => {
        reader.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }, 100);
    } else {
      reader.scrollTop = 0;
    }
  }
};

const prevChapter = () => {
  const idx = currentChapterIndex.value;
  if (idx > 0) {
    const prevChap = flatChapters.value[idx - 1];
    if (prevChap) {
      selectChapter(prevChap.id);
    }
  }
};

const nextChapter = () => {
  const idx = currentChapterIndex.value;
  if (idx >= 0 && idx < flatChapters.value.length - 1) {
    const nextChap = flatChapters.value[idx + 1];
    if (nextChap) {
      selectChapter(nextChap.id);
    }
  }
};
</script>

<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-24 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Container -->
    <div class="relative z-10 max-w-7xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">
      <!-- Header Area -->
      <header class="text-center mb-8 animate-fade-in-down">
        <h1 class="text-3xl md:text-5xl font-serif text-parchment-primary-dark mb-1">Early Church Resources</h1>
        <p class="text-parchment-neutral/50 text-xs uppercase tracking-[0.25em] font-bold">Timeless writings of faith and theology</p>
      </header>

      <!-- Searchable Selector for Writings -->
      <div class="relative max-w-md mx-auto mb-8 w-full animate-fade-in-down z-50">
        <!-- Backdrop overlay for dropdown closing -->
        <div v-if="isSelectorOpen" class="fixed inset-0 z-40" @click="isSelectorOpen = false"></div>

        <!-- Trigger Button -->
        <button 
          @click="isSelectorOpen = !isSelectorOpen"
          class="relative z-50 w-full bg-parchment-neutral-light border border-parchment-border hover:border-parchment-primary rounded-2xl py-3.5 px-5 shadow-sm hover:shadow-md flex items-center justify-between transition-all duration-300 group outline-none"
        >
          <div class="flex items-center gap-3">
            <span class="text-parchment-primary">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"/>
                <path d="M6 6h10M6 10h10M6 14h10"/>
              </svg>
            </span>
            <div class="text-left">
              <span class="text-[10px] text-parchment-neutral/50 tracking-widest uppercase font-bold block">Selected Writing</span>
              <span class="font-serif text-sm md:text-base text-parchment-primary-dark font-semibold group-hover:text-parchment-primary transition-colors">
                {{ selectedWritingTitle }}
              </span>
            </div>
          </div>
          <span class="text-parchment-neutral/40 group-hover:text-parchment-primary transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="transition-transform duration-300" :class="{'rotate-180': isSelectorOpen}">
              <path d="m6 9 6 6 6-6"/>
            </svg>
          </span>
        </button>

        <!-- Dropdown Card -->
        <div 
          v-if="isSelectorOpen"
          class="absolute left-0 right-0 mt-2 bg-parchment-bg border border-parchment-border rounded-2xl shadow-xl z-50 p-4 max-h-[350px] flex flex-col gap-3 animate-fade-in"
        >
          <!-- Filter Search -->
          <div class="relative">
            <input 
              v-model="writingSearchQuery"
              type="text"
              placeholder="Search writings / authors..."
              class="w-full pl-9 pr-4 py-2.5 rounded-xl border border-parchment-border bg-parchment-neutral-light text-parchment-neutral placeholder-parchment-neutral/40 focus:border-parchment-primary shadow-inner transition-all outline-none text-xs"
            >
            <span class="absolute left-3 top-3 text-parchment-neutral/30">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
              </svg>
            </span>
          </div>

          <!-- Writings List -->
          <div class="flex-grow overflow-y-auto space-y-1 pr-1">
            <button
              v-for="writing in filteredWritings"
              :key="writing.id"
              @click="selectWriting(writing.id)"
              class="w-full text-left text-xs py-2.5 px-3.5 rounded-xl transition-colors flex items-center justify-between border-none outline-none"
              :class="selectedDocId === writing.id 
                ? 'bg-parchment-primary text-white font-semibold' 
                : 'text-parchment-neutral/80 hover:bg-parchment-neutral-light hover:text-parchment-primary-dark'"
            >
              <span class="truncate pr-4">
                {{ writing.title }}
                <span v-if="writing.years" class="text-[10px] font-sans ml-1.5" :class="selectedDocId === writing.id ? 'text-white/80' : 'text-parchment-neutral/75'">
                  ({{ writing.years }})
                </span>
              </span>
              <span v-if="selectedDocId === writing.id" class="text-white">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </span>
            </button>
            <div v-if="filteredWritings.length === 0" class="py-6 text-center text-xs text-parchment-neutral/40 italic">
              No writings found matching "{{ writingSearchQuery }}"
            </div>
          </div>
        </div>
      </div>

      <!-- Loading / Error States -->
      <div v-if="isLoading" class="flex-grow flex flex-col items-center justify-center py-20 gap-4">
        <div class="w-12 h-12 border-4 border-parchment-primary/30 border-t-parchment-primary rounded-full animate-spin"></div>
        <p class="font-serif italic text-parchment-neutral/60 text-sm">Loading church father writings...</p>
      </div>

      <div v-else-if="errorMsg" class="flex-grow flex flex-col items-center justify-center py-20 text-center max-w-md mx-auto gap-4">
        <span class="text-red-500">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
        </span>
        <h3 class="font-serif text-lg font-bold text-parchment-neutral">Failed to Load Content</h3>
        <p class="text-xs text-parchment-neutral/60">{{ errorMsg }}</p>
        <AppButton variant="primary" @click="loadWriting(selectedDocId)">Retry Loading</AppButton>
      </div>

      <!-- Main Layout -->
      <main v-else class="flex-grow grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Sidebar Navigation -->
        <aside class="lg:col-span-1 h-fit lg:sticky lg:top-24 flex flex-col gap-6 animate-fade-in-up delay-100">
          <!-- Search & Progress Card -->
          <div class="bg-parchment-neutral-light border border-parchment-border rounded-3xl p-6 shadow-sm flex flex-col gap-4">
            <!-- Progress Tracker -->
            <div class="border-b border-parchment-border/40 pb-4">
              <h3 class="font-serif text-xs text-parchment-primary-dark/80 tracking-widest uppercase font-bold mb-2">Reading Progress</h3>
              <div class="flex justify-between text-[11px] text-parchment-neutral/60 mb-1.5">
                <span>Completed</span>
                <span>{{ readCountText }} ({{ docProgress }}%)</span>
              </div>
              <div class="w-full bg-parchment-bg border border-parchment-border/60 rounded-full h-1.5 overflow-hidden">
                <div class="bg-parchment-primary h-1.5 rounded-full transition-all duration-500" :style="{ width: `${docProgress}%` }"></div>
              </div>
            </div>

            <!-- Search -->
            <div>
              <h3 class="font-serif text-xs text-parchment-neutral/50 tracking-widest uppercase font-bold mb-2">Search Chapters</h3>
              <div class="relative">
                <input 
                  v-model="searchQuery" 
                  type="text" 
                  placeholder="Type title..." 
                  class="w-full pl-3 pr-4 py-2.5 rounded-xl border border-parchment-border bg-parchment-bg text-parchment-neutral placeholder-parchment-neutral/30 focus:bg-parchment-neutral-light focus:border-parchment-primary shadow-inner transition-all outline-none text-xs"
                >
              </div>
            </div>
          </div>

          <!-- Table of Contents List -->
          <div class="bg-parchment-neutral-light border border-parchment-border rounded-3xl p-6 shadow-sm lg:max-h-[50vh] lg:overflow-y-auto scrollbar-hide">
            <h3 class="font-serif text-xs mb-4 text-parchment-primary-dark/80 tracking-widest uppercase font-bold border-b border-parchment-border/40 pb-2.5">Chapters</h3>
            
            <div class="space-y-4">
              <div v-for="section in filteredSections" :key="section.id" class="space-y-2">
                <!-- Section Header (if not default) -->
                <h4 v-if="section.id !== 'default'" class="text-[11px] font-bold text-parchment-neutral/70 mt-3 font-serif italic border-l-2 border-parchment-primary/30 pl-2">
                  {{ section.title }}
                </h4>
                
                <div class="space-y-1">
                  <button
                    v-for="chapter in section.chapters"
                    :key="chapter.id"
                    @click="selectChapter(chapter.id)"
                    class="w-full text-left text-xs transition-colors py-2 px-3 rounded-lg flex items-center justify-between group hover:bg-parchment-bg border-none shadow-none"
                    :class="selectedChapterId === chapter.id 
                      ? 'bg-parchment-primary/10 text-parchment-primary-dark font-semibold border-l-4 border-parchment-primary pl-2' 
                      : readChapters[chapter.id] ? 'text-parchment-neutral/40 hover:text-parchment-neutral/60' : 'text-parchment-neutral/80 hover:text-parchment-neutral'"
                  >
                    <span class="truncate flex-1">
                      <span class="font-serif font-bold text-parchment-neutral/40 group-hover:text-parchment-primary-dark mr-1.5">{{ chapter.numeral }}</span>
                      {{ chapter.title }}
                    </span>
                    <!-- Read Checkbox -->
                    <span 
                      @click.stop="toggleReadStatus(chapter.id)" 
                      class="ml-2 w-4 h-4 rounded border transition-colors flex items-center justify-center cursor-pointer"
                      :class="readChapters[chapter.id] 
                        ? 'bg-parchment-primary/20 border-parchment-primary text-parchment-primary-dark' 
                        : 'border-parchment-border hover:border-parchment-primary text-transparent'"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    </span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </aside>

        <!-- Content Area -->
        <div class="lg:col-span-3 flex flex-col gap-6 animate-fade-in-up delay-150">
          <div class="bg-parchment-neutral-light border border-parchment-border rounded-3xl p-6 md:p-8 flex-grow flex flex-col min-h-[60vh] relative shadow-sm">
            <!-- Reader Area -->
            <div v-if="activeChapter" id="reader-content" class="flex-grow lg:overflow-y-auto space-y-6 pr-2">
              <div class="border-b border-parchment-border/40 pb-4 mb-6">
                <div class="text-[10px] text-parchment-primary font-bold tracking-widest uppercase mb-1">
                  Chapter {{ activeChapter.numeral }}
                </div>
                <h2 class="text-2xl md:text-3xl font-serif text-parchment-neutral font-medium">
                  {{ activeChapter.title }}
                </h2>
              </div>

              <!-- Main text paragraphs -->
              <div class="space-y-6 text-parchment-neutral/90 leading-relaxed font-serif text-base md:text-lg">
                <p v-for="(para, idx) in activeChapter.content" :key="idx" class="indent-6 text-justify">
                  {{ para }}
                </p>
              </div>

              <!-- Checkbox Mark as Read -->
              <div class="pt-8 border-t border-parchment-border/30 mt-8 flex justify-center">
                <AppButton
                  :variant="readChapters[activeChapter.id] ? 'primary' : 'outlined'"
                  @click="toggleReadStatus(activeChapter.id)"
                  custom-class="!px-6 !py-2.5 flex items-center gap-2.5"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="transition-transform duration-300" :class="{'scale-110': readChapters[activeChapter.id]}">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                  <span>{{ readChapters[activeChapter.id] ? 'Completed' : 'Mark as Read' }}</span>
                </AppButton>
              </div>
            </div>

            <div v-else class="flex-grow flex items-center justify-center text-parchment-neutral/40 font-serif italic text-base">
              Select a chapter to begin reading
            </div>

            <!-- Navigation Buttons -->
            <div class="mt-8 pt-6 border-t border-parchment-border/40 flex justify-between items-center text-xs">
              <AppButton 
                variant="outlined"
                @click="prevChapter" 
                :disabled="currentChapterIndex <= 0"
                custom-class="!px-4 !py-1.5"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m15 18-6-6 6-6"/>
                </svg>
                <span>Previous</span>
              </AppButton>

              <AppButton 
                variant="outlined"
                @click="nextChapter" 
                :disabled="currentChapterIndex >= flatChapters.length - 1"
                custom-class="!px-4 !py-1.5"
              >
                <span>Next</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m9 18 6-6-6-6"/>
                </svg>
              </AppButton>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- Global Footer -->
    <BottomNav />
  </div>
</template>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.indent-6 {
  text-indent: 1.5rem;
}
.animate-fade-in-down {
  animation: fadeInDown 0.6s ease-out forwards;
}
.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out forwards;
}
.animate-fade-in {
  animation: fadeIn 0.2s ease-out forwards;
}
.delay-100 {
  animation-delay: 0.1s;
  animation-fill-mode: both;
}
.delay-150 {
  animation-delay: 0.15s;
  animation-fill-mode: both;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-15px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(15px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
