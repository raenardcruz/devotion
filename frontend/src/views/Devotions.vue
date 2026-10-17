<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-6xl mx-auto w-full px-4 pt-8 md:pt-12 pb-12 flex flex-col relative z-10">
      <!-- Page Hero Header -->
      <header class="text-center mb-8 md:mb-10 animate-fade-in-down">
        <div class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#72383D]/10 border border-[#72383D]/20 text-[#72383D] text-[11px] font-bold uppercase tracking-[0.2em] mb-3">
          <span>☦</span>
          <span>Sacred Prayer Sanctuary</span>
        </div>
        <h1 class="font-serif text-3xl md:text-5xl bg-gradient-to-r from-[#72383D] via-[#AC9C8D] to-[#322D29] bg-clip-text text-transparent mb-3 font-bold tracking-tight">
          Catholic Devotions &amp; Chaplets
        </h1>
        <p class="text-[#322D29]/75 text-sm md:text-base max-w-2xl mx-auto leading-relaxed font-serif">
          Enter into contemplative prayer through the sacred mysteries of Christ and His Sorrowful Mother, accompanied by Sacred Scripture, classical sacred art, and Latin tradition.
        </p>
      </header>

      <!-- Liturgical Day Recommendation Banner -->
      <div class="mb-8 animate-fade-in-down">
        <div class="bg-gradient-to-r from-[#72383D] via-[#54292D] to-[#322D29] text-white rounded-3xl p-5 md:p-6 shadow-md border border-white/15 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div class="flex items-start sm:items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-white/15 backdrop-blur-md flex items-center justify-center text-white shrink-0 border border-white/20">
              <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
              </svg>
            </div>
            <div>
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#EFE9E1]/80">
                  {{ currentDayName }} Devotional Rhythm
                </span>
                <span class="px-2.5 py-0.5 rounded-full bg-white/15 text-[10px] font-bold uppercase tracking-wider text-[#EFE9E1]">
                  {{ todayMystery }} Mysteries
                </span>
              </div>
              <h2 class="font-serif text-lg md:text-xl font-bold text-white leading-snug">
                {{ dailyRecommendation.title }}
              </h2>
              <p class="text-xs md:text-sm text-[#EFE9E1]/80 mt-0.5 leading-relaxed">
                {{ dailyRecommendation.subtitle }}
              </p>
            </div>
          </div>

          <router-link
            :to="dailyRecommendation.path"
            class="px-5 py-2.5 rounded-full bg-[#EFE9E1] hover:bg-white text-[#72383D] font-bold text-xs uppercase tracking-wider transition-all shrink-0 text-center shadow-sm active:scale-95 min-h-[44px] inline-flex items-center justify-center gap-1.5"
          >
            <span>{{ dailyRecommendation.cta }}</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="m9 18 6-6-6-6"/>
            </svg>
          </router-link>
        </div>
      </div>

      <!-- Category Filter Pills -->
      <div class="flex items-center justify-start sm:justify-center gap-2 overflow-x-auto no-scrollbar pb-2 mb-6">
        <button
          v-for="cat in categories"
          :key="cat.id"
          @click="selectedCategory = cat.id"
          :class="[
            'px-4 py-2 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-200 border shrink-0 cursor-pointer min-h-[40px]',
            selectedCategory === cat.id
              ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white border-transparent shadow-xs'
              : 'bg-white/75 text-[#322D29]/75 border-[#D1C7BD]/80 hover:bg-white hover:text-[#72383D]'
          ]"
        >
          {{ cat.label }}
        </button>
      </div>

      <!-- Devotions Tiles Grid (2x2 for Major Chaplets + Featured Card) -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 animate-fade-in-up delay-150">
        <router-link
          v-for="devotion in filteredDevotions"
          :key="devotion.path"
          :to="devotion.path"
          class="group block"
        >
          <article class="h-full bg-gradient-to-br from-white/90 via-[#EFE9E1]/85 to-[#D9D9D9]/35 border border-[#D1C7BD]/80 rounded-[2rem] overflow-hidden shadow-sm hover:shadow-xl hover:border-[#72383D]/50 transition-all duration-300 hover:-translate-y-1 flex flex-col justify-between">
            <!-- Tile Artwork Header -->
            <div class="relative h-52 sm:h-56 w-full overflow-hidden bg-[#221D1A]">
              <img
                :src="devotion.image"
                :alt="devotion.title"
                class="w-full h-full object-cover object-center transition-transform duration-700 group-hover:scale-105"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-[#221D1A]/90 via-[#221D1A]/30 to-transparent"></div>

              <!-- Top Badges Row -->
              <div class="absolute top-4 left-4 right-4 flex items-center justify-between gap-2">
                <span class="px-3 py-1 rounded-full bg-black/45 backdrop-blur-md border border-white/20 text-white text-[10px] font-bold uppercase tracking-widest">
                  {{ devotion.tag }}
                </span>
                <span
                  v-if="devotion.badge"
                  class="px-3 py-1 rounded-full bg-[#72383D]/90 backdrop-blur-md border border-white/20 text-white text-[10px] font-bold uppercase tracking-wider shadow-xs"
                >
                  {{ devotion.badge }}
                </span>
              </div>

              <!-- Artwork Bottom Overlay Title -->
              <div class="absolute bottom-4 left-5 right-5">
                <span class="text-[11px] font-serif italic text-[#EFE9E1]/85 block">
                  {{ devotion.latinTitle }}
                </span>
                <h2 class="font-serif text-2xl sm:text-3xl font-bold text-white leading-tight">
                  {{ devotion.title }}
                </h2>
              </div>
            </div>

            <!-- Tile Content Body -->
            <div class="p-5 sm:p-6 flex-1 flex flex-col justify-between space-y-5">
              <div class="space-y-3">
                <p class="text-sm text-[#322D29]/85 leading-relaxed">
                  {{ devotion.description }}
                </p>

                <!-- Feature Pills -->
                <div class="flex flex-wrap gap-1.5 pt-1">
                  <span
                    v-for="feature in devotion.features"
                    :key="feature"
                    class="px-2.5 py-1 rounded-lg bg-[#EFE9E1] border border-[#D1C7BD]/70 text-[#322D29]/80 text-[11px] font-semibold"
                  >
                    {{ feature }}
                  </span>
                </div>
              </div>

              <!-- Tile Footer CTA -->
              <div class="pt-3 border-t border-[#D1C7BD]/50 flex items-center justify-between">
                <div class="flex items-center gap-2 text-xs text-[#322D29]/65 font-semibold">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]">
                    <circle cx="12" cy="12" r="10"></circle>
                    <polyline points="12 6 12 12 16 14"></polyline>
                  </svg>
                  <span>{{ devotion.duration }}</span>
                </div>

                <span class="inline-flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-[#72383D] group-hover:translate-x-1 transition-transform">
                  <span>Enter Devotion</span>
                  <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M5 12h14"></path>
                    <path d="m12 5 7 7-7 7"></path>
                  </svg>
                </span>
              </div>
            </div>
          </article>
        </router-link>
      </div>
    </main>

    <!-- Global Footer Navigation -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import { useDate } from '../composables/useDate';

