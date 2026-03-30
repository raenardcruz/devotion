<template>
    <div class="relative group max-w-2xl mx-auto w-full">
        <div class="absolute inset-0 bg-white/5 blur-xl rounded-[2.5rem] -z-10 transition-all duration-1000 group-hover:bg-white/10"></div>
        <div class="bg-black/20 backdrop-blur-xl border border-white/10 p-8 md:p-14 rounded-[2.5rem] shadow-2xl shadow-black/20 text-center min-h-[350px] flex flex-col justify-center transition-all duration-500 hover:scale-[1.01]">
            
            <transition name="fade-scale" mode="out-in">
                <div :key="(currentStep.title || currentStep.content) + (showLatin ? '-la' : '-en')" class="space-y-6">
                    <h3 v-if="currentStep.title" class="text-amber-200/50 uppercase tracking-[0.3em] text-xs font-bold mb-4">
                        {{ currentStep.title }}
                        <span v-if="currentStep.beadNumber && currentStep.type === 'decade-bead'" class="ml-2 text-amber-200/30">
                            {{ currentStep.beadNumber }}
                        </span>
                    </h3>
                    
                    <!-- Prayer Text -->
                    <InterlinearText 
                        v-if="interlinearWords" 
                        :words="interlinearWords" 
                        latin-class="text-amber-100/90"
                    />

                    <p 
                        v-else-if="currentStep.content"
                        class="text-2xl md:text-3xl font-serif leading-relaxed text-stone-100 drop-shadow-sm whitespace-pre-line transition-all duration-500"
                        :class="{ 'italic text-amber-100/90': showLatin && currentStep.latin }"
                    >
                        {{ (showLatin && currentStep.latin) ? currentStep.latin : currentStep.content }}
                    </p>
 
                    <!-- Latin Translation Reference (Hidden if interlinear is active) -->
                    <p 
                        v-if="showLatin && currentStep.latin && !interlinearWords"
                        class="mt-4 text-sm font-sans text-stone-400 opacity-60 italic"
                    >
                        {{ currentStep.content }}
                    </p>

                    <!-- Scripture Verse (Always English) -->
                    <div v-if="currentStep.verse" class="mt-8 pt-8 border-t border-white/5 space-y-4">
                        <p class="text-lg md:text-xl font-serif italic text-amber-100/60 leading-relaxed italic">
                            "{{ currentStep.verse.text }}"
                        </p>
                        <footer class="text-amber-200/40 font-medium tracking-[0.15em] text-xs uppercase">
                            — {{ currentStep.verse.ref }}
                        </footer>
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
