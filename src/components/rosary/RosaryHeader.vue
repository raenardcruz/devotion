<template>
    <header class="mb-8 text-center space-y-4 pt-8">
        <div class="flex items-center justify-center space-x-2 text-amber-700 mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path
                    d="M11 2a2 2 0 0 0-2 2v5H4a2 2 0 0 0-2 2v2c0 1.1.9 2 2 2h5v5c0 1.1.9 2 2 2h2a2 2 0 0 0 2-2v-5h5a2 2 0 0 0 2-2v-2a2 2 0 0 0-2-2h-5V4a2 2 0 0 0-2-2h-2z" />
            </svg>
            <h1 class="text-sm font-bold uppercase tracking-widest">The Scriptural Rosary</h1>
        </div>

        <div class="flex flex-wrap justify-center gap-2">
            <button v-for="setName in sets" :key="setName" @click="$emit('changeSet', setName)" :class="[
                    'px-4 py-1.5 rounded-full text-sm font-medium transition-all duration-300',
                    currentSetName === setName 
                        ? 'bg-amber-700 text-white shadow-md' 
                        : 'bg-white/60 hover:bg-white text-stone-600 border border-stone-200 shadow-sm'
                ]">
                {{ setName }}
                <span v-if="setName === todayMystery" class="ml-1 text-[10px] opacity-70">(Today)</span>
            </button>
        </div>
    </header>
</template>

<script setup lang="ts">
defineProps<{
    sets: string[];
    currentSetName: string;
    todayMystery: string;
}>();

defineEmits<{
    (e: 'changeSet', setName: string): void;
}>();
</script>
