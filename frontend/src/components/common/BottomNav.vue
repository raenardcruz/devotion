<template>
  <footer class="w-full bg-[#EFE9E1]/85 via-white/95 to-[#EFE9E1]/85 backdrop-blur-xl border-t border-white/60 pt-2 pb-[calc(0.6rem+env(safe-area-inset-bottom))] px-3 flex items-center justify-around fixed bottom-0 left-0 right-0 z-40 shadow-lg touch-none select-none">
    <!-- Home Tab -->
    <router-link to="/" class="flex flex-col items-center group transition-all duration-300 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-12 h-10 flex items-center justify-center rounded-xl transition-all duration-300 shadow-2xs"
        :class="[isHomeActive ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#72383D]/60 group-hover:text-[#72383D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <polyline points="9 22 9 12 15 12 15 22"></polyline>
        </svg>
      </div>
      <span 
        class="text-[10px] font-semibold uppercase tracking-wider mt-0.5 transition-all duration-300"
        :class="[isHomeActive ? 'text-[#72383D] font-bold' : 'text-[#322D29]/70 group-hover:text-[#322D29]']"
      >
        Home
      </span>
    </router-link>

    <!-- Devotions Tab -->
    <router-link to="/daily-readings" class="flex flex-col items-center group transition-all duration-300 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-12 h-10 flex items-center justify-center rounded-xl transition-all duration-300 shadow-2xs"
        :class="[isDevotionsActive ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#72383D]/60 group-hover:text-[#72383D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
          <path d="M6 6h10"></path>
          <path d="M6 10h10"></path>
        </svg>
      </div>
      <span 
        class="text-[10px] font-semibold uppercase tracking-wider mt-0.5 transition-all duration-300"
        :class="[isDevotionsActive ? 'text-[#72383D] font-bold' : 'text-[#322D29]/70 group-hover:text-[#322D29]']"
      >
        Devotions
      </span>
    </router-link>

    <!-- Bible Study Tab -->
    <router-link to="/bible-study" class="flex flex-col items-center group transition-all duration-300 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-12 h-10 flex items-center justify-center rounded-xl transition-all duration-300 shadow-2xs"
        :class="[isBibleStudyActive ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#72383D]/60 group-hover:text-[#72383D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>
        </svg>
      </div>
      <span 
        class="text-[10px] font-semibold uppercase tracking-wider mt-0.5 transition-all duration-300"
        :class="[isBibleStudyActive ? 'text-[#72383D] font-bold' : 'text-[#322D29]/70 group-hover:text-[#322D29]']"
      >
        Study
      </span>
    </router-link>

    <!-- Magisterium AI Tab -->
    <router-link to="/magisterium-chat" class="flex flex-col items-center group transition-all duration-300 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-12 h-10 flex items-center justify-center rounded-xl transition-all duration-300 shadow-2xs"
        :class="[isMagisteriumActive ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#72383D]/60 group-hover:text-[#72383D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>
      </div>
      <span 
        class="text-[10px] font-semibold uppercase tracking-wider mt-0.5 transition-all duration-300"
        :class="[isMagisteriumActive ? 'text-[#72383D] font-bold' : 'text-[#322D29]/70 group-hover:text-[#322D29]']"
      >
        Magisterium
      </span>
    </router-link>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const devotionsRoutes = [
  '/daily-readings',
  '/rosary',
  '/st-michael-rosary',
  '/divine-mercy',
  '/prayers',
  '/catechism',
  '/resources'
];

const isHomeActive = computed(() => route.path === '/');
const isDevotionsActive = computed(() => devotionsRoutes.some(path => route.path.startsWith(path)));
const isBibleStudyActive = computed(() => route.path.startsWith('/bible-study'));
const isMagisteriumActive = computed(() => route.path.startsWith('/magisterium-chat'));
</script>
