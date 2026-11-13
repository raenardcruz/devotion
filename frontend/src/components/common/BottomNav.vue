<template>
  <footer class="lg:hidden w-full bg-[#FAF7F2]/95 backdrop-blur-xl border-t border-[#E6DDD2] pt-2 pb-[calc(0.6rem+env(safe-area-inset-bottom))] px-2 sm:px-4 flex items-center justify-around fixed bottom-0 left-0 right-0 z-40 shadow-[0_-4px_20px_rgba(50,45,41,0.06)] touch-none select-none">
    <!-- Home Tab -->
    <router-link to="/" class="flex flex-col items-center group transition-all duration-200 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-11 h-8 flex items-center justify-center rounded-xl transition-all duration-200"
        :class="[isHomeActive ? 'bg-[#5C282D] text-[#FAF7F2] shadow-xs' : 'text-[#7B6E65] group-hover:text-[#5C282D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
          <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <polyline points="9 22 9 12 15 12 15 22"></polyline>
        </svg>
      </div>
      <span 
        class="text-[10px] tracking-wide mt-0.5 transition-all duration-200"
        :class="[isHomeActive ? 'text-[#5C282D] font-bold' : 'text-[#7B6E65] font-medium group-hover:text-[#322D29]']"
      >
        Home
      </span>
    </router-link>

    <!-- Devotions Hub Tab -->
    <router-link to="/devotions" class="flex flex-col items-center group transition-all duration-200 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-11 h-8 flex items-center justify-center rounded-xl transition-all duration-200"
        :class="[isDevotionsActive ? 'bg-[#5C282D] text-[#FAF7F2] shadow-xs' : 'text-[#7B6E65] group-hover:text-[#5C282D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="9" r="5"></circle>
          <path d="M12 14v7"></path>
          <path d="M9.5 18h5"></path>
        </svg>
      </div>
      <span 
        class="text-[10px] tracking-wide mt-0.5 transition-all duration-200"
        :class="[isDevotionsActive ? 'text-[#5C282D] font-bold' : 'text-[#7B6E65] font-medium group-hover:text-[#322D29]']"
      >
        Devotions
      </span>
    </router-link>

    <!-- Daily Readings Tab -->
    <router-link to="/daily-readings" class="flex flex-col items-center group transition-all duration-200 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-11 h-8 flex items-center justify-center rounded-xl transition-all duration-200"
        :class="[isReadingsActive ? 'bg-[#5C282D] text-[#FAF7F2] shadow-xs' : 'text-[#7B6E65] group-hover:text-[#5C282D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
          <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
          <path d="M12 7v6"></path>
          <path d="M9.5 9.5h5"></path>
        </svg>
      </div>
      <span 
        class="text-[10px] tracking-wide mt-0.5 transition-all duration-200"
        :class="[isReadingsActive ? 'text-[#5C282D] font-bold' : 'text-[#7B6E65] font-medium group-hover:text-[#322D29]']"
      >
        Readings
      </span>
    </router-link>

    <!-- Bible Study Tab -->
    <router-link to="/bible-study" class="flex flex-col items-center group transition-all duration-200 min-w-[56px] min-h-[44px] justify-center active:scale-95">
      <div 
        class="w-11 h-8 flex items-center justify-center rounded-xl transition-all duration-200"
        :class="[isBibleStudyActive ? 'bg-[#5C282D] text-[#FAF7F2] shadow-xs' : 'text-[#7B6E65] group-hover:text-[#5C282D]']"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
          <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
          <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
        </svg>
      </div>
      <span 
        class="text-[10px] tracking-wide mt-0.5 transition-all duration-200"
        :class="[isBibleStudyActive ? 'text-[#5C282D] font-bold' : 'text-[#7B6E65] font-medium group-hover:text-[#322D29]']"
      >
        Study
      </span>
    </router-link>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const devotionRoutes = [
  '/devotions',
  '/rosary',
  '/seven-sorrows',
  '/st-michael-rosary',
  '/divine-mercy',
  '/prayers'
];

const isHomeActive = computed(() => route.path === '/');
const isDevotionsActive = computed(() => devotionRoutes.some(path => route.path.startsWith(path)));
const isReadingsActive = computed(() => route.path.startsWith('/daily-readings') || route.path.startsWith('/catechism') || route.path.startsWith('/resources'));
const isBibleStudyActive = computed(() => route.path.startsWith('/bible-study'));
</script>
