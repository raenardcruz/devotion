<template>
  <div 
    class="relative text-stone-100 flex flex-col overflow-hidden select-none transition-all duration-300 font-sans"
    :class="[
      isDesktopFullscreen 
        ? 'fixed inset-0 z-50 w-screen h-screen' 
        : 'fixed inset-0 z-50 lg:relative lg:inset-auto lg:z-10 w-full h-full lg:rounded-3xl'
    ]"
  >
    <!-- Background: Pure Smooth Animated Gradient -->
    <div class="synth-gradient-bg absolute inset-0 z-0 pointer-events-none">
      <div 
        class="synth-glow-orb absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full transition-all duration-300 pointer-events-none"
        :style="{
          width: isGeminiSpeaking ? '520px' : isUserSpeaking ? '460px' : '360px',
          height: isGeminiSpeaking ? '520px' : isUserSpeaking ? '460px' : '360px',
          background: isGeminiSpeaking 
            ? 'radial-gradient(circle, rgba(34, 211, 238, 0.18) 0%, rgba(197, 160, 89, 0.14) 40%, transparent 70%)' 
            : isUserSpeaking 
              ? 'radial-gradient(circle, rgba(239, 68, 68, 0.18) 0%, rgba(114, 56, 61, 0.16) 40%, transparent 70%)' 
              : 'radial-gradient(circle, rgba(197, 160, 89, 0.14) 0%, rgba(114, 56, 61, 0.12) 50%, transparent 70%)',
          filter: 'blur(60px)'
        }"
      ></div>
    </div>

    <!-- Minimalist Top Bar: Status & Controls -->
    <header class="relative z-20 px-4 sm:px-6 py-3 flex items-center justify-between gap-3 shrink-0">
      
      <!-- Status Pill -->
      <div class="flex items-center gap-2.5 bg-black/30 backdrop-blur-md px-3.5 py-1.5 rounded-full border border-white/10 shadow-sm">
        <span 
          class="w-2 h-2 rounded-full transition-all duration-200"
          :class="isGeminiSpeaking 
            ? 'bg-cyan-400 animate-pulse shadow-[0_0_8px_#22D3EE]' 
            : isQueryingRAG 
              ? 'bg-amber-400 animate-ping shadow-[0_0_8px_#F59E0B]' 
              : isUserSpeaking 
                ? 'bg-red-400 animate-ping shadow-[0_0_8px_#EF4444]' 
                : 'bg-emerald-400 shadow-[0_0_6px_#10B981]'"
        ></span>
        <span class="font-serif font-semibold text-xs tracking-wide text-stone-200">
          {{ liveStatus }}
        </span>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-2">
        <!-- Mute Toggle -->
        <button
          type="button"
          @click="$emit('toggle-mute')"
          class="p-2 sm:px-3 sm:py-1.5 rounded-full text-xs font-serif font-medium border transition-all cursor-pointer flex items-center gap-1.5 backdrop-blur-md active:scale-95"
          :class="isMuted 
            ? 'bg-amber-500/20 text-amber-300 border-amber-500/40 shadow-xs' 
            : 'bg-black/30 text-stone-300 border-white/10 hover:bg-black/50 hover:text-white'"
          :title="isMuted ? 'Unmute microphone' : 'Mute microphone'"
        >
          <span class="text-xs">{{ isMuted ? '🔇' : '🎙️' }}</span>
          <span class="hidden sm:inline">{{ isMuted ? 'Unmute' : 'Mute' }}</span>
        </button>

        <!-- Desktop Fullscreen Toggle -->
        <button
          type="button"
          @click="toggleDesktopFullscreen"
          class="hidden lg:flex p-2 rounded-full text-stone-300 hover:text-white bg-black/30 hover:bg-black/50 border border-white/10 backdrop-blur-md transition-all cursor-pointer active:scale-95"
          :title="isDesktopFullscreen ? 'Restore to chat window' : 'Full screen'"
        >
          <svg v-if="isDesktopFullscreen" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="4 14 10 14 10 20"></polyline>
            <polyline points="20 10 14 10 14 4"></polyline>
            <line x1="14" y1="10" x2="21" y2="3"></line>
            <line x1="3" y1="21" x2="10" y2="14"></line>
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 3 21 3 21 9"></polyline>
            <polyline points="9 21 3 21 3 15"></polyline>
            <line x1="21" y1="3" x2="14" y2="10"></line>
            <line x1="3" y1="21" x2="10" y2="14"></line>
          </svg>
        </button>

        <!-- Exit / Disconnect Button -->
        <button
          type="button"
          @click="$emit('stop-live')"
          class="px-3.5 py-1.5 rounded-full text-xs font-serif font-semibold bg-red-600/80 hover:bg-red-600 text-white border border-red-500/80 backdrop-blur-md transition-all cursor-pointer shadow-sm active:scale-95 flex items-center gap-1.5"
          title="Disconnect Live Voice Session"
        >
          <span>End</span>
          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </header>

    <!-- Center: Pure Synthesizer Waveform (Occupies Entire Center) -->
    <div class="relative z-10 flex-1 flex flex-col items-center justify-center w-full min-h-0 overflow-hidden">
      
      <!-- Live Spoken Subtitle (Floats Cleanly Above the Wave) -->
      <transition name="fade">
        <div 
          v-if="liveTranscript || userLiveTranscript" 
          class="max-w-xl mx-auto px-6 text-center z-20 mb-2 sm:mb-4 pointer-events-none"
        >
          <!-- User Question Subtitle -->
          <p v-if="userLiveTranscript && !liveTranscript" class="text-xs sm:text-sm font-serif italic text-stone-300 drop-shadow-md">
            "{{ userLiveTranscript }}"
          </p>
          <!-- Gemini Response Subtitle -->
          <p v-else-if="liveTranscript" class="text-sm sm:text-base font-serif text-amber-100 drop-shadow-md leading-relaxed">
            {{ liveTranscript }}
            <span v-if="isGeminiSpeaking" class="inline-block w-1.5 h-3.5 ml-1 bg-amber-400 animate-pulse align-middle"></span>
          </p>
        </div>
      </transition>

      <!-- Wave Canvas Container (Occupies Center) -->
      <div class="relative w-full flex-1 flex items-center justify-center min-h-[180px] max-h-[380px]">
        <canvas 
          ref="canvasRef" 
          class="w-full h-full block cursor-pointer"
          @click="onWaveClick"
        ></canvas>

        <!-- RAG Searching Floating Pulse Badge -->
        <transition name="fade">
          <div 
            v-if="isQueryingRAG" 
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 px-4 py-2 rounded-full bg-black/60 border border-amber-400/50 backdrop-blur-md text-amber-200 text-xs font-serif flex items-center gap-2 shadow-lg z-20 pointer-events-none animate-pulse"
          >
            <span class="w-2 h-2 rounded-full bg-amber-400 animate-ping"></span>
            <span>Querying Church Magisterium for: <strong class="text-white">"{{ currentRagQuery }}"</strong>...</span>
          </div>
        </transition>

        <!-- Tap Wave Hint if Idle -->
        <div 
          v-if="!liveTranscript && !userLiveTranscript && !isQueryingRAG && activeCitationsList.length === 0" 
          class="absolute bottom-4 left-1/2 -translate-x-1/2 text-[11px] font-serif text-stone-400/80 pointer-events-none flex items-center gap-1.5"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-amber-400/60 animate-ping"></span>
          <span>Speak naturally to ask any Catholic question</span>
        </div>
      </div>
    </div>

    <!-- Results Display (Positioned at the Bottom as an Elegant Floating Card / Drawer) -->
    <div 
      v-if="activeCitationsList.length > 0" 
      class="relative z-20 w-full px-3 sm:px-6 pb-3 sm:pb-5 max-w-4xl mx-auto shrink-0 transition-all duration-300"
    >
      <div class="bg-black/50 border border-white/15 backdrop-blur-xl rounded-2xl sm:rounded-3xl shadow-2xl overflow-hidden transition-all duration-300">
        
        <!-- Results Header & Toggle -->
        <div 
          @click="isResultsOpen = !isResultsOpen"
          class="px-4 py-2.5 sm:py-3 flex items-center justify-between cursor-pointer select-none border-b border-white/10 hover:bg-white/5 transition-colors"
        >
          <div class="flex items-center gap-2">
            <span class="text-sm">📚</span>
            <span class="font-serif font-bold text-xs sm:text-sm text-[#E4C88B] tracking-wide">
              Magisterium References
            </span>
            <span class="px-2 py-0.5 rounded-full text-[10px] font-mono font-bold bg-[#C5A059]/20 text-[#E4C88B] border border-[#C5A059]/30">
              {{ activeCitationsList.length }} {{ activeCitationsList.length === 1 ? 'Source' : 'Sources' }}
            </span>
          </div>

          <div class="flex items-center gap-2 text-stone-400 text-xs font-serif">
            <span>{{ isResultsOpen ? 'Hide' : 'View Quotes' }}</span>
            <span class="text-xs transition-transform duration-200" :class="{ 'rotate-180': isResultsOpen }">
              ▼
            </span>
          </div>
        </div>

        <!-- Collapsible Results Content -->
        <transition name="slide">
          <div v-if="isResultsOpen" class="p-3 sm:p-4 space-y-2.5 max-h-56 sm:max-h-64 overflow-y-auto scrollbar-thin">
            <div 
              v-for="(cite, cIdx) in activeCitationsList" 
              :key="cIdx"
              class="p-3 rounded-xl bg-white/5 border border-white/10 hover:border-[#C5A059]/50 transition-colors space-y-1.5"
            >
              <!-- Source Title & Ref -->
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-1.5 mb-0.5">
                    <span class="text-[9px] font-mono uppercase font-bold px-1.5 py-0.2 rounded bg-[#72383D]/50 text-amber-200 border border-[#72383D]/60">
                      {{ getDocumentCategory(cite) }}
                    </span>
                    <span v-if="cite.author" class="font-serif text-[11px] text-stone-300 truncate">
                      {{ cite.author }}
                    </span>
                  </div>
                  <h5 class="font-serif font-bold text-xs sm:text-sm text-white">
                    {{ cite.document_title || cite.title || 'Church Document' }}
                  </h5>
                </div>

                <span v-if="cite.ref" class="font-mono text-[10px] font-bold px-2 py-0.5 rounded bg-[#C5A059]/20 text-[#E4C88B] border border-[#C5A059]/30 shrink-0">
                  {{ cite.ref }}
                </span>
              </div>

              <!-- Quote Excerpt -->
              <div v-if="cite.text" class="text-stone-200 font-serif text-xs leading-relaxed italic bg-black/30 p-2.5 rounded-lg border border-white/5">
                "{{ cite.text }}"
              </div>

              <!-- Actions -->
              <div class="flex items-center justify-between text-[10px] font-serif text-stone-400 pt-1">
                <span class="text-emerald-400">✓ Grounded Church Teaching</span>
                <button
                  type="button"
                  @click="copySingleCitation(cite)"
                  class="hover:text-white underline cursor-pointer"
                >
                  {{ copiedCitationIdx === cIdx ? 'Copied ✅' : 'Copy Quote' }}
                </button>
              </div>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import type { MagisteriumLiveCitation } from '../../composables/useGeminiLive';

