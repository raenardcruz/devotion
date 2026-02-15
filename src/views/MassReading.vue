<template>
  <div class="min-h-screen bg-stone-900 text-stone-100 selection:bg-amber-500/30 flex flex-col relative font-sans overflow-hidden">
    <!-- Back Button -->
    <div class="absolute top-6 left-6 z-50">
      <router-link to="/" class="group flex items-center space-x-3 text-white/70 hover:text-white bg-black/20 hover:bg-black/40 backdrop-blur-md px-5 py-2.5 rounded-full transition-all duration-300 border border-white/10 hover:border-white/30">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform">
          <path d="m15 18-6-6 6-6"/>
        </svg>
        <span class="font-medium tracking-wide text-sm">Return Home</span>
      </router-link>
    </div>

    <!-- Background Image -->
    <div class="fixed inset-0 z-0 transition-all duration-1000 ease-in-out" :style="{
            backgroundImage: `url('https://images.unsplash.com/photo-1507692049790-de58293a469d?q=80&w=2070&auto=format&fit=crop')`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            filter: 'brightness(0.4) blur(10px)'
        }"></div>
    
    <!-- Content Container -->
    <div class="relative z-10 max-w-4xl mx-auto px-4 py-20 md:py-24 flex flex-col w-full">
      
      <!-- Header -->
      <header class="text-center mb-12 animate-fade-in-down">
        <span class="text-amber-200 text-xs font-bold uppercase tracking-[0.2em] mb-4 block">Liturgy of the Word</span>
        <h1 class="text-4xl md:text-5xl font-serif text-white mb-2">Daily Mass Readings</h1>
        <p class="text-stone-400 text-lg">{{ formattedDate }}</p>
      </header>

      <!-- Loading State -->
      <div v-if="loading" class="flex flex-col items-center justify-center py-20 animate-pulse">
        <div class="w-12 h-12 border-4 border-amber-200/30 border-t-amber-200 rounded-full animate-spin mb-4"></div>
        <p class="text-stone-400">Loading readings...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="bg-red-900/20 border border-red-500/30 rounded-2xl p-8 text-center animate-fade-in-up">
        <p class="text-red-200 mb-4">Unable to load readings at this time.</p>
        <button @click="fetchReadings" class="px-6 py-2 bg-red-800/40 hover:bg-red-800/60 text-white rounded-full transition-colors border border-red-500/30">
          Try Again
        </button>
        <p class="text-xs text-red-400 mt-4 font-mono">{{ error }}</p>
      </div>

      <!-- Readings Content -->
      <div v-else-if="readings" class="space-y-12 animate-fade-in-up">
        
        <!-- First Reading -->
        <section class="bg-stone-900/60 backdrop-blur-md rounded-3xl p-8 md:p-10 shadow-xl border border-white/5">
          <h2 class="text-2xl font-serif text-amber-200 mb-1">First Reading</h2>
          <p class="text-stone-400 text-sm font-medium uppercase tracking-wider mb-6">{{ readings.first_reading.verse }}</p>
          
          <div class="prose prose-invert prose-lg max-w-none text-stone-200 leading-relaxed mb-8">
            {{ readings.first_reading.text }}
          </div>
          
          <div class="bg-stone-800/50 rounded-xl p-6 border-l-4 border-amber-600/50">
            <h3 class="text-amber-500 text-xs font-bold uppercase tracking-wider mb-2">Context</h3>
            <p class="text-stone-400 text-sm leading-relaxed">{{ readings.first_reading.context }}</p>
          </div>
        </section>

        <!-- Responsorial Psalm -->
        <section class="bg-stone-900/60 backdrop-blur-md rounded-3xl p-8 md:p-10 shadow-xl border border-white/5">
          <h2 class="text-2xl font-serif text-amber-200 mb-1">Responsorial Psalm</h2>
          <p class="text-stone-400 text-sm font-medium uppercase tracking-wider mb-6">{{ readings.responsorial_psalm.verse }}</p>
          
          <div class="prose prose-invert prose-lg max-w-none text-stone-200 leading-relaxed italic text-center">
            {{ readings.responsorial_psalm.text }}
          </div>
        </section>

        <!-- Second Reading (Optional) -->
        <section v-if="readings.second_reading" class="bg-stone-900/60 backdrop-blur-md rounded-3xl p-8 md:p-10 shadow-xl border border-white/5">
          <h2 class="text-2xl font-serif text-amber-200 mb-1">Second Reading</h2>
          <p class="text-stone-400 text-sm font-medium uppercase tracking-wider mb-6">{{ readings.second_reading.verse }}</p>
          
          <div class="prose prose-invert prose-lg max-w-none text-stone-200 leading-relaxed mb-8">
            {{ readings.second_reading.text }}
          </div>
          
          <div class="bg-stone-800/50 rounded-xl p-6 border-l-4 border-amber-600/50">
            <h3 class="text-amber-500 text-xs font-bold uppercase tracking-wider mb-2">Context</h3>
            <p class="text-stone-400 text-sm leading-relaxed">{{ readings.second_reading.context }}</p>
          </div>
        </section>

        <!-- Gospel -->
        <section class="bg-stone-900/60 backdrop-blur-md rounded-3xl p-8 md:p-10 shadow-xl border border-white/5 ring-1 ring-amber-500/20">
          <div class="flex items-center justify-between mb-1">
            <h2 class="text-3xl font-serif text-amber-400">Gospel</h2>
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-amber-500/50">
              <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>
            </svg>
          </div>
          <p class="text-stone-400 text-sm font-medium uppercase tracking-wider mb-6">{{ readings.gospel.verse }}</p>
          
          <div class="prose prose-invert prose-xl max-w-none text-stone-100 leading-relaxed mb-8 font-serif">
            {{ readings.gospel.text }}
          </div>
          
          <div class="bg-stone-800/50 rounded-xl p-6 border-l-4 border-amber-500">
            <h3 class="text-amber-500 text-xs font-bold uppercase tracking-wider mb-2">Values Text / Context</h3>
            <p class="text-stone-400 text-sm leading-relaxed">{{ readings.gospel.context }}</p>
          </div>
        </section>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

