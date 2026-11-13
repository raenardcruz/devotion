<template>
  <div class="min-h-screen bg-[#F7F3EE] text-[#322D29] flex flex-col pb-28 lg:pb-16 selection:bg-[#5C282D]/20 relative z-10">
    <!-- Subtle Warm Ambient Glow -->
    <div class="pointer-events-none fixed inset-0 overflow-hidden z-0" aria-hidden="true">
      <div class="absolute -top-32 left-1/2 -translate-x-1/2 w-[780px] h-[420px] rounded-full bg-[#E8DDD0]/55 blur-3xl"></div>
    </div>

    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-4xl mx-auto w-full px-4 sm:px-6 pt-8 md:pt-11 pb-10 flex flex-col relative z-10">
      <!-- Liturgical Word of the Day Hero -->
      <section class="text-center mb-8 md:mb-10 animate-fade-in-down">
        <div class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#E9DFD3] text-[#5C282D] text-[10px] sm:text-[11px] font-bold uppercase tracking-[0.18em] mb-3.5">
          <span>{{ currentDayName }}</span>
          <span class="opacity-60">&bull;</span>
          <span>{{ todayMystery }} Mysteries</span>
        </div>
        <h1 class="font-serif text-3xl sm:text-4xl md:text-[2.65rem] text-[#4A2525] mb-2.5 leading-tight font-semibold tracking-tight">
          &ldquo;Abide in My love.&rdquo;
        </h1>
        <p class="text-[#322D29]/75 text-xs sm:text-sm md:text-[15px] leading-relaxed max-w-md mx-auto">
          Reflect on the sacred rhythm of the liturgical year and find peace in the quiet moments of the day.
        </p>
      </section>

      <!-- Last Read Featured Hero Banner (Icon breaks outside card edge) -->
      <section class="w-full mb-9 md:mb-11 animate-fade-in-up">
        <router-link :to="lastRead.path" class="group block relative">
          <div class="relative rounded-[1.75rem] bg-gradient-to-br from-[#5B252A] via-[#4D1E23] to-[#351317] text-white shadow-[0_18px_38px_-12px_rgba(62,24,28,0.45)] transition-all duration-300 group-hover:shadow-[0_24px_44px_-12px_rgba(62,24,28,0.55)] group-hover:-translate-y-0.5">
            <!-- Clipped Organic Wave Backdrop -->
            <div class="absolute inset-0 rounded-[1.75rem] overflow-hidden pointer-events-none" aria-hidden="true">
              <svg class="absolute inset-0 w-full h-full opacity-20" viewBox="0 0 800 220" preserveAspectRatio="none">
                <path d="M0 165 C190 110, 420 210, 800 105 L800 220 L0 220 Z" fill="#8C4047" />
                <path d="M0 190 C250 150, 520 225, 800 145 L800 220 L0 220 Z" fill="#2A0D10" fill-opacity="0.45" />
              </svg>
            </div>

            <!-- Content + Overflowing Open Book SVG -->
            <div class="relative z-10 flex items-center justify-between gap-3 px-5 py-5 sm:px-8 sm:py-6">
              <div class="flex flex-col items-start text-left max-w-[62%] sm:max-w-[65%]">
                <div class="inline-flex items-center gap-1.5 text-[11px] sm:text-xs font-medium text-[#F3E7D9]/85 mb-1.5">
                  <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="opacity-85">
                    <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                    <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
                  </svg>
                  <span>Last Read</span>
                </div>

                <h2 class="font-serif text-lg sm:text-2xl font-semibold text-white leading-snug tracking-tight mb-4">
                  {{ lastRead.title }}
                </h2>

                <span class="inline-flex items-center gap-2 px-4 sm:px-5 py-2 rounded-full bg-[#F6EFE6] group-hover:bg-white text-[#4E2226] text-xs font-bold tracking-wide shadow-sm transition-all duration-200 group-hover:gap-2.5">
                  <span>Continue</span>
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M5 12h14"></path>
                    <path d="m12 5 7 7-7 7"></path>
                  </svg>
                </span>
              </div>

              <!-- Open Lectionary Book breaking outside right & bottom-right edge -->
              <div class="relative shrink-0 w-36 sm:w-48 md:w-56 -mr-4 sm:-mr-8 -my-5 sm:-my-7 transition-transform duration-500 ease-out group-hover:scale-105 group-hover:-rotate-2">
                <OpenLectionaryHeroIcon class="w-full h-auto" />
              </div>
            </div>
          </div>
        </router-link>
      </section>

      <!-- Popular Section -->
      <section class="w-full mb-7 animate-fade-in-up delay-150">
        <div class="flex items-center justify-between mb-4 sm:mb-5 px-0.5">
          <h2 class="text-xl sm:text-2xl font-semibold text-[#322D29] tracking-tight">
            Popular
          </h2>
        </div>

        <!-- 2x2 Bento Grid (Icons break outside card boundaries for dimensional depth) -->
        <div class="grid grid-cols-2 gap-4 sm:gap-6">
          <!-- 1. Devotions (Deep Burgundy) -->
          <router-link to="/devotions" class="group block relative">
            <article class="relative h-full min-h-[185px] sm:min-h-[195px] rounded-[1.5rem] sm:rounded-[1.75rem] bg-[#5C282D] text-white p-4 sm:p-6 shadow-[0_14px_30px_-12px_rgba(92,40,45,0.45)] transition-all duration-300 group-hover:-translate-y-1 group-hover:shadow-[0_20px_36px_-12px_rgba(92,40,45,0.55)] flex flex-col justify-between">
              <!-- Subtle inner curve accent -->
              <div class="absolute inset-0 rounded-[1.5rem] sm:rounded-[1.75rem] overflow-hidden pointer-events-none" aria-hidden="true">
                <div class="absolute -bottom-10 -right-10 w-40 h-40 rounded-full bg-white/5"></div>
              </div>

              <h3 class="relative z-10 text-base sm:text-xl font-semibold tracking-tight text-white">
                Devotions
              </h3>

              <!-- Mobile: centered overflowing icon | Desktop: left icon + right copy & ENTRY pill -->
              <div class="relative z-10 flex flex-col sm:flex-row sm:items-end justify-between gap-3 mt-2 sm:mt-3 flex-1">
                <div class="w-28 sm:w-32 mx-auto sm:mx-0 -mb-6 sm:-mb-7 sm:-ml-5 transition-transform duration-500 ease-out group-hover:scale-108 group-hover:-rotate-3 shrink-0">
                  <RosaryIllustration class="w-full h-auto" />
                </div>

                <div class="hidden sm:flex flex-col items-start justify-end flex-1 pb-0.5">
                  <p class="text-xs text-[#F5EBE0]/85 leading-relaxed mb-3.5 line-clamp-3">
                    Pray the Scriptural Rosary, Seven Sorrows of Mary, Divine Mercy, and sacred Catholic chaplets.
                  </p>
                  <span class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-full bg-[#F6EFE6] group-hover:bg-white text-[#4E2226] text-[10px] font-bold uppercase tracking-wider shadow-xs transition-all">
                    <span>Entry</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M5 12h14"></path>
                      <path d="m12 5 7 7-7 7"></path>
                    </svg>
                  </span>
                </div>
              </div>
            </article>
          </router-link>

          <!-- 2. Daily Mass Readings (Warm Sand / Parchment) -->
          <router-link to="/daily-readings" class="group block relative">
            <article class="relative h-full min-h-[185px] sm:min-h-[195px] rounded-[1.5rem] sm:rounded-[1.75rem] bg-[#DDD2C3] text-[#2E2421] p-4 sm:p-6 shadow-[0_14px_30px_-12px_rgba(95,78,60,0.22)] transition-all duration-300 group-hover:-translate-y-1 group-hover:shadow-[0_20px_36px_-12px_rgba(95,78,60,0.3)] flex flex-col justify-between">
              <div class="absolute inset-0 rounded-[1.5rem] sm:rounded-[1.75rem] overflow-hidden pointer-events-none" aria-hidden="true">
                <div class="absolute -top-10 -right-10 w-36 h-36 rounded-full bg-white/25"></div>
              </div>

              <h3 class="relative z-10 text-base sm:text-xl font-semibold tracking-tight text-[#2E2421] leading-snug">
                Daily Mass Readings
              </h3>

              <div class="relative z-10 flex flex-col sm:flex-row sm:items-end justify-between gap-3 mt-2 sm:mt-3 flex-1">
                <div class="w-28 sm:w-32 mx-auto sm:mx-0 -mb-5 sm:-mb-6 sm:-ml-4 transition-transform duration-500 ease-out group-hover:scale-108 group-hover:rotate-2 shrink-0">
                  <MassReadingsIllustration class="w-full h-auto" />
                </div>

                <div class="hidden sm:flex flex-col items-start justify-end flex-1 pb-0.5">
                  <p class="text-xs text-[#2E2421]/80 leading-relaxed mb-3.5 line-clamp-3">
                    Follow the universal Church in its daily scriptural cycle and Traditional Latin Mass propers.
                  </p>
                  <span class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-full bg-[#FAF6F0] group-hover:bg-white text-[#4E2226] text-[10px] font-bold uppercase tracking-wider shadow-xs transition-all">
                    <span>Entry</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M5 12h14"></path>
                      <path d="m12 5 7 7-7 7"></path>
                    </svg>
                  </span>
                </div>
              </div>
            </article>
          </router-link>

          <!-- 3. Bible Study (Terracotta / Sienna) -->
          <router-link to="/bible-study" class="group block relative mt-1 sm:mt-0">
            <article class="relative h-full min-h-[185px] sm:min-h-[195px] rounded-[1.5rem] sm:rounded-[1.75rem] bg-[#A35733] text-white p-4 sm:p-6 shadow-[0_14px_30px_-12px_rgba(163,87,51,0.42)] transition-all duration-300 group-hover:-translate-y-1 group-hover:shadow-[0_20px_36px_-12px_rgba(163,87,51,0.52)] flex flex-col justify-between">
              <div class="absolute inset-0 rounded-[1.5rem] sm:rounded-[1.75rem] overflow-hidden pointer-events-none" aria-hidden="true">
                <div class="absolute -bottom-10 -right-10 w-40 h-40 rounded-full bg-black/10"></div>
              </div>

              <h3 class="relative z-10 text-base sm:text-xl font-semibold tracking-tight text-white">
                Bible Study
              </h3>

              <div class="relative z-10 flex flex-col sm:flex-row sm:items-end justify-between gap-3 mt-2 sm:mt-3 flex-1">
                <div class="w-26 sm:w-32 mx-auto sm:mx-0 -mb-6 sm:-mb-7 sm:-ml-4 transition-transform duration-500 ease-out group-hover:scale-108 group-hover:-rotate-3 shrink-0">
                  <BibleStudyIllustration class="w-full h-auto" />
                </div>

                <div class="hidden sm:flex flex-col items-start justify-end flex-1 pb-0.5">
                  <p class="text-xs text-[#FFF4EC]/85 leading-relaxed mb-3.5 line-clamp-3">
                    Retrieve original Sacred Scripture texts alongside historical and theological Catholic commentary.
                  </p>
                  <span class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-full bg-[#F6EFE6] group-hover:bg-white text-[#4E2226] text-[10px] font-bold uppercase tracking-wider shadow-xs transition-all">
                    <span>Entry</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M5 12h14"></path>
                      <path d="m12 5 7 7-7 7"></path>
                    </svg>
                  </span>
                </div>
              </div>
            </article>
          </router-link>

          <!-- 4. Catechism (Warm Ochre / Gold) -->
          <router-link to="/catechism" class="group block relative mt-1 sm:mt-0">
            <article class="relative h-full min-h-[185px] sm:min-h-[195px] rounded-[1.5rem] sm:rounded-[1.75rem] bg-[#CCA45E] text-[#2E2421] p-4 sm:p-6 shadow-[0_14px_30px_-12px_rgba(160,122,48,0.35)] transition-all duration-300 group-hover:-translate-y-1 group-hover:shadow-[0_20px_36px_-12px_rgba(160,122,48,0.45)] flex flex-col justify-between">
              <div class="absolute inset-0 rounded-[1.5rem] sm:rounded-[1.75rem] overflow-hidden pointer-events-none" aria-hidden="true">
                <div class="absolute -bottom-10 -right-10 w-40 h-40 rounded-full bg-white/15"></div>
              </div>

              <h3 class="relative z-10 text-base sm:text-xl font-semibold tracking-tight text-[#2E2421]">
                Catechism
              </h3>

              <div class="relative z-10 flex flex-col sm:flex-row sm:items-end justify-between gap-3 mt-2 sm:mt-3 flex-1">
                <div class="w-26 sm:w-32 mx-auto sm:mx-0 -mb-6 sm:-mb-7 sm:-ml-4 transition-transform duration-500 ease-out group-hover:scale-108 group-hover:rotate-3 shrink-0">
                  <CatechismIllustration class="w-full h-auto" />
                </div>

                <div class="hidden sm:flex flex-col items-start justify-end flex-1 pb-0.5">
                  <p class="text-xs text-[#2E2421]/85 leading-relaxed mb-3.5 line-clamp-3">
                    Explore the four pillars of the Catholic faith through a searchable treasury of Church teaching.
                  </p>
                  <span class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-full bg-[#FAF6F0] group-hover:bg-white text-[#4E2226] text-[10px] font-bold uppercase tracking-wider shadow-xs transition-all">
                    <span>Entry</span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M5 12h14"></path>
                      <path d="m12 5 7 7-7 7"></path>
                    </svg>
                  </span>
                </div>
              </div>
            </article>
          </router-link>
        </div>
      </section>

      <!-- Early Church Resources Card -->
      <section class="w-full mt-1 animate-fade-in-up delay-150">
        <router-link to="/resources" class="group block relative">
          <article class="relative h-full min-h-[165px] sm:min-h-[180px] rounded-[1.5rem] sm:rounded-[1.75rem] bg-[#4E5B47] text-white p-4 sm:p-6 shadow-[0_14px_30px_-12px_rgba(78,91,71,0.42)] transition-all duration-300 group-hover:-translate-y-1 group-hover:shadow-[0_20px_36px_-12px_rgba(78,91,71,0.52)] flex flex-col justify-between">
            <!-- Subtle inner curve accent -->
            <div class="absolute inset-0 rounded-[1.5rem] sm:rounded-[1.75rem] overflow-hidden pointer-events-none" aria-hidden="true">
              <div class="absolute -bottom-10 -right-10 w-44 h-44 rounded-full bg-white/10"></div>
            </div>

            <h3 class="relative z-10 text-base sm:text-xl font-semibold tracking-tight text-white">
              Early Church Resources
            </h3>

            <div class="relative z-10 flex flex-row items-end justify-between gap-4 mt-2 sm:mt-3 flex-1">
              <div class="w-28 sm:w-34 -mb-6 sm:-mb-7 -ml-3 sm:-ml-4 transition-transform duration-500 ease-out group-hover:scale-108 group-hover:-rotate-2 shrink-0">
                <EarlyChurchIllustration class="w-full h-auto" />
              </div>

              <div class="flex flex-col items-start justify-end flex-1 pb-0.5">
                <p class="text-xs sm:text-[13px] text-[#F5EBE0]/85 leading-relaxed mb-3.5 line-clamp-3">
                  Read the letters of Ignatius of Antioch, Justin Martyr, and writings of the Church Fathers.
                </p>
                <span class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-full bg-[#F6EFE6] group-hover:bg-white text-[#4E2226] text-[10px] font-bold uppercase tracking-wider shadow-xs transition-all">
                  <span>Entry</span>
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M5 12h14"></path>
                    <path d="m12 5 7 7-7 7"></path>
                  </svg>
                </span>
              </div>
            </div>
          </article>
        </router-link>
      </section>
    </main>

    <!-- Global Footer Navigation (Mobile & Tablet) -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import OpenLectionaryHeroIcon from '../components/icons/OpenLectionaryHeroIcon.vue';
import RosaryIllustration from '../components/icons/RosaryIllustration.vue';
import MassReadingsIllustration from '../components/icons/MassReadingsIllustration.vue';
import BibleStudyIllustration from '../components/icons/BibleStudyIllustration.vue';
import CatechismIllustration from '../components/icons/CatechismIllustration.vue';
import EarlyChurchIllustration from '../components/icons/EarlyChurchIllustration.vue';
import { useDate } from '../composables/useDate';
import { useLastRead } from '../composables/useLastRead';

const { getLocalDay } = useDate();
const { storedLastRead } = useLastRead();
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

const lastRead = computed(() => storedLastRead.value ?? {
  title: `Daily Mass Readings - The ${todayMystery.value} Mysteries`,
  path: '/daily-readings',
});
</script>
