<template>
  <div class="flex flex-wrap gap-x-2 gap-y-2 justify-center text-center">
    <div 
      v-for="(pair, index) in words" 
      :key="index" 
      class="flex flex-col items-center min-w-[2.5rem]"
    >
      <span 
        class="text-2xl md:text-3xl font-serif leading-none transition-all duration-500"
        :class="latinClass || 'text-parchment-neutral'"
      >
        {{ pair.la }}
      </span>
      <span 
        class="mt-1.5 text-xs font-sans leading-none text-parchment-neutral/55 tracking-wider"
      >
        {{ pair.en }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  words: { la: string; en: string }[];
  latinClass?: string;
}>();
</script>

<style scoped>
/* Ensure words don't break mid-sentence inappropriately */
.flex-wrap {
  perspective: 1000px;
}
</style>
