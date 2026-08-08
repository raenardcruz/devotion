<template>
  <div class="min-h-screen text-black flex flex-col pb-28 sm:pb-20 selection:bg-[#9333EA]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 py-6 md:py-8 flex flex-col relative z-10">
      
      <!-- Sub-header & Controls Row -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between border-b border-[#E9D5FF]/60 pb-4 mb-4 gap-3.5 animate-fade-in-down">
          <h1 class="font-serif text-2xl md:text-3xl text-[#7E22CE] font-bold">
              {{ displaySetName }}
          </h1>
          
          <!-- Audio Controls & Language Toggle -->
          <div class="flex flex-wrap items-center justify-between sm:justify-end gap-2.5 sm:gap-3.5">
              <!-- Auto Play Toggle Switch -->
              <label class="flex items-center space-x-1.5 text-[10px] font-bold uppercase tracking-wider text-black/60 cursor-pointer select-none">
                  <span>Auto-Play</span>
                  <input type="checkbox" v-model="autoPlay" class="sr-only peer" />
                  <div class="relative w-8 h-4.5 bg-white/60 border border-[#E9D5FF]/80 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[1px] after:start-[1px] after:bg-black/40 peer-checked:after:bg-[#9333EA] after:border-[#E9D5FF] after:border after:rounded-full after:h-3.5 after:w-3.5 after:transition-all peer-checked:bg-[#9333EA]/20"></div>
              </label>

              <!-- Manual Play/Pause Button -->
              <button 
                v-if="currentStep.prayerId && getPrayerAudioUrl(currentStep.prayerId)"
                @click="toggleAudio"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-[#9333EA] border border-[#9333EA] text-white rounded-full hover:bg-[#7E22CE] active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider outline-none shadow-xs cursor-pointer"
              >
                <svg v-if="!isPlaying" xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2">
                  <polygon points="6 3 20 12 6 21 6 3"></polygon>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" x2="18" y1="4" y2="20"></line>
                  <line x1="6" x2="6" y1="4" y2="20"></line>
                </svg>
                <span>{{ isPlaying ? 'Pause' : 'Listen' }}</span>
              </button>

              <!-- Custom Prayers Configuration Button -->
              <button 
                @click="showSettingsModal = true"
                class="flex items-center gap-1.5 px-3.5 py-1.5 bg-white/60 border border-[#E9D5FF]/60 text-[#7E22CE] rounded-full hover:bg-white/80 active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider outline-none backdrop-blur-md cursor-pointer"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.1a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
                <span>Prayers</span>
              </button>

              <!-- Fullscreen Presentation Toggle Button -->
              <button 
                @click="toggleFullscreen(cardContainerRef)"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-white/60 border border-[#E9D5FF]/60 text-[#7E22CE] rounded-full hover:bg-white/80 active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider outline-none backdrop-blur-md cursor-pointer"
                :title="isFullscreen ? 'Exit Fullscreen Card (ESC)' : 'Fullscreen Card Presentation (F)'"
              >
                <svg v-if="!isFullscreen" xmlns="http://www.w3.org/2000/svg" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="15 3 21 3 21 9"></polyline>
                  <polyline points="9 21 3 21 3 15"></polyline>
                  <line x1="21" y1="3" x2="14" y2="10"></line>
                  <line x1="3" y1="21" x2="10" y2="14"></line>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="4 14 10 14 10 20"></polyline>
                  <polyline points="20 10 14 10 14 4"></polyline>
                  <line x1="14" y1="10" x2="21" y2="3"></line>
                  <line x1="10" y1="14" x2="3" y2="21"></line>
                </svg>
                <span>{{ isFullscreen ? 'Exit Fullscreen' : 'Fullscreen Card' }}</span>
              </button>

              <!-- Latin/English Language Toggle -->
              <div class="flex items-center bg-white/60 border border-[#E9D5FF]/60 p-0.5 rounded-full backdrop-blur-md">
                  <button 
                    @click="showLatin = false"
                    :class="[!showLatin ? 'bg-[#9333EA] text-white font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0 cursor-pointer"
                  >
                    EN
                  </button>
                  <button 
                    @click="showLatin = true"
                    :class="[showLatin ? 'bg-[#9333EA] text-white font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0 cursor-pointer"
                  >
                    LA
                  </button>
              </div>
          </div>
      </div>

      <!-- Segmented Mystery Selection Tabs -->
      <AppTabs :tabs="sets" v-model="currentSetName" class="animate-fade-in-down" />

      <!-- Interactive 2-Column Prayer Grid -->
      <div ref="swipeContainer" class="grid grid-cols-1 lg:grid-cols-5 gap-8 mt-6 items-start touch-pan-y animate-fade-in-up delay-150">
        
        <!-- Left Column: Active Card & Controls (3/5 width on desktop) -->
        <div 
            ref="cardContainerRef" 
            :class="[
                'lg:col-span-3 flex flex-col gap-6 rounded-[2.5rem] transition-all duration-300',
                isFullscreen ? 'card-fullscreen-mode !bg-[#FAF8F5] p-6 md:p-12 overflow-y-auto justify-between border-0 shadow-none' : ''
            ]"
        >
            <transition name="fade-slide" mode="out-in">
                <div :key="currentSetName + currentStepIndex + (showLatin ? '-la' : '-en')" class="flex flex-col gap-6 w-full">
                    
                    <!-- Verse Card -->
                    <RosaryCard 
                        :currentStep="currentStep" 
                        :showLatin="showLatin"
                        :mysteryTitle="currentStep.mysteryTitle || getCurrentMystery(currentStep).title"
                        :isFullscreen="isFullscreen"
                        @restart="currentStepIndex = 0"
                        @video-active="stopAudio"
                        @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
                    />

                    <!-- Interactive Bead Tracker -->
                    <RosaryBeads 
                        :beadIndex="currentStep.beadNumber || 0" 
                        :stepType="currentStep.type" 
                        @selectBead="onSelectBead"
                    />
                </div>
            </transition>

            <!-- Navigation Controls -->
            <RosaryControls 
                :totalSteps="steps.length"
                :currentStepIndex="currentStepIndex"
                :isLastStep="currentStepIndex === steps.length - 1"
                @next="nextStep" 
                @prev="prevStep" 
            />

            <!-- Fullscreen Keyboard Helper Hint -->
            <div v-if="isFullscreen" class="text-center text-xs text-black/50 font-sans tracking-wide mt-2">
                Press <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">←</kbd> <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">→</kbd> to Navigate &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">Space</kbd> Play/Pause &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">ESC</kbd> Exit Fullscreen
            </div>
        </div>

        <!-- Right Column: Classical Illustration & Meditative Quote (2/5 width on desktop) -->
        <div class="lg:col-span-2 flex flex-col gap-6">
            <!-- Art Image Container -->
            <div class="w-full h-72 md:h-80 rounded-[2rem] overflow-hidden border border-parchment-border shadow-sm bg-parchment-neutral-light flex items-center justify-center relative group">
                <transition name="fade-slide" mode="out-in">
                  <img 
                      :key="currentVisual.image"
                      :src="currentVisual.image" 
                      :alt="currentStep.mysteryTitle || currentSetName" 
                      class="w-full h-full object-cover transition-all duration-1000 group-hover:scale-105"
                  />
                </transition>
            </div>

            <!-- Meditation Quote Card -->
            <div class="bg-parchment-neutral-light/50 border border-parchment-border border-l-4 border-l-parchment-primary rounded-3xl p-6 shadow-sm">
                <transition name="fade-slide" mode="out-in">
                  <p :key="currentVisual.meditation" class="font-seriftext-sm md:text-base text-parchment-neutral/80 leading-relaxed">
                      {{ currentVisual.meditation }}
                  </p>
                </transition>
            </div>
        </div>

      </div>
    </main>

    <!-- Global Footer -->
    <BottomNav />

    <!-- Custom Prayers Selection Modal -->
    <transition name="fade">
      <div v-if="showSettingsModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-[100] flex items-center justify-center p-4" @click.self="showSettingsModal = false">
        <div class="bg-white/95 backdrop-blur-md border border-[#E9D5FF] w-full max-w-md rounded-[2rem] shadow-2xl flex flex-col max-h-[80vh] overflow-hidden transition-all duration-300 transform scale-100">
          
          <!-- Modal Header -->
          <div class="p-5 border-b border-[#E9D5FF]/60 flex items-center justify-between bg-white/60">
            <div>
              <h3 class="font-serif text-base md:text-lg text-[#7E22CE] font-bold">Custom Prayers</h3>
              <p class="text-[9px] text-[#D97706] font-bold uppercase tracking-wider mt-0.5">Customize your Rosary sequence</p>
            </div>
            <button @click="showSettingsModal = false" class="w-7 h-7 rounded-full border border-[#E9D5FF] flex items-center justify-center text-black/60 hover:bg-[#E9D5FF]/30 hover:text-[#7E22CE] active:scale-95 transition-all p-0 cursor-pointer">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <!-- Search Bar -->
          <div class="px-5 py-3 border-b border-[#E9D5FF]/40 bg-white/40">
            <div class="relative">
              <input 
                type="text" 
                v-model="searchQuery" 
                placeholder="Search prayers by name..." 
                class="w-full px-4 py-2 pl-9 pr-8 bg-white/80 border border-[#E9D5FF] rounded-full text-xs text-black placeholder-black/40 focus:outline-none focus:ring-1 focus:ring-[#9333EA]/30 focus:border-[#9333EA] transition-all font-sans"
              />
              <div class="absolute left-3 top-1/2 -translate-y-1/2 text-[#9333EA] pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                </svg>
              </div>
              <div 
                v-if="searchQuery" 
                @click="searchQuery = ''" 
                class="absolute right-3 top-1/2 -translate-y-1/2 text-black/40 hover:text-[#7E22CE] cursor-pointer transition-all flex items-center justify-center p-0.5 rounded-full hover:bg-black/5"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" x2="6" y1="6" y2="18"></line>
                  <line x1="6" x2="18" y1="6" y2="18"></line>
                </svg>
              </div>
            </div>
          </div>

          <!-- Modal Content (Scrollable) -->
          <div class="p-5 overflow-y-auto space-y-5 flex-grow">
            
            <div v-if="filteredSelectablePrayers.length === 0" class="text-center py-8 flex flex-col items-center justify-center space-y-2 animate-fade-in-up">
              <div class="w-10 h-10 rounded-full bg-[#E9D5FF]/40 flex items-center justify-center text-[#9333EA]">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                </svg>
              </div>
              <div class="text-xs font-serif text-[#7E22CE] font-bold">No prayers found</div>
              <p class="text-[10px] text-black/50 max-w-[200px] text-center leading-normal">
                No results match "{{ searchQuery }}". Try searching for another keyword.
              </p>
            </div>

            <template v-else>
              <!-- Prepend Section -->
              <div>
                <h4 class="font-serif text-xs font-bold text-[#D97706] border-b border-[#E9D5FF]/60 pb-1.5 mb-2">
                  Before the Rosary (Prepend)
                </h4>
                <p class="text-[11px] text-black/60 mb-2 leading-relaxed">
                  Recited immediately after the opening Sign of the Cross.
                </p>
                <div class="space-y-1.5 max-h-36 overflow-y-auto pr-1 border border-[#E9D5FF]/60 rounded-xl p-2 bg-white/50">
                  <label 
                    v-for="prayer in filteredSelectablePrayers" 
                    :key="'before-' + prayer.id" 
                    class="flex items-start space-x-2.5 p-1.5 rounded-lg hover:bg-[#E9D5FF]/20 transition-colors cursor-pointer select-none"
                  >
                    <input 
                      type="checkbox" 
                      :value="prayer.id" 
                      v-model="tempBeforePrayers" 
                      class="mt-0.5 rounded border-[#E9D5FF] text-[#9333EA] focus:ring-[#9333EA]/30 h-3.5 w-3.5 custom-checkbox cursor-pointer" 
                    />
                    <div class="flex-grow min-w-0">
                      <div class="text-xs font-bold text-black truncate">{{ prayer.name }}</div>
                      <div v-if="prayer.latinName" class="text-[9px] text-black/50 truncate">{{ prayer.latinName }}</div>
                    </div>
                  </label>
                </div>
              </div>

              <!-- Append Section -->
              <div>
                <h4 class="font-serif text-xs font-bold text-[#D97706] border-b border-[#E9D5FF]/60 pb-1.5 mb-2">
                  After the Rosary (Append)
                </h4>
                <p class="text-[11px] text-black/60 mb-2 leading-relaxed">
                  Recited after traditional closing prayers, before the final Sign of the Cross.
                </p>
                <div class="space-y-1.5 max-h-36 overflow-y-auto pr-1 border border-[#E9D5FF]/60 rounded-xl p-2 bg-white/50">
                  <label 
                    v-for="prayer in filteredSelectablePrayers" 
                    :key="'after-' + prayer.id" 
                    class="flex items-start space-x-2.5 p-1.5 rounded-lg hover:bg-[#E9D5FF]/20 transition-colors cursor-pointer select-none"
                  >
                    <input 
                      type="checkbox" 
                      :value="prayer.id" 
                      v-model="tempAfterPrayers" 
                      class="mt-0.5 rounded border-[#E9D5FF] text-[#9333EA] focus:ring-[#9333EA]/30 h-3.5 w-3.5 custom-checkbox cursor-pointer" 
                    />
                    <div class="flex-grow min-w-0">
                      <div class="text-xs font-bold text-black truncate">{{ prayer.name }}</div>
                      <div v-if="prayer.latinName" class="text-[9px] text-black/50 truncate">{{ prayer.latinName }}</div>
                    </div>
                  </label>
                </div>
              </div>
            </template>

            <!-- Notice about restarting -->
            <div class="bg-[#FEF3C7] border border-[#FDE68A] rounded-xl p-2.5 flex items-start space-x-2">
              <svg class="text-[#D97706] mt-0.5 shrink-0" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="16" x2="12" y2="12"></line>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
              </svg>
              <p class="text-[9px] text-[#D97706] font-medium leading-normal">
                Applying custom prayers will update the Rosary sequence and restart your current meditation from the beginning.
              </p>
            </div>

          </div>

          <!-- Modal Footer -->
          <div class="p-4 bg-white/80 border-t border-[#E9D5FF]/60 flex items-center justify-end space-x-2">
            <button 
              @click="showSettingsModal = false" 
              class="px-4 py-1.5 border border-[#E9D5FF] hover:bg-white text-black/70 hover:text-[#7E22CE] rounded-full text-[10px] font-bold uppercase tracking-wider outline-none shadow-none cursor-pointer"
            >
              Cancel
            </button>
            <button 
              @click="applySettings" 
              class="px-5 py-1.5 bg-[#9333EA] hover:bg-[#7E22CE] text-white rounded-full text-[10px] font-bold uppercase tracking-wider outline-none shadow-xs transition-all cursor-pointer"
            >
              Apply Prayers
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useSwipe } from '../composables/useSwipe';
import { useDate } from '../composables/useDate';
import { useFullscreen } from '../composables/useFullscreen';
import { ROSARY_DATA, generateRosarySteps, type Mystery, type RosaryStep } from '../components/rosary/rosaryData';
import { getMysteryVisual } from '../utils/mysteryVisuals';
import { getPrayerAudioUrl } from '../utils/audioHelper';
import prayerData from '../data/prayers.json';

import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import AppTabs from '../components/common/AppTabs.vue';
import RosaryBeads from '../components/rosary/RosaryBeads.vue';
import RosaryCard from '../components/rosary/VerseCard.vue';
import RosaryControls from '../components/rosary/RosaryControls.vue';

const currentSetName = ref('Joyful');
const currentStepIndex = ref(0);
const showLatin = ref(false);
const todayMystery = ref('');
const sets = Object.keys(ROSARY_DATA);
const swipeContainer = ref<HTMLElement | null>(null);
const cardContainerRef = ref<HTMLElement | null>(null);

const { isFullscreen, toggleFullscreen } = useFullscreen(cardContainerRef);

// Custom prayers configuration state
const beforePrayers = ref<string[]>([]);
const afterPrayers = ref<string[]>([]);
const showSettingsModal = ref(false);
const tempBeforePrayers = ref<string[]>([]);
const tempAfterPrayers = ref<string[]>([]);
const searchQuery = ref('');

const mysterySetTitleLatinMap: Record<string, string> = {
  'Joyful': 'Mysteria Gaudiosa',
  'Luminous': 'Mysteria Luminosa',
  'Sorrowful': 'Mysteria Dolorosa',
  'Glorious': 'Mysteria Gloriosa'
};

const displaySetName = computed(() => {
  if (showLatin.value) {
    return mysterySetTitleLatinMap[currentSetName.value] || currentSetName.value;
  }
  return `The ${currentSetName.value} Mysteries`;
});

// Audio State
const currentAudio = ref<HTMLAudioElement | null>(null);
const isPlaying = ref(false);
const autoPlay = ref(true);

const { getLocalDay } = useDate();

useSwipe(swipeContainer, {
    onSwipeLeft: () => nextStep(),
    onSwipeRight: () => prevStep(),
});

const getMysteryForDay = (): string => {
    const day = getLocalDay();
    switch (day) {
        case 0: return 'Glorious';
        case 1: return 'Joyful';
        case 2: return 'Sorrowful';
        case 3: return 'Glorious';
        case 4: return 'Luminous';
        case 5: return 'Sorrowful';
        case 6: return 'Joyful';
        default: return 'Joyful';
    }
};

const handleKeyDown = (e: KeyboardEvent) => {
    if (['INPUT', 'TEXTAREA', 'SELECT'].includes((e.target as HTMLElement)?.tagName)) {
        return;
    }
    if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
        e.preventDefault();
        nextStep();
    } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
        e.preventDefault();
        prevStep();
    } else if (e.key === ' ') {
        e.preventDefault();
        toggleAudio();
    } else if (e.key === 'f' || e.key === 'F') {
        e.preventDefault();
        toggleFullscreen(cardContainerRef.value);
    }
};

