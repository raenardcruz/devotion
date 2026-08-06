<template>
  <div 
    class="bg-white/60 border border-[#E9D5FF]/60 backdrop-blur-md rounded-3xl p-6 md:p-8 transition-all duration-300 shadow-sm"
    :class="[
      hoverEffect ? 'hover:scale-[1.01] hover:bg-white/80 hover:shadow-md hover:border-[#9333EA]/40 cursor-pointer' : ''
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
