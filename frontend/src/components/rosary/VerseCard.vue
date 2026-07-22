<template>
    <div class="relative max-w-2xl mx-auto w-full">
        <!-- Parchment Card Wrapper -->
        <div class="bg-parchment-neutral-light border border-parchment-border p-6 md:p-8 rounded-[2rem] shadow-sm flex flex-col justify-between min-h-[300px] transition-all duration-300">
            
            <transition name="fade-scale" mode="out-in">
                <div :key="(currentStep.title || currentStep.content) + (showLatin ? '-la' : '-en')" class="space-y-5 w-full">
                    
                    <!-- Completion Marker Layout -->
                    <div v-if="currentStep.prayerId === 'completion-marker'" class="flex flex-col items-center justify-center text-center py-6 space-y-6">
                        <!-- Golden Glow Cross/Beads Icon -->
                        <div class="w-16 h-16 rounded-full bg-parchment-primary/10 flex items-center justify-center text-parchment-primary-dark shadow-sm ring-4 ring-parchment-primary/20 animate-pulse">
                            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M12 2v20"></path>
                                <path d="M17 5H7"></path>
                                <circle cx="12" cy="12" r="10" stroke-dasharray="3 3"></circle>
                            </svg>
                        </div>

                        <div class="space-y-2">
                            <h3 class="font-serif text-2xl md:text-3xl text-parchment-primary-dark font-medium">
                                {{ showLatin ? 'Rosarium Finitum Est' : 'Rosary Completed' }}
                            </h3>
                            <p class="text-sm md:text-base font-serif italic text-parchment-neutral/70 leading-relaxed max-w-md mx-auto">
                                "{{ showLatin ? currentStep.latin : currentStep.content }}"
                            </p>
                        </div>

                        <!-- Completion Action Buttons -->
                        <div class="flex flex-col sm:flex-row gap-3 pt-4 w-full justify-center">
                            <button 
                                @click="$emit('restart')"
                                class="px-5 py-2.5 bg-parchment-primary hover:bg-parchment-primary-dark text-white rounded-full transition-all duration-300 font-bold uppercase tracking-wider text-xs border border-transparent shadow-sm hover:scale-[1.02] cursor-pointer"
                            >
                                Pray Again
                            </button>
                            <router-link 
                                to="/"
                                class="px-5 py-2.5 bg-white border border-parchment-border hover:bg-parchment-neutral-light text-parchment-neutral rounded-full transition-all duration-300 font-bold uppercase tracking-wider text-xs shadow-sm hover:scale-[1.02] cursor-pointer text-center text-decoration-none"
                            >
                                Return Home
                            </router-link>
                        </div>
                    </div>

                    <template v-else>
                        <!-- Card Header (Type & Ref) -->
                        <div class="flex items-center justify-between border-b border-parchment-border/40 pb-3">
                            <span class="text-parchment-secondary font-bold tracking-[0.15em] text-[11px] uppercase">
                                {{ stepHeading }}
                            </span>
                        </div>
                        
                        <!-- Mystery Title -->
                        <h3 v-if="currentStep.mysteryTitle || currentStep.title" class="text-xl md:text-2xl font-serif text-parchment-neutral font-medium">
                            {{ displayMysteryTitle }}
                        </h3>

                        <!-- Mystery Description & Virtues -->
                        <div v-if="currentStep.type === 'mystery-header'" class="space-y-4 my-4 animate-fade-in-up">
                            <p class="text-base md:text-lg font-serif italic text-parchment-neutral/80 leading-relaxed">
                                {{ currentStep.description }}
                            </p>
                            <div v-if="currentStep.virtues && currentStep.virtues.length" class="flex flex-wrap items-center gap-2 pt-2 border-t border-parchment-border/10">
                                <span class="text-[10px] font-bold uppercase tracking-wider text-parchment-neutral/50">Virtues:</span>
                                <span 
                                    v-for="virtue in currentStep.virtues" 
                                    :key="virtue"
                                    class="px-3 py-1 bg-parchment-primary/10 border border-parchment-primary/20 text-parchment-primary-dark rounded-full text-xs font-semibold"
                                >
                                    {{ virtue }}
                                </span>
                            </div>
                        </div>

                        <!-- Scripture Verse (Mockup blockquote style, always English) -->
                        <div v-if="currentStep.verse" class="border-l-4 border-parchment-primary pl-4 py-1.5 my-4">
                            <p class="text-base md:text-lg font-serif italic text-parchment-neutral/80 leading-relaxed mb-2.5">
                                "{{ currentStep.verse.text }}"
                            </p>
                            <!-- Relocated & Linked Bible Citation -->
                            <div class="flex items-center">
                                <router-link
                                    :to="{ path: '/bible-study', query: { query: currentStep.verse.ref } }"
                                    class="inline-flex items-center space-x-1.5 text-parchment-primary-dark hover:text-parchment-primary transition-colors text-xs font-semibold hover:underline cursor-pointer"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                        <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
                                        <path d="M6 6h10"></path>
                                        <path d="M6 10h10"></path>
                                    </svg>
                                    <span>{{ currentStep.verse.ref }}</span>
                                </router-link>
                            </div>
                        </div>

                        <!-- Audio/Video Selector (Only if prayer has YouTube URL and text) -->
                        <div v-if="currentStep.youtube && (currentStep.content || currentStep.latin)" class="flex justify-center my-3">
                            <div class="flex items-center bg-parchment-neutral-light border border-parchment-border p-0.5 rounded-full w-fit">
                                <button 
                                    @click="handlePlayModeChange('audio')"
                                    :class="[selectedPlayMode === 'audio' ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                                    class="px-3.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                                >
                                    Audio &amp; Text
                                </button>
                                <button 
                                    @click="handlePlayModeChange('video')"
                                    :class="[selectedPlayMode === 'video' ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                                    class="px-3.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                                >
                                    Watch Video
                                </button>
                            </div>
                        </div>

                        <!-- Prayer Text / Latin / Interlinear / Video -->
                        <div v-if="currentStep.type !== 'mystery-header'" class="mt-4 pt-4 border-t border-parchment-border/20">
                            <!-- Video Player -->
                            <div v-if="selectedPlayMode === 'video'" class="relative w-full aspect-video rounded-2xl overflow-hidden border border-parchment-border shadow-md my-4">
                                <iframe
                                    class="absolute inset-0 w-full h-full"
                                    :src="youtubeEmbedUrl"
                                    title="YouTube video player"
                                    frameborder="0"
                                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                                    allowfullscreen
                                ></iframe>
                            </div>

                            <template v-else>
                                <InterlinearText 
                                    v-if="interlinearWords" 
                                    :words="interlinearWords" 
                                    latin-class="text-parchment-primary-dark font-medium"
                                />

                                <p 
                                    v-else-if="currentStep.content"
                                    class="text-base md:text-lg font-sans leading-relaxed text-parchment-neutral/90 whitespace-pre-line"
                                    :class="{ 'font-serif italic text-parchment-primary-dark': showLatin && currentStep.latin }"
                                >
                                    {{ (showLatin && currentStep.latin) ? currentStep.latin : currentStep.content }}
                                </p>
             
                                <!-- Latin Translation Reference (Hidden if interlinear is active) -->
                                <p 
                                    v-if="showLatin && currentStep.latin && !interlinearWords"
                                    class="mt-3 text-xs font-sans text-parchment-neutral/50 italic"
                                >
                                    {{ currentStep.content }}
                                </p>
                            </template>
                        </div>
                    </template>
                </div>
            </transition>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { RosaryStep } from './rosaryData';
