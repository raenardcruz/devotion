<template>
    <div class="text-center mb-6 space-y-2">
        <span class="text-amber-700 font-semibold text-xs uppercase tracking-tighter">
            Mystery {{ mysteryIndex + 1 }} of 5
        </span>
        <h2 class="text-2xl md:text-4xl font-serif text-stone-800 italic">
            {{ currentMystery.title }}
        </h2>
        <p class="text-stone-500 text-sm md:text-base max-w-lg mx-auto leading-relaxed">
            {{ currentMystery.focus }}
        </p>
    </div>
</template>

<script setup lang="ts">
import type { Mystery } from './rosaryData';

defineProps<{
    mysteryIndex: number;
    currentMystery: Mystery;
}>();
</script>
