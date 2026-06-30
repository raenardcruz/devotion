<template>
  <header class="sticky top-0 z-40 w-full bg-parchment-bg/95 backdrop-blur-md border-b border-parchment-border py-4 px-6 md:px-12 flex items-center justify-between">
    <!-- Left Menu Icon -->
    <button 
      @click="toggleSidebar"
      class="text-parchment-neutral hover:text-parchment-primary transition-colors p-1.5 rounded-full hover:bg-parchment-neutral-light/50 outline-none border border-transparent shadow-none"
      id="top-nav-menu-btn"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="4" x2="20" y1="12" y2="12"></line>
        <line x1="4" x2="20" y1="6" y2="6"></line>
        <line x1="4" x2="20" y1="18" y2="18"></line>
      </svg>
    </button>

    <!-- Center Title -->
    <router-link to="/" class="absolute left-1/2 -translate-x-1/2 flex items-center">
      <span class="font-serif text-parchment-primary-dark text-lg md:text-xl font-medium tracking-[0.1em] hover:opacity-80 transition-opacity">
        Digital Liturgy
      </span>
    </router-link>

    <!-- Right Navigation (Desktop) & Profile -->
    <div class="flex items-center space-x-6 md:space-x-8">
      <nav class="hidden md:flex items-center space-x-6">
        <router-link 
          to="/" 
          class="text-sm font-medium tracking-wider uppercase transition-colors"
          :class="[route.path === '/' ? 'text-parchment-primary-dark font-bold' : 'text-parchment-neutral/60 hover:text-parchment-neutral']"
        >
          Home
        </router-link>
        <router-link 
          to="/daily-readings" 
          class="text-sm font-medium tracking-wider uppercase transition-colors"
          :class="[route.path !== '/' ? 'text-parchment-primary-dark font-bold' : 'text-parchment-neutral/60 hover:text-parchment-neutral']"
        >
          Devotions
        </router-link>
        <button 
          @click="showCalendarToast"
          class="text-sm font-medium tracking-wider uppercase text-parchment-neutral/60 hover:text-parchment-neutral bg-transparent border-none p-0 outline-none shadow-none hover:bg-transparent hover:translate-y-0"
        >
          Calendar
        </button>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute();

const toggleSidebar = () => {
  alert('Settings and options menu coming soon!');
};

const showCalendarToast = () => {
  alert('Liturgical Calendar feature is coming soon!');
};
</script>
