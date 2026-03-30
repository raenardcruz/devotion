<template>
    <div class="text-center mb-6 flex flex-col items-center gap-1">
        <div class="flex items-center gap-2">
            <span class="text-rose-300/80 font-bold text-[10px] uppercase tracking-[0.3em]">
                {{ currentStep.category }}
            </span>
            <span v-if="isShort && currentStep.type === 'decade'" class="px-2 py-0.5 bg-rose-500/20 text-rose-300 text-[8px] font-black uppercase tracking-widest rounded-full border border-rose-500/30">
                Short Version
            </span>
        </div>
        <h2 class="text-3xl md:text-5xl font-serif text-stone-100 mt-2 tracking-tight drop-shadow-md">
            {{ currentStep.title }}
        </h2>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';

defineProps<{
    currentStep: Step;
    isShort?: boolean;
}>();
</script>
