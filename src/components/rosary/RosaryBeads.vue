<template>
    <div v-if="shouldShow" class="flex justify-center items-center gap-4 mb-2 md:mb-6 h-10">
        <div v-for="i in maxBeads" :key="i" :class="[
                'rounded-full transition-all duration-700 ease-out relative',
                i <= beadIndex 
                    ? 'w-4 h-4 md:w-5 md:h-5 bg-gradient-to-br from-amber-300 to-amber-600 shadow-lg shadow-amber-900/40 transform scale-110' 
                    : 'w-2 h-2 md:w-3 md:h-3 bg-white/10 border border-white/5',
                i === beadIndex ? 'ring-offset-4 ring-offset-stone-900 ring-2 ring-amber-400/50' : ''
            ]">
            <div v-if="i === beadIndex" class="absolute inset-0 rounded-full bg-amber-400 animate-ping opacity-25"></div>
        </div>
    </div>
    <div v-else class="h-10 mb-2 md:mb-6"></div> <!-- Spacer to prevent layout shift -->
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
    beadIndex: number;
    stepType?: string;
}>();

const shouldShow = computed(() => {
    return props.stepType === 'decade-bead' || 
           (props.stepType === 'intro' && props.beadIndex > 0);
});

const maxBeads = computed(() => {
    if (props.stepType === 'intro') return 3;
    return 10;
});
</script>

<style scoped>
/* Keeping it simple with Tailwind and inline classes for reactivity */
</style>
