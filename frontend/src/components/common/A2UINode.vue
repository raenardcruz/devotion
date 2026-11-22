<template>
  <template v-if="nodeEntry">
    <!-- Column Layout -->
    <div v-if="componentType === 'Column'" :class="columnClasses">
      <A2UINode
        v-for="childId in childIds"
        :key="childId"
        :node-id="childId"
        :surface="surface"
        :theme="theme"
        :visited="nextVisited"
        @citation-click="emit('citation-click', $event)"
      />
    </div>

    <!-- Row Layout -->
    <div v-else-if="componentType === 'Row'" :class="rowClasses">
      <A2UINode
        v-for="childId in childIds"
        :key="childId"
        :node-id="childId"
        :surface="surface"
        :theme="theme"
        :visited="nextVisited"
        @citation-click="emit('citation-click', $event)"
      />
    </div>

    <!-- Responsive Grid Layout (Comparisons, Dual Lists, Four Senses) -->
    <div v-else-if="componentType === 'Grid'" :class="gridClasses">
      <A2UINode
        v-for="childId in childIds"
        :key="childId"
        :node-id="childId"
        :surface="surface"
        :theme="theme"
        :visited="nextVisited"
        @citation-click="emit('citation-click', $event)"
      />
    </div>

    <!-- Card Layout -->
    <div v-else-if="componentType === 'Card'" :class="cardClasses">
      <div v-if="cardTitle || cardSubtitle" class="flex items-start gap-2.5 mb-2.5">
        <A2UIIcon
          v-if="componentProps.icon"
          :name="componentProps.icon"
          size="sm"
          :framed="true"
          :theme="theme"
          class="mt-0.5"
        />
        <span
          v-else
          :class="theme === 'kindle' ? 'w-1.5 h-1.5 bg-black inline-block shrink-0 mt-1.5' : 'w-1.5 h-1.5 rounded-full bg-[#72383D] inline-block shrink-0 mt-1.5'"
        ></span>
        <div class="min-w-0 flex-1">
          <h4
            v-if="cardTitle"
            :class="theme === 'kindle'
              ? 'text-xs font-bold uppercase tracking-wider text-black font-sans leading-snug'
              : 'text-xs font-bold uppercase tracking-wider text-[#72383D] font-serif leading-snug'"
          >
            {{ cardTitle }}
          </h4>
          <p
            v-if="cardSubtitle"
            :class="theme === 'kindle'
              ? 'text-[11px] text-black/80 mt-0.5 font-sans'
              : 'text-[11px] text-[#322D29]/65 mt-0.5 font-serif'"
          >
            {{ cardSubtitle }}
          </p>
        </div>
      </div>
      <div class="space-y-2.5">
        <A2UINode
          v-for="childId in childIds"
          :key="childId"
          :node-id="childId"
          :surface="surface"
          :theme="theme"
          :visited="nextVisited"
          @citation-click="emit('citation-click', $event)"
        />
      </div>
    </div>

    <!-- Dedicated Quote Component (Scripture, Catechism, Councils, Saints) -->
    <figure v-else-if="componentType === 'Quote'" :class="quoteClasses">
      <div class="flex items-start gap-2.5">
        <A2UIIcon
          :name="componentProps.icon || 'quote'"
          size="sm"
          :tone="theme === 'kindle' ? 'inherit' : 'primary'"
          :theme="theme"
          class="mt-0.5 opacity-80"
        />
        <div class="min-w-0 flex-1 space-y-2">
          <blockquote
            :class="theme === 'kindle'
              ? 'italic text-sm leading-relaxed text-black font-serif'
              : 'italic text-sm sm:text-[15px] leading-relaxed text-[#322D29]/90 font-serif'"
            @click="handleInlineClick"
            v-html="formattedQuoteHtml"
          ></blockquote>
          <figcaption
            v-if="quoteAttribution || quoteSource"
            :class="theme === 'kindle'
              ? 'text-xs text-black font-sans not-italic pt-1 border-t border-black/25'
              : 'text-xs text-[#72383D] font-serif not-italic flex flex-wrap items-center gap-x-1.5 gap-y-0.5 pt-1'"
          >
            <span v-if="quoteAttribution" class="font-semibold">&mdash; {{ quoteAttribution }}</span>
            <span v-if="quoteAttribution && quoteSource" class="opacity-50">&middot;</span>
            <span
              v-if="quoteSource"
              :class="theme === 'kindle' ? 'italic' : 'text-[#322D29]/70'"
              @click="handleInlineClick"
              v-html="formatInlineMarkdown(quoteSource)"
            ></span>
          </figcaption>
        </div>
      </div>
    </figure>

    <!-- Accordion / Expandable Detail Section -->
    <div v-else-if="componentType === 'Accordion'" :class="accordionContainerClasses">
      <!-- Kindle renders always expanded without interactive chrome -->
      <template v-if="theme === 'kindle'">
        <div class="flex items-center gap-2 pb-1.5 mb-2 border-b border-black">
          <A2UIIcon v-if="componentProps.icon" :name="componentProps.icon" size="sm" theme="kindle" />
          <span class="text-xs font-bold uppercase tracking-wider text-black font-sans">
            {{ cardTitle || 'Details' }}
          </span>
        </div>
        <div class="space-y-2">
          <A2UINode
            v-for="childId in childIds"
            :key="childId"
            :node-id="childId"
            :surface="surface"
            :theme="theme"
            :visited="nextVisited"
            @citation-click="emit('citation-click', $event)"
          />
        </div>
      </template>
      <template v-else>
        <button
          type="button"
          @click="isAccordionOpen = !isAccordionOpen"
          class="a2ui-btn w-full flex items-center justify-between gap-3 px-4 py-3 text-left transition-colors hover:bg-[#EFE9E1]/55 cursor-pointer"
          :aria-expanded="isAccordionOpen"
        >
          <div class="flex items-center gap-2.5 min-w-0">
            <A2UIIcon
              v-if="componentProps.icon"
              :name="componentProps.icon"
              size="sm"
              :framed="true"
              theme="default"
            />
            <div class="min-w-0">
              <div class="text-xs font-bold uppercase tracking-wider text-[#72383D] font-serif leading-snug">
                {{ cardTitle || 'More Context' }}
              </div>
              <div v-if="cardSubtitle" class="text-[11px] text-[#322D29]/65 font-serif mt-0.5 truncate">
                {{ cardSubtitle }}
              </div>
            </div>
          </div>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-4 h-4 shrink-0 text-[#72383D] transition-transform duration-200"
            :class="{ 'rotate-180': isAccordionOpen }"
          >
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>
        <div
          v-if="isAccordionOpen"
          class="px-4 pb-4 pt-1 border-t border-[#D1C7BD]/55 space-y-2.5"
        >
          <A2UINode
            v-for="childId in childIds"
            :key="childId"
            :node-id="childId"
            :surface="surface"
            :theme="theme"
            :visited="nextVisited"
            @citation-click="emit('citation-click', $event)"
          />
        </div>
      </template>
    </div>

    <!-- KeyValue Metadata / Doctrinal Distinctions Strip -->
    <dl v-else-if="componentType === 'KeyValue' && keyValueItems.length > 0" :class="keyValueClasses">
      <div
        v-for="(item, idx) in keyValueItems"
        :key="idx"
        :class="theme === 'kindle'
          ? 'py-1.5 border-b border-black/30 last:border-b-0 flex flex-col gap-0.5'
          : 'px-3.5 py-2.5 rounded-xl bg-[#FAF7F2]/90 border border-[#D1C7BD]/65 flex flex-col gap-1'"
      >
        <dt
          :class="theme === 'kindle'
            ? 'text-[11px] font-bold uppercase tracking-wider text-black font-sans flex items-center gap-1.5'
            : 'text-[10px] font-bold uppercase tracking-wider text-[#72383D] font-sans flex items-center gap-1.5'"
        >
          <A2UIIcon v-if="item.icon" :name="item.icon" size="xs" :theme="theme" />
          <span>{{ resolveText(item.label) }}</span>
        </dt>
        <dd
          :class="theme === 'kindle'
            ? 'text-sm text-black font-serif leading-snug'
            : 'text-xs sm:text-[13px] text-[#322D29]/90 font-serif leading-relaxed'"
          @click="handleInlineClick"
          v-html="formatInlineMarkdown(resolveText(item.value))"
        ></dd>
      </div>
    </dl>

    <!-- List Layout -->
    <component
      v-else-if="componentType === 'List'"
      :is="componentProps.ordered ? 'ol' : 'ul'"
      :class="[
        'space-y-2 my-1 pl-5',
        componentProps.ordered
          ? 'list-decimal'
          : theme === 'kindle'
          ? 'list-disc marker:text-black'
          : 'list-disc marker:text-[#72383D]'
      ]"
    >
      <li
        v-for="childId in childIds"
        :key="childId"
        :class="theme === 'kindle' ? 'text-sm leading-relaxed text-black pl-1' : 'text-sm leading-relaxed text-[#322D29]/90 pl-1'"
      >
        <A2UINode
          :node-id="childId"
          :surface="surface"
          :theme="theme"
          :visited="nextVisited"
          @citation-click="emit('citation-click', $event)"
        />
      </li>
    </component>

    <!-- Divider -->
    <div
      v-else-if="componentType === 'Divider' && componentProps.variant === 'ornament'"
      class="flex items-center justify-center gap-3 my-2"
      aria-hidden="true"
    >
      <span :class="theme === 'kindle' ? 'h-px flex-1 bg-black/40' : 'h-px flex-1 bg-[#D1C7BD]/70'"></span>
      <A2UIIcon name="cross" size="xs" :tone="theme === 'kindle' ? 'inherit' : 'gold'" :theme="theme" />
      <span :class="theme === 'kindle' ? 'h-px flex-1 bg-black/40' : 'h-px flex-1 bg-[#D1C7BD]/70'"></span>
    </div>
    <hr
      v-else-if="componentType === 'Divider'"
      :class="theme === 'kindle' ? 'border-t border-black my-2' : 'border-t border-[#D1C7BD]/70 my-2'"
    />

    <!-- Tabs -->
    <div v-else-if="componentType === 'Tabs'" class="space-y-3">
      <!-- Kindle e-ink mode renders every tab sequentially so nothing is hidden behind taps -->
      <template v-if="theme === 'kindle'">
        <div
          v-for="(tab, idx) in tabItems"
          :key="idx"
          class="border-l-2 border-black pl-3 py-1 space-y-1.5"
        >
          <div class="flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-black font-sans">
            <A2UIIcon v-if="tab.icon" :name="tab.icon" size="xs" theme="kindle" />
            <span>{{ resolveText(tab.title) || `Section ${idx + 1}` }}</span>
          </div>
          <A2UINode
            v-if="tab.child"
            :node-id="tab.child"
            :surface="surface"
            :theme="theme"
            :visited="nextVisited"
            @citation-click="emit('citation-click', $event)"
          />
        </div>
      </template>
      <template v-else>
        <div
          class="flex flex-wrap gap-1.5 border-b border-[#D1C7BD]/70 pb-2.5"
          role="tablist"
        >
          <button
            v-for="(tab, idx) in tabItems"
            :key="idx"
            type="button"
            role="tab"
            :aria-selected="activeTabIndex === idx"
            @click="activeTabIndex = idx"
            :class="[
              'a2ui-btn inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-serif font-semibold transition-all duration-150 cursor-pointer border',
              activeTabIndex === idx
                ? 'bg-[#72383D] text-[#FAF7F2] border-[#72383D] shadow-xs'
                : 'bg-[#EFE9E1]/80 text-[#322D29]/80 border-[#D1C7BD]/70 hover:bg-[#EFE9E1] hover:text-[#72383D]'
            ]"
          >
            <A2UIIcon
              v-if="tab.icon"
              :name="tab.icon"
              size="xs"
              :tone="activeTabIndex === idx ? 'inherit' : 'auto'"
              theme="default"
            />
            <span>{{ resolveText(tab.title) || `Tab ${idx + 1}` }}</span>
          </button>
        </div>
        <div v-if="activeTabChildId" role="tabpanel">
          <A2UINode
            :node-id="activeTabChildId"
            :surface="surface"
            :theme="theme"
            :visited="nextVisited"
            @citation-click="emit('citation-click', $event)"
          />
        </div>
      </template>
    </div>

    <!-- Text Display -->
    <template v-else-if="componentType === 'Text'">
      <h1
        v-if="usageHint === 'h1'"
        :class="theme === 'kindle'
          ? 'text-lg font-bold font-serif text-black leading-snug pb-1 border-b border-black flex items-center gap-2'
          : 'text-lg sm:text-xl font-bold font-serif text-[#4A2525] leading-snug pb-1.5 border-b border-[#D1C7BD]/60 flex items-center gap-2.5'"
      >
        <A2UIIcon v-if="componentProps.icon" :name="componentProps.icon" size="md" :framed="true" :theme="theme" />
        <span @click="handleInlineClick" v-html="formattedTextHtml"></span>
      </h1>
      <h2
        v-else-if="usageHint === 'h2'"
        :class="theme === 'kindle'
          ? 'text-base font-bold font-serif text-black leading-snug mt-1 flex items-center gap-2'
          : 'text-base sm:text-lg font-bold font-serif text-[#72383D] leading-snug mt-0.5 flex items-center gap-2'"
      >
        <A2UIIcon v-if="componentProps.icon" :name="componentProps.icon" size="sm" :framed="true" :theme="theme" />
        <span @click="handleInlineClick" v-html="formattedTextHtml"></span>
      </h2>
      <h3
        v-else-if="usageHint === 'h3' || usageHint === 'h4' || usageHint === 'h5'"
        :class="theme === 'kindle'
          ? 'text-sm font-bold font-serif text-black leading-snug flex items-center gap-1.5'
          : 'text-sm font-bold font-serif text-[#4A2525] leading-snug flex items-center gap-1.5'"
      >
        <A2UIIcon v-if="componentProps.icon" :name="componentProps.icon" size="xs" :theme="theme" />
        <span @click="handleInlineClick" v-html="formattedTextHtml"></span>
      </h3>
      <blockquote
        v-else-if="usageHint === 'quote'"
        :class="theme === 'kindle'
          ? 'border-l-3 border-black pl-3 py-1 italic text-sm leading-relaxed text-black font-serif'
          : 'border-l-3 border-[#72383D] pl-3.5 py-1.5 italic text-sm leading-relaxed text-[#322D29]/85 bg-[#FAF7F2]/85 rounded-r-xl font-serif'"
        @click="handleInlineClick"
        v-html="formattedTextHtml"
      ></blockquote>
      <p
        v-else-if="usageHint === 'lead'"
        :class="theme === 'kindle'
          ? 'text-[15px] leading-relaxed text-black font-serif'
          : 'text-[15px] sm:text-base leading-relaxed text-[#322D29] font-serif'"
        @click="handleInlineClick"
        v-html="formattedTextHtml"
      ></p>
      <p
        v-else-if="usageHint === 'caption'"
        :class="theme === 'kindle'
          ? 'text-xs text-black/80 leading-relaxed font-sans'
          : 'text-xs text-[#322D29]/65 leading-relaxed font-serif'"
        @click="handleInlineClick"
        v-html="formattedTextHtml"
      ></p>
      <p
        v-else
        :class="theme === 'kindle'
          ? 'text-sm leading-relaxed text-black font-serif'
          : 'text-sm md:text-[15px] leading-relaxed text-[#322D29]/90 font-serif'"
        @click="handleInlineClick"
        v-html="formattedTextHtml"
      ></p>
    </template>

    <!-- Badge / Chip -->
    <span v-else-if="componentType === 'Badge'" :class="badgeClasses">
      <A2UIIcon
        v-if="componentProps.icon"
        :name="componentProps.icon"
        size="xs"
        :tone="componentProps.variant === 'primary' && theme !== 'kindle' ? 'inherit' : 'auto'"
        :theme="theme"
      />
      <span>{{ resolveText(componentProps.label) || resolveText(componentProps.text) }}</span>
    </span>

    <!-- Standalone Icon -->
    <A2UIIcon
      v-else-if="componentType === 'Icon'"
      :name="componentProps.name || componentProps.icon"
      :size="resolvedIconSize"
      :tone="resolvedIconTone"
      :theme="theme"
    />

    <!-- Link -->
    <a
      v-else-if="componentType === 'Link' && componentProps.url"
      :href="String(componentProps.url)"
      target="_blank"
      rel="noopener noreferrer"
      :class="theme === 'kindle'
        ? 'inline-flex items-center gap-1 text-xs font-semibold text-black underline'
        : 'inline-flex items-center gap-1.5 text-xs font-semibold text-[#72383D] hover:text-[#4A2525] underline decoration-[#72383D]/40 hover:decoration-[#4A2525] underline-offset-2 transition-colors'"
    >
      <A2UIIcon v-if="componentProps.icon" :name="componentProps.icon" size="xs" :theme="theme" />
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
        :theme="theme"
        :visited="nextVisited"
        @citation-click="emit('citation-click', $event)"
      />
    </div>
  </template>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import A2UIIcon from './A2UIIcon.vue';
