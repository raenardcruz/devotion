<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 py-6 md:py-8 flex flex-col relative z-10">
      <!-- Clean Unified Devotion Header with Collapsible Settings Menu -->
      <DevotionHeader
        :title="showLatin ? 'Corolla Septem Dolorum Beatae Mariae Virginis' : 'Chaplet of the Seven Sorrows of Mary'"
        subtitle="Mater Dolorosa • Servite Chaplet"
        v-model:showLatin="showLatin"
        v-model:autoPlay="autoPlay"
        :activeIntentionsCount="activeCount"
        :isFullscreen="isFullscreen"
        :isFullscreenSupported="isFullscreenSupported"
        @open-intentions="showIntentionsModal = true"
        @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
      />

      <!-- Quick Sorrow Jump Bar (7 Sorrows + Intro + Closing) -->
      <div class="flex items-center gap-1.5 overflow-x-auto no-scrollbar py-1.5 px-2 mb-4 scroll-smooth devotion-scroll-fade-x animate-fade-in-down">
        <button
          @click="jumpToSection('intro')"
          :class="[
            'px-3.5 py-2 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-200 border shrink-0 cursor-pointer min-h-[40px]',
            currentStep.section === 'intro'
              ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white border-transparent shadow-xs'
              : 'bg-white/70 text-[#322D29]/75 border-[#D1C7BD]/80 hover:bg-white hover:text-[#72383D]'
          ]"
        >
          Opening
        </button>

        <button
          v-for="mystery in SEVEN_SORROWS_MYSTERIES"
          :key="mystery.index"
          @click="jumpToSorrow(mystery.index)"
          :class="[
            'px-3.5 py-2 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-200 border shrink-0 cursor-pointer min-h-[40px]',
            currentStep.sorrowIndex === mystery.index
              ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white border-transparent shadow-xs'
              : 'bg-white/70 text-[#322D29]/75 border-[#D1C7BD]/80 hover:bg-white hover:text-[#72383D]'
          ]"
        >
          {{ mystery.shortTitle }}
        </button>

        <button
          @click="jumpToSection('tears')"
          :class="[
            'px-3.5 py-2 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-200 border shrink-0 cursor-pointer min-h-[40px]',
            currentStep.section === 'tears' || currentStep.section === 'concluding'
              ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white border-transparent shadow-xs'
              : 'bg-white/70 text-[#322D29]/75 border-[#D1C7BD]/80 hover:bg-white hover:text-[#72383D]'
          ]"
        >
          Closing
        </button>
      </div>

      <!-- Main 2-Column Devotion Grid -->
      <div ref="swipeContainer" class="grid grid-cols-1 lg:grid-cols-5 gap-8 items-start touch-pan-y animate-fade-in-up delay-150">
        <!-- Left Column: Active Prayer / Mystery Scripture Card & Controls (3/5 width) -->
        <div
          ref="cardContainerRef"
          :class="[
            'lg:col-span-3 flex flex-col gap-5 rounded-[2.5rem] transition-all duration-300',
            isFullscreen ? 'card-fullscreen-mode !bg-[#EFE9E1] p-6 md:p-12 overflow-y-auto justify-between border-0 shadow-none' : ''
          ]"
        >
          <transition name="fade-slide" mode="out-in">
            <div :key="currentIndex + (showLatin ? '-la' : '-en')" class="flex flex-col gap-5 w-full">
              <!-- Prayer Intentions Step -->
              <PrayerIntentionsCard
                v-if="currentStep.type === 'intentions'"
                :intentions="activeIntentions"
                :showLatin="showLatin"
                :isFullscreen="isFullscreen"
                @toggle-fullscreen="toggleFullscreen(cardContainerRef)"
                @open-intentions-modal="showIntentionsModal = true"
              />

              <!-- Completion Step -->
              <div
                v-else-if="currentStep.type === 'completed'"
                class="bg-parchment-neutral-light border border-parchment-border p-8 md:p-10 rounded-[2rem] shadow-sm flex flex-col items-center justify-center text-center space-y-6 min-h-[360px]"
              >
                <div class="w-16 h-16 rounded-full bg-[#72383D]/10 flex items-center justify-center text-[#72383D] shadow-sm ring-4 ring-[#72383D]/20 animate-pulse">
                  <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"></path>
                  </svg>
                </div>

                <div class="space-y-2">
                  <h3 class="font-serif text-3xl md:text-4xl text-[#72383D] font-bold">
                    {{ showLatin ? currentStep.latinTitle : currentStep.title }}
                  </h3>
                  <p class="text-base md:text-lg font-serif text-[#322D29]/80 leading-relaxed max-w-xl mx-auto">
                    "{{ showLatin ? currentStep.customLatin : currentStep.customText }}"
                  </p>
                </div>

                <div class="flex flex-col sm:flex-row gap-3 pt-4 w-full justify-center">
                  <button
                    @click="currentIndex = 0"
                    class="px-6 py-2.5 bg-gradient-to-r from-[#72383D] to-[#322D29] text-white rounded-full transition-all duration-300 font-bold uppercase tracking-wider text-xs border border-transparent shadow-sm hover:scale-[1.02] cursor-pointer"
                  >
                    Pray Again
                  </button>
                  <router-link
                    to="/devotions"
                    class="px-6 py-2.5 bg-white border border-[#D1C7BD] hover:bg-[#EFE9E1] text-[#322D29] rounded-full transition-all duration-300 font-bold uppercase tracking-wider text-xs shadow-sm hover:scale-[1.02] cursor-pointer text-center flex items-center justify-center"
                  >
                    Back to Devotions
                  </router-link>
                </div>
              </div>

              <!-- Standard / Mystery Header / Prayer Card -->
              <div
                v-else
                :class="[
                  isFullscreen
                    ? 'devotion-card-surface p-8 md:p-12 rounded-[2.5rem] shadow-xl flex flex-col justify-between w-full my-auto relative overflow-hidden'
                    : 'devotion-card-surface p-6 md:p-8 rounded-[2rem] flex flex-col justify-between min-h-[360px] relative overflow-hidden'
                ]"
              >
                <!-- Subtle Top Atmospheric Glow -->
                <div class="pointer-events-none absolute -top-24 left-1/2 -translate-x-1/2 w-3/4 h-40 rounded-full bg-radial from-[#72383D]/10 via-[#AC9C8D]/8 to-transparent blur-2xl" />
                <div class="space-y-4 relative z-10">
                  <!-- Card Header Bar -->
                  <div class="flex flex-wrap items-center justify-between gap-2 border-b border-[#D1C7BD]/60 pb-3">
                    <div>
                      <span class="text-[11px] tracking-[0.18em] text-[#72383D] font-bold uppercase block">
                        {{ stepCategoryHeading }}
                      </span>
                      <h2
                        :class="[isFullscreen ? 'text-2xl md:text-4xl mt-1' : 'text-xl md:text-2xl mt-0.5']"
                        class="font-serif font-bold text-[#322D29]"
                      >
                        {{ showLatin && currentStep.latinTitle ? currentStep.latinTitle : currentStep.title }}
                      </h2>
                    </div>

                    <span
                      v-if="currentStep.beadText"
                      class="px-3 py-1 bg-[#72383D]/10 border border-[#72383D]/30 text-[#72383D] text-[10px] font-bold rounded-full uppercase tracking-wider shrink-0"
                    >
                      {{ currentStep.beadText }}
                    </span>
                  </div>

                  <!-- Mystery Announcement Step (Displayed BEFORE starting each mystery) -->
                  <div v-if="currentStep.type === 'sorrow-header'" class="space-y-5 pt-1">
                    <!-- Prominent Bible Verse Callout Box before starting mystery -->
                    <div
                      v-if="currentStep.verse"
                      class="bg-white/85 border border-[#D1C7BD] border-l-4 border-l-[#72383D] rounded-2xl p-5 shadow-2xs space-y-3"
                    >
                      <div class="flex items-center justify-between">
                        <span class="inline-flex items-center gap-1.5 text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D]">
                          <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                            <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
                          </svg>
                          <span>Scripture Reading Before Mystery</span>
                        </span>
                        <router-link
                          :to="{ path: '/bible-study', query: { query: currentStep.verse.ref } }"
                          class="inline-flex items-center gap-1 text-xs font-bold text-[#72383D] hover:underline"
                        >
                          <span>{{ currentStep.verse.ref }}</span>
                          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M7 17 17 7"></path>
                            <path d="M7 7h10v10"></path>
                          </svg>
                        </router-link>
                      </div>

                      <blockquote
                        :class="[isFullscreen ? 'text-xl md:text-2xl' : 'text-base md:text-lg']"
                        class="font-serif italic text-[#322D29] leading-relaxed"
                      >
                        "{{ currentStep.verse.text }}"
                      </blockquote>
                    </div>

                    <!-- Mystery Description -->
                    <p
                      :class="[isFullscreen ? 'text-lg md:text-xl' : 'text-sm md:text-base']"
                      class="font-serif text-[#322D29]/85 leading-relaxed"
                    >
                      {{ currentStep.description }}
                    </p>

                    <!-- Spiritual Fruits / Virtues -->
                    <div v-if="currentStep.virtues && currentStep.virtues.length" class="flex flex-wrap items-center gap-2 pt-2 border-t border-[#D1C7BD]/40">
                      <span class="text-[10px] font-bold uppercase tracking-wider text-[#322D29]/60">Spiritual Fruit:</span>
                      <span
                        v-for="virtue in currentStep.virtues"
                        :key="virtue"
                        class="px-3 py-1 text-xs bg-[#72383D]/10 border border-[#72383D]/25 text-[#72383D] rounded-full font-semibold"
                      >
                        {{ virtue }}
                      </span>
                    </div>

                    <!-- Action Prompt to Begin Praying the Mystery -->
                    <div class="pt-2 flex justify-end">
                      <button
                        @click="nextStep"
                        class="px-5 py-2.5 rounded-full bg-gradient-to-r from-[#72383D] to-[#322D29] text-white text-xs font-bold uppercase tracking-wider flex items-center gap-2 shadow-sm hover:from-[#8B464C] hover:to-[#453E38] cursor-pointer"
                      >
                        <span>Begin Praying This Sorrow (1 Our Father &amp; 7 Hail Marys)</span>
                        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                          <path d="m9 18 6-6-6-6"/>
                        </svg>
                      </button>
                    </div>
                  </div>

                  <!-- Regular Prayer Step (Our Father, 7 Hail Marys, Invocation, Intro, Closing) -->
                  <div v-else class="space-y-4">
                    <!-- Scripture Verse Blockquote (for Hail Mary beads or Our Father) -->
                    <div
                      v-if="currentStep.verse || currentStep.mysteryVerse"
                      class="border-l-4 border-[#72383D] pl-4 py-1.5 bg-white/50 rounded-r-xl"
                    >
                      <p
                        :class="[isFullscreen ? 'text-lg md:text-2xl' : 'text-sm md:text-base']"
                        class="font-serif italic text-[#322D29]/85 leading-relaxed mb-1.5"
                      >
                        "{{ (currentStep.verse || currentStep.mysteryVerse)?.text }}"
                      </p>
                      <router-link
                        :to="{ path: '/bible-study', query: { query: (currentStep.verse || currentStep.mysteryVerse)?.ref } }"
                        class="inline-flex items-center space-x-1.5 text-xs text-[#72383D] hover:underline font-semibold"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
                          <path d="M6 6h10"></path>
                          <path d="M6 10h10"></path>
                        </svg>
                        <span>{{ (currentStep.verse || currentStep.mysteryVerse)?.ref }}</span>
                      </router-link>
                    </div>

                    <!-- Prayer Text / Interlinear Latin -->
                    <div class="pt-2">
                      <InterlinearText
                        v-if="interlinearWords"
                        :words="interlinearWords"
                        latin-class="text-[#72383D] font-medium"
                      />

                      <p
                        v-else
                        :class="[
                          isFullscreen
                            ? 'text-xl md:text-3xl leading-relaxed font-serif py-3 text-center'
                            : 'text-base md:text-lg font-serif leading-relaxed text-[#322D29]/95',
                          showLatin ? 'text-[#72383D]' : ''
                        ]"
                        class="whitespace-pre-line"
                      >
                        {{ displayPrayerText }}
                      </p>

                      <!-- English Reference below non-interlinear Latin -->
                      <p
                        v-if="showLatin && !interlinearWords && englishPrayerText"
                        class="mt-3 text-xs font-serif text-[#322D29]/60 italic whitespace-pre-line"
                      >
                        {{ englishPrayerText }}
                      </p>
                    </div>
                  </div>
                </div>

                <!-- Interactive Bead Tracker inside Card -->
                <SevenSorrowsBeads
                  :currentStep="currentStep"
                  @selectBead="onSelectSorrowBead"
                  @selectTearsBead="onSelectTearsBead"
                />
              </div>
            </div>
          </transition>

          <!-- Step Navigation & Audio Controls -->
          <DevotionControls
            :totalSteps="steps.length"
            :currentStepIndex="currentIndex"
            :isLastStep="currentIndex === steps.length - 1"
            :hasAudio="Boolean(currentStep.prayerId && getPrayerAudioUrl(currentStep.prayerId))"
            :isPlaying="isPlaying"
            label="Chaplet"
            @next="nextStep"
            @prev="prevStep"
            @toggle-audio="toggleAudio"
          />

          <!-- Fullscreen Keyboard Helper Hint -->
          <div v-if="isFullscreen" class="text-center text-xs text-[#322D29]/60 font-sans tracking-wide mt-2">
            Press <kbd class="px-1.5 py-0.5 bg-[#322D29]/10 rounded font-mono text-[10px]">←</kbd> <kbd class="px-1.5 py-0.5 bg-[#322D29]/10 rounded font-mono text-[10px]">→</kbd> to Navigate &bull; <kbd class="px-1.5 py-0.5 bg-[#322D29]/10 rounded font-mono text-[10px]">Space</kbd> Play/Pause &bull; <kbd class="px-1.5 py-0.5 bg-[#322D29]/10 rounded font-mono text-[10px]">ESC</kbd> Exit Fullscreen
          </div>
        </div>

        <!-- Right Column: Seamless Gradient-Fade Artwork, Meditation & Promises (2/5 width) -->
        <div class="lg:col-span-2 flex flex-col gap-5">
          <!-- Unified Artwork & Contemplative Reflection Card -->
          <div class="devotion-card-surface rounded-[2rem] overflow-hidden flex flex-col relative group">
            <div class="aspect-[4/3] sm:aspect-[16/11] lg:aspect-[4/3] w-full overflow-hidden relative bg-[#221D1A]">
              <transition name="fade-slide" mode="out-in">
                <img
                  :key="currentVisual.image"
                  :src="currentVisual.image"
                  :alt="currentVisual.title"
                  class="w-full h-full object-cover object-top devotion-artwork-fade transition-transform duration-1000 group-hover:scale-105"
                />
              </transition>
              <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-[#EFE9E1] via-[#EFE9E1]/35 to-transparent" />
            </div>

            <div class="relative z-10 -mt-12 px-6 pb-6 pt-2 space-y-3">
              <div>
                <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D]">
                  {{ currentVisual.badge }}
                </span>
                <h3 class="font-serif text-lg md:text-xl font-bold text-[#322D29] leading-snug">
                  {{ currentVisual.title }}
                </h3>
                <p class="text-[11px] text-[#322D29]/60 italic mt-0.5">
                  {{ currentVisual.artist }}
                </p>
              </div>

              <transition name="fade-slide" mode="out-in">
                <p :key="currentVisual.meditation" class="font-serif text-sm md:text-base text-[#322D29]/85 leading-relaxed border-l-2 border-[#72383D]/60 pl-3.5">
                  {{ currentVisual.meditation }}
                </p>
              </transition>
            </div>
          </div>

          <!-- Promises of Our Lady of Sorrows Card -->
          <div class="bg-white/70 border border-[#D1C7BD]/80 rounded-3xl p-5 shadow-2xs backdrop-blur-md">
            <button
              @click="showPromises = !showPromises"
              class="w-full flex items-center justify-between text-left border-none bg-transparent p-0 min-h-0 hover:bg-transparent hover:scale-100 cursor-pointer"
            >
              <div class="flex items-center gap-2">
                <span class="text-[#72383D]">⚔️</span>
                <h4 class="font-serif text-sm md:text-base font-bold text-[#72383D]">
                  Promises of Our Lady of Sorrows
                </h4>
              </div>
              <span class="text-xs text-[#72383D] font-bold">
                {{ showPromises ? 'Hide ▲' : 'Read ▼' }}
              </span>
            </button>

            <div v-if="showPromises" class="mt-3 pt-3 border-t border-[#D1C7BD]/50 space-y-2.5 text-xs text-[#322D29]/85">
              <div
                v-for="promise in SEVEN_SORROWS_PROMISES"
                :key="promise.title"
                class="bg-[#EFE9E1]/60 p-3 rounded-2xl border border-[#D1C7BD]/50"
              >
                <div class="font-bold text-[#72383D] mb-0.5">{{ promise.title }}</div>
                <p class="leading-relaxed">{{ promise.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Global Footer -->
    <BottomNav />

    <!-- Prayer Intentions Modal -->
    <PrayerIntentionsModal
      :isOpen="showIntentionsModal"
      devotionKey="seven_sorrows"
      devotionTitle="Chaplet of the Seven Sorrows"
      @close="showIntentionsModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import DevotionHeader from '../components/devotion/DevotionHeader.vue';
import DevotionControls from '../components/devotion/DevotionControls.vue';
import SevenSorrowsBeads from '../components/sevensorrows/SevenSorrowsBeads.vue';
import InterlinearText from '../components/common/InterlinearText.vue';
import PrayerIntentionsCard from '../components/intentions/PrayerIntentionsCard.vue';
import PrayerIntentionsModal from '../components/intentions/PrayerIntentionsModal.vue';
import { usePrayerIntentions } from '../composables/usePrayerIntentions';
import { useSwipe } from '../composables/useSwipe';
import { useFullscreen } from '../composables/useFullscreen';
import { getPrayerAudioUrl } from '../utils/audioHelper';
import { getInterlinearPairs } from '../utils/interlinearMappers';
import prayersData from '../data/prayers.json';
import {
  SEVEN_SORROWS_MYSTERIES,
  SEVEN_SORROWS_PROMISES,
  generateSevenSorrowsSteps,
  type SevenSorrowsStep
} from '../data/sevenSorrowsData';

const {
  activeIntentions,
  activeCount,
  hasIntentions
} = usePrayerIntentions('seven_sorrows');

const showIntentionsModal = ref(false);
const showLatin = ref(false);
const showPromises = ref(false);
const currentIndex = ref(0);

const steps = computed<SevenSorrowsStep[]>(() => generateSevenSorrowsSteps(hasIntentions.value));

watch(hasIntentions, () => {
  if (currentIndex.value >= steps.value.length) {
    currentIndex.value = 0;
  }
});

const swipeContainer = ref<HTMLElement | null>(null);
const cardContainerRef = ref<HTMLElement | null>(null);

const { isFullscreen, isFullscreenSupported, toggleFullscreen } = useFullscreen(cardContainerRef);

useSwipe(swipeContainer, {
  onSwipeLeft: () => nextStep(),
  onSwipeRight: () => prevStep()
});

const currentStep = computed<SevenSorrowsStep>(() => {
  return steps.value[currentIndex.value] || steps.value[0]!;
});

const currentMystery = computed(() => {
  const sIdx = currentStep.value.sorrowIndex;
  if (sIdx) {
    return SEVEN_SORROWS_MYSTERIES.find(m => m.index === sIdx) || null;
  }
  return null;
});

const currentVisual = computed(() => {
  if (currentMystery.value) {
    return {
      badge: showLatin.value ? currentMystery.value.latinOrdinal : currentMystery.value.ordinal,
      title: showLatin.value ? currentMystery.value.latinTitle : currentMystery.value.title,
      image: currentMystery.value.image,
      artist: currentMystery.value.artist,
      meditation: currentMystery.value.meditation
    };
  }
  return {
    badge: showLatin.value ? 'Mater Dolorosa' : 'Our Lady of Sorrows',
    title: showLatin.value ? 'Beata Maria Virgo Perdolens' : 'Behold Your Mother',
    image: '/images/seven-sorrows/mater_dolorosa.jpg',
    artist: 'El Greco — Mater Dolorosa (c. 1590)',
    meditation:
      'Meditating upon the Seven Sorrows of the Blessed Virgin Mary unites our hearts with her maternal compassion at the foot of the Cross and obtains true contrition for our sins.'
  };
});

const stepCategoryHeading = computed(() => {
  const step = currentStep.value;
  if (showLatin.value) {
    if (step.section === 'intro') return 'Oratio Initialis';
    if (step.type === 'sorrow-header') return `${ currentMystery.value?.latinOrdinal || 'Dolor' } — Lectio Sacra`;
    if (step.section === 'sorrows') return `${ currentMystery.value?.latinOrdinal || 'Dolor' } — ${ currentMystery.value?.latinTitle || '' }`;
    if (step.section === 'tears') return 'In Honorem Lacrimarum Mariae';
    return 'Oratio Conclusiva';
  }
  if (step.section === 'intro') return 'Opening Prayer';
  if (step.type === 'sorrow-header') return 'Mystery Announcement & Scripture';
  if (step.section === 'sorrows') return `${ currentMystery.value?.ordinal || 'Sorrow' } — ${ currentMystery.value?.title || '' }`;
  if (step.section === 'tears') return 'Concluding Three Hail Marys';
  return 'Concluding Prayer';
});

const foundPrayer = computed(() => {
  if (!currentStep.value.prayerId) return null;
  return prayersData.find(p => p.id === currentStep.value.prayerId) || null;
});

const englishPrayerText = computed(() => {
  if (currentStep.value.customText) return currentStep.value.customText;
  return foundPrayer.value?.default || '';
});

const latinPrayerText = computed(() => {
  if (currentStep.value.customLatin) return currentStep.value.customLatin;
  return foundPrayer.value?.latin || '';
});

const displayPrayerText = computed(() => {
  if (showLatin.value && latinPrayerText.value) {
    return latinPrayerText.value;
  }
  return englishPrayerText.value;
});

const interlinearWords = computed(() => {
  if (!showLatin.value || !latinPrayerText.value) return null;
  return getInterlinearPairs(latinPrayerText.value);
});

// Navigation Functions
function nextStep() {
  if (currentIndex.value < steps.value.length - 1) {
    currentIndex.value++;
  } else {
    currentIndex.value = 0;
  }
}

function prevStep() {
  if (currentIndex.value > 0) {
    currentIndex.value--;
  }
}

function jumpToSorrow(sorrowIndex: number) {
  const idx = steps.value.findIndex(s => s.type === 'sorrow-header' && s.sorrowIndex === sorrowIndex);
  if (idx !== -1) {
    currentIndex.value = idx;
  }
}

function jumpToSection(section: 'intro' | 'tears') {
  if (section === 'intro') {
    currentIndex.value = 0;
    return;
  }
  const idx = steps.value.findIndex(s => s.section === 'tears');
  if (idx !== -1) {
    currentIndex.value = idx;
  }
}

function onSelectSorrowBead(beadNum: number) {
  const sIdx = currentStep.value.sorrowIndex;
  if (!sIdx) return;

  if (beadNum === -1) {
    const idx = steps.value.findIndex(s => s.sorrowIndex === sIdx && s.type === 'sorrow-header');
    if (idx !== -1) currentIndex.value = idx;
    return;
  }

  if (beadNum === 0) {
    const idx = steps.value.findIndex(s => s.sorrowIndex === sIdx && s.type === 'sorrow-our-father');
    if (idx !== -1) currentIndex.value = idx;
    return;
  }

  const idx = steps.value.findIndex(
    s => s.sorrowIndex === sIdx && s.type === 'sorrow-hail-mary' && s.beadNumber === beadNum
  );
  if (idx !== -1) currentIndex.value = idx;
}

function onSelectTearsBead(beadNum: number) {
  const idx = steps.value.findIndex(s => s.section === 'tears' && s.beadNumber === beadNum);
  if (idx !== -1) currentIndex.value = idx;
}

// Audio Playback
const currentAudio = ref<HTMLAudioElement | null>(null);
const isPlaying = ref(false);
const autoPlay = ref(true);
let playTimeout: any = null;

function stopAudio() {
  if (playTimeout) {
    clearTimeout(playTimeout);
    playTimeout = null;
  }
  if (currentAudio.value) {
    currentAudio.value.pause();
    currentAudio.value = null;
  }
  isPlaying.value = false;
}

function playAudio(prayerId: string | undefined) {
  if (!prayerId) return;
  const url = getPrayerAudioUrl(prayerId);
  if (!url) return;

  stopAudio();

  const audio = new Audio(url);
  currentAudio.value = audio;

  audio.addEventListener('ended', () => {
    isPlaying.value = false;
    currentAudio.value = null;
  });

  audio.addEventListener('error', () => {
    isPlaying.value = false;
    currentAudio.value = null;
  });

  audio.play()
    .then(() => {
      isPlaying.value = true;
    })
    .catch(() => {
      isPlaying.value = false;
    });
}

function toggleAudio() {
  if (isPlaying.value) {
    stopAudio();
  } else if (currentStep.value.prayerId) {
    playAudio(currentStep.value.prayerId);
  }
}

watch(currentStep, (newStep) => {
  stopAudio();
  if (autoPlay.value && newStep.prayerId) {
    playTimeout = setTimeout(() => {
      if (newStep.prayerId === currentStep.value.prayerId) {
        playAudio(newStep.prayerId);
      }
    }, 300);
  }
}, { deep: true });

watch(autoPlay, (newVal) => {
  localStorage.setItem('seven-sorrows-autoplay', String(newVal));
  if (!newVal) {
    stopAudio();
  } else if (currentStep.value.prayerId) {
    playAudio(currentStep.value.prayerId);
  }
});

function handleKeyDown(e: KeyboardEvent) {
  if (['INPUT', 'TEXTAREA', 'SELECT'].includes((e.target as HTMLElement)?.tagName)) {
    return;
  }
  if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
    e.preventDefault();
    nextStep();
  } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
    e.preventDefault();
    prevStep();
  } else if (e.key === ' ') {
    e.preventDefault();
    toggleAudio();
  } else if (e.key === 'f' || e.key === 'F') {
    e.preventDefault();
    toggleFullscreen(cardContainerRef.value);
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown);
  const savedAutoPlay = localStorage.getItem('seven-sorrows-autoplay');
  if (savedAutoPlay !== null) {
    autoPlay.value = savedAutoPlay === 'true';
  }
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyDown);
  stopAudio();
});
</script>

<style scoped>
:fullscreen, :-webkit-full-screen {
  background-color: #FAF8F5 !important;
  width: 100vw !important;
  height: 100vh !important;
  max-width: none !important;
  max-height: none !important;
  padding: 3rem 2rem !important;
  display: flex !important;
  flex-direction: column !important;
  justify-content: space-between !important;
  align-items: center !important;
  box-sizing: border-box !important;
  overflow-y: auto !important;
}
</style>
