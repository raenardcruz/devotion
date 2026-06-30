<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-24 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <div class="relative z-10 max-w-7xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">
      
      <!-- Header Area -->
      <header class="text-center mb-8 animate-fade-in-down">
          <h1 class="text-3xl md:text-5xl font-serif text-parchment-primary-dark mb-1">Catechism</h1>
          <p class="text-parchment-neutral/50 text-xs uppercase tracking-[0.25em] font-bold">Of The Catholic Church</p>
      </header>

      <main class="flex-grow grid grid-cols-1 lg:grid-cols-4 gap-8">
        
        <!-- Sidebar Navigation (Structure) -->
        <aside class="lg:col-span-1 h-fit sticky top-24 overflow-y-auto max-h-[calc(100vh-8rem)] animate-fade-in-up delay-100 scrollbar-hide">
          <div class="bg-parchment-neutral-light border border-parchment-border rounded-3xl p-6 shadow-sm">
            <h3 class="font-serif text-sm mb-6 text-parchment-primary-dark/80 tracking-widest uppercase font-bold border-b border-parchment-border/40 pb-2.5">Structure</h3>
            
            <div class="space-y-4">
              <div v-for="part in structure.catechism.structure" :key="part.part" class="space-y-2">
                <!-- Part Header -->
                <button 
                  @click="togglePart(part)"
                  class="w-full text-left font-medium text-sm text-parchment-neutral/70 hover:text-parchment-neutral transition-all flex justify-between items-start group p-2 rounded-lg hover:bg-parchment-bg border-none shadow-none"
                  :class="{'bg-parchment-primary/10 text-parchment-primary-dark border border-parchment-primary/20': selectedPart?.part === part.part}"
                >
                  <span class="font-serif font-bold mr-2 text-parchment-primary group-hover:text-parchment-primary-dark transition-colors">Part {{ part.part }}</span>
                  <span class="flex-1">{{ part.title }}</span>
                </button>
                
                <!-- Sections (only if part selected) -->
                <transition name="fade-slide">
                  <div v-if="selectedPart?.part === part.part" class="pl-4 space-y-2 border-l-2 border-parchment-primary/30 ml-3 mt-2">
                    <div v-for="section in part.sections" :key="section.section">
                       <button 
                        @click="toggleSection(section)"
                        class="w-full text-left text-xs text-parchment-neutral/70 hover:text-parchment-neutral transition-colors py-1.5 px-2 rounded block hover:bg-parchment-bg border-none shadow-none"
                        :class="{'text-parchment-primary-dark font-medium bg-parchment-primary/5': selectedSection?.section === section.section}"
                      >
                        {{ section.title }}
                      </button>
  
                      <!-- Chapters (only if section selected) -->
                      <transition name="fade-slide">
                        <div v-if="selectedSection?.section === section.section" class="pl-3 mt-1 space-y-1 border-l border-parchment-primary/10 ml-1.5">
                            <button
                                v-for="chapter in section.chapters"
                                :key="chapter.chapter"
                                @click="toggleChapter(chapter)"
                                class="w-full text-left text-[11px] text-parchment-neutral/50 hover:text-parchment-neutral transition-colors py-1 px-2 rounded block hover:bg-parchment-bg border-none shadow-none"
                                :class="{'text-parchment-primary font-bold bg-parchment-primary/5': selectedChapter?.chapter === chapter.chapter}"
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
        <div class="lg:col-span-3 flex flex-col gap-6 animate-fade-in-up delay-150">
            
            <!-- Search Input -->
            <div class="relative group">
                <input 
                    v-model="searchQuery" 
                    type="text" 
                    placeholder="Search paragraph number (e.g. 5-10, 20)..." 
                    class="w-full pl-12 pr-4 py-4 rounded-2xl border border-parchment-border bg-parchment-neutral-light/50 text-parchment-neutral placeholder-parchment-neutral/40 focus:bg-parchment-neutral-light focus:border-parchment-primary focus:ring-1 focus:ring-parchment-primary shadow-sm transition-all outline-none"
                >
                <div class="absolute left-4 top-1/2 -translate-y-1/2 text-parchment-neutral/30">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="11" cy="11" r="8"></circle>
                    <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                  </svg>
                </div>
            </div>

            <!-- Active Path Display (Breadcrumbs) -->
            <div v-if="!searchQuery && (selectedPart || selectedSection || selectedChapter)" class="bg-parchment-neutral-light/40 backdrop-blur-sm border border-parchment-border rounded-2xl p-5 text-sm flex flex-wrap gap-2 items-center shadow-sm">
                <span v-if="selectedPart" class="font-serif font-bold text-parchment-primary-dark text-base">Part {{ selectedPart.part }}: {{ selectedPart.title }}</span>
                <span v-if="selectedSection" class="text-parchment-neutral/30 mx-1">/</span>
                <span v-if="selectedSection" class="font-medium text-parchment-neutral/80">{{ selectedSection.title }}</span>
                <span v-if="selectedChapter" class="text-parchment-neutral/30 mx-1">/</span>
                <span v-if="selectedChapter" class="text-parchment-neutral/60">{{ selectedChapter.title }}</span>
            </div>

          <!-- Paragraphs List -->
          <div class="space-y-6">
            <template v-if="filteredParagraphs.length > 0">
                <div v-for="(paragraph, index) in filteredParagraphs" :key="paragraph.id" 
                     class="bg-parchment-neutral-light border border-parchment-border rounded-3xl p-6 md:p-8 shadow-sm transition-all duration-300 group/card"
                     :class="{ 'opacity-60': readParagraphs.has(paragraph.id) }"
                     :style="{ animationDelay: `${index * 50}ms` }">
                  <div class="flex items-start gap-5">
                    <!-- Read Checkbox Circle -->
                    <span class="flex-shrink-0 w-10 h-10 rounded-full border-2 border-parchment-primary text-parchment-primary-dark flex items-center justify-center font-serif font-bold text-sm bg-parchment-bg relative cursor-pointer transition-all duration-300 hover:bg-parchment-primary hover:text-white"
                        :class="{ '!bg-parchment-primary !text-white': readParagraphs.has(paragraph.id) }"
                        @click="toggleRead(paragraph.id)">
                      <span class="transition-opacity duration-200" :class="{ 'opacity-0': readParagraphs.has(paragraph.id) }">{{ paragraph.id }}</span>
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" 
                        class="absolute opacity-0 transition-opacity duration-200 text-white"
                        :class="{ '!opacity-100': readParagraphs.has(paragraph.id) }">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    </span>
                    
                    <div class="prose max-w-none flex-1">
                      <div class="text-[10px] text-parchment-neutral/40 mb-1.5 font-bold tracking-wider uppercase select-none">{{ getBreadcrumb(paragraph.id) }}</div>
                      <!-- Paragraph text rendering with reference helper link -->
                      <CatechismText :text="paragraph.text" :paragraphs="paragraphs" @show-reference="openModal" />
                    </div>
                  </div>
                </div>
            </template>

            <!-- Empty State / Instructions -->
             <div v-else class="text-center py-24 bg-parchment-neutral-light/30 rounded-3xl border border-dashed border-parchment-border">
                <div class="w-16 h-16 bg-parchment-neutral-light rounded-full flex items-center justify-center mx-auto mb-4 shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/40">
                        <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/>
                        <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>
                    </svg>
                </div>
                <h3 class="text-lg font-serif text-parchment-neutral font-medium mb-1.5">
                    {{ searchQuery ? `No paragraph found for "${searchQuery}"` : 'Begin Your Study' }}
                </h3>
                <p class="text-parchment-neutral/50 max-w-sm mx-auto leading-relaxed text-xs md:text-sm">
                    {{ searchQuery ? 'Try searching for a valid paragraph number.' : 'Select a section from the structure on the left or search for a specific paragraph number above.' }}
                </p>
             </div>
          </div>
        </div>
      </main>
    </div>

    <!-- Reference Modal -->
    <transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100" leave-to-class="opacity-0">
        <div v-if="showModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4" @click="showModal = false">
            <!-- Backdrop -->
            <div class="absolute inset-0 bg-parchment-neutral/30 backdrop-blur-sm"></div>
            
            <!-- Modal Content -->
            <div class="relative w-full max-w-2xl bg-parchment-bg border border-parchment-border rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[80vh]" @click.stop>
                <!-- Header -->
                <div class="p-5 border-b border-parchment-border flex justify-between items-center bg-parchment-neutral-light/50">
                    <h3 class="text-lg font-serif text-parchment-primary-dark font-medium">Referenced Paragraphs</h3>
                    <button @click="showModal = false" class="p-2 hover:bg-parchment-neutral-light rounded-full transition-colors text-parchment-neutral/50 hover:text-parchment-neutral border-none shadow-none">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="18" x2="6" y1="6" y2="18"></line>
                            <line x1="6" x2="18" y1="6" y2="18"></line>
                        </svg>
                    </button>
                </div>

                <!-- Body -->
                <div class="p-6 overflow-y-auto space-y-4">
                    <div v-for="p in modalContent" :key="p.id" class="p-5 rounded-2xl bg-parchment-neutral-light/55 border border-parchment-border">
                        <div class="flex items-start gap-4">
                            <span class="text-parchment-secondary font-bold text-lg font-serif opacity-85 select-none">§{{ p.id }}</span>
                            <p class="text-parchment-neutral font-serif leading-relaxed text-sm md:text-base">{{ p.text }}</p>
                        </div>
                    </div>
                </div>
                
                <!-- Footer -->
                <div class="p-4 bg-parchment-neutral-light/50 text-center border-t border-parchment-border">
                     <span class="text-[10px] text-parchment-neutral/40 uppercase tracking-widest font-semibold">Catechism of the Catholic Church</span>
                </div>
            </div>
        </div>
    </transition>

    <!-- Global Footer -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import catechismData from '../components/catechism/catechism.json';