onMounted(() => {
    window.addEventListener('keydown', handleKeyDown);
    const daily = getMysteryForDay();
    todayMystery.value = daily;
    currentSetName.value = daily;

    // Load Autoplay Preference
    const savedAutoPlay = localStorage.getItem('rosary-autoplay');
    if (savedAutoPlay !== null) {
        autoPlay.value = savedAutoPlay === 'true';
    }

    // Load custom prayers preference
    const savedBefore = localStorage.getItem('rosary-before-prayers');
    if (savedBefore) {
        try {
            beforePrayers.value = JSON.parse(savedBefore);
        } catch (e) {
            console.error('Failed to parse saved before prayers', e);
        }
    }
    const savedAfter = localStorage.getItem('rosary-after-prayers');
    if (savedAfter) {
        try {
            afterPrayers.value = JSON.parse(savedAfter);
        } catch (e) {
            console.error('Failed to parse saved after prayers', e);
        }
    }
});

const steps = computed(() => generateRosarySteps(currentSetName.value, beforePrayers.value, afterPrayers.value));

const selectablePrayers = computed(() => {
    // Filter out Divine Mercy specific prayers to keep selection options clean and focused.
    const dmIds = ['dm-opening-1', 'dm-opening-2', 'eternal-father', 'sorrowful-passion', 'holy-god', 'dm-closing', 'jesus-i-trust-in-you'];
    return prayerData.filter(p => !dmIds.includes(p.id) && p.id !== 'sign-of-the-cross');
});

