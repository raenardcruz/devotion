<template>
  <div class="min-h-screen bg-parchment-bg text-parchment-neutral flex flex-col pb-20 selection:bg-parchment-primary/20">
    <TopNav />

    <main class="max-w-7xl mx-auto px-4 md:px-8 py-6 md:py-8 flex-grow w-full flex flex-col">
      <!-- Page Header -->
      <div class="text-center mb-6">
        <div class="inline-flex items-center space-x-2 px-3.5 py-1 rounded-full bg-parchment-neutral-light border border-parchment-border/70 mb-2.5 text-xs text-parchment-primary font-semibold uppercase tracking-wider shadow-2xs">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
          <span>Grounded in Sacred Tradition & Church Doctrine</span>
        </div>
        <h1 class="text-3xl md:text-4xl font-serif font-bold text-parchment-primary-dark tracking-tight mb-1">
          Magisterium AI Sanctuary
        </h1>
        <p class="text-parchment-neutral/70 text-xs md:text-sm max-w-2xl mx-auto font-serif">
          Engage in Catholic Q&A, Biblical Reflection, and Magisterial Search powered by Magisterium AI.
        </p>

        <!-- Status & Rate Limit Bar -->
        <div class="mt-4 max-w-2xl mx-auto bg-amber-50/80 border border-amber-200/90 rounded-2xl p-3.5 text-xs text-amber-950 flex items-center justify-between shadow-sm backdrop-blur-xs">
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 rounded-xl bg-amber-200/60 border border-amber-300 flex items-center justify-center text-amber-900 font-serif font-bold text-sm shadow-2xs">
              📜
            </div>
            <div class="text-left">
              <span class="font-bold text-amber-950 uppercase tracking-wider text-[10px] block">Magisterium API Quota & Usage</span>
              <span class="text-[11px] text-amber-900/90 font-serif">
                <span v-if="apiUsage?.remaining !== undefined">
                  Remaining Capacity: <strong>{{ apiUsage.remaining }}</strong> <span v-if="apiUsage.limit">/ {{ apiUsage.limit }}</span> queries
                </span>
                <span v-else>
                  Rate Limits: <strong>15 req/min</strong> (Search) • <strong>2 req/min</strong> (Q&A Pipeline)
                </span>
              </span>
            </div>
          </div>
          <a 
            href="https://www.magisterium.com" 
            target="_blank" 
            rel="noopener noreferrer" 
            class="text-[11px] bg-parchment-primary text-white font-medium px-3 py-1.5 rounded-xl hover:bg-parchment-primary-dark transition-all shadow-xs flex items-center space-x-1"
          >
            <span>Console</span>
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
          </a>
        </div>
      </div>

      <!-- Main Layout: Sidebar & Chat Container -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 flex-grow items-start">
        
        <!-- Sidebar Navigation: Chat Sessions -->
        <aside class="lg:col-span-4 xl:col-span-3 bg-parchment-neutral-light/80 border border-parchment-border rounded-3xl p-4 flex flex-col shadow-sm h-full max-h-[750px] backdrop-blur-xs">
          <!-- Sidebar Header -->
          <div class="flex items-center justify-between pb-3 mb-3 border-b border-parchment-border/60">
            <div class="flex items-center space-x-2">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-parchment-primary">
                <path d="M12 20h9"></path>
                <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
              </svg>
              <span class="font-serif font-bold text-xs uppercase tracking-wider text-parchment-primary-dark">Saved Conversations</span>
            </div>
            <button 
              @click="createNewSession" 
              class="text-xs bg-parchment-primary text-white px-3 py-1.5 rounded-xl font-semibold hover:bg-parchment-primary-dark transition-all shadow-xs flex items-center space-x-1"
            >
              <span>+ New Chat</span>
            </button>
          </div>

          <!-- Sessions Scroll List -->
          <div class="overflow-y-auto space-y-2.5 flex-1 scrollbar-thin pr-1 max-h-[620px]">
            <div 
              v-for="session in chatSessions" 
              :key="session.id"
              @click="selectSession(session.id)"
              class="p-3 rounded-2xl border text-left cursor-pointer transition-all flex items-center justify-between group relative overflow-hidden"
              :class="[activeSessionId === session.id ? 'bg-amber-100/90 border-amber-300 shadow-xs' : 'bg-parchment-bg/70 border-parchment-border/50 hover:bg-parchment-bg hover:border-parchment-border']"
            >
              <div class="min-w-0 flex-1 pr-2">
                <div class="text-xs font-semibold text-parchment-primary-dark truncate font-serif">
                  {{ session.title || 'New Conversation' }}
                </div>
                <div class="text-[10px] text-parchment-neutral/60 mt-1 flex items-center space-x-2">
                  <span>{{ formatDate(session.updatedAt) }}</span>
                  <span class="w-1 h-1 rounded-full bg-parchment-neutral/30"></span>
                  <span>{{ session.messages.length }} msgs</span>
                </div>
              </div>

              <!-- Delete Button -->
              <button 
                @click.stop="deleteSession(session.id)"
                class="opacity-0 group-hover:opacity-100 text-stone-400 hover:text-red-600 transition-all p-1.5 rounded-lg hover:bg-red-50"
                title="Delete Chat"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                </svg>
              </button>
            </div>

            <div v-if="chatSessions.length === 0" class="text-center py-8 px-4 text-xs text-parchment-neutral/50 italic space-y-2">
              <div class="text-2xl opacity-40">✍️</div>
              <p>No saved conversations yet. Click "+ New Chat" to begin a discussion.</p>
            </div>
          </div>
        </aside>

        <!-- Main Chat Box Panel -->
        <div class="lg:col-span-8 xl:col-span-9 flex flex-col bg-parchment-neutral-light/60 border border-parchment-border rounded-3xl shadow-sm overflow-hidden h-[750px]">
          
          <!-- Scrollable Chat Window -->
          <div ref="chatContainer" class="flex-1 p-5 md:p-8 overflow-y-auto space-y-6 scrollbar-thin">
            
            <!-- Welcome Card (Shows when session is empty) -->
            <div v-if="activeMessages.length === 0" class="max-w-2xl mx-auto my-6 space-y-6 text-center animate-fade-in">
              <div class="w-16 h-16 rounded-3xl bg-parchment-primary-dark text-white flex items-center justify-center font-serif text-2xl font-bold mx-auto shadow-md">
                M
              </div>
              <div class="bg-parchment-neutral-light border border-parchment-border rounded-3xl p-6 shadow-xs space-y-3">
                <h3 class="font-serif text-lg font-bold text-parchment-primary-dark">Pax Vobiscum! Welcome to Magisterium AI</h3>
                <p class="text-sm text-parchment-neutral/80 font-serif leading-relaxed">
                  Search through the Catechism of the Catholic Church, Ecumenical Council texts, Papal Encyclicals, and Sacred Scripture with AI answers grounded in Church Magisterium.
                </p>
                
                <div class="pt-3 border-t border-parchment-border/50 text-left space-y-2">
                  <span class="text-xs font-bold uppercase tracking-wider text-parchment-primary-dark block mb-2">Suggested Reflection Questions:</span>
                  <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
                    <button 
                      @click="askSuggested('What does the Catechism teach about prayer?')" 
                      class="text-xs text-left p-3 rounded-2xl bg-parchment-bg hover:bg-amber-100/70 text-parchment-primary-dark border border-parchment-border/70 transition-all font-serif flex items-center justify-between group shadow-2xs"
                    >
                      <span>"What does the Catechism teach about prayer?"</span>
                      <span class="opacity-0 group-hover:opacity-100 transition-opacity">→</span>
                    </button>
                    <button 
                      @click="askSuggested('Explain the relationship between Faith and Reason (Fides et Ratio).')" 
                      class="text-xs text-left p-3 rounded-2xl bg-parchment-bg hover:bg-amber-100/70 text-parchment-primary-dark border border-parchment-border/70 transition-all font-serif flex items-center justify-between group shadow-2xs"
                    >
                      <span>"Explain Faith and Reason (Fides et Ratio)."</span>
                      <span class="opacity-0 group-hover:opacity-100 transition-opacity">→</span>
                    </button>
                    <button 
                      @click="askSuggested('What are the corporal and spiritual works of mercy?')" 
                      class="text-xs text-left p-3 rounded-2xl bg-parchment-bg hover:bg-amber-100/70 text-parchment-primary-dark border border-parchment-border/70 transition-all font-serif flex items-center justify-between group shadow-2xs"
                    >
                      <span>"What are the works of mercy?"</span>
                      <span class="opacity-0 group-hover:opacity-100 transition-opacity">→</span>
                    </button>
                    <button 
                      @click="askSuggested('How does Saint Thomas Aquinas explain Grace and Free Will?')" 
                      class="text-xs text-left p-3 rounded-2xl bg-parchment-bg hover:bg-amber-100/70 text-parchment-primary-dark border border-parchment-border/70 transition-all font-serif flex items-center justify-between group shadow-2xs"
                    >
                      <span>"How does St. Thomas explain Grace?"</span>
                      <span class="opacity-0 group-hover:opacity-100 transition-opacity">→</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Active Messages Stream -->
            <div v-for="(msg, index) in activeMessages" :key="index" class="flex flex-col space-y-2 animate-fade-in">
              <!-- User Query Bubble -->
              <div v-if="msg.role === 'user'" class="flex items-start justify-end space-x-3 max-w-3xl ml-auto">
                <div class="bg-parchment-primary text-white rounded-3xl rounded-tr-none px-5 py-3.5 text-sm leading-relaxed shadow-sm font-sans">
                  {{ msg.content }}
                </div>
                <div class="w-9 h-9 rounded-2xl bg-parchment-primary-dark text-white flex items-center justify-center font-serif text-xs font-bold flex-shrink-0 shadow-xs border border-white/20">
                  You
                </div>
              </div>

              <!-- Magisterium Assistant Bubble -->
              <div v-else class="flex items-start space-x-3.5 max-w-5xl">
                <div class="w-9 h-9 rounded-2xl bg-parchment-primary-dark text-white flex items-center justify-center font-serif text-sm font-bold flex-shrink-0 shadow-sm border border-amber-300/30">
                  M
                </div>
                <div class="bg-white/90 border border-parchment-border rounded-3xl rounded-tl-none p-5 md:p-6 text-sm text-parchment-neutral shadow-sm space-y-4 w-full backdrop-blur-xs">
                  
                  <!-- Assistant Answer Content -->
                  <div class="whitespace-pre-wrap font-serif text-[15px] leading-relaxed text-stone-800 space-y-2" v-html="formatResponse(msg.content)"></div>

                  <!-- Structured Search Results & Citations (Formatted Schema) -->
                  <div v-if="msg.citations && msg.citations.length > 0" class="mt-5 pt-4 border-t border-amber-200/70 bg-amber-50/50 rounded-2xl p-4 space-y-3">
                    <div class="flex items-center justify-between border-b border-amber-200/60 pb-2.5">
                      <span class="text-xs font-bold text-amber-950 uppercase tracking-wider flex items-center space-x-1.5">
                        <span>📚 Magisterium Citations</span>
                        <span class="bg-amber-200/80 text-amber-900 text-[10px] px-2 py-0.5 rounded-full font-mono font-bold">{{ msg.citations.length }}</span>
                      </span>
                      <span class="text-[10px] text-amber-800/70 font-semibold uppercase tracking-wider">Church Sources</span>
                    </div>

                    <div class="grid grid-cols-1 gap-3">
                      <div 
                        v-for="(cite, idx) in msg.citations" 
                        :key="idx" 
                        class="p-3.5 bg-white border border-amber-200 rounded-2xl shadow-2xs space-y-2 transition-all hover:border-amber-400"
                      >
                        <!-- Header Title & Metadata -->
                        <div class="flex items-start justify-between gap-2">
                          <div class="space-y-0.5 min-w-0 flex-1">
                            <h4 class="text-xs font-bold text-amber-950 font-serif leading-snug">
                              {{ cite.document_title || cite.title || 'Magisterial Citation' }}
                            </h4>
                            <div class="flex flex-wrap items-center gap-x-2.5 text-[11px] text-stone-600 font-serif">
                              <span v-if="cite.author" class="font-semibold text-amber-900">Author: {{ cite.author }}</span>
                              <span v-if="cite.ref" class="bg-amber-100 text-amber-900 px-2 py-0.5 rounded-md text-[10px] font-mono font-semibold">Ref: {{ cite.ref }}</span>
                            </div>
                          </div>

                          <!-- Match Score Badge -->
                          <span v-if="cite.score" class="text-[10px] font-mono font-bold bg-amber-100 text-amber-900 px-2 py-1 rounded-full border border-amber-200 flex-shrink-0">
                            {{ (cite.score * 100).toFixed(1) }}% Match
                          </span>
                        </div>

                        <!-- Snippet Quote -->
                        <p v-if="cite.text" class="text-xs text-stone-700 font-serif leading-relaxed border-l-2 border-amber-500 pl-3 py-1 bg-amber-50/40 rounded-r-xl italic">
                          "{{ cite.text }}"
                        </p>

                        <!-- External Document URL Link -->
                        <div v-if="cite.url" class="pt-1 text-right">
                          <a 
                            :href="cite.url" 
                            target="_blank" 
                            rel="noopener noreferrer" 
                            class="inline-flex items-center space-x-1 text-[11px] font-semibold text-amber-900 hover:text-amber-700 underline transition-colors"
                          >
                            <span>Read Original Document</span>
                            <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
                          </a>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Loading Thinking State -->
            <div v-if="isLoading" class="flex items-start space-x-3.5 max-w-3xl">
              <div class="w-9 h-9 rounded-2xl bg-parchment-primary-dark text-white flex items-center justify-center font-serif text-sm font-bold flex-shrink-0 animate-pulse border border-amber-300/30">
                M
              </div>
              <div class="bg-white/90 border border-parchment-border rounded-3xl rounded-tl-none p-4 text-sm text-parchment-neutral flex items-center space-x-3 shadow-2xs">
                <span class="text-xs font-serif text-stone-600 italic">Searching Magisterium AI database...</span>
                <span class="flex space-x-1">
                  <span class="w-2 h-2 bg-parchment-primary rounded-full animate-bounce"></span>
                  <span class="w-2 h-2 bg-parchment-primary rounded-full animate-bounce [animation-delay:0.2s]"></span>
                  <span class="w-2 h-2 bg-parchment-primary rounded-full animate-bounce [animation-delay:0.4s]"></span>
                </span>
              </div>
            </div>

            <!-- Error Banner -->
            <div v-if="chatError" class="p-4 bg-red-50 border border-red-200 rounded-2xl text-xs text-red-800 flex items-start space-x-3 shadow-2xs">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="flex-1">
                <span class="font-bold">Error Processing Request:</span> {{ chatError }}
                <div v-if="chatError.includes('API Key')" class="mt-1.5">
                  <router-link to="/admin" class="text-amber-900 underline font-semibold hover:text-amber-700">
                    Configure Magisterium API Key in Admin Settings →
                  </router-link>
                </div>
              </div>
            </div>
          </div>

          <!-- Bottom Chat Input Bar -->
          <div class="p-4 bg-parchment-neutral-light/90 border-t border-parchment-border backdrop-blur-xs">
            <form @submit.prevent="sendMessage" class="flex items-center space-x-3">
              <input
                v-model="inputQuery"
                type="text"
                placeholder="Ask a question on Catholic scripture, doctrine, or Church teaching..."
                :disabled="isLoading"
                class="flex-1 px-5 py-3.5 bg-white border border-parchment-border rounded-2xl text-sm text-stone-800 placeholder-parchment-neutral/40 focus:outline-none focus:ring-2 focus:ring-parchment-primary focus:border-transparent transition-all shadow-2xs disabled:opacity-60 font-serif"
              />
              <AppButton
                type="submit"
                variant="primary"
                :disabled="isLoading || !inputQuery.trim()"
                class="py-3.5 px-6 rounded-2xl flex items-center space-x-2 font-semibold shadow-xs"
              >
                <span>Send</span>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="22" y1="2" x2="11" y2="13"></line>
                  <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
                </svg>
              </AppButton>
            </form>
            <div class="mt-2.5 text-center flex items-center justify-between text-[11px] text-parchment-neutral/50 px-2 font-serif">
              <span>Auto-saved to local browser storage</span>
              <router-link to="/admin" class="hover:text-parchment-primary transition-colors flex items-center space-x-1">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
                <span>API Settings</span>
              </router-link>
            </div>
          </div>
        </div>

      </div>
    </main>

    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import AppButton from '../components/common/AppButton.vue';
