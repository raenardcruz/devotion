<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-28 sm:pb-20 selection:bg-parchment-primary/20">
    <!-- Global Header -->
    <TopNav />

    <!-- Content Area -->
    <div class="relative z-10 max-w-6xl mx-auto px-4 py-12 w-full flex-grow">
      <header class="mb-10 text-center animate-fade-in-down">
        <h1 class="text-3xl md:text-5xl font-serif text-parchment-primary-dark mb-2">Prayers Gallery</h1>
        <p class="text-parchment-neutral/60 text-sm md:text-base max-w-2xl mx-auto leading-relaxed">A collection of traditional Catholic prayers for reflection and devotion.</p>
      </header>
      
      <!-- Search Bar -->
      <div class="mb-8 max-w-md mx-auto animate-fade-in-down delay-75">
        <div class="relative">
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="Search prayers by name or content..." 
            class="w-full px-5 py-3 pl-12 pr-10 bg-parchment-neutral-light/50 border border-parchment-border rounded-full text-parchment-neutral placeholder-parchment-neutral/40 focus:outline-none focus:ring-2 focus:ring-parchment-primary/50 focus:border-parchment-primary transition-all font-sans text-sm shadow-sm"
          />
          <!-- Search Icon (Magnifying Glass SVG) -->
          <div class="absolute left-4 top-1/2 -translate-y-1/2 text-parchment-neutral/40 pointer-events-none">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="8"></circle>
              <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
            </svg>
          </div>
          <!-- Clear Button (X Icon) -->
          <div 
            v-if="searchQuery" 
            @click="searchQuery = ''" 
            role="button"
            tabindex="0"
            @keyup.enter="searchQuery = ''"
            class="absolute right-4 top-1/2 -translate-y-1/2 text-parchment-neutral/40 hover:text-parchment-neutral p-1 rounded-full cursor-pointer hover:bg-parchment-neutral/5 transition-all flex items-center justify-center select-none"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" x2="6" y1="6" y2="18"></line>
              <line x1="6" x2="18" y1="6" y2="18"></line>
            </svg>
          </div>
        </div>
      </div>
      
      <!-- Prayers Grid -->
      <div v-if="filteredPrayers.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 animate-fade-in-up delay-150">
        <div 
          v-for="prayer in filteredPrayers" 
          :key="prayer.id"
          @click="openPrayer(prayer)"
          class="group"
        >
          <!-- Using Reusable Parchment Card with hover animation -->
          <ParchmentCard hover-effect class="h-full flex flex-col justify-between">
            <div>
              <h3 class="text-xl md:text-2xl font-serif text-parchment-neutral font-medium mb-3 group-hover:text-parchment-primary-dark transition-colors">
                {{ prayer.name }}
                <span v-if="prayer.latinName" class="block text-xs font-sans text-parchment-neutral/40font-normal mt-1">{{ prayer.latinName }}</span>
              </h3>
              <p class="text-parchment-neutral/70 text-sm line-clamp-3 leading-relaxed">{{ prayer.default }}</p>
            </div>
            
            <div class="mt-5 flex items-center justify-between">
              <div class="flex items-center text-parchment-primary font-bold text-xs uppercase tracking-wider group-hover:text-parchment-primary-dark transition-all select-none">
                <span>{{ (prayer.default || prayer.latin) ? 'Read Prayer' : 'Watch Video' }}</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="ml-2 group-hover:translate-x-1.5 transition-transform">
                  <path d="M5 12h14"></path>
                  <path d="m12 5 7 7-7 7"></path>
                </svg>
              </div>
              
              <div class="flex items-center gap-1.5">
                <!-- Audio Indicator Badge -->
                <span v-if="hasAudio(prayer.id)" class="text-[10px] bg-parchment-primary/10 text-parchment-primary-dark px-2.5 py-0.5 rounded-full font-bold uppercase tracking-wider flex items-center gap-1 select-none">
                  <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon>
                    <path d="M15.54 8.46a5 5 0 0 1 0 7.07"></path>
                  </svg>
                  Audio
                </span>

                <!-- Video Indicator Badge -->
                <span v-if="prayer.youtube" class="text-[10px] bg-parchment-secondary/10 text-parchment-secondary px-2.5 py-0.5 rounded-full font-bold uppercase tracking-wider flex items-center gap-1 select-none">
                  <svg xmlns="http://www.w3.org/2000/svg" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M23 7a2 2 0 0 0-2.45-1.45L16 7V5a2 2 0 0 0-2-2H2a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2l4.55 1.45A2 2 0 0 0 23 17V7Z"></path>
                  </svg>
                  Video
                </span>
              </div>
            </div>
          </ParchmentCard>
        </div>
      </div>

      <!-- Empty State -->
      <div 
        v-else 
        class="text-center py-16 px-4 animate-fade-in-up flex flex-col items-center justify-center space-y-4"
      >
        <div class="w-16 h-16 rounded-full bg-parchment-primary/10 flex items-center justify-center text-parchment-primary mb-2 shadow-sm ring-4 ring-parchment-primary/5 animate-pulse">
          <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
            <line x1="8" x2="14" y1="11" y2="11"></line>
          </svg>
        </div>
        <h3 class="font-serif text-2xl text-parchment-primary-dark font-medium">No prayers found</h3>
        <p class="text-parchment-neutral/60 text-sm max-w-sm mx-auto leading-relaxed font-sans">
          We couldn't find any prayers matching "{{ searchQuery }}". Try searching for another keyword or language.
        </p>
        <button 
          @click="searchQuery = ''" 
          class="btn-primary mt-2 font-sans font-bold text-xs uppercase tracking-wider px-5 py-2.5 hover:scale-[1.02] active:scale-95 transition-all outline-none shadow-sm cursor-pointer"
        >
          Clear Search
        </button>
      </div>
    </div>

    <!-- Prayer Modal -->
    <transition name="modal">
      <div v-if="selectedPrayer" class="fixed inset-0 z-[100] flex items-center justify-center p-4" @click.self="closePrayer">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-parchment-neutral/30 backdrop-blur-sm transition-opacity" @click="closePrayer"></div>
        
        <!-- Modal Card Container -->
        <div class="relative bg-parchment-bg border border-parchment-border rounded-[2.5rem] p-6 md:p-10 max-w-2xl w-full max-h-[85vh] overflow-y-auto shadow-2xl transform transition-all">
          <button @click="closePrayer" class="absolute top-5 right-5 text-parchment-neutral/40 hover:text-parchment-neutral transition-colors p-2 rounded-full hover:bg-parchment-neutral-light border-none shadow-none">
            <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" x2="6" y1="6" y2="18"></line>
              <line x1="6" x2="18" y1="6" y2="18"></line>
            </svg>
          </button>
          
          <div class="flex flex-col items-center mb-8 mt-4">
            <h2 class="text-2xl md:text-3xl font-serif text-parchment-primary-dark mb-4 text-center font-medium">
              {{ showLatin && selectedPrayer.latinName ? selectedPrayer.latinName : selectedPrayer.name }}
            </h2>
            
            <div class="flex flex-col sm:flex-row items-center gap-3">
              <!-- Latin/English Toggle (Only if Latin exists and we are in audio mode) -->
              <div v-if="selectedPrayer.latin && selectedPlayMode === 'audio'" class="flex items-center bg-parchment-neutral-light border border-parchment-border p-0.5 rounded-full w-fit">
                <button 
                  @click="showLatin = false"
                  :class="[!showLatin ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                  class="px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                >
                  English
                </button>
                <button 
                  @click="showLatin = true"
                  :class="[showLatin ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                  class="px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                >
                  Latin
                </button>
              </div>

              <!-- Play Audio Button (Only in audio mode) -->
              <button 
                v-if="hasAudio(selectedPrayer.id) && selectedPlayMode === 'audio'"
                @click="toggleAudio(selectedPrayer.id)"
                class="flex items-center gap-2 px-4 py-1.5 bg-parchment-primary/10 border border-parchment-primary/30 text-parchment-primary-dark rounded-full hover:bg-parchment-primary/20 active:scale-95 transition-all text-xs font-bold uppercase tracking-wider outline-none shadow-none"
              >
                <svg v-if="!isPlaying" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2">
                  <polygon points="6 3 20 12 6 21 6 3"></polygon>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" x2="18" y1="4" y2="20"></line>
                  <line x1="6" x2="6" y1="4" y2="20"></line>
                </svg>
                <span>{{ isPlaying ? 'Pause Audio' : 'Play Audio' }}</span>
              </button>

              <!-- Play Mode Toggle (Only if video is available and prayer has text) -->
              <div v-if="selectedPrayer.youtube && (selectedPrayer.default || selectedPrayer.latin)" class="flex items-center bg-parchment-neutral-light border border-parchment-border p-0.5 rounded-full w-fit">
                <button 
                  @click="selectedPlayMode = 'audio'"
                  :class="[selectedPlayMode === 'audio' ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                  class="px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                >
                  Audio &amp; Text
                </button>
                <button 
                  @click="selectedPlayMode = 'video'; stopAudio()"
                  :class="[selectedPlayMode === 'video' ? 'bg-parchment-primary text-white font-bold' : 'text-parchment-neutral/50 hover:text-parchment-neutral']"
                  class="px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all duration-300 border-none outline-none shadow-none hover:bg-transparent hover:translate-y-0"
                >
                  Watch Video
                </button>
              </div>
            </div>
          </div>

          <!-- Content Block -->
          <transition name="fade-content" mode="out-in">
            <div :key="selectedPlayMode + (showLatin ? 'latin' : 'english')" class="max-w-xl mx-auto text-parchment-neutral leading-loose">
              <!-- Video Player -->
              <div v-if="selectedPlayMode === 'video'" class="relative w-full aspect-video rounded-2xl overflow-hidden border border-parchment-border shadow-lg my-4">
                <iframe
                  class="absolute inset-0 w-full h-full"
                  :src="youtubeEmbedUrl"
                  title="YouTube video player"
                  frameborder="0"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                  allowfullscreen
                ></iframe>
              </div>
              
              <!-- Text Content -->
              <p 
                v-else
                :class="[showLatin ? 'font-seriftext-parchment-primary-dark text-lg' : 'font-serif text-base md:text-lg text-parchment-neutral/90']" 
                class="whitespace-pre-line text-center"
              >
                {{ showLatin ? selectedPrayer.latin : selectedPrayer.default }}
              </p>
            </div>
          </transition>
        </div>
      </div>
    </transition>

    <!-- Global Footer -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount } from 'vue';
