<template>
  <div 
    class="bg-parchment-neutral-light/60 border border-parchment-border rounded-3xl p-6 md:p-8 transition-all duration-300 shadow-sm"
    :class="[
      hoverEffect ? 'hover:scale-[1.01] hover:bg-parchment-neutral-light hover:shadow-md hover:border-parchment-primary/30 cursor-pointer' : ''
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
