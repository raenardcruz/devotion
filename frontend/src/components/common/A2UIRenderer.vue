<template>
  <div class="a2ui-renderer" :class="theme === 'kindle' ? 'a2ui-theme-kindle' : 'a2ui-theme-default'">
    <!-- Render A2UI Surface when valid root component is ready -->
    <A2UINode
      v-if="parsed.isA2UI && activeRootId && parsed.activeSurface"
      :node-id="activeRootId"
      :surface="parsed.activeSurface"
      :theme="theme"
      @citation-click="emit('citation-click', $event)"
    />

    <!-- Streaming A2UI JSONL placeholder before first complete surfaceUpdate arrives -->
    <div
      v-else-if="parsed.isStreamingPayload && isStreaming"
      class="flex items-center gap-2.5 py-2 text-xs font-serif text-[#72383D]/85 italic"
    >
      <A2UIIcon name="scripture" size="sm" :framed="true" :theme="theme" class="animate-pulse" />
      <span>Composing structured response...</span>
    </div>

    <!-- Graceful prose recovery if A2UI payload was truncated or malformed -->
    <div
      v-else-if="parsed.isStreamingPayload && recoveredProseHtml"
      :class="fallbackClass"
      v-html="recoveredProseHtml"
    ></div>

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
import A2UIIcon from './A2UIIcon.vue';
import { extractA2UIFallbackProse, parseA2UISurface } from '../../utils/a2ui';
import { renderMarkdown, type MarkdownOptions } from '../../utils/markdown';

const props = withDefaults(
  defineProps<{
    content: string;
    isStreaming?: boolean;
    theme?: 'default' | 'kindle';
    markdownOptions?: MarkdownOptions;
    fallbackClass?: string;
  }>(),
  {
    isStreaming: false,
    theme: 'default',
    fallbackClass: '',
  }
);

const emit = defineEmits<{
  (e: 'citation-click', citationNumber: number): void;
}>();

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
  const firstId = surface.componentOrder[0] || Object.keys(surface.components)[0];
  return firstId || '';
});

const recoveredProseHtml = computed(() => {
  if (!parsed.value.isStreamingPayload || activeRootId.value) return '';
  const recovered = extractA2UIFallbackProse(props.content || '');
  if (!recovered) return '';
  return renderMarkdown(recovered, props.markdownOptions);
});

const fallbackHtml = computed(() => {
  return renderMarkdown(props.content || '', props.markdownOptions);
});
</script>
