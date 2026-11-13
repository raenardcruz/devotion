<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 py-6 md:py-8 flex flex-col relative z-10">
      
      <!-- Clean Unified Devotion Header with Collapsible Settings Menu -->
      <DevotionHeader
        title="Chaplet of St. Michael the Archangel"
        subtitle="Nine Angelic Choirs"
        v-model:showLatin="showLatin"
        :supportsAutoPlay="false"
        :activeIntentionsCount="activeCount"
        :isFullscreen="isFullscreen"
        :isFullscreenSupported="isFullscreenSupported"
        @open-intentions="showIntentionsModal = true"
        @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
      />

      <!-- Main Step View & Controls -->
      <div ref="swipeContainer" class="grid grid-cols-1 lg:grid-cols-5 gap-8 items-start touch-pan-y animate-fade-in-up delay-150">
        
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
              <!-- Prayer Intentions Card -->
              <PrayerIntentionsCard 
                v-if="currentStep.id === 'prayer-intentions'"
                :intentions="activeIntentions"
                :showLatin="showLatin"
                :isFullscreen="isFullscreen"
                @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
                @open-intentions-modal="showIntentionsModal = true"
              />

              <!-- Standard Prayer Card (Unified Gradient & Fade Surface) -->
              <div
                v-else
                :class="[
                  isFullscreen
                    ? 'devotion-card-surface p-8 md:p-12 rounded-[2.5rem] shadow-xl flex flex-col justify-between w-full my-auto relative overflow-hidden'
                    : 'devotion-card-surface p-6 md:p-8 rounded-[2rem] min-h-[380px] flex flex-col justify-between relative overflow-hidden'
                ]"
              >
                <!-- Subtle Top Atmospheric Glow -->
                <div class="pointer-events-none absolute -top-24 left-1/2 -translate-x-1/2 w-3/4 h-40 rounded-full bg-radial from-[#72383D]/10 via-[#AC9C8D]/8 to-transparent blur-2xl" />

                <div class="relative z-10">
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
                <div class="mt-4 relative z-10">
                  <StMichaelBeads 
                    :currentStep="currentStep" 
                    :activeBeadIndex="activeBeadIndex"
                    @selectBead="onSelectBead"
                  />
                </div>
              </div>

              <!-- Card Bottom Navigation & Audio Controls -->
              <DevotionControls
                :currentStepIndex="currentIndex"
                :totalSteps="steps.length"
                :isLastStep="currentIndex === steps.length - 1"
                :disableNextWhenLast="true"
                :hasAudio="Boolean(currentStep.prayerId && getPrayerAudioUrl(currentStep.prayerId))"
                :isPlaying="isPlaying"
                label="Chaplet"
                @next="nextStep"
                @prev="prevStep"
                @toggle-audio="toggleAudio"
              />
            </div>
          </transition>

          <!-- Fullscreen Keyboard Helper Hint -->
          <div v-if="isFullscreen" class="text-center text-xs text-black/50 font-sans tracking-wide mt-2">
            Press <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">←</kbd> <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">→</kbd> to Navigate &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">Space</kbd> Play/Pause &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">ESC</kbd> Exit Fullscreen
          </div>
        </div>

        <!-- Right: Seamless Gradient-Fade Artwork & Chaplet Sections (2/5 width) -->
        <div class="lg:col-span-2 space-y-5">
          <!-- Unified Artwork & Promises Card -->
          <div class="devotion-card-surface rounded-[2rem] overflow-hidden flex flex-col relative group">
            <div class="aspect-[4/3] w-full overflow-hidden relative bg-[#221D1A]">
              <transition name="fade-slide" mode="out-in">
                <img
                  :key="currentStep.choirIndex ? 'choir-' + currentStep.choirIndex : 'st-michael'"
                  :src="currentStep.choirIndex ? '/nine_angelic_choirs.jpg' : '/st_michael_archangel.jpg'"
                  alt="St. Michael the Archangel"
                  class="w-full h-full object-cover object-top devotion-artwork-fade transition-transform duration-1000 group-hover:scale-105"
                />
              </transition>
              <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-[#EFE9E1] via-[#EFE9E1]/35 to-transparent" />
            </div>

            <div class="relative z-10 -mt-10 px-6 pb-5 pt-2 space-y-2.5">
              <div>
                <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D]">
                  Nine Angelic Choirs
                </span>
                <h3 class="font-serif text-lg md:text-xl font-bold text-[#322D29] leading-snug">
                  {{ currentStep.choirName || 'St. Michael the Archangel' }}
                </h3>
              </div>
              <p class="font-serif text-xs md:text-sm text-[#322D29]/80 leading-relaxed border-l-2 border-[#72383D]/60 pl-3">
                Whoever honors St. Michael with these nine salutations before Holy Communion will be escorted by an angel from each of the nine celestial choirs.
              </p>
            </div>
          </div>

          <div class="devotion-card-surface p-4 rounded-[2rem]">
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

    <!-- Prayer Intentions Modal -->
    <PrayerIntentionsModal 
      :isOpen="showIntentionsModal"
      devotionKey="st_michael"
      devotionTitle="Chaplet of St. Michael the Archangel"
      @close="showIntentionsModal = false"
    />

    <!-- Global Footer Navigation -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import DevotionHeader from '../components/devotion/DevotionHeader.vue';
import DevotionControls from '../components/devotion/DevotionControls.vue';
import StMichaelBeads from '../components/stmichael/StMichaelBeads.vue';
import PrayerIntentionsCard from '../components/intentions/PrayerIntentionsCard.vue';
import PrayerIntentionsModal from '../components/intentions/PrayerIntentionsModal.vue';
import { usePrayerIntentions } from '../composables/usePrayerIntentions';
import prayersData from '../data/prayers.json';
import { generateStMichaelSteps, type StMichaelStep } from '../data/stMichaelData';
import { useSwipe } from '../composables/useSwipe';
import { useFullscreen } from '../composables/useFullscreen';
import { getPrayerAudioUrl } from '../utils/audioHelper';

const {
  activeIntentions,
  activeCount,
  hasIntentions
} = usePrayerIntentions('st_michael');

const showIntentionsModal = ref(false);
const showLatin = ref(false);
const currentIndex = ref(0);
const expandedSection = ref<'intro' | 'choirs' | 'pendant' | 'concluding'>('intro');
const steps = computed<StMichaelStep[]>(() => generateStMichaelSteps(hasIntentions.value));

watch(hasIntentions, () => {
  if (currentIndex.value >= steps.value.length) {
    currentIndex.value = 0;
  }
});

const swipeContainer = ref<HTMLElement | null>(null);
const cardContainerRef = ref<HTMLElement | null>(null);

const { isFullscreen, isFullscreenSupported, toggleFullscreen } = useFullscreen(cardContainerRef);

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
