<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Chaplet Container -->
    <main class="flex-grow max-w-4xl mx-auto px-4 py-8 flex flex-col w-full relative z-10">

        <!-- Sub-header & Toggles -->
        <div class="flex flex-col md:flex-row md:items-center md:justify-between border-b border-[#D1C7BD]/60 pb-4 mb-6 gap-4 animate-fade-in-down">
            <h1 class="font-serif text-2xl md:text-3xl bg-gradient-to-r from-[#72383D] via-[#AC9C8D] to-[#322D29] bg-clip-text text-transparent font-bold">
                The Divine Mercy Chaplet
            </h1>
            
            <div class="flex flex-wrap items-center gap-3">
                <!-- Auto Play Toggle Switch -->
                <label class="flex items-center space-x-2 text-[10px] font-bold uppercase tracking-wider text-[#322D29]/70 cursor-pointer select-none">
                    <span>Auto-Play</span>
                    <input type="checkbox" v-model="autoPlay" class="sr-only peer" />
                    <div class="relative w-8 h-4.5 bg-white/60 border border-[#D1C7BD]/80 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[1px] after:start-[1px] after:bg-[#322D29]/40 peer-checked:after:bg-[#72383D] after:border-[#D1C7BD] after:border after:rounded-full after:h-3.5 after:w-3.5 after:transition-all peer-checked:bg-[#72383D]/20"></div>
                </label>

                <!-- Manual Play/Pause Button -->
                <button 
                  v-if="currentStep.id && getPrayerAudioUrl(currentStep.id)"
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

                <!-- Short/Full Version Toggle -->
                <div class="flex items-center bg-white/70 border border-[#D1C7BD]/80 p-0.5 rounded-full w-fit backdrop-blur-md">
                    <button 
                      @click="isShortVersion = false"
                      :class="[!isShortVersion ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold' : 'text-[#322D29]/70 hover:text-[#72383D]']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0 cursor-pointer"
                    >
                      Full
                    </button>
                    <button 
                      @click="isShortVersion = true"
                      :class="[isShortVersion ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold' : 'text-[#322D29]/70 hover:text-[#72383D]']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0 cursor-pointer"
                    >
                      Short
                    </button>
                </div>

                <!-- Fullscreen Presentation Toggle Button -->
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

                <!-- Latin/English Toggle -->
                <div class="flex items-center bg-white/70 border border-[#D1C7BD]/80 p-0.5 rounded-full w-fit backdrop-blur-md">
                    <button 
                      @click="showLatin = false"
                      :class="[!showLatin ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold' : 'text-[#322D29]/70 hover:text-[#72383D]']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0 cursor-pointer"
                    >
                      EN
                    </button>
                    <button 
                      @click="showLatin = true"
                      :class="[showLatin ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white font-bold' : 'text-[#322D29]/70 hover:text-[#72383D]']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0 cursor-pointer"
                    >
                      LA
                    </button>
                </div>
            </div>
        </div>

        <!-- Main Prayer Card Area -->
        <section 
            ref="cardContainerRef" 
            :class="[
                'flex-grow flex flex-col justify-between gap-6 py-6 touch-pan-y animate-fade-in-up delay-150 transition-all duration-300',
                isFullscreen ? 'card-fullscreen-mode !bg-[#EFE9E1] p-6 md:p-12 overflow-y-auto' : ''
            ]"
        >
            <transition name="fade-slide" mode="out-in">
                <div :key="currentStepIndex + (showLatin ? '-la' : '-en') + (isShortVersion ? '-short' : '-full')" class="w-full flex flex-col gap-6">

                    <!-- Phase label / Category -->
                    <PhaseLabel :currentStep="currentStep" :isShort="isShortVersion" :showLatin="showLatin" />

                    <!-- Chaplet Beads Progress -->
                    <div class="min-h-[50px] flex items-center justify-center">
                         <DivineMercyBeads 
                            :currentStep="currentStep" 
                            :beadInDecade="beadInDecade" 
                        />
                    </div>

                    <!-- Prayer Card -->
                    <PrayerCard 
                        :currentStep="currentStep" 
                        :showLatin="showLatin" 
                        :isFullscreen="isFullscreen"
                        :isFullscreenSupported="isFullscreenSupported"
                        @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
                    />
                </div>
            </transition>

            <!-- Navigation Controls -->
            <DivineMercyControls 
                :currentStepIndex="currentStepIndex"
                :currentStep="currentStep" 
                :decadeIndex="decadeIndex" 
                :totalSteps="steps.length" 
                :isShort="isShortVersion"
                @next="next" 
                @prev="prev" 
            />

            <!-- Fullscreen Keyboard Helper Hint -->
            <div v-if="isFullscreen" class="text-center text-xs text-black/50 font-sans tracking-wide mt-2">
                Press <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">←</kbd> <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">→</kbd> to Navigate &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">Space</kbd> Play/Pause &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">ESC</kbd> Exit Fullscreen
            </div>
        </section>

        <footer class="mt-8 text-center text-parchment-secondary/40 text-xs font-bold py-4 uppercase tracking-[0.25em] select-none">
            <p>Jesus, I Trust In You</p>
            <div v-if="isFullscreen" class="mt-3 text-center text-xs text-black/50 font-sans normal-case tracking-wide">
                Press <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">←</kbd> <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">→</kbd> to Navigate &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">Space</kbd> Play/Pause &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">F</kbd> Fullscreen &bull; <kbd class="px-1.5 py-0.5 bg-black/10 rounded font-mono text-[10px]">ESC</kbd> Exit
            </div>
        </footer>
    </main>

    <!-- Global Footer -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useSwipe } from '../composables/useSwipe';
