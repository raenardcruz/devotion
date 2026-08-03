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

        <!-- Copy Toast Banner -->
        <transition name="fade">
          <div v-if="copyNotice" class="mt-3 max-w-md mx-auto p-2.5 bg-emerald-100 border border-emerald-300 text-emerald-950 font-bold text-xs rounded-xl shadow-xs flex items-center justify-center space-x-2">
            <span>✅</span>
            <span>{{ copyNotice }}</span>
          </div>
        </transition>
      </div>

      <!-- Mobile Navigation Trigger Button -->
      <div class="lg:hidden mb-4 flex items-center justify-between bg-parchment-neutral-light border border-parchment-border rounded-2xl p-3 shadow-sm">
        <div class="flex items-center gap-2 overflow-hidden pr-2">
          <span class="text-parchment-primary">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
          </span>
          <span class="font-serif text-xs font-semibold text-parchment-primary-dark truncate">
            {{ activeSession ? activeSession.title : 'Chat Conversations' }}
          </span>
        </div>
        <button 
          @click="isMobileMenuOpen = true" 
          class="px-3 py-1.5 bg-parchment-primary text-white text-xs font-medium rounded-xl hover:bg-parchment-primary-dark transition-all shadow-xs shrink-0 flex items-center gap-1.5"
        >
          <span>My Chats</span>
        </button>
      </div>

      <!-- Mobile Navigation Drawer -->
      <transition 
        enter-active-class="transition duration-300 ease-out" 
        enter-from-class="opacity-0" 
        enter-to-class="opacity-100" 
        leave-active-class="transition duration-200 ease-in" 
        leave-from-class="opacity-100" 
        leave-to-class="opacity-0"
      >
        <div v-if="isMobileMenuOpen" class="fixed inset-0 z-50 flex justify-start lg:hidden" @click="isMobileMenuOpen = false">
          <!-- Backdrop -->
          <div class="absolute inset-0 bg-parchment-neutral/40 backdrop-blur-xs"></div>
          
          <!-- Panel Content -->
          <div class="relative w-4/5 max-w-xs bg-parchment-bg h-full shadow-2xl p-4 flex flex-col z-10 overflow-y-auto border-r border-parchment-border" @click.stop>
            <div class="flex items-center justify-between pb-3 border-b border-parchment-border mb-3">
              <h3 class="font-serif text-base text-parchment-primary-dark font-bold">Magisterium Sanctuary</h3>
              <button @click="isMobileMenuOpen = false" class="p-1.5 text-parchment-neutral/50 hover:text-parchment-neutral rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <!-- Sidebar Tabs inside Mobile Drawer -->
            <div class="grid grid-cols-2 gap-1 p-1 bg-parchment-neutral-light border border-parchment-border/60 rounded-2xl mb-3 text-center text-xs font-serif font-bold">
              <button 
                @click="sidebarTab = 'local'"
                class="py-1.5 rounded-xl transition-all"
                :class="[sidebarTab === 'local' ? 'bg-parchment-primary text-white shadow-2xs' : 'text-parchment-neutral/70 hover:text-parchment-primary-dark']"
              >
                My Chats
              </button>
              <button 
                @click="sidebarTab = 'public'"
                class="py-1.5 rounded-xl transition-all"
                :class="[sidebarTab === 'public' ? 'bg-parchment-primary text-white shadow-2xs' : 'text-parchment-neutral/70 hover:text-parchment-primary-dark']"
              >
                Public Sanctuary
              </button>
            </div>

            <!-- Sidebar Header for Local -->
            <div v-if="sidebarTab === 'local'" class="flex items-center justify-between pb-3 mb-3 border-b border-parchment-border/60">
              <span class="font-serif font-bold text-[11px] uppercase tracking-wider text-parchment-primary-dark">Saved Conversations</span>
              <button 
                @click="createNewSession(); isMobileMenuOpen = false" 
                class="text-xs bg-parchment-primary text-white px-3 py-1.5 rounded-xl font-semibold hover:bg-parchment-primary-dark transition-all shadow-xs flex items-center space-x-1"
              >
                <span>+ New Chat</span>
              </button>
            </div>

            <!-- Sidebar Header for Public -->
            <div v-else class="flex items-center justify-between pb-3 mb-3 border-b border-parchment-border/60">
              <span class="font-serif font-bold text-[11px] uppercase tracking-wider text-amber-950">Community Discussions</span>
              <button 
                @click="fetchPublicConversations" 
                class="text-[11px] text-amber-900 hover:text-amber-950 underline font-semibold flex items-center space-x-1"
              >
                <span>Refresh</span>
              </button>
            </div>

            <!-- Local Sessions Scroll List -->
            <div v-if="sidebarTab === 'local'" class="overflow-y-auto space-y-2.5 flex-1 scrollbar-thin pr-1">
              <div 
                v-for="session in chatSessions" 
                :key="session.id"
                @click="selectSession(session.id); isMobileMenuOpen = false"
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

            <!-- Public Conversations Scroll List -->
            <div v-else class="overflow-y-auto space-y-2.5 flex-1 scrollbar-thin pr-1">
              <div 
                v-for="pub in publicConversations" 
                :key="pub.id"
                @click="loadPublicSession(pub); isMobileMenuOpen = false"
                class="p-3 rounded-2xl border text-left cursor-pointer transition-all flex items-center justify-between group relative overflow-hidden bg-amber-50/70 border-amber-200/80 hover:bg-amber-100/70"
                :class="[activeSessionId === pub.id ? 'ring-2 ring-amber-500 bg-amber-100/90' : '']"
              >
                <div class="min-w-0 flex-1 pr-2">
                  <div class="text-xs font-bold text-amber-950 truncate font-serif">
                    {{ pub.title || 'Public Discussion' }}
                  </div>
                  <div class="text-[10px] text-amber-900/80 mt-1 flex items-center space-x-1.5 font-serif">
                    <span class="font-semibold text-amber-950">By {{ pub.author_name }}</span>
                    <span class="w-1 h-1 rounded-full bg-amber-400"></span>
                    <span>{{ formatIsoDate(pub.created_at) }}</span>
                  </div>
                </div>
                <button 
                  @click.stop="copyPublicShareLink(pub.id)"
                  class="text-[10px] bg-amber-200/90 hover:bg-amber-300 text-amber-950 font-bold px-2 py-1 rounded-lg border border-amber-300 flex items-center space-x-1 transition-colors"
                  title="Copy Public Share Link"
                >
                  <span>🔗 Share</span>
                </button>
              </div>

              <div v-if="publicConversations.length === 0" class="text-center py-8 px-4 text-xs text-amber-900/60 italic space-y-2 font-serif">
                <div class="text-2xl opacity-40">🌐</div>
                <p>No public community chats published yet. Be the first to share your Q&A!</p>
              </div>
            </div>
          </div>
        </div>
      </transition>

      <!-- Main Layout: Sidebar & Chat Container -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 flex-grow items-start">
        
        <!-- Desktop Sidebar Navigation: Local & Public Chat Sessions -->
        <aside class="hidden lg:flex lg:col-span-4 xl:col-span-3 bg-parchment-neutral-light/80 border border-parchment-border rounded-3xl p-4 flex-col shadow-sm h-auto max-h-[400px] lg:max-h-[750px] lg:h-full backdrop-blur-xs">
          <!-- Sidebar Tabs -->
          <div class="grid grid-cols-2 gap-1 p-1 bg-parchment-bg border border-parchment-border/60 rounded-2xl mb-3 text-center text-xs font-serif font-bold">
            <button 
              @click="sidebarTab = 'local'"
              class="py-1.5 rounded-xl transition-all"
              :class="[sidebarTab === 'local' ? 'bg-parchment-primary text-white shadow-2xs' : 'text-parchment-neutral/70 hover:text-parchment-primary-dark']"
            >
              My Chats
            </button>
            <button 
              @click="sidebarTab = 'public'"
              class="py-1.5 rounded-xl transition-all"
              :class="[sidebarTab === 'public' ? 'bg-parchment-primary text-white shadow-2xs' : 'text-parchment-neutral/70 hover:text-parchment-primary-dark']"
            >
              Public Sanctuary
            </button>
          </div>

          <!-- Sidebar Header for Local -->
          <div v-if="sidebarTab === 'local'" class="flex items-center justify-between pb-3 mb-3 border-b border-parchment-border/60">
            <span class="font-serif font-bold text-[11px] uppercase tracking-wider text-parchment-primary-dark">Saved Conversations</span>
            <button 
              @click="createNewSession" 
              class="text-xs bg-parchment-primary text-white px-3 py-1.5 rounded-xl font-semibold hover:bg-parchment-primary-dark transition-all shadow-xs flex items-center space-x-1"
            >
              <span>+ New Chat</span>
            </button>
          </div>

          <!-- Sidebar Header for Public -->
          <div v-else class="flex items-center justify-between pb-3 mb-3 border-b border-parchment-border/60">
            <span class="font-serif font-bold text-[11px] uppercase tracking-wider text-amber-950">Community Discussions</span>
            <button 
              @click="fetchPublicConversations" 
              class="text-[11px] text-amber-900 hover:text-amber-950 underline font-semibold flex items-center space-x-1"
            >
              <span>Refresh</span>
            </button>
          </div>

          <!-- Local Sessions Scroll List -->
          <div v-if="sidebarTab === 'local'" class="overflow-y-auto space-y-2.5 flex-1 scrollbar-thin pr-1 max-h-[300px] lg:max-h-[620px]">
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

          <!-- Public Conversations Scroll List -->
          <div v-else class="overflow-y-auto space-y-2.5 flex-1 scrollbar-thin pr-1 max-h-[300px] lg:max-h-[620px]">
            <div 
              v-for="pub in publicConversations" 
              :key="pub.id"
              @click="loadPublicSession(pub)"
              class="p-3 rounded-2xl border text-left cursor-pointer transition-all flex items-center justify-between group relative overflow-hidden bg-amber-50/70 border-amber-200/80 hover:bg-amber-100/70"
              :class="[activeSessionId === pub.id ? 'ring-2 ring-amber-500 bg-amber-100/90' : '']"
            >
              <div class="min-w-0 flex-1 pr-2">
                <div class="text-xs font-bold text-amber-950 truncate font-serif">
                  {{ pub.title || 'Public Discussion' }}
                </div>
                <div class="text-[10px] text-amber-900/80 mt-1 flex items-center space-x-1.5 font-serif">
                  <span class="font-semibold text-amber-950">By {{ pub.author_name }}</span>
                  <span class="w-1 h-1 rounded-full bg-amber-400"></span>
                  <span>{{ formatIsoDate(pub.created_at) }}</span>
                </div>
              </div>
              <button 
                @click.stop="copyPublicShareLink(pub.id)"
                class="text-[10px] bg-amber-200/90 hover:bg-amber-300 text-amber-950 font-bold px-2 py-1 rounded-lg border border-amber-300 flex items-center space-x-1 transition-colors"
                title="Copy Public Share Link"
              >
                <span>🔗 Share</span>
              </button>
            </div>

            <div v-if="publicConversations.length === 0" class="text-center py-8 px-4 text-xs text-amber-900/60 italic space-y-2 font-serif">
              <div class="text-2xl opacity-40">🌐</div>
              <p>No public community chats published yet. Be the first to share your Q&A!</p>
            </div>
          </div>
        </aside>

        <!-- Main Chat Box Panel -->
        <div class="lg:col-span-8 xl:col-span-9 flex flex-col bg-parchment-neutral-light/60 border border-parchment-border rounded-3xl shadow-sm overflow-hidden h-[550px] sm:h-[650px] lg:h-[750px]">
          
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
            <div 
              v-for="(msg, index) in activeMessages" 
              :key="index" 
              :ref="el => { if (el) messageRefs[index] = el as HTMLElement }"
              class="flex flex-col space-y-2 animate-fade-in"
            >
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
                    <div 
                      @click="toggleCitations(index)"
                      class="flex items-center justify-between cursor-pointer select-none border-b border-amber-200/60 pb-2.5 group"
                    >
                      <span class="text-xs font-bold text-amber-950 uppercase tracking-wider flex items-center space-x-1.5">
                        <span>📚 Magisterium Citations</span>
                        <span class="bg-amber-200/80 text-amber-900 text-[10px] px-2 py-0.5 rounded-full font-mono font-bold">{{ msg.citations.length }}</span>
                      </span>
                      <button 
                        type="button" 
                        class="text-[11px] text-amber-900 hover:text-amber-950 font-semibold uppercase tracking-wider flex items-center space-x-1 border-none shadow-none py-0 px-1 bg-transparent hover:bg-amber-200/50 rounded-lg transition-colors"
                      >
                        <span>{{ expandedCitations.has(index) ? 'Hide Citations' : 'Show Citations' }}</span>
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
                          class="transition-transform duration-200"
                          :class="{ 'rotate-180': expandedCitations.has(index) }"
                        >
                          <path d="m6 9 6 6 6-6"/>
                        </svg>
                      </button>
                    </div>

                    <div v-if="expandedCitations.has(index)" class="grid grid-cols-1 gap-3 animate-fade-in">
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
          <div class="p-4 bg-parchment-neutral-light/90 border-t border-parchment-border backdrop-blur-xs space-y-2.5">
            <!-- Mode Selector Bar -->
            <div class="flex items-center justify-between px-1 text-xs font-serif">
              <span class="text-stone-600 font-semibold flex items-center space-x-1">
                <span>Response Mode:</span>
              </span>
              <div class="inline-flex p-0.5 bg-stone-200/80 border border-parchment-border/70 rounded-xl shadow-inner">
                <button
                  type="button"
                  @click="completionMode = 'magisterium'"
                  class="px-3 py-1 rounded-lg text-xs font-bold transition-all flex items-center space-x-1"
                  :class="[completionMode === 'magisterium' ? 'bg-amber-800 text-white shadow-2xs' : 'text-stone-700 hover:text-amber-900']"
                  title="Use direct Magisterium AI answer engine"
                >
                  <span>📜 Direct Magisterium AI</span>
                </button>
                <button
                  type="button"
                  @click="completionMode = 'llm_summary'"
                  class="px-3 py-1 rounded-lg text-xs font-bold transition-all flex items-center space-x-1"
                  :class="[completionMode === 'llm_summary' ? 'bg-amber-800 text-white shadow-2xs' : 'text-stone-700 hover:text-amber-900']"
                  title="Search citations and generate summary with configured LLM (Ollama or Gemini)"
                >
                  <span>🤖 Custom LLM Summary</span>
                </button>
              </div>
            </div>

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
            <div class="mt-3 flex items-center justify-between text-[11px] text-parchment-neutral/60 px-1 font-serif border-t border-parchment-border/40 pt-2.5">
              <div class="flex items-center space-x-2">
                <span>Auto-saved to local browser storage</span>
                <button 
                  v-if="activeMessages.length > 0"
                  @click="openPublishModal"
                  class="ml-2 px-2.5 py-1 rounded-lg bg-amber-100/90 text-amber-950 border border-amber-300/80 font-bold hover:bg-amber-200 transition-all flex items-center space-x-1 shadow-2xs"
                >
                  <span>🌐 Publish to Public Sanctuary</span>
                </button>
              </div>

              <router-link to="/admin" class="hover:text-parchment-primary transition-colors flex items-center space-x-1">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
                <span>API Settings</span>
              </router-link>
            </div>
          </div>
        </div>

      </div>

      <!-- Publish Conversation to Database Modal -->
      <Teleport to="body">
        <div v-if="showPublishModal" class="fixed inset-0 z-50 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4">
          <div class="bg-parchment-neutral-light border border-amber-300 rounded-3xl p-6 max-w-lg w-full shadow-2xl space-y-4 font-serif relative">
            <div class="flex items-center justify-between border-b border-parchment-border/60 pb-3">
              <div class="flex items-center space-x-2">
                <span class="text-xl">🌐</span>
                <h3 class="font-serif text-lg font-bold text-parchment-primary-dark">Publish Conversation to Database</h3>
              </div>
              <button @click="showPublishModal = false" class="text-stone-400 hover:text-stone-700 text-lg">✕</button>
            </div>

            <!-- Public Warning Banner -->
            <div class="bg-amber-100/70 border border-amber-300/80 rounded-2xl p-4 text-xs text-amber-950 space-y-1.5">
              <div class="font-bold flex items-center space-x-1 text-amber-900 uppercase tracking-wider">
                <span>⚠️ Public Community Disclosure</span>
              </div>
              <p class="leading-relaxed">
                By publishing this discussion, it will be stored securely in the PostgreSQL database and <strong>made publicly visible to all users</strong> on the site under the <strong>Public Sanctuary</strong> tab.
              </p>
            </div>

            <form @submit.prevent="submitPublishConversation" class="space-y-4 pt-1">
              <div>
                <label class="block text-xs font-bold uppercase tracking-wider text-amber-950 mb-1">Your Name / Title</label>
                <input 
                  v-model="authorName" 
                  type="text" 
                  required 
                  placeholder="e.g. Brother John, Maria, Pilgrim..."
                  class="w-full px-4 py-2.5 bg-white border border-amber-300 rounded-xl text-sm text-stone-800 focus:outline-none focus:ring-2 focus:ring-amber-500 font-sans"
                />
                <p class="text-[11px] text-stone-500 mt-1">This name will be attributed alongside your published conversation.</p>
              </div>

              <div v-if="publishMessage" class="p-3 bg-emerald-100 text-emerald-900 border border-emerald-300 rounded-xl text-xs font-semibold">
                {{ publishMessage }}
              </div>

              <div class="flex justify-end space-x-3 pt-2">
                <AppButton type="button" variant="secondary" @click="showPublishModal = false">
                  Cancel
                </AppButton>
                <AppButton type="submit" variant="primary" :disabled="isPublishing || !authorName.trim()">
                  {{ isPublishing ? 'Publishing...' : 'Confirm & Publish Publicly' }}
                </AppButton>
              </div>
            </form>
          </div>
        </div>
      </Teleport>
    </main>

    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue';
