<template>
  <div
    class="magisterium-live-root relative text-[#EFE9E1] flex flex-col overflow-hidden select-none transition-all duration-300 font-sans"
    :class="[
      isDesktopFullscreen
        ? 'fixed inset-0 z-50 w-screen h-screen'
        : 'fixed inset-0 z-50 lg:relative lg:inset-auto lg:z-10 w-full h-full lg:rounded-3xl'
    ]"
  >
    <!-- Retained Animated Gradient Background -->
    <div class="synth-gradient-bg absolute inset-0 z-0 pointer-events-none">
      <!-- Bottom-up Luminous Glow (Video-inspired, App Theme Burgundy & Sacred Gold) -->
      <div
        class="synth-bottom-glow absolute inset-x-0 bottom-0 h-[65%] transition-all duration-700 pointer-events-none"
        :class="{
          'is-active': isGeminiSpeaking || isQueryingRAG || isUserSpeaking
        }"
      ></div>

      <!-- Ambient Center Glow behind Sphere -->
      <div
        class="synth-glow-orb absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full transition-all duration-500 pointer-events-none"
        :style="{
          width: isGeminiSpeaking ? '440px' : isUserSpeaking ? '390px' : isQueryingRAG ? '420px' : '320px',
          height: isGeminiSpeaking ? '440px' : isUserSpeaking ? '390px' : isQueryingRAG ? '420px' : '320px',
          background: isGeminiSpeaking
            ? 'radial-gradient(circle, rgba(228, 200, 139, 0.24) 0%, rgba(114, 56, 61, 0.18) 45%, transparent 72%)'
            : isUserSpeaking
              ? 'radial-gradient(circle, rgba(139, 70, 76, 0.28) 0%, rgba(197, 160, 89, 0.16) 45%, transparent 72%)'
              : isQueryingRAG
                ? 'radial-gradient(circle, rgba(197, 160, 89, 0.28) 0%, rgba(114, 56, 61, 0.2) 45%, transparent 72%)'
                : 'radial-gradient(circle, rgba(197, 160, 89, 0.15) 0%, rgba(114, 56, 61, 0.14) 50%, transparent 72%)',
          filter: 'blur(55px)'
        }"
      ></div>
    </div>

    <!-- Minimalist Top Bar: Status & Controls -->
    <header class="relative z-20 px-4 sm:px-6 pt-3.5 pb-2 flex items-center justify-between gap-3 shrink-0">
      <!-- Status Pill -->
      <div class="glass-pill flex items-center gap-2.5 px-3.5 py-1.5 rounded-full">
        <span
          class="w-2 h-2 rounded-full transition-all duration-200"
          :class="isGeminiSpeaking
            ? 'bg-[#E4C88B] animate-pulse shadow-[0_0_10px_#E4C88B]'
            : isQueryingRAG
              ? 'bg-amber-400 animate-ping shadow-[0_0_10px_#F59E0B]'
              : isUserSpeaking
                ? 'bg-[#8B464C] animate-ping shadow-[0_0_10px_#8B464C]'
                : 'bg-emerald-400 shadow-[0_0_8px_#10B981]'"
        ></span>
        <span class="font-serif font-medium text-xs tracking-wide text-[#EFE9E1]/90">
          {{ liveStatus }}
        </span>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-2">
        <!-- Mute Toggle -->
        <button
          type="button"
          @click="$emit('toggle-mute')"
          class="glass-pill min-h-[38px] px-3 py-1.5 rounded-full text-xs font-serif font-medium transition-all cursor-pointer flex items-center gap-1.5 active:scale-95"
          :class="isMuted
            ? 'bg-amber-500/20! text-amber-200 border-amber-400/40!'
            : 'text-[#EFE9E1]/85 hover:text-white hover:bg-white/10'"
          :title="isMuted ? 'Unmute microphone' : 'Mute microphone'"
        >
          <svg v-if="isMuted" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="1" y1="1" x2="23" y2="23"></line>
            <path d="M9 9v3a3 3 0 0 0 5.12 2.12M15 9.34V4a3 3 0 0 0-5.94-.6"></path>
            <path d="M17 16.95A7 7 0 0 1 5 12v-2m14 0v2a7 7 0 0 1-.11 1.23"></path>
            <line x1="12" y1="19" x2="12" y2="23"></line>
            <line x1="8" y1="23" x2="16" y2="23"></line>
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z"></path>
            <path d="M19 10v2a7 7 0 0 1-14 0v-2"></path>
            <line x1="12" y1="19" x2="12" y2="23"></line>
            <line x1="8" y1="23" x2="16" y2="23"></line>
          </svg>
          <span class="hidden sm:inline">{{ isMuted ? 'Unmute' : 'Mute' }}</span>
        </button>

        <!-- Desktop Fullscreen Toggle -->
        <button
          type="button"
          @click="toggleDesktopFullscreen"
          class="glass-pill hidden lg:flex min-h-[38px] min-w-[38px] items-center justify-center rounded-full text-[#EFE9E1]/80 hover:text-white hover:bg-white/10 transition-all cursor-pointer active:scale-95"
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
          class="min-h-[38px] px-3.5 py-1.5 rounded-full text-xs font-serif font-semibold bg-[#72383D]/85 hover:bg-[#8B464C] text-[#EFE9E1] border border-[#E4C88B]/25 backdrop-blur-md transition-all cursor-pointer shadow-sm active:scale-95 flex items-center gap-1.5"
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

    <!-- Center Stage: Live Spoken Subtitle + Animated 3D Sphere (Always Visible) -->
    <div class="relative z-10 flex-1 flex flex-col items-center justify-center w-full min-h-0 overflow-hidden px-4 sm:px-6">
      <!-- Live Spoken Subtitle (Floats Cleanly Above the Sphere, just like previously) -->
      <transition name="fade-up" mode="out-in">
        <div
          v-if="liveTranscript || userLiveTranscript"
          key="live-subtitles"
          class="max-w-2xl mx-auto px-4 text-center z-20 mb-2 sm:mb-3 pointer-events-none shrink-0"
        >
          <!-- User Question Subtitle -->
          <p v-if="userLiveTranscript && !liveTranscript" class="text-xs sm:text-sm font-serif italic text-[#D1C7BD] drop-shadow-md">
            “{{ userLiveTranscript }}”
          </p>
          <!-- AI Live Spoken Subtitle -->
          <p v-else-if="liveTranscript" class="text-sm sm:text-base font-serif text-[#F5E6C8] drop-shadow-md leading-relaxed max-h-28 overflow-y-auto scrollbar-thin">
            {{ liveTranscript }}
            <span v-if="isGeminiSpeaking" class="inline-block w-1.5 h-3.5 ml-1 bg-[#E4C88B] animate-pulse align-middle rounded-xs"></span>
          </p>
        </div>
      </transition>

      <!-- Animated 3D Energy Sphere Canvas (Never covered by sources) -->
      <div
        class="relative flex items-center justify-center shrink-0 transition-all duration-300"
        :class="activeCitationsList.length > 0 && isResultsOpen
          ? 'w-[140px] h-[140px] sm:w-[175px] sm:h-[175px]'
          : 'w-[190px] h-[190px] sm:w-[235px] sm:h-[235px]'"
      >
        <canvas
          ref="canvasRef"
          class="w-full h-full block cursor-pointer"
          @click="onSphereClick"
          :title="isGeminiSpeaking ? 'Tap to interrupt' : isMuted ? 'Tap to unmute' : 'Tap to mute'"
        ></canvas>
      </div>

      <!-- Idle / Searching Headline Below Sphere -->
      <transition name="fade-up" mode="out-in">
        <div v-if="isQueryingRAG" key="searching" class="mt-2 text-center max-w-md px-4 shrink-0">
          <p class="text-[11px] sm:text-xs font-sans tracking-wide text-[#E4C88B]/85 mb-0.5">
            Searching Magisterium
          </p>
          <p v-if="currentRagQuery" class="font-serif text-sm sm:text-base font-medium text-[#EFE9E1] truncate">
            “{{ currentRagQuery }}”
          </p>
        </div>

        <div
          v-else-if="!liveTranscript && !userLiveTranscript && activeCitationsList.length === 0"
          key="idle"
          class="mt-2 text-center max-w-md px-4 shrink-0"
        >
          <p class="text-[11px] sm:text-xs font-sans tracking-wide text-[#D1C7BD]/75 mb-0.5">
            {{ activeSessionTitle && activeSessionTitle !== 'New Sanctuary' ? activeSessionTitle : 'Magisterium Live' }}
          </p>
          <h2 class="font-serif text-xl sm:text-2xl font-semibold text-[#EFE9E1] tracking-tight leading-tight">
            How can I help you today?
          </h2>
        </div>
      </transition>
    </div>

    <!-- Bottom Stage: Glass RAG Sources (Below Sphere) + Floating Glass Input Pill -->
    <div class="relative z-20 w-full px-3 sm:px-6 pb-3.5 sm:pb-5 max-w-3xl mx-auto shrink-0 flex flex-col gap-2.5">
      <!-- RAG Results Glass Drawer (Docked cleanly below the sphere, never covering it) -->
      <div
        v-if="activeCitationsList.length > 0"
        class="glass-sheet-surface rounded-2xl sm:rounded-3xl overflow-hidden transition-all duration-300"
      >
        <!-- Results Header & Collapse Toggle -->
        <div
          @click="isResultsOpen = !isResultsOpen"
          class="px-4 py-2.5 flex items-center justify-between gap-2 cursor-pointer select-none border-b border-white/10 hover:bg-white/5 transition-colors"
        >
          <div class="flex items-center gap-2 min-w-0">
            <span class="font-serif font-semibold text-xs sm:text-sm text-[#E4C88B] tracking-wide truncate">
              Magisterium Sources
            </span>
            <span class="px-2 py-0.5 rounded-full text-[10px] font-mono font-semibold bg-white/10 text-[#E4C88B] border border-[#E4C88B]/25 shrink-0">
              {{ activeCitationsList.length }} {{ activeCitationsList.length === 1 ? 'Source' : 'Sources' }}
            </span>
          </div>

          <div class="flex items-center gap-2.5 shrink-0">
            <!-- User Query Pill (Matching video top-right query pill) -->
            <span
              v-if="displayQuery"
              class="hidden sm:inline-block px-2.5 py-0.5 rounded-full bg-white/10 border border-white/15 text-[11px] font-sans text-[#EFE9E1]/90 max-w-[200px] truncate"
              :title="displayQuery"
            >
              {{ displayQuery }}
            </span>
            <span class="text-[#D1C7BD] text-xs font-sans flex items-center gap-1">
              <span class="hidden sm:inline">{{ isResultsOpen ? 'Hide' : 'Show' }}</span>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="13"
                height="13"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="transition-transform duration-200"
                :class="{ 'rotate-180': isResultsOpen }"
              >
                <polyline points="6 9 12 15 18 9"></polyline>
              </svg>
            </span>
          </div>
        </div>

        <!-- Collapsible Frosted-Glass RAG Cards -->
        <transition name="slide">
          <div v-if="isResultsOpen" class="p-3 sm:p-3.5 space-y-2.5 max-h-48 sm:max-h-56 overflow-y-auto scrollbar-thin">
            <article
              v-for="(cite, cIdx) in activeCitationsList"
              :key="cIdx"
              class="glass-rag-card p-3 sm:p-3.5 rounded-2xl transition-all duration-200 space-y-1.5"
            >
              <!-- Source Title & Reference Pill -->
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-1.5 mb-0.5 flex-wrap">
                    <span class="text-[9px] font-mono uppercase font-bold px-2 py-0.5 rounded-full bg-[#72383D]/70 text-[#EFE9E1] border border-[#E4C88B]/25">
                      {{ getDocumentCategory(cite) }}
                    </span>
                    <span v-if="cite.author" class="font-serif text-[11px] text-[#D1C7BD] truncate">
                      {{ cite.author }}
                    </span>
                  </div>
                  <h5 class="font-serif font-bold text-xs sm:text-sm text-white leading-snug">
                    {{ cite.document_title || cite.title || 'Church Document' }}
                  </h5>
                </div>

                <span
                  v-if="cite.ref"
                  class="font-mono text-[10px] font-semibold px-2 py-0.5 rounded-full bg-white/10 text-[#E4C88B] border border-white/15 shrink-0"
                >
                  {{ cite.ref }}
                </span>
              </div>

              <!-- Quote Excerpt -->
              <div
                v-if="cite.text"
                class="text-[#EFE9E1]/90 font-serif text-xs leading-relaxed italic bg-black/25 px-3 py-2 rounded-xl border border-white/8"
              >
                “{{ cite.text }}”
              </div>

              <!-- Footer Actions -->
              <div class="flex items-center justify-between text-[10px] font-sans text-[#D1C7BD]/80 pt-0.5">
                <span class="text-emerald-300/90 flex items-center gap-1">
                  <span>✓</span>
                  <span>Grounded Church Teaching</span>
                </span>
                <div class="flex items-center gap-3">
                  <a
                    v-if="cite.url"
                    :href="cite.url"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="hover:text-[#E4C88B] underline underline-offset-2 transition-colors"
                  >
                    Source ↗
                  </a>
                  <button
                    type="button"
                    @click="copySingleCitation(cite)"
                    class="glass-text-btn hover:text-white underline underline-offset-2 cursor-pointer transition-colors"
                  >
                    {{ copiedCitationIdx === cIdx ? 'Copied ✓' : 'Copy' }}
                  </button>
                </div>
              </div>
            </article>
          </div>
        </transition>
      </div>

      <!-- Morphing Bottom Glass Pill ("Searching..." vs "Ask anything...") -->
      <div class="flex justify-center w-full">
        <transition name="pill-morph" mode="out-in">
          <div
            v-if="isQueryingRAG"
            key="searching-pill"
            class="glass-input-pill px-6 py-2.5 rounded-full flex items-center gap-2.5 shadow-xl"
          >
            <span class="w-2 h-2 rounded-full bg-[#E4C88B] animate-ping"></span>
            <span class="text-xs sm:text-sm font-sans font-medium text-[#EFE9E1]/90 tracking-wide">
              Searching…
            </span>
          </div>

          <form
            v-else
            key="input-pill"
            @submit.prevent="handleTextSubmit"
            class="glass-input-pill w-full min-h-[48px] pl-4 sm:pl-5 pr-1.5 py-1 rounded-full flex items-center gap-2 shadow-2xl transition-all duration-300 focus-within:border-[#E4C88B]/45"
          >
            <input
              v-model="draftQuery"
              type="text"
              placeholder="Ask anything..."
              class="flex-1 bg-transparent border-none outline-none text-sm text-[#EFE9E1] placeholder-[#D1C7BD]/60 font-sans min-w-0"
              :disabled="!isLiveConnected && isConnecting"
            />
            <button
              v-if="draftQuery.trim()"
              type="submit"
              class="min-w-[36px] min-h-[36px] rounded-full bg-[#72383D] hover:bg-[#8B464C] text-[#EFE9E1] border border-[#E4C88B]/30 flex items-center justify-center transition-all cursor-pointer active:scale-95 shrink-0"
              title="Send question"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="22" y1="2" x2="11" y2="13"></line>
                <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
              </svg>
            </button>
            <button
              v-else
              type="button"
              @click="onSphereClick"
              class="glass-icon-btn min-w-[36px] min-h-[36px] rounded-full text-[#EFE9E1]/90 flex items-center justify-center transition-all cursor-pointer active:scale-95 shrink-0"
              :title="isGeminiSpeaking ? 'Interrupt voice' : isMuted ? 'Unmute microphone' : 'Listening (tap to mute)'"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="{ 'text-[#E4C88B] animate-pulse': isUserSpeaking || isGeminiSpeaking }">
                <line x1="4" y1="10" x2="4" y2="14"></line>
                <line x1="8" y1="7" x2="8" y2="17"></line>
                <line x1="12" y1="4" x2="12" y2="20"></line>
                <line x1="16" y1="7" x2="16" y2="17"></line>
                <line x1="20" y1="10" x2="20" y2="14"></line>
              </svg>
            </button>
          </form>
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
const draftQuery = ref('');

