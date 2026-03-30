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
        <p class="text-stone-400 text-sm md:text-base max-w-lg mx-auto leading-relaxed italic opacity-80">
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
