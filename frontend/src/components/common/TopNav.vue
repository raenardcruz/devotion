<template>
  <header class="sticky top-0 z-40 w-full bg-[#F7F3EE]/90 backdrop-blur-xl border-b border-[#E6DDD2]/80 py-3.5 px-4 md:px-10 flex items-center justify-between">
    <!-- Left Menu Icon -->
    <button 
      @click="toggleSidebar"
      class="text-[#5C282D] hover:text-[#322D29] transition-colors p-2.5 rounded-full hover:bg-[#E6DDD2]/50 outline-none border border-transparent shadow-none min-w-[44px] min-h-[44px] flex items-center justify-center active:scale-95 cursor-pointer"
      id="top-nav-menu-btn"
      aria-label="Open Navigation Menu"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="21" height="21" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
        <line x1="4" x2="20" y1="12" y2="12"></line>
        <line x1="4" x2="20" y1="6" y2="6"></line>
        <line x1="4" x2="20" y1="18" y2="18"></line>
      </svg>
    </button>

    <!-- Center Title -->
    <router-link to="/" class="absolute left-1/2 -translate-x-1/2 flex items-center min-h-[44px] justify-center">
      <span class="font-serif text-lg md:text-xl font-bold tracking-[0.04em] text-[#4A2525] hover:opacity-80 transition-opacity whitespace-nowrap">
        Faith and Devotion
      </span>
    </router-link>

    <!-- Right Navigation (Desktop) & Quick Action (Mobile) -->
    <div class="flex items-center space-x-4 md:space-x-7">
      <nav class="hidden lg:flex items-center space-x-6">
        <router-link 
          to="/" 
          class="text-[11px] font-semibold tracking-[0.12em] uppercase transition-colors py-1"
          :class="[isHomeActive ? 'text-[#5C282D] font-bold border-b-2 border-[#5C282D]' : 'text-[#322D29]/65 hover:text-[#5C282D]']"
        >
          Home
        </router-link>
        <router-link 
          to="/daily-readings" 
          class="text-[11px] font-semibold tracking-[0.12em] uppercase transition-colors py-1"
          :class="[isReadingsActive ? 'text-[#5C282D] font-bold border-b-2 border-[#5C282D]' : 'text-[#322D29]/65 hover:text-[#5C282D]']"
        >
          Readings
        </router-link>
        <router-link 
          to="/devotions" 
          class="text-[11px] font-semibold tracking-[0.12em] uppercase transition-colors py-1"
          :class="[isDevotionsActive ? 'text-[#5C282D] font-bold border-b-2 border-[#5C282D]' : 'text-[#322D29]/65 hover:text-[#5C282D]']"
        >
          Devotions
        </router-link>
        <router-link 
          to="/bible-study" 
          class="text-[11px] font-semibold tracking-[0.12em] uppercase transition-colors py-1"
          :class="[isBibleStudyActive ? 'text-[#5C282D] font-bold border-b-2 border-[#5C282D]' : 'text-[#322D29]/65 hover:text-[#5C282D]']"
        >
          Bible Study
        </router-link>
      </nav>

      <!-- Mobile Quick Link to Magisterium / Search -->
      <router-link
        to="/catechism"
        class="lg:hidden text-[#5C282D] hover:text-[#322D29] p-2 rounded-full hover:bg-[#E6DDD2]/50 transition-colors min-w-[40px] min-h-[40px] flex items-center justify-center"
        aria-label="Search Catechism"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="19" height="19" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.3-4.3"></path>
        </svg>
      </router-link>
    </div>
  </header>

  <!-- Sidebar Overlay & Drawer via Teleport -->
  <Teleport to="body">
    <!-- Backdrop -->
    <Transition name="fade">
      <div 
        v-if="isSidebarOpen"
        @click="closeSidebar"
        class="fixed inset-0 z-50 bg-[#322D29]/45 backdrop-blur-md cursor-pointer"
      ></div>
    </Transition>

    <!-- Drawer Panel -->
    <Transition name="slide">
      <div 
        v-if="isSidebarOpen"
        @touchstart="onTouchStart"
        @touchmove="onTouchMove"
        @touchend="onTouchEnd"
        class="fixed top-0 left-0 bottom-0 z-50 w-80 max-w-[85vw] bg-[#EFE9E1] border-r border-[#D1C7BD] shadow-2xl flex flex-col touch-pan-y"
      >
        <!-- Header -->
        <div class="p-5 border-b border-[#D1C7BD]/60 flex items-center justify-between">
          <span class="font-serif text-lg font-bold tracking-[0.08em] bg-gradient-to-r from-[#72383D] via-[#AC9C8D] to-[#322D29] bg-clip-text text-transparent">
            Faith and Devotion
          </span>
          <button 
            @click="closeSidebar"
            class="hover:text-[#72383D] transition-colors p-2 rounded-full hover:bg-[#D1C7BD]/30 border border-transparent outline-none cursor-pointer text-[#322D29] min-w-[44px] min-h-[44px] flex items-center justify-center active:scale-95"
            aria-label="Close Navigation Menu"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" x2="6" y1="6" y2="18"></line>
              <line x1="6" x2="18" y1="6" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Navigation Links -->
        <nav class="flex-1 overflow-y-auto py-5 px-4 space-y-2">
          <router-link 
            v-for="item in navItems" 
            :key="item.path" 
            :to="item.path"
            @click="closeSidebar"
            class="flex items-center space-x-4 p-3 rounded-2xl transition-all duration-200 group border min-h-[48px] active:scale-[0.98]"
            :class="[isNavItemActive(item.path) ? 'bg-gradient-to-r from-[#72383D]/15 to-[#AC9C8D]/20 border-[#D1C7BD] text-[#72383D] font-bold shadow-xs' : 'text-[#322D29]/80 hover:text-[#72383D] border-transparent hover:bg-white/60']"
          >
            <!-- Icon Wrapper -->
            <div 
              class="w-10 h-10 rounded-xl flex items-center justify-center transition-all duration-300 shadow-2xs shrink-0"
              :class="[isNavItemActive(item.path) ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white' : 'bg-[#D1C7BD]/40 text-[#72383D] group-hover:bg-gradient-to-r group-hover:from-[#72383D] group-hover:to-[#322D29] group-hover:text-white']"
            >
              <!-- Icons matching Home page illustrations -->
              <svg v-if="item.icon === 'home'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
              <RosaryIllustration v-else-if="item.icon === 'rosary'" class="w-8 h-8" />
              <MassReadingsIllustration v-else-if="item.icon === 'book'" class="w-8 h-8" />
              <OpenLectionaryHeroIcon v-else-if="item.icon === 'kindle'" class="w-8 h-8" />
              <CatechismIllustration v-else-if="item.icon === 'shield'" class="w-8 h-8" />
              <EarlyChurchIllustration v-else-if="item.icon === 'scroll'" class="w-8 h-8" />
              <BibleStudyIllustration v-else-if="item.icon === 'context'" class="w-8 h-8" />
            </div>
            <div class="text-left flex-1 min-w-0">
              <div class="text-sm font-semibold tracking-wide leading-none mb-1">{{ item.label }}</div>
              <div class="text-[11px] text-[#322D29]/70 truncate">{{ item.desc }}</div>
            </div>
          </router-link>
        </nav>

        <!-- Footer / AMDG -->
        <div class="p-5 border-t border-[#D1C7BD]/60 text-center">
          <span class="text-xs text-[#72383D] tracking-wider uppercase font-semibold">
            Ad Maiorem Dei Gloriam
          </span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import RosaryIllustration from '../icons/RosaryIllustration.vue';
