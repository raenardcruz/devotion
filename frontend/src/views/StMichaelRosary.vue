<template>
  <div class="min-h-screen text-black flex flex-col pb-28 sm:pb-20 selection:bg-[#9333EA]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 py-6 md:py-8 flex flex-col relative z-10">
      
      <!-- Sub-header & Controls Row -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between border-b border-[#E9D5FF]/60 pb-4 mb-4 gap-3.5 animate-fade-in-down">
          <div>
            <span class="text-xs font-bold tracking-widest text-[#7E22CE] uppercase block mb-1">Angelic Chaplet</span>
            <h1 class="font-serif text-2xl md:text-3xl text-black font-bold">
                Chaplet of St. Michael the Archangel
            </h1>
          </div>
          
          <!-- Controls -->
          <div class="flex flex-wrap items-center justify-between sm:justify-end gap-2.5 sm:gap-3.5">
              <!-- Latin/English Language Toggle -->
              <div class="flex items-center bg-white/60 border border-[#E9D5FF]/60 p-0.5 rounded-full backdrop-blur-md">
                  <button 
                    @click="showLatin = false"
                    :class="[!showLatin ? 'bg-[#9333EA] text-white font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none cursor-pointer"
                  >
                    EN
                  </button>
                  <button 
                    @click="showLatin = true"
                    :class="[showLatin ? 'bg-[#9333EA] text-white font-bold' : 'text-black/60 hover:text-[#7E22CE]']"
                    class="px-3.5 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none cursor-pointer"
                  >
                    LA
                  </button>
              </div>
          </div>
      </div>

      <!-- Hero Banner & Visual Artwork -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8 items-center bg-white/70 border border-[#E9D5FF]/80 p-5 rounded-3xl backdrop-blur-md shadow-sm">
        <div class="relative rounded-2xl overflow-hidden shadow-md aspect-square max-w-[280px] mx-auto md:max-w-none w-full">
          <img 
            :src="currentStep.choirIndex ? '/nine_angelic_choirs.jpg' : '/st_michael_archangel.jpg'" 
            alt="St. Michael the Archangel" 
            class="w-full h-full object-cover transition-transform duration-700 hover:scale-105"
          />
          <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent flex items-end p-3">
            <span class="text-white text-xs font-semibold tracking-wide">
              {{ currentStep.choirName || 'St. Michael the Archangel' }}
            </span>
          </div>
        </div>

        <div class="md:col-span-2 flex flex-col justify-center space-y-3">
          <span class="text-[11px] font-bold tracking-widest text-[#9333EA] uppercase">
            Devotion of the Nine Angelic Choirs
          </span>
          <h2 class="font-serif text-xl md:text-2xl font-bold text-black">
            Invocation to the Archangel & Angelic Hosts
          </h2>
          <p class="text-xs md:text-sm text-black/80 leading-relaxed">
            The Chaplet of St. Michael is a powerful prayer revealed by St. Michael the Archangel to Antonia d'Astonac. St. Michael promised that whoever honors him with these 9 salutations before Holy Communion will be escorted by an angel from each of the 9 celestial choirs.
          </p>
          <div class="flex items-center gap-2 pt-1 text-xs text-[#7E22CE] font-semibold">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
            <span>Promises: Escort of 9 Angels, Continual Guidance, Deliverance from Purgatory</span>
          </div>

          <!-- Official PDF Booklet Link -->
          <div class="pt-2">
            <a 
              href="https://www.dropbox.com/scl/fi/jqqsamecvha2ib5i4xz6o/4-x6-e-Booklet-The-Saint-Michael-Rosary-Prayer-Chaplet-1-.pdf.pdf?rlkey=516siw7pgeuy35vcw1rf1aky2&e=1&dl=0" 
              target="_blank" 
              rel="noopener noreferrer"
              class="inline-flex items-center gap-2 px-3.5 py-1.5 bg-[#9333EA]/10 border border-[#9333EA]/30 text-[#7E22CE] rounded-full hover:bg-[#9333EA]/20 transition-all text-xs font-bold"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                <polyline points="14 2 14 8 20 8"></polyline>
                <line x1="12" y1="18" x2="12" y2="12"></line>
                <line x1="9" y1="15" x2="12" y2="18"></line>
                <line x1="15" y1="15" x2="12" y2="18"></line>
              </svg>
              <span>Download Printable 4"x6" e-Booklet Guide (PDF)</span>
            </a>
          </div>
        </div>
      </div>

      <!-- Main Step View & Controls -->
      <div class="grid grid-cols-1 lg:grid-cols-5 gap-8 items-start">
        
        <!-- Left: Current Prayer Card (3/5 width) -->
        <div class="lg:col-span-3 flex flex-col space-y-4">
          <ParchmentCard class="!p-6 md:!p-8 relative min-h-[380px] flex flex-col justify-between shadow-lg border-[#E9D5FF]">
            <div>
              <div class="flex items-center justify-between border-b border-[#E9D5FF]/60 pb-3 mb-4">
                <div>
                  <span class="text-xs font-bold text-[#9333EA] uppercase tracking-wider block">
                    Step {{ currentIndex + 1 }} of {{ steps.length }}
                  </span>
                  <h3 class="font-serif text-xl md:text-2xl font-bold text-black">
                    {{ currentStep.title }}
                  </h3>
                </div>
                <span v-if="currentStep.beadText" class="px-3 py-1 bg-[#9333EA]/10 border border-[#9333EA]/30 text-[#7E22CE] text-[10px] font-bold rounded-full uppercase tracking-wider">
                  {{ currentStep.beadText }}
                </span>
              </div>

              <!-- Prayer Subtitle / Virtue -->
              <p v-if="currentStep.subtitle" class="text-xs font-semibold text-black/60 italic mb-4">
                {{ currentStep.subtitle }}
              </p>

              <!-- Prayer Content -->
              <div class="text-base md:text-lg leading-relaxed text-black/90 whitespace-pre-line font-serif py-2">
                {{ displayPrayerContent }}
              </div>
            </div>

            <!-- Card Bottom Navigation Controls -->
            <div class="flex items-center justify-between pt-6 border-t border-[#E9D5FF]/60 mt-6">
              <button 
                @click="prevStep" 
                :disabled="currentIndex === 0"
                class="px-4 py-2 rounded-xl bg-white border border-[#E9D5FF] text-black font-semibold text-xs disabled:opacity-40 disabled:cursor-not-allowed hover:bg-[#E9D5FF]/20 active:scale-95 transition-all flex items-center gap-1.5 cursor-pointer"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
                Previous
              </button>

              <div class="text-xs font-bold text-[#7E22CE]">
                {{ currentIndex + 1 }} / {{ steps.length }}
              </div>

              <button 
                @click="nextStep" 
                :disabled="currentIndex === steps.length - 1"
                class="px-4 py-2 rounded-xl bg-[#9333EA] border border-[#9333EA] text-white font-semibold text-xs disabled:opacity-40 disabled:cursor-not-allowed hover:bg-[#7E22CE] active:scale-95 transition-all flex items-center gap-1.5 cursor-pointer shadow-sm"
              >
                Next
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6 6-6"/></svg>
              </button>
            </div>
          </ParchmentCard>
        </div>

        <!-- Right: Accordion Collapsible Sequence List (2/5 width) -->
        <div class="lg:col-span-2 space-y-3">
          <div class="bg-white/70 border border-[#E9D5FF]/80 p-4 rounded-2xl backdrop-blur-md">
            <div class="flex items-center justify-between mb-3 border-b border-[#E9D5FF]/60 pb-2">
              <h4 class="font-serif font-bold text-base text-black">
                Chaplet Sections
              </h4>
              <span class="text-[10px] font-bold uppercase tracking-wider text-[#7E22CE] bg-[#9333EA]/10 px-2 py-0.5 rounded-full">
                {{ steps.length }} Steps
              </span>
            </div>

            <!-- Accordion Sections -->
            <div class="space-y-2.5 max-h-[520px] overflow-y-auto pr-1">
              <!-- Section 1: Opening Prayers -->
              <div class="border border-[#E9D5FF]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('intro')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#9333EA]/5 hover:bg-[#9333EA]/10 text-[#7E22CE] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>✝</span>
                    <span>1. Opening Prayers</span>
                  </span>
                  <span class="text-[10px] font-semibold text-black/50">
                    {{ expandedSection === 'intro' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'intro'" class="p-2 space-y-1 bg-white/40 border-t border-[#E9D5FF]/40">
                  <button 
                    v-for="step in introSteps" 
                    :key="step.id"
                    @click="selectStepByGlobalIndex(step.globalIdx)"
                    class="w-full text-left p-2 rounded-lg text-xs transition-all flex items-center justify-between cursor-pointer"
                    :class="[
                      currentIndex === step.globalIdx 
                        ? 'bg-[#9333EA] text-white font-bold shadow-xs' 
                        : 'hover:bg-white/90 text-black/80'
                    ]"
                  >
                    <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                  </button>
                </div>
              </div>

              <!-- Section 2: The 9 Celestial Choirs -->
              <div class="border border-[#E9D5FF]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('choirs')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#9333EA]/5 hover:bg-[#9333EA]/10 text-[#7E22CE] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>👼</span>
                    <span>2. The 9 Angelic Choirs</span>
                  </span>
                  <span class="text-[10px] font-semibold text-black/50">
                    {{ expandedSection === 'choirs' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'choirs'" class="p-2 space-y-1.5 bg-white/40 border-t border-[#E9D5FF]/40 max-h-[300px] overflow-y-auto">
                  <div v-for="choirGroup in choirGroups" :key="choirGroup.choirIndex" class="border border-[#E9D5FF]/50 rounded-lg p-1.5 bg-white/60">
                    <span class="text-[10px] font-bold text-[#7E22CE] uppercase tracking-wider block mb-1">
                      {{ choirGroup.choirName }}
                    </span>
                    <div class="space-y-0.5">
                      <button 
                        v-for="step in choirGroup.steps" 
                        :key="step.id"
                        @click="selectStepByGlobalIndex(step.globalIdx)"
                        class="w-full text-left p-1.5 rounded text-[11px] transition-all flex items-center justify-between cursor-pointer"
                        :class="[
                          currentIndex === step.globalIdx 
                            ? 'bg-[#9333EA] text-white font-bold shadow-xs' 
                            : 'hover:bg-[#9333EA]/10 text-black/80'
                        ]"
                      >
                        <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Section 3: Pendant Beads (4 Our Fathers) -->
              <div class="border border-[#E9D5FF]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('pendant')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#9333EA]/5 hover:bg-[#9333EA]/10 text-[#7E22CE] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>📿</span>
                    <span>3. Pendant 4 Our Fathers</span>
                  </span>
                  <span class="text-[10px] font-semibold text-black/50">
                    {{ expandedSection === 'pendant' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'pendant'" class="p-2 space-y-1 bg-white/40 border-t border-[#E9D5FF]/40">
                  <button 
                    v-for="step in pendantSteps" 
                    :key="step.id"
                    @click="selectStepByGlobalIndex(step.globalIdx)"
                    class="w-full text-left p-2 rounded-lg text-xs transition-all flex items-center justify-between cursor-pointer"
                    :class="[
                      currentIndex === step.globalIdx 
                        ? 'bg-[#9333EA] text-white font-bold shadow-xs' 
                        : 'hover:bg-white/90 text-black/80'
                    ]"
                  >
                    <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                  </button>
                </div>
              </div>

              <!-- Section 4: Concluding Prayers -->
              <div class="border border-[#E9D5FF]/70 rounded-xl overflow-hidden bg-white/60">
                <button 
                  @click="toggleSection('concluding')" 
                  class="w-full p-2.5 text-left text-xs font-bold flex items-center justify-between bg-[#9333EA]/5 hover:bg-[#9333EA]/10 text-[#7E22CE] cursor-pointer"
                >
                  <span class="flex items-center gap-1.5">
                    <span>🛡️</span>
                    <span>4. Concluding Prayers</span>
                  </span>
                  <span class="text-[10px] font-semibold text-black/50">
                    {{ expandedSection === 'concluding' ? '▲' : '▼' }}
                  </span>
                </button>

                <div v-if="expandedSection === 'concluding'" class="p-2 space-y-1 bg-white/40 border-t border-[#E9D5FF]/40">
                  <button 
                    v-for="step in concludingSteps" 
                    :key="step.id"
                    @click="selectStepByGlobalIndex(step.globalIdx)"
                    class="w-full text-left p-2 rounded-lg text-xs transition-all flex items-center justify-between cursor-pointer"
                    :class="[
                      currentIndex === step.globalIdx 
                        ? 'bg-[#9333EA] text-white font-bold shadow-xs' 
                        : 'hover:bg-white/90 text-black/80'
                    ]"
                  >
                    <span class="truncate">{{ step.globalIdx + 1 }}. {{ step.title }}</span>
                  </button>
                </div>
              </div>

            </div>
          </div>
        </div>

      </div>

    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import prayersData from '../data/prayers.json';
import { generateStMichaelSteps, type StMichaelStep } from '../data/stMichaelData';

const showLatin = ref(false);
const currentIndex = ref(0);
const expandedSection = ref<'intro' | 'choirs' | 'pendant' | 'concluding'>('intro');
const steps = ref<StMichaelStep[]>(generateStMichaelSteps());

const currentStep = computed<StMichaelStep>(() => {
  return steps.value[currentIndex.value]!;
});

// Auto-expand current section when step changes
watch(currentStep, (newStep) => {
  if (newStep && newStep.section) {
    expandedSection.value = newStep.section;
  }
}, { immediate: true });

function toggleSection(section: 'intro' | 'choirs' | 'pendant' | 'concluding') {
  expandedSection.value = expandedSection.value === section ? 'intro' : section;
}

function selectStepByGlobalIndex(index: number) {
  currentIndex.value = index;
}

const introSteps = computed(() => {
  return steps.value
    .map((step, idx) => ({ ...step, globalIdx: idx }))
    .filter(s => s.section === 'intro');
});

const pendantSteps = computed(() => {
  return steps.value
    .map((step, idx) => ({ ...step, globalIdx: idx }))
    .filter(s => s.section === 'pendant');
});

const concludingSteps = computed(() => {
  return steps.value
    .map((step, idx) => ({ ...step, globalIdx: idx }))
    .filter(s => s.section === 'concluding');
});

const choirGroups = computed(() => {
  const groups: { choirIndex: number; choirName: string; steps: Array<StMichaelStep & { globalIdx: number }> }[] = [];
  steps.value.forEach((step, idx) => {
    if (step.section === 'choirs' && step.choirIndex) {
      let group = groups.find(g => g.choirIndex === step.choirIndex);
      if (!group) {
        group = { choirIndex: step.choirIndex, choirName: step.choirName || `Choir ${step.choirIndex}`, steps: [] };
        groups.push(group);
      }
      group.steps.push({ ...step, globalIdx: idx });
    }
  });
  return groups;
});

const displayPrayerContent = computed(() => {
  const step = currentStep.value;
  if (!step) return '';
  if (step.customText && !showLatin.value) return step.customText;
  if (step.customLatin && showLatin.value) return step.customLatin;

  const found = prayersData.find((p: any) => p.id === step.prayerId);
  if (!found) return '';

  if (showLatin.value && found.latin) {
    return found.latin;
  }
  return found.default || found.latin || '';
});

const nextStep = () => {
  if (currentIndex.value < steps.value.length - 1) {
    currentIndex.value++;
  }
};

const prevStep = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--;
  }
};
</script>
