<template>
  <div
    class="magisterium-sanctuary w-full flex flex-1 min-h-0 relative overflow-hidden bg-[#F7F3EE] text-[#322D29] font-sans"
    :class="isModal ? 'h-full' : 'h-[calc(100dvh-7rem)] min-h-[560px] rounded-3xl border border-[#D1C7BD]/80 shadow-sm'"
  >
    <!-- Ambient Top Radial Glow (Sacred Gold & Burgundy Warmth) -->
    <div
      class="pointer-events-none absolute inset-0 z-0"
      style="background: radial-gradient(ellipse 68% 46% at 52% 0%, rgba(223, 184, 120, 0.24), rgba(114, 56, 61, 0.06) 46%, transparent 76%);"
    ></div>

    <!-- Copy / Action Toast Banner -->
    <transition name="fade">
      <div
        v-if="copyNotice"
        class="fixed bottom-6 right-6 z-[80] px-4 py-2.5 bg-[#FAF7F2]/95 border border-emerald-400/50 text-emerald-950 font-semibold text-xs rounded-2xl shadow-xl backdrop-blur-xl flex items-center gap-2"
      >
        <span class="w-2 h-2 rounded-full bg-emerald-600"></span>
        <span>{{ copyNotice }}</span>
      </div>
    </transition>

    <!-- Slim Left Icon Rail (Desktop) -->
    <aside
      class="hidden lg:flex flex-col items-center justify-between w-[60px] shrink-0 py-4 border-r border-[#D1C7BD]/70 bg-[#FAF7F2]/85 backdrop-blur-xl z-30"
    >
      <!-- Top Brand Emblem -->
      <button
        type="button"
        @click="createNewSession(); isSidebarDrawerOpen = false"
        title="New Magisterium Chat"
        class="sanctuary-icon-btn w-10 h-10 rounded-xl flex items-center justify-center hover:scale-105 transition-transform cursor-pointer"
      >
        <MagisteriumAvatarIcon class="w-9 h-9 shrink-0" />
      </button>

      <!-- Middle Rail Navigation Icons -->
      <nav class="flex flex-col items-center gap-2.5" aria-label="Magisterium Navigation">
        <!-- New Chat -->
        <button
          type="button"
          @click="createNewSession(); isSidebarDrawerOpen = false"
          title="New Chat"
          class="sanctuary-icon-btn w-10 h-10 rounded-xl flex items-center justify-center text-[#322D29]/70 hover:text-[#72383D] hover:bg-[#EFE9E1] transition-all cursor-pointer"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-[18px] h-[18px] shrink-0">
            <path d="M12 5v14M5 12h14"/>
          </svg>
        </button>

        <!-- My Chats Drawer Toggle -->
        <button
          type="button"
          @click="openDrawerTab('local')"
          title="Saved Conversations"
          class="sanctuary-icon-btn relative w-10 h-10 rounded-xl flex items-center justify-center transition-all cursor-pointer"
          :class="isSidebarDrawerOpen && sidebarTab === 'local'
            ? 'bg-[#72383D]/12 text-[#72383D] border border-[#72383D]/25'
            : 'text-[#322D29]/70 hover:text-[#72383D] hover:bg-[#EFE9E1]'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-[18px] h-[18px] shrink-0">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
          <span
            v-if="chatSessions.length > 0"
            class="absolute top-1.5 right-1.5 w-1.5 h-1.5 rounded-full bg-[#B5935A]"
          ></span>
        </button>

        <!-- Public Sanctuary Drawer Toggle -->
        <button
          type="button"
          @click="openDrawerTab('public')"
          title="Public Sanctuary"
          class="sanctuary-icon-btn w-10 h-10 rounded-xl flex items-center justify-center transition-all cursor-pointer"
          :class="isSidebarDrawerOpen && sidebarTab === 'public'
            ? 'bg-[#72383D]/12 text-[#72383D] border border-[#72383D]/25'
            : 'text-[#322D29]/70 hover:text-[#72383D] hover:bg-[#EFE9E1]'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-[18px] h-[18px] shrink-0">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="2" y1="12" x2="22" y2="12"></line>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
          </svg>
        </button>

        <!-- Magisterium Live Voice Toggle -->
        <button
          type="button"
          @click="toggleLiveSession"
          :disabled="isLiveConnecting"
          title="Magisterium Live Voice"
          class="sanctuary-icon-btn w-10 h-10 rounded-xl flex items-center justify-center transition-all cursor-pointer disabled:opacity-40"
          :class="isLiveConnected
            ? 'bg-red-500/15 text-red-700 border border-red-400/40 animate-pulse'
            : 'text-[#322D29]/70 hover:text-[#72383D] hover:bg-[#EFE9E1]'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-[18px] h-[18px] shrink-0">
            <path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3Z"></path>
            <path d="M19 10v2a7 7 0 0 1-14 0v-2"></path>
            <line x1="12" y1="19" x2="12" y2="22"></line>
          </svg>
        </button>
      </nav>

      <!-- Bottom Close / Exit -->
      <div class="flex flex-col items-center gap-2">
        <button
          v-if="isModal"
          type="button"
          @click="emit('close')"
          title="Close Sanctuary"
          class="sanctuary-icon-btn w-9 h-9 rounded-full flex items-center justify-center bg-[#EFE9E1] hover:bg-[#72383D] text-[#322D29] hover:text-white border border-[#D1C7BD] shadow-2xs transition-all cursor-pointer"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </aside>

    <!-- Slide-Over Sessions & Public Sanctuary Drawer (Desktop & Mobile) -->
    <transition
      enter-active-class="transition duration-250 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isSidebarDrawerOpen || isMobileMenuOpen"
        class="fixed inset-0 z-[70] flex justify-start"
        @click="closeAllDrawers"
      >
        <div class="absolute inset-0 bg-[#322D29]/35 backdrop-blur-xs"></div>

        <div
          class="relative w-[84vw] max-w-[340px] lg:ml-[60px] bg-[#FAF7F2]/95 backdrop-blur-2xl h-full shadow-2xl p-4 flex flex-col z-10 border-r border-[#D1C7BD]"
          @click.stop
        >
          <div class="flex items-center justify-between pb-3.5 mb-3 border-b border-[#D1C7BD]/70">
            <div class="flex items-center gap-2.5">
              <MagisteriumAvatarIcon class="w-7 h-7 shrink-0" />
              <h3 class="font-serif text-sm font-bold text-[#4A2525] tracking-tight">Magisterium History</h3>
            </div>
            <button
              type="button"
              @click="closeAllDrawers"
              class="sanctuary-icon-btn w-8 h-8 text-[#322D29]/60 hover:text-[#72383D] hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer flex items-center justify-center"
              aria-label="Close history drawer"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <!-- Drawer Tabs -->
          <div class="grid grid-cols-2 gap-1 p-1 bg-white/80 border border-[#D1C7BD]/70 rounded-xl mb-3 text-xs font-serif font-bold">
            <button
              type="button"
              @click="sidebarTab = 'local'"
              class="sanctuary-tab-btn py-1.5 rounded-lg transition-all cursor-pointer"
              :class="sidebarTab === 'local' ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-2xs' : 'text-[#322D29]/70 hover:text-[#72383D]'"
            >
              My Chats
            </button>
            <button
              type="button"
              @click="sidebarTab = 'public'"
              class="sanctuary-tab-btn py-1.5 rounded-lg transition-all cursor-pointer"
              :class="sidebarTab === 'public' ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white shadow-2xs' : 'text-[#322D29]/70 hover:text-[#72383D]'"
            >
              Public Sanctuary
            </button>
          </div>

          <!-- Drawer Subheader -->
          <div v-if="sidebarTab === 'local'" class="flex items-center justify-between pb-2.5 mb-2.5 border-b border-[#D1C7BD]/60">
            <span class="font-serif text-[11px] uppercase tracking-wider text-[#72383D] font-bold">Saved Conversations</span>
            <button
              type="button"
              @click="createNewSession(); closeAllDrawers()"
              class="sanctuary-pill-btn text-xs bg-gradient-to-r from-[#72383D] to-[#322D29] text-white px-3 py-1 rounded-lg font-semibold hover:from-[#8B464C] hover:to-[#453E38] transition-all cursor-pointer"
            >
              + New Chat
            </button>
          </div>
          <div v-else class="flex items-center justify-between pb-2.5 mb-2.5 border-b border-[#D1C7BD]/60">
            <span class="font-serif text-[11px] uppercase tracking-wider text-[#322D29] font-bold">Community Discussions</span>
            <button
              type="button"
              @click="fetchPublicConversations"
              class="sanctuary-link-btn text-[11px] text-[#72383D] hover:text-[#322D29] font-semibold underline cursor-pointer"
            >
              Refresh
            </button>
          </div>

          <!-- Local Sessions List -->
          <div v-if="sidebarTab === 'local'" class="overflow-y-auto space-y-2 flex-1 scrollbar-thin pr-1">
            <div
              v-for="session in chatSessions"
              :key="session.id"
              @click="selectSession(session.id); closeAllDrawers()"
              class="p-3 rounded-2xl border text-left cursor-pointer transition-all flex items-center justify-between group"
              :class="activeSessionId === session.id
                ? 'bg-[#EFE9E1] border-[#72383D]/35 text-[#322D29] shadow-2xs'
                : 'bg-white/70 border-[#D1C7BD]/60 text-[#322D29]/85 hover:bg-white hover:border-[#D1C7BD]'"
            >
              <div class="min-w-0 flex-1 pr-2">
                <div class="text-xs font-bold font-serif truncate">
                  {{ session.title || 'New Conversation' }}
                </div>
                <div class="text-[10px] text-[#322D29]/55 mt-1 flex items-center gap-1.5">
                  <span>{{ formatDate(session.updatedAt) }}</span>
                  <span>&middot;</span>
                  <span>{{ session.messages.length }} msgs</span>
                </div>
              </div>

              <button
                type="button"
                @click.stop="deleteSession(session.id)"
                class="sanctuary-icon-btn opacity-0 group-hover:opacity-100 text-stone-400 hover:text-red-600 transition-all w-7 h-7 rounded-lg hover:bg-red-50 cursor-pointer flex items-center justify-center shrink-0"
                title="Delete Chat"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0">
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                </svg>
              </button>
            </div>

            <div v-if="chatSessions.length === 0" class="text-center py-10 px-4 text-xs text-[#322D29]/50 font-serif">
              <p>No saved conversations yet. Start a new chat to begin.</p>
            </div>
          </div>

          <!-- Public Conversations List -->
          <div v-else class="overflow-y-auto space-y-2 flex-1 scrollbar-thin pr-1">
            <div
              v-for="pub in publicConversations"
              :key="pub.id"
              @click="loadPublicSession(pub); closeAllDrawers()"
              class="p-3 rounded-2xl border text-left cursor-pointer transition-all flex items-center justify-between group"
              :class="activeSessionId === pub.id
                ? 'bg-[#EFE9E1] border-[#72383D] text-[#322D29] shadow-2xs'
                : 'bg-white/70 border-[#D1C7BD]/60 text-[#322D29]/85 hover:bg-white hover:border-[#D1C7BD]'"
            >
              <div class="min-w-0 flex-1 pr-2">
                <div class="text-xs font-bold font-serif truncate">
                  {{ pub.title || 'Public Discussion' }}
                </div>
                <div class="text-[10px] text-[#72383D] mt-1 flex items-center gap-1.5 font-serif">
                  <span class="font-semibold">By {{ pub.author_name }}</span>
                  <span>&middot;</span>
                  <span>{{ formatIsoDate(pub.created_at) }}</span>
                </div>
              </div>
              <div class="flex items-center gap-1 shrink-0">
                <button
                  type="button"
                  @click.stop="copyPublicShareLink(pub.id)"
                  class="sanctuary-pill-btn text-[10px] bg-[#EFE9E1] hover:bg-[#D1C7BD] text-[#72383D] font-bold px-2 py-1 rounded-lg border border-[#D1C7BD] transition-colors cursor-pointer"
                  title="Copy Public Share Link"
                >
                  Share
                </button>
                <button
                  type="button"
                  @click.stop="deletePublicConversation(pub.id)"
                  :disabled="isDeletingPublic === pub.id"
                  class="sanctuary-icon-btn opacity-0 group-hover:opacity-100 text-stone-400 hover:text-red-600 transition-all w-7 h-7 rounded-lg hover:bg-red-50 cursor-pointer disabled:opacity-40 flex items-center justify-center shrink-0"
                  title="Delete Public Conversation"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </button>
              </div>
            </div>

            <div v-if="publicConversations.length === 0" class="text-center py-10 px-4 text-xs text-[#322D29]/50 font-serif">
              <p>No public community chats published yet.</p>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- Main Centered Chat Canvas -->
    <div class="flex-1 flex flex-col min-w-0 min-h-0 relative z-10">
      <!-- Magisterium Live Synthesizer Workstation -->
      <MagisteriumSynthesizer
        v-if="isLiveConnected || isLiveConnecting"
        :is-live-connected="isLiveConnected"
        :is-connecting="isLiveConnecting"
        :is-muted="isLiveMuted"
        :is-gemini-speaking="isGeminiSpeaking"
        :is-user-speaking="isUserSpeaking"
        :is-querying-r-a-g="isQueryingRAG"
        :current-rag-query="currentRagQuery"
        :current-citations="liveCitations"
        :live-transcript="liveTranscript"
        :user-live-transcript="userLiveTranscript"
        :live-status="liveStatus"
        :live-error="liveError"
        :audio-level="audioLevel"
        :session-messages="activeSession?.messages || []"
        :active-session-title="activeSession?.title"
        @toggle-mute="toggleLiveMute"
        @stop-live="stopLiveSession"
        @send-text="sendLiveTextMessage"
        @interrupt="interruptLivePlayback"
      />

      <!-- Standard Chat Canvas (When Live Mode is Inactive) -->
      <template v-else>
        <!-- Minimal Floating Top Bar -->
        <header class="h-14 px-3 sm:px-6 flex items-center justify-between gap-3 shrink-0 z-20">
          <!-- Left: Mobile Drawer Trigger + Active Session Title -->
          <div class="flex items-center gap-2.5 min-w-0">
            <button
              type="button"
              @click="isMobileMenuOpen = true"
              class="sanctuary-icon-btn lg:hidden w-9 h-9 rounded-xl bg-white/80 hover:bg-[#EFE9E1] border border-[#D1C7BD] flex items-center justify-center text-[#72383D] cursor-pointer shrink-0 shadow-2xs"
              aria-label="Open chat history"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0">
                <line x1="3" y1="12" x2="21" y2="12"></line>
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <line x1="3" y1="18" x2="21" y2="18"></line>
              </svg>
            </button>

            <div v-if="activeMessages.length > 0" class="flex items-center gap-2 min-w-0">
              <span class="font-serif text-xs sm:text-sm font-semibold text-[#4A2525] truncate">
                {{ activeSession?.title || 'Magisterium Discussion' }}
              </span>
              <span
                v-if="isCurrentSessionPublic"
                class="shrink-0 px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider rounded-full bg-[#72383D]/10 text-[#72383D] border border-[#72383D]/20"
              >
                Public
              </span>
            </div>
          </div>

          <!-- Right: Compact Status & Action Pills -->
          <div class="flex items-center gap-2 shrink-0">
            <span
              v-if="!isAdmin && hasAccess"
              class="hidden sm:inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-medium bg-white/85 border border-[#D1C7BD] text-[#5C282D] shadow-2xs"
              title="Daily guest request limit"
            >
              <span class="w-1.5 h-1.5 rounded-full" :class="remainingToday > 0 ? 'bg-emerald-600' : 'bg-[#9B2C2C]'"></span>
              {{ remainingToday }}/{{ dailyLimit }} left
            </span>
            <button
              v-if="isCurrentSessionPublic"
              type="button"
              @click="deletePublicConversation(activeSessionId)"
              :disabled="isDeletingPublic === activeSessionId"
              class="sanctuary-pill-btn px-3 py-1.5 rounded-full text-xs font-serif font-semibold text-red-700 hover:bg-red-50 bg-white/85 border border-red-200 transition-all cursor-pointer flex items-center gap-1.5 disabled:opacity-50"
              title="Delete this discussion from Public Sanctuary"
            >
              <span class="hidden sm:inline">Delete Public</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
            </button>

            <button
              v-if="activeMessages.length > 0"
              type="button"
              @click="openPublishModal"
              :disabled="isLoading || isStreaming"
              class="sanctuary-pill-btn px-3 py-1.5 rounded-full text-xs font-serif font-semibold text-[#322D29] hover:text-[#72383D] bg-white/85 hover:bg-[#EFE9E1] border border-[#D1C7BD] transition-all cursor-pointer flex items-center gap-1.5 disabled:opacity-40 shadow-2xs"
              title="Publish Conversation to Public Sanctuary"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="2" y1="12" x2="22" y2="12"></line>
                <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
              </svg>
              <span class="hidden sm:inline">Publish</span>
            </button>

            <button
              type="button"
              @click="toggleLiveSession"
              :disabled="isLiveConnecting"
              class="sanctuary-pill-btn px-3.5 py-1.5 rounded-full text-xs font-serif font-bold text-[#72383D] bg-white/85 hover:bg-[#EFE9E1] border border-[#D1C7BD] transition-all cursor-pointer flex items-center gap-1.5 shadow-2xs"
              title="Start Gemini Live Voice Session"
            >
              <span class="w-2 h-2 rounded-full bg-[#C89B54] shrink-0"></span>
              <span>Magisterium Live</span>
            </button>

            <button
              v-if="isModal"
              type="button"
              @click="emit('close')"
              aria-label="Close Magisterium AI modal"
              title="Close"
              class="sanctuary-icon-btn lg:hidden w-8 h-8 rounded-full flex items-center justify-center bg-[#EFE9E1] hover:bg-[#72383D] text-[#322D29] hover:text-white border border-[#D1C7BD] transition-all cursor-pointer shrink-0"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
        </header>

        <!-- Scrollable Chat Stream -->
        <div ref="chatContainer" class="flex-1 overflow-y-auto scrollbar-thin px-4 sm:px-6">
          <!-- Welcome / Empty State (Reference-inspired centered composition in LUXE Parchment) -->
          <div
            v-if="activeMessages.length === 0"
            class="min-h-full flex flex-col items-center justify-center max-w-[740px] mx-auto py-8 text-center animate-fade-in"
          >
            <!-- Glowing Emblem -->
            <div class="relative mb-5">
              <div class="absolute inset-0 rounded-full bg-[#DFB878]/35 blur-2xl scale-150 pointer-events-none"></div>
              <div class="relative w-16 h-16 rounded-2xl bg-[#FAF7F2] border border-[#D1C7BD] flex items-center justify-center shadow-[0_12px_30px_-8px_rgba(114,56,61,0.22)]">
                <MagisteriumAvatarIcon class="w-12 h-12 shrink-0" />
              </div>
            </div>

            <p class="text-xs font-serif font-semibold tracking-wide text-[#72383D]/80 mb-1.5">
              Welcome to Magisterium AI
            </p>
            <h1 class="font-serif text-3xl sm:text-4xl font-semibold tracking-tight text-[#4A2525] mb-2">
              How can I help?
            </h1>
            <p class="text-xs sm:text-sm text-[#322D29]/70 font-serif max-w-md mx-auto leading-relaxed">
              Grounded in the Catechism, Ecumenical Councils, Papal Encyclicals, and Sacred Scripture.
            </p>

            <!-- 3 Horizontal Suggestion Cards -->
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 w-full mt-8">
              <button
                type="button"
                @click="askSuggested('What does the Catechism teach about prayer?')"
                class="sanctuary-card-btn group relative overflow-hidden text-left p-4 rounded-2xl border border-[#D1C7BD]/80 bg-gradient-to-br from-[#FAF7F2] via-[#F4ECE2] to-[#ECE0D1] hover:border-[#72383D]/45 hover:-translate-y-0.5 transition-all duration-200 cursor-pointer min-h-[104px] flex flex-col justify-between shadow-xs"
              >
                <span class="font-serif text-xs font-semibold text-[#322D29] leading-snug pr-3">
                  What does the Catechism teach about prayer?
                </span>
                <span class="text-[11px] text-[#72383D] font-serif font-semibold flex items-center gap-1 opacity-85 group-hover:opacity-100">
                  <span>Catechism</span>
                  <span class="transition-transform group-hover:translate-x-0.5">&rarr;</span>
                </span>
              </button>

              <button
                type="button"
                @click="askSuggested('Explain the relationship between Faith and Reason (Fides et Ratio).')"
                class="sanctuary-card-btn group relative overflow-hidden text-left p-4 rounded-2xl border border-[#D1C7BD]/80 bg-gradient-to-br from-[#FAF7F2] via-[#F3EAE8] to-[#EAE0DE] hover:border-[#72383D]/45 hover:-translate-y-0.5 transition-all duration-200 cursor-pointer min-h-[104px] flex flex-col justify-between shadow-xs"
              >
                <span class="font-serif text-xs font-semibold text-[#322D29] leading-snug pr-3">
                  Explain Faith and Reason (Fides et Ratio)
                </span>
                <span class="text-[11px] text-[#72383D] font-serif font-semibold flex items-center gap-1 opacity-85 group-hover:opacity-100">
                  <span>Encyclicals</span>
                  <span class="transition-transform group-hover:translate-x-0.5">&rarr;</span>
                </span>
              </button>

              <button
                type="button"
                @click="askSuggested('What are the corporal and spiritual works of mercy?')"
                class="sanctuary-card-btn group relative overflow-hidden text-left p-4 rounded-2xl border border-[#D1C7BD]/80 bg-gradient-to-br from-[#FAF7F2] via-[#F6EFE3] to-[#EFE2CE] hover:border-[#72383D]/45 hover:-translate-y-0.5 transition-all duration-200 cursor-pointer min-h-[104px] flex flex-col justify-between shadow-xs"
              >
                <span class="font-serif text-xs font-semibold text-[#322D29] leading-snug pr-3">
                  What are the corporal &amp; spiritual works of mercy?
                </span>
                <span class="text-[11px] text-[#8B6B35] font-serif font-semibold flex items-center gap-1 opacity-85 group-hover:opacity-100">
                  <span>Tradition</span>
                  <span class="transition-transform group-hover:translate-x-0.5">&rarr;</span>
                </span>
              </button>
            </div>
          </div>

          <!-- Active Messages Stream -->
          <div v-else class="max-w-[740px] mx-auto py-6 space-y-8">
            <div
              v-for="(msg, index) in activeMessages"
              :key="index"
              :ref="el => { if (el) messageRefs[index] = el as HTMLElement }"
              class="animate-fade-in"
            >
              <!-- User Prompt (Centered Editorial Header Style) -->
              <div v-if="msg.role === 'user'" class="text-center space-y-2 py-2">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-[#72383D]/10 text-[#72383D] border border-[#72383D]/20">
                  You
                </span>
                <h2 class="font-serif text-lg sm:text-xl font-semibold text-[#4A2525] leading-snug max-w-2xl mx-auto">
                  {{ msg.content }}
                </h2>
              </div>

              <!-- Assistant Response Card -->
              <div v-else class="space-y-3">
                <div class="bg-white/80 border border-[#D1C7BD]/80 rounded-3xl p-5 sm:p-7 backdrop-blur-md shadow-xs">
                  <!-- Prose Answer -->
                  <div class="text-[15px] leading-relaxed text-[#322D29]/95 font-serif">
                    <span v-html="formatResponse(msg.content)"></span>
                    <span
                      v-if="isStreaming && activeMessages.indexOf(msg) === activeMessages.length - 1"
                      class="inline-block w-2 h-4 ml-1 bg-[#72383D] animate-pulse align-middle rounded-xs"
                    ></span>
                  </div>

                  <!-- Collapsible Citations Drawer -->
                  <div v-if="msg.citations && msg.citations.length > 0 && expandedCitations.has(index)" class="mt-6 pt-5 border-t border-[#D1C7BD]/60 space-y-3 animate-fade-in">
                    <div class="flex items-center justify-between">
                      <span class="text-[11px] font-bold uppercase tracking-wider text-[#72383D]">
                        Magisterium Sources ({{ msg.citations.length }})
                      </span>
                      <button
                        type="button"
                        @click="toggleCitations(index)"
                        class="sanctuary-link-btn text-[11px] text-[#322D29]/60 hover:text-[#72383D] font-semibold cursor-pointer"
                      >
                        Hide
                      </button>
                    </div>

                    <div class="grid grid-cols-1 gap-2.5">
                      <div
                        v-for="(cite, idx) in msg.citations"
                        :key="idx"
                        class="p-3.5 bg-[#FAF7F2] border border-[#D1C7BD]/70 rounded-2xl space-y-2 hover:border-[#72383D]/40 transition-colors"
                      >
                        <div class="flex items-start justify-between gap-2">
                          <div class="space-y-0.5 min-w-0 flex-1">
                            <h4 class="text-xs font-bold text-[#72383D] font-serif leading-snug">
                              {{ cite.document_title || cite.title || 'Magisterial Citation' }}
                            </h4>
                            <div class="flex flex-wrap items-center gap-2 text-[11px] text-[#322D29]/70 font-serif">
                              <span v-if="cite.author" class="text-[#72383D] font-semibold">{{ cite.author }}</span>
                              <span v-if="cite.ref" class="bg-[#EFE9E1] text-[#72383D] px-2 py-0.5 rounded-md text-[10px] font-mono font-semibold">
                                {{ cite.ref }}
                              </span>
                            </div>
                          </div>
                          <span
                            v-if="cite.score"
                            class="text-[10px] font-mono font-bold bg-[#EFE9E1] text-[#72383D] px-2 py-0.5 rounded-full border border-[#D1C7BD] shrink-0"
                          >
                            {{ (cite.score * 100).toFixed(0) }}%
                          </span>
                        </div>

                        <p v-if="cite.text" class="text-xs text-[#322D29]/80 font-serif leading-relaxed border-l-2 border-[#72383D] pl-3 py-0.5 italic">
                          "{{ cite.text }}"
                        </p>

                        <div v-if="cite.url" class="text-right">
                          <a
                            :href="cite.url"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="inline-flex items-center gap-1 text-[11px] font-semibold text-[#72383D] hover:text-[#322D29] underline"
                          >
                            <span>Read Original Document</span>
                            <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3 h-3 shrink-0"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
                          </a>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Floating Action Toolbar Pill (Reference-style) -->
                <div class="flex items-center justify-center">
                  <div class="inline-flex items-center gap-1 p-1 rounded-xl bg-[#FAF7F2]/95 border border-[#D1C7BD] shadow-sm backdrop-blur-md">
                    <button
                      type="button"
                      @click="copyAssistantResponse(msg.content)"
                      title="Copy response"
                      class="sanctuary-icon-btn w-7 h-7 rounded-lg text-[#322D29]/65 hover:text-[#72383D] hover:bg-[#EFE9E1] transition-colors cursor-pointer flex items-center justify-center"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0">
                        <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                      </svg>
                    </button>

                    <button
                      v-if="!isLoading && !isStreaming"
                      type="button"
                      @click="regenerateFromIndex(index)"
                      title="Ask again"
                      class="sanctuary-icon-btn w-7 h-7 rounded-lg text-[#322D29]/65 hover:text-[#72383D] hover:bg-[#EFE9E1] transition-colors cursor-pointer flex items-center justify-center"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0">
                        <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
                        <path d="M3 3v5h5"></path>
                        <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"></path>
                        <path d="M16 16h5v5"></path>
                      </svg>
                    </button>

                    <button
                      v-if="msg.citations && msg.citations.length > 0"
                      type="button"
                      @click="toggleCitations(index)"
                      class="sanctuary-pill-btn px-2.5 py-1 rounded-lg text-[11px] font-serif font-bold transition-colors cursor-pointer flex items-center gap-1.5"
                      :class="expandedCitations.has(index)
                        ? 'bg-[#72383D] text-white'
                        : 'text-[#72383D] hover:bg-[#EFE9E1]'"
                    >
                      <span>Sources</span>
                      <span
                        class="px-1.5 rounded-full text-[10px] font-mono"
                        :class="expandedCitations.has(index) ? 'bg-white/20 text-white' : 'bg-[#EFE9E1] text-[#72383D]'"
                      >{{ msg.citations.length }}</span>
                    </button>

                    <button
                      type="button"
                      @click="openPublishModal"
                      title="Publish to Public Sanctuary"
                      class="sanctuary-icon-btn w-7 h-7 rounded-lg text-[#322D29]/65 hover:text-[#72383D] hover:bg-[#EFE9E1] transition-colors cursor-pointer flex items-center justify-center"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 shrink-0">
                        <path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path>
                        <polyline points="16 6 12 2 8 6"></polyline>
                        <line x1="12" y1="2" x2="12" y2="15"></line>
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Loading / Generating State -->
            <div v-if="isLoading" class="bg-white/80 border border-[#D1C7BD] rounded-3xl p-6 sm:p-8 flex flex-col items-center justify-center gap-3 text-center shadow-2xs">
              <div class="flex items-center gap-1.5">
                <span class="w-2 h-2 bg-[#72383D] rounded-full animate-bounce"></span>
                <span class="w-2 h-2 bg-[#72383D] rounded-full animate-bounce [animation-delay:0.15s]"></span>
                <span class="w-2 h-2 bg-[#72383D] rounded-full animate-bounce [animation-delay:0.3s]"></span>
              </div>
              <span class="text-xs text-[#322D29]/75 font-serif italic">
                {{ streamStatus || 'Searching Magisterium AI database...' }}
              </span>
            </div>

            <!-- Error Banner -->
            <div v-if="chatError" class="p-4 bg-red-50 border border-red-200 rounded-2xl text-xs text-red-900 flex items-start gap-3 shadow-2xs">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-red-600 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="flex-1">
                <span class="font-bold">Error Processing Request:</span> {{ chatError }}
                <div v-if="chatError.includes('API Key')" class="mt-1.5">
                  <router-link to="/admin" @click="emit('close')" class="text-[#72383D] underline font-semibold hover:text-[#322D29]">
                    Configure Magisterium API Key in Admin Settings &rarr;
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Floating Bottom Composer -->
        <div class="w-full max-w-[740px] mx-auto px-4 sm:px-6 pb-4 sm:pb-6 pt-2 shrink-0 relative z-20">
          <!-- Floating Focus Selector Popover (Above Composer) -->
          <transition
            enter-active-class="transition duration-150 ease-out"
            enter-from-class="opacity-0 translate-y-2 scale-95"
            enter-to-class="opacity-100 translate-y-0 scale-100"
            leave-active-class="transition duration-100 ease-in"
            leave-from-class="opacity-100 translate-y-0 scale-100"
            leave-to-class="opacity-0 translate-y-2 scale-95"
          >
            <div
              v-if="isFocusMenuOpen"
              class="absolute right-4 sm:right-6 bottom-[126px] w-64 p-2 rounded-2xl bg-[#FAF7F2]/95 border border-[#D1C7BD] shadow-xl backdrop-blur-2xl z-30"
            >
              <div class="px-2.5 py-1.5 text-[10px] font-bold uppercase tracking-wider text-[#72383D]/70">
                Source Focus
              </div>
              <button
                v-for="opt in focusOptions"
                :key="opt.id"
                type="button"
                @click="selectedFocus = opt.id; isFocusMenuOpen = false"
                class="sanctuary-menu-btn w-full text-left px-3 py-2 rounded-xl text-xs font-serif flex items-center justify-between transition-colors cursor-pointer"
                :class="selectedFocus === opt.id
                  ? 'bg-[#72383D]/10 text-[#72383D] font-bold'
                  : 'text-[#322D29]/80 hover:bg-[#EFE9E1] hover:text-[#322D29]'"
              >
                <span>{{ opt.label }}</span>
                <span v-if="selectedFocus === opt.id" class="w-2 h-2 rounded-full bg-[#72383D]"></span>
              </button>
            </div>
          </transition>

          <!-- Top-Right Mode / Focus Pill Bar -->
          <div class="flex items-center justify-end gap-2 mb-2">
            <button
              type="button"
              @click="isFocusMenuOpen = !isFocusMenuOpen"
              class="sanctuary-pill-btn inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-serif font-semibold bg-white/90 hover:bg-[#EFE9E1] text-[#322D29] border border-[#D1C7BD] backdrop-blur-xl transition-all cursor-pointer shadow-xs"
            >
              <span class="w-2 h-2 rounded-full bg-[#B5935A] shrink-0"></span>
              <span>{{ activeFocusLabel }}</span>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="w-3 h-3 shrink-0 text-[#72383D] transition-transform duration-150"
                :class="{ 'rotate-180': isFocusMenuOpen }"
              >
                <path d="m18 15-6-6-6 6"/>
              </svg>
            </button>
          </div>

          <!-- Glass Composer Card -->
          <form
            @submit.prevent="sendMessage"
            class="bg-white/90 backdrop-blur-2xl border border-[#D1C7BD] focus-within:border-[#72383D]/60 rounded-[22px] p-3.5 sm:p-4 shadow-[0_16px_40px_-14px_rgba(50,45,41,0.16)] transition-all"
          >
            <textarea
              ref="textareaRef"
              v-model="inputQuery"
              rows="1"
              :placeholder="isLiveConnected ? 'Speak into your microphone, or type here to prompt Gemini Live...' : (isLoading || isStreaming) ? 'Magisterium is responding...' : 'Ask a question on Catholic scripture, doctrine, or Church teaching...'"
              :disabled="isLoading || isStreaming"
              @keydown.enter.exact.prevent="sendMessage"
              @input="adjustTextareaHeight"
              class="w-full bg-transparent text-sm sm:text-[15px] text-[#322D29] placeholder-[#322D29]/40 focus:outline-none disabled:opacity-60 font-serif resize-none min-h-[28px] max-h-36 overflow-y-auto leading-relaxed scrollbar-thin"
            ></textarea>

            <!-- Composer Bottom Action Row -->
            <div class="flex items-center justify-between pt-2.5 mt-1">
              <div class="flex items-center gap-1.5">
                <!-- Mic / Live Toggle -->
                <button
                  type="button"
                  @click="toggleLiveSession"
                  :disabled="isLiveConnecting || isStreaming"
                  class="sanctuary-icon-btn w-8 h-8 rounded-full flex items-center justify-center transition-all cursor-pointer disabled:opacity-40"
                  :class="isLiveConnected
                    ? 'bg-red-500/15 text-red-700 border border-red-400/40 animate-pulse'
                    : 'text-[#72383D]/75 hover:text-[#72383D] hover:bg-[#EFE9E1]'"
                  :title="isLiveConnected ? 'Disconnect Gemini Live' : 'Start Gemini Live Voice Session'"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0">
                    <path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3Z"></path>
                    <path d="M19 10v2a7 7 0 0 1-14 0v-2"></path>
                    <line x1="12" y1="19" x2="12" y2="22"></line>
                  </svg>
                </button>

                <!-- Publish Toggle -->
                <button
                  type="button"
                  @click="openPublishModal"
                  :disabled="activeMessages.length === 0 || isLoading || isStreaming"
                  class="sanctuary-icon-btn w-8 h-8 rounded-full flex items-center justify-center text-[#72383D]/75 hover:text-[#72383D] hover:bg-[#EFE9E1] transition-all cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed"
                  title="Publish Conversation to Public Sanctuary"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="2" y1="12" x2="22" y2="12"></line>
                    <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
                  </svg>
                </button>
              </div>

              <!-- Stop or Send Button -->
              <button
                v-if="isLoading || isStreaming"
                type="button"
                @click="stopGeneration"
                class="sanctuary-pill-btn px-3.5 h-8 rounded-full bg-red-600 hover:bg-red-700 text-white text-xs font-semibold flex items-center gap-1.5 transition-all cursor-pointer active:scale-95"
                title="Stop generating"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="11" height="11" viewBox="0 0 24 24" fill="currentColor" class="w-3 h-3 shrink-0">
                  <rect x="4" y="4" width="16" height="16" rx="2"></rect>
                </svg>
                <span>Stop</span>
              </button>

              <button
                v-else
                type="submit"
                :disabled="!inputQuery.trim()"
                aria-label="Send message"
                class="sanctuary-icon-btn w-9 h-9 rounded-full flex items-center justify-center transition-all cursor-pointer active:scale-95 disabled:cursor-not-allowed"
                :class="inputQuery.trim()
                  ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white hover:from-[#8B464C] hover:to-[#453E38] shadow-sm'
                  : 'bg-[#EFE9E1] text-[#322D29]/35'"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 shrink-0">
                  <line x1="12" y1="19" x2="12" y2="5"></line>
                  <polyline points="5 12 12 5 19 12"></polyline>
                </svg>
              </button>
            </div>
          </form>
        </div>
      </template>
    </div>

    <!-- Publish Conversation to Database Modal (LUXE Sacred Parchment) -->
    <Teleport to="body">
      <div v-if="showPublishModal" class="fixed inset-0 z-[80] bg-[#322D29]/50 backdrop-blur-sm flex items-center justify-center p-4" @click="showPublishModal = false">
        <div class="bg-[#FAF7F2] border border-[#D1C7BD] rounded-3xl p-6 max-w-md w-full shadow-2xl space-y-4 text-[#322D29] font-serif" @click.stop>
          <div class="flex items-center justify-between pb-3 border-b border-[#D1C7BD]/70">
            <h3 class="text-base font-bold text-[#4A2525] tracking-tight">Publish to Public Sanctuary</h3>
            <button
              type="button"
              @click="showPublishModal = false"
              class="sanctuary-icon-btn w-7 h-7 rounded-full bg-[#EFE9E1] hover:bg-[#72383D] text-[#322D29] hover:text-white flex items-center justify-center cursor-pointer transition-colors"
            >
              &times;
            </button>
          </div>

          <div class="bg-[#EFE9E1]/80 border border-[#D1C7BD] rounded-2xl p-3.5 text-xs text-[#322D29]/85 leading-relaxed">
            <span class="font-bold text-[#72383D] block mb-1">Public Community Disclosure</span>
            This discussion will be stored in the database and made publicly visible under the <strong>Public Sanctuary</strong> tab.
          </div>

          <form @submit.prevent="submitPublishConversation" class="space-y-4">
            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-[#72383D] mb-1.5">Your Name / Title</label>
              <input
                v-model="authorName"
                type="text"
                required
                placeholder="e.g. Brother John, Maria, Pilgrim..."
                class="w-full px-3.5 py-2.5 bg-white border border-[#D1C7BD] rounded-xl text-sm text-[#322D29] placeholder-[#322D29]/35 focus:outline-none focus:ring-2 focus:ring-[#72383D] font-sans"
              />
            </div>

            <div v-if="publishMessage" class="p-3 bg-emerald-100 text-emerald-900 border border-emerald-300 rounded-xl text-xs font-semibold">
              {{ publishMessage }}
            </div>

            <div class="flex justify-end gap-2.5 pt-1">
              <AppButton type="button" variant="secondary" @click="showPublishModal = false">
                Cancel
              </AppButton>
              <AppButton type="submit" variant="primary" :disabled="isPublishing || !authorName.trim()">
                {{ isPublishing ? 'Publishing...' : 'Confirm & Publish' }}
              </AppButton>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue';
import AppButton from '../common/AppButton.vue';
import { useDevotionApi } from '../../composables/useDevotionApi';
import { useAccessControl } from '../../composables/useAccessControl';
import { useGeminiLive, type LiveTurnResult } from '../../composables/useGeminiLive';
import { useMagisteriumModal } from '../../composables/useMagisteriumModal';
import MagisteriumSynthesizer from './MagisteriumSynthesizer.vue';
import MagisteriumAvatarIcon from '../icons/MagisteriumAvatarIcon.vue';

const props = withDefaults(defineProps<{
  isModal?: boolean;
}>(), {
  isModal: false,
});

const emit = defineEmits<{
  (e: 'close'): void;
}>();

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

interface PublicConversation {
  id: string;
  author_name: string;
  title: string;
  messages: ChatMessage[];
  created_at: string;
}

const { fetchWithAuth, getApiUrl, getAuthHeaders } = useDevotionApi();
const { sharedSessionId, clearSharedSessionId } = useMagisteriumModal();
const {
  isAdmin,
  hasAccess,
  remainingToday,
  dailyLimit,
  refreshUsage,
  clearAccessCode,
  openAccessModal,
} = useAccessControl();

const LOCAL_STORAGE_KEY = 'magisterium_chat_sessions_v1';

const sidebarTab = ref<'local' | 'public'>('local');
const chatSessions = ref<ChatSession[]>([]);
const publicConversations = ref<PublicConversation[]>([]);
const activeSessionId = ref<string>('');
const isMobileMenuOpen = ref(false);
const isSidebarDrawerOpen = ref(false);
const isFocusMenuOpen = ref(false);

type FocusId = 'all' | 'catechism' | 'scripture' | 'encyclicals' | 'fathers';
const selectedFocus = ref<FocusId>('all');
const focusOptions: { id: FocusId; label: string; hint?: string }[] = [
  { id: 'all', label: 'All Magisterium' },
  { id: 'catechism', label: 'Catechism (CCC)', hint: 'Focus primarily on the Catechism of the Catholic Church.' },
  { id: 'scripture', label: 'Sacred Scripture', hint: 'Focus primarily on Sacred Scripture and biblical exegesis.' },
  { id: 'encyclicals', label: 'Papal Encyclicals', hint: 'Focus primarily on Papal Encyclicals and Apostolic Letters.' },
  { id: 'fathers', label: 'Church Fathers', hint: 'Focus primarily on the writings of the Church Fathers and Doctors.' },
];

const activeFocusLabel = computed(() => {
  return focusOptions.find(o => o.id === selectedFocus.value)?.label || 'All Magisterium';
});

const openDrawerTab = (tab: 'local' | 'public') => {
  if (isSidebarDrawerOpen.value && sidebarTab.value === tab) {
    isSidebarDrawerOpen.value = false;
    return;
  }
  sidebarTab.value = tab;
  isSidebarDrawerOpen.value = true;
};

const closeAllDrawers = () => {
  isSidebarDrawerOpen.value = false;
  isMobileMenuOpen.value = false;
};
const inputQuery = ref('');
const isLoading = ref(false);
const isStreaming = ref(false);
const streamStatus = ref<string>('');
const chatError = ref<string | null>(null);
let abortController: AbortController | null = null;
let scrollRafId: number | null = null;
const chatContainer = ref<HTMLDivElement | null>(null);
const textareaRef = ref<HTMLTextAreaElement | null>(null);
const messageRefs = ref<HTMLElement[]>([]);
const apiUsage = ref<ApiUsage | null>(null);

const adjustTextareaHeight = () => {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto';
    const newHeight = Math.min(textareaRef.value.scrollHeight, 144);
    textareaRef.value.style.height = `${newHeight}px`;
  }
};

watch(inputQuery, (val) => {
  if (!val) {
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto';
    }
  } else {
    nextTick(adjustTextareaHeight);
  }
});

// Toast Banner / Copy state
const copyNotice = ref<string | null>(null);

// Collapsible Citations State (Hidden by default)
const expandedCitations = ref<Set<number>>(new Set());

const toggleCitations = (msgIndex: number) => {
  if (expandedCitations.value.has(msgIndex)) {
    expandedCitations.value.delete(msgIndex);
  } else {
    expandedCitations.value.add(msgIndex);
  }
};

// Modal & Publish State
const showPublishModal = ref(false);
const authorName = ref('');
const isPublishing = ref(false);
const publishMessage = ref<string | null>(null);

const activeSession = computed(() => {
  return chatSessions.value.find(s => s.id === activeSessionId.value) || null;
});

const isCurrentSessionPublic = computed(() => {
  return publicConversations.value.some(p => p.id === activeSessionId.value);
});

const activeMessages = computed(() => {
  const current = activeSession.value;
  return current ? current.messages : [];
});

onMounted(async () => {
  loadSessionsFromLocalStorage();
  await fetchPublicConversations();

  // Check if a shared public conversation ID was requested
  const paramId = sharedSessionId.value;
  if (paramId) {
    clearSharedSessionId();
    const pub = publicConversations.value.find(p => p.id === paramId);
    if (pub) {
      sidebarTab.value = 'public';
      loadPublicSession(pub);
      return;
    }
  }

  // Always default to creating a new chat session when landing on the page
  createNewSession();
});

const copyPublicShareLink = (id: string) => {
  const baseUrl = window.location.origin + window.location.pathname;
  const shareUrl = `${baseUrl}#/magisterium-chat/${id}`;
  
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(shareUrl);
  } else {
    const textArea = document.createElement("textarea");
    textArea.value = shareUrl;
    document.body.appendChild(textArea);
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);
  }

  copyNotice.value = 'Public Share Link copied to clipboard!';
  setTimeout(() => {
    copyNotice.value = null;
  }, 3000);
};