interface ChatMessage {
  role: 'user' | 'assistant';
  content: string;
  citations?: any[];
}

const props = defineProps<{
  isLiveConnected: boolean;
  isConnecting: boolean;
  isMuted: boolean;
  isGeminiSpeaking: boolean;
  isUserSpeaking: boolean;
  isQueryingRAG: boolean;
  currentRagQuery: string;
  currentCitations: MagisteriumLiveCitation[];
  liveTranscript: string;
  userLiveTranscript: string;
  liveStatus: string;
  liveError: string | null;
  audioLevel: number;
  sessionMessages: ChatMessage[];
  activeSessionTitle?: string;
}>();

const emit = defineEmits<{
  (e: 'toggle-mute'): void;
  (e: 'stop-live'): void;
  (e: 'send-text', query: string): void;
  (e: 'interrupt'): void;
}>();

// UI States
const isDesktopFullscreen = ref(false);
const isResultsOpen = ref(true);
const copiedCitationIdx = ref<number | null>(null);

// Canvas Wave
const canvasRef = ref<HTMLCanvasElement | null>(null);
let animationFrameId: number | null = null;
let wavePhase = 0;

// Gather all citations from current turn and session history
const activeCitationsList = computed<MagisteriumLiveCitation[]>(() => {
  if (props.currentCitations && props.currentCitations.length > 0) {
    return props.currentCitations;
  }
  // Fallback to latest session citations
  const all: MagisteriumLiveCitation[] = [];
  const seen = new Set<string>();
  for (let i = props.sessionMessages.length - 1; i >= 0; i--) {
    const msg = props.sessionMessages[i];
    if (msg && msg.role === 'assistant' && msg.citations && Array.isArray(msg.citations)) {
      msg.citations.forEach((c: any) => {
        const key = (c.document_title || c.title || '') + (c.ref || '');
        if (!seen.has(key)) {
          seen.add(key);
          all.push(c);
        }
      });
    }
  }
  return all;
});

