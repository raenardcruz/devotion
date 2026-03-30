<template>
  <div class="min-h-screen bg-stone-900 text-stone-100 selection:bg-amber-500/30 flex flex-col relative font-sans overflow-hidden">
    <!-- Back Button -->
    <div class="absolute top-6 left-6 z-50">
      <router-link to="/" class="group flex items-center space-x-3 text-white/70 hover:text-white bg-black/20 hover:bg-black/40 backdrop-blur-md px-5 py-2.5 rounded-full transition-all duration-300 border border-white/10 hover:border-white/30">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform">
          <path d="m15 18-6-6 6-6"/>
        </svg>
        <span class="font-medium tracking-wide text-sm">Return Home</span>
      </router-link>
    </div>

    <!-- Latin Toggle -->
    <div class="absolute top-6 right-6 z-50">
      <div class="flex items-center bg-black/20 backdrop-blur-md p-1 rounded-full border border-white/10">
        <button 
          @click="showLatin = false"
          :class="[!showLatin ? 'bg-amber-500/20 text-amber-100 border-amber-500/30' : 'text-stone-400 hover:text-stone-200']"
          class="px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-widest transition-all duration-300 border border-transparent"
        >
          EN
        </button>
        <button 
          @click="showLatin = true"
          :class="[showLatin ? 'bg-amber-500/20 text-amber-100 border-amber-500/30' : 'text-stone-400 hover:text-stone-200']"
          class="px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-widest transition-all duration-300 border border-transparent"
        >
          LA
        </button>
      </div>
    </div>

    <!-- Dynamic Background Blur -->
    <transition name="bg-fade" mode="out-in">
        <div :key="currentSetName" class="fixed inset-0 z-0 transition-all duration-1000 ease-in-out" :style="{
                backgroundImage: `url('https://images.unsplash.com/photo-1569845177077-2a37322a60c7?q=80&w=2248&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                filter: 'brightness(0.6) blur(20px) contrast(1.1)'
            }"></div>
    </transition>
    
    <!-- Gradient Overlay -->
    <div class="fixed inset-0 z-0 bg-gradient-to-t from-stone-900 via-stone-900/40 to-stone-900/20 pointer-events-none"></div>

    <div class="relative z-10 max-w-5xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">

        <RosaryHeader 
            :sets="sets" 
            :currentSetName="currentSetName" 
            :todayMystery="todayMystery" 
            @changeSet="changeSet" 
        />

        <!-- Main Content Area -->
        <main ref="swipeContainer" class="flex-grow flex flex-col justify-center gap-10 md:gap-14 py-8 touch-pan-y">


            <transition name="fade-slide" mode="out-in">
                <div :key="currentSetName + currentStepIndex + (showLatin ? '-la' : '-en')" class="flex flex-col gap-8 md:gap-12 w-full">
                    <MysteryInfo 
                        :mysteryIndex="currentStep.decadeNumber ? currentStep.decadeNumber - 1 : 0" 
                        :currentMystery="getCurrentMystery(currentStep)"
                        v-if="currentStep.type !== 'intro' && currentStep.type !== 'opening' && currentStep.type !== 'closing'"
                    />

                    <RosaryBeads :beadIndex="currentStep.beadNumber || 0" :stepType="currentStep.type" />

                    <RosaryCard 
                        :currentStep="currentStep" 
                        :showLatin="showLatin"
                    />
                </div>
            </transition>

            <RosaryControls 
                :totalSteps="steps.length"
                :currentStepIndex="currentStepIndex"
                @next="nextStep" 
                @prev="prevStep" 
            />
        </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useSwipe } from '../composables/useSwipe';
import { useDate } from '../composables/useDate';
import { ROSARY_DATA, generateRosarySteps, type Mystery, type RosaryStep } from '../components/rosary/rosaryData';
import prayerData from '../data/prayers.json';

import RosaryHeader from '../components/rosary/RosaryHeader.vue';
import MysteryInfo from '../components/rosary/MysteryInfo.vue';
import RosaryBeads from '../components/rosary/RosaryBeads.vue';
import RosaryCard from '../components/rosary/VerseCard.vue'; // Renaming for better context
import RosaryControls from '../components/rosary/RosaryControls.vue';

const currentSetName = ref('Joyful');
const currentStepIndex = ref(0);
const showLatin = ref(false);
const todayMystery = ref('');
const sets = Object.keys(ROSARY_DATA);
const swipeContainer = ref<HTMLElement | null>(null);

const { getLocalDay } = useDate();

useSwipe(swipeContainer, {
    onSwipeLeft: () => nextStep(),
    onSwipeRight: () => prevStep(),
});

const getMysteryForDay = (): string => {
    const day = getLocalDay();
    switch (day) {
        case 0: return 'Glorious'; // Sunday
        case 1: return 'Joyful';   // Monday
        case 2: return 'Sorrowful';// Tuesday
        case 3: return 'Glorious'; // Wednesday
        case 4: return 'Luminous'; // Thursday
        case 5: return 'Sorrowful';// Friday
        case 6: return 'Joyful';   // Saturday
        default: return 'Joyful';
    }
};

onMounted(() => {
    const daily = getMysteryForDay();
    todayMystery.value = daily;
    currentSetName.value = daily;
});

const steps = computed(() => generateRosarySteps(currentSetName.value));

const currentStep = computed((): RosaryStep => {
    const step = steps.value[currentStepIndex.value];
    if (!step) return steps.value[0]!;

    // Enrich with prayer data
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

const nextStep = () => {
    if (currentStepIndex.value < steps.value.length - 1) {
        currentStepIndex.value++;
    } else {
        // Reset
        currentStepIndex.value = 0;
    }
};

const prevStep = () => {
    if (currentStepIndex.value > 0) {
        currentStepIndex.value--;
    }
};

const changeSet = (setName: string) => {
    currentSetName.value = setName;
    currentStepIndex.value = 0;
};

// Reset index when set changes
watch(currentSetName, () => {
    currentStepIndex.value = 0;
});
</script>

<style scoped>
/* Transitions moved to global style.css */
</style>