<template>
  <footer class="w-full bg-parchment-bg border-t border-parchment-border py-2 px-6 flex items-center justify-around fixed bottom-0 left-0 right-0 z-40">
    <!-- Home Tab -->
    <router-link to="/" class="flex flex-col items-center group transition-all duration-300">
      <div 
        class="w-12 h-10 flex items-center justify-center rounded-xl transition-all duration-300"
        :class="[route.path === '/' ? 'bg-parchment-primary-dark text-white' : 'text-parchment-neutral/50 group-hover:text-parchment-neutral']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <polyline points="9 22 9 12 15 12 15 22"></polyline>
        </svg>
      </div>
      <span 
        class="text-[11px] font-medium uppercase tracking-wider mt-1 transition-all duration-300"
        :class="[route.path === '/' ? 'text-parchment-primary-dark font-bold' : 'text-parchment-neutral/50 group-hover:text-parchment-neutral']"
      >
        Home
      </span>
    </router-link>

    <!-- Devotions Tab -->
    <router-link to="/daily-readings" class="flex flex-col items-center group transition-all duration-300">
      <div 
        class="w-12 h-10 flex items-center justify-center rounded-xl transition-all duration-300"
        :class="[route.path !== '/' && route.path !== '/search' && route.path !== '/profile' ? 'bg-parchment-primary-dark text-white' : 'text-parchment-neutral/50 group-hover:text-parchment-neutral']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
          <path d="M6 6h10"></path>
          <path d="M6 10h10"></path>
        </svg>
      </div>
      <span 
        class="text-[11px] font-medium uppercase tracking-wider mt-1 transition-all duration-300"
        :class="[route.path !== '/' && route.path !== '/search' && route.path !== '/profile' ? 'text-parchment-primary-dark font-bold' : 'text-parchment-neutral/50 group-hover:text-parchment-neutral']"
      >
        Devotions
      </span>
    </router-link>

    <!-- Search Tab -->
    <button 
      @click="triggerSearch"
      class="flex flex-col items-center group transition-all duration-300 bg-transparent border-none p-0 outline-none shadow-none hover:bg-transparent hover:translate-y-0"
    >
      <div 
        class="w-12 h-10 flex items-center justify-center rounded-xl text-parchment-neutral/50 group-hover:text-parchment-neutral transition-all duration-300"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
        </svg>
      </div>
      <span 
        class="text-[11px] font-medium uppercase tracking-wider mt-1 text-parchment-neutral/50 group-hover:text-parchment-neutral transition-all duration-300"
      >
        Search
      </span>
    </button>
  </footer>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

const triggerSearch = () => {
  const query = prompt('What devotional or paragraph are you looking for?');
  if (query) {
    router.push({ path: '/catechism', query: { q: query.trim() } });
  }
};

</script>