import {
  extractChildIds,
  resolveA2UIString,
  type A2UIComponentProps,
  type A2UIKeyValueItem,
  type A2UIStringValue,
  type A2UISurfaceState,
} from '../../utils/a2ui';
import type { A2UIIconTone } from '../../utils/a2uiIcons';

const props = withDefaults(
  defineProps<{
    nodeId: string;
    surface: A2UISurfaceState;
    theme?: 'default' | 'kindle';
    visited?: string[];
  }>(),
  {
    theme: 'default',
    visited: () => [],
  }
);

const emit = defineEmits<{
  (e: 'citation-click', citationNumber: number): void;
}>();

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

const isAccordionOpen = ref(Boolean(componentProps.value.defaultOpen));
watch(
  () => componentProps.value.defaultOpen,
  (val) => {
    if (typeof val === 'boolean') {
      isAccordionOpen.value = val;
    }
  }
);

const childIds = computed<string[]>(() => extractChildIds(componentProps.value));

const resolveText = (val: A2UIStringValue | string | undefined) => {
  return resolveA2UIString(val, props.surface.dataModel);
};

const cardTitle = computed(() => resolveText(componentProps.value.title));
const cardSubtitle = computed(() => resolveText(componentProps.value.subtitle));
const cardVariant = computed(() => componentProps.value.variant || 'default');
const usageHint = computed(() => componentProps.value.usageHint || componentProps.value.variant || 'body');

