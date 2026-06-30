<template>
  <div class="flex flex-wrap gap-2 justify-center my-4">
    <button
      v-for="tab in tabs"
      :key="typeof tab === 'string' ? tab : tab.id"
      @click="selectTab(tab)"
      :class="[
        'px-5 py-2 rounded-full text-xs md:text-sm font-bold uppercase tracking-wider transition-all duration-300 border shadow-none hover:shadow-none',
        isActive(tab)
          ? 'bg-parchment-primary text-white border-transparent'
          : 'bg-parchment-bg text-parchment-neutral/60 border-parchment-border hover:bg-parchment-neutral-light hover:text-parchment-neutral'
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
