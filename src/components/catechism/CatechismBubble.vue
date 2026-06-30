<template>
  <div class="relative">
    <!-- Click outside handler: invisible overlay to close popup when clicking outside -->
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-40 bg-transparent" 
      @click="isOpen = false"
    ></div>

    <!-- Floating Chat Bubble Button -->
    <button 
      @click="toggleOpen"
      class="fixed bottom-24 right-6 z-50 w-14 h-14 bg-parchment-primary-dark hover:bg-[#4a3c1d] text-white rounded-full flex items-center justify-center shadow-xl hover:shadow-[0_0_20px_rgba(197,160,89,0.4)] hover:-translate-y-1 transform active:scale-95 transition-all duration-300 group border border-parchment-primary/30"
      :aria-label="isOpen ? 'Close Catechism Lookup' : 'Open Catechism Lookup'"
      :title="isOpen ? 'Close Catechism Lookup' : 'Quick Catechism Lookup'"
    >
      <!-- Chat Bubble / Book Icon -->
      <span class="relative">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="transition-transform duration-300" :class="{ 'rotate-90 scale-90': isOpen }">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>
        <span class="absolute -top-1 -right-1.5 bg-parchment-primary text-[9px] font-bold px-1 rounded-md scale-90 select-none group-hover:bg-white group-hover:text-parchment-primary-dark transition-colors duration-200">CCC</span>
      </span>
    </button>

    <!-- Popover Menu -->
    <transition name="popover-fade">
      <div 
        v-if="isOpen"
        class="fixed bottom-40 right-6 z-50 w-[340px] sm:w-[380px] bg-parchment-neutral-light border border-parchment-border rounded-3xl shadow-[0_20px_50px_rgba(92,75,37,0.15)] overflow-hidden flex flex-col max-h-[460px] animate-fade-in"
      >
        <!-- Popover Header -->
        <header class="px-5 py-4 border-b border-parchment-border/40 flex justify-between items-center bg-white/50 backdrop-blur-sm">
          <div class="flex items-center gap-2">
            <span class="text-parchment-primary-dark">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
                <path d="M6 6h10"></path>
                <path d="M6 10h10"></path>
              </svg>
            </span>
            <h3 class="text-sm font-serif font-bold text-parchment-primary-dark">Quick Catechism</h3>
          </div>
          
          <button 
            @click="isOpen = false" 
            class="p-1 hover:bg-parchment-neutral-light/50 rounded-full transition-colors text-parchment-neutral/40 hover:text-parchment-neutral border-none bg-transparent shadow-none"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" x2="6" y1="6" y2="18"></line>
              <line x1="6" x2="18" y1="6" y2="18"></line>
            </svg>
          </button>
        </header>

        <!-- Search Input Block -->
        <div class="p-4 bg-white/30 border-b border-parchment-border/20 flex gap-2">
          <div class="relative flex-grow">
            <input 
              v-model="searchQuery" 
              type="text" 
              inputmode="numeric"
              pattern="[0-9]*"
              placeholder="Enter paragraph (e.g. 748)..." 
              @keyup.enter="performSearch"
              class="w-full pl-9 pr-4 py-2.5 rounded-xl border border-parchment-border/60 bg-white/70 text-parchment-neutral placeholder-parchment-neutral/40 focus:border-parchment-primary focus:bg-white shadow-inner transition-all outline-none text-xs font-semibold"
            >
            <div class="absolute left-3 top-1/2 -translate-y-1/2 text-parchment-neutral/30">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="11" cy="11" r="8"></circle>
                <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
              </svg>
            </div>
          </div>
          <button 
            @click="performSearch"
            class="px-4 py-2.5 bg-parchment-primary-dark hover:bg-[#4a3c1d] text-white rounded-xl text-xs font-bold transition-all border border-parchment-primary/30 shadow-sm"
          >
            Find
          </button>
        </div>

        <!-- Popover Body (Content) -->
        <div class="flex-grow p-5 overflow-y-auto min-h-[160px] bg-[#FDFBF7]/60">
          
          <!-- History Back Button -->
          <div v-if="history.length > 0" class="mb-3">
            <button 
              @click="goBack"
              class="flex items-center gap-1 text-[11px] font-bold text-parchment-primary-dark hover:text-parchment-primary transition-colors bg-transparent border-none p-0 shadow-none hover:bg-transparent"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <line x1="19" x2="5" y1="12" y2="12"></line>
                <polyline points="12 19 5 12 12 5"></polyline>
              </svg>
              Back to §{{ history[history.length - 1] }}
            </button>
          </div>

          <!-- Paragraph Display -->
          <div v-if="currentParagraph" class="space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-parchment-secondary font-serif select-none">§ {{ currentParagraph.id }}</span>
              <button 
                @click="viewInFullCatechism"
                class="text-[9px] font-bold text-parchment-primary uppercase tracking-wider hover:opacity-80 transition-opacity bg-transparent border-none p-0 shadow-none hover:bg-transparent"
                title="View in full Catechism explorer"
              >
                Go to full page
              </button>
            </div>
            
            <div class="prose max-w-none text-xs md:text-sm font-serif leading-relaxed text-parchment-neutral/95">
              <!-- Using original CatechismText component for parsing references -->
              <CatechismText 
                :text="currentParagraph.text" 
                :paragraphs="paragraphs" 
                @show-reference="handleShowReference"
              />
            </div>
          </div>

          <!-- Empty/Welcome State -->
          <div v-else-if="!errorMsg" class="flex flex-col items-center justify-center text-center h-full py-8 space-y-3">
            <div class="w-12 h-12 bg-white/60 rounded-full flex items-center justify-center shadow-inner">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/40">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" x2="12" y1="16" y2="12"></line>
                <line x1="12" x2="12.01" y1="8" y2="8"></line>
              </svg>
            </div>
            <div class="space-y-1">
              <h4 class="text-xs font-serif font-bold text-parchment-neutral">Quick Lookup</h4>
              <p class="text-[10px] text-parchment-neutral/50 max-w-[220px] leading-normal">
                Enter a paragraph number (1 to 2865) to quickly consult the Catechism during your readings.
              </p>
            </div>
          </div>

          <!-- Error State -->
          <div v-else class="flex flex-col items-center justify-center text-center h-full py-8 space-y-2">
            <div class="w-10 h-10 bg-red-50 text-red-500 rounded-full flex items-center justify-center shadow-inner">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" x2="12" y1="8" y2="12"></line>
                <line x1="12" x2="12.01" y1="16" y2="16"></line>
              </svg>
            </div>
            <p class="text-xs text-red-700 font-semibold">{{ errorMsg }}</p>
          </div>
        </div>

        <!-- Popover Footer -->
        <footer class="px-5 py-3 bg-[#F5EFE6] border-t border-parchment-border/40 text-center">
          <span class="text-[9px] text-parchment-neutral/40 uppercase tracking-widest font-semibold">Catechism of the Catholic Church</span>
        </footer>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import catechismData from './catechism.json';