interface Reading {
  verse: string;
  text: string;
  context: string;
}

interface Psalm {
    verse: string;
    text: string;
}

interface MassReadingsResponse {
  first_reading: Reading;
  responsorial_psalm: Psalm;
  second_reading?: Reading | null;
  gospel: Reading;
  error?: string;
}

const readings = ref<MassReadingsResponse | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

const formattedDate = new Date().toLocaleDateString('en-US', {
  weekday: 'long',
  year: 'numeric',
  month: 'long',
  day: 'numeric'
});

const fetchReadings = async () => {
  loading.value = true;
  error.value = null;

  const todayKey = new Date().toISOString().split('T')[0];
  const storageKey = `mass-readings-${todayKey}`;

  // Check localStorage for today's readings
  const cachedData = localStorage.getItem(storageKey);
  if (cachedData) {
    readings.value = JSON.parse(cachedData);
    loading.value = false;
    return;
  }

  // Clear old data from localStorage
  const keysToRemove: string[] = [];
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i);
    if (key && key.startsWith('mass-readings-')) {
      keysToRemove.push(key);
    }
  }
  keysToRemove.forEach(key => localStorage.removeItem(key));


  try {
    const response = await fetch('/mass-readings');
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    if (data.error) {
        throw new Error(data.error);
    }
    
    // Save new data to localStorage
    localStorage.setItem(storageKey, JSON.stringify(data));
    readings.value = data;
  } catch (err: any) {
    console.error('Error fetching readings:', err);
    error.value = err.message || 'Failed to load readings.';
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchReadings();
});
</script>

<style scoped>
/* Scoped styles if needed, mostly using Tailwind */
</style>
