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
    const data = await getDevotion<MassReadings>(getLocalISOString());
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

onMounted(() => {
  fetchReadings();
});

// Dynamic date strings for header
const formattedDate = computed(() => {
  return new Date().toLocaleDateString('en-US', { 
    weekday: 'long', 
    month: 'long', 
    day: 'numeric' 
  });
});

const formattedWeek = computed(() => {
  // Format details like: "Wednesday of the Liturgical Week"
  const weekday = new Date().toLocaleDateString('en-US', { weekday: 'long' });
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
      
      <!-- Liturgical Header (Image 4 spec) -->
      <header class="text-center mb-8 animate-fade-in-down">
          <span class="text-parchment-secondary font-bold tracking-[0.25em] text-xs uppercase block mb-1">
            Ordinary Time
          </span>
          <h1 class="text-3xl md:text-5xl font-serif text-parchment-primary-dark font-medium my-2">
            {{ formattedDate }}
          </h1>
          <p class="text-parchment-neutral/70text-sm">
            {{ formattedWeek }}
          </p>
      </header>

      <!-- Content Area -->
      <div v-if="loading" class="flex-grow flex flex-col items-center justify-center py-20 space-y-4">
         <div class="w-12 h-12 border-4 border-parchment-border border-t-parchment-primary rounded-full animate-spin"></div>
         <p class="text-parchment-neutral/50 animate-pulse font-medium text-sm">Loading today's liturgy...</p>
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
              <header class="mb-6 border-b border-parchment-border/40 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-parchment-primary font-bold tracking-[0.2em] text-[10px] uppercase block">First Reading</span>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100/80 text-amber-900 border border-amber-300/60">📖 Sacred Scripture</span>
                  </div>
                  <router-link
                    :to="{ path: '/bible-study', query: { query: readings.first_reading.citation } }"
                    class="inline-flex items-center space-x-2 group cursor-pointer"
                  >
                    <h2 class="text-2xl md:text-3xl font-serif text-parchment-neutral group-hover:text-parchment-primary-dark transition-colors font-medium">{{ readings.first_reading.citation }}</h2>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/30 group-hover:text-parchment-primary transition-colors">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" x2="21" y1="14" y2="3"></line>
                    </svg>
                  </router-link>
              </header>
              
              <div class="prose max-w-none mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-parchment-neutral/90" v-html="readings.first_reading.text"></div>
              </div>

              <div v-if="readings.first_reading.context" class="bg-parchment-bg border-l-4 border-parchment-primary rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-parchment-primary-dark text-[10px] font-bold uppercase tracking-wider">Contextual Meditation</h3>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-50 text-blue-900 border border-blue-200/80">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-parchment-neutral/70text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.first_reading.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.first_reading.citations" :rawContextHtml="readings.first_reading.context" />
              </div>
            </ParchmentCard>

            <!-- Second Reading Card -->
            <ParchmentCard v-if="activeTab === 'second_reading' && readings.second_reading" class="shadow-sm">
              <header class="mb-6 border-b border-parchment-border/40 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-parchment-primary font-bold tracking-[0.2em] text-[10px] uppercase block">Second Reading</span>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100/80 text-amber-900 border border-amber-300/60">📖 Sacred Scripture</span>
                  </div>
                  <router-link
                    :to="{ path: '/bible-study', query: { query: readings.second_reading.citation } }"
                    class="inline-flex items-center space-x-2 group cursor-pointer"
                  >
                    <h2 class="text-2xl md:text-3xl font-serif text-parchment-neutral group-hover:text-parchment-primary-dark transition-colors font-medium">{{ readings.second_reading.citation }}</h2>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/30 group-hover:text-parchment-primary transition-colors">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" x2="21" y1="14" y2="3"></line>
                    </svg>
                  </router-link>
              </header>
              
              <div class="prose max-w-none mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-parchment-neutral/90" v-html="readings.second_reading.text"></div>
              </div>

              <div v-if="readings.second_reading.context" class="bg-parchment-bg border-l-4 border-parchment-primary rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-parchment-primary-dark text-[10px] font-bold uppercase tracking-wider">Contextual Meditation</h3>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-50 text-blue-900 border border-blue-200/80">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-parchment-neutral/70text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.second_reading.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.second_reading.citations" :rawContextHtml="readings.second_reading.context" />
              </div>
            </ParchmentCard>

            <!-- Responsorial Psalm Card -->
            <ParchmentCard v-if="activeTab === 'responsorial_psalm' && readings.responsorial_psalm" class="shadow-sm">
              <header class="mb-6 border-b border-parchment-border/40 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-parchment-neutral/40 font-bold tracking-[0.2em] text-[10px] uppercase block">Responsorial Psalm</span>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100/80 text-amber-900 border border-amber-300/60">📖 Sacred Scripture</span>
                  </div>
                  <div class="text-center">
                    <router-link
                      :to="{ path: '/bible-study', query: { query: readings.responsorial_psalm.citation } }"
                      class="inline-flex items-center space-x-2 group cursor-pointer"
                    >
                      <h2 class="text-xl md:text-2xl font-serif text-parchment-neutral group-hover:text-parchment-primary-dark transition-colors font-medium">{{ readings.responsorial_psalm.citation }}</h2>
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/30 group-hover:text-parchment-primary transition-colors">
                        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                        <polyline points="15 3 21 3 21 9"></polyline>
                        <line x1="10" x2="21" y1="14" y2="3"></line>
                      </svg>
                    </router-link>
                  </div>
              </header>
              
              <div class="prose max-w-none text-center mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg font-serif text-parchment-neutral/90 leading-loose" v-html="readings.responsorial_psalm.text"></div>
              </div>

              <div v-if="readings.responsorial_psalm.context" class="bg-parchment-bg border-l-4 border-parchment-primary rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-parchment-primary-dark text-[10px] font-bold uppercase tracking-wider">Psalm Context</h3>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-50 text-blue-900 border border-blue-200/80">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-parchment-neutral/70text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.responsorial_psalm.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.responsorial_psalm.citations" :rawContextHtml="readings.responsorial_psalm.context" />
              </div>
            </ParchmentCard>

            <!-- Gospel Card -->
            <ParchmentCard v-if="activeTab === 'gospel' && readings.gospel" class="shadow-sm">
              <header class="mb-6 border-b border-parchment-border/40 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-parchment-primary font-bold tracking-[0.2em] text-[10px] uppercase block">Gospel Reading</span>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100/80 text-amber-900 border border-amber-300/60">📖 Sacred Gospel</span>
                  </div>
                  <router-link
                    :to="{ path: '/bible-study', query: { query: readings.gospel.citation } }"
                    class="inline-flex items-center space-x-2 group cursor-pointer"
                  >
                    <h2 class="text-2xl md:text-3xl font-serif text-parchment-neutral group-hover:text-parchment-primary-dark transition-colors font-medium">{{ readings.gospel.citation }}</h2>
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-neutral/30 group-hover:text-parchment-primary transition-colors">
                      <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                      <polyline points="15 3 21 3 21 9"></polyline>
                      <line x1="10" x2="21" y1="14" y2="3"></line>
                    </svg>
                  </router-link>
              </header>
              
              <div class="prose max-w-none mb-8">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-parchment-neutral/95" v-html="readings.gospel.text"></div>
              </div>

              <div v-if="readings.gospel.context" class="bg-parchment-bg border-l-4 border-parchment-primary rounded-xl p-5 shadow-sm mt-6">
                  <div class="flex items-center justify-between mb-2">
                    <h3 class="text-parchment-primary-dark text-[10px] font-bold uppercase tracking-wider">Gospel Reflection</h3>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-50 text-blue-900 border border-blue-200/80">🤖 AI-Assisted Reflection</span>
                  </div>
                  <div class="text-parchment-neutral/70text-xs md:text-sm leading-relaxed" v-html="renderMarkdown(cleanContext(readings.gospel.context), contextMarkdownOptions)"></div>
                  <FactCheckCitations :citations="readings.gospel.citations" :rawContextHtml="readings.gospel.context" />
              </div>
            </ParchmentCard>

            <!-- Pope Quote/Reflection Card -->
            <ParchmentCard v-if="activeTab === 'pope_quote' && readings.pope_quote" class="shadow-sm">
              <header class="mb-6 border-b border-parchment-border/40 pb-4">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-parchment-secondary font-bold tracking-[0.2em] text-[10px] uppercase block">Words of the Popes</span>
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-purple-100/80 text-purple-900 border border-purple-300/60">🏛️ Papal Magisterium</span>
                  </div>
                  <h2 class="text-xl md:text-2xl font-serif text-parchment-neutral font-medium">Holy Father's Reflection</h2>
              </header>
              
              <div class="prose max-w-none mb-6">
                  <div class="whitespace-pre-line text-base md:text-lg leading-relaxed font-serif text-parchment-neutral/90" v-html="readings.pope_quote"></div>
              </div>

              <div class="text-center text-[10px] text-parchment-neutral/30 tracking-[0.1em] border-t border-parchment-border/20 pt-4">
                Source: Vatican News Daily Liturgy
              </div>
            </ParchmentCard>

          </div>
        </transition>

        <!-- Pastoral & Sacramental Engagement Callout (Magnifica Humanitas) -->
        <div class="bg-amber-50/70 border border-amber-200/80 rounded-2xl p-5 text-center shadow-xs mt-2">
          <div class="flex items-center justify-center space-x-2 text-amber-900 font-bold text-xs uppercase tracking-wider mb-1">
            <span>Pastoral Call to Community & Prayer</span>
          </div>
          <p class="text-amber-950/80 text-xs md:text-sm font-serif leading-relaxed">
            "Digital tools serve to support human reflection, but can never replace personal prayer, attendance at the Holy Sacrifice of the Mass, or the grace of the Sacraments in your local parish community."
          </p>
          <span class="text-[10px] text-amber-800/60 block mt-1.5 uppercase font-sans font-semibold">
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