const filteredSelectablePrayers = computed(() => {
    const query = searchQuery.value.trim().toLowerCase();
    if (!query) return selectablePrayers.value;
    return selectablePrayers.value.filter(prayer => {
        return (
            prayer.name.toLowerCase().includes(query) ||
            (prayer.latinName && prayer.latinName.toLowerCase().includes(query))
        );
    });
});

watch(showSettingsModal, (isOpen) => {
    if (isOpen) {
        tempBeforePrayers.value = [...beforePrayers.value];
        tempAfterPrayers.value = [...afterPrayers.value];
        searchQuery.value = '';
    }
});

const applySettings = () => {
    beforePrayers.value = [...tempBeforePrayers.value];
    afterPrayers.value = [...tempAfterPrayers.value];
    
    localStorage.setItem('rosary-before-prayers', JSON.stringify(beforePrayers.value));
    localStorage.setItem('rosary-after-prayers', JSON.stringify(afterPrayers.value));
    
    currentStepIndex.value = 0;
    showSettingsModal.value = false;
};

const currentStep = computed((): RosaryStep => {
    const step = steps.value[currentStepIndex.value];
    if (!step) return steps.value[0]!;

    if (step.prayerId) {
        const prayer = prayerData.find(p => p.id === step.prayerId);
        if (prayer) {
            return {
                ...step,
                title: step.title || prayer.name,
                content: step.content || prayer.default,
                latin: step.latin || prayer.latin,
                youtube: (prayer as any).youtube
            };
        }
    }
    
    return step;
});

