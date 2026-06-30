<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useDevotionApi } from '../composables/useDevotionApi';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import AppButton from '../components/common/AppButton.vue';
import AppTabs from '../components/common/AppTabs.vue';
import CatechismBubble from '../components/catechism/CatechismBubble.vue';

interface SavedVerse {
  id: string;
  citation: string;
  content: string;
  context: string;
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
}

const { fetchWithAuth } = useDevotionApi();

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
onMounted(() => {
  try {
    const raw = localStorage.getItem('saved_verses');
    if (raw) {
      savedVerses.value = JSON.parse(raw);
      // Auto-expand the first item if there is any history
      if (savedVerses.value.length > 0) {
        const firstItem = savedVerses.value[0];
        if (firstItem) {
          expandedIds.value.add(firstItem.id);
        }
      }
    }
  } catch (err) {
    console.error('Error parsing saved verses from localStorage:', err);
  }
});

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

const renderMarkdown = (markdown: string): string => {
  if (!markdown) return '';
  
  let html = markdown
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
  
  // Replace headers
  html = html.replace(/^### (.*?)$/gm, '<h5 class="text-xs font-bold text-parchment-primary-dark mt-4 mb-2 font-serif uppercase tracking-wider">$1</h5>');
  html = html.replace(/^## (.*?)$/gm, '<h4 class="text-sm font-bold text-parchment-primary-dark mt-5 mb-2 font-serif">$1</h4>');
  html = html.replace(/^# (.*?)$/gm, '<h3 class="text-base font-bold text-parchment-primary-dark mt-6 mb-3 font-serif">$1</h3>');
  
  // Replace bold
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong class="font-bold text-parchment-primary-dark">$1</strong>');
  html = html.replace(/__(.*?)__/g, '<strong class="font-bold text-parchment-primary-dark">$1</strong>');
  
  // Replace italic
  html = html.replace(/\*(.*?)\*/g, '<em class="italic">$1</em>');
  html = html.replace(/_(.*?)_/g, '<em class="italic">$1</em>');
  
  // Replace code
  html = html.replace(/`(.*?)`/g, '<code class="bg-parchment-neutral-light/70 px-1 py-0.5 rounded font-mono text-xs text-parchment-secondary">$1</code>');
  
  // Replace blockquotes
  html = html.replace(/^&gt; (.*?)$/gm, '<blockquote class="border-l-2 border-parchment-primary/40 pl-3 italic text-parchment-neutral/80 my-2">$1</blockquote>');
  
  // Process line by line for paragraphs and lists
  const lines = html.split('\n');
  let inList = false;
  const processedLines: string[] = [];
  
  for (let line of lines) {
    const trimmed = line.trim();
    
    // Check for bullet list item
    const listMatch = trimmed.match(/^[\*\-]\s+(.*)$/);
    if (listMatch) {
      if (!inList) {
        processedLines.push('<ul class="list-disc pl-5 space-y-1.5 my-3">');
        inList = true;
      }
      processedLines.push(`<li class="text-sm md:text-[14px] text-parchment-neutral/85">${listMatch[1]}</li>`);
      continue;
    } else {
      if (inList) {
        processedLines.push('</ul>');
        inList = false;
      }
    }
    
    if (trimmed === '') {
      processedLines.push('');
      continue;
    }
    
    // If it's already an HTML block tag, keep it
    if (trimmed.startsWith('<h') || trimmed.startsWith('<blockquote') || trimmed.startsWith('</blockquote') || trimmed.startsWith('<ul') || trimmed.startsWith('</ul')) {
      processedLines.push(trimmed);
    } else {
      processedLines.push(`<p class="text-sm md:text-[15px] leading-relaxed text-parchment-neutral/80 mb-3">${trimmed}</p>`);
    }
  }
  
  if (inList) {
    processedLines.push('</ul>');
  }
  
  return processedLines.filter(l => l !== '').join('\n');
};
</script>

<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-24 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Layout Container -->
    <main class="flex-grow max-w-3xl mx-auto w-full px-4 py-8 flex flex-col">
      <!-- Title Header -->
      <header class="text-center mb-8 animate-fade-in-down">
        <span class="text-parchment-primary font-bold tracking-[0.25em] text-xs uppercase block mb-1">
          Scripture Study
        </span>
        <h1 class="text-3xl md:text-5xl font-serif text-parchment-primary-dark font-medium my-2">
          Context Explorer
        </h1>
        <p class="text-parchment-neutral/70 italic text-sm max-w-md mx-auto leading-relaxed">
          Unveil the historical context, theological layers, and original scriptures of any biblical passage.
        </p>
      </header>

      <!-- Selector & Input Section -->
      <ParchmentCard class="shadow-sm mb-8 animate-fade-in-up">
        <AppTabs :tabs="tabs" v-model="activeTab" class="mb-6" />

        <transition name="fade" mode="out-in">
          <div :key="activeTab">
            <!-- TAB 1: STRUCTURED SELECTOR -->
            <div v-if="activeTab === 'selector'" class="space-y-4">
              <div class="grid grid-cols-1 md:grid-cols-4 gap-4 items-end">
                <!-- Book Dropdown Input -->
                <div class="relative md:col-span-2">
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-2">
                    Bible Book
                  </label>
                  <button 
                    type="button"
                    @click="showBookDropdown = !showBookDropdown"
                    class="w-full text-left bg-parchment-bg border border-parchment-border/70 hover:border-parchment-primary/50 px-4 py-3 rounded-2xl flex items-center justify-between text-sm transition-all focus:outline-none focus:ring-1 focus:ring-parchment-primary"
                  >
                    <span class="font-medium text-parchment-neutral">{{ selectedBook || 'Select a Book' }}</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/50 transition-transform duration-200" :class="{ 'rotate-180': showBookDropdown }">
                      <path d="m6 9 6 6 6-6"></path>
                    </svg>
                  </button>

                  <!-- Custom Searchable Dropdown Overlay -->
                  <div v-if="showBookDropdown" class="absolute z-50 left-0 right-0 mt-2 bg-parchment-bg border border-parchment-border rounded-2xl shadow-xl max-h-72 overflow-y-auto p-2">
                    <div class="sticky top-0 bg-parchment-bg pb-2 pt-1 border-b border-parchment-border/30 mb-2">
                      <div class="relative">
                        <input 
                          type="text"
                          v-model="bookSearchQuery"
                          placeholder="Search books..."
                          class="w-full bg-parchment-neutral-light/50 border border-parchment-border/40 text-xs px-3 py-2 pl-8 rounded-xl focus:outline-none focus:border-parchment-primary text-parchment-neutral"
                        />
                        <svg class="absolute left-2.5 top-2.5 text-parchment-neutral/40" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="11" cy="11" r="8"></circle>
                          <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                        </svg>
                      </div>
                    </div>
                    
                    <div v-if="filteredBooks.length === 0" class="text-center text-xs text-parchment-neutral/40 py-4">
                      No books found
                    </div>
                    
                    <div v-else-if="!bookSearchQuery" class="space-y-3">
                      <div v-for="cat in categories" :key="cat.name">
                        <div class="text-[9px] font-bold text-parchment-primary uppercase tracking-wider px-2 mb-1">{{ cat.name }}</div>
                        <div class="grid grid-cols-2 gap-1">
                          <button 
                            v-for="book in cat.books" 
                            :key="book"
                            @click="selectBook(book)"
                            class="text-left text-xs px-2.5 py-1.5 rounded-lg hover:bg-parchment-neutral-light transition-all text-parchment-neutral"
                            :class="{ '!bg-parchment-primary-dark !text-white font-bold': selectedBook === book }"
                          >
                            {{ book }}
                          </button>
                        </div>
                      </div>
                    </div>

                    <div v-else class="grid grid-cols-2 gap-1">
                      <button 
                        v-for="book in filteredBooks" 
                        :key="book"
                        @click="selectBook(book)"
                        class="text-left text-xs px-2.5 py-1.5 rounded-lg hover:bg-parchment-neutral-light transition-all text-parchment-neutral"
                        :class="{ '!bg-parchment-primary-dark !text-white font-bold': selectedBook === book }"
                      >
                        {{ book }}
                      </button>
                    </div>
                  </div>
                </div>

                <!-- Chapter Input -->
                <div>
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-2">
                    Chapter
                  </label>
                  <input 
                    type="number"
                    v-model="chapterInput"
                    placeholder="Ch"
                    min="1"
                    class="w-full bg-parchment-bg border border-parchment-border/70 focus:border-parchment-primary/60 px-4 py-3 rounded-2xl text-sm text-center text-parchment-neutral focus:outline-none"
                  />
                </div>

                <!-- Verse Selector -->
                <div>
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-2">
                    Verse(s)
                  </label>
                  <div class="flex items-center space-x-1.5">
                    <input 
                      type="number"
                      v-model="verseStart"
                      placeholder="Start"
                      min="1"
                      class="w-full bg-parchment-bg border border-parchment-border/70 focus:border-parchment-primary/60 px-2 py-3 rounded-2xl text-sm text-center text-parchment-neutral focus:outline-none"
                    />
                    <span class="text-parchment-neutral/40 text-xs font-bold">—</span>
                    <input 
                      type="number"
                      v-model="verseEnd"
                      placeholder="End"
                      min="1"
                      class="w-full bg-parchment-bg border border-parchment-border/70 focus:border-parchment-primary/60 px-2 py-3 rounded-2xl text-sm text-center text-parchment-neutral focus:outline-none"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- TAB 2: MANUAL SEARCH -->
            <div v-else class="space-y-2">
              <label class="block text-[11px] font-bold uppercase tracking-wider text-parchment-neutral/60 mb-2">
                Bible Passage Citation
              </label>
              <div class="relative">
                <input 
                  type="text"
                  v-model="manualQuery"
                  @keyup.enter="getContext"
                  placeholder="e.g. John 3:16, Genesis 1:1-5, Romans 8:28"
                  class="w-full bg-parchment-bg border border-parchment-border/70 focus:border-parchment-primary/60 px-4 py-3 rounded-2xl text-sm text-parchment-neutral focus:outline-none pl-11 shadow-inner"
                />
                <svg class="absolute left-4 top-3.5 text-parchment-neutral/40" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                </svg>
              </div>
            </div>
          </div>
        </transition>

        <!-- Action Button & Inline Error -->
        <div class="mt-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pt-4 border-t border-parchment-border/30">
          <div v-if="error" class="text-xs text-[#8B2635] font-semibold flex items-center gap-1.5">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" x2="12" y1="8" y2="12"></line>
              <line x1="12" x2="12.01" y1="16" y2="16"></line>
            </svg>
            {{ error }}
          </div>
          <div v-else class="text-xs text-parchment-neutral/40 italic">
            * Fetches scripture translation and Catholic theological commentary.
          </div>
          <AppButton 
            variant="primary" 
            @click="getContext"
            :disabled="loading"
            class="sm:w-auto w-full justify-center flex items-center gap-2 shadow-sm hover:shadow active:scale-95 transition-all !px-6"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" x2="12" y1="15" y2="3"></line>
            </svg>
            Fetch Context
          </AppButton>
        </div>
      </ParchmentCard>

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
          <p class="text-parchment-neutral/40 italic text-xs">Unearthing historical contexts and theological commentary</p>
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
          <p class="text-sm font-serif italic text-parchment-neutral/50">Your Scripture history is empty.</p>
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
                  <span class="text-[9px] font-bold text-parchment-primary uppercase tracking-widest block mb-0.5">
                    {{ formatDate(item.timestamp) }}
                  </span>
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
                <p class="font-serif italic text-sm text-parchment-neutral/70 leading-relaxed">
                  "{{ getSnippet(item.content) }}"
                </p>
              </div>

              <!-- Expanded Content Block -->
              <transition name="expand">
                <div v-if="expandedIds.has(item.id)" class="mt-5 space-y-5 border-t border-parchment-border/40 pt-4 cursor-default" @click.stop>
                  <!-- Scripture Text -->
                  <div>
                    <h4 class="text-[9px] font-bold text-parchment-neutral/40 uppercase tracking-widest mb-2">Scripture Text</h4>
                    <div class="prose max-w-none">
                      <blockquote class="border-l-2 border-parchment-border pl-4 font-serif text-base md:text-lg text-parchment-neutral/95 leading-relaxed whitespace-pre-line" v-html="item.content"></blockquote>
                      <p v-if="item.copyright" class="text-[9px] text-parchment-neutral/35 mt-2 ml-4">
                        {{ item.copyright }}
                      </p>
                    </div>
                  </div>

                  <!-- Theological Context -->
                  <div class="bg-parchment-bg border-l-4 border-parchment-primary rounded-2xl p-5 shadow-inner">
                    <div class="flex items-center gap-1.5 mb-2">
                      <span class="text-parchment-primary-dark">
                        <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <circle cx="12" cy="12" r="10"></circle>
                          <line x1="12" x2="12" y1="16" y2="12"></line>
                          <line x1="12" x2="12.01" y1="8" y2="8"></line>
                        </svg>
                      </span>
                      <h4 class="text-[10px] font-bold text-parchment-primary-dark uppercase tracking-wider">Historical & Theological Commentary</h4>
                    </div>
                    <div class="prose max-w-none text-parchment-neutral/85" v-html="renderMarkdown(item.context)"></div>
                  </div>
                </div>
              </transition>
            </ParchmentCard>
          </div>
        </transition-group>
      </div>
    </main>

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
