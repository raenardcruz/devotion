<template>
    <div v-if="currentStep.type === 'decade'"
        class="flex justify-center items-center gap-1.5 mb-8 flex-wrap max-w-sm mx-auto">
        <!-- Large Bead for Eternal Father -->
        <div :class="[
                'w-5 h-5 rounded-full border-2 transition-all duration-500',
                beadInDecade >= 0 ? 'bg-rose-700 border-rose-900 shadow-md' : 'bg-white border-stone-200'
            ]"></div>
        <div class="w-4 h-0.5 bg-stone-200"></div>
        <!-- Small Beads for Passion -->
        <div v-for="i in 10" :key="i" :class="[
                'w-3 h-3 md:w-3.5 md:h-3.5 rounded-full transition-all duration-500',
                i <= beadInDecade ? 'bg-blue-600 shadow-sm' : 'bg-stone-200',
                i === beadInDecade ? 'ring-4 ring-blue-100 animate-pulse' : ''
            ]"></div>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';

defineProps<{
    currentStep: Step;
    beadInDecade: number;
}>();
</script>
