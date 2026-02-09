<template>
  <span class="leading-loose text-stone-300 group-hover/card:text-stone-100 transition-colors">
    <template v-for="(part, index) in parsedContent" :key="index">
      <template v-if="part.type === 'text'">{{ part.content }}</template>
      <CatechismReference 
        v-else-if="part.type === 'reference'" 
        :reference="part.content" 
        :paragraphs="paragraphs" 
        @show-reference="(p) => $emit('show-reference', p)"
      />
    </template>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import CatechismReference from './CatechismReference.vue';

interface Paragraph {
  id: number;
  text: string;
}

const props = defineProps<{
  text: string;
  paragraphs: Paragraph[];
}>();

const emit = defineEmits<{
    (e: 'show-reference', paragraphs: Paragraph[]): void
}>();

// Token types
type TokenType = 'text' | 'reference';
interface Token {
  type: TokenType;
  content: string;
}

const parsedContent = computed(() => {
    const text = props.text || '';
    const tokens: Token[] = [];
    
    // Regex to find parenthetical groups potentially containing references
    const regex = /\(([\d\s,–-]+)\)/g;
    
    let lastIndex = 0;
    let match;

    while ((match = regex.exec(text)) !== null) {
        // Add text before the match
        if (match.index > lastIndex) {
            tokens.push({ type: 'text', content: text.substring(lastIndex, match.index) });
        }

        const innerContent = match[1]; // e.g., "843, 2566"
        if (!innerContent) continue;

        // We want to preserve the parens as text, but make the numbers interactive.
        tokens.push({ type: 'text', content: '(' });

        // Split inner content by comma to find individual references
        // We need to keep the delimiters (comma + whitespace) as text
        const refs = innerContent.split(/(,)/);
        
        for (const refPart of refs) {
            if (refPart === ',') {
                tokens.push({ type: 'text', content: ',' });
            } else {
                 // Check if it's purely whitespace
                if (!refPart.trim()) {
                     tokens.push({ type: 'text', content: refPart });
                     continue;
                }

                // It likely contains a reference number, potentially surrounded by whitespace
                // Let's separate the whitespace from the reference
                const whitespaceMatch = refPart.match(/^(\s*)(.*?)(\s*)$/);
                if (whitespaceMatch) {
                    const prefix = whitespaceMatch[1];
                    const coreRef = whitespaceMatch[2];
                    const suffix = whitespaceMatch[3];

                    if (prefix) tokens.push({ type: 'text', content: prefix });
                    
                    // Validate if `coreRef` looks like a reference (digits, maybe range)
                    if (coreRef && /^[\d–-]+$/.test(coreRef)) {
                         tokens.push({ type: 'reference', content: coreRef });
                    } else if (coreRef) {
                        tokens.push({ type: 'text', content: coreRef });
                    }

                    if (suffix) tokens.push({ type: 'text', content: suffix });
                } else {
                     tokens.push({ type: 'text', content: refPart });
                }
            }
        }

        tokens.push({ type: 'text', content: ')' });

        lastIndex = regex.lastIndex;
    }

    // Add remaining text
    if (lastIndex < text.length) {
        tokens.push({ type: 'text', content: text.substring(lastIndex) });
    }

    return tokens;
});
</script>
