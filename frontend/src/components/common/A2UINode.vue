<template>
  <template v-if="nodeEntry">
    <!-- Column Layout -->
    <div v-if="componentType === 'Column'" class="flex flex-col gap-3">
      <A2UINode
        v-for="childId in childIds"
        :key="childId"
        :node-id="childId"
        :surface="surface"
        :visited="nextVisited"
      />
    </div>

    <!-- Row Layout -->
    <div v-else-if="componentType === 'Row'" class="flex flex-wrap items-center gap-2">
      <A2UINode
        v-for="childId in childIds"
        :key="childId"
        :node-id="childId"
        :surface="surface"
        :visited="nextVisited"
      />
    </div>

    <!-- Card Layout -->
    <div
      v-else-if="componentType === 'Card'"
      :class="[
        'rounded-2xl p-4 sm:p-5 transition-colors',
        cardVariant === 'highlight'
          ? 'bg-[#FAF7F2] border border-[#72383D]/30 shadow-xs'
          : cardVariant === 'quote'
          ? 'bg-[#FAF7F2]/90 border-l-4 border-[#72383D] border-y border-r border-[#D1C7BD]/70'
          : cardVariant === 'callout'
          ? 'bg-[#EFE9E1]/65 border border-[#D1C7BD]'
          : 'bg-white/75 border border-[#D1C7BD]/80 shadow-xs'
      ]"
    >
      <h4
        v-if="cardTitle"
        class="text-xs font-bold uppercase tracking-wider text-[#72383D] font-serif mb-2.5 flex items-center gap-1.5"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-[#72383D] inline-block shrink-0"></span>
        <span>{{ cardTitle }}</span>
      </h4>
      <div class="space-y-2.5">
        <A2UINode
          v-for="childId in childIds"
          :key="childId"
          :node-id="childId"
          :surface="surface"
          :visited="nextVisited"
        />
      </div>
    </div>

    <!-- List Layout -->
    <component
      v-else-if="componentType === 'List'"
      :is="componentProps.ordered ? 'ol' : 'ul'"
      :class="[
        'space-y-2 my-1 pl-5',
        componentProps.ordered ? 'list-decimal' : 'list-disc marker:text-[#72383D]'
      ]"
    >
      <li
        v-for="childId in childIds"
        :key="childId"
        class="text-sm leading-relaxed text-[#322D29]/90 pl-1"
      >
        <A2UINode
          :node-id="childId"
          :surface="surface"
          :visited="nextVisited"
        />
      </li>
    </component>

    <!-- Divider -->
    <hr
      v-else-if="componentType === 'Divider'"
      class="border-t border-[#D1C7BD]/70 my-2"
    />

    <!-- Tabs -->
    <div v-else-if="componentType === 'Tabs'" class="space-y-3">
      <div class="flex flex-wrap gap-1.5 border-b border-[#D1C7BD]/70 pb-2">
        <button
          v-for="(tab, idx) in tabItems"
          :key="idx"
          type="button"
          @click="activeTabIndex = idx"
          :class="[
            'px-3 py-1.5 rounded-xl text-xs font-serif font-semibold transition-colors cursor-pointer',
            activeTabIndex === idx
              ? 'bg-[#72383D] text-[#FAF7F2] shadow-xs'
              : 'bg-[#EFE9E1]/80 text-[#322D29]/80 hover:bg-[#EFE9E1] hover:text-[#72383D]'
          ]"
        >
          {{ resolveText(tab.title) || `Tab ${idx + 1}` }}
        </button>
      </div>
      <div v-if="activeTabChildId">
        <A2UINode
          :node-id="activeTabChildId"
          :surface="surface"
          :visited="nextVisited"
        />
      </div>
    </div>

    <!-- Text Display -->
    <template v-else-if="componentType === 'Text'">
      <h1
        v-if="usageHint === 'h1'"
        class="text-lg sm:text-xl font-bold font-serif text-[#4A2525] leading-snug pb-1 border-b border-[#D1C7BD]/60"
        v-html="formattedTextHtml"
      ></h1>
      <h2
        v-else-if="usageHint === 'h2'"
        class="text-base sm:text-lg font-bold font-serif text-[#72383D] leading-snug mt-1"
        v-html="formattedTextHtml"
      ></h2>
      <h3
        v-else-if="usageHint === 'h3' || usageHint === 'h4' || usageHint === 'h5'"
        class="text-sm font-bold font-serif text-[#4A2525] leading-snug"
        v-html="formattedTextHtml"
      ></h3>
      <blockquote
        v-else-if="usageHint === 'quote'"
        class="border-l-3 border-[#72383D] pl-3.5 py-1 italic text-sm leading-relaxed text-[#322D29]/85 bg-[#FAF7F2]/80 rounded-r-xl font-serif"
        v-html="formattedTextHtml"
      ></blockquote>
      <p
        v-else-if="usageHint === 'caption'"
        class="text-xs text-[#322D29]/65 leading-relaxed"
        v-html="formattedTextHtml"
      ></p>
      <p
        v-else
        class="text-sm md:text-[15px] leading-relaxed text-[#322D29]/90 font-serif"
        v-html="formattedTextHtml"
      ></p>
    </template>

    <!-- Badge / Chip -->
    <span
      v-else-if="componentType === 'Badge'"
      :class="[
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-[11px] font-mono font-semibold border',
        componentProps.variant === 'primary'
          ? 'bg-[#72383D] text-[#FAF7F2] border-[#72383D]'
          : 'bg-[#EFE9E1] text-[#72383D] border-[#D1C7BD]'
      ]"
    >
      {{ resolveText(componentProps.label) || resolveText(componentProps.text) }}
    </span>

    <!-- Link -->
    <a
      v-else-if="componentType === 'Link' && componentProps.url"
      :href="String(componentProps.url)"
      target="_blank"
      rel="noopener noreferrer"
      class="inline-flex items-center gap-1 text-xs font-semibold text-[#72383D] hover:text-[#4A2525] underline"
    >
      <span>{{ resolveText(componentProps.label) || componentProps.url }}</span>
      <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3 h-3 shrink-0">
        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
        <polyline points="15 3 21 3 21 9"></polyline>
        <line x1="10" y1="14" x2="21" y2="3"></line>
      </svg>
    </a>

    <!-- Fallback for unknown catalog components with children -->
    <div v-else-if="childIds.length > 0" class="space-y-2">
      <A2UINode
        v-for="childId in childIds"
        :key="childId"
        :node-id="childId"
        :surface="surface"
        :visited="nextVisited"
      />
    </div>
  </template>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import {
  resolveA2UIString,
  type A2UIComponentProps,
  type A2UIStringValue,
  type A2UISurfaceState,
} from '../../utils/a2ui';

