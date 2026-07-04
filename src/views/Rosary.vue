<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-24 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 py-8 flex flex-col">
      
      <!-- Sub-header & Controls Row -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between border-b border-parchment-border pb-4 mb-4 gap-4 animate-fade-in-down">
          <h1 class="font-serif text-2xl md:text-3xl text-parchment-neutral font-medium">
              {{ displaySetName }}
          </h1>
          
          <!-- Audio Controls & Language Toggle -->
          <div class="flex flex-wrap items-center gap-3 sm:gap-4">
              <!-- Auto Play Toggle Switch -->
              <label class="flex items-center space-x-2 text-[10px] font-bold uppercase tracking-wider text-parchment-neutral/50 cursor-pointer select-none">
                  <span>Auto-Play</span>
                  <input type="checkbox" v-model="autoPlay" class="sr-only peer" />
                  <div class="relative w-8 h-4.5 bg-parchment-neutral-light border border-parchment-border rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[1px] after:start-[1px] after:bg-parchment-neutral/40 peer-checked:after:bg-parchment-primary after:border-parchment-border after:border after:rounded-full after:h-3.5 after:w-3.5 after:transition-all peer-checked:bg-parchment-primary/10"></div>
              </label>

              <!-- Manual Play/Pause Button -->
              <button 
                v-if="currentStep.prayerId && getPrayerAudioUrl(currentStep.prayerId)"
                @click="toggleAudio"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-parchment-primary/10 border border-parchment-primary/30 text-parchment-primary-dark rounded-full hover:bg-parchment-primary/20 active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider outline-none shadow-none"
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
                class="flex items-center gap-1.5 px-3 py-1.5 bg-parchment-neutral-light border border-parchment-border text-parchment-neutral rounded-full hover:bg-parchment-primary/10 active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider outline-none shadow-none"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.1a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path>
                  <circle cx="12" cy="12" r="3"></circle>
                </svg>
                <span>Prayers</span>
              </button>

              <!-- Latin/English Language Toggle -->
              <div class="flex items-center bg-parchment-neutral-light border border-parchment-border p-0.5 rounded-full">
                  <button 
                    @click="showLatin = false"
                    :class="[!showLatin ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                  >
                    EN
                  </button>
                  <button 
                    @click="showLatin = true"
                    :class="[showLatin ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
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
        <div class="lg:col-span-3 flex flex-col gap-6">
            <transition name="fade-slide" mode="out-in">
                <div :key="currentSetName + currentStepIndex + (showLatin ? '-la' : '-en')" class="flex flex-col gap-6 w-full">
                    
                    <!-- Verse Card -->
                    <RosaryCard 
                        :currentStep="currentStep" 
                        :showLatin="showLatin"
                        :mysteryTitle="currentStep.mysteryTitle || getCurrentMystery(currentStep).title"
                        @restart="currentStepIndex = 0"
                        @video-active="stopAudio"
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
                  <p :key="currentVisual.meditation" class="font-serif italic text-sm md:text-base text-parchment-neutral/80 leading-relaxed">
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
      <div v-if="showSettingsModal" class="fixed inset-0 bg-parchment-neutral/65 backdrop-blur-sm z-[100] flex items-center justify-center p-4" @click.self="showSettingsModal = false">
        <div class="bg-[#FDFBF7] border border-parchment-border w-full max-w-md rounded-[2rem] shadow-2xl flex flex-col max-h-[80vh] overflow-hidden transition-all duration-300 transform scale-100">
          
          <!-- Modal Header -->
          <div class="p-5 border-b border-parchment-border/40 flex items-center justify-between bg-parchment-neutral-light/20">
            <div>
              <h3 class="font-serif text-base md:text-lg text-parchment-primary-dark font-medium">Custom Prayers</h3>
              <p class="text-[9px] text-parchment-neutral/50 font-bold uppercase tracking-wider mt-0.5">Customize your Rosary sequence</p>
            </div>
            <button @click="showSettingsModal = false" class="w-7 h-7 rounded-full border border-parchment-border flex items-center justify-center text-parchment-neutral hover:bg-parchment-neutral-light hover:text-parchment-primary active:scale-95 transition-all p-0">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <!-- Search Bar -->
          <div class="px-5 py-3 border-b border-parchment-border/20 bg-parchment-neutral-light/5">
            <div class="relative">
              <input 
                type="text" 
                v-model="searchQuery" 
                placeholder="Search prayers by name..." 
                class="w-full px-4 py-2 pl-9 pr-8 bg-parchment-neutral-light/50 border border-parchment-border rounded-full text-xs text-parchment-neutral placeholder-parchment-neutral/40 focus:outline-none focus:ring-1 focus:ring-parchment-primary/50 focus:border-parchment-primary transition-all font-sans"
              />
              <div class="absolute left-3 top-1/2 -translate-y-1/2 text-parchment-neutral/40 pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                </svg>
              </div>
              <div 
                v-if="searchQuery" 
                @click="searchQuery = ''" 
                class="absolute right-3 top-1/2 -translate-y-1/2 text-parchment-neutral/40 hover:text-parchment-neutral cursor-pointer transition-all flex items-center justify-center p-0.5 rounded-full hover:bg-parchment-neutral/5"
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
              <div class="w-10 h-10 rounded-full bg-parchment-primary/10 flex items-center justify-center text-parchment-primary">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                </svg>
              </div>
              <div class="text-xs font-serif text-parchment-primary-dark font-medium">No prayers found</div>
              <p class="text-[10px] text-parchment-neutral/50 max-w-[200px] text-center leading-normal">
                No results match "{{ searchQuery }}". Try searching for another keyword.
              </p>
            </div>

            <template v-else>
              <!-- Prepend Section -->
              <div>
                <h4 class="font-serif text-xs font-semibold text-parchment-secondary border-b border-parchment-border/30 pb-1.5 mb-2">
                  Before the Rosary (Prepend)
                </h4>
                <p class="text-[11px] text-parchment-neutral/60 mb-2 leading-relaxed">
                  Recited immediately after the opening Sign of the Cross.
                </p>
                <div class="space-y-1.5 max-h-36 overflow-y-auto pr-1 border border-parchment-border/20 rounded-xl p-2 bg-parchment-neutral-light/10">
                  <label 
                    v-for="prayer in filteredSelectablePrayers" 
                    :key="'before-' + prayer.id" 
                    class="flex items-start space-x-2.5 p-1.5 rounded-lg hover:bg-parchment-neutral-light/50 transition-colors cursor-pointer select-none"
                  >
                    <input 
                      type="checkbox" 
                      :value="prayer.id" 
                      v-model="tempBeforePrayers" 
                      class="mt-0.5 rounded border-parchment-border text-parchment-primary focus:ring-parchment-primary/30 h-3.5 w-3.5 custom-checkbox" 
                    />
                    <div class="flex-grow min-w-0">
                      <div class="text-xs font-bold text-parchment-neutral truncate">{{ prayer.name }}</div>
                      <div v-if="prayer.latinName" class="text-[9px] text-parchment-neutral/40 italic truncate">{{ prayer.latinName }}</div>
                    </div>
                  </label>
                </div>
              </div>

              <!-- Append Section -->
              <div>
                <h4 class="font-serif text-xs font-semibold text-parchment-secondary border-b border-parchment-border/30 pb-1.5 mb-2">
                  After the Rosary (Append)
                </h4>
                <p class="text-[11px] text-parchment-neutral/60 mb-2 leading-relaxed">
                  Recited after traditional closing prayers, before the final Sign of the Cross.
                </p>
                <div class="space-y-1.5 max-h-36 overflow-y-auto pr-1 border border-parchment-border/20 rounded-xl p-2 bg-parchment-neutral-light/10">
                  <label 
                    v-for="prayer in filteredSelectablePrayers" 
                    :key="'after-' + prayer.id" 
                    class="flex items-start space-x-2.5 p-1.5 rounded-lg hover:bg-parchment-neutral-light/50 transition-colors cursor-pointer select-none"
                  >
                    <input 
                      type="checkbox" 
                      :value="prayer.id" 
                      v-model="tempAfterPrayers" 
                      class="mt-0.5 rounded border-parchment-border text-parchment-primary focus:ring-parchment-primary/30 h-3.5 w-3.5 custom-checkbox" 
                    />
                    <div class="flex-grow min-w-0">
                      <div class="text-xs font-bold text-parchment-neutral truncate">{{ prayer.name }}</div>
                      <div v-if="prayer.latinName" class="text-[9px] text-parchment-neutral/40 italic truncate">{{ prayer.latinName }}</div>
                    </div>
                  </label>
                </div>
              </div>
            </template>

            <!-- Notice about restarting -->
            <div class="bg-parchment-primary/5 border border-parchment-primary/20 rounded-xl p-2.5 flex items-start space-x-2">
              <svg class="text-parchment-primary-dark mt-0.5 shrink-0" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="16" x2="12" y2="12"></line>
                <line x1="12" y1="8" x2="12.01" y2="8"></line>
              </svg>
              <p class="text-[9px] text-parchment-neutral/70 leading-normal">
                Applying custom prayers will update the Rosary sequence and restart your current meditation from the beginning.
              </p>
            </div>

          </div>

          <!-- Modal Footer -->
          <div class="p-4 bg-parchment-neutral-light/30 border-t border-parchment-border/40 flex items-center justify-end space-x-2">
            <button 
              @click="showSettingsModal = false" 
              class="px-4 py-1.5 border border-parchment-border hover:bg-parchment-neutral-light text-parchment-neutral rounded-full text-[10px] font-bold uppercase tracking-wider outline-none shadow-none cursor-pointer"
            >
              Cancel
            </button>
            <button 
              @click="applySettings" 
              class="px-4 py-1.5 bg-parchment-primary hover:bg-parchment-primary-dark text-white rounded-full text-[10px] font-bold uppercase tracking-wider border border-transparent shadow-sm hover:scale-[1.02] cursor-pointer"
            >
              Apply &amp; Restart
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

onMounted(() => {
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
    const title = currentStep.value.mysteryTitle || getCurrentMystery(currentStep.value).title;
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
</style>