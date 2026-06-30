<template>
    <div class="relative max-w-2xl mx-auto w-full">
        <!-- Parchment Card Wrapper -->
        <div class="bg-parchment-neutral-light border border-parchment-border p-6 md:p-8 rounded-[2rem] shadow-sm flex flex-col justify-between min-h-[300px] transition-all duration-300">
            
            <transition name="fade-scale" mode="out-in">
                <div :key="(currentStep.title || currentStep.content) + (showLatin ? '-la' : '-en')" class="space-y-5">
                    
                    <!-- Card Header (Type & Ref) -->
                    <div class="flex items-center justify-between border-b border-parchment-border/40 pb-3">
                        <span class="text-parchment-secondary font-bold tracking-[0.15em] text-[11px] uppercase">
                            {{ stepHeading }}
                        </span>
                        <div v-if="currentStep.verse" class="flex items-center space-x-1.5 text-parchment-primary-dark text-xs font-semibold">
                            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
                                <path d="M6 6h10"></path>
                                <path d="M6 10h10"></path>
                            </svg>
                            <span>{{ currentStep.verse.ref }}</span>
                        </div>
                    </div>
                    
                    <!-- Mystery Title -->
                    <h3 v-if="currentStep.mysteryTitle || currentStep.title" class="text-xl md:text-2xl font-serif text-parchment-neutral font-medium">
                        {{ currentStep.mysteryTitle || currentStep.title }}
                    </h3>

                    <!-- Scripture Verse (Mockup blockquote style, always English) -->
                    <div v-if="currentStep.verse" class="border-l-4 border-parchment-primary pl-4 py-1.5 my-4">
                        <p class="text-base md:text-lg font-serif italic text-parchment-neutral/80 leading-relaxed">
                            "{{ currentStep.verse.text }}"
                        </p>
                    </div>

                    <!-- Prayer Text / Latin / Interlinear -->
                    <div class="mt-4 pt-4 border-t border-parchment-border/20">
                        <InterlinearText 
                            v-if="interlinearWords" 
                            :words="interlinearWords" 
                            latin-class="text-parchment-primary-dark font-medium"
                        />

                        <p 
                            v-else-if="currentStep.content"
                            class="text-base md:text-lg font-sans leading-relaxed text-parchment-neutral/90 whitespace-pre-line"
                            :class="{ 'font-serif italic text-parchment-primary-dark': showLatin && currentStep.latin }"
                        >
                            {{ (showLatin && currentStep.latin) ? currentStep.latin : currentStep.content }}
                        </p>
     
                        <!-- Latin Translation Reference (Hidden if interlinear is active) -->
                        <p 
                            v-if="showLatin && currentStep.latin && !interlinearWords"
                            class="mt-3 text-xs font-sans text-parchment-neutral/50 italic"
                        >
                            {{ currentStep.content }}
                        </p>
                    </div>
                </div>
            </transition>

        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { RosaryStep } from './rosaryData';
import { getInterlinearPairs } from '../../utils/interlinearMappers';
import InterlinearText from '../common/InterlinearText.vue';

const props = defineProps<{
    currentStep: RosaryStep;
    showLatin?: boolean;
}>();

const interlinearWords = computed(() => {
    if (!props.showLatin || !props.currentStep.latin) return null;
    return getInterlinearPairs(props.currentStep.latin);
});

const stepHeading = computed(() => {
    const step = props.currentStep;
    if (step.type === 'intro' || step.type === 'opening') {
        return 'Opening Prayer';
    }
    if (step.type === 'closing') {
        return 'Closing Prayer';
    }
    if (step.type === 'mystery-header') {
        return 'Mystery Announcement';
    }
    
    // Decade bead number (1-10) or decade start/end
    if (step.decadeNumber) {
        const ordinals = ['First', 'Second', 'Third', 'Fourth', 'Fifth'];
        const name = ordinals[step.decadeNumber - 1] || 'Decade';
        if (step.type === 'decade-start') {
            return `${name} Decade - Our Father`;
        }
        if (step.type === 'decade-end') {
            return `${name} Decade - Doxology`;
        }
        return `${name} Mystery - Bead ${step.beadNumber || 1}`;
    }
    
    return step.title || 'Rosary Devotion';
});
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.4s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.98);
  filter: blur(2px);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(1.02);
  filter: blur(2px);
}
</style>