// Canvas Animated Sphere
const canvasRef = ref<HTMLCanvasElement | null>(null);
let animationFrameId: number | null = null;
let spherePhase = 0;
let smoothedAudio = 0;

// Latest user query from session history (fallback after turn_complete)
const latestSessionUserQuery = computed(() => {
  for (let i = props.sessionMessages.length - 1; i >= 0; i--) {
    const msg = props.sessionMessages[i];
    if (msg && msg.role === 'user' && msg.content) {
      return msg.content;
    }
  }
  return '';
});

// Gather all citations from current turn or latest session turn
const activeCitationsList = computed<MagisteriumLiveCitation[]>(() => {
  if (props.currentCitations && props.currentCitations.length > 0) {
    return props.currentCitations;
  }
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
      if (all.length > 0) break;
    }
  }
  return all;
});

const displayQuery = computed(() => {
  return props.userLiveTranscript || props.currentRagQuery || latestSessionUserQuery.value;
});

// Auto-open results drawer when new citations arrive
watch(() => props.currentCitations.length, (newLen) => {
  if (newLen > 0) {
    isResultsOpen.value = true;
  }
});

const toggleDesktopFullscreen = () => {
  isDesktopFullscreen.value = !isDesktopFullscreen.value;
};

const onSphereClick = () => {
  if (props.isGeminiSpeaking) {
    emit('interrupt');
  } else {
    emit('toggle-mute');
  }
};

