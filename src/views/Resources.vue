<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import ignatiusRaw from '../components/resources/letter_of_ignatius.md?raw';
import apologyRaw from '../components/resources/apology_of_justin_martyr.md?raw';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import AppButton from '../components/common/AppButton.vue';
import AppTabs from '../components/common/AppTabs.vue';

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

// State
const documents = ref<Document[]>([]);
const selectedDocId = ref<string>('ignatius');
const selectedChapterId = ref<string>('');
const readChapters = ref<Record<string, boolean>>({});
const searchQuery = ref<string>('');

// Parse both resources on load
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

// Initialize
onMounted(() => {
  documents.value = [
    parseMarkdown('ignatius', ignatiusRaw),
    parseMarkdown('justin', apologyRaw)
  ];
  
  if (documents.value.length > 0) {
    const firstDoc = documents.value[0];
    if (firstDoc && firstDoc.sections.length > 0) {
      const firstSection = firstDoc.sections[0];
      if (firstSection && firstSection.chapters.length > 0) {
        const firstChapter = firstSection.chapters[0];
        if (firstChapter) {
          selectedChapterId.value = firstChapter.id;
        }
      }
    }
  }

  const stored = localStorage.getItem('resource_read_progress');
  if (stored) {
    try {
      readChapters.value = JSON.parse(stored);
    } catch (e) {
      console.error('Error parsing reading progress', e);
    }
  }
});

// Watch document selection to update active chapter
watch(selectedDocId, (newDocId) => {
  const doc = documents.value.find(d => d.id === newDocId);
  if (doc && doc.sections.length > 0) {
    const firstSection = doc.sections[0];
    if (firstSection && firstSection.chapters.length > 0) {
      const firstChapter = firstSection.chapters[0];
      if (firstChapter) {
        selectedChapterId.value = firstChapter.id;
      }
    }
  }
});

const activeDoc = computed(() => {
  return documents.value.find(d => d.id === selectedDocId.value) || null;
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

const docTabs = [
  { id: 'ignatius', label: 'Letters of Ignatius' },
  { id: 'justin', label: 'First Apology of Justin' }
];
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

      <!-- Document Switcher Tabs using AppTabs -->
      <AppTabs :tabs="docTabs" v-model="selectedDocId" class="animate-fade-in-down" />

      <!-- Main Layout -->
      <main class="flex-grow grid grid-cols-1 lg:grid-cols-4 gap-8">
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
                  {{ section.title.replace('The Epistle of Ignatius to the ', '') }}
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
</style>
