<template>
    <div class="relative group max-w-2xl mx-auto w-full">
        <div class="absolute inset-0 bg-white/5 blur-xl rounded-[2.5rem] -z-10 transition-all duration-1000 group-hover:bg-white/10"></div>
        <div class="bg-black/20 backdrop-blur-xl border border-white/10 p-8 md:p-14 rounded-[2.5rem] shadow-2xl shadow-black/20 text-center min-h-[300px] flex flex-col justify-center transition-all duration-500 hover:scale-[1.01]">
            
            <transition name="fade-scale" mode="out-in">
                <div :key="currentStep.label || currentStep.content" class="space-y-6">
                    <h3 v-if="currentStep.label" class="text-rose-200/50 uppercase tracking-[0.3em] text-xs font-bold mb-4">
                        {{ currentStep.label }}
                    </h3>
                    
                    <p class="text-2xl md:text-3xl font-serif leading-relaxed text-stone-100 drop-shadow-sm whitespace-pre-line">
                        {{ currentStep.content }}
                    </p>
                </div>
            </transition>

        </div>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';

defineProps<{
    currentStep: Step;
}>();
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
