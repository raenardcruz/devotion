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
      return 'bg-[#9333EA] text-white border-transparent hover:bg-[#7E22CE] hover:-translate-y-0.5 shadow-sm shadow-[#9333EA]/30';
    case 'inverted':
      return 'bg-black text-white border-transparent hover:bg-black/90 hover:-translate-y-0.5';
    case 'outlined':
      return 'bg-transparent text-[#7E22CE] border-[#E9D5FF]/80 hover:border-[#9333EA] hover:bg-white/50';
    case 'secondary':
    default:
      return 'bg-white/60 text-black border-[#E9D5FF]/60 hover:bg-white/80 hover:text-[#7E22CE] hover:-translate-y-0.5 backdrop-blur-md';
  }
});
</script>
