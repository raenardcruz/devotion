<template>
    <div class="relative group max-w-2xl mx-auto w-full">
        <div class="absolute inset-0 bg-white/30 blur-2xl rounded-[3rem] -z-10 transition-all duration-1000 group-hover:bg-white/40"></div>
        <div
            class="bg-white/70 backdrop-blur-xl border border-white/60 p-10 md:p-16 rounded-[3rem] shadow-2xl text-center min-h-[320px] flex flex-col justify-center transition-all duration-500 hover:shadow-amber-900/5">
            
            <div class="mb-8 flex justify-center text-amber-800/20">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24"
                    fill="currentColor" class="transform transition-transform duration-700 group-hover:scale-110">
                    <path d="M14.017 21L14.017 18C14.017 16.8954 13.1216 16 12.017 16H6.01723C4.91266 16 4.01723 16.8954 4.01723 18V21H2.01723V18C2.01723 15.7909 3.8081 14 6.01723 14H12.017C14.2261 14 16.017 15.7909 16.017 18V21H14.017ZM20.0172 10H14.017C11.8081 10 10.0172 8.20914 10.0172 6V3H2.01723V21H8.01723V19H4.01723V5H10.0172V6C10.0172 8.20914 11.8081 10 14.017 10ZM12.017 5V6C12.017 4.89543 12.9124 4 14.017 4H18.0172C19.1218 4 20.0172 4.89543 20.0172 6V10H14.017C12.9124 10 12.017 9.10457 12.017 8V5Z"></path>
                </svg>
            </div>

            <transition name="fade-scale" mode="out-in">
                <blockquote :key="currentVerse.text" class="space-y-6">
                    <p class="text-2xl md:text-3xl font-serif leading-relaxed text-stone-800 drop-shadow-sm">
                        "{{ currentVerse.text }}"
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
