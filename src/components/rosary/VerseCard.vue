<template>
    <div class="relative group">
        <div class="absolute inset-0 bg-white/40 blur-xl rounded-[2rem] -z-10"></div>
        <div
            class="bg-white/80 backdrop-blur-md border border-white p-8 md:p-12 rounded-[2rem] shadow-xl text-center min-h-[300px] flex flex-col justify-center">
            <div class="mb-6 flex justify-center text-amber-600/30">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
                </svg>
            </div>

            <blockquote class="space-y-4">
                <p class="text-xl md:text-3xl font-serif leading-snug text-stone-800">
                    "{{ currentVerse.text }}"
                </p>
                <footer class="text-amber-700 font-medium tracking-wide">
                    — {{ currentVerse.ref }}
                </footer>
            </blockquote>

            <div class="mt-8 text-stone-400 text-xs font-semibold tracking-widest uppercase">
                Hail Mary {{ beadIndex + 1 }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Verse } from './rosaryData';

defineProps<{
    currentVerse: Verse;
    beadIndex: number;
}>();
</script>