const fetchPublicConversations = async () => {
  try {
    const res = await fetchWithAuth<{ conversations: PublicConversation[] }>('/magisterium/conversations/public');
    if (res && res.conversations) {
      publicConversations.value = res.conversations;
    }
  } catch (err) {
    console.error('Failed to fetch public conversations:', err);
  }
};

const openPublishModal = () => {
  publishMessage.value = null;
  showPublishModal.value = true;
};

const submitPublishConversation = async () => {
  const current = chatSessions.value.find(s => s.id === activeSessionId.value);
  if (!current || current.messages.length === 0 || !authorName.value.trim()) return;

  isPublishing.value = true;
  publishMessage.value = null;

  try {
    const payload = {
      id: current.id,
      author_name: authorName.value.trim(),
      title: current.title || 'Magisterium AI Reflection',
      messages: current.messages,
    };

    await fetchWithAuth('/magisterium/conversations/save', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    publishMessage.value = 'Conversation successfully published! It is now visible to all users under Public Sanctuary.';
    fetchPublicConversations();
    setTimeout(() => {
      showPublishModal.value = false;
      sidebarTab.value = 'public';
    }, 1500);
  } catch (err: any) {
    chatError.value = err.message || 'Failed to publish conversation to database.';
  } finally {
    isPublishing.value = false;
  }
};

const isDeletingPublic = ref<string | null>(null);

const deletePublicConversation = async (id: string) => {
  if (!window.confirm('Are you sure you want to remove this conversation from the Public Sanctuary?')) {
    return;
  }
  isDeletingPublic.value = id;
  try {
    await fetchWithAuth('/magisterium/conversations/delete', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id }),
    });

    // Remove from public conversations
    publicConversations.value = publicConversations.value.filter(p => p.id !== id);

    // If active session is the deleted public session, remove from local sessions list if present
    if (activeSessionId.value === id) {
      chatSessions.value = chatSessions.value.filter(s => s.id !== id);
      saveSessionsToLocalStorage();
      const first = chatSessions.value[0];
      if (first) {
        activeSessionId.value = first.id;
      } else {
        createNewSession();
      }
    }

    copyNotice.value = 'Conversation successfully removed from Public Sanctuary!';
    setTimeout(() => {
      copyNotice.value = null;
    }, 3000);
  } catch (err: any) {
    console.error('Failed to delete public conversation:', err);
    chatError.value = err.message || 'Failed to remove conversation from Public Sanctuary.';
  } finally {
    isDeletingPublic.value = null;
  }
};

