<template>
  <div class="w-full max-w-2xl mx-auto pt-3">
    <div class="flex items-center justify-between gap-3 px-2 py-2 rounded-full bg-white/55 border border-[#D1C7BD]/70 backdrop-blur-xl shadow-[0_12px_32px_-12px_rgba(50,45,41,0.12)]">
      <!-- Previous Button -->
      <button
        type="button"
        :disabled="currentStepIndex === 0"
        @click="$emit('prev')"
        class="flex items-center gap-1.5 px-4 py-2.5 min-h-[44px] rounded-full bg-white/80 hover:bg-white border border-[#D1C7BD]/80 text-[#322D29] text-xs font-bold uppercase tracking-wider transition-all disabled:opacity-35 disabled:pointer-events-none cursor-pointer active:scale-95"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="m15 18-6-6 6-6" />
        </svg>
        <span class="hidden sm:inline">Previous</span>
      </button>

      <!-- Center: Audio Play/Pause Ring + Step Counter -->
      <div class="flex items-center gap-3">
        <button
          v-if="hasAudio"
          type="button"
          @click="$emit('toggle-audio')"
          :aria-label="isPlaying ? 'Pause prayer audio' : 'Listen to prayer audio'"
          :title="isPlaying ? 'Pause Audio (Space)' : 'Listen to Prayer (Space)'"
          class="relative w-12 h-12 min-h-[48px] p-0 rounded-full bg-gradient-to-br from-[#72383D] to-[#322D29] text-white border-2 border-white/80 shadow-[0_8px_20px_-4px_rgba(114,56,61,0.45)] hover:from-[#8B464C] hover:to-[#453E38] flex items-center justify-center transition-all cursor-pointer active:scale-95"
        >
          <!-- Subtle active pulse ring when playing -->
          <span
            v-if="isPlaying"
            class="absolute -inset-1.5 rounded-full border border-[#72383D]/45 animate-ping pointer-events-none"
          />
          <svg v-if="!isPlaying" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="ml-0.5">
            <polygon points="6 3 20 12 6 21 6 3" />
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.8" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" x2="18" y1="5" y2="19" />
            <line x1="6" x2="6" y1="5" y2="19" />
          </svg>
        </button>

        <div class="flex flex-col items-center min-w-[68px]">
          <span class="text-[9px] font-bold uppercase tracking-[0.2em] text-[#72383D]/80">
            {{ label }}
          </span>
          <div class="text-xs md:text-sm font-bold tabular-nums tracking-wider text-[#322D29]">
            {{ currentStepIndex + 1 }}
            <span class="text-[#AC9C8D] mx-0.5">/</span>
            {{ totalSteps }}
          </div>
          <!-- Thin progress bar -->
          <div class="w-16 h-1 rounded-full bg-[#D1C7BD]/50 overflow-hidden mt-1">
            <div
              class="h-full rounded-full bg-gradient-to-r from-[#72383D] to-[#AC9C8D] transition-all duration-300"
              :style="{ width: `${progressPercent}%` }"
            />
          </div>
        </div>
      </div>

      <!-- Next / Done Button -->
      <button
        type="button"
        :disabled="disableNextWhenLast && isLastStep"
        @click="$emit('next')"
        class="flex items-center gap-1.5 px-4 sm:px-5 py-2.5 min-h-[44px] rounded-full bg-gradient-to-r from-[#72383D] to-[#322D29] hover:from-[#8B464C] hover:to-[#453E38] text-white border border-transparent text-xs font-bold uppercase tracking-wider shadow-sm transition-all disabled:opacity-35 disabled:pointer-events-none cursor-pointer active:scale-95"
      >
        <span class="hidden sm:inline">{{ isLastStep ? 'Done' : 'Next' }}</span>
        <span class="sm:hidden">{{ isLastStep ? 'Done' : 'Next' }}</span>
        <svg v-if="!isLastStep" xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="m9 18 6-6-6-6" />
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="20 6 9 17 4 12" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  currentStepIndex: number;
  totalSteps: number;
  isLastStep?: boolean;
  disableNextWhenLast?: boolean;
  hasAudio?: boolean;
  isPlaying?: boolean;
  label?: string;
}>(), {
  isLastStep: false,
  disableNextWhenLast: false,
  hasAudio: false,
  isPlaying: false,
  label: 'Progress',
});

defineEmits<{
  (e: 'next'): void;
  (e: 'prev'): void;
  (e: 'toggle-audio'): void;
}>();

const progressPercent = computed(() => {
  if (!props.totalSteps) return 0;
  return Math.min(100, Math.max(0, ((props.currentStepIndex + 1) / props.totalSteps) * 100));
});
</script>
