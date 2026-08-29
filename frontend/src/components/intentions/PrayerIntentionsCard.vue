<template>
  <div :class="[isFullscreen ? 'w-full max-w-4xl mx-auto my-auto flex-grow flex flex-col justify-center items-center px-4' : 'relative max-w-2xl mx-auto w-full']">
    <!-- Parchment Card Container -->
    <div 
      :class="[
        isFullscreen 
          ? 'bg-parchment-neutral-light border border-parchment-border p-8 md:p-12 rounded-[2.5rem] shadow-xl w-full my-auto flex flex-col justify-between' 
          : 'bg-parchment-neutral-light border border-parchment-border p-6 md:p-8 rounded-[2rem] shadow-sm flex flex-col justify-between min-h-[320px]'
      ]" 
      class="transition-all duration-300"
    >
      <div class="space-y-5 w-full">
        <!-- Card Header -->
        <div class="flex items-center justify-between border-b border-parchment-border/40 pb-3">
          <div class="flex items-center space-x-2">
            <span class="text-base">🙏</span>
            <span :class="[isFullscreen ? 'text-xs md:text-sm tracking-[0.2em]' : 'text-[11px] tracking-[0.15em]']" class="text-parchment-secondary font-bold uppercase">
              {{ showLatin ? 'Intentiones Orationis' : 'Prayer Intentions' }}
            </span>
          </div>

          <div class="flex items-center space-x-2">
            <!-- Edit Intentions Shortcut Button -->
            <button 
              @click="$emit('open-intentions-modal')"
              class="inline-flex items-center gap-1 px-3 py-1 bg-white/70 border border-parchment-border text-parchment-primary hover:bg-white active:scale-95 transition-all text-[10px] font-bold uppercase tracking-wider rounded-full cursor-pointer"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"></path>
                <path d="m15 5 4 4"></path>
              </svg>
              <span>Edit</span>
            </button>

            <!-- Fullscreen Toggle -->
            <button 
              @click="$emit('toggle-fullscreen')"
              class="text-parchment-neutral/50 hover:text-parchment-primary transition-colors p-1.5 rounded-full hover:bg-parchment-primary/10 cursor-pointer"
              :title="isFullscreen ? 'Exit Fullscreen Card (ESC)' : 'Fullscreen Card Presentation (F)'"
            >
              <svg v-if="!isFullscreen" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="15 3 21 3 21 9"></polyline>
                <polyline points="9 21 3 21 3 15"></polyline>
                <line x1="21" y1="3" x2="14" y2="10"></line>
                <line x1="3" y1="21" x2="10" y2="14"></line>
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="4 14 10 14 10 20"></polyline>
                <polyline points="20 10 14 10 14 4"></polyline>
                <line x1="14" y1="10" x2="21" y2="3"></line>
                <line x1="10" y1="14" x2="3" y2="21"></line>
              </svg>
            </button>
          </div>
        </div>

        <!-- Sacred Title & Invocation -->
        <div class="text-center space-y-1.5 py-1">
          <h3 :class="[isFullscreen ? 'text-2xl md:text-3xl lg:text-4xl font-bold' : 'text-xl md:text-2xl font-medium']" class="font-serif text-parchment-neutral">
            {{ showLatin ? 'Oblatio Intentionum' : 'Offering of Prayer Intentions' }}
          </h3>
          <p :class="[isFullscreen ? 'text-sm md:text-base' : 'text-xs md:text-sm']" class="font-serif italic text-parchment-neutral/70 max-w-xl mx-auto leading-relaxed">
            {{ showLatin 
              ? 'Has intentiones ante Thronum Gratiae deponimus et hanc sanctam devotionem Domino nostro Iesu Christo offerimus:' 
              : 'We place these petitions before the Throne of Grace and offer this holy devotion to Our Lord:' 
            }}
          </p>
        </div>

        <!-- Intentions Tiles List -->
        <div class="space-y-2.5 max-h-[380px] overflow-y-auto pr-1">
          <div 
            v-for="item in intentions" 
            :key="item.id"
            class="bg-white/70 border border-parchment-border/80 rounded-2xl p-3.5 flex items-start space-x-3.5 shadow-2xs transition-all hover:bg-white/90"
          >
            <!-- Icon Badge -->
            <div class="w-9 h-9 rounded-xl bg-parchment-primary/10 border border-parchment-primary/20 flex items-center justify-center text-base shrink-0 shadow-2xs">
              {{ item.icon || '🙏' }}
            </div>

            <!-- Content -->
            <div class="flex-grow min-w-0">
              <div class="flex items-center gap-2">
                <h4 class="font-serif text-sm md:text-base font-bold text-parchment-neutral">
                  {{ item.title }}
                </h4>
                <span v-if="getCategoryName(item.category)" class="text-[9px] font-bold uppercase tracking-wider text-parchment-primary-dark bg-parchment-primary/10 px-2 py-0.5 rounded-full shrink-0">
                  {{ getCategoryName(item.category) }}
                </span>
              </div>

              <p v-if="item.description" class="text-xs text-parchment-neutral/80 mt-1 leading-relaxed whitespace-pre-line font-sans">
                {{ item.description }}
              </p>
            </div>
          </div>
        </div>

        <!-- Footer Guidance -->
        <div class="pt-2 text-center">
          <p class="text-[11px] font-sans text-parchment-neutral/50">
            {{ showLatin 
              ? 'Domine, exaudi orationes nostras et clamor noster ad te veniat.' 
              : 'Lord, hear our prayer, and let our cry come unto You.' 
            }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PrayerIntention } from '../../composables/usePrayerIntentions';
import { PRAYER_CATEGORIES } from '../../composables/usePrayerIntentions';

defineProps<{
  intentions: PrayerIntention[];
  showLatin?: boolean;
  isFullscreen?: boolean;
}>();

defineEmits<{
  (e: 'toggle-fullscreen'): void;
  (e: 'open-intentions-modal'): void;
}>();

function getCategoryName(categoryId: string) {
  const found = PRAYER_CATEGORIES.find(c => c.id === categoryId);
  return found ? found.name : '';
}
</script>