const { getLocalDay } = useDate();

const selectedCategory = ref<'all' | 'marian' | 'chaplets' | 'treasury'>('all');

const categories = [
  { id: 'all', label: 'All Devotions' },
  { id: 'marian', label: 'Marian Devotions' },
  { id: 'chaplets', label: 'Bead Chaplets' },
  { id: 'treasury', label: 'Prayer Treasury' }
] as const;

const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
const currentDayIndex = getLocalDay();
const currentDayName = dayNames[currentDayIndex] || 'Today';

const todayMystery = computed(() => {
  switch (currentDayIndex) {
    case 0: return 'Glorious';
    case 1: return 'Joyful';
    case 2: return 'Sorrowful';
    case 3: return 'Glorious';
    case 4: return 'Luminous';
    case 5: return 'Sorrowful';
    case 6: return 'Joyful';
    default: return 'Joyful';
  }
});

const dailyRecommendation = computed(() => {
  if (currentDayIndex === 5) {
    return {
      title: 'Friday Meditation: The Seven Sorrows & Sorrowful Mysteries',
      subtitle: 'Fridays are traditionally dedicated to the Passion of Our Lord and the Seven Sorrows of the Blessed Virgin Mary.',
      path: '/seven-sorrows',
      cta: 'Pray Seven Sorrows'
    };
  }
  return {
    title: `Pray Today’s ${todayMystery.value} Mysteries of the Holy Rosary`,
    subtitle: `Join the universal Church in meditating on the ${todayMystery.value} Mysteries with verse-by-verse Sacred Scripture.`,
    path: '/rosary',
    cta: 'Pray Rosary'
  };
});

