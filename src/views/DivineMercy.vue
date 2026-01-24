<template>
  <div class="min-h-screen bg-stone-50 text-stone-900 selection:bg-rose-100 flex flex-col relative font-sans overflow-hidden">
    <!-- Back Button -->
    <div class="absolute top-4 left-4 z-50">
      <router-link to="/" class="flex items-center space-x-2 text-rose-900/80 hover:text-rose-900 bg-white/40 hover:bg-white/60 backdrop-blur-sm px-4 py-2 rounded-full transition-all text-sm font-medium border border-rose-100/20">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="m15 18-6-6 6-6"/>
        </svg>
        <span>Home</span>
      </router-link>
    </div>

    <!-- Meditative Background -->
    <div class="fixed inset-0 z-0 opacity-20 transition-all duration-1000 ease-in-out bg-rays"
        style="background-image: url('https://images.unsplash.com/photo-1621610212025-da775e84bea9?q=80&w=1740&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'); background-size: cover; background-position: center;">
    </div>
    <!-- Light Rays Overlay -->
    <div class="fixed inset-0 z-0 bg-rays pointer-events-none"></div>

    <div class="relative z-10 max-w-4xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">

        <DivineMercyHeader />

        <!-- Main Prayer Area -->
        <main class="flex-grow flex flex-col justify-center gap-8">

            <transition name="fade" mode="out-in">
                <div :key="currentStepIndex" class="w-full">

                    <PhaseLabel :currentStep="currentStep" />

                    <DivineMercyBeads 
                        :currentStep="currentStep" 
                        :beadInDecade="beadInDecade" 
                    />

                    <PrayerCard :currentStep="currentStep" />
                </div>
            </transition>

            <DivineMercyControls 
                :currentStepIndex="currentStepIndex"
                :currentStep="currentStep" 
                :decadeIndex="decadeIndex" 
                :totalSteps="steps.length" 
                @next="next" 
                @prev="prev" 
            />
        </main>

        <footer class="mt-12 text-center text-stone-400 text-[10px] py-4 uppercase tracking-[0.2em]">
            <p>Jesus, I Trust In You</p>
        </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { DIVINE_MERCY_STEPS as steps, type Step } from '../components/divinemercy/divineMercyData';
import DivineMercyHeader from '../components/divinemercy/DivineMercyHeader.vue';
import PhaseLabel from '../components/divinemercy/PhaseLabel.vue';
import DivineMercyBeads from '../components/divinemercy/DivineMercyBeads.vue';
import PrayerCard from '../components/divinemercy/PrayerCard.vue';
import DivineMercyControls from '../components/divinemercy/DivineMercyControls.vue';

const currentStepIndex = ref(0);
const decadeIndex = ref(1);
const beadInDecade = ref(-1); // -1 is Eternal Father, 1-10 are Passion beads

const currentStep = computed((): Step => {
    const step = steps[currentStepIndex.value];
    if (!step) return steps[0]!; // Fallback
    
    if (step.type === 'decade') {
        if (beadInDecade.value === -1) {
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
            if (decadeIndex.value < 5) {
                decadeIndex.value++;
                beadInDecade.value = -1;
            } else {
                // Decades finished, move to closing
                currentStepIndex.value = 8;
            }
        }
    } else {
        if (currentStepIndex.value === 5) { // Apostles creed finished
            currentStepIndex.value = 6; // Start decades
            decadeIndex.value = 1;
            beadInDecade.value = -1;
        } else if (currentStepIndex.value === steps.length - 1) {
            // Reset
            currentStepIndex.value = 0;
            decadeIndex.value = 1;
            beadInDecade.value = -1;
        } else {
            currentStepIndex.value++;
        }
    }
};

const prev = () => {
    if (beadInDecade.value > -1) {
        beadInDecade.value--;
    } else if (decadeIndex.value > 1) {
        decadeIndex.value--;
        beadInDecade.value = 10;
    } else {
        if (currentStepIndex.value > 0) {
            currentStepIndex.value--;
            if (currentStepIndex.value === 7) currentStepIndex.value = 5; // Skip placeholder decade steps backward
        }
    }
};
</script>

<style scoped>
.bg-rays {
    background: radial-gradient(circle at 50% 50%, rgba(255, 255, 255, 0.8) 0%, rgba(255, 220, 220, 0.2) 30%, rgba(220, 230, 255, 0.2) 60%, transparent 100%);
}

.fade-enter-active,
.fade-leave-active {
    transition: all 0.4s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}
</style>