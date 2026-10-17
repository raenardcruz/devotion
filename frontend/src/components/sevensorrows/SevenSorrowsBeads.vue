<template>
  <div v-if="shouldShow" class="flex flex-col items-center gap-2.5 pt-3 pb-1 select-none">
    <!-- Septet Beads: 1 Medal/Our Father Bead + 7 Hail Mary Beads -->
    <div v-if="currentStep.section === 'sorrows'" class="flex flex-wrap justify-center items-center gap-2.5 sm:gap-3 min-h-[44px]">
      <!-- Scripture / Mystery Intro Pill -->
      <button
        @click="$emit('selectBead', -1)"
        title="Scripture Reading before Mystery"
        :class="[
          'px-2.5 py-1 rounded-full border text-[10px] font-bold uppercase tracking-wider cursor-pointer transition-all duration-300 outline-none flex items-center gap-1 min-h-[32px]',
          currentStep.type === 'sorrow-header'
            ? 'bg-[#72383D] text-white border-[#72383D] ring-4 ring-[#72383D]/20 shadow-xs'
            : 'bg-white/80 text-[#72383D] border-[#D1C7BD] hover:border-[#72383D]'
        ]"
      >
        <span>📖</span>
        <span class="hidden sm:inline">Verse</span>
      </button>

      <!-- Divider -->
      <div class="w-3 h-0.5 bg-[#D1C7BD]"></div>

      <!-- Large Bead (Our Father) -->
      <button
        @click="$emit('selectBead', 0)"
        title="Our Father"
        :class="[
          'w-7 h-7 rounded-full border-2 border-[#72383D] bg-white flex items-center justify-center cursor-pointer transition-all duration-300 hover:scale-110 outline-none shadow-2xs',
          currentStep.type === 'sorrow-our-father' ? 'ring-4 ring-[#72383D]/30 scale-110' : ''
        ]"
      >
        <div
          :class="[
            'w-3 h-3 rounded-full transition-colors',
            currentStep.type === 'sorrow-our-father' || (currentStep.beadNumber && currentStep.beadNumber > 0) || currentStep.type === 'sorrow-invocation'
              ? 'bg-[#72383D]'
              : 'bg-[#72383D]/30'
          ]"
        ></div>
      </button>

      <!-- Divider -->
      <div class="w-3 h-0.5 bg-[#D1C7BD]"></div>

      <!-- 7 Small Beads (7 Hail Marys) -->
      <div class="flex items-center gap-2">
        <button
          v-for="i in 7"
          :key="i"
          @click="$emit('selectBead', i)"
          :title="'Hail Mary ' + i + ' of 7'"
          :class="[
            'w-5 h-5 rounded-full border-2 border-[#72383D]/80 bg-white flex items-center justify-center cursor-pointer transition-all duration-300 hover:scale-115 outline-none shadow-2xs',
            currentStep.type === 'sorrow-hail-mary' && currentStep.beadNumber === i
              ? 'ring-4 ring-[#72383D]/30 scale-115 border-[#72383D]'
              : ''
          ]"
        >
          <div
            :class="[
              'w-2 h-2 rounded-full transition-colors',
              (currentStep.type === 'sorrow-hail-mary' && (currentStep.beadNumber || 0) >= i) || currentStep.type === 'sorrow-invocation'
                ? 'bg-[#72383D]'
                : 'bg-[#AC9C8D]/35'
            ]"
          ></div>
        </button>
      </div>
    </div>

    <!-- 3 Concluding Tears Beads -->
    <div v-else-if="currentStep.section === 'tears'" class="flex justify-center items-center gap-3.5 min-h-[44px]">
      <button
        v-for="i in 3"
        :key="i"
        @click="$emit('selectTearsBead', i)"
        :title="'Tears of Mary — Hail Mary ' + i + ' of 3'"
        :class="[
          'w-6 h-6 rounded-full border-2 border-[#72383D] bg-white flex items-center justify-center cursor-pointer transition-all duration-300 hover:scale-110 outline-none shadow-2xs',
          currentStep.beadNumber === i ? 'ring-4 ring-[#72383D]/30 scale-110' : ''
        ]"
      >
        <div
          :class="[
            'w-2.5 h-2.5 rounded-full transition-colors',
            (currentStep.beadNumber || 0) >= i ? 'bg-[#72383D]' : 'bg-[#72383D]/30'
          ]"
        ></div>
      </button>
    </div>

    <span class="text-[10px] text-[#322D29]/55 font-bold uppercase tracking-[0.18em]">
      {{ currentStep.section === 'tears' ? '3 Hail Marys in Honor of Mary’s Tears' : 'Tap verse or bead to jump within this Sorrow' }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { SevenSorrowsStep } from '../../data/sevenSorrowsData';

const props = defineProps<{
  currentStep: SevenSorrowsStep;
}>();

defineEmits<{
  (e: 'selectBead', beadNum: number): void;
  (e: 'selectTearsBead', beadNum: number): void;
}>();

const shouldShow = computed(() => {
  return props.currentStep.section === 'sorrows' || props.currentStep.section === 'tears';
});
</script>
