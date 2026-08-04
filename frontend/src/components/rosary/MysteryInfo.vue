<template>
    <div class="text-center space-y-3 mb-6">
        <div class="flex flex-col items-center gap-1.5">
            <span class="text-amber-200/40 font-bold text-[10px] uppercase tracking-[0.3em]">
                Mystery {{ mysteryIndex + 1 }} of 5
            </span>
            <div class="h-0.5 w-8 bg-amber-500/20 rounded-full"></div>
        </div>
        <h2 class="text-3xl md:text-5xl font-serif text-stone-100 tracking-tight drop-shadow-md">
            {{ currentMystery.title }}
        </h2>
        <p class="text-stone-400 text-sm md:text-base max-w-lg mx-auto leading-relaxed opacity-80">
            {{ currentMystery.description }}
        </p>
        <div v-if="currentMystery.virtues && currentMystery.virtues.length" class="flex flex-wrap justify-center gap-2 pt-1">
            <span 
                v-for="virtue in currentMystery.virtues" 
                :key="virtue"
                class="px-2.5 py-1 bg-amber-500/10 border border-amber-500/20 text-amber-200 text-xs rounded-full font-semibold"
            >
                {{ virtue }}
            </span>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Mystery } from './rosaryData';

defineProps<{
    mysteryIndex: number;
    currentMystery: Mystery;
}>();
</script>
