<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-24 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 py-8 flex flex-col">
      
      <!-- Sub-header & Controls Row -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between border-b border-parchment-border pb-4 mb-4 gap-4 animate-fade-in-down">
          <h1 class="font-serif text-2xl md:text-3xl text-parchment-neutral font-medium">
              The {{ currentSetName }} Mysteries
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
});

const steps = computed(() => generateRosarySteps(currentSetName.value));

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
                latin: step.latin || prayer.latin
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
const playAudio = (prayerId: string | undefined) => {
    if (!prayerId) return;
    const url = getPrayerAudioUrl(prayerId);
    if (!url) return;
    
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
    if (currentAudio.value) {
        currentAudio.value.pause();
        currentAudio.value = null;
    }
    isPlaying.value = false;
};

const toggleAudio = () => {
    if (isPlaying.value) {
        stopAudio();
    } else if (currentStep.value.prayerId) {
        playAudio(currentStep.value.prayerId);
    }
};

// Monitor currentStep change to trigger audio auto-play
watch(currentStep, (newStep) => {
    stopAudio();
    if (autoPlay.value && newStep.prayerId) {
        // Delay slightly to give page transition room to render nicely
        setTimeout(() => {
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