import { getInterlinearPairs } from '../../utils/interlinearMappers';
import InterlinearText from '../common/InterlinearText.vue';
import prayerData from '../../data/prayers.json';
import { getYouTubeEmbedUrl } from '../../utils/youtube';

const props = defineProps<{
    currentStep: RosaryStep;
    showLatin?: boolean;
    mysteryTitle?: string;
}>();

const emit = defineEmits<{
    (e: 'restart'): void;
    (e: 'video-active'): void;
}>();

const selectedPlayMode = ref<'audio' | 'video'>('audio');

watch(() => props.currentStep, (newStep) => {
    if (newStep.youtube && !newStep.content && !newStep.latin) {
        selectedPlayMode.value = 'video';
        emit('video-active');
    } else {
        selectedPlayMode.value = 'audio';
    }
}, { immediate: true });

const youtubeEmbedUrl = computed(() => {
    return getYouTubeEmbedUrl(props.currentStep.youtube);
});

const handlePlayModeChange = (mode: 'audio' | 'video') => {
    selectedPlayMode.value = mode;
    if (mode === 'video') {
        emit('video-active');
    }
};

const interlinearWords = computed(() => {
    if (!props.showLatin || !props.currentStep.latin) return null;
    return getInterlinearPairs(props.currentStep.latin);
});

