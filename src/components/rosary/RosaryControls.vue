<template>
    <div class="flex items-center justify-between mt-auto pt-6 w-full max-w-2xl mx-auto">
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

        <!-- Progress Indicator -->
        <div class="flex flex-col items-center">
            <span class="text-parchment-neutral/40 text-[9px] font-bold uppercase tracking-[0.2em]">Meditation Progress</span>
            <div class="text-parchment-neutral text-xs md:text-sm font-bold tabular-nums tracking-widest mt-0.5">
                {{ currentStepIndex + 1 }} <span class="text-parchment-neutral/30 mx-1">/</span> {{ totalSteps }}
            </div>
        </div>

        <!-- Next Button -->
        <AppButton 
            variant="primary" 
            @click="$emit('next')"
            custom-class="!px-5 !py-2"
        >
            <span class="font-bold uppercase tracking-wider text-[10px]">Next</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6" />
            </svg>
        </AppButton>
    </div>
</template>

<script setup lang="ts">
import AppButton from '../common/AppButton.vue';

defineProps<{
    currentStepIndex: number;
    totalSteps: number;
}>();

defineEmits<{
    (e: 'next'): void;
    (e: 'prev'): void;
}>();
</script>
