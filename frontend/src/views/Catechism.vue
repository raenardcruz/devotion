<template>
  <div class="min-h-screen text-[#322D29] flex flex-col pb-28 sm:pb-20 selection:bg-[#72383D]/20 relative z-10">
    <!-- Global Header -->
    <TopNav />

    <!-- Main Content Container -->
    <div class="relative z-10 max-w-7xl mx-auto px-4 py-6 md:py-10 flex flex-col min-h-screen w-full">
      
      <!-- Top Header Area -->
      <header class="mb-8 animate-fade-in-down">
        <div class="flex items-center gap-2 mb-2">
          <span class="text-[#72383D] text-[11px] font-bold uppercase tracking-[0.2em]">OFFICIAL TEXT</span>
        </div>
        <h1 class="text-3xl md:text-5xl font-serif bg-gradient-to-r from-[#72383D] via-[#AC9C8D] to-[#322D29] bg-clip-text text-transparent font-bold mb-3 tracking-tight">
          Catechism of the Catholic Church
        </h1>
        <p class="text-[#322D29]/70 italic text-sm md:text-base max-w-3xl leading-relaxed">
          "It is the task of the Catechism to present the essential and fundamental content of Catholic doctrine."
        </p>
      </header>

      <!-- Overview Landing View (Shown when no part or search active) -->
      <div v-if="!selectedPart && !searchQuery.trim()" class="space-y-8 animate-fade-in-up">
        
        <!-- Universal Search Card (Positioned at Top) -->
        <div class="bg-white/80 border border-[#D1C7BD]/80 rounded-3xl p-6 md:p-8 shadow-sm backdrop-blur-md flex flex-col md:flex-row items-stretch md:items-center justify-between gap-6">
          <div class="shrink-0">
            <h3 class="text-xl font-serif font-bold text-[#72383D] mb-1">Universal Search</h3>
            <p class="text-[#322D29]/70 text-xs md:text-sm">Find specific paragraphs or theological topics across all four parts.</p>
          </div>
          <div class="flex items-center gap-3 flex-1 w-full max-w-2xl">
            <div class="relative flex-1">
              <input 
                v-model="universalSearchInput"
                @keyup.enter="executeUniversalSearch"
                type="text" 
                placeholder="Search paragraph number (e.g. 27, 1324, 5-10)..." 
                class="w-full pl-11 pr-4 py-3.5 rounded-2xl border border-[#D1C7BD] bg-white text-sm text-[#322D29] placeholder-[#322D29]/40 focus:border-[#72383D] focus:ring-2 focus:ring-[#72383D]/20 focus:outline-none shadow-2xs transition-all"
              />
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-[#72383D]">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                </svg>
              </div>
            </div>
            <button 
              @click="executeUniversalSearch" 
              class="px-6 py-3.5 bg-gradient-to-r from-[#72383D] to-[#322D29] hover:from-[#8B464C] hover:to-[#453E38] text-white text-xs font-bold rounded-2xl transition-all shadow-xs cursor-pointer shrink-0"
            >
              Jump to Reference
            </button>
          </div>
        </div>

        <!-- Part One Hero Featured Card -->
        <div class="bg-white/80 border border-[#D1C7BD]/80 rounded-3xl p-6 md:p-8 shadow-sm backdrop-blur-md grid grid-cols-1 md:grid-cols-3 gap-8 items-center group hover:shadow-md transition-all duration-300">
          <div class="md:col-span-2 flex flex-col justify-between h-full space-y-4">
            <div>
              <div class="flex items-center justify-between mb-4">
                <span class="text-[#72383D] text-xs font-serif font-bold tracking-wider uppercase">Part One</span>
                <span class="text-[#72383D]/60">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/>
                    <path d="M22 3h-6a4 4 0 0 1 4 4v14a3 3 0 0 1 3-3h7z"/>
                  </svg>
                </span>
              </div>
              <h2 class="text-2xl md:text-3xl font-serif font-bold text-[#72383D] mb-3">The Profession of Faith</h2>
              <p class="text-[#322D29]/80 text-sm md:text-base leading-relaxed mb-6">
                Focusing on the Creed—the summary of our core beliefs. Exploring the nature of God, the mystery of the Incarnation, and the Holy Spirit's role in the Church.
              </p>
            </div>

            <!-- Topic Pills -->
            <div class="flex flex-wrap gap-2.5 pt-2">
              <button 
                @click="openPart(1)" 
                class="px-4 py-2 rounded-xl bg-[#EFE9E1] text-[#72383D] font-bold text-[11px] uppercase tracking-wider hover:bg-[#D1C7BD] transition-all cursor-pointer shadow-2xs border border-[#D1C7BD]/50"
              >
                I BELIEVE
              </button>
              <button 
                @click="openPart(1)" 
                class="px-4 py-2 rounded-xl bg-[#EFE9E1] text-[#72383D] font-bold text-[11px] uppercase tracking-wider hover:bg-[#D1C7BD] transition-all cursor-pointer shadow-2xs border border-[#D1C7BD]/50"
              >
                THE TRINITY
              </button>
              <button 
                @click="openPart(1)" 
                class="px-4 py-2 rounded-xl bg-[#EFE9E1] text-[#72383D] font-bold text-[11px] uppercase tracking-wider hover:bg-[#D1C7BD] transition-all cursor-pointer shadow-2xs border border-[#D1C7BD]/50"
              >
                ETERNAL LIFE
              </button>
            </div>
          </div>

          <!-- Hero Image -->
          <div class="md:col-span-1 h-56 md:h-full overflow-hidden rounded-2xl shadow-sm border border-[#D1C7BD]/40 relative">
            <img 
              src="/images/catechism_hero.jpg" 
              alt="Catechism sacred illuminated book" 
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
            />
          </div>
        </div>

        <!-- 3-Column Grid for Parts Two, Three, Four -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          
          <!-- Part Two Card -->
          <div class="bg-white/80 border border-[#D1C7BD]/80 rounded-3xl p-6 shadow-sm backdrop-blur-md flex flex-col justify-between hover:shadow-md transition-all duration-300">
            <div>
              <span class="text-[#72383D] text-xs font-serif font-bold tracking-wider uppercase block mb-3">Part Two</span>
              <h3 class="text-lg md:text-xl font-serif font-bold text-[#72383D] mb-2">The Celebration of the Christian Mystery</h3>
              <p class="text-[#322D29]/75 text-xs md:text-sm leading-relaxed mb-6">
                The sacramental economy: baptism, eucharist, and the life of the liturgy as a bridge between heaven and earth.
              </p>
            </div>
            <button 
              @click="openPart(2)" 
              class="inline-flex items-center gap-2 text-[#72383D] font-bold text-xs hover:text-[#322D29] transition-colors cursor-pointer self-start group/btn"
            >
              <span>Open Volume</span>
              <svg class="group-hover/btn:translate-x-1 transition-transform" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <line x1="5" y1="12" x2="19" y2="12"></line>
                <polyline points="12 5 19 12 12 19"></polyline>
              </svg>
            </button>
          </div>

          <!-- Part Three Card -->
          <div class="bg-white/80 border border-[#D1C7BD]/80 rounded-3xl p-6 shadow-sm backdrop-blur-md flex flex-col justify-between hover:shadow-md transition-all duration-300">
            <div>
              <span class="text-[#72383D] text-xs font-serif font-bold tracking-wider uppercase block mb-3">Part Three</span>
              <h3 class="text-lg md:text-xl font-serif font-bold text-[#72383D] mb-2">Life in Christ</h3>
              <p class="text-[#322D29]/75 text-xs md:text-sm leading-relaxed mb-6">
                The moral life, the beatitudes, and the dignity of the human person. A guide to living the faith through action.
              </p>
            </div>
            <button 
              @click="openPart(3)" 
              class="inline-flex items-center gap-2 text-[#72383D] font-bold text-xs hover:text-[#322D29] transition-colors cursor-pointer self-start group/btn"
            >
              <span>Open Volume</span>
              <svg class="group-hover/btn:translate-x-1 transition-transform" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <line x1="5" y1="12" x2="19" y2="12"></line>
                <polyline points="12 5 19 12 12 19"></polyline>
              </svg>
            </button>
          </div>          <!-- Part Four Card -->
          <div class="bg-white/80 border border-[#D1C7BD]/80 rounded-3xl p-6 shadow-sm backdrop-blur-md flex flex-col justify-between hover:shadow-md transition-all duration-300">
            <div>
              <span class="text-[#72383D] text-xs font-serif font-bold tracking-wider uppercase block mb-3">Part Four</span>
              <h3 class="text-lg md:text-xl font-serif font-bold text-[#72383D] mb-2">Christian Prayer</h3>
              <p class="text-[#322D29]/75 text-xs md:text-sm leading-relaxed mb-6">
                The meaning and importance of prayer in the Christian life, with a deep dive into the Lord's Prayer.
              </p>
            </div>
            <button 
              @click="openPart(4)" 
              class="inline-flex items-center gap-2 text-[#72383D] font-bold text-xs hover:text-[#322D29] transition-colors cursor-pointer self-start group/btn"
            >
              <span>Open Volume</span>
              <svg class="group-hover/btn:translate-x-1 transition-transform" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <line x1="5" y1="12" x2="19" y2="12"></line>
                <polyline points="12 5 19 12 12 19"></polyline>
              </svg>
            </button>
          </div>
        </div>

        <!-- Recently Viewed Paragraphs Section -->
        <div>
          <div class="flex items-center gap-2 mb-4">
            <span class="w-1 h-5 bg-gradient-to-r from-[#72383D] to-[#322D29] rounded-full inline-block"></span>
            <h3 class="font-serif text-sm font-bold text-[#72383D] tracking-wide">Recently Viewed Paragraphs</h3>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div 
              v-for="item in recentParagraphs" 
              :key="item.id"
              @click="jumpToParagraph(item.id)"
              class="bg-white/80 border border-[#D1C7BD]/60 rounded-2xl p-5 shadow-2xs hover:shadow-md transition-all cursor-pointer group"
            >
              <span class="text-xs font-serif font-bold text-[#72383D] block mb-2 group-hover:text-[#322D29]">Para. {{ item.id }}</span>
              <p class="text-[#322D29]/75 text-xs line-clamp-3 leading-relaxed italic">"{{ item.text }}"</p>
            </div>
          </div>
        </div>

      </div>

      <!-- Reading View Header Controls (Shown when part or search active) -->
      <div v-else class="mb-4 flex items-center justify-between">
        <button 
          @click="resetSelection" 
          class="inline-flex items-center gap-2 px-4 py-2 bg-white/70 border border-[#D1C7BD] text-[#72383D] hover:text-[#322D29] hover:bg-white rounded-xl text-xs font-bold shadow-2xs transition-all cursor-pointer"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="19" y1="12" x2="5" y2="12"></line>
            <polyline points="12 19 5 12 12 5"></polyline>
          </svg>
          <span>Catechism Overview</span>
        </button>

        <span class="text-xs text-[#322D29]/60 font-medium">
          {{ selectedPart ? `Part ${selectedPart.part}` : 'Search Results' }}
        </span>
      </div>

      <!-- Mobile Navigation Floating Trigger Button -->
      <div v-if="selectedPart || searchQuery.trim()" class="lg:hidden mb-4 flex items-center justify-between bg-white/60 border border-[#D1C7BD]/60 rounded-2xl p-3 shadow-sm backdrop-blur-md">
        <div class="flex items-center gap-2 overflow-hidden pr-2">
          <span class="text-[#72383D]">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="3" y1="12" x2="21" y2="12"></line>
              <line x1="3" y1="6" x2="21" y2="6"></line>
              <line x1="3" y1="18" x2="21" y2="18"></line>
            </svg>
          </span>
          <span class="font-serif text-xs font-semibold text-[#72383D] truncate">
            {{ selectedPart ? `Part ${selectedPart.part}: ${selectedPart.title}` : 'Structure & Navigation' }}
          </span>
        </div>
        <button 
          @click="isMobileMenuOpen = true" 
          class="px-3.5 py-1.5 bg-gradient-to-r from-[#72383D] to-[#322D29] text-white text-xs font-medium rounded-xl hover:from-[#8B464C] hover:to-[#453E38] transition-all shadow-xs shrink-0 flex items-center gap-1.5 cursor-pointer"
        >
          <span>Contents</span>
        </button>
      </div>

      <!-- Mobile Navigation Slide-Over Drawer -->
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
          <div class="absolute inset-0 bg-[#322D29]/40 backdrop-blur-sm"></div>
          
          <!-- Drawer Content Panel -->
          <div class="relative w-4/5 max-w-xs bg-white h-full shadow-2xl p-5 flex flex-col z-10 overflow-y-auto border-r border-[#D1C7BD]" @click.stop>
            <div class="flex items-center justify-between pb-4 border-b border-[#D1C7BD]/60 mb-4">
              <h3 class="font-serif text-base text-[#72383D] font-bold">Catechism Structure</h3>
              <button @click="isMobileMenuOpen = false" class="p-1.5 text-[#72383D] hover:text-[#322D29] rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <div class="space-y-4 flex-grow">
              <div v-for="part in structure.catechism.structure" :key="part.part" class="space-y-2">
                <!-- Part Header -->
                <button 
                  @click="togglePart(part)"
                  class="w-full text-left font-medium text-sm text-[#322D29]/70 hover:text-[#72383D] transition-all flex justify-between items-start group p-2 rounded-lg hover:bg-[#D1C7BD]/20 border-none shadow-none cursor-pointer"
                  :class="{'bg-[#EFE9E1]/80 text-[#72383D] border border-[#D1C7BD]': selectedPart?.part === part.part}"
                >
                  <span class="font-serif font-bold mr-2 text-[#72383D] group-hover:text-[#322D29] transition-colors">Part {{ part.part }}</span>
                  <span class="flex-1">{{ part.title }}</span>
                </button>
                
                <!-- Sections (only if part selected) -->
                <transition name="fade-slide">
                  <div v-if="selectedPart?.part === part.part" class="pl-4 space-y-2 border-l-2 border-[#72383D]/30 ml-3 mt-2">
                    <div v-for="section in part.sections" :key="section.section">
                       <button 
                        @click="toggleSection(section)"
                        class="w-full text-left text-xs text-[#322D29]/70 hover:text-[#72383D] transition-colors py-1.5 px-2 rounded block hover:bg-white/60 border-none shadow-none cursor-pointer"
                        :class="{'text-[#72383D] font-medium bg-[#72383D]/5': selectedSection?.section === section.section}"
                      >
                        {{ section.title }}
                      </button>
  
                      <!-- Chapters (only if section selected) -->
                      <transition name="fade-slide">
                        <div v-if="selectedSection?.section === section.section" class="pl-3 mt-1 space-y-1 border-l border-[#72383D]/10 ml-1.5">
                            <button
                                v-for="chapter in section.chapters"
                                :key="chapter.chapter"
                                @click="toggleChapter(chapter); isMobileMenuOpen = false"
                                class="w-full text-left text-[11px] text-[#322D29]/60 hover:text-[#72383D] transition-colors py-1 px-2 rounded block hover:bg-white/60 border-none shadow-none cursor-pointer"
                                :class="{'text-[#72383D] font-bold bg-[#72383D]/10': selectedChapter?.chapter === chapter.chapter}"
                            >
                                {{ chapter.title }}
                            </button>
                        </div>
                      </transition>
                    </div>
                  </div>
                </transition>
              </div>
            </div>
          </div>
        </div>
      </transition>

      <!-- Active Reader Layout (Sidebar + Content) -->
      <main v-if="selectedPart || searchQuery.trim()" class="flex-grow grid grid-cols-1 lg:grid-cols-4 gap-8">
        
        <!-- Desktop Sidebar Navigation (Structure) -->
        <aside class="hidden lg:block lg:col-span-1 h-fit lg:sticky lg:top-24 lg:overflow-y-auto lg:max-h-[calc(100vh-8rem)] animate-fade-in-up delay-100 scrollbar-hide">
          <div class="bg-white/60 border border-[#D1C7BD]/60 rounded-3xl p-6 shadow-xs backdrop-blur-md">
            <h3 class="font-serif text-xs mb-6 text-[#72383D] tracking-widest uppercase font-bold border-b border-[#D1C7BD]/60 pb-2.5">Structure</h3>
            
            <div class="space-y-4">
              <div v-for="part in structure.catechism.structure" :key="part.part" class="space-y-2">
                <!-- Part Header -->
                <button 
                  @click="togglePart(part)"
                  class="w-full text-left font-medium text-sm text-[#322D29]/70 hover:text-[#72383D] transition-all flex justify-between items-start group p-2 rounded-xl hover:bg-white/80 border-none shadow-none cursor-pointer"
                  :class="{'bg-[#EFE9E1]/80 text-[#72383D] border border-[#D1C7BD] font-bold': selectedPart?.part === part.part}"
                >
                  <span class="font-serif font-bold mr-2 text-[#72383D] group-hover:text-[#322D29] transition-colors">Part {{ part.part }}</span>
                  <span class="flex-1">{{ part.title }}</span>
                </button>
                
                <!-- Sections (only if part selected) -->
                <transition name="fade-slide">
                  <div v-if="selectedPart?.part === part.part" class="pl-4 space-y-2 border-l-2 border-[#72383D]/30 ml-3 mt-2">
                    <div v-for="section in part.sections" :key="section.section">
                       <button 
                        @click="toggleSection(section)"
                        class="w-full text-left text-xs text-[#322D29]/70 hover:text-[#72383D] transition-colors py-1.5 px-2 rounded-lg block hover:bg-white/60 border-none shadow-none cursor-pointer"
                        :class="{'text-[#72383D] font-bold bg-[#EFE9E1]/80': selectedSection?.section === section.section}"
                      >
                        {{ section.title }}
                      </button>
  
                      <!-- Chapters (only if section selected) -->
                      <transition name="fade-slide">
                        <div v-if="selectedSection?.section === section.section" class="pl-3 mt-1 space-y-1 border-l border-[#72383D]/20 ml-1.5">
                            <button
                                v-for="chapter in section.chapters"
                                :key="chapter.chapter"
                                @click="toggleChapter(chapter)"
                                class="w-full text-left text-[11px] text-[#322D29]/60 hover:text-[#72383D] transition-colors py-1 px-2 rounded block hover:bg-white/60 border-none shadow-none cursor-pointer"
                                :class="{'text-[#72383D] font-bold bg-[#EFE9E1]': selectedChapter?.chapter === chapter.chapter}"
                            >
                                {{ chapter.title }}
                            </button>
                        </div>
                      </transition>
                    </div>
                  </div>
                </transition>
              </div>
            </div>
          </div>
        </aside>

        <!-- Content Area -->
        <div id="catechism-content" class="lg:col-span-3 flex flex-col gap-6 animate-fade-in-up delay-150">
            
            <!-- Search Input -->
            <div class="relative group">
                <input 
                    v-model="searchQuery" 
                    type="text" 
                    placeholder="Search paragraph number (e.g. 5-10, 20)..." 
                    class="w-full pl-12 pr-4 py-4 rounded-2xl border border-[#D1C7BD]/60 bg-white/60 text-[#322D29] placeholder-[#322D29]/40 focus:bg-white/90 focus:border-[#72383D] shadow-xs backdrop-blur-md transition-all outline-none"
                >
                <div class="absolute left-4 top-1/2 -translate-y-1/2 text-[#72383D]">
                  <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="11" cy="11" r="8"></circle>
                    <line x1="21" x2="16.65" y1="21" y2="16.65"></line>
                  </svg>
                </div>
            </div>

            <!-- Active Path Display (Breadcrumbs) -->
            <div v-if="!searchQuery && (selectedPart || selectedSection || selectedChapter)" class="bg-white/60 backdrop-blur-md border border-[#D1C7BD]/60 rounded-2xl p-5 text-sm flex flex-wrap gap-2 items-center shadow-xs">
                <span v-if="selectedPart" class="font-serif font-bold text-[#72383D] text-base">Part {{ selectedPart.part }}: {{ selectedPart.title }}</span>
                <span v-if="selectedSection" class="text-[#322D29]/30 mx-1">/</span>
                <span v-if="selectedSection" class="font-medium text-[#322D29]/80">{{ selectedSection.title }}</span>
                <span v-if="selectedChapter" class="text-[#322D29]/30 mx-1">/</span>
                <span v-if="selectedChapter" class="text-[#322D29]/60">{{ selectedChapter.title }}</span>
            </div>

          <!-- Paragraphs List -->
          <div class="space-y-6">
            <template v-if="filteredParagraphs.length > 0">
                <div v-for="(paragraph, index) in filteredParagraphs" :key="paragraph.id" 
                     class="bg-white/60 border border-[#D1C7BD]/60 backdrop-blur-md rounded-3xl p-6 md:p-8 shadow-xs transition-all duration-300 group/card"
                     :class="{ 'opacity-60': readParagraphs.has(paragraph.id) }"
                     :style="{ animationDelay: `${index * 50}ms` }">
                  <div class="flex items-start gap-5">
                    <!-- Read Checkbox Circle -->
                    <span class="flex-shrink-0 w-10 h-10 rounded-full border-2 border-[#72383D] text-[#72383D] flex items-center justify-center font-serif font-bold text-sm bg-white relative cursor-pointer transition-all duration-300 hover:bg-[#72383D] hover:text-white"
                        :class="{ '!bg-[#72383D] !text-white': readParagraphs.has(paragraph.id) }"
                        @click="toggleRead(paragraph.id)">
                      <span class="transition-opacity duration-200" :class="{ 'opacity-0': readParagraphs.has(paragraph.id) }">{{ paragraph.id }}</span>
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" 
                        class="absolute opacity-0 transition-opacity duration-200 text-white"
                        :class="{ '!opacity-100': readParagraphs.has(paragraph.id) }">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    </span>
                    
                    <div class="prose max-w-none flex-1">
                      <div class="text-[10px] text-[#72383D] mb-1.5 font-bold tracking-wider uppercase select-none">{{ getBreadcrumb(paragraph.id) }}</div>
                      <!-- Paragraph text rendering with reference helper link -->
                      <CatechismText :text="paragraph.text" :paragraphs="paragraphs" @show-reference="openModal" />
                    </div>
                  </div>
                </div>
            </template>

            <!-- Empty State / Instructions -->
             <div v-else class="text-center py-24 bg-white/40 backdrop-blur-md rounded-3xl border border-dashed border-[#D1C7BD]">
                <div class="w-16 h-16 bg-white/70 rounded-full flex items-center justify-center mx-auto mb-4 shadow-xs">
                    <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="text-[#72383D]">
                        <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/>
                        <path d="M22 3h-6a4 4 0 0 1 4 4v14a3 3 0 0 1 3-3h7z"/>
                    </svg>
                </div>
                <h3 class="text-lg font-serif text-[#72383D] font-bold mb-1.5">
                    {{ searchQuery ? `No paragraph found for "${searchQuery}"` : 'Begin Your Study' }}
                </h3>
                <p class="text-[#322D29]/60 max-w-sm mx-auto leading-relaxed text-xs md:text-sm">
                    {{ searchQuery ? 'Try searching for a valid paragraph number.' : 'Select a section from the structure on the left or search for a specific paragraph number above.' }}
                </p>
             </div>
          </div>
        </div>
      </main>
    </div>

    <!-- Floating Bookmark Action Button -->
    <button 
      @click="isMobileMenuOpen = true"
      class="fixed bottom-24 right-6 sm:bottom-10 sm:right-8 z-40 bg-gradient-to-r from-[#72383D] to-[#322D29] hover:from-[#8B464C] hover:to-[#453E38] text-white p-4 rounded-2xl shadow-xl transition-all duration-300 cursor-pointer flex items-center justify-center border border-white/20 hover:scale-105"
      title="Open Contents Navigation"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor" stroke="none">
        <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path>
      </svg>
    </button>

    <!-- Reference Modal -->
    <transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100" leave-to-class="opacity-0">
        <div v-if="showModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4" @click="showModal = false">
            <!-- Backdrop -->
            <div class="absolute inset-0 bg-[#322D29]/40 backdrop-blur-sm"></div>
            
            <!-- Modal Content -->
            <div class="relative w-full max-w-2xl bg-white/95 backdrop-blur-md border border-[#D1C7BD] rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[80vh]" @click.stop>
                <!-- Header -->
                <div class="p-5 border-b border-[#D1C7BD]/60 flex justify-between items-center bg-white/70">
                    <h3 class="text-lg font-serif text-[#72383D] font-bold">Referenced Paragraphs</h3>
                    <button @click="showModal = false" class="p-2 hover:bg-[#D1C7BD]/30 rounded-full transition-colors text-[#322D29]/50 hover:text-[#72383D] border-none shadow-none cursor-pointer">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <line x1="18" x2="6" y1="6" y2="18"></line>
                            <line x1="6" x2="18" y1="6" y2="18"></line>
                        </svg>
                    </button>
                </div>

                <!-- Body -->
                <div class="p-6 overflow-y-auto space-y-4">
                    <div v-for="p in modalContent" :key="p.id" class="p-5 rounded-2xl bg-white/70 border border-[#D1C7BD]/60">
                        <div class="flex items-start gap-4">
                            <span class="text-[#72383D] font-bold text-lg font-serif opacity-85 select-none">§{{ p.id }}</span>
                            <p class="text-[#322D29] font-serif leading-relaxed text-sm md:text-base">{{ p.text }}</p>
                        </div>
                    </div>
                </div>
                
                <!-- Footer -->
                <div class="p-4 bg-white/70 text-center border-t border-[#D1C7BD]/60">
                     <span class="text-[10px] text-[#72383D] uppercase tracking-widest font-semibold">Catechism of the Catholic Church</span>
                </div>
            </div>
        </div>
    </transition>

    <!-- Global Footer -->
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import catechismData from '../components/catechism/catechism.json';
import structureData from '../components/catechism/catechism-sections.json';
import CatechismText from '../components/catechism/CatechismText.vue';

import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';

// Types
interface Paragraph {
  id: number;
  text: string;
}

interface Chapter {
  chapter: number;
  title: string;
  paragraphs: string;
}

interface Section {
  section: number;
  title: string;
  paragraph_range: string;
  chapters: Chapter[];
}

interface Part {
  part: number;
  title: string;
  paragraph_range: string;
  sections: Section[];
}

interface Structure {
    catechism: {
        title: string;
        structure: Part[];
    }
}

const route = useRoute();

// State
const structure = ref<Structure>(structureData as unknown as Structure);
const paragraphs = ref<Paragraph[]>(catechismData as unknown as Paragraph[]);
const searchQuery = ref('');
const universalSearchInput = ref('');
const isMobileMenuOpen = ref(false);

const selectedPart = ref<Part | null>(null);
const selectedSection = ref<Section | null>(null);
const selectedChapter = ref<Chapter | null>(null);

// Progress Tracking
const readParagraphs = ref<Set<number>>(new Set());

// Featured/Recently Viewed Paragraphs
const recentParagraphs = [
  { id: 27, text: 'The desire for God is written in the human heart, because man is created by God and for God...' },
  { id: 1324, text: 'The Eucharist is the source and summit of the Christian life.' },
  { id: 2558, text: 'Great is the mystery of the faith! The Church professes this mystery in the Apostles\' Creed...' }
];

// Modal State
const showModal = ref(false);
const modalContent = ref<Paragraph[]>([]);

const openModal = (paragraphs: Paragraph[]) => {
    modalContent.value = paragraphs;
    showModal.value = true;
};

onMounted(() => {
    // Sync initial search from router query if present (from BottomNav search query)
    if (route.query.q) {
        searchQuery.value = String(route.query.q);
    }

    const saved = localStorage.getItem('catechism-read-progress');
    if (saved) {
        try {
            const parsed = JSON.parse(saved);
            if (Array.isArray(parsed)) {
                readParagraphs.value = new Set(parsed);
            }
        } catch (e) {
            console.error('Failed to load catechism progress', e);
        }
    }
});

const toggleRead = (id: number) => {
    const newSet = new Set(readParagraphs.value);
    if (newSet.has(id)) {
        newSet.delete(id);
    } else {
        newSet.add(id);
    }
    readParagraphs.value = newSet;
    localStorage.setItem('catechism-read-progress', JSON.stringify(Array.from(newSet)));
};

// Methods
const scrollToContent = () => {
  setTimeout(() => {
    const contentArea = document.getElementById('catechism-content');
    if (contentArea) {
      contentArea.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }, 100);
};

const openPart = (partNum: number) => {
  const targetPart = structure.value.catechism.structure.find(p => p.part === partNum);
  if (targetPart) {
    selectedPart.value = targetPart;
    selectedSection.value = null;
    selectedChapter.value = null;
    scrollToContent();
  }
};

const executeUniversalSearch = () => {
  if (universalSearchInput.value.trim()) {
    searchQuery.value = universalSearchInput.value.trim();
    scrollToContent();
  }
};

const jumpToParagraph = (paraId: number) => {
  searchQuery.value = String(paraId);
  scrollToContent();
};

const resetSelection = () => {
  selectedPart.value = null;
  selectedSection.value = null;
  selectedChapter.value = null;
  searchQuery.value = '';
  universalSearchInput.value = '';
};

const togglePart = (part: Part) => {
  if (selectedPart.value?.part === part.part) {
    selectedPart.value = null;
    selectedSection.value = null;
    selectedChapter.value = null;
  } else {
    selectedPart.value = part;
    selectedSection.value = null; 
    selectedChapter.value = null;
    scrollToContent();
  }
};

const toggleSection = (section: Section) => {
    if (selectedSection.value?.section === section.section) {
        selectedSection.value = null;
        selectedChapter.value = null;
    } else {
        selectedSection.value = section;
        selectedChapter.value = null;
        if (!section.chapters || section.chapters.length === 0) {
            scrollToContent();
        }
    }
}

const toggleChapter = (chapter: Chapter) => {
    if (selectedChapter.value?.chapter === chapter.chapter) {
        selectedChapter.value = null;
    } else {
        selectedChapter.value = chapter;
        scrollToContent();
    }
}

// Helper to parse range string "1-10" to object {start, end}
const parseRange = (rangeStr: string): { start: number, end: number } => {
  if (!rangeStr) return { start: 0, end: 0 };
  
  const parts = rangeStr.split('-');
  if (parts.length === 2) {
    const startStr = parts[0];
    const endStr = parts[1];
    if (typeof startStr === 'string' && typeof endStr === 'string') {
        const start = parseInt(startStr.trim());
        const end = parseInt(endStr.trim());
        return { start, end };
    }
  }
  return { start: 0, end: 0 };
};

const getBreadcrumb = (id: number): string => {
  const s = structure.value.catechism.structure;
  for (const part of s) {
    const pRange = parseRange(part.paragraph_range);
    if (id >= pRange.start && id <= pRange.end) {
      const parts = [part.title];
      
      if (part.sections) {
        for (const section of part.sections) {
          const sRange = parseRange(section.paragraph_range);
          if (id >= sRange.start && id <= sRange.end) {
            parts.push(section.title);
            
            if (section.chapters) {
              for (const chapter of section.chapters) {
                const cRange = parseRange(chapter.paragraphs);
                if (id >= cRange.start && id <= cRange.end) {
                  parts.push(chapter.title);
                  break;
                }
              }
            }
            break;
          }
        }
      }
      return parts.join(' / ');
    }
  }
  return '';
};

// Computed
const filteredParagraphs = computed(() => {
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.trim();
    const desiredIds = new Set<number>();
    const parts = query.split(',');
    
    for (const part of parts) {
        const trimmed = part.trim();
        if (!trimmed) continue;
        
        if (trimmed.includes('-')) {
            const rangeParts = trimmed.split('-').map(s => s.trim());
            if (rangeParts.length === 2) {
                const start = parseInt(rangeParts[0]!);
                const end = parseInt(rangeParts[1]!);
                
                if (!isNaN(start) && !isNaN(end) && start <= end) {
                    for (let i = start; i <= end; i++) {
                        desiredIds.add(i);
                    }
                }
            }
        } else {
            const num = parseInt(trimmed);
            if (!isNaN(num)) {
                desiredIds.add(num);
            }
        }
    }
    
    if (desiredIds.size > 0) {
        return paragraphs.value.filter(p => desiredIds.has(p.id));
    }
    
    return [];
  }

  let range = { start: 0, end: 0 };

  if (selectedChapter.value) {
      range = parseRange(selectedChapter.value.paragraphs || '');
  } else if (selectedSection.value) {
      range = parseRange(selectedSection.value.paragraph_range || '');
  } else if (selectedPart.value) {
      range = parseRange(selectedPart.value.paragraph_range || '');
  } else {
      return []; 
  }

  if (range.start === 0 && range.end === 0) return [];

  return paragraphs.value.filter(p => p.id >= range.start && p.id <= range.end);
});

// Watch query from query params
watch(() => route.query.q, (newVal) => {
    if (newVal !== undefined) {
        searchQuery.value = String(newVal);
    }
});
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}
.animate-fade-in-down {
  animation: fadeInDown 0.6s ease-out forwards;
}
.animate-fade-in-up {
  animation: fadeInUp 0.6s ease-out forwards;
}
.delay-100 {
  animation-delay: 0.1s;
  animation-fill-mode: both;
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

