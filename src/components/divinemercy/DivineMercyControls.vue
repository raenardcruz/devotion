<template>
    <div class="flex items-center justify-between mt-auto pt-6 max-w-2xl mx-auto w-full">
        <!-- Previous Button -->
        <AppButton 
            variant="outlined" 
            :disabled="currentStepIndex === 0"
            @click="$emit('prev')"
            custom-class="!px-5 !py-2"
        >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m15 18-6-6 6-6" />
            </svg>
            <span class="font-bold uppercase tracking-wider text-[10px]">Previous</span>
        </AppButton>

        <!-- Progress/Decade Status -->
        <div v-if="currentStep.type === 'decade'" class="flex flex-col items-center">
            <span class="text-parchment-secondary font-bold tracking-[0.2em] text-[9px] uppercase">Decade</span>
            <span class="text-parchment-neutral text-xs md:text-sm font-bold tabular-nums tracking-widest mt-0.5">
                {{ decadeIndex }} <span class="text-parchment-neutral/30 mx-1">/</span> {{ isShort ? '1' : '5' }}
            </span>
        </div>
        <div v-else class="flex flex-col items-center opacity-40">
            <span class="text-parchment-neutral/40 font-bold tracking-[0.2em] text-[9px] uppercase">Preparation</span>
            <div class="flex gap-1 h-3 items-center mt-1">
                <div class="w-1 h-1 rounded-full bg-parchment-neutral/50"></div>
                <div class="w-1.5 h-1.5 rounded-full bg-parchment-neutral/80"></div>
                <div class="w-1 h-1 rounded-full bg-parchment-neutral/50"></div>
            </div>
        </div>

        <!-- Next / Finish Button -->
        <AppButton 
            variant="primary" 
            @click="$emit('next')"
            custom-class="!px-5 !py-2 !bg-parchment-primary-dark hover:!bg-[#4a3c1d]"
        >
            <span class="font-bold uppercase tracking-wider text-[10px]">
                {{ currentStepIndex === totalSteps - 1 ? 'Finish' : 'Next' }}
            </span>
            <svg v-if="currentStepIndex === totalSteps - 1" xmlns="http://www.w3.org/2000/svg"
                width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
                <path d="M3 3v5h5" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"
                fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"
                stroke-linejoin="round">
                <path d="m9 18 6-6-6-6" />
            </svg>
        </AppButton>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';
import AppButton from '../common/AppButton.vue';

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
