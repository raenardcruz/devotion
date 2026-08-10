<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useDate } from '../composables/useDate';
import { useDevotionApi } from '../composables/useDevotionApi';
import { renderMarkdown } from '../utils/markdown';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import AppTabs from '../components/common/AppTabs.vue';
import CatechismBubble from '../components/catechism/CatechismBubble.vue';
import FactCheckCitations, { type MagisteriumCitation } from '../components/common/FactCheckCitations.vue';

interface Reading {
  citation: string;
  text: string;
  context?: string;
  citations?: MagisteriumCitation[];
}

interface ResponsorialPsalm {
  citation: string;
  text: string;
  context?: string;
  citations?: MagisteriumCitation[];
}

interface MassReadings {
  first_reading: Reading;
  responsorial_psalm: ResponsorialPsalm;
  second_reading?: Reading;
  gospel: Reading;
  pope_quote?: string;
}

const readings = ref<MassReadings | null>(null);
const activeTab = ref('first_reading');

const { getLocalISOString } = useDate();
const { getDevotion, loading, error } = useDevotionApi();
const selectedDate = ref(getLocalISOString());
const modalDate = ref(getLocalISOString());
const isModalOpen = ref(false);

const openModal = () => {
  modalDate.value = selectedDate.value;
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
};

const handleApplyDate = () => {
  selectedDate.value = modalDate.value;
  isModalOpen.value = false;
  fetchReadings();
};

const contextMarkdownOptions = {
  paragraphClass: 'mb-2 last:mb-0',
  listClass: 'list-disc pl-4 space-y-1 my-2',
  listItemClass: 'mb-1 last:mb-0',
};

const cleanContext = (contextStr?: string) => {
  if (!contextStr) return '';
  return contextStr.replace(/<details[\s\S]*?<\/details>/gi, '').trim();
};

// Set initial loading to true for onMounted fetch
loading.value = true;

const fetchReadings = async () => {
  try {
    const data = await getDevotion<MassReadings>(selectedDate.value);
    readings.value = data;
    
    // Set default active tab based on what's available
    if (data.first_reading) {
      activeTab.value = 'first_reading';
    } else if (data.gospel) {
      activeTab.value = 'gospel';
    }
  } catch (err: any) {
    console.error("Error fetching readings:", err);
  }
};

const selectedDateValue = computed(() => new Date(`${selectedDate.value}T12:00:00`));

onMounted(() => {
  fetchReadings();
});

// Dynamic date strings for header
const formattedDate = computed(() => {
  return selectedDateValue.value.toLocaleDateString('en-US', {
    weekday: 'long', 
    month: 'long', 
    day: 'numeric' 
  });
});

const formattedWeek = computed(() => {
  // Format details like: "Wednesday of the Liturgical Week"
  const weekday = selectedDateValue.value.toLocaleDateString('en-US', { weekday: 'long' });
  return `${weekday} of the Liturgical Year — Daily Reflection`;
});

// Build tabs list dynamically
const availableTabs = computed(() => {
  const tabs = [];
  if (!readings.value) return [];
  
  if (readings.value.first_reading) {
    tabs.push({ id: 'first_reading', label: 'First Reading' });
  }
  if (readings.value.second_reading) {
    tabs.push({ id: 'second_reading', label: 'Second Reading' });
  }
  if (readings.value.responsorial_psalm) {
    tabs.push({ id: 'responsorial_psalm', label: 'Psalm' });
  }
  if (readings.value.gospel) {
    tabs.push({ id: 'gospel', label: 'Gospel' });
  }
  if (readings.value.pope_quote) {
    tabs.push({ id: 'pope_quote', label: 'Words of the Popes' });
  }
  return tabs;
});
</script>

