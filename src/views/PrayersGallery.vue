<template>
  <div class="min-h-screen bg-stone-900 text-stone-100 selection:bg-cyan-500/30 flex flex-col relative font-sans overflow-hidden">
    <!-- Back Button -->
    <div class="absolute top-6 left-6 z-50">
      <router-link to="/" class="group flex items-center space-x-3 text-white/70 hover:text-white bg-black/20 hover:bg-black/40 backdrop-blur-md px-5 py-2.5 rounded-full transition-all duration-300 border border-white/10 hover:border-white/30">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform">
          <path d="m15 18-6-6 6-6"/>
        </svg>
        <span class="font-medium tracking-wide text-sm">Return Home</span>
      </router-link>
    </div>

    <!-- Background -->
    <div class="fixed inset-0 z-0 bg-[url('https://images.unsplash.com/photo-1438232992991-995b7058bbb3?q=80&w=2073&auto=format&fit=crop')] bg-cover bg-center filter brightness-[0.3] contrast-125"></div>
    <div class="fixed inset-0 z-0 bg-gradient-to-t from-stone-900 via-stone-900/50 to-transparent"></div>

    <!-- Content -->
    <div class="relative z-10 max-w-6xl mx-auto px-4 py-24 w-full">
      <header class="mb-12 text-center">
        <h1 class="text-4xl md:text-5xl font-serif text-cyan-100 mb-4 animate-fade-in-down">Prayers Gallery</h1>
        <p class="text-stone-400 max-w-2xl mx-auto animate-fade-in-up delay-100">A collection of traditional Catholic prayers for reflection and devotion.</p>
      </header>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 animate-fade-in-up delay-200">
        <div 
          v-for="prayer in prayers" 
          :key="prayer.id"
          class="group bg-stone-800/40 backdrop-blur-md border border-white/10 rounded-2xl p-6 hover:bg-stone-800/60 transition-all duration-300 hover:scale-[1.02] hover:shadow-xl hover:border-cyan-500/30 cursor-pointer"
          @click="openPrayer(prayer)"
        >
          <h3 class="text-2xl font-serif text-cyan-200 mb-3 group-hover:text-cyan-100 transition-colors">{{ prayer.title }}</h3>
          <p class="text-stone-300 text-sm line-clamp-3 leading-relaxed">{{ prayer.content }}</p>
          <div class="mt-4 flex items-center text-cyan-400/80 text-xs font-bold uppercase tracking-wider group-hover:text-cyan-300">
            <span>Read Prayer</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-2 group-hover:translate-x-1 transition-transform">
              <path d="M5 12h14"></path>
              <path d="m12 5 7 7-7 7"></path>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Prayer Modal -->
    <transition name="modal">
      <div v-if="selectedPrayer" class="fixed inset-0 z-[100] flex items-center justify-center p-4" @click.self="closePrayer">
        <div class="absolute inset-0 bg-black/80 backdrop-blur-sm transition-opacity" @click="closePrayer"></div>
        <div class="relative bg-stone-900 border border-stone-700 rounded-3xl p-8 md:p-12 max-w-2xl w-full max-h-[90vh] overflow-y-auto shadow-2xl transform transition-all">
          <button @click="closePrayer" class="absolute top-4 right-4 text-stone-400 hover:text-white transition-colors p-2 rounded-full hover:bg-white/10">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 6 6 18"></path>
              <path d="m6 6 18 18"></path>
            </svg>
          </button>
          
          <h2 class="text-3xl md:text-4xl font-serif text-cyan-100 mb-6 text-center">{{ selectedPrayer.title }}</h2>
          <div class="prose prose-invert prose-lg mx-auto text-stone-200 leading-loose">
            <p class="whitespace-pre-line">{{ selectedPrayer.content }}</p>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import prayersData from '../components/prayers/prayers.json';

interface Prayer {
  id: string;
  title: string;
  content: string;
}

const prayers = ref<Prayer[]>(prayersData);
const selectedPrayer = ref<Prayer | null>(null);

const openPrayer = (prayer: Prayer) => {
  selectedPrayer.value = prayer;
  document.body.style.overflow = 'hidden';
};

const closePrayer = () => {
  selectedPrayer.value = null;
  document.body.style.overflow = '';
};
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .transform,
.modal-leave-active .transform {
  transition: all 0.3s ease-out;
}

.modal-enter-from .transform,
.modal-leave-to .transform {
  opacity: 0;
  transform: scale(0.95);
}
</style>
