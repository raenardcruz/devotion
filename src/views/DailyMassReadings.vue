<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useDate } from '../composables/useDate';

interface Reading {
  citation: string;
  text: string;
  context?: string;
}

interface ResponsorialPsalm {
  citation: string;
  text: string;
  context?: string;
}

interface MassReadings {
  first_reading: Reading;
  responsorial_psalm: ResponsorialPsalm;
  second_reading: Reading;
  gospel: Reading;
  pope_quote: string;
}

const readings = ref<MassReadings | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);

const { getLocalISOString } = useDate();

const fetchReadings = async () => {
  try {
    const response = await fetch('https://devotionapi.raenardcruz.com/devotion?date=' + getLocalISOString());
    if (!response.ok) {
        // Construct a more informative error message
        const message = `Failed to fetch readings: ${response.status} ${response.statusText}`;
        throw new Error(message);
    }
    const data = await response.json();
    readings.value = data;
  } catch (err: any) {
    error.value = err.message || 'Failed to load readings.';
    console.error("Error fetching readings:", err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchReadings();
});
</script>

<template>
  <div class="min-h-screen bg-stone-50 text-stone-900 font-sans selection:bg-amber-100 flex flex-col">
    <!-- Header -->
    <header class="pt-12 pb-8 text-center px-4">
        <router-link to="/" class="inline-block mb-4 text-stone-400 hover:text-stone-600 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
        </router-link>
      <h1 class="text-4xl md:text-5xl font-serif text-stone-800 tracking-tight mb-2">
        Daily Mass Readings
      </h1>
       <p class="text-stone-500 text-sm uppercase tracking-[0.2em] font-medium">
        Liturgy of the Word
      </p>
    </header>

    <main class="flex-grow container mx-auto px-4 max-w-4xl pb-20">
      <div v-if="loading" class="flex flex-col items-center justify-center py-20 space-y-4">
         <div class="w-12 h-12 border-4 border-stone-200 border-t-amber-600 rounded-full animate-spin"></div>
         <p class="text-stone-500 animate-pulse">Loading readings...</p>
      </div>

      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-8 text-center text-red-800">
        <p class="font-medium mb-2">Recall error</p>
        <p class="text-sm opacity-80">{{ error }}</p>
        <button @click="fetchReadings" class="mt-4 px-4 py-2 bg-red-100 hover:bg-red-200 text-red-800 rounded-lg transition-colors text-sm font-medium">
            Try Again
        </button>
      </div>

      <div v-else-if="readings" class="space-y-12 animate-fade-in-up">
        
        <!-- First Reading -->
        <section class="bg-white rounded-[2rem] shadow-sm p-8 md:p-12" v-if="!!readings.first_reading">
            <header class="mb-8 border-b border-stone-100 pb-6">
                <span class="text-amber-600 text-xs font-bold uppercase tracking-[0.2em] mb-2 block">First Reading</span>
                <h2 class="text-3xl font-serif text-stone-800">{{ readings.first_reading.citation }}</h2>
            </header>
            
            <div class="prose prose-stone max-w-none mb-8">
                 <div class="whitespace-pre-line text-lg leading-relaxed text-stone-700" v-html="readings.first_reading.text"></div>
            </div>

             <div class="bg-stone-50 rounded-xl p-6 border-l-4 border-amber-200">
                <h3 class="text-stone-500 text-xs font-bold uppercase tracking-wider mb-2">Context</h3>
                <p class="text-stone-600 italic text-sm leading-relaxed">
                    {{ readings.first_reading.context }}
                </p>
            </div>
        </section>

        <!-- Responsorial Psalm -->
         <section class="bg-stone-100 rounded-[2rem] p-8 md:p-12 shadow-inner" v-if="!!readings.responsorial_psalm">
            <header class="mb-6 text-center">
                <span class="text-stone-500 text-xs font-bold uppercase tracking-[0.2em] mb-2 block">Responsorial Psalm</span>
                <h2 class="text-2xl font-serif text-stone-800">{{ readings.responsorial_psalm.citation }}</h2>
            </header>
            
            <div class="prose prose-stone max-w-none text-center">
                 <div class="whitespace-pre-line text-xl font-serif italic text-stone-700 leading-loose" v-html="readings.responsorial_psalm.text"></div>
            </div>
            <div class="bg-stone-50 rounded-xl p-6 border-l-4 border-amber-200">
                <h3 class="text-stone-500 text-xs font-bold uppercase tracking-wider mb-2">Context</h3>
                <p class="text-stone-600 italic text-sm leading-relaxed">
                    {{ readings.responsorial_psalm.context }}
                </p>
            </div>
        </section>

        <!-- Second reading -->
         <section class="bg-white rounded-[2rem] shadow-sm p-8 md:p-12" v-if="!!readings.second_reading">
            <header class="mb-8 border-b border-stone-100 pb-6">
                <span class="text-amber-600 text-xs font-bold uppercase tracking-[0.2em] mb-2 block">Second Reading</span>
                <h2 class="text-3xl font-serif text-stone-800">{{ readings.second_reading.citation }}</h2>
            </header>
            
            <div class="prose prose-stone max-w-none mb-8">
                 <div class="whitespace-pre-line text-lg leading-relaxed text-stone-700" v-html="readings.second_reading.text"></div>
            </div>

             <div class="bg-stone-50 rounded-xl p-6 border-l-4 border-amber-200">
                <h3 class="text-stone-500 text-xs font-bold uppercase tracking-wider mb-2">Context</h3>
                <p class="text-stone-600 italic text-sm leading-relaxed">
                    {{ readings.second_reading.context }}
                </p>
            </div>
        </section>

        <!-- Gospel -->
        <section class="bg-white rounded-[2rem] shadow-xl p-8 md:p-12 relative overflow-hidden border border-stone-100" v-if="!!readings.gospel">
             <div class="absolute top-0 right-0 p-8 opacity-5">
                <svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 24 24" fill="currentColor" class="text-amber-900"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-1.07 3.97-2.9 5.34z"></path></svg>
            </div>

            <header class="mb-8 border-b border-stone-100 pb-6 relative z-10">
                <span class="text-amber-600 text-xs font-bold uppercase tracking-[0.2em] mb-2 block">Gospel</span>
                <h2 class="text-3xl font-serif text-stone-800">{{ readings.gospel.citation }}</h2>
            </header>
            
            <div class="prose prose-stone max-w-none mb-8 relative z-10">
                 <div class="whitespace-pre-line text-lg leading-relaxed text-stone-800" v-html="readings.gospel.text"></div>
            </div>

             <div class="bg-amber-50 rounded-xl p-6 border-l-4 border-amber-600 relative z-10">
                <h3 class="text-amber-800/60 text-xs font-bold uppercase tracking-wider mb-2">Context</h3>
                <p class="text-stone-700 italic text-sm leading-relaxed">
                    {{ readings.gospel.context }}
                </p>
            </div>
        </section>

        <!-- Words of the Popes -->
        <section class="bg-white rounded-[2rem] shadow-xl p-8 md:p-12 relative overflow-hidden border border-stone-100" v-if="!!readings.pope_quote">
             <div class="absolute top-0 right-0 p-8 opacity-5">
                <svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 24 24" fill="currentColor" class="text-amber-900"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-1.07 3.97-2.9 5.34z"></path></svg>
            </div>

            <header class="mb-8 border-b border-stone-100 pb-6 relative z-10">
                <span class="text-amber-600 text-xs font-bold uppercase tracking-[0.2em] mb-2 block">Words of the Popes</span>
            </header>

            <div class="prose prose-stone max-w-none mb-8 relative z-10">
                 <div class="whitespace-pre-line text-lg leading-relaxed text-stone-800" v-html="readings.pope_quote"></div>
            </div>

            <div class="text-center text-stone-400 text-xs tracking-[0.2em] mb-2">https://www.vaticannews.va/en/word-of-the-day</div>
        </section>

        <section>
          <div class="text-center text-stone-400 text-xs tracking-[0.2em] mb-2">Holy Bible, New Living Translation Catholic Edition, copyright © 2016 by Tyndale House Foundation. All rights reserved. Used by permission of Tyndale House Publishers, Carol Stream, Illinois 60188. All rights reserved.</div>
        </section>
      </div>
    </main>
  </div>
</template>

<style scoped>
.animate-fade-in-up {
  animation: fadeInUp 0.8s ease-out forwards;
}
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

<style>
.v {
  font-size: 12px;
  color: grey;
  margin-right: 4px;
}
</style>