const streamScrollToBottom = () => {
  if (scrollRafId !== null) return;
  scrollRafId = requestAnimationFrame(() => {
    scrollRafId = null;
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
    }
  });
};

const stopGeneration = () => {
  if (abortController) {
    abortController.abort();
    abortController = null;
  }
  if (scrollRafId !== null) {
    cancelAnimationFrame(scrollRafId);
    scrollRafId = null;
  }
  isLoading.value = false;
  isStreaming.value = false;
  streamStatus.value = '';

  const current = activeSession.value;
  if (current && current.messages.length > 0) {
    const last = current.messages[current.messages.length - 1];
    if (last && last.role === 'assistant' && !last.content.trim() && (!last.citations || last.citations.length === 0)) {
      current.messages.pop();
    }
  }
  saveSessionsToLocalStorage();
};

const loadPublicSession = (pub: PublicConversation) => {
  stopGeneration();
  // Load public conversation as a readable local chat session
  let existing = chatSessions.value.find(s => s.id === pub.id);
  if (!existing) {
    existing = {
      id: pub.id,
      title: `[Public] ${pub.title} (${pub.author_name})`,
      updatedAt: new Date(pub.created_at).getTime() || Date.now(),
      messages: pub.messages || [],
    };
    chatSessions.value.unshift(existing);
  }
  activeSessionId.value = existing.id;
  scrollToTop();
};

