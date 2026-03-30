<template>
  <div class="min-h-screen bg-stone-900 text-stone-100 selection:bg-rose-500/30 flex flex-col relative font-sans overflow-hidden">
    <!-- Back Button -->
    <div class="absolute top-6 left-6 z-50">
      <router-link to="/" class="group flex items-center space-x-3 text-white/70 hover:text-white bg-black/20 hover:bg-black/40 backdrop-blur-md px-5 py-2.5 rounded-full transition-all duration-300 border border-white/10 hover:border-white/30">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform">
          <path d="m15 18-6-6 6-6"/>
        </svg>
        <span class="font-medium tracking-wide text-sm">Return Home</span>
      </router-link>
    </div>

    <!-- Latin & Version Toggles -->
    <div class="absolute top-6 right-6 z-50 flex items-center gap-3">
      <!-- Short Version Toggle -->
      <div class="flex items-center bg-black/20 backdrop-blur-md p-1 rounded-full border border-white/10">
        <button 
          @click="isShortVersion = false"
          :class="[!isShortVersion ? 'bg-rose-500/20 text-rose-100 border-rose-500/30' : 'text-stone-400 hover:text-stone-200']"
          class="px-4 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all duration-300 border border-transparent whitespace-nowrap"
        >
          Full
        </button>
        <button 
          @click="isShortVersion = true"
          :class="[isShortVersion ? 'bg-rose-500/20 text-rose-100 border-rose-500/30' : 'text-stone-400 hover:text-stone-200']"
          class="px-4 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all duration-300 border border-transparent whitespace-nowrap"
        >
          Short
        </button>
      </div>

      <!-- Latin Toggle -->
      <div class="flex items-center bg-black/20 backdrop-blur-md p-1 rounded-full border border-white/10">
        <button 
          @click="showLatin = false"
          :class="[!showLatin ? 'bg-rose-500/20 text-rose-100 border-rose-500/30' : 'text-stone-400 hover:text-stone-200']"
          class="px-4 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all duration-300 border border-transparent"
        >
          EN
        </button>
        <button 
          @click="showLatin = true"
          :class="[showLatin ? 'bg-rose-500/20 text-rose-100 border-rose-500/30' : 'text-stone-400 hover:text-stone-200']"
          class="px-4 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all duration-300 border border-transparent"
        >
          LA
        </button>
      </div>
    </div>

    <!-- Meditative Background -->
    <div class="fixed inset-0 z-0 transition-all duration-1000 ease-in-out" :style="{
                backgroundImage: `url('https://images.unsplash.com/photo-1621610212025-da775e84bea9?q=80&w=1740&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                filter: 'brightness(0.6) blur(20px) contrast(1.1)'
            }"></div>

    <!-- Gradient Overlay -->
    <div class="fixed inset-0 z-0 bg-gradient-to-t from-stone-900 via-stone-900/40 to-stone-900/20 pointer-events-none"></div>

    <div class="relative z-10 max-w-4xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">

        <DivineMercyHeader />

        <!-- Main Prayer Area -->
        <main ref="swipeContainer" class="flex-grow flex flex-col justify-center gap-6 md:gap-10 py-6 touch-pan-y">


            <transition name="fade-slide" mode="out-in">
                <div :key="currentStepIndex + (showLatin ? '-la' : '-en') + (isShortVersion ? '-short' : '-full')" class="w-full flex flex-col gap-6 md:gap-8">

                    <PhaseLabel :currentStep="currentStep" :isShort="isShortVersion" />

                    <div class="min-h-[60px] flex items-center justify-center">
                         <DivineMercyBeads 
                            :currentStep="currentStep" 
                            :beadInDecade="beadInDecade" 
                        />
                    </div>

                    <PrayerCard :currentStep="currentStep" :showLatin="showLatin" />
                </div>
            </transition>

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

        <footer class="mt-8 text-center text-rose-100/40 text-xs font-bold py-4 uppercase tracking-[0.3em] opacity-70">
            <p>Jesus, I Trust In You</p>
        </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useSwipe } from '../composables/useSwipe';
import { DIVINE_MERCY_STEPS as steps, type Step } from '../components/divinemercy/divineMercyData';

import DivineMercyHeader from '../components/divinemercy/DivineMercyHeader.vue';
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
            // Check short version
            if (isShortVersion.value || decadeIndex.value >= 5) {
                // Done with decades
                currentStepIndex.value = 8;
            } else {
                decadeIndex.value++;
                beadInDecade.value = 0;
            }
        }
    } else {
        if (currentStepIndex.value === 5) { // Apostles creed finished
            currentStepIndex.value = 6; // Start decades
            decadeIndex.value = 1;
            beadInDecade.value = 0;
        } else if (currentStepIndex.value === steps.length - 1) {
            // Reset
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
            // If going back from closing to decades
            if (currentStepIndex.value === 8) {
                currentStepIndex.value = 6;
                // In short version, going back from closing means going to bead 10 of first decade
                // Else it's bead 10 of 5th decade
                decadeIndex.value = isShortVersion.value ? 1 : 5;
                beadInDecade.value = 10;
            } else if (currentStepIndex.value === 6) {
                // Going back from decades to intro
                currentStepIndex.value = 5; // Apostles creed
            } else {
                currentStepIndex.value--;
            }
        }
    }
};
</script>

<style scoped>
/* Transitions moved to global style.css */
</style>