const quoteAttribution = computed(() => resolveText(componentProps.value.attribution));
const quoteSource = computed(() => resolveText(componentProps.value.source || componentProps.value.ref));

const tabItems = computed(() => {
  return Array.isArray(componentProps.value.tabItems) ? componentProps.value.tabItems : [];
});

const activeTabChildId = computed(() => {
  const item = tabItems.value[activeTabIndex.value];
  return item?.child || '';
});

const keyValueItems = computed<A2UIKeyValueItem[]>(() => {
  return Array.isArray(componentProps.value.items) ? componentProps.value.items : [];
});

const columnClasses = computed(() => {
  const gap = componentProps.value.gap;
  if (gap === 'sm') return 'flex flex-col gap-2';
  if (gap === 'lg') return 'flex flex-col gap-4';
  return 'flex flex-col gap-3';
});

const rowClasses = computed(() => {
  const gap = componentProps.value.gap;
  return gap === 'md' ? 'flex flex-wrap items-center gap-3' : 'flex flex-wrap items-center gap-2';
});

const gridClasses = computed(() => {
  const cols = Number(componentProps.value.columns) || 2;
  if (props.theme === 'kindle') {
    return 'grid grid-cols-1 gap-3';
  }
  if (cols >= 3) {
    return 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3';
  }
  return 'grid grid-cols-1 sm:grid-cols-2 gap-3';
});

