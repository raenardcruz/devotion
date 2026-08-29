<template>
    <div class="text-center mb-4 flex flex-col items-center gap-1">
        <div class="flex items-center gap-2">
            <span class="text-parchment-secondary font-bold text-[10px] uppercase tracking-[0.25em]">
                {{ displayCategory }}
            </span>
            <span v-if="isShort && currentStep.type === 'decade'" class="px-2 py-0.5 bg-parchment-secondary/10 text-parchment-secondary text-[8px] font-bold uppercase tracking-wider rounded-full border border-parchment-secondary/20">
                Short Version
            </span>
        </div>
        <h2 class="text-2xl md:text-3xl font-serif text-parchment-neutral mt-1 font-medium">
            {{ displayTitle }}
        </h2>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Step } from './divineMercyData';
import prayerData from '../../data/prayers.json';

const props = defineProps<{
    currentStep: Step;
    isShort?: boolean;
    showLatin?: boolean;
}>();

const displayCategory = computed(() => {
    const category = props.currentStep.category;
    if (props.showLatin) {
        const catMap: Record<string, string> = {
            'Introduction': 'Introductio',
            'Intentions': 'Intentiones',
            'Opening Prayer': 'Oratio Initialis',
            'Foundation': 'Fundamentum',
            'The Decade': 'Decas',
            'Closing Prayer': 'Oratio Conclusiva'
        };
        return catMap[category] || category;
    }
    return category;
});

const displayTitle = computed(() => {
    const title = props.currentStep.title;
    if (props.showLatin) {
        if (title === 'Prayer Intentions') {
            return 'Intentiones Orationis';
        }
        // Find if this title is in prayers.json and has a latinName
        const prayer = prayerData.find(p => p.name === title || (p.id === 'dm-closing' && title === 'Optional Closing'));
        if (prayer && (prayer as any).latinName) {
            return (prayer as any).latinName;
        }
    }
    return title;
});
</script>