<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-28 sm:pb-20 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Readings Container -->
    <main class="flex-grow max-w-3xl mx-auto w-full px-4 py-8 flex flex-col">
      
      <!-- Page Header -->
      <header class="text-center mb-8 animate-fade-in-down">
          <span class="text-[#72383D] font-bold tracking-[0.25em] text-xs uppercase block mb-2">
            Liturgical Cycle
          </span>
          <h1 class="text-3xl md:text-5xl font-serif font-bold my-2 bg-gradient-to-r from-[#72383D] via-[#AC9C8D] to-[#322D29] bg-clip-text text-transparent">
            {{ formattedDate }}
          </h1>
          <p class="text-[#322D29]/70 text-sm">
            {{ formattedWeek }}
          </p>
          <div class="mt-5 flex items-center justify-center">
            <button
              @click="openModal"
              class="inline-flex items-center gap-2 px-4 py-2 bg-white/80 hover:bg-white text-[#72383D] border border-[#D1C7BD] rounded-xl text-xs font-bold uppercase tracking-wider transition-all shadow-xs hover:shadow-sm focus:outline-none focus:ring-2 focus:ring-[#72383D]/30 active:scale-95 cursor-pointer"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="16" y1="2" x2="16" y2="6"></line>
                <line x1="8" y1="2" x2="8" y2="6"></line>
                <line x1="3" y1="10" x2="21" y2="10"></line>
              </svg>
              <span>Change Date</span>
            </button>
          </div>
      </header>

      <!-- Content Area -->
      <div v-if="loading" class="flex-grow flex flex-col items-center justify-center py-20 space-y-4">
         <div class="w-12 h-12 border-4 border-[#D1C7BD] border-t-[#72383D] rounded-full animate-spin"></div>
         <p class="text-[#322D29]/60 animate-pulse font-medium text-sm">Loading the liturgy...</p>
      </div>

      <div v-else-if="error" class="bg-red-50/50 border border-red-200 rounded-2xl p-8 text-center text-red-800 animate-fade-in-up">
        <p class="font-bold mb-2">Liturgy Fetch Error</p>
        <p class="text-sm opacity-80">{{ error }}</p>
        <button @click="fetchReadings" class="mt-4 px-5 py-2.5 bg-red-100 hover:bg-red-200 text-red-800 rounded-full transition-colors text-xs font-bold uppercase tracking-wider">
            Try Again
        </button>
      </div>

      <div v-else-if="readings" class="flex-grow flex flex-col gap-6 animate-fade-in-up delay-150">
        
        <!-- Segmented Navigation Tabs -->
        <AppTabs :tabs="availableTabs" v-model="activeTab" />

        <!-- Readings Cards -->
        <transition name="fade-slide" mode="out-in">
          <div :key="activeTab" class="w-full">
            
            <!-- First Reading Card -->
            <ParchmentCard v-if="activeTab === 'first_reading' && readings.first_reading" class="shadow-sm">
              <header class="mb-6 border-b border-[#D1C7BD]/60 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-[#72383D] font-bold tracking-[0.2em] text-[10px] uppercase block">First Reading</span>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">📖 Sacred Scripture</span>
                  </div>
                  <router-link
                    :to="{ path: '/bible-study', query: { query: readings.first_reading.citation } }"
                    class="inline-flex items-center space-x-2 group cursor-pointer"
                  >
                    <h2 class="text-2xl md:text-3xl font-serif text-[#72383D] group-hover:text-[#322D29] transition-colors font-bold">{{ readings.first_reading.citation }}</h2>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" x2="21" y1="14" y2="3"></line>
                    </svg>
                  </router-link>
              </header>
              
              <div class="prose max-w-none mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-[#322D29]" v-html="readings.first_reading.text"></div>
              </div>

              <div v-if="readings.first_reading.context" class="bg-white/60 border-l-4 border-[#72383D] rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-[#72383D] text-[10px] font-bold uppercase tracking-wider">Contextual Meditation</h3>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-[#322D29]/80 text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.first_reading.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.first_reading.citations" :rawContextHtml="readings.first_reading.context" />
              </div>
            </ParchmentCard>

            <!-- Second Reading Card -->
            <ParchmentCard v-if="activeTab === 'second_reading' && readings.second_reading" class="shadow-sm">
              <header class="mb-6 border-b border-[#D1C7BD]/60 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-[#72383D] font-bold tracking-[0.2em] text-[10px] uppercase block">Second Reading</span>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">📖 Sacred Scripture</span>
                  </div>
                  <router-link
                    :to="{ path: '/bible-study', query: { query: readings.second_reading.citation } }"
                    class="inline-flex items-center space-x-2 group cursor-pointer"
                  >
                    <h2 class="text-2xl md:text-3xl font-serif text-[#72383D] group-hover:text-[#322D29] transition-colors font-bold">{{ readings.second_reading.citation }}</h2>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" x2="21" y1="14" y2="3"></line>
                    </svg>
                  </router-link>
              </header>
              
              <div class="prose max-w-none mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-[#322D29]" v-html="readings.second_reading.text"></div>
              </div>

              <div v-if="readings.second_reading.context" class="bg-white/60 border-l-4 border-[#72383D] rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-[#72383D] text-[10px] font-bold uppercase tracking-wider">Contextual Meditation</h3>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-[#322D29]/80 text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.second_reading.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.second_reading.citations" :rawContextHtml="readings.second_reading.context" />
              </div>
            </ParchmentCard>

            <!-- Responsorial Psalm Card -->
            <ParchmentCard v-if="activeTab === 'responsorial_psalm' && readings.responsorial_psalm" class="shadow-sm">
              <header class="mb-6 border-b border-[#D1C7BD]/60 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-[#72383D] font-bold tracking-[0.2em] text-[10px] uppercase block">Responsorial Psalm</span>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">📖 Sacred Scripture</span>
                  </div>
                  <div class="text-center">
                    <router-link
                      :to="{ path: '/bible-study', query: { query: readings.responsorial_psalm.citation } }"
                      class="inline-flex items-center space-x-2 group cursor-pointer"
                    >
                      <h2 class="text-xl md:text-2xl font-serif text-[#72383D] group-hover:text-[#322D29] transition-colors font-bold">{{ readings.responsorial_psalm.citation }}</h2>
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors">
                        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                        <polyline points="15 3 21 3 21 9"></polyline>
                        <line x1="10" x2="21" y1="14" y2="3"></line>
                      </svg>
                    </router-link>
                  </div>
              </header>
              
              <div class="prose max-w-none text-center mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg font-serif text-[#322D29] leading-loose" v-html="readings.responsorial_psalm.text"></div>
              </div>

              <div v-if="readings.responsorial_psalm.context" class="bg-white/60 border-l-4 border-[#72383D] rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-[#72383D] text-[10px] font-bold uppercase tracking-wider">Psalm Context</h3>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-[#322D29]/80 text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.responsorial_psalm.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.responsorial_psalm.citations" :rawContextHtml="readings.responsorial_psalm.context" />
              </div>
            </ParchmentCard>

            <!-- Gospel Card -->
            <ParchmentCard v-if="activeTab === 'gospel' && readings.gospel" class="shadow-sm">
              <header class="mb-6 border-b border-[#D1C7BD]/60 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-[#72383D] font-bold tracking-[0.2em] text-[10px] uppercase block">Gospel Reading</span>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">📖 Sacred Gospel</span>
                  </div>
                  <router-link
                    :to="{ path: '/bible-study', query: { query: readings.gospel.citation } }"
                    class="inline-flex items-center space-x-2 group cursor-pointer"
                  >
                    <h2 class="text-2xl md:text-3xl font-serif text-[#72383D] group-hover:text-[#322D29] transition-colors font-bold">{{ readings.gospel.citation }}</h2>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" x2="21" y1="14" y2="3"></line>
                    </svg>
                  </router-link>
              </header>
              
              <div class="prose max-w-none mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-[#322D29]" v-html="readings.gospel.text"></div>
              </div>

              <div v-if="readings.gospel.context" class="bg-white/60 border-l-4 border-[#72383D] rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-[#72383D] text-[10px] font-bold uppercase tracking-wider">Gospel Reflection</h3>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-[#322D29]/80 text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.gospel.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.gospel.citations" :rawContextHtml="readings.gospel.context" />
              </div>
            </ParchmentCard>

            <!-- Pope Quote/Reflection Card -->
            <ParchmentCard v-if="activeTab === 'pope_quote' && readings.pope_quote" class="shadow-sm">
              <header class="mb-6 border-b border-[#D1C7BD]/60 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-[#72383D] font-bold tracking-[0.2em] text-[10px] uppercase block">Words of the Popes</span>
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">🏛️ Papal Magisterium</span>
                  </div>
                  <h2 class="text-xl md:text-2xl font-serif text-[#72383D] font-bold">Holy Father's Reflection</h2>
              </header>
              
              <div class="prose max-w-none mb-6">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-[#322D29]" v-html="readings.pope_quote"></div>
              </div>

              <div class="text-center text-[10px] text-[#322D29]/50 tracking-[0.1em] border-t border-[#D1C7BD]/60 pt-4">
                Source: Vatican News Daily Liturgy
              </div>
            </ParchmentCard>

          </div>
        </transition>

        <!-- Pastoral & Sacramental Engagement Callout -->
        <div class="bg-white/60 border border-[#D1C7BD]/60 rounded-2xl p-5 text-center shadow-xs mt-2 backdrop-blur-md">
          <div class="flex items-center justify-center space-x-2 text-[#72383D] font-bold text-xs uppercase tracking-wider mb-1">
            <span>Pastoral Call to Community & Prayer</span>
          </div>
          <p class="text-[#322D29]/80 text-xs md:text-sm font-serif leading-relaxed">
            "Digital tools serve to support human reflection, but can never replace personal prayer, attendance at the Holy Sacrifice of the Mass, or the grace of the Sacraments in your local parish community."
          </p>
          <span class="text-[10px] text-[#72383D] block mt-1.5 uppercase font-sans font-semibold">
            Guided by Pope Leo XIV's Encyclical <em>Magnifica Humanitas</em>
          </span>
        </div>

        <!-- Copyright Footnote -->
        <span class="text-center text-parchment-neutral/30 text-[10px] leading-relaxed max-w-2xl mx-auto mt-4">
          Holy Bible, New Living Translation Catholic Edition, copyright © 2016 by Tyndale House Foundation. All rights reserved. Used by permission of Tyndale House Publishers.
        </span>

      </div>
    </main>

    <!-- Global Footer -->
    <BottomNav />

    <!-- Floating Catechism Bubble -->
    <CatechismBubble />

    <!-- Date Picker Modal -->
    <Teleport to="body">
      <transition name="fade">
        <div
          v-if="isModalOpen"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-xs"
          @click.self="closeModal"
        >
          <div class="bg-[#FBF8F3] border border-[#D1C7BD] rounded-2xl p-6 max-w-sm w-full shadow-xl animate-fade-in-up">
            <div class="flex items-center justify-between mb-4 border-b border-[#D1C7BD]/60 pb-3">
              <h3 class="text-lg font-serif font-bold text-[#72383D]">Choose Liturgical Date</h3>
              <button
                @click="closeModal"
                class="text-[#322D29]/50 hover:text-[#322D29] transition-colors p-1 rounded-lg cursor-pointer"
                aria-label="Close modal"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <p class="text-xs text-[#322D29]/70 mb-4">
              Select a date to view readings and contextual meditations.
            </p>

            <form @submit.prevent="handleApplyDate" class="space-y-4">
              <div>
                <label for="modal-reading-date" class="block text-xs font-bold uppercase tracking-wider text-[#322D29]/70 mb-1.5">
                  Date
                </label>
                <input
                  id="modal-reading-date"
                  v-model="modalDate"
                  type="date"
                  class="w-full rounded-xl border border-[#D1C7BD] bg-white px-3 py-2.5 text-sm text-[#322D29] shadow-xs focus:border-[#72383D] focus:outline-none focus:ring-2 focus:ring-[#72383D]/20"
                >
              </div>

              <div class="flex items-center justify-end gap-2 pt-2">
                <button
                  type="button"
                  @click="closeModal"
                  class="px-4 py-2 bg-transparent hover:bg-[#EFE9E1] text-[#322D29]/80 rounded-xl text-xs font-bold uppercase tracking-wider transition-colors cursor-pointer"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  class="px-5 py-2 bg-[#72383D] hover:bg-[#5c2d31] text-white rounded-xl text-xs font-bold uppercase tracking-wider transition-colors shadow-xs focus:outline-none focus:ring-2 focus:ring-[#72383D]/40 active:scale-95 cursor-pointer"
                >
                  Submit
                </button>
              </div>
            </form>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped>
.animate-fade-in-down {
  animation: fadeInDown 0.6s ease-out forwards;
}
.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out forwards;
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

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.4s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
