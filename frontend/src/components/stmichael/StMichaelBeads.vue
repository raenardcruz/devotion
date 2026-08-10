<template>
  <div v-if="shouldShow" class="flex flex-col items-center gap-2 mb-2 md:mb-4">
    <!-- Choir Beads: 1 Big Bead (Our Father) + 3 Small Beads (Hail Marys) -->
    <div v-if="currentStep.section === 'choirs'" class="flex justify-center items-center gap-3.5 h-10">
      <!-- Big Bead (Our Father) -->
      <button 
        @click="$emit('selectBead', 0)"
        :title="'Our Father (' + (currentStep.choirName || '') + ')'"
        :class="[
          'w-7 h-7 rounded-full border-2 border-[#72383D] bg-white flex items-center justify-center cursor-pointer transition-all duration-300 hover:scale-110 outline-none shadow-xs',
          activeBeadIndex === 0 ? 'ring-4 ring-[#72383D]/30 scale-110' : ''
        ]"
      >
        <div :class="['w-3 h-3 rounded-full', activeBeadIndex === 0 ? 'bg-[#72383D]' : 'bg-[#72383D]/40']"></div>
      </button>

      <!-- Divider line -->
      <div class="w-4 h-0.5 bg-[#D1C7BD]"></div>

      <!-- 3 Small Beads (Hail Marys) -->
      <div class="flex items-center gap-2.5">
        <button 
          v-for="i in 3" 
          :key="i"
          @click="$emit('selectBead', i)"
          :title="'Hail Mary ' + i + ' of 3'"
          :class="[
            'w-5 h-5 rounded-full border-2 border-[#AC9C8D] bg-white flex items-center justify-center cursor-pointer transition-all duration-300 hover:scale-110 outline-none shadow-xs',
            activeBeadIndex === i ? 'ring-4 ring-[#AC9C8D]/30 scale-110' : ''
          ]"
        >
          <div :class="['w-2 h-2 rounded-full', activeBeadIndex >= i ? 'bg-[#AC9C8D]' : 'bg-[#AC9C8D]/30']"></div>
        </button>
      </div>
    </div>

    <!-- Pendant Beads: 4 Our Fathers -->
    <div v-else-if="currentStep.section === 'pendant'" class="flex justify-center items-center gap-3 h-10">
      <button 
        v-for="i in 4" 
        :key="i"
        @click="$emit('selectBead', i - 1)"
        :title="'Pendant Bead ' + i + ' of 4'"
        :class="[
          'w-6 h-6 rounded-full border-2 border-[#72383D] bg-white flex items-center justify-center cursor-pointer transition-all duration-300 hover:scale-110 outline-none shadow-xs',
          activeBeadIndex === i - 1 ? 'ring-4 ring-[#72383D]/30 scale-110' : ''
        ]"
      >
        <div :class="['w-2.5 h-2.5 rounded-full', activeBeadIndex >= i - 1 ? 'bg-[#72383D]' : 'bg-[#72383D]/30']"></div>
      </button>
    </div>

    <span class="text-[10px] text-[#322D29]/60 font-bold uppercase tracking-[0.2em] select-none">
      Tap a bead to navigate
    </span>
  </div>
  <div v-else class="h-4"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { StMichaelStep } from '../../data/stMichaelData';

const props = defineProps<{
  currentStep: StMichaelStep;
  activeBeadIndex: number;
}>();

defineEmits<{
  (e: 'selectBead', beadIndex: number): void;
}>();

const shouldShow = computed(() => {
  return props.currentStep.section === 'choirs' || props.currentStep.section === 'pendant';
});
</script>
