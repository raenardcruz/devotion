<template>
    <div class="flex items-center justify-between mt-auto pt-8 max-w-2xl mx-auto w-full">
        <button @click="$emit('prev')" :disabled="currentStepIndex === 0"
            class="group flex items-center space-x-2 text-stone-400 hover:text-rose-800 transition-colors disabled:opacity-0 disabled:pointer-events-none cursor-pointer">
            <div class="p-3 rounded-full bg-white shadow-sm group-hover:shadow-md transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m15 18-6-6 6-6" />
                </svg>
            </div>
            <span class="hidden sm:inline font-medium text-sm">Previous</span>
        </button>

        <div v-if="currentStep.type === 'decade'" class="flex flex-col items-center">
            <span class="text-[10px] text-stone-400 uppercase tracking-widest font-bold">Decade</span>
            <span class="text-lg font-serif italic text-stone-600">{{ decadeIndex }} / 5</span>
        </div>
        <div v-else class="text-stone-300">
            •••
        </div>

        <button @click="$emit('next')"
            class="group flex items-center space-x-2 text-stone-400 hover:text-rose-800 transition-colors cursor-pointer">
            <span class="hidden sm:inline font-medium text-sm">{{ currentStepIndex === totalSteps - 1 ?
                'Start Over' : 'Next' }}</span>
            <div class="p-3 rounded-full bg-rose-800 text-white shadow-lg group-hover:scale-110 transition-all">
                <svg v-if="currentStepIndex === totalSteps - 1" xmlns="http://www.w3.org/2000/svg"
                    width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                    stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                    <path d="M3 3v5h5" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m9 18 6-6-6-6" />
                </svg>
            </div>
        </button>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';

defineProps<{
    currentStepIndex: number;
    currentStep: Step;
    decadeIndex: number;
    totalSteps: number;
}>();

defineEmits<{
    (e: 'next'): void;
    (e: 'prev'): void;
}>();
</script>
