<template>
  <header class="relative z-30 border-b border-[#D1C7BD]/55 pb-4 mb-5 animate-fade-in-down">
    <div class="flex items-start sm:items-center justify-between gap-3">
      <!-- Left: Back Link & Title -->
      <div class="min-w-0 flex-1">
        <div class="flex flex-wrap items-center gap-x-3 gap-y-1 mb-1">
          <router-link
            to="/devotions"
            class="inline-flex items-center gap-1.5 text-[11px] font-bold uppercase tracking-[0.16em] text-[#72383D]/80 hover:text-[#72383D] transition-colors group"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-0.5 transition-transform">
              <path d="m15 18-6-6 6-6" />
            </svg>
            <span>All Devotions</span>
          </router-link>

          <span
            v-if="subtitle"
            class="hidden sm:inline-flex items-center gap-1.5 text-[10px] font-bold uppercase tracking-[0.18em] text-[#AC9C8D]"
          >
            <span>&bull;</span>
            <span class="text-[#72383D]/75">{{ subtitle }}</span>
          </span>
        </div>

        <h1 class="font-serif text-2xl md:text-3xl bg-gradient-to-r from-[#72383D] via-[#5A2D31] to-[#322D29] bg-clip-text text-transparent font-bold tracking-tight truncate">
          {{ title }}
        </h1>
      </div>

      <!-- Right: Collapsible Settings Menu Trigger -->
      <div ref="menuWrapperRef" class="relative shrink-0 pt-0.5">
        <button
          type="button"
          @click="isOpen = !isOpen"
          :aria-expanded="isOpen"
          aria-haspopup="true"
          aria-label="Prayer settings"
          class="flex items-center gap-2 px-3.5 py-2 min-h-[42px] rounded-full bg-white/70 hover:bg-white/95 border border-[#D1C7BD]/80 text-[#322D29] shadow-xs backdrop-blur-md transition-all cursor-pointer active:scale-95"
        >
          <!-- Sliders Icon (matches reference design) -->
          <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]">
            <line x1="4" x2="20" y1="6" y2="6" />
            <circle cx="9" cy="6" r="2" fill="currentColor" />
            <line x1="4" x2="20" y1="12" y2="12" />
            <circle cx="15" cy="12" r="2" fill="currentColor" />
            <line x1="4" x2="20" y1="18" y2="18" />
            <circle cx="10" cy="18" r="2" fill="currentColor" />
          </svg>
          <span class="hidden sm:inline text-xs font-bold uppercase tracking-wider text-[#322D29]/85">
            Settings
          </span>
          <span class="px-1.5 py-0.5 rounded-full bg-[#72383D]/10 text-[#72383D] text-[10px] font-bold uppercase tracking-wider">
            {{ showLatin ? 'LA' : 'EN' }}
          </span>
          <span
            v-if="activeIntentionsCount > 0"
            class="w-2 h-2 rounded-full bg-[#72383D] ring-2 ring-white"
            :title="`${activeIntentionsCount} active prayer intentions`"
          />
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="13"
            height="13"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="text-[#322D29]/50 transition-transform duration-200"
            :class="{ 'rotate-180': isOpen }"
          >
            <path d="m6 9 6 6 6-6" />
          </svg>
        </button>

        <!-- Collapsible Settings Menu Panel -->
        <transition name="settings-pop">
          <div
            v-if="isOpen"
            class="absolute right-0 mt-2.5 w-72 sm:w-80 rounded-[1.75rem] devotion-glass-panel p-4 z-50 space-y-3.5"
          >
            <!-- Menu Header -->
            <div class="flex items-center justify-between pb-2 border-b border-[#D1C7BD]/50">
              <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D]">
                Prayer Settings
              </span>
              <button
                type="button"
                @click="isOpen = false"
                class="w-7 h-7 min-h-0 p-0 rounded-full border border-[#D1C7BD]/60 bg-white/60 text-[#322D29]/60 hover:text-[#72383D] hover:bg-white flex items-center justify-center cursor-pointer"
                aria-label="Close settings"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" x2="6" y1="6" y2="18" />
                  <line x1="6" x2="18" y1="6" y2="18" />
                </svg>
              </button>
            </div>

            <!-- Language Toggle -->
            <div class="flex items-center justify-between gap-3">
              <div>
                <div class="text-xs font-bold text-[#322D29]">Language</div>
                <div class="text-[10px] text-[#322D29]/55">English or Latin text</div>
              </div>
              <div class="flex items-center bg-white/75 border border-[#D1C7BD]/80 p-0.5 rounded-full">
                <button
                  type="button"
                  @click="$emit('update:showLatin', false)"
                  :class="[!showLatin ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#322D29]/70 hover:text-[#72383D] bg-transparent']"
                  class="px-3 py-1.5 min-h-[32px] rounded-full text-[10px] font-bold uppercase tracking-wider transition-all border-none cursor-pointer"
                >
                  EN
                </button>
                <button
                  type="button"
                  @click="$emit('update:showLatin', true)"
                  :class="[showLatin ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#322D29]/70 hover:text-[#72383D] bg-transparent']"
                  class="px-3 py-1.5 min-h-[32px] rounded-full text-[10px] font-bold uppercase tracking-wider transition-all border-none cursor-pointer"
                >
                  LA
                </button>
              </div>
            </div>

            <!-- Optional Short/Full Mode (Divine Mercy) -->
            <div v-if="hasShortMode" class="flex items-center justify-between gap-3 pt-1 border-t border-[#D1C7BD]/35">
              <div>
                <div class="text-xs font-bold text-[#322D29]">Chaplet Format</div>
                <div class="text-[10px] text-[#322D29]/55">Full decades or shortened</div>
              </div>
              <div class="flex items-center bg-white/75 border border-[#D1C7BD]/80 p-0.5 rounded-full">
                <button
                  type="button"
                  @click="$emit('update:isShortVersion', false)"
                  :class="[!isShortVersion ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#322D29]/70 hover:text-[#72383D] bg-transparent']"
                  class="px-3 py-1.5 min-h-[32px] rounded-full text-[10px] font-bold uppercase tracking-wider transition-all border-none cursor-pointer"
                >
                  Full
                </button>
                <button
                  type="button"
                  @click="$emit('update:isShortVersion', true)"
                  :class="[isShortVersion ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-xs' : 'text-[#322D29]/70 hover:text-[#72383D] bg-transparent']"
                  class="px-3 py-1.5 min-h-[32px] rounded-full text-[10px] font-bold uppercase tracking-wider transition-all border-none cursor-pointer"
                >
                  Short
                </button>
              </div>
            </div>

            <!-- Auto-Play Audio Toggle -->
            <label
              v-if="supportsAutoPlay"
              class="flex items-center justify-between gap-3 pt-2 border-t border-[#D1C7BD]/35 cursor-pointer select-none"
            >
              <div>
                <div class="text-xs font-bold text-[#322D29]">Auto-Play Audio</div>
                <div class="text-[10px] text-[#322D29]/55">Automatically advance with audio</div>
              </div>
              <div class="relative shrink-0">
                <input
                  type="checkbox"
                  :checked="autoPlay"
                  @change="$emit('update:autoPlay', ($event.target as HTMLInputElement).checked)"
                  class="sr-only peer"
                />
                <div class="w-10 h-5.5 bg-white/80 border border-[#D1C7BD] rounded-full peer peer-checked:bg-[#72383D]/25 peer-checked:border-[#72383D]/50 transition-colors" />
                <div class="absolute top-0.5 left-0.5 w-4.5 h-4.5 rounded-full bg-[#322D29]/45 peer-checked:bg-[#72383D] peer-checked:translate-x-[18px] transition-all shadow-2xs" />
              </div>
            </label>

            <!-- Action Rows -->
            <div class="pt-2 border-t border-[#D1C7BD]/40 space-y-1.5">
              <!-- Prayer Intentions -->
              <button
                type="button"
                @click="handleAction('open-intentions')"
                class="w-full flex items-center justify-between px-3 py-2.5 min-h-[42px] rounded-2xl bg-white/60 hover:bg-white border border-[#D1C7BD]/55 text-left transition-all cursor-pointer"
              >
                <div class="flex items-center gap-2.5">
                  <span class="text-sm">🙏</span>
                  <span class="text-xs font-bold text-[#322D29]">Prayer Intentions</span>
                </div>
                <span
                  v-if="activeIntentionsCount > 0"
                  class="px-2 py-0.5 rounded-full bg-[#72383D] text-white text-[10px] font-bold"
                >
                  {{ activeIntentionsCount }}
                </span>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-[#322D29]/40">
                  <path d="m9 18 6-6-6-6" />
                </svg>
              </button>

              <!-- Custom Prayers (Rosary) -->
              <button
                v-if="hasCustomPrayers"
                type="button"
                @click="handleAction('open-custom-prayers')"
                class="w-full flex items-center justify-between px-3 py-2.5 min-h-[42px] rounded-2xl bg-white/60 hover:bg-white border border-[#D1C7BD]/55 text-left transition-all cursor-pointer"
              >
                <div class="flex items-center gap-2.5">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]">
                    <path d="M12 20h9" />
                    <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" />
                  </svg>
                  <span class="text-xs font-bold text-[#322D29]">Custom Prayers</span>
                </div>
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-[#322D29]/40">
                  <path d="m9 18 6-6-6-6" />
                </svg>
              </button>

              <!-- Fullscreen Mode -->
              <button
                v-if="isFullscreenSupported"
                type="button"
                @click="handleAction('toggle-fullscreen')"
                class="w-full flex items-center justify-between px-3 py-2.5 min-h-[42px] rounded-2xl bg-white/60 hover:bg-white border border-[#D1C7BD]/55 text-left transition-all cursor-pointer"
              >
                <div class="flex items-center gap-2.5">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]">
                    <polyline v-if="!isFullscreen" points="15 3 21 3 21 9" />
                    <polyline v-if="!isFullscreen" points="9 21 3 21 3 15" />
                    <line v-if="!isFullscreen" x1="21" y1="3" x2="14" y2="10" />
                    <line v-if="!isFullscreen" x1="3" y1="21" x2="10" y2="14" />
                    <polyline v-if="isFullscreen" points="4 14 10 14 10 20" />
                    <polyline v-if="isFullscreen" points="20 10 14 10 14 4" />
                    <line v-if="isFullscreen" x1="14" y1="10" x2="21" y2="3" />
                    <line v-if="isFullscreen" x1="10" y1="14" x2="3" y2="21" />
                  </svg>
                  <span class="text-xs font-bold text-[#322D29]">
                    {{ isFullscreen ? 'Exit Fullscreen' : 'Fullscreen Mode' }}
                  </span>
                </div>
                <span class="text-[10px] font-mono text-[#322D29]/45">F</span>
              </button>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';