const cardClasses = computed(() => {
  if (props.theme === 'kindle') {
    if (cardVariant.value === 'quote') {
      return 'border-l-4 border-black border-y border-r border-black/40 p-3 bg-white';
    }
    if (cardVariant.value === 'highlight') {
      return 'border-2 border-black p-3.5 bg-[#f5f5f5]';
    }
    return 'border border-black p-3 bg-white';
  }

  const base = 'rounded-2xl p-4 sm:p-4.5 transition-colors';
  switch (cardVariant.value) {
    case 'highlight':
      return `${base} bg-gradient-to-br from-[#FAF7F2] via-[#F7F1E7] to-[#EFE7D8]/75 border border-[#72383D]/30 shadow-2xs`;
    case 'quote':
      return `${base} bg-[#FAF7F2]/90 border-l-4 border-[#72383D] border-y border-r border-[#D1C7BD]/70`;
    case 'callout':
      return `${base} bg-[#EFE9E1]/60 border border-[#D1C7BD]/85`;
    case 'subtle':
      return `${base} bg-[#FAF7F2]/65 border border-[#D1C7BD]/60`;
    default:
      return `${base} bg-white/80 border border-[#D1C7BD]/80 shadow-2xs`;
  }
});

const quoteClasses = computed(() => {
  if (props.theme === 'kindle') {
    return 'border-l-4 border-black pl-3.5 py-2 my-1 bg-[#fafafa]';
  }
  return 'rounded-2xl bg-[#FAF7F2]/95 border border-[#D1C7BD]/75 border-l-4 border-l-[#72383D] px-4 py-3.5 my-0.5';
});

