<template>
    <div class="text-center mb-6">
        <span class="text-rose-300/80 font-semibold text-xs uppercase tracking-widest">
            {{ currentStep.category }}
        </span>
        <h2 class="text-2xl md:text-4xl font-serif text-stone-100 italic mt-2">
            {{ currentStep.title }}
        </h2>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';

defineProps<{
    currentStep: Step;
}>();
</script>