// Auto-open results when new citations arrive
watch(() => props.currentCitations.length, (newLen) => {
  if (newLen > 0) {
    isResultsOpen.value = true;
  }
});

const toggleDesktopFullscreen = () => {
  isDesktopFullscreen.value = !isDesktopFullscreen.value;
};

const onWaveClick = () => {
  if (props.isGeminiSpeaking) {
    emit('interrupt');
  } else {
    emit('toggle-mute');
  }
};

const getDocumentCategory = (cite: MagisteriumLiveCitation): string => {
  const title = (cite.document_title || cite.title || '').toLowerCase();
  if (title.includes('catechism') || title.includes('ccc')) return 'Catechism';
  if (title.includes('council') || title.includes('vatican') || title.includes('trent')) return 'Council';
  if (title.includes('encyclical') || (cite.author || '').toLowerCase().includes('pope')) return 'Encyclical';
  if (title.includes('gospel') || title.includes('epistle') || title.includes('genesis') || title.includes('matthew')) return 'Scripture';
  return 'Doctrine';
};

const copySingleCitation = async (cite: MagisteriumLiveCitation) => {
  const text = `${cite.document_title || cite.title || 'Church Document'} (${cite.ref || ''}):\n"${cite.text || ''}"`;
  try {
    await navigator.clipboard.writeText(text);
    const idx = activeCitationsList.value.indexOf(cite);
    copiedCitationIdx.value = idx;
    setTimeout(() => {
      copiedCitationIdx.value = null;
    }, 2000);
  } catch {}
};

