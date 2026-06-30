<template>
    <div class="text-center mb-4 flex flex-col items-center gap-1">
        <div class="flex items-center gap-2">
            <span class="text-parchment-secondary font-bold text-[10px] uppercase tracking-[0.25em]">
                {{ currentStep.category }}
            </span>
            <span v-if="isShort && currentStep.type === 'decade'" class="px-2 py-0.5 bg-parchment-secondary/10 text-parchment-secondary text-[8px] font-bold uppercase tracking-wider rounded-full border border-parchment-secondary/20">
                Short Version
            </span>
        </div>
        <h2 class="text-2xl md:text-3xl font-serif text-parchment-neutral mt-1 font-medium">
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
