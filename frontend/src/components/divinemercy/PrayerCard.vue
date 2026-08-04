<template>
    <div class="relative max-w-2xl mx-auto w-full">
        <!-- Parchment Card Container -->
        <div class="bg-parchment-neutral-light border border-parchment-border p-6 md:p-8 rounded-[2rem] shadow-sm text-center min-h-[280px] flex flex-col justify-center transition-all duration-300">
            
            <transition name="fade-scale" mode="out-in">
                <div :key="(currentStep.id || currentStep.title || currentStep.content) + (showLatin ? '-la' : '-en')" class="space-y-6">
                    <h3 v-if="currentStep.title" class="text-parchment-secondary font-bold tracking-[0.2em] text-xs uppercase">
                        {{ displayTitle }}
                    </h3>
                    
                    <InterlinearText 
                        v-if="interlinearWords" 
                        :words="interlinearWords" 
                        latin-class="text-parchment-secondary font-medium"
                    />

                    <p 
                        v-else
                        class="text-xl md:text-2xl font-serif leading-relaxed text-parchment-neutral whitespace-pre-line transition-all duration-300"
                        :class="{ 'text-parchment-secondary': showLatin && currentStep.latin }"
                    >
                        {{ (showLatin && currentStep.latin) ? currentStep.latin : currentStep.content }}
                    </p>

                    <!-- Latin Translation Reference (Hidden if interlinear is active) -->
                    <p 
                        v-if="showLatin && currentStep.latin && !interlinearWords"
                        class="mt-4 text-xs font-sans leading-relaxed text-parchment-neutral/50 max-w-xl mx-auto opacity-70 border-t border-parchment-border/40 pt-4"
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
