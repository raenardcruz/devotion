<script setup lang="ts">
import { ref, computed } from 'vue';
import { renderMarkdown } from '../../utils/markdown';

export interface MagisteriumCitation {
  title?: string;
  document_title?: string;
  author?: string;
  ref?: string;
  score?: number;
  url?: string;
  text?: string;
}

const props = defineProps<{
  citations?: MagisteriumCitation[];
  rawContextHtml?: string;
}>();

const isOpen = ref(false);

const parsedCitations = computed<MagisteriumCitation[]>(() => {
  if (props.citations && props.citations.length > 0) {
    return props.citations;
  }

  // Fallback parser for legacy embedded <details> HTML/Markdown in context
  if (props.rawContextHtml && props.rawContextHtml.includes('<details')) {
    const citations: MagisteriumCitation[] = [];
    const detailsMatch = props.rawContextHtml.match(/<details[\s\S]*?<\/details>/i);
    if (detailsMatch) {
      const detailsContent = detailsMatch[0];
      // Match pattern like **1. Title** — *Author* (Ref)\n> Text
      const itemRegex = /\*\*\d+\.\s*(.*?)\*\*(?:\s*—\s*\*?(.*?)\*?)?(?:\s*\((.*?)\))?\s*\n(?:>\s*(.*?)(?=\n\n|\n\*\*|\n<\/details>|$))/gs;
      let match;
      while ((match = itemRegex.exec(detailsContent)) !== null) {
        citations.push({
          title: match[1]?.trim() || '',
          document_title: match[1]?.trim() || '',
          author: match[2]?.trim() || '',
          ref: match[3]?.trim() || '',
          text: match[4]?.trim() || '',
        });
      }
    }
    return citations;
  }
  return [];
});

const citationMarkdownOptions = {
  paragraphClass: 'text-xs leading-relaxed text-amber-950/90 mb-1.5 font-serif',
  listClass: 'list-disc pl-4 space-y-1 my-1 text-xs font-serif',
  listItemClass: 'text-xs text-amber-950/90 font-serif',
};
</script>

<template>
  <div v-if="parsedCitations.length > 0" class="mt-4 border-t border-[#D1C7BD]/60 pt-3">
    <button 
      @click="isOpen = !isOpen"
      type="button"
      class="w-full flex items-center justify-between py-2.5 px-3.5 bg-gradient-to-r from-[#EFE9E1] via-white to-[#D9D9D9]/50 hover:bg-[#EFE9E1] border border-[#D1C7BD] rounded-xl transition-all text-left group cursor-pointer shadow-xs"
    >
      <div class="flex items-center space-x-2">
        <span class="text-base">📜</span>
        <span class="text-xs font-semibold text-[#72383D] font-serif">
          Fact-Checked Magisterium Citations 
          <span class="ml-1.5 text-[10px] px-2 py-0.5 rounded-full bg-[#72383D] text-white font-sans font-bold border border-[#72383D]">
            {{ parsedCitations.length }}
          </span>
        </span>
      </div>
      <svg 
        xmlns="http://www.w3.org/2000/svg" 
        width="16" 
        height="16" 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        stroke-width="2" 
        stroke-linecap="round" 
        stroke-linejoin="round"
        class="text-[#72383D] group-hover:text-[#322D29] transition-transform duration-200"
        :class="{ 'rotate-180': isOpen }"
      >
        <polyline points="6 9 12 15 18 9"></polyline>
      </svg>
    </button>

    <transition name="expand">
      <div v-if="isOpen" class="mt-2.5 space-y-3 pl-0.5 pr-0.5">
        <div 
          v-for="(item, idx) in parsedCitations" 
          :key="idx"
          class="bg-white/80 border border-[#D1C7BD]/80 rounded-xl p-3.5 shadow-2xs space-y-1.5"
        >
          <div class="flex items-start justify-between gap-2">
            <h5 class="text-xs font-bold font-serif text-[#72383D] flex items-center gap-1.5 leading-snug">
              <span class="text-[#AC9C8D] font-mono text-[11px]">{{ idx + 1 }}.</span>
              <span>{{ item.document_title || item.title || 'Magisterial Document' }}</span>
            </h5>
            <span v-if="item.ref" class="text-[10px] font-mono px-1.5 py-0.5 bg-[#EFE9E1] text-[#72383D] rounded border border-[#D1C7BD] shrink-0 font-semibold">
              {{ item.ref }}
            </span>
          </div>

          <div v-if="item.author" class="text-[11px] text-[#AC9C8D] font-serif">
            — {{ item.author }}
          </div>

          <div v-if="item.text" class="mt-2 text-xs border-l-2 border-[#72383D] pl-3 py-1 bg-white/60 rounded-r-lg">
            <div class="prose max-w-none text-xs leading-relaxed text-[#322D29]" v-html="renderMarkdown(item.text, citationMarkdownOptions)"></div>
          </div>

          <a 
            v-if="item.url" 
            :href="item.url" 
            target="_blank" 
            rel="noopener noreferrer"
            class="inline-flex items-center gap-1 text-[10px] text-[#72383D] hover:text-[#322D29] font-medium underline mt-1"
          >
            <span>View Source Document</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="3" x2="21" y2="14"></line></svg>
          </a>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s ease-out;
  max-height: 2000px;
  opacity: 1;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  max-height: 0;
  opacity: 0;
  padding-top: 0;
  padding-bottom: 0;
  margin-top: 0;
}
</style>
