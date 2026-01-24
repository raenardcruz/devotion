<template>
  <div class="min-h-screen bg-stone-50 text-stone-900 font-sans selection:bg-amber-100 flex flex-col">
    <!-- Header -->
    <header class="pt-16 pb-8 text-center px-4">
      <h1 class="text-4xl md:text-6xl font-serif text-stone-800 tracking-tight mb-4">Devotion</h1>
      <p class="text-stone-500 text-sm md:text-base uppercase tracking-widest font-medium">Daily Catholic Prayers & Meditations</p>
    </header>

    <!-- Main Content -->
    <main class="flex-grow flex items-center justify-center p-4">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl w-full">
        
        <!-- Rosary Card -->
        <router-link to="/rosary" class="group relative overflow-hidden rounded-3xl shadow-xl hover:shadow-2xl transition-all duration-500 hover:-translate-y-1 h-80 md:h-96">
          <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1554089638-52a6a0c7b2a1?q=80&w=772&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')] bg-cover bg-center transition-transform duration-700 group-hover:scale-105"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-stone-900/90 via-stone-900/40 to-transparent"></div>
          
          <div class="relative h-full flex flex-col justify-end p-8 text-white">
            <div class="mb-4 transform translate-y-2 group-hover:translate-y-0 transition-transform duration-500">
              <span class="text-amber-200 text-xs font-bold uppercase tracking-widest mb-2 block">Meditation</span>
              <h2 class="text-3xl font-serif italic mb-2">Scriptural Rosary</h2>
              <p class="text-stone-300 text-sm line-clamp-2 opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-100">
                Meditate on the mysteries of Christ's life with scriptural verses for each Hail Mary.
              </p>
            </div>
            
            <div class="flex items-center text-amber-200 text-sm font-medium opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-200">
              <span>Start Prayer</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-2 group-hover:translate-x-1 transition-transform">
                <path d="M5 12h14"></path>
                <path d="m12 5 7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </router-link>

        <!-- Divine Mercy Card -->
        <router-link to="/divine-mercy" class="group relative overflow-hidden rounded-3xl shadow-xl hover:shadow-2xl transition-all duration-500 hover:-translate-y-1 h-80 md:h-96">
          <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1696261803619-9cae57748d90?q=80&w=774&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')] bg-cover bg-center transition-transform duration-700 group-hover:scale-105"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-rose-950/90 via-rose-900/40 to-transparent"></div>
          
          <div class="relative h-full flex flex-col justify-end p-8 text-white">
            <div class="mb-4 transform translate-y-2 group-hover:translate-y-0 transition-transform duration-500">
              <span class="text-rose-200 text-xs font-bold uppercase tracking-widest mb-2 block">Chaplet</span>
              <h2 class="text-3xl font-serif italic mb-2">Divine Mercy</h2>
              <p class="text-stone-300 text-sm line-clamp-2 opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-100">
                Pray the Chaplet of Divine Mercy for grace and mercy upon the whole world.
              </p>
            </div>
            
            <div class="flex items-center text-rose-200 text-sm font-medium opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-200">
              <span>Start Prayer</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-2 group-hover:translate-x-1 transition-transform">
                <path d="M5 12h14"></path>
                <path d="m12 5 7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </router-link>

      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
// No script needed for this simple landing page
</script>

<style scoped>
/* Scoped styles if needed, mostly using Tailwind */
</style>