const formatIsoDate = (isoStr: string) => {
  if (!isoStr) return '';
  try {
    return new Date(isoStr).toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
    });
  } catch {
    return isoStr;
  }
};

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
  stopGeneration();
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
  if (activeSessionId.value !== id) {
    stopGeneration();
  }
  activeSessionId.value = id;
  scrollToBottom();
};

const deleteSession = (id: string) => {
  if (activeSessionId.value === id) {
    stopGeneration();
  }
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

const scrollToTop = async () => {
  await nextTick();
  if (chatContainer.value) {
    chatContainer.value.scrollTop = 0;
  }
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

const copyAssistantResponse = async (content: string) => {
  if (!content) return;
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(content);
    } else {
      const textArea = document.createElement('textarea');
      textArea.value = content;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
    }
    copyNotice.value = 'Response copied to clipboard';
    setTimeout(() => {
      copyNotice.value = null;
    }, 2500);
  } catch (err) {
    console.error('Failed to copy response:', err);
  }
};

const regenerateFromIndex = (assistantIndex: number) => {
  const currentSession = chatSessions.value.find(s => s.id === activeSessionId.value);
  if (!currentSession || isLoading.value || isStreaming.value) return;

  let userPrompt = '';
  for (let i = assistantIndex - 1; i >= 0; i--) {
    const msg = currentSession.messages[i];
    if (msg && msg.role === 'user') {
      userPrompt = msg.content;
      break;
    }
  }
  if (!userPrompt) return;

  inputQuery.value = userPrompt;
  sendMessage();
};

