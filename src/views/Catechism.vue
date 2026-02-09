<template>
  <div class="min-h-screen bg-stone-900 text-stone-100 selection:bg-emerald-500/30 flex flex-col relative font-sans overflow-hidden">
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
                backgroundImage: `url('https://images.unsplash.com/photo-1544253328-910f17122115?q=80&w=2070&auto=format&fit=crop')`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                filter: 'brightness(0.4) blur(10px) contrast(1.1)'
            }"></div>

    <!-- Gradient Overlay -->
    <div class="fixed inset-0 z-0 bg-gradient-to-t from-stone-900 via-stone-900/60 to-stone-900/40 pointer-events-none"></div>

    <div class="relative z-10 max-w-7xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">
      
      <!-- Header Area -->
      <header class="text-center mb-12 mt-16 animate-fade-in-down">
          <h1 class="text-4xl md:text-5xl font-serif text-emerald-100 mb-2 drop-shadow-lg">Catechism</h1>
          <p class="text-stone-400 text-sm uppercase tracking-[0.2em] font-medium">Of The Catholic Church</p>
      </header>

      <main class="flex-grow grid grid-cols-1 lg:grid-cols-4 gap-8">
        
        <!-- Search Bar (Mobile: top, Desktop: sidebar top or content top?) -->
        <!-- Placing search in content area for now to match layout sort of, or maybe simplified -->
        
        <!-- Sidebar Navigation -->
        <aside class="lg:col-span-1 h-fit sticky top-8 overflow-y-auto max-h-[calc(100vh-4rem)] animate-fade-in-up delay-100 scrollbar-hide">
          <div class="bg-black/20 backdrop-blur-md rounded-3xl p-6 border border-white/10">
            <h3 class="font-serif text-lg mb-6 text-emerald-200/80 tracking-wider uppercase text-xs font-bold border-b border-white/5 pb-2">Structure</h3>
            
            <div class="space-y-4">
              <div v-for="part in structure.catechism.structure" :key="part.part" class="space-y-2">
                <!-- Part Header -->
                <button 
                  @click="togglePart(part)"
                  class="w-full text-left font-medium text-sm text-stone-400 hover:text-emerald-300 transition-all flex justify-between items-start group p-2 rounded-lg hover:bg-white/5"
                  :class="{'bg-white/10 text-emerald-300': selectedPart?.part === part.part}"
                >
                  <span class="font-serif font-bold mr-2 text-stone-500 group-hover:text-emerald-400 transition-colors">Part {{ part.part }}</span>
                  <span class="flex-1">{{ part.title }}</span>
                </button>
                
                <!-- Sections (only if part selected) -->
                <transition name="fade-slide">
                  <div v-if="selectedPart?.part === part.part" class="pl-4 space-y-2 border-l-2 border-emerald-500/20 ml-3 mt-2">
                    <div v-for="section in part.sections" :key="section.section">
                       <button 
                        @click="toggleSection(section)"
                        class="w-full text-left text-xs text-stone-500 hover:text-emerald-200 transition-colors py-1.5 px-2 rounded block hover:bg-white/5"
                        :class="{'text-emerald-300 font-medium bg-white/5': selectedSection?.section === section.section}"
                      >
                        {{ section.title }}
                      </button>
  
                      <!-- Chapters (only if section selected) -->
                      <transition name="fade-slide">
                        <div v-if="selectedSection?.section === section.section" class="pl-3 mt-1 space-y-1 border-l border-emerald-500/10 ml-1.5">
                            <button
                                v-for="chapter in section.chapters"
                                :key="chapter.chapter"
                                @click="toggleChapter(chapter)"
                                class="w-full text-left text-[11px] text-stone-600 hover:text-emerald-200 transition-colors py-1 px-2 rounded block hover:bg-white/5"
                                :class="{'text-emerald-300 font-medium bg-white/5': selectedChapter?.chapter === chapter.chapter}"
                            >
                                {{ chapter.title }}
                            </button>
                        </div>
                      </transition>
                    </div>
                  </div>
                </transition>
              </div>
            </div>
          </div>
        </aside>

        <!-- Content Area -->
        <div class="lg:col-span-3 flex flex-col gap-6 animate-fade-in-up delay-200">
            
            <!-- Search -->
            <div class="relative group">
                <input 
                    v-model="searchQuery" 
                    type="text" 
                    placeholder="Search paragraph number..." 
                    class="w-full pl-12 pr-4 py-4 rounded-2xl border border-white/10 bg-black/20 backdrop-blur-md text-stone-200 focus:bg-black/40 focus:border-emerald-500/50 focus:ring-1 focus:ring-emerald-500/50 shadow-lg transition-all outline-none"
                >
            </div>

            <!-- Active Path Display -->
            <div v-if="!searchQuery && (selectedPart || selectedSection || selectedChapter)" class="bg-emerald-950/30 backdrop-blur-sm border border-emerald-500/20 rounded-2xl p-6 text-sm flex flex-wrap gap-2 items-center shadow-sm">
                <span v-if="selectedPart" class="font-bold text-emerald-400 text-lg">Part {{ selectedPart.part }}: {{ selectedPart.title }}</span>
                <span v-if="selectedSection" class="text-stone-600 mx-2">/</span>
                <span v-if="selectedSection" class="font-medium text-emerald-200/80">{{ selectedSection.title }}</span>
                <span v-if="selectedChapter" class="text-stone-600 mx-2">/</span>
                <span v-if="selectedChapter" class="text-stone-400">{{ selectedChapter.title }}</span>
            </div>

          <!-- Paragraphs List -->
          <div class="space-y-6">
            <template v-if="filteredParagraphs.length > 0">
                <div v-for="(paragraph, index) in filteredParagraphs" :key="paragraph.id" 
                     class="bg-black/20 backdrop-blur-md rounded-3xl p-8 shadow-xl border border-white/5 hover:bg-black/30 hover:border-emerald-500/20 transition-all duration-500 group/card"
                     :class="{ 'opacity-50 grayscale': readParagraphs.has(paragraph.id) }"
                     :style="{ animationDelay: `${index * 50}ms` }">
                  <div class="flex items-start gap-6">
                    <span class="flex-shrink-0 w-10 h-10 rounded-full bg-emerald-900/50 text-emerald-400 flex items-center justify-center font-serif font-bold text-sm ring-2 ring-emerald-500/20 relative cursor-pointer transition-all duration-300 hover:bg-emerald-500 hover:text-white hover:shadow-[0_0_15px_rgba(16,185,129,0.5)]"
                        :class="{ '!bg-green-500 !text-white !ring-green-400 shadow-[0_0_10px_rgba(34,197,94,0.4)]': readParagraphs.has(paragraph.id) }"
                        @click="toggleRead(paragraph.id)">
                      <span class="transition-opacity duration-200" :class="{ 'opacity-0': readParagraphs.has(paragraph.id) }">{{ paragraph.id }}</span>
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" 
                        class="absolute opacity-0 transition-opacity duration-200 text-white"
                        :class="{ '!opacity-100': readParagraphs.has(paragraph.id) }">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    </span>
                    <div class="prose prose-invert prose-lg max-w-none">
                      <div class="text-xs text-stone-500 mb-2 font-bold tracking-wider uppercase opacity-60">{{ getBreadcrumb(paragraph.id) }}</div>
                      <p class="whitespace-pre-line leading-loose text-stone-300 group-hover/card:text-stone-100 transition-colors">{{ paragraph.text }}</p>
                    </div>
                  </div>
                </div>
            </template>

            <!-- Empty State / Instructions -->
             <div v-else class="text-center py-32 bg-black/10 backdrop-blur-sm rounded-[3rem] border-2 border-dashed border-white/5">
                <div class="w-20 h-20 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-6 shadow-inner">
                    <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-stone-600">
                        <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/>
                        <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>
                    </svg>
                </div>
                <h3 class="text-xl font-serif text-stone-300 mb-2">
                    {{ searchQuery ? `No paragraph found for "${searchQuery}"` : 'Begin Your Study' }}
                </h3>
                <p class="text-stone-500 max-w-sm mx-auto leading-relaxed">
                    {{ searchQuery ? 'Try searching for a valid paragraph number.' : 'Select a section from the structure on the left or search for a specific paragraph number above.' }}
                </p>
             </div>
          </div>
        </div>
    </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import catechismData from '../components/catechism/catechism.json';
