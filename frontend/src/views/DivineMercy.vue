<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-28 sm:pb-20 selection:bg-parchment-secondary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Chaplet Container -->
    <main class="flex-grow max-w-4xl mx-auto px-4 py-8 flex flex-col w-full">
        <!-- Sub-header & Toggles -->
        <div class="flex flex-col md:flex-row md:items-center md:justify-between border-b border-parchment-border pb-4 mb-6 gap-4 animate-fade-in-down">
            <h1 class="font-serif text-2xl md:text-3xl text-parchment-neutral font-medium">
                The Divine Mercy Chaplet
            </h1>
            
            <div class="flex flex-wrap items-center gap-3">
                <!-- Auto Play Toggle Switch -->
                <label class="flex items-center space-x-2 text-[10px] font-bold uppercase tracking-wider text-parchment-neutral/50 cursor-pointer select-none">
                    <span>Auto-Play</span>
                    <input type="checkbox" v-model="autoPlay" class="sr-only peer" />
                    <div class="relative w-8 h-4.5 bg-parchment-neutral-light border border-parchment-border rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[1px] after:start-[1px] after:bg-parchment-neutral/40 peer-checked:after:bg-parchment-secondary after:border-parchment-border after:border after:rounded-full after:h-3.5 after:w-3.5 after:transition-all peer-checked:bg-parchment-secondary/10"></div>
                </label>

                <!-- Manual Play/Pause Button -->
                <button 
                  v-if="currentStep.id && getPrayerAudioUrl(currentStep.id)"
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

                <!-- Short/Full Version Toggle -->
                <div class="flex items-center bg-parchment-neutral-light border border-parchment-border p-0.5 rounded-full w-fit">
                    <button 
                      @click="isShortVersion = false"
                      :class="[!isShortVersion ? 'bg-parchment-secondary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                    >
                      Full
                    </button>
                    <button 
                      @click="isShortVersion = true"
                      :class="[isShortVersion ? 'bg-parchment-secondary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                    >
                      Short
                    </button>
                </div>

                <!-- Latin/English Toggle -->
                <div class="flex items-center bg-parchment-neutral-light border border-parchment-border p-0.5 rounded-full w-fit">
                    <button 
                      @click="showLatin = false"
                      :class="[!showLatin ? 'bg-parchment-secondary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                    >
                      EN
                    </button>
                    <button 
                      @click="showLatin = true"
                      :class="[showLatin ? 'bg-parchment-secondary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                      class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                    >
                      LA
                    </button>
                </div>
            </div>
        </div>

        <!-- Main Prayer Card Area -->
        <main ref="swipeContainer" class="flex-grow flex flex-col justify-center gap-6 py-6 touch-pan-y animate-fade-in-up delay-150">
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
                    <PrayerCard :currentStep="currentStep" :showLatin="showLatin" />
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
        </main>

        <footer class="mt-8 text-center text-parchment-secondary/40 text-xs font-bold py-4 uppercase tracking-[0.25em] select-none">
            <p>Jesus, I Trust In You</p>
        </footer>
    </main>

    <!-- Global Footer -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useSwipe } from '../composables/useSwipe';
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

onMounted(() => {
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
</style>