const formatResponse = (text: string) => {
  if (!text) return '';
  
  let html = text;

  // Escape unsafe HTML entities first
  html = html
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');

  // Fenced Code Blocks ```lang ... ```
  html = html.replace(/```([\s\S]*?)```/g, (_match, code) => {
    return `<pre class="bg-stone-900 text-stone-100 p-3 rounded-xl my-2 overflow-x-auto font-mono text-xs"><code>${code.trim()}</code></pre>`;
  });

  // Inline code `code`
  html = html.replace(/`([^`]+)`/g, '<code class="bg-amber-100/80 text-amber-950 px-1.5 py-0.5 rounded font-mono text-xs">$1</code>');

  // Headings
  html = html.replace(/^### (.*$)/gim, '<h3 class="text-base font-bold font-serif text-amber-950 mt-4 mb-1">$1</h3>');
  html = html.replace(/^## (.*$)/gim, '<h2 class="text-lg font-bold font-serif text-amber-950 mt-4 mb-1.5 pb-1 border-b border-amber-200">$1</h2>');
  html = html.replace(/^# (.*$)/gim, '<h1 class="text-xl font-bold font-serif text-amber-950 mt-5 mb-2 pb-1 border-b border-amber-300">$1</h1>');

  // Blockquotes
  html = html.replace(/^\&gt; (.*$)/gim, '<blockquote class="border-l-4 border-amber-500 pl-3 py-1 my-2 text-stone-700 bg-amber-50/50 rounded-r-lg">$1</blockquote>');

  // Bold & Italics
  html = html.replace(/\*\*\*(.*?)\*\*\*/g, '<strong><em>$1</em></strong>');
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong class="font-bold text-amber-950">$1</strong>');
  html = html.replace(/\*(.*?)\*/g, '<em class="italic">$1</em>');

  // Unordered Lists (- or *)
  html = html.replace(/^[\*\-] (.*$)/gim, '<li class="ml-4 list-disc text-stone-800 my-0.5">$1</li>');

  // Ordered Lists (1. 2.)
  html = html.replace(/^\d+\. (.*$)/gim, '<li class="ml-4 list-decimal text-stone-800 my-0.5">$1</li>');

  // Wrap contiguous <li> tags in <ul> or <ol>
  html = html.replace(/((?:<li class="ml-4 list-disc[^>]*>.*?<\/li>\s*)+)/gs, '<ul class="my-2 space-y-1">$1</ul>');
  html = html.replace(/((?:<li class="ml-4 list-decimal[^>]*>.*?<\/li>\s*)+)/gs, '<ol class="my-2 space-y-1">$1</ol>');

  // Markdown Links [title](url)
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-amber-900 font-semibold underline hover:text-amber-700">$1</a>');

  // Paragraph breaks & newlines
  html = html.replace(/\n\n+/g, '</p><p class="my-2.5">');
  html = html.replace(/\n/g, '<br/>');

  return `<p class="my-1">${html}</p>`;
};

const onLiveTurnComplete = (result: LiveTurnResult) => {
  const currentSession = chatSessions.value.find(s => s.id === activeSessionId.value);
  if (!currentSession) return;

  if (result.userQuery) {
    currentSession.messages.push({
      role: 'user',
      content: result.userQuery,
    });
  }

  currentSession.messages.push({
    role: 'assistant',
    content: result.response,
    citations: result.citations,
  });

  if (currentSession.messages.length <= 2 && result.userQuery) {
    currentSession.title = result.userQuery.length > 28 ? result.userQuery.substring(0, 28) + '...' : result.userQuery;
  }
  currentSession.updatedAt = Date.now();
  saveSessionsToLocalStorage();
  scrollToBottom();
};

const {
  isLiveConnected,
  isConnecting: isLiveConnecting,
  isMuted: isLiveMuted,
  isGeminiSpeaking,
  isUserSpeaking,
  isQueryingRAG,
  currentRagQuery,
  currentCitations: liveCitations,
  liveTranscript,
  userLiveTranscript,
  liveStatus,
  liveError,
  audioLevel,
  startLiveSession,
  stopLiveSession,
  toggleMute: toggleLiveMute,
  sendTextMessage: sendLiveTextMessage,
  interruptPlayback: interruptLivePlayback,
} = useGeminiLive(onLiveTurnComplete);

const toggleLiveSession = async () => {
  if (isLiveConnected.value || isLiveConnecting.value) {
    stopLiveSession();
  } else {
    chatError.value = null;
    await startLiveSession();
  }
};

watch(liveError, (val) => {
  if (val) {
    chatError.value = val;
  }
});

onUnmounted(() => {
  stopGeneration();
  stopLiveSession();
});

const sendMessage = async () => {
  const query = inputQuery.value.trim();
  if (!query || isLoading.value || isStreaming.value) return;

  if (isLiveConnected.value) {
    sendLiveTextMessage(query);
    inputQuery.value = '';
    return;
  }

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
  isStreaming.value = false;
  streamStatus.value = 'Searching Magisterium AI database...';
  saveSessionsToLocalStorage();
  await scrollToBottom();

  abortController = new AbortController();
  let assistantMsg: ChatMessage | null = null;

  try {
    const focusOption = focusOptions.find(o => o.id === selectedFocus.value);
    const payloadMessages = currentSession.messages.map((m, idx) => {
      const isLastUser = idx === currentSession.messages.length - 1 && m.role === 'user';
      if (isLastUser && focusOption?.hint) {
        return {
          role: m.role,
          content: `${m.content}\n\n[${focusOption.hint}]`,
        };
      }
      return {
        role: m.role,
        content: m.content,
      };
    });

    const url = getApiUrl('/magisterium/chat');
    const headers = getAuthHeaders({
      'Content-Type': 'application/json',
      'Accept': 'text/event-stream, application/json',
    });

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify({
        messages: payloadMessages,
        mode: 'llm_summary',
        stream: true,
      }),
      signal: abortController.signal,
    });

    if (response.status === 401) {
      throw new Error('Unauthorized access. Please check your API Token configuration.');
    }

    if (!response.ok) {
      let errorMsg = `Request failed: ${response.status} ${response.statusText}`;
      try {
        const errJson = await response.json();
        if (errJson && errJson.error) {
          errorMsg = errJson.error;
        }
      } catch {
        // ignore
      }
      throw new Error(errorMsg);
    }

    const contentType = response.headers.get('content-type') || '';
    if (contentType.includes('text/event-stream') && response.body) {
      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';

      const ensureAssistantMsg = () => {
        if (!assistantMsg) {
          assistantMsg = {
            role: 'assistant',
            content: '',
            citations: [],
          };
          currentSession.messages.push(assistantMsg);
        }
        return assistantMsg;
      };

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');
        buffer = lines.pop() || '';

        for (const line of lines) {
          const trimmed = line.trim();
          if (!trimmed || !trimmed.startsWith('data:')) continue;
          const dataStr = trimmed.replace(/^data:\s*/, '').trim();
          if (dataStr === '[DONE]') continue;

          try {
            const event = JSON.parse(dataStr);
            if (event.type === 'status') {
              streamStatus.value = event.message || '';
            } else if (event.type === 'citations') {
              const msg = ensureAssistantMsg();
              msg.citations = event.citations || [];
              if (event.usage) {
                apiUsage.value = event.usage;
              }
              streamStatus.value = 'Synthesizing Catholic response...';
              streamScrollToBottom();
            } else if (event.type === 'chunk') {
              const msg = ensureAssistantMsg();
              msg.content += event.delta || '';
              isLoading.value = false;
              isStreaming.value = true;
              streamScrollToBottom();
            } else if (event.type === 'done') {
              const msg = ensureAssistantMsg();
              if (event.response && !msg.content) {
                msg.content = event.response;
              }
              isStreaming.value = false;
              isLoading.value = false;
            } else if (event.type === 'error') {
              throw new Error(event.error || 'Magisterium streaming error');
            }
          } catch (parseErr: any) {
            if (parseErr.message && !parseErr.message.includes('JSON')) {
              throw parseErr;
            }
          }
        }
      }

      if (buffer.trim().startsWith('data:')) {
        const dataStr = buffer.trim().replace(/^data:\s*/, '').trim();
        if (dataStr && dataStr !== '[DONE]') {
          try {
            const event = JSON.parse(dataStr);
            if (event.type === 'chunk') {
              const msg = ensureAssistantMsg();
              msg.content += event.delta || '';
            } else if (event.type === 'done') {
              const msg = ensureAssistantMsg();
              if (event.response && !msg.content) {
                msg.content = event.response;
              }
            }
          } catch {
            // ignore
          }
        }
      }

      if (assistantMsg) {
        currentSession.updatedAt = Date.now();
        saveSessionsToLocalStorage();
      } else {
        throw new Error('No response was generated by Magisterium AI.');
      }
      await refreshUsage();
    } else {
      // Non-streaming fallback
      const res = await response.json() as MagisteriumResponse;
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
        await refreshUsage();
      } else {
        throw new Error('Received an empty response from Magisterium AI.');
      }
    }
  } catch (err: any) {
    if (err.name === 'AbortError') {
      return;
    }
    console.error('Magisterium chat error:', err);
    const msg = err.message || 'An error occurred while calling the Magisterium AI endpoint.';
    chatError.value = msg;
    if (msg.toLowerCase().includes('access code')) {
      clearAccessCode();
      openAccessModal();
    } else if (msg.toLowerCase().includes('daily limit')) {
      await refreshUsage();
    }
    const targetMsg = assistantMsg as ChatMessage | null;
    if (targetMsg && !targetMsg.content.trim() && (!targetMsg.citations || targetMsg.citations.length === 0)) {
      const idx = currentSession.messages.indexOf(targetMsg);
      if (idx !== -1) currentSession.messages.splice(idx, 1);
    }
  } finally {
    isLoading.value = false;
    isStreaming.value = false;
    streamStatus.value = '';
    abortController = null;
    await scrollToBottom();
  }
};
</script>

<style scoped>
/* Neutralize global style.css button defaults so custom icon/pill/card buttons render accurately */
.magisterium-sanctuary button.sanctuary-icon-btn,
.magisterium-sanctuary button.sanctuary-pill-btn,
.magisterium-sanctuary button.sanctuary-tab-btn,
.magisterium-sanctuary button.sanctuary-link-btn,
.magisterium-sanctuary button.sanctuary-card-btn,
.magisterium-sanctuary button.sanctuary-menu-btn {
  min-height: 0;
  font-family: inherit;
}

.magisterium-sanctuary button.sanctuary-icon-btn {
  padding: 0;
}

.magisterium-sanctuary button.sanctuary-icon-btn:not([class*="bg-"]):not([class*="from-"]) {
  background-color: transparent;
}

.magisterium-sanctuary button.sanctuary-icon-btn:not([class*="border-"]) {
  border-color: transparent;
}

.magisterium-sanctuary button.sanctuary-link-btn {
  padding: 0;
  border: none;
  background-color: transparent;
}

.magisterium-sanctuary button.sanctuary-tab-btn:not([class*="border-"]),
.magisterium-sanctuary button.sanctuary-menu-btn:not([class*="border-"]) {
  border-color: transparent;
}

.magisterium-sanctuary button.sanctuary-tab-btn:not([class*="bg-"]),
.magisterium-sanctuary button.sanctuary-menu-btn:not([class*="bg-"]) {
  background-color: transparent;
}

.magisterium-sanctuary button.sanctuary-card-btn {
  border-radius: 1rem;
}

@media (hover: hover) {
  .magisterium-sanctuary button.sanctuary-icon-btn:hover,
  .magisterium-sanctuary button.sanctuary-link-btn:hover,
  .magisterium-sanctuary button.sanctuary-tab-btn:hover,
  .magisterium-sanctuary button.sanctuary-menu-btn:hover,
  .magisterium-sanctuary button.sanctuary-card-btn:hover {
    transform: none;
  }

  .magisterium-sanctuary button.sanctuary-icon-btn:not([class*="hover:bg-"]):hover,
  .magisterium-sanctuary button.sanctuary-link-btn:hover {
    background-color: transparent;
    border-color: transparent;
  }
}
</style>
