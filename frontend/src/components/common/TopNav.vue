<template>
  <header class="sticky top-0 z-40 w-full bg-white/60 backdrop-blur-md border-b border-[#E9D5FF]/60 py-4 px-6 md:px-12 flex items-center justify-between">
    <!-- Left Menu Icon -->
    <button 
      @click="toggleSidebar"
      class="text-[#9333EA] hover:text-[#7E22CE] transition-colors p-1.5 rounded-full hover:bg-[#E9D5FF]/30 outline-none border border-transparent shadow-none"
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
      <span class="font-serif text-lg md:text-xl font-semibold tracking-[0.1em] hover:opacity-80 transition-opacity">
        Faith and Devotion
      </span>
    </router-link>

    <!-- Right Navigation (Desktop) & Profile -->
    <div class="flex items-center space-x-6 md:space-x-8">
      <nav class="hidden md:flex items-center space-x-6">
        <router-link 
          to="/" 
          class="text-sm font-medium tracking-wider uppercase transition-colors"
          :class="[route.path === '/' ? 'text-[#7E22CE] font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
        >
          Home
        </router-link>
        <router-link 
          to="/daily-readings" 
          class="text-sm font-medium tracking-wider uppercase transition-colors"
          :class="[route.path !== '/' ? 'text-[#7E22CE] font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
        >
          Devotions
        </router-link>
        <router-link 
          to="/bible-study" 
          class="text-sm font-medium tracking-wider uppercase transition-colors"
          :class="[route.path === '/bible-study' ? 'text-[#7E22CE] font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
        >
          Bible Study
        </router-link>
        <router-link 
          to="/magisterium-chat" 
          class="text-sm font-medium tracking-wider uppercase transition-colors"
          :class="[route.path === '/magisterium-chat' ? 'text-[#7E22CE] font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
        >
          Magisterium
        </router-link>
      </nav>
    </div>
  </header>

  <!-- Sidebar Overlay & Drawer via Teleport -->
  <Teleport to="body">
    <!-- Backdrop -->
    <Transition name="fade">
      <div 
        v-if="isSidebarOpen"
        @click="closeSidebar"
        class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm cursor-pointer"
      ></div>
    </Transition>

    <!-- Drawer Panel -->
    <Transition name="slide">
      <div 
        v-if="isSidebarOpen"
        class="fixed top-0 left-0 bottom-0 z-50 w-80 max-w-[85vw] bg-white/90 backdrop-blur-md border-r border-[#E9D5FF]/60 shadow-2xl flex flex-col"
      >
        <!-- Header -->
        <div class="p-6 border-b border-[#E9D5FF]/60 flex items-center justify-between">
          <span class="font-serif text-lg font-bold tracking-[0.1em]">
            Faith and Devotion
          </span>
          <button 
            @click="closeSidebar"
            class="hover:text-[#7E22CE] transition-colors p-1.5 rounded-full hover:bg-[#E9D5FF]/30 border border-transparent outline-none cursor-pointer"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" x2="6" y1="6" y2="18"></line>
              <line x1="6" x2="18" y1="6" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Navigation Links -->
        <nav class="flex-1 overflow-y-auto py-6 px-4 space-y-2">
          <router-link 
            v-for="item in navItems" 
            :key="item.path" 
            :to="item.path"
            @click="closeSidebar"
            class="flex items-center space-x-4 p-3 rounded-2xl transition-all duration-200 group border"
            :class="[route.path === item.path ? 'bg-[#E9D5FF]/40 border-[#E9D5FF] text-[#7E22CE] font-bold shadow-xs' : 'text-black/80 hover:text-[#7E22CE] border-transparent hover:bg-white/60']"
          >
            <!-- Icon Wrapper -->
            <div 
              class="w-10 h-10 rounded-xl flex items-center justify-center transition-all duration-300"
              :class="[route.path === item.path ? 'bg-[#9333EA] text-white shadow-xs' : 'bg-[#E9D5FF]/30 text-[#9333EA] group-hover:bg-[#9333EA] group-hover:text-white']"
            >
              <!-- Inline SVGs depending on icon name -->
              <svg v-if="item.icon === 'home'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
              <svg v-else-if="item.icon === 'book'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path><path d="M6 6h10"></path><path d="M6 10h10"></path></svg>
              <svg v-else-if="item.icon === 'rosary'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22a7 7 0 0 0 7-7c0-4.3-7-11-7-11S5 10.7 5 15a7 7 0 0 0 7 7z"></path><circle cx="12" cy="15" r="2"></circle></svg>
              <svg v-else-if="item.icon === 'heart'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"></path></svg>
              <svg v-else-if="item.icon === 'shield'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.42 10.922a1 1 0 0 0-.019-1.838L12.83 5.18a2 2 0 0 0-1.66 0L2.6 9.08a1 1 0 0 0 0 1.832l8.57 3.908a2 2 0 0 0 1.66 0z"></path><path d="M6 12v5c0 2 2 3 6 3s6-1 6-3v-5"></path></svg>
              <svg v-else-if="item.icon === 'star'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.912 5.886L4.2 9.08l4.956 3.6L7.243 18.57 12 15l4.757 3.57-1.913-5.887 4.957-3.6-5.887-.193Z"></path></svg>
              <svg v-else-if="item.icon === 'scroll'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path></svg>
              <svg v-else-if="item.icon === 'context'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path></svg>
              <svg v-else-if="item.icon === 'chat'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
            </div>
            <div class="text-left flex-1 min-w-0">
              <div class="text-sm font-semibold tracking-wide leading-none mb-1">{{ item.label }}</div>
              <div class="text-[11px] text-black/60 truncate">{{ item.desc }}</div>
            </div>
          </router-link>
        </nav>

        <!-- Footer / AMDG -->
        <div class="p-6 border-t border-[#E9D5FF]/60 text-center">
          <span class="text-xs text-[#D97706] tracking-wider uppercase font-semibold">
            Ad Maiorem Dei Gloriam
          </span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const isSidebarOpen = ref(false);

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

const closeSidebar = () => {
  isSidebarOpen.value = false;
};

const navItems = [
  {
    path: '/',
    label: 'Home',
    desc: 'Return to the main sanctuary',
    icon: 'home'
  },
  {
    path: '/daily-readings',
    label: 'Daily Readings',
    desc: 'Reflect on today\'s liturgical scriptures',
    icon: 'book'
  },
  {
    path: '/rosary',
    label: 'Scriptural Rosary',
    desc: 'Meditate on the Holy Rosary mysteries',
    icon: 'rosary'
  },
  {
    path: '/st-michael-rosary',
    label: 'St. Michael Chaplet',
    desc: 'Pray the Chaplet of St. Michael & the 9 Angelic Choirs',
    icon: 'shield'
  },
  {
    path: '/divine-mercy',
    label: 'Divine Mercy Chaplet',
    desc: 'Pray the Chaplet of Divine Mercy',
    icon: 'heart'
  },
  {
    path: '/catechism',
    label: 'Catechism',
    desc: 'Explore the teachings of the Church',
    icon: 'shield'
  },
  {
    path: '/prayers',
    label: 'Traditional Prayers',
    desc: 'Novenas, litanies, and common prayers',
    icon: 'star'
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
  },
  {
    path: '/magisterium-chat',
    label: 'Magisterium AI Chat',
    desc: 'Ask questions on Church doctrine & teachings',
    icon: 'chat'
  }
];
</script>

<style scoped>
/* Sidebar Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}
</style>