const getCurrentMystery = (step: RosaryStep): Mystery => {
    const mysteries = ROSARY_DATA[currentSetName.value];
    if (step.decadeNumber && mysteries && mysteries[step.decadeNumber - 1]) {
        return mysteries[step.decadeNumber - 1]!;
    }
    return mysteries![0]!;
};

const currentVisual = computed(() => {
    const mystery = getCurrentMystery(currentStep.value);
    const title = currentStep.value.mysteryTitle || mystery.title;
    if (mystery && mystery.image) {
        return {
            image: mystery.image,
            meditation: getMysteryVisual(title).meditation
        };
    }
    return getMysteryVisual(title);
});

const nextStep = () => {
    if (currentStepIndex.value < steps.value.length - 1) {
        currentStepIndex.value++;
    } else {
        currentStepIndex.value = 0;
    }
};

const prevStep = () => {
    if (currentStepIndex.value > 0) {
        currentStepIndex.value--;
    }
};

const onSelectBead = (index: number) => {
    if (currentStep.value.type === 'intro') {
        const targetIdx = steps.value.findIndex(
            s => s.type === 'intro' && s.beadNumber === index
        );
        if (targetIdx !== -1) {
            currentStepIndex.value = targetIdx;
        }
        return;
    }

    const decNum = currentStep.value.decadeNumber;
    if (!decNum) return;
    const targetIdx = steps.value.findIndex(
        s => s.type === 'decade-bead' && s.decadeNumber === decNum && s.beadNumber === index
    );
    if (targetIdx !== -1) {
        currentStepIndex.value = targetIdx;
    }
};

