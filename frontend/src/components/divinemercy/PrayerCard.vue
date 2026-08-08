<template>
    <div :class="[isFullscreen ? 'w-full max-w-4xl mx-auto my-auto flex-grow flex flex-col justify-center items-center px-4' : 'relative max-w-2xl mx-auto w-full']">
        <!-- Parchment Card Container -->
        <div :class="[isFullscreen ? 'bg-parchment-neutral-light border border-parchment-border p-8 md:p-12 rounded-[2.5rem] shadow-xl text-center w-full my-auto flex flex-col justify-center items-center' : 'bg-parchment-neutral-light border border-parchment-border p-6 md:p-8 rounded-[2rem] shadow-sm text-center min-h-[280px] flex flex-col justify-center']" class="transition-all duration-300">
            
            <transition name="fade-scale" mode="out-in">
                <div :key="(currentStep.id || currentStep.title || currentStep.content) + (showLatin ? '-la' : '-en')" class="space-y-6 w-full">
                    <div class="flex items-center justify-between border-b border-parchment-border/40 pb-3 mb-4">
                        <h3 v-if="currentStep.title" :class="[isFullscreen ? 'text-xs md:text-sm tracking-[0.2em]' : 'text-xs tracking-[0.2em]']" class="text-parchment-secondary font-bold uppercase">
                            {{ displayTitle }}
                        </h3>
                        <button 
                            @click="$emit('toggle-fullscreen')"
                            class="text-parchment-neutral/50 hover:text-parchment-primary transition-colors p-1.5 rounded-full hover:bg-parchment-primary/10 cursor-pointer ml-auto"
                            :title="isFullscreen ? 'Exit Fullscreen Card (ESC)' : 'Fullscreen Card Presentation (F)'"
                        >
                            <svg v-if="!isFullscreen" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                <polyline points="15 3 21 3 21 9"></polyline>
                                <polyline points="9 21 3 21 3 15"></polyline>
                                <line x1="21" y1="3" x2="14" y2="10"></line>
                                <line x1="3" y1="21" x2="10" y2="14"></line>
                            </svg>
                            <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                <polyline points="4 14 10 14 10 20"></polyline>
                                <polyline points="20 10 14 10 14 4"></polyline>
                                <line x1="14" y1="10" x2="21" y2="3"></line>
                                <line x1="10" y1="14" x2="3" y2="21"></line>
                            </svg>
                        </button>
                    </div>
                    
                    <InterlinearText 
                        v-if="interlinearWords" 
                        :words="interlinearWords" 
                        latin-class="text-parchment-secondary font-medium"
                    />

                    <p 
                        v-else
                        :class="[
                            isFullscreen 
                                ? 'text-2xl md:text-4xl lg:text-5xl font-serif leading-relaxed text-parchment-neutral py-4' 
                                : 'text-xl md:text-2xl font-serif leading-relaxed text-parchment-neutral',
                            (showLatin && currentStep.latin) ? 'text-parchment-secondary' : ''
                        ]"
                        class="whitespace-pre-line transition-all duration-300"
                    >
                        {{ (showLatin && currentStep.latin) ? currentStep.latin : currentStep.content }}
                    </p>

                    <!-- Latin Translation Reference (Hidden if interlinear is active) -->
                    <p 
                        v-if="showLatin && currentStep.latin && !interlinearWords"
                        :class="[isFullscreen ? 'mt-4 text-sm md:text-lg' : 'mt-4 text-xs']"
                        class="font-sans leading-relaxed text-parchment-neutral/50 max-w-xl mx-auto opacity-70 border-t border-parchment-border/40 pt-4"
                    >
                        {{ currentStep.content }}
                    </p>
                </div>
            </transition>

        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Step } from './divineMercyData';
import { getInterlinearPairs } from '../../utils/interlinearMappers';
import InterlinearText from '../common/InterlinearText.vue';
import prayerData from '../../data/prayers.json';

const props = defineProps<{
    currentStep: Step;
    showLatin?: boolean;
    isFullscreen?: boolean;
}>();

defineEmits<{
    (e: 'toggle-fullscreen'): void;
}>();

const interlinearWords = computed(() => {
    if (!props.showLatin || !props.currentStep.latin) return null;
    return getInterlinearPairs(props.currentStep.latin);
});

const displayTitle = computed(() => {
    const title = props.currentStep.title;
    if (props.showLatin) {
        const prayer = prayerData.find(p => p.name === title || (p.id === 'dm-closing' && title === 'Optional Closing'));
        if (prayer && (prayer as any).latinName) {
            return (prayer as any).latinName;
        }
    }
    return title;
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