import { useFullscreen } from '../composables/useFullscreen';
import { DIVINE_MERCY_STEPS as steps, type Step } from '../components/divinemercy/divineMercyData';
import { getPrayerAudioUrl } from '../utils/audioHelper';

import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import PhaseLabel from '../components/divinemercy/PhaseLabel.vue';
import DivineMercyBeads from '../components/divinemercy/DivineMercyBeads.vue';
import PrayerCard from '../components/divinemercy/PrayerCard.vue';
import DivineMercyControls from '../components/divinemercy/DivineMercyControls.vue';

const currentStepIndex = ref(0);
const decadeIndex = ref(1);
const showLatin = ref(false);
const isShortVersion = ref(false);

const beadInDecade = ref(0); // 0 is Eternal Father, 1-10 are Passion beads
const swipeContainer = ref<HTMLElement | null>(null);
const cardContainerRef = ref<HTMLElement | null>(null);

const { isFullscreen, isFullscreenSupported, toggleFullscreen } = useFullscreen(cardContainerRef);

// Audio State
const currentAudio = ref<HTMLAudioElement | null>(null);
const isPlaying = ref(false);
const autoPlay = ref(true);

useSwipe(swipeContainer, {
    onSwipeLeft: () => next(),
    onSwipeRight: () => prev(),
});

const currentStep = computed((): Step => {
    const step = steps[currentStepIndex.value];
    if (!step) return steps[0]!; // Fallback
    
    if (step.type === 'decade') {
        if (beadInDecade.value === 0) {
            return steps[6]!; // Eternal Father step
        } else {
            return steps[7]!; // Passion step
        }
    }
    return step!;
});

const next = () => {
    const step = steps[currentStepIndex.value];
    if (!step) return;

    if (step.type === 'decade') {
        if (beadInDecade.value < 10) {
            beadInDecade.value++;
        } else {
            if (isShortVersion.value || decadeIndex.value >= 5) {
                currentStepIndex.value = 8;
            } else {
                decadeIndex.value++;
                beadInDecade.value = 0;
            }
        }
    } else {
        if (currentStepIndex.value === 5) {
            currentStepIndex.value = 6;
            decadeIndex.value = 1;
            beadInDecade.value = 0;
        } else if (currentStepIndex.value === steps.length - 1) {
            currentStepIndex.value = 0;
            decadeIndex.value = 1;
            beadInDecade.value = 0;
        } else {
            currentStepIndex.value++;
        }
    }
};

const prev = () => {
    if (beadInDecade.value > 0) {
        beadInDecade.value--;
    } else if (decadeIndex.value > 1 && !isShortVersion.value) {
        decadeIndex.value--;
        beadInDecade.value = 10;
    } else {
        if (currentStepIndex.value > 0) {
            if (currentStepIndex.value === 8) {
                currentStepIndex.value = 6;
                decadeIndex.value = isShortVersion.value ? 1 : 5;
                beadInDecade.value = 10;
            } else if (currentStepIndex.value === 6) {
                currentStepIndex.value = 5;
            } else {
                currentStepIndex.value--;
            }
        }
    }
};

// Audio Controls
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
    } else if (currentStep.value.id) {
        playAudio(currentStep.value.id);
    }
};

const handleKeyDown = (e: KeyboardEvent) => {
    if (['INPUT', 'TEXTAREA', 'SELECT'].includes((e.target as HTMLElement)?.tagName)) {
        return;
    }
    if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
        e.preventDefault();
        next();
    } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
        e.preventDefault();
        prev();
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
    const savedAutoPlay = localStorage.getItem('dm-autoplay');
    if (savedAutoPlay !== null) {
        autoPlay.value = savedAutoPlay === 'true';
    }
});

// Watch currentStepIndex and beadInDecade change to trigger audio auto-play
watch([currentStepIndex, beadInDecade], () => {
    stopAudio();
    if (autoPlay.value && currentStep.value.id) {
        playTimeout = setTimeout(() => {
            playAudio(currentStep.value.id);
        }, 300);
    }
});

// Sync Autoplay Preference
watch(autoPlay, (newVal) => {
    localStorage.setItem('dm-autoplay', String(newVal));
    if (!newVal) {
        stopAudio();
    } else if (currentStep.value.id) {
        playAudio(currentStep.value.id);
    }
});

onBeforeUnmount(() => {
    window.removeEventListener('keydown', handleKeyDown);
    stopAudio();
});
</script>

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