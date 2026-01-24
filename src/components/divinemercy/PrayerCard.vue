<template>
    <div class="relative group max-w-2xl mx-auto">
        <div class="absolute inset-0 bg-white/40 blur-2xl rounded-[2rem] -z-10"></div>
        <div
            class="bg-white/90 backdrop-blur-lg border border-white p-8 md:p-12 rounded-[2rem] shadow-xl text-center min-h-[320px] flex flex-col justify-center transition-all duration-500">

            <div class="mb-6 flex justify-center opacity-20 text-rose-900">
                <svg v-if="currentStep.type === 'intro'" xmlns="http://www.w3.org/2000/svg" width="48" height="48"
                    viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round"
                    stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10" />
                    <path d="M12 8v8" />
                    <path d="M8 12h8" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
                    <path
                        d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" />
                </svg>
            </div>

            <div class="space-y-6">
                <p class="text-xl md:text-2xl font-serif leading-relaxed text-stone-800 whitespace-pre-line">
                    {{ currentStep.content }}
                </p>

                <div v-if="currentStep.instruction"
                    class="text-rose-800 text-xs font-bold uppercase tracking-widest pt-4">
                    {{ currentStep.instruction }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';

defineProps<{
    currentStep: Step;
}>();
</script>
