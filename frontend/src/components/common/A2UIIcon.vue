<template>
  <span
    v-if="iconDef"
    :class="wrapperClasses"
    :title="iconDef.label"
    :aria-label="decorative ? undefined : iconDef.label"
    :aria-hidden="decorative ? 'true' : undefined"
    role="img"
  >
    <svg
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      :class="svgClasses"
    >
      <!-- Subtle duotone wash accent in default theme (suppressed on Kindle e-ink) -->
      <g v-if="theme !== 'kindle' && iconDef.fills?.length" :class="fillClasses">
        <path
          v-for="(d, idx) in iconDef.fills"
          :key="`f-${idx}`"
          :d="d"
          fill="currentColor"
        />
      </g>
      <!-- Fine-line vector strokes -->
      <g
        stroke="currentColor"
        :stroke-width="strokeWidth"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path
          v-for="(d, idx) in iconDef.strokes"
          :key="`s-${idx}`"
          :d="d"
        />
      </g>
    </svg>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { getA2UIIconDefinition, type A2UIIconTone } from '../../utils/a2uiIcons';

const props = withDefaults(
  defineProps<{
    name?: string | null;
    size?: 'xs' | 'sm' | 'md' | 'lg';
    tone?: A2UIIconTone | 'inherit' | 'auto';
    framed?: boolean;
    theme?: 'default' | 'kindle';
    decorative?: boolean;
  }>(),
  {
    name: null,
    size: 'sm',
    tone: 'auto',
    framed: false,
    theme: 'default',
    decorative: true,
  }
);

const iconDef = computed(() => getA2UIIconDefinition(props.name));

const resolvedTone = computed<A2UIIconTone | 'inherit'>(() => {
  if (props.tone === 'auto') {
    return iconDef.value?.defaultTone || 'primary';
  }
  return props.tone;
});

const svgClasses = computed(() => {
  switch (props.size) {
    case 'xs':
      return 'w-3 h-3 shrink-0';
    case 'sm':
      return 'w-3.5 h-3.5 shrink-0';
    case 'lg':
      return 'w-5 h-5 shrink-0';
    case 'md':
    default:
      return 'w-4 h-4 shrink-0';
  }
});

const strokeWidth = computed(() => {
  if (props.theme === 'kindle') return '1.9';
  return props.size === 'xs' ? '1.9' : '1.7';
});

const fillClasses = computed(() => {
  if (resolvedTone.value === 'gold') return 'opacity-18';
  if (resolvedTone.value === 'neutral') return 'opacity-12';
  return 'opacity-14';
});

const wrapperClasses = computed(() => {
  const base = 'inline-flex items-center justify-center shrink-0 align-middle';

  if (props.theme === 'kindle') {
    if (props.framed) {
      return `${base} w-6 h-6 rounded-sm border border-black bg-white text-black`;
    }
    return `${base} text-black`;
  }

  if (props.framed) {
    const frameSize =
      props.size === 'lg'
        ? 'w-8 h-8 rounded-xl'
        : props.size === 'md'
        ? 'w-7 h-7 rounded-lg'
        : 'w-6 h-6 rounded-lg';

    if (resolvedTone.value === 'gold') {
      return `${base} ${frameSize} bg-[#F5ECD9] text-[#8B6B35] border border-[#D8C299]/80`;
    }
    if (resolvedTone.value === 'neutral') {
      return `${base} ${frameSize} bg-[#EFE9E1] text-[#4A2525] border border-[#D1C7BD]`;
    }
    return `${base} ${frameSize} bg-[#72383D]/10 text-[#72383D] border border-[#72383D]/20`;
  }

  if (resolvedTone.value === 'inherit') return base;
  if (resolvedTone.value === 'gold') return `${base} text-[#8B6B35]`;
  if (resolvedTone.value === 'neutral') return `${base} text-[#4A2525]/80`;
  return `${base} text-[#72383D]`;
});
</script>