withDefaults(defineProps<{
  title: string;
  subtitle?: string;
  showLatin: boolean;
  autoPlay?: boolean;
  supportsAutoPlay?: boolean;
  activeIntentionsCount?: number;
  isFullscreen?: boolean;
  isFullscreenSupported?: boolean;
  hasShortMode?: boolean;
  isShortVersion?: boolean;
  hasCustomPrayers?: boolean;
}>(), {
  supportsAutoPlay: true,
  activeIntentionsCount: 0,
  isFullscreen: false,
  isFullscreenSupported: true,
  hasShortMode: false,
  isShortVersion: false,
  hasCustomPrayers: false,
});

const emit = defineEmits<{
  (e: 'update:showLatin', val: boolean): void;
  (e: 'update:autoPlay', val: boolean): void;
  (e: 'update:isShortVersion', val: boolean): void;
  (e: 'open-intentions'): void;
  (e: 'open-custom-prayers'): void;
  (e: 'toggle-fullscreen'): void;
}>();

const isOpen = ref(false);
const menuWrapperRef = ref<HTMLElement | null>(null);

const handleAction = (eventName: 'open-intentions' | 'open-custom-prayers' | 'toggle-fullscreen') => {
  isOpen.value = false;
  if (eventName === 'open-intentions') emit('open-intentions');
  else if (eventName === 'open-custom-prayers') emit('open-custom-prayers');
  else emit('toggle-fullscreen');
};

const handleClickOutside = (e: MouseEvent) => {
  if (isOpen.value && menuWrapperRef.value && !menuWrapperRef.value.contains(e.target as Node)) {
    isOpen.value = false;
  }
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isOpen.value) {
    isOpen.value = false;
  }
};

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside);
  window.addEventListener('keydown', handleKeyDown);
});

onBeforeUnmount(() => {
  document.removeEventListener('mousedown', handleClickOutside);
  window.removeEventListener('keydown', handleKeyDown);
});
</script>