const accordionContainerClasses = computed(() => {
  if (props.theme === 'kindle') {
    return 'border border-black p-3 bg-white';
  }
  return 'rounded-2xl border border-[#D1C7BD]/80 bg-[#FAF7F2]/70 overflow-hidden';
});

const keyValueClasses = computed(() => {
  if (props.theme === 'kindle') {
    return 'border border-black px-3 py-1 bg-white';
  }
  return 'grid grid-cols-1 sm:grid-cols-2 gap-2';
});

const badgeClasses = computed(() => {
  const base = 'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold border';
  if (props.theme === 'kindle') {
    return 'inline-flex items-center gap-1 px-2 py-0.5 text-[11px] font-sans font-bold border border-black bg-[#f5f5f5] text-black';
  }
  if (componentProps.value.variant === 'primary') {
    return `${base} bg-[#72383D] text-[#FAF7F2] border-[#72383D]`;
  }
  if (componentProps.value.variant === 'gold') {
    return `${base} bg-[#F5ECD9] text-[#7A5B25] border-[#D8C299]`;
  }
  return `${base} bg-[#EFE9E1] text-[#72383D] border-[#D1C7BD]`;
});

const resolvedIconSize = computed<'xs' | 'sm' | 'md' | 'lg'>(() => {
  const s = componentProps.value.size;
  if (s === 'xs' || s === 'sm' || s === 'md' || s === 'lg') return s;
  return 'md';
});

