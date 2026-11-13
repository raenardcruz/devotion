<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Chaplet Container -->
    <main class="flex-grow max-w-4xl mx-auto px-4 py-8 flex flex-col w-full relative z-10">

        <!-- Clean Unified Devotion Header with Collapsible Settings Menu -->
        <DevotionHeader
          title="The Divine Mercy Chaplet"
          subtitle="Chaplet of Divine Mercy"
          v-model:showLatin="showLatin"
          v-model:autoPlay="autoPlay"
          v-model:isShortVersion="isShortVersion"
          :hasShortMode="true"
          :activeIntentionsCount="activeCount"
          :isFullscreen="isFullscreen"
          :isFullscreenSupported="isFullscreenSupported"
          @open-intentions="showIntentionsModal = true"
          @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
        />

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

                    <!-- Intentions Card or Regular Prayer Card -->
                    <PrayerIntentionsCard 
                        v-if="currentStep.id === 'prayer-intentions'"
                        :intentions="activeIntentions"
                        :showLatin="showLatin"
                        :isFullscreen="isFullscreen"
                        @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
                        @open-intentions-modal="showIntentionsModal = true"
                    />

                    <!-- Regular Prayer Card -->
                    <PrayerCard 
                        v-else
                        :currentStep="currentStep" 
                        :showLatin="showLatin" 
                        :isFullscreen="isFullscreen"
                        :isFullscreenSupported="isFullscreenSupported"
                        @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
                    />
                </div>
            </transition>

            <!-- Navigation & Audio Controls -->
            <DevotionControls
                :currentStepIndex="currentStepIndex"
                :totalSteps="steps.length"
                :isLastStep="currentStepIndex === steps.length - 1"
                :hasAudio="Boolean(currentStep.id && getPrayerAudioUrl(currentStep.id))"
                :isPlaying="isPlaying"
                label="Chaplet"
                @next="next"
                @prev="prev"
                @toggle-audio="toggleAudio"
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

    <!-- Prayer Intentions Modal -->
    <PrayerIntentionsModal 
      :isOpen="showIntentionsModal"
      devotionKey="divine_mercy"
      devotionTitle="The Divine Mercy Chaplet"
      @close="showIntentionsModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useSwipe } from '../composables/useSwipe';
import { useFullscreen } from '../composables/useFullscreen';
import { usePrayerIntentions } from '../composables/usePrayerIntentions';
import { generateDivineMercySteps, type Step } from '../components/divinemercy/divineMercyData';
import { getPrayerAudioUrl } from '../utils/audioHelper';

import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import PhaseLabel from '../components/divinemercy/PhaseLabel.vue';
import DivineMercyBeads from '../components/divinemercy/DivineMercyBeads.vue';
import PrayerCard from '../components/divinemercy/PrayerCard.vue';
import PrayerIntentionsCard from '../components/intentions/PrayerIntentionsCard.vue';
import PrayerIntentionsModal from '../components/intentions/PrayerIntentionsModal.vue';
import DevotionHeader from '../components/devotion/DevotionHeader.vue';
import DevotionControls from '../components/devotion/DevotionControls.vue';

const {
    activeIntentions,
    activeCount,
    hasIntentions
} = usePrayerIntentions('divine_mercy');

const showIntentionsModal = ref(false);
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

const steps = computed(() => generateDivineMercySteps(hasIntentions.value));

const currentStep = computed((): Step => {
    const stepList = steps.value;
    const step = stepList[currentStepIndex.value];
    if (!step) return stepList[0]!; // Fallback
    
    if (step.type === 'decade') {
        if (beadInDecade.value === 0) {
            return stepList.find(s => s.id === 'eternal-father') || step;
        } else {
            return stepList.find(s => s.id === 'sorrowful-passion') || step;
        }
    }
    return step;
});

const next = () => {
    const stepList = steps.value;
    const step = stepList[currentStepIndex.value];
    if (!step) return;

    if (step.type === 'decade') {
        if (beadInDecade.value < 10) {
            beadInDecade.value++;
        } else {
            if (isShortVersion.value || decadeIndex.value >= 5) {
                const holyGodIdx = stepList.findIndex(s => s.id === 'holy-god');
                currentStepIndex.value = holyGodIdx !== -1 ? holyGodIdx : stepList.length - 3;
            } else {
                decadeIndex.value++;
                beadInDecade.value = 0;
            }
        }
    } else {
        const decadeStartIdx = stepList.findIndex(s => s.type === 'decade');
        if (decadeStartIdx !== -1 && currentStepIndex.value === decadeStartIdx - 1) {
            currentStepIndex.value = decadeStartIdx;
            decadeIndex.value = 1;
            beadInDecade.value = 0;
        } else if (currentStepIndex.value === stepList.length - 1) {
            currentStepIndex.value = 0;
            decadeIndex.value = 1;
            beadInDecade.value = 0;
        } else {
            currentStepIndex.value++;
        }
    }
};

const prev = () => {
    const stepList = steps.value;
    if (beadInDecade.value > 0) {
        beadInDecade.value--;
    } else if (decadeIndex.value > 1 && !isShortVersion.value) {
        decadeIndex.value--;
        beadInDecade.value = 10;
    } else {
        if (currentStepIndex.value > 0) {
            const holyGodIdx = stepList.findIndex(s => s.id === 'holy-god');
            const decadeStartIdx = stepList.findIndex(s => s.type === 'decade');
            if (holyGodIdx !== -1 && currentStepIndex.value === holyGodIdx) {
                currentStepIndex.value = decadeStartIdx;
                decadeIndex.value = isShortVersion.value ? 1 : 5;
                beadInDecade.value = 10;
            } else if (decadeStartIdx !== -1 && currentStepIndex.value === decadeStartIdx) {
                currentStepIndex.value = decadeStartIdx - 1;
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