const props = withDefaults(
  defineProps<{
    nodeId: string;
    surface: A2UISurfaceState;
    visited?: string[];
  }>(),
  {
    visited: () => [],
  }
);

const activeTabIndex = ref(0);

const isCycle = computed(() => props.visited.includes(props.nodeId));
const nextVisited = computed(() => [...props.visited, props.nodeId]);

const node = computed(() => {
  if (isCycle.value) return undefined;
  return props.surface.components[props.nodeId];
});

const nodeEntry = computed<[string, A2UIComponentProps] | null>(() => {
  if (!node.value?.component) return null;
  const entries = Object.entries(node.value.component);
  return entries[0] ?? null;
});

const componentType = computed(() => nodeEntry.value?.[0] || '');
const componentProps = computed<A2UIComponentProps>(() => nodeEntry.value?.[1] || {});

const childIds = computed<string[]>(() => {
  const list: string[] = [];
  if (componentProps.value.child && typeof componentProps.value.child === 'string') {
    list.push(componentProps.value.child);
  }
  const explicit = componentProps.value.children?.explicitList;
  if (Array.isArray(explicit)) {
    for (const id of explicit) {
      if (typeof id === 'string' && !list.includes(id)) {
        list.push(id);
      }
    }
  }
  return list;
});

const resolveText = (val: A2UIStringValue | string | undefined) => {
  return resolveA2UIString(val, props.surface.dataModel);
};

const cardTitle = computed(() => resolveText(componentProps.value.title));
const cardVariant = computed(() => componentProps.value.variant || 'default');
const usageHint = computed(() => componentProps.value.usageHint || 'body');

const tabItems = computed(() => {
  return Array.isArray(componentProps.value.tabItems) ? componentProps.value.tabItems : [];
});

const activeTabChildId = computed(() => {
  const item = tabItems.value[activeTabIndex.value];
  return item?.child || '';
});

const formatInlineMarkdown = (raw: string): string => {
  if (!raw) return '';
  let html = raw
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');

  html = html.replace(
    /`([^`]+)`/g,
    '<code class="bg-[#EFE9E1] text-[#72383D] px-1.5 py-0.5 rounded font-mono text-xs">$1</code>'
  );
  html = html.replace(/\*\*\*(.*?)\*\*\*/g, '<strong><em>$1</em></strong>');
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong class="font-bold text-[#4A2525]">$1</strong>');
  html = html.replace(/\*(.*?)\*/g, '<em class="italic">$1</em>');
  html = html.replace(
    /\[([^\]]+)\]\((https?:\/\/[^)]+)\)/g,
    '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-[#72383D] font-semibold underline hover:text-[#4A2525]">$1</a>'
  );
  html = html.replace(
    /\[(\d+)\]/g,
    '<span class="inline-flex items-center justify-center px-1.5 py-0.2 mx-0.5 text-[10px] font-mono font-bold rounded bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD]">$1</span>'
  );
  html = html.replace(/\n/g, '<br/>');
  return html;
};

const formattedTextHtml = computed(() => {
  return formatInlineMarkdown(resolveText(componentProps.value.text));
});
</script>