import { useRoute } from 'vue-router';
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

interface PublicConversation {
  id: string;
  author_name: string;
  title: string;
  messages: ChatMessage[];
  created_at: string;
}

const route = useRoute();
const { fetchWithAuth } = useDevotionApi();

const LOCAL_STORAGE_KEY = 'magisterium_chat_sessions_v1';

const sidebarTab = ref<'local' | 'public'>('local');
const chatSessions = ref<ChatSession[]>([]);
const publicConversations = ref<PublicConversation[]>([]);
const activeSessionId = ref<string>('');
const isMobileMenuOpen = ref(false);
const inputQuery = ref('');
const completionMode = ref<'magisterium' | 'llm_summary'>('llm_summary');
const isLoading = ref(false);
const chatError = ref<string | null>(null);
const chatContainer = ref<HTMLDivElement | null>(null);
const messageRefs = ref<HTMLElement[]>([]);
const apiUsage = ref<ApiUsage | null>(null);

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

const activeMessages = computed(() => {
  const current = activeSession.value;
  return current ? current.messages : [];
});

onMounted(async () => {
  loadSessionsFromLocalStorage();
  await fetchPublicConversations();

  // Check if route has shared public conversation ID param
  const paramId = route.params.id as string;
  if (paramId) {
    const pub = publicConversations.value.find(p => p.id === paramId);
    if (pub) {
      sidebarTab.value = 'public';
      loadPublicSession(pub);
      return;
    }
  }

  const first = chatSessions.value[0];
  if (first) {
    activeSessionId.value = first.id;
  } else {
    createNewSession();
  }
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

const loadPublicSession = (pub: PublicConversation) => {
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
  html = html.replace(/^\&gt; (.*$)/gim, '<blockquote class="border-l-4 border-amber-500 pl-3 py-1 my-2 italic text-stone-700 bg-amber-50/50 rounded-r-lg">$1</blockquote>');

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
      body: JSON.stringify({
        messages: payloadMessages,
        mode: completionMode.value,
      }),
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
