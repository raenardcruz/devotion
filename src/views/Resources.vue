<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import ignatiusRaw from '../components/resources/letter_of_ignatius.md?raw';
import apologyRaw from '../components/resources/apology_of_justin_martyr.md?raw';

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

    // Document Title (#)
    if (line.startsWith('# ')) {
      title = line.substring(2).trim();
      continue;
    }

    // Section (##)
    if (line.startsWith('## ')) {
      // If we have an active chapter, save it
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

    // Chapter (#### Chapter)
    const chapMatch = line.match(/^#### Chapter\s+([IVXLCDM\d]+):\s*(.*)/i);
    if (chapMatch) {
      if (currentChapter) {
        if (!currentSection) {
          // If no section yet (like Justin Martyr), create a default section
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

    // Paragraph
    if (currentChapter) {
      currentChapter.content.push(line);
    }
  }

  // Flush remaining
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
  
  // Set default active chapter
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

  // Load reading progress from localStorage
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

// Selected document computed
const activeDoc = computed(() => {
  return documents.value.find(d => d.id === selectedDocId.value) || null;
});

// All chapters in current document for search and stats
const allDocChapters = computed(() => {
  if (!activeDoc.value) return [];
  const list: Chapter[] = [];
  activeDoc.value.sections.forEach(s => {
    list.push(...s.chapters);
  });
  return list;
});

// Selected chapter computed
const activeChapter = computed(() => {
  if (!activeDoc.value) return null;
  for (const s of activeDoc.value.sections) {
    const found = s.chapters.find(c => c.id === selectedChapterId.value);
    if (found) return found;
  }
  return null;
});

// Flat list of chapters for navigation (next/prev)
const flatChapters = computed(() => {
  return allDocChapters.value;
});

const currentChapterIndex = computed(() => {
  return flatChapters.value.findIndex(c => c.id === selectedChapterId.value);
});

// Filtered sections and chapters based on search query
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

// Stats
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
  // Scroll reader pane to top
  const reader = document.getElementById('reader-content');
  if (reader) {
    reader.scrollTop = 0;
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
  <div class="min-h-screen bg-stone-900 text-stone-100 selection:bg-amber-500/30 flex flex-col relative font-sans overflow-hidden">
    <!-- Back Button -->
    <div class="absolute top-6 left-6 z-50">
      <router-link to="/" class="group flex items-center space-x-3 text-white/70 hover:text-white bg-black/20 hover:bg-black/40 backdrop-blur-md px-5 py-2.5 rounded-full transition-all duration-300 border border-white/10 hover:border-white/30">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform">
          <path d="m15 18-6-6 6-6"/>
        </svg>
        <span class="font-medium tracking-wide text-sm">Return Home</span>
      </router-link>
    </div>

    <!-- Background Image -->
    <div class="fixed inset-0 z-0 transition-all duration-1000 ease-in-out" :style="{
      backgroundImage: `url('https://images.unsplash.com/photo-1457369804613-52c61a468e7d?q=80&w=2070&auto=format&fit=crop')`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      filter: 'brightness(0.3) blur(8px) contrast(1.1)'
    }"></div>

    <!-- Gradient Overlay -->
    <div class="fixed inset-0 z-0 bg-gradient-to-t from-stone-900 via-stone-900/60 to-stone-900/40 pointer-events-none"></div>

    <div class="relative z-10 max-w-7xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">
      <!-- Header Area -->
      <header class="text-center mb-8 mt-16 animate-fade-in-down">
        <h1 class="text-4xl md:text-5xl font-serif text-amber-100 mb-2 drop-shadow-lg">Early Church Resources</h1>
        <p class="text-stone-400 text-sm uppercase tracking-[0.2em] font-medium">Timeless writings of faith and theology</p>
      </header>

      <!-- Document Switcher Tabs -->
      <div class="flex justify-center space-x-2 md:space-x-4 mb-8">
        <button 
          @click="selectedDocId = 'ignatius'"
          class="px-6 py-3 rounded-full text-sm font-semibold transition-all border duration-300"
          :class="selectedDocId === 'ignatius' 
            ? 'bg-amber-500/20 text-amber-200 border-amber-500/40 shadow-lg shadow-amber-500/10' 
            : 'bg-black/20 text-stone-400 border-white/10 hover:text-stone-200 hover:bg-black/35'"
        >
          Letters of Ignatius
        </button>
        <button 
          @click="selectedDocId = 'justin'"
          class="px-6 py-3 rounded-full text-sm font-semibold transition-all border duration-300"
          :class="selectedDocId === 'justin' 
            ? 'bg-amber-500/20 text-amber-200 border-amber-500/40 shadow-lg shadow-amber-500/10' 
            : 'bg-black/20 text-stone-400 border-white/10 hover:text-stone-200 hover:bg-black/35'"
        >
          First Apology of Justin Martyr
        </button>
      </div>

      <!-- Main Layout -->
      <main class="flex-grow grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Sidebar Navigation -->
        <aside class="lg:col-span-1 h-fit sticky top-8 flex flex-col gap-6">
          <!-- Search & Progress Card -->
          <div class="bg-black/20 backdrop-blur-md rounded-3xl p-6 border border-white/10 flex flex-col gap-4">
            <!-- Progress Tracker -->
            <div class="border-b border-white/5 pb-4">
              <h3 class="font-serif text-sm text-amber-200/80 tracking-wider uppercase font-bold mb-2">Reading Progress</h3>
              <div class="flex justify-between text-xs text-stone-400 mb-1">
                <span>Completed</span>
                <span>{{ readCountText }} ({{ docProgress }}%)</span>
              </div>
              <div class="w-full bg-white/5 rounded-full h-1.5 overflow-hidden">
                <div class="bg-amber-400 h-1.5 rounded-full transition-all duration-500" :style="{ width: `${docProgress}%` }"></div>
              </div>
            </div>

            <!-- Search -->
            <div>
              <h3 class="font-serif text-sm text-stone-400 tracking-wider uppercase font-bold mb-2">Search Chapters</h3>
              <div class="relative">
                <input 
                  v-model="searchQuery" 
                  type="text" 
                  placeholder="Type title..." 
                  class="w-full pl-3 pr-4 py-2.5 rounded-xl border border-white/10 bg-black/20 text-stone-200 focus:bg-black/40 focus:border-amber-500/50 shadow-md transition-all outline-none text-xs"
                >
              </div>
            </div>
          </div>

          <!-- Table of Contents List -->
          <div class="bg-black/20 backdrop-blur-md rounded-3xl p-6 border border-white/10 max-h-[50vh] overflow-y-auto scrollbar-hide">
            <h3 class="font-serif text-sm mb-4 text-amber-200/80 tracking-wider uppercase font-bold border-b border-white/5 pb-2">Chapters</h3>
            
            <div class="space-y-4">
              <div v-for="section in filteredSections" :key="section.id" class="space-y-2">
                <!-- Section Header (if not default) -->
                <h4 v-if="section.id !== 'default'" class="text-xs font-bold text-stone-400 mt-3 font-serif italic border-l-2 border-amber-500/20 pl-2">
                  {{ section.title.replace('The Epistle of Ignatius to the ', '') }}
                </h4>
                
                <div class="space-y-1">
                  <button
                    v-for="chapter in section.chapters"
                    :key="chapter.id"
                    @click="selectChapter(chapter.id)"
                    class="w-full text-left text-xs transition-colors py-2 px-3 rounded-lg flex items-center justify-between group hover:bg-white/5"
                    :class="selectedChapterId === chapter.id 
                      ? 'bg-amber-500/10 text-amber-300 font-semibold border-l-4 border-amber-400 pl-2' 
                      : readChapters[chapter.id] ? 'text-stone-500 hover:text-stone-300' : 'text-stone-300 hover:text-white'"
                  >
                    <span class="truncate flex-1">
                      <span class="font-serif font-bold text-stone-400 group-hover:text-amber-300 mr-1.5">{{ chapter.numeral }}</span>
                      {{ chapter.title }}
                    </span>
                    <!-- Read Checkbox -->
                    <span 
                      @click.stop="toggleReadStatus(chapter.id)" 
                      class="ml-2 w-4 h-4 rounded border transition-colors flex items-center justify-center cursor-pointer"
                      :class="readChapters[chapter.id] 
                        ? 'bg-emerald-500/20 border-emerald-500/40 text-emerald-400' 
                        : 'border-white/10 hover:border-white/30 text-transparent'"
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
        <div class="lg:col-span-3 flex flex-col gap-6">
          <div class="bg-black/20 backdrop-blur-md rounded-3xl p-6 md:p-8 border border-white/10 flex-grow flex flex-col min-h-[60vh] relative">
            <!-- Reader Area -->
            <div v-if="activeChapter" id="reader-content" class="flex-grow overflow-y-auto space-y-6 pr-2">
              <div class="border-b border-white/5 pb-4 mb-6">
                <div class="text-xs text-amber-400 font-semibold tracking-widest uppercase mb-1">
                  Chapter {{ activeChapter.numeral }}
                </div>
                <h2 class="text-2xl md:text-3xl font-serif text-stone-100 leading-tight">
                  {{ activeChapter.title }}
                </h2>
              </div>

              <!-- Main text paragraphs -->
              <div class="space-y-6 text-stone-300 leading-relaxed font-serif text-sm md:text-base selection:bg-amber-500/20">
                <p v-for="(para, idx) in activeChapter.content" :key="idx" class="indent-6 text-justify">
                  {{ para }}
                </p>
              </div>

              <!-- Checkbox Mark as Read -->
              <div class="pt-8 border-t border-white/5 mt-8 flex justify-center">
                <button
                  @click="toggleReadStatus(activeChapter.id)"
                  class="flex items-center space-x-2.5 px-6 py-3 rounded-full text-xs font-semibold tracking-wider uppercase transition-all duration-300 border"
                  :class="readChapters[activeChapter.id]
                    ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40 hover:bg-emerald-500/30'
                    : 'bg-amber-500/10 text-amber-200 border-amber-500/30 hover:bg-amber-500/20'"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="transition-transform duration-300" :class="{'scale-110': readChapters[activeChapter.id]}">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                  <span>{{ readChapters[activeChapter.id] ? 'Completed' : 'Mark as Read' }}</span>
                </button>
              </div>
            </div>

            <div v-else class="flex-grow flex items-center justify-center text-stone-500 font-serif italic">
              Select a chapter to begin reading
            </div>

            <!-- Navigation Buttons -->
            <div class="mt-8 pt-6 border-t border-white/5 flex justify-between items-center text-xs">
              <button 
                @click="prevChapter" 
                :disabled="currentChapterIndex <= 0"
                class="flex items-center space-x-2 text-stone-400 hover:text-white transition-colors disabled:opacity-30 disabled:hover:text-stone-400"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m15 18-6-6 6-6"/>
                </svg>
                <span>Previous</span>
              </button>

              <button 
                @click="nextChapter" 
                :disabled="currentChapterIndex >= flatChapters.length - 1"
                class="flex items-center space-x-2 text-stone-400 hover:text-white transition-colors disabled:opacity-30 disabled:hover:text-stone-400"
              >
                <span>Next</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m9 18 6-6-6-6"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
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
</style>
