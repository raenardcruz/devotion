<template>
  <span 
    @click="handleClick"
    class="relative inline-block cursor-pointer text-emerald-400 hover:text-emerald-300 transition-colors border-b border-dashed border-emerald-500/50 hover:border-emerald-400 select-none"
    title="Click to view reference"
  >
    {{ reference }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Paragraph {
  id: number;
  text: string;
}

const props = defineProps<{
  reference: string;
  paragraphs: Paragraph[];
}>();

const emit = defineEmits<{
  (e: 'show-reference', paragraphs: Paragraph[]): void
}>();

const resolveParagraphs = () => {
    const ref = props.reference.trim();
    if (!ref) return [];

    // Check for range "123-125" or "123–125"
    const rangeMatch = ref.match(/^(\d+)[\-–](\d+)$/);
    if (rangeMatch && rangeMatch[1] && rangeMatch[2]) {
        const start = parseInt(rangeMatch[1]);
        const end = parseInt(rangeMatch[2]);
        if (!isNaN(start) && !isNaN(end)) {
            return props.paragraphs.filter(p => p.id >= start && p.id <= end);
        }
    }

    // Single number
    const id = parseInt(ref);
    if (!isNaN(id)) {
        const p = props.paragraphs.find(p => p.id === id);
        return p ? [p] : [];
    }

    return [];
};

const handleClick = (e: MouseEvent) => {
    e.stopPropagation();
    const resolved = resolveParagraphs();
    if (resolved.length > 0) {
        emit('show-reference', resolved);
    }
};
</script>