const handleTextSubmit = () => {
  const trimmed = draftQuery.value.trim();
  if (!trimmed) return;
  emit('send-text', trimmed);
  draftQuery.value = '';
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

// Center Animated 3D Energy Sphere (App Theme Gold/Burgundy + Video-inspired Plasma Rings)
const renderSphere = () => {
  const canvas = canvasRef.value;
  if (!canvas) {
    animationFrameId = requestAnimationFrame(renderSphere);
    return;
  }

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const rect = canvas.getBoundingClientRect();
  const dpr = window.devicePixelRatio || 1;
  const displayWidth = Math.max(1, Math.round(rect.width * dpr));
  const displayHeight = Math.max(1, Math.round(rect.height * dpr));

  if (canvas.width !== displayWidth || canvas.height !== displayHeight) {
    canvas.width = displayWidth;
    canvas.height = displayHeight;
  }

  const width = canvas.width;
  const height = canvas.height;
  const cx = width / 2;
  const cy = height / 2;

  ctx.clearRect(0, 0, width, height);

  const rawLevel = props.isMuted ? 0 : props.audioLevel;
  smoothedAudio += (rawLevel - smoothedAudio) * 0.16;

  const isSpeaking = props.isGeminiSpeaking;
  const isQuerying = props.isQueryingRAG;
  const isUser = props.isUserSpeaking;

  const speed = isQuerying ? 0.045 : isSpeaking ? 0.032 : isUser ? 0.028 : 0.014;
  spherePhase += speed;

  const baseRadius = Math.min(width, height) * 0.36;
  const breathe = Math.sin(spherePhase * 1.3) * (baseRadius * 0.025);
  const voiceBoost = isSpeaking
    ? Math.sin(spherePhase * 5.2) * (baseRadius * 0.045)
    : isUser
      ? smoothedAudio * baseRadius * 0.22
      : 0;
  const R = baseRadius + breathe + voiceBoost;

  const primaryGlow = isSpeaking
    ? 'rgba(228, 200, 139, 0.85)'
    : isQuerying
      ? 'rgba(245, 196, 81, 0.9)'
      : isUser
        ? 'rgba(180, 92, 100, 0.9)'
        : 'rgba(228, 200, 139, 0.72)';

  const secondaryGlow = isSpeaking
    ? 'rgba(255, 244, 214, 0.75)'
    : isQuerying
      ? 'rgba(255, 220, 140, 0.8)'
      : isUser
        ? 'rgba(228, 200, 139, 0.65)'
        : 'rgba(139, 70, 76, 0.65)';

  // 1. Soft Outer & Inner Core Sphere Gradient
  const coreGrad = ctx.createRadialGradient(
    cx - R * 0.2,
    cy - R * 0.22,
    R * 0.08,
    cx,
    cy,
    R * 1.18
  );
  coreGrad.addColorStop(0, isSpeaking ? 'rgba(255, 245, 220, 0.28)' : 'rgba(228, 200, 139, 0.16)');
  coreGrad.addColorStop(0.45, isUser ? 'rgba(114, 56, 61, 0.28)' : 'rgba(114, 56, 61, 0.2)');
  coreGrad.addColorStop(0.85, 'rgba(30, 18, 16, 0.35)');
  coreGrad.addColorStop(1, 'rgba(18, 15, 14, 0)');

  ctx.fillStyle = coreGrad;
  ctx.beginPath();
  ctx.arc(cx, cy, R * 1.18, 0, Math.PI * 2);
  ctx.fill();

  // 2. Subtle Glass Rim around the Sphere
  ctx.save();
  ctx.beginPath();
  ctx.arc(cx, cy, R, 0, Math.PI * 2);
  ctx.strokeStyle = isSpeaking
    ? 'rgba(255, 244, 214, 0.35)'
    : isQuerying
      ? 'rgba(245, 196, 81, 0.4)'
      : 'rgba(228, 200, 139, 0.22)';
  ctx.lineWidth = 1.5 * dpr;
  ctx.stroke();
  ctx.restore();

  // 3. 3D Orbital Energy Ribbons
  ctx.save();
  ctx.globalCompositeOperation = 'lighter';

  const ribbons = [
    { tiltX: 0.55, tiltZ: 0.3, speedMult: 1.0, phaseOff: 0, color: primaryGlow, width: 2.4 },
    { tiltX: -0.68, tiltZ: -0.5, speedMult: -0.85, phaseOff: 1.8, color: secondaryGlow, width: 2.0 },
    { tiltX: 1.05, tiltZ: 0.85, speedMult: 1.25, phaseOff: 3.4, color: 'rgba(255, 248, 232, 0.55)', width: 1.5 },
    { tiltX: -0.25, tiltZ: 1.2, speedMult: 0.7, phaseOff: 4.9, color: primaryGlow, width: 1.8 },
  ];

  const steps = 84;
  ribbons.forEach((ribbon, rIdx) => {
    const rotY = spherePhase * ribbon.speedMult + ribbon.phaseOff;
    const cosY = Math.cos(rotY);
    const sinY = Math.sin(rotY);
    const cosX = Math.cos(ribbon.tiltX + Math.sin(spherePhase * 0.5 + rIdx) * 0.18);
    const sinX = Math.sin(ribbon.tiltX + Math.sin(spherePhase * 0.5 + rIdx) * 0.18);
    const cosZ = Math.cos(ribbon.tiltZ);
    const sinZ = Math.sin(ribbon.tiltZ);

    for (let i = 0; i < steps; i++) {
      const t1 = (i / steps) * Math.PI * 2;
      const t2 = ((i + 1) / steps) * Math.PI * 2;

      const project = (t: number) => {
        const waveAmp = isSpeaking || isUser
          ? Math.sin(t * 4 + spherePhase * 4) * (0.04 + smoothedAudio * 0.08)
          : Math.sin(t * 3 + spherePhase * 2) * 0.025;
        const rLocal = R * (0.96 + waveAmp);

        const x0 = Math.cos(t) * rLocal;
        const y0 = Math.sin(t * 2 + spherePhase) * rLocal * 0.06;
        const z0 = Math.sin(t) * rLocal;

        const x1 = x0 * cosY + z0 * sinY;
        const z1 = -x0 * sinY + z0 * cosY;
        const y1 = y0;

        const y2 = y1 * cosX - z1 * sinX;
        const z2 = y1 * sinX + z1 * cosX;
        const x2 = x1;

        const x3 = x2 * cosZ - y2 * sinZ;
        const y3 = x2 * sinZ + y2 * cosZ;

        return { x: cx + x3, y: cy + y3, z: z2 };
      };

      const p1 = project(t1);
      const p2 = project(t2);

      const depthNorm = ((p1.z / R) + 1) * 0.5;
      const alpha = 0.08 + Math.pow(Math.max(0, depthNorm), 1.8) * 0.85;

      ctx.beginPath();
      ctx.moveTo(p1.x, p1.y);
      ctx.lineTo(p2.x, p2.y);
      ctx.strokeStyle = ribbon.color;
      ctx.globalAlpha = alpha;
      ctx.lineWidth = (ribbon.width * (0.55 + depthNorm * 0.75)) * dpr;
      ctx.lineCap = 'round';
      ctx.stroke();
    }
  });

  // Bright Specular Sparkle Node on Front Arc
  const sparkAngle = spherePhase * 1.4;
  const sx = cx + Math.cos(sparkAngle) * R * 0.78;
  const sy = cy + Math.sin(sparkAngle * 0.7) * R * 0.42;
  const sparkGrad = ctx.createRadialGradient(sx, sy, 0, sx, sy, R * 0.22);
  sparkGrad.addColorStop(0, 'rgba(255, 250, 235, 0.85)');
  sparkGrad.addColorStop(0.35, 'rgba(228, 200, 139, 0.35)');
  sparkGrad.addColorStop(1, 'rgba(228, 200, 139, 0)');
  ctx.globalAlpha = 0.9;
  ctx.fillStyle = sparkGrad;
  ctx.beginPath();
  ctx.arc(sx, sy, R * 0.22, 0, Math.PI * 2);
  ctx.fill();

  ctx.restore();

  animationFrameId = requestAnimationFrame(renderSphere);
};

onMounted(() => {
  animationFrameId = requestAnimationFrame(renderSphere);
});

onUnmounted(() => {
  if (animationFrameId !== null) {
    cancelAnimationFrame(animationFrameId);
  }
});
</script>

<style scoped>
/* Retained Animated Gradient Background */
.synth-gradient-bg {
  background: linear-gradient(-45deg, #120f0e, #36191d, #141110, #2b1f13);
  background-size: 300% 300%;
  animation: ambientGradient 16s ease-in-out infinite alternate;
}

/* Bottom-up Luminous Glow matching Video Composition in App Burgundy/Gold */
.synth-bottom-glow {
  background: radial-gradient(
    ellipse 90% 62% at 50% 108%,
    rgba(114, 56, 61, 0.58) 0%,
    rgba(197, 160, 89, 0.22) 38%,
    rgba(50, 45, 41, 0.08) 65%,
    transparent 85%
  );
  opacity: 0.85;
}

.synth-bottom-glow.is-active {
  background: radial-gradient(
    ellipse 95% 68% at 50% 105%,
    rgba(139, 70, 76, 0.72) 0%,
    rgba(228, 200, 139, 0.34) 40%,
    transparent 82%
  );
  opacity: 1;
}

/* Glassmorphic Surfaces (Apple HIG Glassmorphism) */
.glass-pill {
  background: rgba(25, 20, 19, 0.42);
  backdrop-filter: blur(18px) saturate(160%);
  -webkit-backdrop-filter: blur(18px) saturate(160%);
  border: 1px solid rgba(255, 255, 255, 0.14);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.22), inset 0 1px 0 rgba(255, 255, 255, 0.16);
}

.magisterium-live-root button.glass-pill:hover {
  background-color: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.24);
}