// Reset index when set changes
watch(currentSetName, () => {
    currentStepIndex.value = 0;
});

// Audio Control Functions
let playTimeout: any = null;

const playAudio = (prayerId: string | undefined) => {
    if (!prayerId) return;
    const url = getPrayerAudioUrl(prayerId);
    if (!url) return;
    
    stopAudio();
    
    const audio = new Audio(url);
    currentAudio.value = audio;
    
    audio.addEventListener('ended', () => {
        isPlaying.value = false;
        currentAudio.value = null;
    });
    
    audio.addEventListener('error', () => {
        isPlaying.value = false;
        currentAudio.value = null;
    });

    audio.play()
      .then(() => {
          isPlaying.value = true;
      })
      .catch(err => {
          // Autoplay blocked by browser policy
          console.log('Audio playback postponed until user interaction:', err);
          isPlaying.value = false;
      });
};

const stopAudio = () => {
    if (playTimeout) {
        clearTimeout(playTimeout);
        playTimeout = null;
    }
    if (currentAudio.value) {
        currentAudio.value.pause();
        currentAudio.value = null;
    }
    isPlaying.value = false;
};

const pauseAudio = () => {
    if (playTimeout) {
        clearTimeout(playTimeout);
        playTimeout = null;
    }
    if (currentAudio.value) {
        currentAudio.value.pause();
    }
    isPlaying.value = false;
};