import prayersData from '../data/prayers.json';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import { getPrayerAudioUrl, AVAILABLE_AUDIO_PRAYERS } from '../utils/audioHelper';
import { getYouTubeEmbedUrl } from '../utils/youtube';

interface Prayer {
  id: string;
  name: string;
  latinName?: string;
  default: string;
  latin?: string;
  youtube?: string;
}

const prayers = ref<Prayer[]>(prayersData);
const searchQuery = ref('');

const filteredPrayers = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return prayers.value;
  
  return prayers.value.filter(prayer => {
    return (
      prayer.name.toLowerCase().includes(query) ||
      (prayer.latinName && prayer.latinName.toLowerCase().includes(query)) ||
      (prayer.default && prayer.default.toLowerCase().includes(query)) ||
      (prayer.latin && prayer.latin.toLowerCase().includes(query))
    );
  });
});
const selectedPrayer = ref<Prayer | null>(null);
const showLatin = ref(false);
const selectedPlayMode = ref<'audio' | 'video'>('audio');

const youtubeEmbedUrl = computed(() => {
  return selectedPrayer.value ? getYouTubeEmbedUrl(selectedPrayer.value.youtube) : '';
});

// Audio Player State
const currentAudio = ref<HTMLAudioElement | null>(null);
const isPlaying = ref(false);