const resolvedIconTone = computed<A2UIIconTone | 'auto'>(() => {
  const t = componentProps.value.tone;
  if (t === 'primary' || t === 'gold' || t === 'neutral') return t;
  return 'auto';
});

const formatInlineMarkdown = (raw: string): string => {
  if (!raw) return '';
  let html = raw
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');

  const codeClass =
    props.theme === 'kindle'
      ? 'border border-black px-1 py-0.5 font-mono text-xs text-black'
      : 'bg-[#EFE9E1] text-[#72383D] px-1.5 py-0.5 rounded font-mono text-xs';

  const strongClass =
    props.theme === 'kindle' ? 'font-bold text-black' : 'font-bold text-[#4A2525]';

  const linkClass =
    props.theme === 'kindle'
      ? 'text-black font-semibold underline'
      : 'text-[#72383D] font-semibold underline hover:text-[#4A2525]';

  const citeChipClass =
    props.theme === 'kindle'
      ? 'inline-flex items-center justify-center px-1 mx-0.5 text-[10px] font-mono font-bold border border-black text-black'
      : 'a2ui-cite-chip inline-flex items-center justify-center px-1.5 py-0.2 mx-0.5 text-[10px] font-mono font-bold rounded bg-[#EFE9E1] text-[#72383D] border border-[#D1C7BD] hover:bg-[#72383D] hover:text-[#FAF7F2] hover:border-[#72383D] transition-colors cursor-pointer';

  html = html.replace(/`([^`]+)`/g, `<code class="${codeClass}">$1</code>`);
  html = html.replace(/\*\*\*(.*?)\*\*\*/g, '<strong><em>$1</em></strong>');
  html = html.replace(/\*\*(.*?)\*\*/g, `<strong class="${strongClass}">$1</strong>`);
  html = html.replace(/\*(.*?)\*/g, '<em class="italic">$1</em>');
  html = html.replace(
    /\[([^\]]+)\]\((https?:\/\/[^)]+)\)/g,
    `<a href="$2" target="_blank" rel="noopener noreferrer" class="${linkClass}">$1</a>`
  );
  html = html.replace(
    /\[(\d+)\]/g,
    `<button type="button" data-a2ui-cite="$1" title="View source [$1]" class="${citeChipClass}">$1</button>`
  );
  html = html.replace(/\n/g, '<br/>');
  return html;
};

const formattedTextHtml = computed(() => {
  return formatInlineMarkdown(resolveText(componentProps.value.text));
});

const formattedQuoteHtml = computed(() => {
  const raw = resolveText(componentProps.value.text || componentProps.value.quote).trim();
  const stripped = raw.replace(/^["“]+|["”]+$/g, '');
  return `&ldquo;${formatInlineMarkdown(stripped)}&rdquo;`;
});

const handleInlineClick = (event: MouseEvent) => {
  const target = event.target as HTMLElement | null;
  const citeBtn = target?.closest<HTMLElement>('[data-a2ui-cite]');
  if (!citeBtn) return;
  const num = Number(citeBtn.getAttribute('data-a2ui-cite'));
  if (Number.isFinite(num) && num > 0) {
    event.preventDefault();
    emit('citation-click', num);
  }
};
</script>

<style scoped>
/* Neutralize any global button padding/min-height overrides cleanly */
button.a2ui-btn {
  min-height: 0;
}
:deep(button.a2ui-cite-chip) {
  min-height: 0;
  line-height: 1.35;
  vertical-align: baseline;
}
</style>