import CatechismText from './CatechismText.vue';

interface Paragraph {
  id: number;
  text: string;
}

const router = useRouter();
const paragraphs = ref<Paragraph[]>(catechismData as unknown as Paragraph[]);

const isOpen = ref(false);
const searchQuery = ref('');
const currentParagraph = ref<Paragraph | null>(null);
const errorMsg = ref('');
const history = ref<number[]>([]);

const toggleOpen = () => {
  isOpen.value = !isOpen.value;
  if (isOpen.value) {
    // Reset states on open
    searchQuery.value = '';
    currentParagraph.value = null;
    errorMsg.value = '';
    history.value = [];
  }
};

const performSearch = () => {
  const cleaned = searchQuery.value.trim();
  if (!cleaned) {
    currentParagraph.value = null;
    errorMsg.value = '';
    return;
  }

  const id = parseInt(cleaned);
  if (!isNaN(id)) {
    const p = paragraphs.value.find(p => p.id === id);
    if (p) {
      currentParagraph.value = p;
      errorMsg.value = '';
      history.value = []; // Reset history when executing a new search query manually
    } else {
      currentParagraph.value = null;
      errorMsg.value = 'Paragraph not found (1-2865)';
    }
  } else {
    currentParagraph.value = null;
    errorMsg.value = 'Please enter a valid number';
  }
};

// Handle clicks on referenced paragraphs in CatechismText
const handleShowReference = (resolvedParagraphs: Paragraph[]) => {
  if (resolvedParagraphs.length > 0) {
    const target = resolvedParagraphs[0];
    if (target) {
      if (currentParagraph.value) {
        history.value.push(currentParagraph.value.id);
      }
      searchQuery.value = String(target.id);
      currentParagraph.value = target;
      errorMsg.value = '';
    }
  }
};

const goBack = () => {
  if (history.value.length > 0) {
    const prevId = history.value.pop();
    if (prevId) {
      searchQuery.value = String(prevId);
      const p = paragraphs.value.find(p => p.id === prevId);
      if (p) {
        currentParagraph.value = p;
        errorMsg.value = '';
      }
    }
  }
};

const viewInFullCatechism = () => {
  if (currentParagraph.value) {
    isOpen.value = false;
    router.push({ path: '/catechism', query: { q: String(currentParagraph.value.id) } });
  }
};

// Auto search while typing full paragraph IDs (1-2865)
watch(searchQuery, (newVal) => {
  const cleaned = newVal.trim();
  if (/^\d+$/.test(cleaned)) {
    const id = parseInt(cleaned);
    if (id >= 1 && id <= 2865) {
      const p = paragraphs.value.find(p => p.id === id);
      if (p) {
        currentParagraph.value = p;
        errorMsg.value = '';
        return;
      }
    }
  }
  if (!cleaned) {
    currentParagraph.value = null;
    errorMsg.value = '';
  }
});
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.25s ease-out forwards;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.popover-fade-enter-active,
.popover-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.popover-fade-enter-from,
.popover-fade-leave-to {
  opacity: 0;
  transform: translateY(10px) scale(0.98);
}
</style>
