<template>
    <div class="relative group mx-auto">
        <div class="absolute inset-0 bg-white/30 blur-2xl rounded-[3rem] -z-10 transition-all duration-1000 group-hover:bg-white/40"></div>
        <div
            class="bg-white/70 backdrop-blur-xl border border-white/60 p-10 md:p-16 rounded-[3rem] shadow-2xl text-center min-h-[320px] flex flex-col justify-center transition-all duration-500 hover:shadow-amber-900/5">

            <transition name="fade-scale" mode="out-in">
                <blockquote :key="currentVerse.text" class="space-y-6">
                    <p class="text-2xl md:text-3xl font-serif leading-relaxed text-stone-800 drop-shadow-sm">
                        {{ currentVerse.text }}
                    </p>
                    <footer class="text-amber-700 font-medium tracking-[0.15em] text-sm uppercase opacity-80">
                        — {{ currentVerse.ref }}
                    </footer>
                </blockquote>
            </transition>

            <div class="mt-10 flex items-center justify-center gap-3 text-stone-400 text-[10px] font-bold tracking-[0.3em] uppercase opacity-60">
                <span class="h-px w-8 bg-stone-300"></span>
                Hail Mary {{ beadIndex + 1 }}
                <span class="h-px w-8 bg-stone-300"></span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Verse } from './rosaryData';

defineProps<{
    currentVerse: Verse;
    beadIndex: number;
}>();
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.5s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
  filter: blur(4px);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(1.05);
  filter: blur(4px);
}
</style>