import { useDevotionApi } from '../composables/useDevotionApi';

interface Citation {
  title?: string;
  document_title?: string;
  author?: string;
  ref?: string;
  score?: number;
  url?: string;
  text?: string;
}

interface ChatMessage {
  role: 'user' | 'assistant';
  content: string;
  citations?: Citation[];
}

interface ChatSession {
  id: string;
  title: string;
  updatedAt: number;
  messages: ChatMessage[];
}

interface ApiUsage {
  limit?: string;
  remaining?: string;
  reset?: string;
}

interface MagisteriumResponse {
  response: string;
  citations?: Citation[];
  usage?: ApiUsage;
}

const { fetchWithAuth } = useDevotionApi();

const LOCAL_STORAGE_KEY = 'magisterium_chat_sessions_v1';

const chatSessions = ref<ChatSession[]>([]);
const activeSessionId = ref<string>('');
const inputQuery = ref('');
const isLoading = ref(false);
const chatError = ref<string | null>(null);
const chatContainer = ref<HTMLDivElement | null>(null);
const apiUsage = ref<ApiUsage | null>(null);

const activeMessages = computed(() => {
  const current = chatSessions.value.find(s => s.id === activeSessionId.value);
  return current ? current.messages : [];
});

onMounted(() => {
  loadSessionsFromLocalStorage();
  const first = chatSessions.value[0];
  if (first) {
    activeSessionId.value = first.id;
  } else {
    createNewSession();
  }
});

