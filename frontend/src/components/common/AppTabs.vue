<template>
  <div class="flex items-center gap-2 overflow-x-auto no-scrollbar max-w-full py-1.5 my-3 px-1 sm:justify-center sm:flex-wrap scroll-smooth">
    <button
      v-for="tab in tabs"
      :key="typeof tab === 'string' ? tab : tab.id"
      @click="selectTab(tab)"
      :class="[
        'px-4 sm:px-5 py-2 rounded-full text-xs md:text-sm font-bold uppercase tracking-wider transition-all duration-300 border shadow-none hover:shadow-none shrink-0 whitespace-nowrap',
        isActive(tab)
          ? 'bg-parchment-primary text-white border-transparent shadow-xs'
          : 'bg-parchment-bg text-parchment-neutral/70 border-parchment-border/80 hover:bg-parchment-neutral-light hover:text-parchment-neutral'
      ]"
    >
      {{ typeof tab === 'string' ? tab : tab.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  tabs: string[] | { id: string; label: string }[];
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}>();

const selectTab = (tab: string | { id: string; label: string }) => {
  const value = typeof tab === 'string' ? tab : tab.id;
  emit('update:modelValue', value);
  emit('change', value);
};

const isActive = (tab: string | { id: string; label: string }) => {
  const value = typeof tab === 'string' ? tab : tab.id;
  return value === props.modelValue;
};
</script>
