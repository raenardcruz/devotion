<template>
  <!-- Floating Animated AI Avatar Trigger (Available on all pages) -->
  <transition name="avatar-appear">
    <div
      v-if="!isOpen"
      class="fixed bottom-[5.25rem] right-4 sm:right-6 lg:bottom-7 lg:right-7 z-40 flex items-center gap-3 group"
    >
      <!-- Hover Prompt Pill (Desktop) -->
      <div
        class="hidden md:flex items-center gap-2 px-3.5 py-2 rounded-2xl bg-[#FAF7F2]/95 border border-[#D1C7BD] shadow-lg backdrop-blur-md text-[#322D29] pointer-events-none opacity-0 translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 ease-out"
      >
        <span class="w-2 h-2 rounded-full bg-[#B5935A] animate-pulse"></span>
        <span class="font-serif text-xs font-semibold tracking-wide text-[#4A2525] whitespace-nowrap">
          Ask Magisterium AI
        </span>
      </div>

      <!-- Animated AI Avatar Button -->
      <button
        type="button"
        @click="openModal()"
        aria-label="Open Magisterium AI Chat"
        title="Ask Magisterium AI"
        class="relative !p-0 !border-0 !bg-transparent w-16 h-16 sm:w-[74px] sm:h-[74px] flex items-center justify-center cursor-pointer focus:outline-none focus-visible:ring-4 focus-visible:ring-[#B5935A]/50 rounded-full animate-avatar-float transition-transform duration-300 ease-[cubic-bezier(0.34,1.56,0.64,1)] hover:scale-110 active:scale-95"
      >
        <!-- Outer Breathing Sacred Aura -->
        <span
          class="absolute inset-1 rounded-full bg-gradient-to-tr from-[#6B2D33]/35 via-[#DFB878]/35 to-[#6B2D33]/35 blur-md animate-aura-pulse pointer-events-none"
        ></span>

        <!-- Home-Page Style Illustrated Medallion AI Avatar -->
        <MagisteriumAvatarIcon class="w-full h-full relative z-10" />

        <!-- Live / AI Ready Status Dot -->
        <span
          class="absolute top-1 right-1 z-20 w-4 h-4 rounded-full bg-[#FAF7EE] border border-[#6B2D33]/30 flex items-center justify-center shadow-xs"
        >
          <span class="w-2.5 h-2.5 rounded-full bg-[#DFB878] animate-ping absolute opacity-70"></span>
          <span class="w-2 h-2 rounded-full bg-[#C89B54] relative"></span>
        </span>
      </button>
    </div>
  </transition>

  <!-- Full-Screen Magisterium AI Modal -->
  <Teleport to="body">
    <transition name="modal-pop">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-[60] bg-[#F7F3EE] text-[#322D29] flex flex-col overflow-hidden selection:bg-[#72383D]/20"
        role="dialog"
        aria-modal="true"
        aria-label="Magisterium AI Sanctuary"
      >
        <!-- Full-Screen Chat Body (Edge-to-Edge Obsidian Canvas) -->
        <main class="flex-1 min-h-0 w-full flex flex-col">
          <MagisteriumSanctuary :is-modal="true" @close="closeModal" />
        </main>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { watch, onMounted, onBeforeUnmount } from 'vue';
import MagisteriumSanctuary from './MagisteriumSanctuary.vue';
import MagisteriumAvatarIcon from '../icons/MagisteriumAvatarIcon.vue';
import { useMagisteriumModal } from '../../composables/useMagisteriumModal';

const {
  isMagisteriumModalOpen: isOpen,
  openMagisteriumModal: openModal,
  closeMagisteriumModal: closeModal,
} = useMagisteriumModal();

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isOpen.value) {
    closeModal();
  }
};

// Lock body scroll while full-screen modal is open
watch(isOpen, (open) => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = open ? 'hidden' : '';
  }
});

onMounted(() => {
  window.addEventListener('keydown', handleKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown);
  if (typeof document !== 'undefined') {
    document.body.style.overflow = '';
  }
});
</script>

<style scoped>
/* Floating levitation */
@keyframes avatarFloat {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-5px);
  }
}

.animate-avatar-float {
  animation: avatarFloat 4s ease-in-out infinite;
}

.animate-avatar-float:hover {
  animation-play-state: paused;
}

/* Outer breathing aura */
@keyframes auraPulse {
  0%, 100% {
    opacity: 0.55;
    transform: scale(0.96);
  }
  50% {
    opacity: 0.95;
    transform: scale(1.08);
  }
}

.animate-aura-pulse {
  animation: auraPulse 3.2s ease-in-out infinite;
}

/* Slow rotating orbital ring */
@keyframes spinSlow {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin-slow {
  animation: spinSlow 18s linear infinite;
}

/* Periodic shimmer sweep across medallion */
@keyframes shimmerSweep {
  0%, 65% {
    transform: translateX(-120%) skewX(-15deg);
  }
  100% {
    transform: translateX(160%) skewX(-15deg);
  }
}

.animate-shimmer-sweep {
  animation: shimmerSweep 5s ease-in-out infinite;
}

/* Avatar eye blink */
@keyframes avatarBlink {
  0%, 44%, 48%, 92%, 96%, 100% {
    transform: scaleY(1);
  }
  46%, 94% {
    transform: scaleY(0.12);
  }
}

.animate-avatar-blink {
  animation: avatarBlink 5.5s ease-in-out infinite;
}

/* Halo breath */
@keyframes haloBreathe {
  0%, 100% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
}

.animate-halo-breathe {
  animation: haloBreathe 3s ease-in-out infinite;
}

/* Twinkling motes */
@keyframes twinkle {
  0%, 100% {
    opacity: 0.35;
    transform: scale(0.85);
  }
  50% {
    opacity: 1;
    transform: scale(1.15);
  }
}

.animate-twinkle {
  animation: twinkle 2.6s ease-in-out infinite;
}

.animate-twinkle-delayed {
  animation: twinkle 2.6s ease-in-out 1.3s infinite;
}

/* Avatar appear transition */
.avatar-appear-enter-active,
.avatar-appear-leave-active {
  transition: opacity 0.28s ease, transform 0.32s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.avatar-appear-enter-from,
.avatar-appear-leave-to {
  opacity: 0;
  transform: scale(0.75) translateY(12px);
}

/* Full-screen modal pop transition */
.modal-pop-enter-active {
  transition: opacity 0.25s ease-out, transform 0.3s cubic-bezier(0.22, 1, 0.36, 1);
}

.modal-pop-leave-active {
  transition: opacity 0.2s ease-in, transform 0.2s ease-in;
}

.modal-pop-enter-from,
.modal-pop-leave-to {
  opacity: 0;
  transform: scale(0.97);
}
</style>