import structureData from '../components/catechism/catechism-sections.json';
import CatechismText from '../components/catechism/CatechismText.vue';

import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';

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

const route = useRoute();

// State
const structure = ref<Structure>(structureData as unknown as Structure);
const paragraphs = ref<Paragraph[]>(catechismData as unknown as Paragraph[]);
const searchQuery = ref('');

const selectedPart = ref<Part | null>(null);
const selectedSection = ref<Section | null>(null);
const selectedChapter = ref<Chapter | null>(null);

// Progress Tracking
const readParagraphs = ref<Set<number>>(new Set());

// Modal State
const showModal = ref(false);
const modalContent = ref<Paragraph[]>([]);

const openModal = (paragraphs: Paragraph[]) => {
    modalContent.value = paragraphs;
    showModal.value = true;
};

onMounted(() => {
    // Sync initial search from router query if present (from BottomNav search query)
    if (route.query.q) {
        searchQuery.value = String(route.query.q);
    }

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
    selectedSection.value = null; 
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
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.trim();
    const desiredIds = new Set<number>();
    const parts = query.split(',');
    
    for (const part of parts) {
        const trimmed = part.trim();
        if (!trimmed) continue;
        
        if (trimmed.includes('-')) {
            const rangeParts = trimmed.split('-').map(s => s.trim());
            if (rangeParts.length === 2) {
                const start = parseInt(rangeParts[0]!);
                const end = parseInt(rangeParts[1]!);
                
                if (!isNaN(start) && !isNaN(end) && start <= end) {
                    for (let i = start; i <= end; i++) {
                        desiredIds.add(i);
                    }
                }
            }
        } else {
            const num = parseInt(trimmed);
            if (!isNaN(num)) {
                desiredIds.add(num);
            }
        }
    }
    
    if (desiredIds.size > 0) {
        return paragraphs.value.filter(p => desiredIds.has(p.id));
    }
    
    return [];
  }

  let range = { start: 0, end: 0 };

  if (selectedChapter.value) {
      range = parseRange(selectedChapter.value.paragraphs || '');
  } else if (selectedSection.value) {
      range = parseRange(selectedSection.value.paragraph_range || '');
  } else if (selectedPart.value) {
      range = parseRange(selectedPart.value.paragraph_range || '');
  } else {
      return []; 
  }

  if (range.start === 0 && range.end === 0) return [];

  return paragraphs.value.filter(p => p.id >= range.start && p.id <= range.end);
});

// Watch query from query params
watch(() => route.query.q, (newVal) => {
    if (newVal !== undefined) {
        searchQuery.value = String(newVal);
    }
});
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
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