const toggleAudio = () => {
    if (isPlaying.value) {
        pauseAudio();
    } else if (currentAudio.value) {
        currentAudio.value.play()
            .then(() => {
                isPlaying.value = true;
            })
            .catch(err => {
                console.log('Audio playback postponed:', err);
                isPlaying.value = false;
            });
    } else if (currentStep.value.prayerId) {
        playAudio(currentStep.value.prayerId);
    }
};

// Monitor currentStep change to trigger audio auto-play
watch(currentStep, (newStep) => {
    stopAudio();
    if (autoPlay.value && newStep.prayerId) {
        // Delay slightly to give page transition room to render nicely
        playTimeout = setTimeout(() => {
            if (newStep.prayerId === currentStep.value.prayerId) {
                playAudio(newStep.prayerId);
            }
        }, 300);
    }
}, { deep: true });

// Sync Autoplay Preference
watch(autoPlay, (newVal) => {
    localStorage.setItem('rosary-autoplay', String(newVal));
    if (!newVal) {
        stopAudio();
    } else if (currentStep.value.prayerId) {
        playAudio(currentStep.value.prayerId);
    }
});

onBeforeUnmount(() => {
    window.removeEventListener('keydown', handleKeyDown);
    stopAudio();
});
</script>

<style scoped>
.custom-checkbox {
  accent-color: var(--parchment-primary);
}
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

:fullscreen, :-webkit-full-screen {
  background-color: #FAF8F5 !important;
  width: 100vw !important;
  height: 100vh !important;
  max-width: none !important;
  max-height: none !important;
  padding: 3rem 2rem !important;
  display: flex !important;
  flex-direction: column !important;
  justify-content: space-between !important;
  align-items: center !important;
  box-sizing: border-box !important;
  overflow-y: auto !important;
}

:fullscreen :deep(.max-w-4xl),
:-webkit-full-screen :deep(.max-w-4xl) {
  max-width: 56rem !important;
  width: 100% !important;
  margin-left: auto !important;
  margin-right: auto !important;
}
</style>