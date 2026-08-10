<template>
  <div 
    class="bg-gradient-to-br from-white/85 via-[#EFE9E1]/75 to-[#D9D9D9]/30 border border-[#D1C7BD]/60 backdrop-blur-md rounded-3xl p-6 md:p-8 transition-all duration-300 shadow-sm"
    :class="[
      hoverEffect ? 'hover:scale-[1.01] hover:shadow-md hover:border-[#72383D]/50 cursor-pointer' : ''
    ]"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  hoverEffect?: boolean;
}>();
</script>