.magisterium-live-root button.glass-icon-btn {
  background-color: transparent;
  border-color: transparent;
  padding: 0;
  min-height: 36px;
}

.magisterium-live-root button.glass-icon-btn:hover {
  background-color: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.18);
}

.magisterium-live-root button.glass-text-btn {
  background-color: transparent;
  border-color: transparent;
  padding: 0;
  min-height: 0;
  color: inherit;
}

.magisterium-live-root button.glass-text-btn:hover {
  background-color: transparent;
  border-color: transparent;
}

.glass-input-pill {
  background: rgba(255, 255, 255, 0.09);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.18);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.35), inset 0 1px 0 rgba(255, 255, 255, 0.24);
}

.glass-sheet-surface {
  background: linear-gradient(
    180deg,
    rgba(38, 28, 26, 0.52) 0%,
    rgba(28, 20, 19, 0.64) 100%
  );
  backdrop-filter: blur(26px) saturate(180%);
  -webkit-backdrop-filter: blur(26px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.18);
  box-shadow:
    0 20px 50px rgba(0, 0, 0, 0.45),
    inset 0 1px 0 rgba(255, 255, 255, 0.22);
}

.glass-rag-card {
  background: rgba(255, 255, 255, 0.07);
  backdrop-filter: blur(14px) saturate(150%);
  -webkit-backdrop-filter: blur(14px) saturate(150%);
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.18), inset 0 1px 0 rgba(255, 255, 255, 0.14);
}

