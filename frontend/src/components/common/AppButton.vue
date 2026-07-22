<template>
  <button
    :class="[
      'rounded-full px-6 py-2.5 font-medium tracking-wide text-sm transition-all duration-300 transform active:scale-95 outline-none focus:ring-2 focus:ring-parchment-primary/50 flex items-center justify-center gap-2 shadow-none hover:shadow-sm border',
      variantClass,
      customClass
    ]"
    @click="$emit('click', $event)"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  variant?: 'primary' | 'secondary' | 'outlined' | 'inverted';
  customClass?: string;
}>(), {
  variant: 'secondary',
  customClass: ''
});

defineEmits<{
  (e: 'click', event: MouseEvent): void;
}>();

const variantClass = computed(() => {
  switch (props.variant) {
    case 'primary':
      return 'bg-parchment-primary-dark text-white border-transparent hover:bg-[#4a3c1d] hover:-translate-y-0.5';
    case 'inverted':
      return 'bg-parchment-neutral text-white border-transparent hover:bg-black hover:-translate-y-0.5';
    case 'outlined':
      return 'bg-transparent text-parchment-neutral border-parchment-border hover:border-parchment-primary hover:bg-parchment-neutral-light/50';
    case 'secondary':
    default:
      return 'bg-parchment-neutral-light text-parchment-neutral border-parchment-border hover:bg-[#ECE4D8] hover:-translate-y-0.5';
  }
});
</script>
