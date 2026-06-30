<template>
    <div v-if="shouldShow" class="flex flex-col items-center gap-2 mb-2 md:mb-4">
        <div class="flex justify-center items-center gap-3 h-10">
            <button v-for="i in maxBeads" :key="i" 
               @click="$emit('selectBead', i)"
               :class="[
                    'w-6 h-6 rounded-full border-2 border-parchment-primary bg-parchment-bg flex items-center justify-center cursor-pointer transition-all duration-300 hover:scale-115 outline-none shadow-none',
                    i === beadIndex ? 'ring-2 ring-parchment-primary/30 ring-offset-2' : ''
               ]"
            >
                <div v-if="i === beadIndex" class="w-2.5 h-2.5 rounded-full bg-parchment-primary-dark"></div>
            </button>
        </div>
        <span class="text-[10px] text-parchment-neutral/40 font-bold uppercase tracking-[0.2em] select-none">
            Tap a bead to mark your meditation
        </span>
    </div>
    <div v-else class="h-10 mb-2 md:mb-4"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
    beadIndex: number;
    stepType?: string;
}>();

defineEmits<{
    (e: 'selectBead', index: number): void;
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
