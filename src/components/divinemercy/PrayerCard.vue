<template>
    <div class="relative group max-w-2xl mx-auto w-full">
        <div class="absolute inset-0 bg-white/5 blur-xl rounded-[2.5rem] -z-10 transition-all duration-1000 group-hover:bg-white/10"></div>
        <div class="bg-black/20 backdrop-blur-xl border border-white/10 p-8 md:p-14 rounded-[2.5rem] shadow-2xl shadow-black/20 text-center min-h-[300px] flex flex-col justify-center transition-all duration-500 hover:scale-[1.01]">
            
            <transition name="fade-scale" mode="out-in">
                <div :key="(currentStep.id || currentStep.title || currentStep.content) + (showLatin ? '-la' : '-en')" class="space-y-6">
                    <h3 v-if="currentStep.title" class="text-rose-200/50 uppercase tracking-[0.3em] text-xs font-bold mb-4">
                        {{ currentStep.title }}
                    </h3>
                    
                    <InterlinearText 
                        v-if="interlinearWords" 
                        :words="interlinearWords" 
                    />

                    <p 
                        v-else
                        class="text-2xl md:text-3xl font-serif leading-relaxed text-stone-100 drop-shadow-sm whitespace-pre-line transition-all duration-500"
                        :class="{ 'italic text-rose-100/90': showLatin && currentStep.latin }"
                    >
                        {{ (showLatin && currentStep.latin) ? currentStep.latin : currentStep.content }}
                    </p>

                    <!-- Latin Translation Reference (Hidden if interlinear is active) -->
                    <p 
                        v-if="showLatin && currentStep.latin && !interlinearWords"
                        class="mt-6 text-sm md:text-base font-sans leading-relaxed text-stone-400 max-w-xl mx-auto opacity-60 italic border-t border-white/5 pt-6"
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

const props = defineProps<{
    currentStep: Step;
    showLatin?: boolean;
}>();

const interlinearWords = computed(() => {
    if (!props.showLatin || !props.currentStep.latin) return null;
    return getInterlinearPairs(props.currentStep.latin);
});
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.5s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
  filter: blur(4px);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(1.05);
  filter: blur(4px);
}
</style>