const loadSessionsFromLocalStorage = () => {
  try {
    const raw = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (raw) {
      const parsed: ChatSession[] = JSON.parse(raw);
      chatSessions.value = parsed.filter(s => s.messages && s.messages.length > 0);
    }
  } catch (err) {
    console.error('Failed to load chat sessions from localStorage:', err);
  }
};

const saveSessionsToLocalStorage = () => {
  try {
    const sessionsToSave = chatSessions.value.filter(s => s.messages && s.messages.length > 0);
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(sessionsToSave));
  } catch (err) {
    console.error('Failed to save chat sessions to localStorage:', err);
  }
};

const createNewSession = () => {
  const existingEmpty = chatSessions.value.find(s => s.messages.length === 0);
  if (existingEmpty) {
    activeSessionId.value = existingEmpty.id;
    return;
  }

  const newSession: ChatSession = {
    id: `session_${Date.now()}`,
    title: 'New Conversation',
    updatedAt: Date.now(),
    messages: [],
  };
  chatSessions.value.unshift(newSession);
  activeSessionId.value = newSession.id;
};

const selectSession = (id: string) => {
  activeSessionId.value = id;
  scrollToBottom();
};

const deleteSession = (id: string) => {
  chatSessions.value = chatSessions.value.filter(s => s.id !== id);
  if (activeSessionId.value === id) {
    const first = chatSessions.value[0];
    if (first) {
      activeSessionId.value = first.id;
    } else {
      createNewSession();
    }
  }
  saveSessionsToLocalStorage();
};

