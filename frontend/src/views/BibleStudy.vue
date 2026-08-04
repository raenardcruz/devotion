<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useDevotionApi } from '../composables/useDevotionApi';
import { renderMarkdown } from '../utils/markdown';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import AppButton from '../components/common/AppButton.vue';
import AppTabs from '../components/common/AppTabs.vue';
import FactCheckCitations, { type MagisteriumCitation } from '../components/common/FactCheckCitations.vue';

interface SavedVerse {
  id: string;
  citation: string;
  content: string;
  context: string;
  citations?: MagisteriumCitation[];
  copyright?: string;
  timestamp: number;
}

interface BibleResponse {
  content: string;
  copyright: string;
}

interface ContextResponse {
  citation: string;
  context: string;
  citations?: MagisteriumCitation[];
}

const { fetchWithAuth } = useDevotionApi();
const route = useRoute();

// Form input state
const activeTab = ref('selector');
const selectedBook = ref('John');
const chapterInput = ref<number | ''>(3);
const verseStart = ref<number | ''>(16);
const verseEnd = ref<number | ''>('');
const manualQuery = ref('');
const showBookDropdown = ref(false);
const bookSearchQuery = ref('');

// UI state
const loading = ref(false);
const error = ref<string | null>(null);
const savedVerses = ref<SavedVerse[]>([]);
const expandedIds = ref<Set<string>>(new Set());
const isMobileSelectorOpen = ref(false);

const cleanContext = (contextStr?: string) => {
  if (!contextStr) return '';
  return contextStr.replace(/<details[\s\S]*?<\/details>/gi, '').trim();
};

// Book data for Catholic Bible (73 books)
const categories = [
  {
    name: 'Gospels & Acts',
    books: ['Matthew', 'Mark', 'Luke', 'John', 'Acts']
  },
  {
    name: 'New Testament Letters',
    books: [
      'Romans', '1 Corinthians', '2 Corinthians', 'Galatians', 'Ephesians',
      'Philippians', 'Colossians', '1 Thessalonians', '2 Thessalonians',
      '1 Timothy', '2 Timothy', 'Titus', 'Philemon', 'Hebrews', 'James',
      '1 Peter', '2 Peter', '1 John', '2 John', '3 John', 'Jude', 'Revelation'
    ]
  },
  {
    name: 'Old Testament Pentateuch & Historical',
    books: [
      'Genesis', 'Exodus', 'Leviticus', 'Numbers', 'Deuteronomy',
      'Joshua', 'Judges', 'Ruth', '1 Samuel', '2 Samuel', '1 Kings', '2 Kings',
      '1 Chronicles', '2 Chronicles', 'Ezra', 'Nehemiah', 'Tobit', 'Judith',
      'Esther', '1 Maccabees', '2 Maccabees'
    ]
  },
  {
    name: 'Old Testament Wisdom & Prophets',
    books: [
      'Job', 'Psalms', 'Proverbs', 'Ecclesiastes', 'Song of Songs', 'Wisdom',
      'Sirach', 'Isaiah', 'Jeremiah', 'Lamentations', 'Baruch', 'Ezekiel',
      'Daniel', 'Hosea', 'Joel', 'Amos', 'Obadiah', 'Jonah', 'Micah',
      'Nahum', 'Habakkuk', 'Zephaniah', 'Haggai', 'Zechariah', 'Malachi'
    ]
  }
];

// Flat list of books for search filter
const allBooks = computed(() => {
  return categories.flatMap(c => c.books);
});

const filteredBooks = computed(() => {
  const query = bookSearchQuery.value.trim().toLowerCase();
  if (!query) return allBooks.value;
  return allBooks.value.filter(b => b.toLowerCase().includes(query));
});

const tabs = [
  { id: 'selector', label: 'Selector' },
  { id: 'search', label: 'Direct Search' }
];

// Initialize from localStorage
onMounted(async () => {
  try {
    const raw = localStorage.getItem('saved_verses');
    if (raw) {
      savedVerses.value = JSON.parse(raw);
    }
  } catch (err) {
    console.error('Error parsing saved verses from localStorage:', err);
  }

  // Check query parameter 'query' or 'passage'
  const queryParam = route.query.query || route.query.passage;
  if (queryParam) {
    activeTab.value = 'search';
    manualQuery.value = String(queryParam);
    await getContext();
  } else {
    // Auto-expand the first item if there is any history
    if (savedVerses.value.length > 0) {
      const firstItem = savedVerses.value[0];
      if (firstItem) {
        expandedIds.value.add(firstItem.id);
      }
    }
  }
});

