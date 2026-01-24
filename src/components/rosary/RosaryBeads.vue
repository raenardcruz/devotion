<template>
    <div class="flex justify-center items-center gap-3 mb-10 perspective-1000">
        <div v-for="i in 10" :key="i" :class="[
                'rounded-full transition-all duration-700 ease-out relative',
                (i - 1) <= beadIndex 
                    ? 'w-5 h-5 md:w-6 md:h-6 bg-gradient-to-br from-amber-400 to-amber-700 shadow-lg shadow-amber-900/40 transform scale-110' 
                    : 'w-3 h-3 md:w-4 md:h-4 bg-gradient-to-br from-stone-200 to-stone-400 opacity-60',
                (i - 1) === beadIndex ? 'ring-offset-4 ring-offset-transparent ring-2 ring-amber-300 animate-pulse-slow' : ''
            ]">
            <div v-if="(i - 1) === beadIndex" class="absolute inset-0 rounded-full bg-amber-400 animate-ping opacity-20"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    beadIndex: number;
}>();
</script>

<style scoped>
.perspective-1000 {
    perspective: 1000px;
}
.animate-pulse-slow {
    animation: pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
@keyframes pulse {
  0%, 100% {
    opacity: 1;
    box-shadow: 0 0 0 0px rgba(245, 158, 11, 0.7);
  }
  50% {
    opacity: .8;
    box-shadow: 0 0 0 10px rgba(245, 158, 11, 0);
  }
}
</style>
