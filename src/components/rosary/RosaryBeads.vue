<template>
    <div class="flex justify-center items-center gap-2 mb-8">
        <div v-for="i in 10" :key="i" :class="[
                'w-3 h-3 md:w-4 md:h-4 rounded-full transition-all duration-500',
                (i - 1) <= beadIndex 
                    ? 'bg-amber-600 shadow-[0_0_8px_rgba(180,83,9,0.5)]' 
                    : 'bg-stone-200',
                (i - 1) === beadIndex ? 'ring-4 ring-amber-100 animate-pulse' : ''
            ]"></div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    beadIndex: number;
}>();
</script>