watch(
  () => route.query.query || route.query.passage,
  async (newQuery) => {
    if (newQuery) {
      activeTab.value = 'search';
      manualQuery.value = String(newQuery);
      await getContext();
    }
  }
);

const toggleExpand = (id: string) => {
  if (expandedIds.value.has(id)) {
    expandedIds.value.delete(id);
  } else {
    expandedIds.value.add(id);
  }
};

const deleteVerse = (id: string, event: Event) => {
  event.stopPropagation();
  savedVerses.value = savedVerses.value.filter(item => item.id !== id);
  localStorage.setItem('saved_verses', JSON.stringify(savedVerses.value));
  expandedIds.value.delete(id);
};

const clearAllHistory = () => {
  if (confirm('Are you sure you want to clear all your Bible study history?')) {
    savedVerses.value = [];
    localStorage.removeItem('saved_verses');
    expandedIds.value.clear();
  }
};

const selectBook = (book: string) => {
  selectedBook.value = book;
  showBookDropdown.value = false;
  bookSearchQuery.value = '';
};

// Main fetch logic
const getContext = async () => {
  error.value = null;
  let passage = '';

  if (activeTab.value === 'selector') {
    if (!selectedBook.value) {
      error.value = 'Please select a book.';
      return;
    }
    if (!chapterInput.value) {
      error.value = 'Please enter a chapter number.';
      return;
    }
    
    passage = `${selectedBook.value} ${chapterInput.value}`;
    if (verseStart.value) {
      passage += `:${verseStart.value}`;
      if (verseEnd.value) {
        passage += `-${verseEnd.value}`;
      }
    }
  } else {
    if (!manualQuery.value.trim()) {
      error.value = 'Please enter a Bible passage citation (e.g. Genesis 1:1).';
      return;
    }
    passage = manualQuery.value.trim();
  }

  loading.value = true;

  try {
    const encodedPassage = encodeURIComponent(passage);
    
    // Fetch scripture text and theological context in parallel
    const [bibleRes, contextRes] = await Promise.all([
      fetchWithAuth<BibleResponse>(`/bible?passage=${encodedPassage}`),
      fetchWithAuth<ContextResponse>(`/context?passage=${encodedPassage}`)
    ]);

    if (!bibleRes || !bibleRes.content) {
      throw new Error(`Could not find text for passage: ${passage}`);
    }

    const newItem: SavedVerse = {
      id: Date.now().toString(),
      citation: contextRes.citation || passage,
      content: bibleRes.content,
      context: contextRes.context || 'Context commentary unavailable.',
      citations: contextRes.citations,
      copyright: bibleRes.copyright,
      timestamp: Date.now()
    };

    // Prepend to history list
    savedVerses.value.unshift(newItem);
    localStorage.setItem('saved_verses', JSON.stringify(savedVerses.value));

    // Expand the newly added card
    expandedIds.value.add(newItem.id);

    // Reset direct search if using that
    if (activeTab.value === 'search') {
      manualQuery.value = '';
    }

  } catch (err: any) {
    console.error('Error fetching context:', err);
    error.value = err.message || 'An error occurred while fetching context. Please check the spelling or format.';
  } finally {
    loading.value = false;
  }
};

