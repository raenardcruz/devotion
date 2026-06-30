<template>
    <div v-if="currentStep.type === 'decade'"
        class="flex justify-center items-center gap-3 mb-4 flex-wrap max-w-lg mx-auto">
        
        <!-- Large Bead for Eternal Father (Crimson accent) -->
        <div :class="[
                'w-6 h-6 rounded-full border-2 border-parchment-secondary bg-parchment-bg flex items-center justify-center transition-all duration-300',
                beadInDecade >= 0 ? 'ring-2 ring-parchment-secondary/20 ring-offset-2' : ''
            ]">
            <div v-if="beadInDecade >= 0" class="w-2.5 h-2.5 rounded-full bg-parchment-secondary"></div>
        </div>
        
        <!-- Divider Line -->
        <div class="w-6 h-0.5 bg-parchment-border"></div>
        
        <!-- Small Beads for Passion (Slate Blue accent) -->
        <div v-for="i in 10" :key="i" :class="[
                'w-4 h-4 rounded-full border-2 border-parchment-tertiary bg-parchment-bg flex items-center justify-center transition-all duration-300',
                i === beadInDecade ? 'ring-2 ring-parchment-tertiary/30 ring-offset-2' : ''
            ]">
             <div v-if="i <= beadInDecade" class="w-1.5 h-1.5 rounded-full bg-parchment-tertiary"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Step } from './divineMercyData';

defineProps<{
    currentStep: Step;
    beadInDecade: number;
}>();
</script>