// Canvas Synthesizer Wave Animation (Harmonic Smooth Fluid Wave)
const renderWave = () => {
  const canvas = canvasRef.value;
  if (!canvas) {
    animationFrameId = requestAnimationFrame(renderWave);
    return;
  }

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const rect = canvas.getBoundingClientRect();
  const dpr = window.devicePixelRatio || 1;
  const displayWidth = Math.round(rect.width * dpr);
  const displayHeight = Math.round(rect.height * dpr);

  if (canvas.width !== displayWidth || canvas.height !== displayHeight) {
    canvas.width = displayWidth;
    canvas.height = displayHeight;
  }

  const width = canvas.width;
  const height = canvas.height;
  const centerY = height / 2;

  ctx.clearRect(0, 0, width, height);

  wavePhase += 0.045;

  const level = props.isMuted ? 0 : props.audioLevel;
  const isSpeaking = props.isGeminiSpeaking;
  const isQuerying = props.isQueryingRAG;
  const isUser = props.isUserSpeaking;

  // Render 3 Interwoven Harmonic Wave Strands for Silky Smooth Organic Flow
  const strands = [
    {
      color: isSpeaking 
        ? 'rgba(34, 211, 238, 0.85)' 
        : isUser 
          ? 'rgba(239, 68, 68, 0.85)' 
          : isQuerying 
            ? 'rgba(245, 158, 11, 0.85)' 
            : 'rgba(228, 200, 139, 0.75)',
      lineWidth: 3 * dpr,
      freq: 1.0,
      phaseOffset: 0,
      ampMultiplier: 1.0,
    },
    {
      color: isSpeaking 
        ? 'rgba(197, 160, 89, 0.65)' 
        : isUser 
          ? 'rgba(248, 113, 113, 0.6)' 
          : isQuerying 
            ? 'rgba(251, 191, 36, 0.65)' 
            : 'rgba(197, 160, 89, 0.45)',
      lineWidth: 2 * dpr,
      freq: 1.8,
      phaseOffset: 1.4,
      ampMultiplier: 0.65,
    },
    {
      color: isSpeaking 
        ? 'rgba(6, 182, 212, 0.45)' 
        : isUser 
          ? 'rgba(220, 38, 38, 0.45)' 
          : 'rgba(245, 230, 200, 0.3)',
      lineWidth: 1.5 * dpr,
      freq: 0.6,
      phaseOffset: 2.8,
      ampMultiplier: 0.45,
    },
  ];

  strands.forEach((strand) => {
    ctx.beginPath();
    ctx.lineWidth = strand.lineWidth;
    ctx.strokeStyle = strand.color;

    const points = 120;
    const sliceWidth = width / points;

    for (let i = 0; i <= points; i++) {
      const normalizedX = (i / points); // 0 to 1
      
      // Windowing envelope so edges taper smoothly to zero at ends
      const envelope = Math.sin(normalizedX * Math.PI);

      let rawAmp = 0;
      if (isSpeaking) {
        // Formant synthesis wave
        rawAmp = (Math.sin(normalizedX * 10 * strand.freq + wavePhase * 2 + strand.phaseOffset) * 45 +
                  Math.cos(normalizedX * 16 * strand.freq - wavePhase * 1.5) * 22) * dpr;
      } else if (isUser && level > 0.02) {
        // Voice microphone wave
        rawAmp = (Math.sin(normalizedX * 12 * strand.freq + wavePhase * 3.5 + strand.phaseOffset) * (level * 95) +
                  Math.sin(normalizedX * 24 * strand.freq - wavePhase * 2) * (level * 40)) * dpr;
      } else if (isQuerying) {
        // Church search radar pulse
        rawAmp = Math.sin(normalizedX * 8 * strand.freq + wavePhase * 2.8 + strand.phaseOffset) * 32 * dpr;
      } else {
        // Soothing natural breathing baseline
        rawAmp = Math.sin(normalizedX * 4 * strand.freq + wavePhase * 0.8 + strand.phaseOffset) * 14 * dpr;
      }

      const y = centerY + (rawAmp * envelope * strand.ampMultiplier);
      const x = i * sliceWidth;

      if (i === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    }

    ctx.stroke();
  });

  animationFrameId = requestAnimationFrame(renderWave);
};

onMounted(() => {
  animationFrameId = requestAnimationFrame(renderWave);
});

onUnmounted(() => {
  if (animationFrameId !== null) {
    cancelAnimationFrame(animationFrameId);
  }
});
</script>

<style scoped>
/* Ultra-smooth Animated Gradient Background */
.synth-gradient-bg {
  background: linear-gradient(-45deg, #120f0e, #36191d, #141110, #2b1f13);
  background-size: 300% 300%;
  animation: ambientGradient 16s ease-in-out infinite alternate;
}

@keyframes ambientGradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease-out;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
  transform: translateY(8px);
}
</style>