const formatDate = (timestamp: number) => {
  return new Date(timestamp).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const getSnippet = (text: string) => {
  const clean = text.replace(/<[^>]*>/g, '');
  if (clean.length <= 110) return clean;
  return clean.substring(0, 110) + '...';
};

</script>

<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-28 sm:pb-20 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Layout Container -->
    <div class="relative z-10 max-w-7xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">
      <!-- Title Header -->
      <header class="text-center mb-8 animate-fade-in-down">
        <span class="text-parchment-primary font-bold tracking-[0.25em] text-xs uppercase block mb-1">
          Scripture Study
        </span>
        <h1 class="text-3xl md:text-5xl font-serif text-parchment-primary-dark font-medium my-2">
          Context Explorer
        </h1>
        <p class="text-parchment-neutral/70text-sm max-w-md mx-auto leading-relaxed">
          Unveil the historical context, theological layers, and original scriptures of any biblical passage.
        </p>
      </header>

      <!-- Mobile Navigation Trigger Button -->
      <div class="lg:hidden mb-4 flex items-center justify-between bg-parchment-neutral-light border border-parchment-border rounded-2xl p-3 shadow-sm">
        <div class="flex items-center gap-2 overflow-hidden pr-2">
          <span class="text-parchment-primary">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/>
              <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>
            </svg>
          </span>
          <span class="font-serif text-xs font-semibold text-parchment-primary-dark truncate">
            {{ selectedBook }} {{ chapterInput ? chapterInput : '' }}{{ verseStart ? ':' + verseStart : '' }}
          </span>
        </div>
        <button 
          @click="isMobileSelectorOpen = true" 
          class="px-3 py-1.5 bg-parchment-primary text-white text-xs font-medium rounded-xl hover:bg-parchment-primary-dark transition-all shadow-xs shrink-0 flex items-center gap-1.5"
        >
          <span>Select Passage</span>
        </button>
      </div>

      <!-- Mobile Navigation Drawer -->
      <transition 
        enter-active-class="transition duration-300 ease-out" 
        enter-from-class="opacity-0" 
        enter-to-class="opacity-100" 
        leave-active-class="transition duration-200 ease-in" 
        leave-from-class="opacity-100" 
        leave-to-class="opacity-0"
      >
        <div v-if="isMobileSelectorOpen" class="fixed inset-0 z-50 flex justify-start lg:hidden" @click="isMobileSelectorOpen = false">
          <!-- Backdrop -->
          <div class="absolute inset-0 bg-parchment-neutral/40 backdrop-blur-xs"></div>
          
          <!-- Panel Content -->
          <div class="relative w-4/5 max-w-xs bg-parchment-bg h-full shadow-2xl p-4 flex flex-col z-10 overflow-y-auto border-r border-parchment-border space-y-4" @click.stop>
            <div class="flex items-center justify-between pb-3 border-b border-parchment-border">
              <h3 class="font-serif text-base text-parchment-primary-dark font-bold">Passage Selector</h3>
              <button @click="isMobileSelectorOpen = false" class="p-1.5 text-parchment-neutral/50 hover:text-parchment-neutral rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <!-- Mobile Selector Form Card inside Drawer -->
            <ParchmentCard class="shadow-sm">
              <AppTabs :tabs="tabs" v-model="activeTab" class="mb-4" />
              <div v-if="activeTab === 'selector'" class="space-y-3">
                <!-- Book Dropdown -->
                <div>
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1">Bible Book</label>
                  <button 
                    type="button"
                    @click="showBookDropdown = !showBookDropdown"
                    class="w-full text-left bg-parchment-bg border border-parchment-border/70 hover:border-parchment-primary/50 px-3 py-2 rounded-xl flex items-center justify-between text-xs transition-all"
                  >
                    <span class="font-medium text-parchment-neutral truncate">{{ selectedBook || 'Select a Book' }}</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/50 shrink-0" :class="{ 'rotate-180': showBookDropdown }">
                      <path d="m6 9 6 6 6-6"></path>
                    </svg>
                  </button>
                  <div v-if="showBookDropdown" class="mt-2 bg-parchment-bg border border-parchment-border rounded-xl shadow-xl max-h-56 overflow-y-auto p-2 space-y-2">
                    <input type="text" v-model="bookSearchQuery" placeholder="Search books..." class="w-full bg-parchment-neutral-light/50 border border-parchment-border/40 text-xs px-2.5 py-1.5 rounded-lg text-parchment-neutral outline-none mb-1"/>
                    <div class="grid grid-cols-1 gap-1">
                      <button v-for="book in filteredBooks" :key="book" type="button" @click="selectBook(book)" class="text-left text-xs px-2 py-1 rounded hover:bg-parchment-neutral-light truncate" :class="{ 'bg-parchment-primary/10 text-parchment-primary-dark font-bold': selectedBook === book }">{{ book }}</button>
                    </div>
                  </div>
                </div>

                <!-- Chapter Input -->
                <div>
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1">Chapter</label>
                  <input type="number" v-model="chapterInput" placeholder="e.g. 3" min="1" class="w-full bg-parchment-bg border border-parchment-border/70 px-3 py-2 rounded-xl text-xs text-parchment-neutral outline-none"/>
                </div>

                <!-- Verse Selector -->
                <div>
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1">Verse Range</label>
                  <div class="flex items-center space-x-1.5">
                    <input type="number" v-model="verseStart" placeholder="Start" min="1" class="w-full bg-parchment-bg border border-parchment-border/70 px-2 py-2 rounded-xl text-xs text-center text-parchment-neutral outline-none"/>
                    <span class="text-parchment-neutral/40 text-xs font-bold">—</span>
                    <input type="number" v-model="verseEnd" placeholder="End" min="1" class="w-full bg-parchment-bg border border-parchment-border/70 px-2 py-2 rounded-xl text-xs text-center text-parchment-neutral outline-none"/>
                  </div>
                </div>
              </div>

              <!-- Direct Search -->
              <div v-else class="space-y-2">
                <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1">Citation</label>
                <input type="text" v-model="manualQuery" @keyup.enter="getContext(); isMobileSelectorOpen = false" placeholder="e.g. John 3:16" class="w-full bg-parchment-bg border border-parchment-border/70 px-3 py-2 rounded-xl text-xs text-parchment-neutral outline-none"/>
              </div>

              <div class="mt-4 pt-3 border-t border-parchment-border/30">
                <AppButton variant="primary" :disabled="loading" @click="getContext(); isMobileSelectorOpen = false" class="w-full text-xs py-2">Explore Context</AppButton>
              </div>
            </ParchmentCard>
          </div>
        </div>
      </transition>

      <!-- Main Layout Grid -->
      <main class="flex-grow grid grid-cols-1 lg:grid-cols-12 gap-8">
        
        <!-- Sidebar Navigation (Desktop) -->
        <aside class="hidden lg:block lg:col-span-4 xl:col-span-4 h-fit lg:sticky lg:top-24 flex flex-col gap-6 animate-fade-in-up delay-100">
          <ParchmentCard class="shadow-sm">
            <h3 class="font-serif text-sm mb-4 text-parchment-primary-dark/80 tracking-widest uppercase font-bold border-b border-parchment-border/40 pb-2.5">Passage Selector</h3>
            <AppTabs :tabs="tabs" v-model="activeTab" class="mb-4" />

            <div v-if="activeTab === 'selector'" class="space-y-4">
              <!-- Book Dropdown -->
              <div class="relative">
                <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1.5">Bible Book</label>
                <button type="button" @click="showBookDropdown = !showBookDropdown" class="w-full text-left bg-parchment-bg border border-parchment-border/70 hover:border-parchment-primary/50 px-3.5 py-2.5 rounded-xl flex items-center justify-between text-xs transition-all">
                  <span class="font-medium text-parchment-neutral truncate">{{ selectedBook || 'Select a Book' }}</span>
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/50 shrink-0" :class="{ 'rotate-180': showBookDropdown }"><path d="m6 9 6 6 6-6"></path></svg>
                </button>
                <div v-if="showBookDropdown" class="absolute z-50 left-0 right-0 mt-1 bg-parchment-bg border border-parchment-border rounded-xl shadow-xl max-h-60 overflow-y-auto p-2 space-y-2">
                  <input type="text" v-model="bookSearchQuery" placeholder="Search books..." class="w-full bg-parchment-neutral-light/50 border border-parchment-border/40 text-xs px-2.5 py-1.5 rounded-lg text-parchment-neutral outline-none mb-1"/>
                  <div class="grid grid-cols-1 gap-1">
                    <button v-for="book in filteredBooks" :key="book" type="button" @click="selectBook(book)" class="text-left text-xs px-2.5 py-1.5 rounded-lg hover:bg-parchment-neutral-light truncate" :class="{ 'bg-parchment-primary/10 text-parchment-primary-dark font-bold': selectedBook === book }">{{ book }}</button>
                  </div>
                </div>
              </div>

              <!-- Chapter Input -->
              <div>
                <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1.5">Chapter</label>
                <input type="number" v-model="chapterInput" placeholder="e.g. 3" min="1" class="w-full bg-parchment-bg border border-parchment-border/70 px-3.5 py-2.5 rounded-xl text-xs text-parchment-neutral outline-none"/>
              </div>

              <!-- Verse Range -->
              <div>
                <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1.5">Verse Range</label>
                <div class="flex items-center space-x-2">
                  <input type="number" v-model="verseStart" placeholder="Start" min="1" class="w-full bg-parchment-bg border border-parchment-border/70 px-2.5 py-2 rounded-xl text-xs text-center text-parchment-neutral outline-none"/>
                  <span class="text-parchment-neutral/40 text-xs font-bold">—</span>
                  <input type="number" v-model="verseEnd" placeholder="End" min="1" class="w-full bg-parchment-bg border border-parchment-border/70 px-2.5 py-2 rounded-xl text-xs text-center text-parchment-neutral outline-none"/>
                </div>
              </div>
            </div>

            <!-- Direct Search -->
            <div v-else class="space-y-3">
              <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-1.5">Passage Citation</label>
              <input type="text" v-model="manualQuery" @keyup.enter="getContext" placeholder="e.g. John 3:16" class="w-full bg-parchment-bg border border-parchment-border/70 px-3.5 py-2.5 rounded-xl text-xs text-parchment-neutral outline-none"/>
            </div>

            <div v-if="error" class="text-xs text-[#8B2635] font-semibold mt-3">
              {{ error }}
            </div>

            <div class="mt-5 pt-4 border-t border-parchment-border/30">
              <AppButton variant="primary" :disabled="loading" @click="getContext" class="w-full text-xs py-2.5">
                <span v-if="loading">Exploring...</span>
                <span v-else>Explore Context</span>
              </AppButton>
            </div>
          </ParchmentCard>
        </aside>

        <!-- Main Content Area -->
        <div class="lg:col-span-8 xl:col-span-8 flex flex-col gap-6 animate-fade-in-up delay-150">

      <!-- Custom Loading Animation Overlay / Section -->
      <div v-if="loading" class="flex-grow flex flex-col items-center justify-center py-20 space-y-5 animate-fade-in-up">
        <!-- Elegant pulsing book/cross loader -->
        <div class="relative flex items-center justify-center">
          <div class="w-16 h-16 border-2 border-parchment-primary/30 border-t-parchment-primary rounded-full animate-spin"></div>
          <div class="absolute w-8 h-8 text-parchment-primary animate-pulse">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
        </div>
        <div class="text-center space-y-1">
          <p class="text-parchment-neutral/80 font-medium text-sm">Searching the Sacred Scriptures...</p>
          <p class="text-parchment-neutral/40text-xs">Unearthing historical contexts and theological commentary</p>
        </div>
      </div>

      <!-- Saved History Section -->
      <div v-else class="space-y-4">
        <div class="flex items-center justify-between px-2 mb-2">
          <h3 class="text-xs font-bold text-parchment-neutral/50 uppercase tracking-widest">
            Search History ({{ savedVerses.length }})
          </h3>
          <button 
            v-if="savedVerses.length > 0" 
            @click="clearAllHistory"
            class="text-[10px] font-bold text-parchment-secondary uppercase tracking-wider hover:opacity-85 transition-opacity bg-transparent border-none p-0 outline-none hover:bg-transparent"
          >
            Clear All
          </button>
        </div>

        <div v-if="savedVerses.length === 0" class="text-center py-16 border border-dashed border-parchment-border rounded-3xl p-8 bg-parchment-neutral-light/20">
          <svg class="mx-auto text-parchment-neutral/30 mb-3" xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"></path>
            <path d="M12 8v4l3 3"></path>
          </svg>
          <p class="text-sm font-seriftext-parchment-neutral/50">Your Scripture history is empty.</p>
          <p class="text-xs text-parchment-neutral/40 mt-1">Select a passage above to research its context.</p>
        </div>

        <!-- Interactive Collapsible Cards -->
        <transition-group name="list" tag="div" class="space-y-4">
          <div 
            v-for="item in savedVerses" 
            :key="item.id"
            class="overflow-hidden"
          >
            <ParchmentCard 
              @click="toggleExpand(item.id)"
              :hover-effect="!expandedIds.has(item.id)"
              class="relative !p-6 shadow-sm border border-parchment-border/70 hover:border-parchment-primary/30 transition-all duration-300 select-none"
              :class="{ '!bg-parchment-neutral-light/40': expandedIds.has(item.id) }"
            >
              <!-- Card Header -->
              <header class="flex items-start justify-between gap-4">
                <div class="text-left">
                  <div class="flex items-center gap-2 mb-1 flex-wrap">
                    <span class="text-[9px] font-bold text-parchment-primary uppercase tracking-widest">
                      {{ formatDate(item.timestamp) }}
                    </span>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100/80 text-amber-900 border border-amber-300/60">
                      📖 Sacred Scripture
                    </span>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-50 text-blue-900 border border-blue-200/80">
                      🤖 AI-Generated
                    </span>
                  </div>
                  <h3 class="font-serif text-lg md:text-xl font-bold text-parchment-primary-dark">
                    {{ item.citation }}
                  </h3>
                </div>

                <div class="flex items-center space-x-2">
                  <!-- Delete button -->
                  <button 
                    @click="deleteVerse(item.id, $event)"
                    class="text-parchment-neutral/30 hover:text-parchment-secondary hover:bg-parchment-secondary/10 p-1.5 rounded-full transition-colors border border-transparent bg-transparent outline-none hover:translate-y-0"
                    title="Delete item"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M3 6h18"></path>
                      <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
                      <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path>
                    </svg>
                  </button>

                  <!-- Expand/Collapse Chevron -->
                  <div 
                    class="text-parchment-neutral/50 transition-transform duration-300"
                    :class="{ 'rotate-180 text-parchment-primary-dark': expandedIds.has(item.id) }"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="m6 9 6 6 6-6"></path>
                    </svg>
                  </div>
                </div>
              </header>

              <!-- Collapsed Snippet -->
              <div v-if="!expandedIds.has(item.id)" class="text-left mt-3">
                <p class="font-seriftext-sm text-parchment-neutral/70 leading-relaxed">
                  "{{ getSnippet(item.content) }}"
                </p>
              </div>

              <!-- Expanded Content Block -->
              <transition name="expand">
                <div v-if="expandedIds.has(item.id)" class="mt-5 space-y-5 border-t border-parchment-border/40 pt-4 cursor-default" @click.stop>
                  <!-- Scripture Text -->
                  <div>
                    <div class="flex items-center justify-between mb-2">
                      <h4 class="text-[9px] font-bold text-parchment-neutral/40 uppercase tracking-widest">Scripture Text</h4>
                      <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100/80 text-amber-900 border border-amber-300/60">
                        📖 Sacred Scripture
                      </span>
                    </div>
                    <div class="prose max-w-none">
                      <blockquote class="border-l-2 border-parchment-border pl-4 font-serif text-base md:text-lg text-parchment-neutral/95 leading-relaxed whitespace-pre-line" v-html="item.content"></blockquote>
                      <p v-if="item.copyright" class="text-[9px] text-parchment-neutral/35 mt-2 ml-4">
                        {{ item.copyright }}
                      </p>
                    </div>
                  </div>

                  <!-- Theological Context -->
                  <div class="bg-parchment-bg border-l-4 border-parchment-primary rounded-2xl p-5 shadow-inner">
                    <div class="flex items-center justify-between mb-2">
                      <div class="flex items-center gap-1.5">
                        <span class="text-parchment-primary-dark">
                          <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="12" cy="12" r="10"></circle>
                            <line x1="12" x2="12" y1="16" y2="12"></line>
                            <line x1="12" x2="12.01" y1="8" y2="8"></line>
                          </svg>
                        </span>
                        <h4 class="text-[10px] font-bold text-parchment-primary-dark uppercase tracking-wider">Historical & Theological Commentary</h4>
                      </div>
                      <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-50 text-blue-900 border border-blue-200/80">
                        🤖 AI-Generated
                      </span>
                    </div>
                    <div class="prose max-w-none text-parchment-neutral/85" v-html="renderMarkdown(cleanContext(item.context))"></div>
                    <FactCheckCitations :citations="item.citations" :rawContextHtml="item.context" />
                  </div>
                </div>
              </transition>
            </ParchmentCard>
          </div>
        </transition-group>
      </div>
        </div>
      </main>
    </div>

    <!-- Global Footer -->
    <BottomNav />

    <!-- Floating Catechism Bubble -->
    <CatechismBubble />
  </div>
</template>

<style scoped>
.animate-fade-in-down {
  animation: fadeInDown 0.6s ease-out forwards;
}
.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out forwards;
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

/* Custom list items transition animations */
.list-enter-active,
.list-leave-active {
  transition: all 0.4s ease;
}
.list-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}
.list-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

/* Accordion expand transition */
.expand-enter-active,
.expand-leave-active {
  transition: all 0.3s ease-out;
  max-height: 1000px;
  opacity: 1;
}
.expand-enter-from,
.expand-leave-to {
  max-height: 0;
  opacity: 0;
  overflow: hidden;
}
</style>