const MYSTERY_LATIN_MAP: Record<string, string> = {
  "The Annunciation": "Annuntiatio Beatae Mariae Virginis",
  "The Visitation": "Visitatio Beatae Mariae Virginis",
  "The Nativity": "Nativitas Domini nostri Iesu Christi",
  "The Presentation": "Praesentatio Iesu in Templo",
  "The Finding in the Temple": "Inventio Iesu in Templo",
  "The Baptism of Jesus": "Baptismus Iesu in Iordane",
  "The Wedding Feast at Cana": "Manifestatio Iesu in Nuptiis Canensibus",
  "Proclamation of the Kingdom": "Proclamatio Regni Dei",
  "The Transfiguration": "Transfiguratio Domini nostri Iesu Christi",
  "Institution of the Eucharist": "Institutio Sanctissimae Eucharistiae",
  "Agony in the Garden": "Agonia Iesu in Horto",
  "Scourging at the Pillar": "Flagellatio Iesu ad Columnam",
  "Crowning with Thorns": "Coronatio Spinis Iesu",
  "Carrying of the Cross": "Baiulatio Crucis Iesu",
  "The Crucifixion": "Crucifixio et Mors Iesu",
  "The Resurrection": "Resurrectio Domini",
  "The Ascension": "Ascensio Domini",
  "Descent of the Holy Spirit": "Descensus Spiritus Sancti",
  "Assumption of Mary": "Assumptio Beatae Mariae Virginis in Caelum",
  "Coronation of Mary": "Coronatio Beatae Mariae Virginis in Caelum"
};

const displayMysteryTitle = computed(() => {
    const rawTitle = props.currentStep.mysteryTitle || props.currentStep.title;
    if (!rawTitle) return '';
    if (props.showLatin) {
        if (MYSTERY_LATIN_MAP[rawTitle]) {
            return MYSTERY_LATIN_MAP[rawTitle];
        }
        const prayer = prayerData.find(p => p.name === rawTitle);
        if (prayer && (prayer as any).latinName) {
            return (prayer as any).latinName;
        }
    }
    return rawTitle;
});

const stepHeading = computed(() => {
    const step = props.currentStep;
    
    if (props.showLatin) {
        if (step.type === 'intro' || step.type === 'opening') {
            return 'Oratio Initialis';
        }
        if (step.type === 'closing') {
            return 'Oratio Conclusiva';
        }
        if (step.type === 'mystery-header') {
            return 'Denuntiatio Mysterii';
        }
        
        if (step.decadeNumber) {
            const ordinals = ['Prima', 'Secunda', 'Tertia', 'Quarta', 'Quinta'];
            const name = ordinals[step.decadeNumber - 1] || 'Decas';
            const latinMystery = props.mysteryTitle ? (MYSTERY_LATIN_MAP[props.mysteryTitle] || props.mysteryTitle) : '';
            const mysteryPart = latinMystery ? ` - ${latinMystery}` : '';
            if (step.type === 'decade-start') {
                return `${name} Decas${mysteryPart} - Pater Noster`;
            }
            if (step.type === 'decade-end') {
                const prayer = prayerData.find(p => p.id === step.prayerId);
                const prayerPart = prayer && (prayer as any).latinName ? (prayer as any).latinName : 'Gloria Patri';
                return `${name} Decas${mysteryPart} - ${prayerPart}`;
            }
            return `${name} Decas${mysteryPart} - Calculus ${step.beadNumber || 1}`;
        }
        
        if (step.title) {
            const prayer = prayerData.find(p => p.name === step.title);
            if (prayer && (prayer as any).latinName) {
                return (prayer as any).latinName;
            }
        }
        return 'Devotio Rosarii';
    }

    if (step.type === 'intro' || step.type === 'opening') {
        return 'Opening Prayer';
    }
    if (step.type === 'closing') {
        return 'Closing Prayer';
    }
    if (step.type === 'mystery-header') {
        return 'Mystery Announcement';
    }
    
    // Decade bead number (1-10) or decade start/end
    if (step.decadeNumber) {
        const ordinals = ['First', 'Second', 'Third', 'Fourth', 'Fifth'];
        const name = ordinals[step.decadeNumber - 1] || 'Decade';
        const mysteryPart = props.mysteryTitle ? ` - ${props.mysteryTitle}` : '';
        if (step.type === 'decade-start') {
            return `${name} Decade${mysteryPart} - Our Father`;
        }
        if (step.type === 'decade-end') {
            const prayer = prayerData.find(p => p.id === step.prayerId);
            const prayerPart = prayer ? prayer.name : 'Doxology';
            return `${name} Decade${mysteryPart} - ${prayerPart}`;
        }
        return `${name} Decade${mysteryPart} - Bead ${step.beadNumber || 1}`;
    }
    
    return step.title || 'Rosary Devotion';
});
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.4s ease;
}

.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.98);
  filter: blur(2px);
}

.fade-scale-leave-to {
  opacity: 0;
  transform: scale(1.02);
  filter: blur(2px);
}
</style>
