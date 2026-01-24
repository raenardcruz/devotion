<template>
    <div v-if="currentStep.type === 'decade'"
        class="flex justify-center items-center gap-2 mb-10 flex-wrap max-w-lg mx-auto perspective-1000">
        
        <!-- Large Bead for Eternal Father -->
        <div :class="[
                'w-6 h-6 rounded-full border border-white/20 transition-all duration-700 relative z-10',
                beadInDecade >= 0 
                    ? 'bg-gradient-to-br from-rose-700 to-rose-950 shadow-lg shadow-rose-900/50 scale-110' 
                    : 'bg-gradient-to-br from-stone-600 to-stone-800 shadow-inner'
            ]">
            <div v-if="beadInDecade >= 0" class="absolute inset-0 bg-rose-500 rounded-full blur-sm opacity-50 animate-pulse"></div>
        </div>
        
        <div class="w-8 h-px bg-white/20"></div>
        
        <!-- Small Beads for Passion -->
        <div v-for="i in 10" :key="i" :class="[
                'w-3.5 h-3.5 rounded-full transition-all duration-500 relative',
                i <= beadInDecade 
                    ? 'bg-gradient-to-br from-blue-500 to-blue-700 shadow-[0_0_10px_rgba(59,130,246,0.5)] transform scale-110' 
                    : 'bg-stone-700 opacity-60',
                i === beadInDecade ? 'ring-2 ring-blue-300 ring-offset-2 ring-offset-transparent animate-pulse-fast' : ''
            ]">
             <!-- Ray effects -->
             <div v-if="i <= beadInDecade" :class="[
                'absolute -inset-1 blur-md rounded-full opacity-40',
                i % 2 === 0 ? 'bg-rose-400' : 'bg-blue-300' 
             ]"></div>
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
