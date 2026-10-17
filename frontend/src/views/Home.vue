<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-3xl mx-auto w-full px-4 pt-10 md:pt-12 pb-10 flex flex-col items-center relative z-10">
      <!-- Word of the Day Hero -->
      <div class="text-center mb-10 w-full animate-fade-in-down">
        <div class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#72383D]/10 border border-[#72383D]/20 text-[#72383D] text-[11px] font-bold uppercase tracking-[0.2em] mb-3">
          <span>{{ currentDayName }}</span>
          <span>&bull;</span>
          <span>{{ todayMystery }} Mysteries</span>
        </div>
        <h1 class="font-serif text-3xl md:text-4xl bg-gradient-to-r from-[#72383D] via-[#AC9C8D] to-[#322D29] bg-clip-text text-transparent mb-3 leading-snug font-bold">
          "Abide in My love."
        </h1>
        <p class="text-[#322D29]/80 text-sm md:text-base leading-relaxed max-w-lg mx-auto font-serif">
          Reflect on the sacred rhythms of the liturgical year and find peace in the quiet moments of the day.
        </p>
      </div>

      <!-- Featured Devotions & Chaplets Hub Card -->
      <div class="w-full mb-6 animate-fade-in-up">
        <div class="bg-gradient-to-br from-white/95 via-[#EFE9E1]/85 to-[#D9D9D9]/40 border border-[#D1C7BD] rounded-[2rem] p-6 md:p-7 shadow-md backdrop-blur-md">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-5 border-b border-[#D1C7BD]/60">
            <div class="flex items-start space-x-4">
              <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-[#72383D] to-[#322D29] flex items-center justify-center text-white shadow-sm shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 22a7 7 0 0 0 7-7c0-4.3-7-11-7-11S5 10.7 5 15a7 7 0 0 0 7 7z"></path>
                  <circle cx="12" cy="15" r="2"></circle>
                </svg>
              </div>
              <div class="text-left">
                <span class="text-[10px] font-bold uppercase tracking-[0.2em] text-[#72383D] block mb-0.5">
                  Sacred Prayer Hub
                </span>
                <h2 class="font-serif text-xl md:text-2xl font-bold text-[#322D29]">
                  Devotions &amp; Chaplets
                </h2>
                <p class="text-xs md:text-sm text-[#322D29]/80 mt-0.5 leading-snug max-w-md">
                  Pray the Scriptural Rosary, Seven Sorrows of Mary, Divine Mercy, St. Michael Chaplet, and traditional prayers.
                </p>
              </div>
            </div>

            <router-link
              to="/devotions"
              class="px-5 py-2.5 rounded-full bg-gradient-to-r from-[#72383D] to-[#322D29] hover:from-[#8B464C] hover:to-[#453E38] text-white text-xs font-bold uppercase tracking-wider transition-all shrink-0 inline-flex items-center justify-center gap-1.5 shadow-xs min-h-[44px]"
            >
              <span>Explore All</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </router-link>
          </div>

          <!-- Quick Launch Devotion Tiles Row -->
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2.5 pt-4">
            <router-link
              to="/rosary"
              class="p-3 rounded-2xl bg-white/75 hover:bg-white border border-[#D1C7BD]/70 hover:border-[#72383D]/50 transition-all group flex flex-col justify-between"
            >
              <span class="text-[10px] font-bold uppercase tracking-wider text-[#72383D]/75">Rosary</span>
              <span class="font-serif text-sm font-bold text-[#322D29] group-hover:text-[#72383D] mt-1 flex items-center justify-between">
                <span>Holy Rosary</span>
                <span>→</span>
              </span>
            </router-link>

            <router-link
              to="/seven-sorrows"
              class="p-3 rounded-2xl bg-white/75 hover:bg-white border border-[#D1C7BD]/70 hover:border-[#72383D]/50 transition-all group flex flex-col justify-between"
            >
              <span class="text-[10px] font-bold uppercase tracking-wider text-[#72383D]/75">7 Mysteries</span>
              <span class="font-serif text-sm font-bold text-[#322D29] group-hover:text-[#72383D] mt-1 flex items-center justify-between">
                <span>Seven Sorrows</span>
                <span>→</span>
              </span>
            </router-link>

            <router-link
              to="/divine-mercy"
              class="p-3 rounded-2xl bg-white/75 hover:bg-white border border-[#D1C7BD]/70 hover:border-[#72383D]/50 transition-all group flex flex-col justify-between"
            >
              <span class="text-[10px] font-bold uppercase tracking-wider text-[#72383D]/75">3:00 PM Hour</span>
              <span class="font-serif text-sm font-bold text-[#322D29] group-hover:text-[#72383D] mt-1 flex items-center justify-between">
                <span>Divine Mercy</span>
                <span>→</span>
              </span>
            </router-link>

            <router-link
              to="/st-michael-rosary"
              class="p-3 rounded-2xl bg-white/75 hover:bg-white border border-[#D1C7BD]/70 hover:border-[#72383D]/50 transition-all group flex flex-col justify-between"
            >
              <span class="text-[10px] font-bold uppercase tracking-wider text-[#72383D]/75">9 Choirs</span>
              <span class="font-serif text-sm font-bold text-[#322D29] group-hover:text-[#72383D] mt-1 flex items-center justify-between">
                <span>St. Michael</span>
                <span>→</span>
              </span>
            </router-link>
          </div>
        </div>
      </div>

      <!-- Liturgy, Study & Doctrine List -->
      <div class="w-full space-y-4 animate-fade-in-up delay-150">
        <!-- Daily Mass Readings -->
        <router-link to="/daily-readings" class="block group">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-[#72383D]/15 to-[#AC9C8D]/20 flex items-center justify-center text-[#72383D] shadow-xs border border-[#D1C7BD]/50 shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
                  <path d="M6 6h10"></path>
                  <path d="M6 10h10"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-[#72383D]">Daily Mass Readings</h3>
                <p class="text-xs md:text-sm text-[#322D29]/80 mt-0.5 leading-snug max-w-md">
                  Follow the universal Church in its daily scriptural cycle. Reflect on today's holy scriptures and Traditional Latin Mass propers.
                </p>
              </div>
            </div>
            <div class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Bible Study -->
        <router-link to="/bible-study" class="block group">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-[#72383D]/15 to-[#AC9C8D]/20 flex items-center justify-center text-[#72383D] shadow-xs border border-[#D1C7BD]/50 shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-[#72383D]">Bible Study</h3>
                <p class="text-xs md:text-sm text-[#322D29]/80 mt-0.5 leading-snug max-w-md">
                  Retrieve original Bible texts alongside historical and theological Catholic commentary.
                </p>
              </div>
            </div>
            <div class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Catechism -->
        <router-link to="/catechism" class="block group">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-[#72383D]/15 to-[#AC9C8D]/20 flex items-center justify-center text-[#72383D] shadow-xs border border-[#D1C7BD]/50 shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21.42 10.922a1 1 0 0 0-.019-1.838L12.83 5.18a2 2 0 0 0-1.66 0L2.6 9.08a1 1 0 0 0 0 1.832l8.57 3.908a2 2 0 0 0 1.66 0z"></path>
                  <path d="M6 12v5c0 2 2 3 6 3s6-1 6-3v-5"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-[#72383D]">Catechism</h3>
                <p class="text-xs md:text-sm text-[#322D29]/80 mt-0.5 leading-snug max-w-md">
                  Explore the four pillars of the Catholic faith through a searchable and navigable treasury of Church teaching.
                </p>
              </div>
            </div>
            <div class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Early Church Resources -->
        <router-link to="/resources" class="block group">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-[#72383D]/15 to-[#AC9C8D]/20 flex items-center justify-center text-[#72383D] shadow-xs border border-[#D1C7BD]/50 shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                  <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-[#72383D]">Early Church Resources</h3>
                <p class="text-xs md:text-sm text-[#322D29]/80 mt-0.5 leading-snug max-w-md">
                  Read the letters of Ignatius of Antioch, the apologies of Justin Martyr, and writings of the Church Fathers.
                </p>
              </div>
            </div>
            <div class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Magisterium AI Chat -->
        <router-link to="/magisterium-chat" class="block group">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-[#72383D]/15 to-[#AC9C8D]/20 flex items-center justify-center text-[#72383D] shadow-xs border border-[#D1C7BD]/50 shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#72383D" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-[#72383D]">Magisterium AI Chat</h3>
                <p class="text-xs md:text-sm text-[#322D29]/80 mt-0.5 leading-snug max-w-md">
                  Ask questions on Church doctrine and explore Catholic teachings with Magisterium-grounded citations.
                </p>
              </div>
            </div>
            <div class="text-[#72383D]/60 group-hover:text-[#72383D] transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>
      </div>
    </main>

    <!-- Global Footer Navigation -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import { useDate } from '../composables/useDate';

const { getLocalDay } = useDate();
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
</script>