const hasAudio = (prayerId: string) => {
  return AVAILABLE_AUDIO_PRAYERS.has(prayerId);
};

const toggleAudio = (prayerId: string) => {
  if (isPlaying.value) {
    currentAudio.value?.pause();
    isPlaying.value = false;
  } else {
    if (currentAudio.value) {
      currentAudio.value.play();
      isPlaying.value = true;
    } else {
      const url = getPrayerAudioUrl(prayerId);
      if (!url) return;
      
      const audio = new Audio(url);
      currentAudio.value = audio;
      audio.addEventListener('ended', () => {
        isPlaying.value = false;
        currentAudio.value = null;
      });
      audio.addEventListener('error', () => {
        isPlaying.value = false;
        currentAudio.value = null;
        alert('Could not play audio file.');
      });
      audio.play();
      isPlaying.value = true;
    }
  }
};

const stopAudio = () => {
  if (currentAudio.value) {
    currentAudio.value.pause();
    currentAudio.value = null;
  }
  isPlaying.value = false;
};

const openPrayer = (prayer: Prayer) => {
  stopAudio();
  selectedPrayer.value = prayer;
  showLatin.value = false;
  
  if (prayer.youtube && !prayer.default && !prayer.latin) {
    selectedPlayMode.value = 'video';
  } else {
    selectedPlayMode.value = 'audio';
  }
  
  document.body.style.overflow = 'hidden';
};

const closePrayer = () => {
  stopAudio();
  selectedPrayer.value = null;
  document.body.style.overflow = '';
};

onBeforeUnmount(() => {
  stopAudio();
});
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .transform,
.modal-leave-active .transform {
  transition: all 0.3s ease-out;
}

.modal-enter-from .transform,
.modal-leave-to .transform {
  opacity: 0;
  transform: scale(0.98);
}

.fade-content-enter-active,
.fade-content-leave-active {
  transition: opacity 0.25s ease;
}

.fade-content-enter-from,
.fade-content-leave-to {
  opacity: 0;
}

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
