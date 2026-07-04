<template>
  <transition name="fade">
    <button
      v-show="isVisible"
      @click="scrollToTop"
      class="fixed bottom-20 left-1/2 z-40 px-5 py-2.5 rounded-full bg-parchment-primary hover:bg-parchment-primary-dark text-white font-bold uppercase tracking-wider text-[10px] shadow-lg active:scale-95 transition-all duration-300 flex items-center gap-1.5 cursor-pointer border border-transparent outline-none focus:ring-2 focus:ring-parchment-primary/30"
      aria-label="Back to top"
      id="back-to-top-btn"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="12"
        height="12"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="3"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="animate-bounce-slow"
      >
        <polyline points="18 15 12 9 6 15"></polyline>
      </svg>
      <span>Back to Top</span>
    </button>
  </transition>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';

const isVisible = ref(false);

const checkScroll = () => {
  isVisible.value = window.scrollY > 300;
};

const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth',
  });
};

onMounted(() => {
  window.addEventListener('scroll', checkScroll);
  checkScroll();
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', checkScroll);
});
</script>

<style scoped>
#back-to-top-btn {
  transform: translate(-50%, 0) scale(1);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translate(-50%, 15px) scale(0.9) !important;
}

.animate-bounce-slow {
  animation: bounceSlow 2s infinite;
}

@keyframes bounceSlow {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-4px);
  }
  60% {
    transform: translateY(-2px);
  }
}
</style>
