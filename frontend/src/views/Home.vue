<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-24 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content -->
    <main class="flex-grow max-w-2xl mx-auto w-full px-4 pt-12 pb-8 flex flex-col items-center">
      <!-- Word of the Day Hero -->
      <div class="text-center mb-10 w-full animate-fade-in-down">
        <span class="text-parchment-primary font-bold tracking-[0.25em] text-xs uppercase block mb-3">
          The Word of the Day
        </span>
        <h2 class="font-serif text-2xl md:text-3xl text-parchment-neutral italic mb-3 leading-snug">
          "Abide in My love."
        </h2>
        <p class="text-parchment-neutral/70 italic text-sm leading-relaxed max-w-md mx-auto">
          Reflect on the sacred rhythms of the liturgical year and find peace in the quiet moments of the day.
        </p>
      </div>

      <!-- Devotionals List -->
      <div class="w-full space-y-4 animate-fade-in-up delay-150">
        <!-- Daily Mass Readings -->
        <router-link to="/daily-readings" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1-2.5-2.5Z"></path>
                  <path d="M6 6h10"></path>
                  <path d="M6 10h10"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Daily Mass Readings</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  Follow the universal Church in its daily scriptural cycle. Today: Memorial of St. Polycarp.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Scriptural Rosary -->
        <router-link to="/rosary" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 22a7 7 0 0 0 7-7c0-4.3-7-11-7-11S5 10.7 5 15a7 7 0 0 0 7 7z"></path>
                  <circle cx="12" cy="15" r="2"></circle>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Scriptural Rosary</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  Meditate on the life of Christ through the mysteries of the most holy Rosary, accompanied by scripture.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Divine Mercy -->
        <router-link to="/divine-mercy" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Divine Mercy</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  Pray the Chaplet and enter into the depths of God's unfathomable mercy. Ideal for the 3:00 PM hour.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Catechism -->
        <router-link to="/catechism" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21.42 10.922a1 1 0 0 0-.019-1.838L12.83 5.18a2 2 0 0 0-1.66 0L2.6 9.08a1 1 0 0 0 0 1.832l8.57 3.908a2 2 0 0 0 1.66 0z"></path>
                  <path d="M6 12v5c0 2 2 3 6 3s6-1 6-3v-5"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Catechism</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  Explore the pillar of faith through a searchable and navigable treasury of Church teaching.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Traditional Prayers -->
        <router-link to="/prayers" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="m12 3-1.912 5.886L4.2 9.08l4.956 3.6L7.243 18.57 12 15l4.757 3.57-1.913-5.887 4.957-3.6-5.887-.193Z"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Traditional Prayers</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  An anthology of ancient prayers, novenas, and litanies passed down through generations.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Early Church Resources -->
        <router-link to="/resources" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                  <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Early Church Resources</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  Read the letters of Ignatius, the apologies of Justin Martyr, and other early writings.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Scripture Context Explorer -->
        <router-link to="/bible-study" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Bible Study</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  Retrieve original Bible texts alongside historical and theological Catholic commentary.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="m9 18 6-6-6-6"></path>
              </svg>
            </div>
          </ParchmentCard>
        </router-link>

        <!-- Magisterium AI Chat -->
        <router-link to="/magisterium-chat" class="block">
          <ParchmentCard hover-effect class="!p-5 flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-xl bg-parchment-neutral-light/80 flex items-center justify-center text-parchment-primary-dark shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                </svg>
              </div>
              <div class="text-left">
                <h3 class="font-serif text-base md:text-lg font-bold text-parchment-neutral">Magisterium AI Chat</h3>
                <p class="text-xs md:text-sm text-parchment-neutral/60 mt-0.5 leading-snug max-w-md">
                  Ask questions on Church doctrine and explore Catholic teachings with AI assistance.
                </p>
              </div>
            </div>
            <div class="text-parchment-neutral/30 group-hover:text-parchment-neutral/60 transition-colors pl-2">
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
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
</script>

<style scoped>
/* Page Animations */
.animate-fade-in-down {
  animation: fadeInDown 0.6s ease-out forwards;
}
.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out forwards;
}
.delay-150 {
  animation-delay: 0.15s;
  animation-fill-mode: both;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-15px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(15px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>