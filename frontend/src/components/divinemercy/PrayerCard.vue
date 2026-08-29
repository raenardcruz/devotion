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
