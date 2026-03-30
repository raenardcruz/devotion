<template>
    <div class="flex items-center justify-between mt-auto pt-8">
        <button 
            @click="$emit('prev')" 
            :disabled="currentStepIndex === 0"
            class="group flex items-center space-x-3 text-stone-400 hover:text-amber-200 transition-all disabled:opacity-20 disabled:pointer-events-none cursor-pointer"
        >
            <div class="p-3 rounded-full bg-white/5 border border-white/10 shadow-xl group-hover:bg-amber-500/10 group-hover:border-amber-500/30 transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m15 18-6-6 6-6" />
                </svg>
            </div>
            <span class="hidden sm:inline font-bold uppercase tracking-widest text-[10px]">Previous</span>
        </button>

        <div class="flex flex-col items-center gap-1">
            <div class="text-amber-200/40 text-[10px] font-black uppercase tracking-[0.2em]">Progress</div>
            <div class="text-stone-100 text-sm font-bold tabular-nums tracking-widest">
                {{ currentStepIndex + 1 }} <span class="text-stone-600 mx-1">/</span> {{ totalSteps }}
            </div>
        </div>

        <button 
            @click="$emit('next')" 
            class="group flex items-center space-x-3 text-stone-400 hover:text-amber-200 transition-all cursor-pointer"
        >
            <span class="hidden sm:inline font-bold uppercase tracking-widest text-[10px]">Next</span>
            <div class="p-3 rounded-full bg-amber-500/10 border border-amber-500/30 shadow-xl shadow-amber-900/20 group-hover:bg-amber-500/20 group-active:scale-95 transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"
                    fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m9 18 6-6-6-6" />
                </svg>
            </div>
        </button>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    currentStepIndex: number;
    totalSteps: number;
}>();

defineEmits<{
    (e: 'next'): void;
    (e: 'prev'): void;
}>();
</script>
