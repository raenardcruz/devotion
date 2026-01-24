<template>
  <div class="min-h-screen bg-stone-50 text-stone-900 selection:bg-amber-100 flex flex-col relative font-sans">
    <!-- Back Button -->
    <div class="absolute top-4 left-4 z-50">
      <router-link to="/" class="flex items-center space-x-2 text-white/80 hover:text-white bg-black/20 hover:bg-black/40 backdrop-blur-sm px-4 py-2 rounded-full transition-all text-sm font-medium">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="m15 18-6-6 6-6"/>
        </svg>
        <span>Home</span>
      </router-link>
    </div>

    <!-- Dynamic Background Blur -->
    <div class="fixed inset-0 z-0 opacity-20 transition-all duration-1000 ease-in-out" :style="{
            backgroundImage: `url('https://images.unsplash.com/photo-1569845177077-2a37322a60c7?q=80&w=2248&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            filter: 'blur(40px) scale(1.1)'
        }"></div>

    <div class="relative z-10 max-w-4xl mx-auto px-4 py-8 md:py-12 flex flex-col min-h-screen w-full">

        <RosaryHeader 
            :sets="sets" 
            :currentSetName="currentSetName" 
            :todayMystery="todayMystery" 
            @changeSet="changeSet" 
        />

        <!-- Main Content Area -->
        <main class="flex-grow flex flex-col justify-center gap-8">

            <transition name="fade" mode="out-in">
                <div :key="currentSetName + mysteryIndex + beadIndex">
                    <MysteryInfo 
                        :mysteryIndex="mysteryIndex" 
                        :currentMystery="currentMystery" 
                    />

                    <RosaryBeads :beadIndex="beadIndex" />

                    <VerseCard 
                        :currentVerse="currentVerse" 
                        :beadIndex="beadIndex" 
                    />
                </div>
            </transition>

            <RosaryControls 
                :mysteryIndex="mysteryIndex" 
                :beadIndex="beadIndex" 
                @next="nextBead" 
                @prev="prevBead" 
            />
        </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { ROSARY_DATA, type Mystery, type Verse } from '../components/rosary/rosaryData';
import RosaryHeader from '../components/rosary/RosaryHeader.vue';
import MysteryInfo from '../components/rosary/MysteryInfo.vue';
import RosaryBeads from '../components/rosary/RosaryBeads.vue';
import VerseCard from '../components/rosary/VerseCard.vue';
import RosaryControls from '../components/rosary/RosaryControls.vue';

const currentSetName = ref('Joyful');
const mysteryIndex = ref(0);
const beadIndex = ref(0);
const todayMystery = ref('');
const sets = Object.keys(ROSARY_DATA);

const getMysteryForDay = (): string => {
    const day = new Date().getDay();
    switch (day) {
        case 0: return 'Glorious'; // Sunday
        case 1: return 'Joyful';   // Monday
        case 2: return 'Sorrowful';// Tuesday
        case 3: return 'Glorious'; // Wednesday
        case 4: return 'Luminous'; // Thursday
        case 5: return 'Sorrowful';// Friday
        case 6: return 'Joyful';   // Saturday
        default: return 'Joyful';
    }
};

onMounted(() => {
    const daily = getMysteryForDay();
    todayMystery.value = daily;
    currentSetName.value = daily;
});

const currentMystery = computed((): Mystery => {
    const set = ROSARY_DATA[currentSetName.value];
    if (set && set[mysteryIndex.value]) {
        return set[mysteryIndex.value]!;
    }
    // Fallback or safe return
    return ROSARY_DATA['Joyful']![0]!;
});

const currentVerse = computed((): Verse => {
    if (currentMystery.value?.verses?.[beadIndex.value]) {
        return currentMystery.value.verses[beadIndex.value]!;
    }
    return { text: '', ref: '' };
});

const nextBead = () => {
    if (beadIndex.value < 9) {
        beadIndex.value++;
    } else if (mysteryIndex.value < 4) {
        mysteryIndex.value++;
        beadIndex.value = 0;
    }
};

const prevBead = () => {
    if (beadIndex.value > 0) {
        beadIndex.value--;
    } else if (mysteryIndex.value > 0) {
        mysteryIndex.value--;
        beadIndex.value = 9;
    }
};

const changeSet = (setName: string) => {
    currentSetName.value = setName;
    mysteryIndex.value = 0;
    beadIndex.value = 0;
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
    transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: scale(0.98);
}
</style>