.glass-rag-card:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(228, 200, 139, 0.35);
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

/* Fluid Spring-inspired Transitions */
.fade-up-enter-active,
.fade-up-leave-active {
  transition: opacity 0.32s cubic-bezier(0.16, 1, 0.3, 1), transform 0.32s cubic-bezier(0.16, 1, 0.3, 1);
}
.fade-up-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.fade-up-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

.pill-morph-enter-active,
.pill-morph-leave-active {
  transition: opacity 0.28s cubic-bezier(0.16, 1, 0.3, 1), transform 0.28s cubic-bezier(0.16, 1, 0.3, 1);
}
.pill-morph-enter-from,
.pill-morph-leave-to {
  opacity: 0;
  transform: scale(0.94) translateY(6px);
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
  transform: translateY(8px);
}

/* Accessibility: Reduced Motion & Reduced Transparency */
@media (prefers-reduced-motion: reduce) {
  .synth-gradient-bg {
    animation: none !important;
  }
  *,
  ::before,
  ::after {
    transition-duration: 0.01ms !important;
  }
}

@media (prefers-reduced-transparency: reduce) {
  .glass-pill,
  .glass-input-pill,
  .glass-sheet-surface,
  .glass-rag-card {
    backdrop-filter: none !important;
    -webkit-backdrop-filter: none !important;
    background: #1d1715 !important;
  }
}
</style>
