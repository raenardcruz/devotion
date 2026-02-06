<template>
    <div class="flex items-center justify-between mt-auto pt-8">
        <button @click="$emit('prev')" :disabled="mysteryIndex === 0 && beadIndex === 0"
            class="group flex items-center space-x-2 text-stone-400 transition-colors disabled:opacity-30 disabled:pointer-events-none cursor-pointer">
            <div class="p-3 rounded-full bg-white shadow-sm transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m15 18-6-6 6-6" />
                </svg>
            </div>
            <span class="hidden sm:inline font-medium text-sm">Previous</span>
        </button>

        <div class="text-stone-400 text-sm font-medium tabular-nums">
            {{ beadIndex + 1 }} <span class="text-stone-300 mx-1">/</span> 10
        </div>

        <button @click="$emit('next')" :disabled="mysteryIndex === 4 && beadIndex === 9"
            class="group flex items-center space-x-2 text-stone-400 transition-colors disabled:opacity-30 disabled:pointer-events-none cursor-pointer">
            <span class="hidden sm:inline font-medium text-sm">Next</span>
            <div class="p-3 rounded-full bg-white shadow-sm transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m9 18 6-6-6-6" />
                </svg>
            </div>
        </button>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    mysteryIndex: number;
    beadIndex: number;
}>();

defineEmits<{
    (e: 'next'): void;
    (e: 'prev'): void;
}>();
</script>