import MassReadingsIllustration from '../icons/MassReadingsIllustration.vue';
import OpenLectionaryHeroIcon from '../icons/OpenLectionaryHeroIcon.vue';
import CatechismIllustration from '../icons/CatechismIllustration.vue';
import EarlyChurchIllustration from '../icons/EarlyChurchIllustration.vue';
import BibleStudyIllustration from '../icons/BibleStudyIllustration.vue';

const route = useRoute();
const isSidebarOpen = ref(false);

// Touch gesture state for swipe-to-close drawer
let touchStartX = 0;
let touchCurrentX = 0;

const onTouchStart = (e: TouchEvent) => {
  if (e.touches[0]) {
    touchStartX = e.touches[0].clientX;
    touchCurrentX = touchStartX;
  }
};

const onTouchMove = (e: TouchEvent) => {
  if (e.touches[0]) {
    touchCurrentX = e.touches[0].clientX;
  }
};

const onTouchEnd = () => {
  const deltaX = touchCurrentX - touchStartX;
  if (deltaX < -50) {
    closeSidebar();
  }
};

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

const closeSidebar = () => {
  isSidebarOpen.value = false;
};

const devotionChildRoutes = [
  '/devotions',
  '/rosary',
  '/seven-sorrows',
  '/st-michael-rosary',
  '/divine-mercy',
  '/prayers'
];

const isHomeActive = computed(() => route.path === '/');
const isReadingsActive = computed(() => route.path.startsWith('/daily-readings'));
const isDevotionsActive = computed(() => devotionChildRoutes.some(path => route.path.startsWith(path)));
const isBibleStudyActive = computed(() => route.path.startsWith('/bible-study'));

const isNavItemActive = (itemPath: string) => {
  if (itemPath === '/') return route.path === '/';
  if (itemPath === '/devotions') return isDevotionsActive.value;
  return route.path.startsWith(itemPath);
};

const navItems = [
  {
    path: '/',
    label: 'Home',
    desc: 'Return to the main sanctuary',
    icon: 'home'
  },
  {
    path: '/devotions',
    label: 'Devotions',
    desc: 'Rosary, Seven Sorrows, Divine Mercy & St. Michael',
    icon: 'rosary'
  },
  {
    path: '/daily-readings',
    label: 'Daily Readings',
    desc: 'Reflect on today\'s liturgical scriptures',
    icon: 'book'
  },
  {
    path: '/kindle',
    label: 'Kindle Readings',
    desc: 'E-Ink friendly daily readings without animations',
    icon: 'kindle'
  },
  {
    path: '/catechism',
    label: 'Catechism',
    desc: 'Explore the teachings of the Church',
    icon: 'shield'
  },
  {
    path: '/resources',
    label: 'Early Church Resources',
    desc: 'Writings of the Apostolic Fathers',
    icon: 'scroll'
  },
  {
    path: '/bible-study',
    label: 'Bible Study',
    desc: 'Retrieve scripture and commentary',
    icon: 'context'
  }
];
</script>

<style scoped>
/* Sidebar Transitions - Apple Spring Simulation */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}
</style>
