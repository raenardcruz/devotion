<template>
    <div class="flex items-center justify-between mt-auto pt-8 max-w-2xl mx-auto w-full">
        <button 
            @click="$emit('prev')" 
            :disabled="currentStepIndex === 0"
            class="group flex items-center space-x-3 text-stone-400 hover:text-rose-200 transition-all disabled:opacity-0 disabled:pointer-events-none cursor-pointer"
        >
            <div class="p-3 rounded-full bg-white/5 border border-white/10 shadow-xl group-hover:bg-rose-500/10 group-hover:border-rose-500/30 transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m15 18-6-6 6-6" />
                </svg>
            </div>
            <span class="hidden sm:inline font-bold uppercase tracking-widest text-[10px]">Previous</span>
        </button>

        <div v-if="currentStep.type === 'decade'" class="flex flex-col items-center gap-1">
            <span class="text-rose-300/40 text-[10px] font-black uppercase tracking-[0.2em]">Decade</span>
            <span class="text-stone-100 text-sm font-bold tabular-nums tracking-widest">
                {{ decadeIndex }} <span class="text-stone-600 mx-1">/</span> {{ isShort ? '1' : '5' }}
            </span>
        </div>
        <div v-else class="flex flex-col items-center gap-1 opacity-20">
            <span class="text-rose-300/40 text-[10px] font-black uppercase tracking-[0.2em]">Preparing</span>
            <div class="flex gap-1.5 h-4 items-center">
                <div class="w-1 h-1 rounded-full bg-stone-500"></div>
                <div class="w-1.5 h-1.5 rounded-full bg-stone-400"></div>
                <div class="w-1 h-1 rounded-full bg-stone-500"></div>
            </div>
        </div>

        <button 
            @click="$emit('next')"
            class="group flex items-center space-x-3 text-stone-400 hover:text-rose-200 transition-all cursor-pointer"
        >
            <span class="hidden sm:inline font-bold uppercase tracking-widest text-[10px]">
                {{ currentStepIndex === totalSteps - 1 ? 'Finish' : 'Next' }}
            </span>
            <div class="p-3 rounded-full bg-rose-500/10 border border-rose-500/30 shadow-xl shadow-rose-900/20 group-hover:bg-rose-500/20 group-active:scale-95 transition-all">
                <svg v-if="currentStepIndex === totalSteps - 1" xmlns="http://www.w3.org/2000/svg"
                    width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                    stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                    <path d="M3 3v5h5" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"
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
    isShort?: boolean;
}>();

defineEmits<{
    (e: 'next'): void;
    (e: 'prev'): void;
}>();
</script>