interface DevotionTile {
  path: string;
  title: string;
  latinTitle: string;
  tag: string;
  badge?: string;
  category: ('marian' | 'chaplets' | 'treasury')[];
  description: string;
  features: string[];
  duration: string;
  image: string;
}

const devotions = computed<DevotionTile[]>(() => [
  {
    path: '/rosary',
    title: 'Scriptural Rosary',
    latinTitle: 'Rosarium Virginis Mariae',
    tag: 'Marian Rosary',
    badge: `Today: ${todayMystery.value}`,
    category: ['marian', 'chaplets'],
    description:
      'Meditate on the twenty mysteries of the life of Christ with Our Lady, featuring a Scripture verse before every Hail Mary, classical artwork, and audio guidance.',
    features: ['4 Mystery Sets', '200 Scripture Verses', 'Audio & Latin Interlinear', 'Custom Prayers'],
    duration: '15–20 mins • 5 Decades',
    image: '/images/rosary/annunciation.png'
  },
  {
    path: '/seven-sorrows',
    title: 'Seven Sorrows of Mary',
    latinTitle: 'Corolla Septem Dolorum',
    tag: 'Servite Chaplet',
    badge: '7 Mysteries & Verses',
    category: ['marian', 'chaplets'],
    description:
      'Contemplate the seven swords of sorrow that pierced the Immaculate Heart of Mary, with a full Bible passage displayed before starting each mystery and classical museum paintings.',
    features: ['7 Sorrow Mysteries', 'Scripture Before Each Mystery', '7 Classical Paintings', '1 Our Father + 7 Hail Marys'],
    duration: '15 mins • 7 Septets',
    image: '/images/seven-sorrows/sorrow_6_pieta.jpg'
  },
  {
    path: '/divine-mercy',
    title: 'Divine Mercy Chaplet',
    latinTitle: 'Corolla Divinae Misericordiae',
    tag: 'Mercy Chaplet',
    badge: '3:00 PM Hour of Mercy',
    category: ['chaplets'],
    description:
      'Revealed by Our Lord to St. Faustina Kowalska to implore God’s unfathomable mercy upon us and the whole world through the sorrowful Passion of Christ.',
    features: ['Full & Short Versions', 'Spoken Audio Guide', 'Interactive Beads', 'Latin & English'],
    duration: '7–10 mins • 5 Decades',
    image: '/images/divine_mercy.jpg'
  },
  {
    path: '/st-michael-rosary',
    title: 'St. Michael Chaplet',
    latinTitle: 'Corolla Sancti Michaelis Archangeli',
    tag: 'Angelic Chaplet',
    badge: '9 Angelic Choirs',
    category: ['chaplets'],
    description:
      'Invoke St. Michael the Archangel and the nine celestial Choirs of Angels for spiritual defense, purity of heart, and angelic escort at Holy Communion.',
    features: ['9 Angelic Salutations', 'Promises of St. Michael', 'Interactive Choir Beads', 'Latin & English'],
    duration: '10–12 mins • 9 Salutations',
    image: '/st_michael_archangel.jpg'
  },
  {
    path: '/prayers',
    title: 'Traditional Prayers Gallery',
    latinTitle: 'Thesaurus Orationum',
    tag: 'Prayer Anthology',
    badge: 'Audio & Latin',
    category: ['marian', 'treasury'],
    description:
      'Explore an illuminated treasury of essential Catholic prayers, Marian antiphons, creeds, acts of consecration, and litanies in English and Latin.',
    features: ['Searchable Treasury', 'Spoken Audio Recordings', 'Chant Videos', 'Orthodox Rubric Styling'],
    duration: 'Self-paced • 30+ Prayers',
    image: '/images/rosary/coronation.png'
  }
]);

const filteredDevotions = computed(() => {
  if (selectedCategory.value === 'all') return devotions.value;
  return devotions.value.filter(d => d.category.includes(selectedCategory.value as any));
});
</script>