const formatDate = (timestamp: number) => {
  return new Date(timestamp).toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

const scrollToBottom = async () => {
  await nextTick();
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};

const askSuggested = (question: string) => {
  inputQuery.value = question;
  sendMessage();
};

const formatResponse = (text: string) => {
  if (!text) return '';
  return text
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n/g, '<br/>');
};

const sendMessage = async () => {
  const query = inputQuery.value.trim();
  if (!query || isLoading.value) return;

  const currentSession = chatSessions.value.find(s => s.id === activeSessionId.value);
  if (!currentSession) return;

  chatError.value = null;
  const userMsg: ChatMessage = { role: 'user', content: query };
  currentSession.messages.push(userMsg);

  if (currentSession.messages.length === 1) {
    currentSession.title = query.length > 28 ? query.substring(0, 28) + '...' : query;
  }
  currentSession.updatedAt = Date.now();

  inputQuery.value = '';
  isLoading.value = true;
  saveSessionsToLocalStorage();
  await scrollToBottom();

  try {
    const payloadMessages = currentSession.messages.map(m => ({
      role: m.role,
      content: m.content,
    }));

    const res = await fetchWithAuth<MagisteriumResponse>('/magisterium/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ messages: payloadMessages }),
    });

    if (res && res.response) {
      currentSession.messages.push({
        role: 'assistant',
        content: res.response,
        citations: res.citations || [],
      });
      if (res.usage) {
        apiUsage.value = res.usage;
      }
      currentSession.updatedAt = Date.now();
      saveSessionsToLocalStorage();
    } else {
      throw new Error('Received an empty response from Magisterium AI.');
    }
  } catch (err: any) {
    console.error('Magisterium chat error:', err);
    chatError.value = err.message || 'An error occurred while calling the Magisterium AI endpoint.';
  } finally {
    isLoading.value = false;
    await scrollToBottom();
  }
};
</script>