import structureData from '../components/catechism/catechism-sections.json';

// Types
interface Paragraph {
  id: number;
  text: string;
}

interface Chapter {
  chapter: number;
  title: string;
  paragraphs: string;
}

interface Section {
  section: number;
  title: string;
  paragraph_range: string;
  chapters: Chapter[];
}

interface Part {
  part: number;
  title: string;
  paragraph_range: string;
  sections: Section[];
}

interface Structure {
    catechism: {
        title: string;
        structure: Part[];
    }
}

// State
// Explicitly cast imported JSON to types
const structure = ref<Structure>(structureData as unknown as Structure);
const paragraphs = ref<Paragraph[]>(catechismData as unknown as Paragraph[]);
const searchQuery = ref('');

const selectedPart = ref<Part | null>(null);
const selectedSection = ref<Section | null>(null);
const selectedChapter = ref<Chapter | null>(null);

// Progress Tracking
const readParagraphs = ref<Set<number>>(new Set());

onMounted(() => {
    const saved = localStorage.getItem('catechism-read-progress');
    if (saved) {
        try {
            const parsed = JSON.parse(saved);
            if (Array.isArray(parsed)) {
                readParagraphs.value = new Set(parsed);
            }
        } catch (e) {
            console.error('Failed to load catechism progress', e);
        }
    }
});

