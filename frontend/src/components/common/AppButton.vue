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
      return 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white border-transparent hover:from-[#8B464C] hover:to-[#453E38] hover:-translate-y-0.5 shadow-sm shadow-[#72383D]/30';
    case 'inverted':
      return 'bg-[#322D29] text-white border-transparent hover:bg-[#322D29]/90 hover:-translate-y-0.5';
    case 'outlined':
      return 'bg-transparent text-[#72383D] border-[#D1C7BD] hover:border-[#72383D] hover:bg-[#EFE9E1]/60';
    case 'secondary':
    default:
      return 'bg-white/70 text-[#322D29] border-[#D1C7BD]/80 hover:bg-white hover:text-[#72383D] hover:border-[#72383D]/60 hover:-translate-y-0.5 backdrop-blur-md';
  }
});
</script>
