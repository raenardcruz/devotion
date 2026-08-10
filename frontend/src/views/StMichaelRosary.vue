<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 py-6 md:py-8 flex flex-col relative z-10">
      
      <!-- Sub-header & Controls Row -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between border-b border-[#D1C7BD]/60 pb-4 mb-4 gap-3.5 animate-fade-in-down">
          <div>
            <span class="text-xs font-bold tracking-widest text-[#72383D] uppercase block mb-1">Angelic Chaplet</span>
            <h1 class="font-serif text-2xl md:text-3xl bg-gradient-to-r from-[#72383D] via-[#AC9C8D] to-[#322D29] bg-clip-text text-transparent font-bold">
                Chaplet of St. Michael the Archangel
            </h1>
          </div>
          
          <!-- Controls -->
          <div class="flex flex-wrap items-center justify-between sm:justify-end gap-2.5 sm:gap-3.5">
              <!-- Audio Toggle & Audio Player Button -->
              <button 
                v-if="currentStep.prayerId && getPrayerAudioUrl(currentStep.prayerId)"
                @click="toggleAudio"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-gradient-to-r from-[#72383D] to-[#322D29] border border-transparent text-white rounded-full hover:from-[#8B464C] hover:to-[#453E38] active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider outline-none shadow-xs cursor-pointer"
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

              <!-- Fullscreen Card Toggle Button -->
              <button 
                @click="toggleFullscreen(cardContainerRef)"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-white/70 border border-[#D1C7BD]/80 text-[#72383D] rounded-full hover:bg-white active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider outline-none backdrop-blur-md cursor-pointer"
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
              <div class="flex items-center bg-white/70 border border-[#D1C7BD]/80 p-0.5 rounded-full backdrop-blur-md">
                  <button 
                    @click="showLatin = false"
                    :class="[!showLatin ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold' : 'text-[#322D29]/70 hover:text-[#72383D]']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none cursor-pointer"
                  >
                    EN
                  </button>
                  <button 
                    @click="showLatin = true"
                    :class="[showLatin ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold' : 'text-[#322D29]/70 hover:text-[#72383D]']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none cursor-pointer"
                  >
                    LA
                  </button>
              </div>
          </div>
      </div>

      <!-- Hero Banner & Visual Artwork -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8 items-center bg-gradient-to-r from-white/80 via-[#EFE9E1]/80 to-[#D9D9D9]/30 border border-[#D1C7BD]/80 p-5 rounded-3xl backdrop-blur-md shadow-sm">
        <div class="relative rounded-2xl overflow-hidden shadow-md aspect-square max-w-[280px] mx-auto md:max-w-none w-full">
          <transition name="fade-slide" mode="out-in">
            <img 
              :key="currentStep.choirIndex ? 'choir-' + currentStep.choirIndex : 'st-michael'"
              :src="currentStep.choirIndex ? '/nine_angelic_choirs.jpg' : '/st_michael_archangel.jpg'" 
              alt="St. Michael the Archangel" 
              class="w-full h-full object-cover transition-transform duration-700 hover:scale-105"
            />
          </transition>
          <div class="absolute inset-0 bg-gradient-to-t from-[#322D29]/70 via-transparent to-transparent flex items-end p-3">
            <span class="text-white text-xs font-semibold tracking-wide">
              {{ currentStep.choirName || 'St. Michael the Archangel' }}
            </span>
          </div>
        </div>

        <div class="md:col-span-2 flex flex-col justify-center space-y-3">
          <span class="text-[11px] font-bold tracking-widest text-[#72383D] uppercase">
            Devotion of the Nine Angelic Choirs
          </span>
          <h2 class="font-serif text-xl md:text-2xl font-bold text-[#322D29]">
            Invocation to the Archangel & Angelic Hosts
          </h2>
          <p class="text-xs md:text-sm text-[#322D29]/80 leading-relaxed">
            The Chaplet of St. Michael is a powerful prayer revealed by St. Michael the Archangel to Antonia d'Astonac. St. Michael promised that whoever honors him with these 9 salutations before Holy Communion will be escorted by an angel from each of the 9 celestial choirs.
          </p>
          <div class="flex items-center gap-2 pt-1 text-xs text-[#72383D] font-semibold">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
            <span>Promises: Escort of 9 Angels, Continual Guidance, Deliverance from Purgatory</span>
          </div>

          <!-- Official PDF Booklet Link -->
          <div class="pt-2">
            <a 
              href="https://www.dropbox.com/scl/fi/jqqsamecvha2ib5i4xz6o/4-x6-e-Booklet-The-Saint-Michael-Rosary-Prayer-Chaplet-1-.pdf.pdf?rlkey=516siw7pgeuy35vcw1rf1aky2&e=1&dl=0" 
              target="_blank" 
              rel="noopener noreferrer"
              class="inline-flex items-center gap-2 px-3.5 py-1.5 bg-[#72383D]/10 border border-[#72383D]/30 text-[#72383D] rounded-full hover:bg-[#72383D]/20 transition-all text-xs font-bold"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="12" y1="18" x2="12" y2="12"></line>
                <line x1="9" y1="15" x2="12" y2="18"></line>
                <line x1="15" y1="15" x2="12" y2="18"></line>
              </svg>
              <span>Download Printable 4"x6" e-Booklet Guide (PDF)</span>
            </a>
          </div>
        </div>
      </div>

      <!-- Main Step View & Controls -->
      <div ref="swipeContainer" class="grid grid-cols-1 lg:grid-cols-5 gap-8 items-start touch-pan-y">
        
        <!-- Left: Current Prayer Card & Controls (3/5 width) -->
        <div 
          ref="cardContainerRef"
          :class="[
            'lg:col-span-3 flex flex-col space-y-4 rounded-[2.5rem] transition-all duration-300',
            isFullscreen ? 'card-fullscreen-mode !bg-[#EFE9E1] p-6 md:p-12 overflow-y-auto justify-between border-0 shadow-none' : ''
          ]"
        >
          <transition name="fade-slide" mode="out-in">
            <div :key="currentIndex + (showLatin ? '-la' : '-en')" class="flex flex-col gap-4 w-full">
              <ParchmentCard class="!p-6 md:!p-8 relative min-h-[380px] flex flex-col justify-between shadow-lg border-[#D1C7BD]">
                <div>
                  <div class="flex items-center justify-between border-b border-[#D1C7BD]/60 pb-3 mb-4">
                    <div>
                      <span class="text-xs font-bold text-[#72383D] uppercase tracking-wider block">
                        Step {{ currentIndex + 1 }} of {{ steps.length }}
                      </span>
                      <h3 class="font-serif text-xl md:text-2xl font-bold text-[#322D29]">
                        {{ currentStep.title }}
                      </h3>
                    </div>
                    <span v-if="currentStep.beadText" class="px-3 py-1 bg-[#72383D]/10 border border-[#72383D]/30 text-[#72383D] text-[10px] font-bold rounded-full uppercase tracking-wider">
                      {{ currentStep.beadText }}
                    </span>
                  </div>

                  <!-- Prayer Subtitle / Virtue -->
                  <p v-if="currentStep.subtitle" class="text-xs font-semibold text-[#322D29]/70 italic mb-4">
                    {{ currentStep.subtitle }}
                  </p>

                  <!-- Prayer Content -->
                  <div class="text-base md:text-lg leading-relaxed text-[#322D29] whitespace-pre-line font-serif py-2">
                    {{ displayPrayerContent }}
                  </div>
                </div>

                <!-- Interactive Bead Indicator -->
                <div class="mt-4">
                  <StMichaelBeads 
                    :currentStep="currentStep" 
                    :activeBeadIndex="activeBeadIndex"
                    @selectBead="onSelectBead"
                  />
                </div>

                <!-- Card Bottom Navigation Controls -->
                <div class="flex items-center justify-between pt-6 border-t border-[#D1C7BD]/60 mt-4">
                  <button 
                    @click="prevStep" 
                    :disabled="currentIndex === 0"
                    class="px-4 py-2 rounded-xl bg-white border border-[#D1C7BD] text-[#322D29] font-semibold text-xs disabled:opacity-40 disabled:cursor-not-allowed hover:bg-[#D1C7BD]/20 active:scale-95 transition-all flex items-center gap-1.5 cursor-pointer"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
                    Previous
                  </button>

                  <div class="text-xs font-bold text-[#72383D]">
                    {{ currentIndex + 1 }} / {{ steps.length }}
                  </div>

                  <button 
                    @click="nextStep" 
                    :disabled="currentIndex === steps.length - 1"
                    class="px-4 py-2 rounded-xl bg-gradient-to-r from-[#72383D] to-[#322D29] border border-transparent text-white font-semibold text-xs disabled:opacity-40 disabled:cursor-not-allowed hover:from-[#8B464C] hover:to-[#453E38] active:scale-95 transition-all flex items-center gap-1.5 cursor-pointer shadow-sm"
                  >
                    Next
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                  </button>
                </div>
              </ParchmentCard>
            </div>
          </transition>

          <!-- Fullscreen Keyboard Helper Hint -->
          <div v-if="isFullscreen" class="text-center text-xs text-black/50 font-sans tracking-wide mt-2">
            Press <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">←</kbd> <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">→</kbd> to Navigate &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">Space</kbd> Play/Pause &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">ESC</kbd> Exit Fullscreen
          </div>
        </div>

        <!-- Right: Accordion Collapsible Sequence List (2/5 width) -->
        <div class="lg:col-span-2 space-y-3">
          <div class="bg-white/80 border border-[#D1C7BD]/80 p-4 rounded-2xl backdrop-blur-md">
            <div class="flex items-center justify-between mb-3 border-b border-[#D1C7BD]/60 pb-2">
              <h4 class="font-serif font-bold text-base text-[#322D29]">
                Chaplet Sections
              </h4>
              <span class="text-[10px] font-bold uppercase tracking-wider text-[#72383D] bg-[#72383D]/10 px-2 py-0.5 rounded-full">
                {{ steps.length }} Steps
              </span>
            </div>

            <!-- Accordion Sections -->
            <div class="space-y-2.5 max-h-[520px] overflow-y-auto pr-1">
              <!-- Section 1: Opening Prayers -->
              <div class="border border-[#D1C7BD]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('intro')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#72383D]/5 hover:bg-[#72383D]/10 text-[#72383D] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>✝</span>
                    <span>1. Opening Prayers</span>
                  </span>
                  <span class="text-[10px] font-semibold text-[#322D29]/50">
                    {{ expandedSection === 'intro' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'intro'" class="p-2 space-y-1 bg-white/40 border-t border-[#D1C7BD]/40">
                  <button 
                    v-for="step in introSteps" 
                    :key="step.id"
                    @click="selectStepByGlobalIndex(step.globalIdx)"
                    class="w-full text-left p-2 rounded-lg text-xs transition-all flex items-center justify-between cursor-pointer"
                    :class="[
                      currentIndex === step.globalIdx 
                        ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold shadow-xs' 
                        : 'hover:bg-white/90 text-[#322D29]/80'
                    ]"
                  >
                    <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                  </button>
                </div>
              </div>

              <!-- Section 2: The 9 Celestial Choirs -->
              <div class="border border-[#D1C7BD]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('choirs')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#72383D]/5 hover:bg-[#72383D]/10 text-[#72383D] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>👼</span>
                    <span>2. The 9 Angelic Choirs</span>
                  </span>
                  <span class="text-[10px] font-semibold text-[#322D29]/50">
                    {{ expandedSection === 'choirs' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'choirs'" class="p-2 space-y-1.5 bg-white/40 border-t border-[#D1C7BD]/40 max-h-[300px] overflow-y-auto">
                  <div v-for="choirGroup in choirGroups" :key="choirGroup.choirIndex" class="border border-[#D1C7BD]/50 rounded-lg p-1.5 bg-white/60">
                    <span class="text-[10px] font-bold text-[#72383D] uppercase tracking-wider block mb-1">
                      {{ choirGroup.choirName }}
                    </span>
                    <div class="space-y-0.5">
                      <button 
                        v-for="step in choirGroup.steps" 
                        :key="step.id"
                        @click="selectStepByGlobalIndex(step.globalIdx)"
                        class="w-full text-left p-1.5 rounded text-[11px] transition-all flex items-center justify-between cursor-pointer"
                        :class="[
                          currentIndex === step.globalIdx 
                            ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold shadow-xs' 
                            : 'hover:bg-[#72383D]/10 text-[#322D29]/80'
                        ]"
                      >
                        <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Section 3: Pendant Beads (4 Our Fathers) -->
              <div class="border border-[#D1C7BD]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('pendant')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#72383D]/5 hover:bg-[#72383D]/10 text-[#72383D] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>📿</span>
                    <span>3. Pendant 4 Our Fathers</span>
                  </span>
                  <span class="text-[10px] font-semibold text-[#322D29]/50">
                    {{ expandedSection === 'pendant' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'pendant'" class="p-2 space-y-1 bg-white/40 border-t border-[#D1C7BD]/40">
                  <button 
                    v-for="step in pendantSteps" 
                    :key="step.id"
                    @click="selectStepByGlobalIndex(step.globalIdx)"
                    class="w-full text-left p-2 rounded-lg text-xs transition-all flex items-center justify-between cursor-pointer"
                    :class="[
                      currentIndex === step.globalIdx 
                        ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold shadow-xs' 
                        : 'hover:bg-white/90 text-[#322D29]/80'
                    ]"
                  >
                    <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                  </button>
                </div>
              </div>

              <!-- Section 4: Concluding Prayers -->
              <div class="border border-[#D1C7BD]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('concluding')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#72383D]/5 hover:bg-[#72383D]/10 text-[#72383D] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>🛡️</span>
                    <span>4. Concluding Prayers</span>
                  </span>
                  <span class="text-[10px] font-semibold text-[#322D29]/50">
                    {{ expandedSection === 'concluding' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'concluding'" class="p-2 space-y-1 bg-white/40 border-t border-[#D1C7BD]/40">
                  <button 
                    v-for="step in concludingSteps" 
                    :key="step.id"
                    @click="selectStepByGlobalIndex(step.globalIdx)"
                    class="w-full text-left p-2 rounded-lg text-xs transition-all flex items-center justify-between cursor-pointer"
                    :class="[
                      currentIndex === step.globalIdx 
                        ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold shadow-xs' 
                        : 'hover:bg-white/90 text-[#322D29]/80'
                    ]"
                  >
                    <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                  </button>
                </div>
              </div>

            </div>
          </div>
        </div>

      </div>

    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import StMichaelBeads from '../components/stmichael/StMichaelBeads.vue';
import prayersData from '../data/prayers.json';
import { generateStMichaelSteps, type StMichaelStep } from '../data/stMichaelData';
import { useSwipe } from '../composables/useSwipe';
import { useFullscreen } from '../composables/useFullscreen';
import { getPrayerAudioUrl } from '../utils/audioHelper';

const showLatin = ref(false);
const currentIndex = ref(0);
const expandedSection = ref<'intro' | 'choirs' | 'pendant' | 'concluding'>('intro');
const steps = ref<StMichaelStep[]>(generateStMichaelSteps());

const swipeContainer = ref<HTMLElement | null>(null);
const cardContainerRef = ref<HTMLElement | null>(null);

const { isFullscreen, toggleFullscreen } = useFullscreen(cardContainerRef);

// Touch Swipe navigation
useSwipe(swipeContainer, {
  onSwipeLeft: () => nextStep(),
  onSwipeRight: () => prevStep(),
});

// Audio State
const currentAudio = ref<HTMLAudioElement | null>(null);
const isPlaying = ref(false);

const currentStep = computed<StMichaelStep>(() => {
  return steps.value[currentIndex.value]!;
});

// Active bead index within choir (0: Our Father / Salutation, 1-3: Hail Marys) or pendant (0-3: 4 Our Fathers)
const activeBeadIndex = computed(() => {
  const step = currentStep.value;
  if (!step) return 0;

  if (step.section === 'choirs') {
    if (step.id.endsWith('-hail-mary-1')) return 1;
    if (step.id.endsWith('-hail-mary-2')) return 2;
    if (step.id.endsWith('-hail-mary-3')) return 3;
    return 0;
  }

  if (step.section === 'pendant') {
    if (step.id === 'pendant-our-father-1') return 0;
    if (step.id === 'pendant-our-father-2') return 1;
    if (step.id === 'pendant-our-father-3') return 2;
    if (step.id === 'pendant-our-father-4') return 3;
  }

  return 0;
});

function stopAudio() {
  if (currentAudio.value) {
    currentAudio.value.pause();
    currentAudio.value.currentTime = 0;
    currentAudio.value = null;
  }
  isPlaying.value = false;
}

function playAudio(prayerId: string) {
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
    .catch((err) => {
      console.error('Audio playback failed:', err);
      isPlaying.value = false;
      currentAudio.value = null;
    });
}

function toggleAudio() {
  if (!currentStep.value.prayerId) return;
  if (isPlaying.value) {
    stopAudio();
  } else {
    playAudio(currentStep.value.prayerId);
  }
}

function nextStep() {
  if (currentIndex.value < steps.value.length - 1) {
    currentIndex.value++;
  }
}

function prevStep() {
  if (currentIndex.value > 0) {
    currentIndex.value--;
  }
}

function onSelectBead(beadIdx: number) {
  const step = currentStep.value;
  if (!step) return;

  if (step.section === 'choirs' && step.choirIndex) {
    const choirNum = step.choirIndex;
    let targetId = `choir-${choirNum}-our-father`;
    if (beadIdx === 1) targetId = `choir-${choirNum}-hail-mary-1`;
    if (beadIdx === 2) targetId = `choir-${choirNum}-hail-mary-2`;
    if (beadIdx === 3) targetId = `choir-${choirNum}-hail-mary-3`;

    const targetIndex = steps.value.findIndex(s => s.id === targetId);
    if (targetIndex !== -1) {
      currentIndex.value = targetIndex;
    }
    return;
  }

  if (step.section === 'pendant') {
    const targetId = `pendant-our-father-${beadIdx + 1}`;
    const targetIndex = steps.value.findIndex(s => s.id === targetId);
    if (targetIndex !== -1) {
      currentIndex.value = targetIndex;
    }
  }
}

// Keyboard Arrow Controls & Audio / Fullscreen shortcuts
function handleKeyDown(e: KeyboardEvent) {
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
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyDown);
  stopAudio();
});

// Auto-expand current section when step changes
watch(currentStep, (newStep) => {
  if (newStep && newStep.section) {
    expandedSection.value = newStep.section;
  }
  stopAudio();
}, { immediate: true });

function toggleSection(section: 'intro' | 'choirs' | 'pendant' | 'concluding') {
  expandedSection.value = expandedSection.value === section ? 'intro' : section;
}

function selectStepByGlobalIndex(index: number) {
  currentIndex.value = index;
}

const introSteps = computed(() => {
  return steps.value
    .map((step, idx) => ({ ...step, globalIdx: idx }))
    .filter(s => s.section === 'intro');
});

const pendantSteps = computed(() => {
  return steps.value
    .map((step, idx) => ({ ...step, globalIdx: idx }))
    .filter(s => s.section === 'pendant');
});

const concludingSteps = computed(() => {
  return steps.value
    .map((step, idx) => ({ ...step, globalIdx: idx }))
    .filter(s => s.section === 'concluding');
});

const choirGroups = computed(() => {
  const groups: { choirIndex: number; choirName: string; steps: Array<StMichaelStep & { globalIdx: number }> }[] = [];
  steps.value.forEach((step, idx) => {
    if (step.section === 'choirs' && step.choirIndex) {
      let group = groups.find(g => g.choirIndex === step.choirIndex);
      if (!group) {
        group = { choirIndex: step.choirIndex, choirName: step.choirName || `Choir ${step.choirIndex}`, steps: [] };
        groups.push(group);
      }
      group.steps.push({ ...step, globalIdx: idx });
    }
  });
  return groups;
});

const displayPrayerContent = computed(() => {
  const step = currentStep.value;
  if (!step) return '';
  if (step.customText && !showLatin.value) return step.customText;
  if (step.customLatin && showLatin.value) return step.customLatin;

  const found = prayersData.find((p: any) => p.id === step.prayerId);
  if (!found) return '';

  if (showLatin.value && found.latin) {
    return found.latin;
  }
  return found.default || found.latin || '';
});
</script>