const toggleRead = (id: number) => {
    const newSet = new Set(readParagraphs.value);
    if (newSet.has(id)) {
        newSet.delete(id);
    } else {
        newSet.add(id);
    }
    readParagraphs.value = newSet;
    localStorage.setItem('catechism-read-progress', JSON.stringify(Array.from(newSet)));
};

// Methods
const togglePart = (part: Part) => {
  if (selectedPart.value?.part === part.part) {
    selectedPart.value = null;
    selectedSection.value = null;
    selectedChapter.value = null;
  } else {
    selectedPart.value = part;
    selectedSection.value = null; // Reset child selections when changing parent
    selectedChapter.value = null;
  }
};

const toggleSection = (section: Section) => {
    if (selectedSection.value?.section === section.section) {
        selectedSection.value = null;
        selectedChapter.value = null;
    } else {
        selectedSection.value = section;
        selectedChapter.value = null;
    }
}

const toggleChapter = (chapter: Chapter) => {
    if (selectedChapter.value?.chapter === chapter.chapter) {
        selectedChapter.value = null;
    } else {
        selectedChapter.value = chapter;
    }
}

// Helper to parse range string "1-10" to object {start, end}
const parseRange = (rangeStr: string): { start: number, end: number } => {
  if (!rangeStr) return { start: 0, end: 0 };
  
  // Handle single numbers or malformed strings if necessary, though data seems consistent "X-Y"
  const parts = rangeStr.split('-');
  if (parts.length === 2) {
    const startStr = parts[0];
    const endStr = parts[1];
    if (typeof startStr === 'string' && typeof endStr === 'string') {
        const start = parseInt(startStr.trim());
        const end = parseInt(endStr.trim());
        return { start, end };
    }
  }
  return { start: 0, end: 0 };
};

const getBreadcrumb = (id: number): string => {
  const s = structure.value.catechism.structure;
  for (const part of s) {
    const pRange = parseRange(part.paragraph_range);
    if (id >= pRange.start && id <= pRange.end) {
      const parts = [part.title];
      
      if (part.sections) {
        for (const section of part.sections) {
          const sRange = parseRange(section.paragraph_range);
          if (id >= sRange.start && id <= sRange.end) {
            parts.push(section.title);
            
            if (section.chapters) {
              for (const chapter of section.chapters) {
                const cRange = parseRange(chapter.paragraphs);
                if (id >= cRange.start && id <= cRange.end) {
                  parts.push(chapter.title);
                  break;
                }
              }
            }
            break;
          }
        }
      }
      return parts.join(' / ');
    }
  }
  return '';
};

// Computed
const filteredParagraphs = computed(() => {
  // 1. Search filter (highest priority) - if there is a search query, ignore navigation
  if (searchQuery.value.trim()) {
    const queryNum = parseInt(searchQuery.value);
    if (!isNaN(queryNum)) {
      return paragraphs.value.filter(p => p.id === queryNum);
    }
    return [];
  }

  // 2. Navigation filter
  let range = { start: 0, end: 0 };

  if (selectedChapter.value) {
      range = parseRange(selectedChapter.value.paragraphs);
  } else if (selectedSection.value) {
      range = parseRange(selectedSection.value.paragraph_range);
  } else if (selectedPart.value) {
      range = parseRange(selectedPart.value.paragraph_range);
  } else {
      return []; 
  }

  if (range.start === 0 && range.end === 0) return [];

  return paragraphs.value.filter(p => p.id >= range.start && p.id <= range.end);
});

// Watch inputs to auto-expand navigation if possible? 
// For now, let's just make sure state is consistent
watch(searchQuery, (newVal) => {
    if (newVal) {
        // defined behavior: search overrides navigation, so we don't necessarily need to clear selection,
        // but user might find it confusing if selection remains highlighted while content shows search result.
        // Let's keep selection state for now, as clearing it might be annoying if they just typed a number to check something quick.
    }
});

</script>

<style scoped>
/* Transitions moved to global style.css */
</style>
