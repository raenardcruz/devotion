<template>
  <div class="a2ui-renderer">
    <!-- Render A2UI Surface when valid root component is ready -->
    <A2UINode
      v-if="parsed.isA2UI && activeRootId && parsed.activeSurface"
      :node-id="activeRootId"
      :surface="parsed.activeSurface"
    />

    <!-- Streaming A2UI JSONL placeholder before first complete surfaceUpdate arrives -->
    <div
      v-else-if="parsed.isStreamingPayload && isStreaming"
      class="flex items-center gap-2.5 py-2 text-xs font-serif text-[#72383D]/80 italic"
    >
      <span class="inline-block w-2 h-2 rounded-full bg-[#72383D] animate-pulse"></span>
      <span>Composing structured A2UI response...</span>
    </div>

    <!-- Fallback to Markdown rendering for plain markdown or legacy cached responses -->
    <div
      v-else
      :class="fallbackClass"
      v-html="fallbackHtml"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import A2UINode from './A2UINode.vue';
import { parseA2UISurface } from '../../utils/a2ui';
import { renderMarkdown, type MarkdownOptions } from '../../utils/markdown';

const props = withDefaults(
  defineProps<{
    content: string;
    isStreaming?: boolean;
    markdownOptions?: MarkdownOptions;
    fallbackClass?: string;
  }>(),
  {
    isStreaming: false,
    fallbackClass: '',
  }
);

const parsed = computed(() => parseA2UISurface(props.content || ''));

const activeRootId = computed(() => {
  const surface = parsed.value.activeSurface;
  if (!surface) return '';
  if (surface.rootId && surface.components[surface.rootId]) {
    return surface.rootId;
  }
  if (surface.components['root']) {
    return 'root';
  }
  const firstId = Object.keys(surface.components)[0];
  return firstId || '';
});

const fallbackHtml = computed(() => {
  return renderMarkdown(props.content || '', props.markdownOptions);
});
</script>
