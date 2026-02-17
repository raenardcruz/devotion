<template>
  <div class="min-h-screen bg-stone-50 text-stone-900 font-sans selection:bg-amber-100 flex flex-col overflow-hidden">
    <!-- Header -->
    <header class="pt-20 pb-12 text-center px-4">
      <h1 class="text-5xl md:text-7xl font-serif text-stone-800 tracking-tight mb-6 animate-fade-in-down">
        Devotion
      </h1>
      <p class="text-stone-500 text-sm md:text-base uppercase tracking-[0.2em] font-medium animate-fade-in-up delay-100">
        Daily Catholic Prayers & Meditations
      </p>
    </header>

    <!-- Main Content -->
    <main class="flex-grow flex items-center justify-center p-4">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-7xl w-full">
        
        <!-- Daily Mass Readings Card -->
        <router-link to="/daily-readings" class="group relative overflow-hidden rounded-[2rem] shadow-xl hover:shadow-2xl transition-all duration-700 hover:-translate-y-2 h-[28rem] animate-fade-in-up delay-150">
          <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1497621122273-f5cfb6065c56?q=80&w=1374&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')] bg-cover bg-center transition-transform duration-1000 group-hover:scale-110"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-stone-900 via-stone-900/40 to-transparent opacity-80 group-hover:opacity-70 transition-opacity duration-700"></div>
          
          <div class="relative h-full flex flex-col justify-end p-10 text-white">
            <div class="mb-6 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-500 ease-out">
              <span class="text-amber-200 text-xs font-bold uppercase tracking-[0.2em] mb-3 block opacity-80">Liturgy</span>
              <h2 class="text-4xl font-serif italic mb-3">Daily Readings</h2>
              <p class="text-stone-200 text-sm leading-relaxed opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-100 max-w-sm">
                Read and meditate on today's Mass readings: First Reading, Psalm, and Gospel.
              </p>
            </div>
            
            <div class="flex items-center text-amber-200 text-sm font-medium tracking-wide opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-200">
              <span class="border-b border-transparent group-hover:border-amber-200 transition-all duration-300">Read Now</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-3 group-hover:translate-x-2 transition-transform duration-300">
                <path d="M5 12h14"></path>
                <path d="m12 5 7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </router-link>

        <!-- Rosary Card -->
        <router-link to="/rosary" class="group relative overflow-hidden rounded-[2rem] shadow-xl hover:shadow-2xl transition-all duration-700 hover:-translate-y-2 h-[28rem] animate-fade-in-up delay-200">
          <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1554089638-52a6a0c7b2a1?q=80&w=772&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')] bg-cover bg-center transition-transform duration-1000 group-hover:scale-110"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-stone-900 via-stone-900/40 to-transparent opacity-80 group-hover:opacity-70 transition-opacity duration-700"></div>
          
          <div class="relative h-full flex flex-col justify-end p-10 text-white">
            <div class="mb-6 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-500 ease-out">
              <span class="text-amber-200 text-xs font-bold uppercase tracking-[0.2em] mb-3 block opacity-80">Meditation</span>
              <h2 class="text-4xl font-serif italic mb-3">Scriptural Rosary</h2>
              <p class="text-stone-200 text-sm leading-relaxed opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-100 max-w-sm">
                Meditate on the mysteries of Christ's life with scriptural verses for each Hail Mary.
              </p>
            </div>
            
            <div class="flex items-center text-amber-200 text-sm font-medium tracking-wide opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-200">
              <span class="border-b border-transparent group-hover:border-amber-200 transition-all duration-300">Start Prayer</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-3 group-hover:translate-x-2 transition-transform duration-300">
                <path d="M5 12h14"></path>
                <path d="m12 5 7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </router-link>

        <!-- Divine Mercy Card -->
        <router-link to="/divine-mercy" class="group relative overflow-hidden rounded-[2rem] shadow-xl hover:shadow-2xl transition-all duration-700 hover:-translate-y-2 h-[28rem] animate-fade-in-up delay-300">
          <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1696261803619-9cae57748d90?q=80&w=774&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')] bg-cover bg-center transition-transform duration-1000 group-hover:scale-110"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-rose-950 via-rose-900/40 to-transparent opacity-80 group-hover:opacity-70 transition-opacity duration-700"></div>
          
          <div class="relative h-full flex flex-col justify-end p-10 text-white">
            <div class="mb-6 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-500 ease-out">
              <span class="text-rose-200 text-xs font-bold uppercase tracking-[0.2em] mb-3 block opacity-80">Chaplet</span>
              <h2 class="text-4xl font-serif italic mb-3">Divine Mercy</h2>
              <p class="text-stone-200 text-sm leading-relaxed opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-100 max-w-sm">
                Pray the Chaplet of Divine Mercy for grace and mercy upon the whole world.
              </p>
            </div>
            
            <div class="flex items-center text-rose-200 text-sm font-medium tracking-wide opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-200">
              <span class="border-b border-transparent group-hover:border-rose-200 transition-all duration-300">Start Prayer</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-3 group-hover:translate-x-2 transition-transform duration-300">
                <path d="M5 12h14"></path>
                <path d="m12 5 7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </router-link>

        <!-- Catechism Card -->
        <router-link to="/catechism" class="group relative overflow-hidden rounded-[2rem] shadow-xl hover:shadow-2xl transition-all duration-700 hover:-translate-y-2 h-[28rem] animate-fade-in-up delay-400 md:col-span-2 lg:col-span-1">
          <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1633963792796-8e35bb51675a?q=80&w=687&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')] bg-cover bg-center transition-transform duration-1000 group-hover:scale-110"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-emerald-950 via-emerald-900/40 to-transparent opacity-80 group-hover:opacity-70 transition-opacity duration-700"></div>
          
          <div class="relative h-full flex flex-col justify-end p-10 text-white">
            <div class="mb-6 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-500 ease-out">
              <span class="text-emerald-200 text-xs font-bold uppercase tracking-[0.2em] mb-3 block opacity-80">Teaching</span>
              <h2 class="text-4xl font-serif italic mb-3">Catechism</h2>
              <p class="text-stone-200 text-sm leading-relaxed opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-100 max-w-sm">
                Explore the Catechism of the Catholic Church by paragraph or structure.
              </p>
            </div>
            
            <div class="flex items-center text-emerald-200 text-sm font-medium tracking-wide opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-200">
              <span class="border-b border-transparent group-hover:border-emerald-200 transition-all duration-300">Read Now</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-3 group-hover:translate-x-2 transition-transform duration-300">
                <path d="M5 12h14"></path>
                <path d="m12 5 7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </router-link>

        <!-- Prayers Gallery Card -->
        <router-link to="/prayers" class="group relative overflow-hidden rounded-[2rem] shadow-xl hover:shadow-2xl transition-all duration-700 hover:-translate-y-2 h-[28rem] animate-fade-in-up delay-500 md:col-span-2 lg:col-span-1">
          <div class="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1437603568260-1950d3ca6eab?q=80&w=1170&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')] bg-cover bg-center transition-transform duration-1000 group-hover:scale-110"></div>
          <div class="absolute inset-0 bg-gradient-to-t from-cyan-950 via-cyan-900/40 to-transparent opacity-80 group-hover:opacity-70 transition-opacity duration-700"></div>
          
          <div class="relative h-full flex flex-col justify-end p-10 text-white">
            <div class="mb-6 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-500 ease-out">
              <span class="text-cyan-200 text-xs font-bold uppercase tracking-[0.2em] mb-3 block opacity-80">Collection</span>
              <h2 class="text-4xl font-serif italic mb-3">Traditional Prayers</h2>
              <p class="text-stone-200 text-sm leading-relaxed opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-100 max-w-sm">
                Explore a gallery of essential Catholic prayers for every occasion and need.
              </p>
            </div>
            
            <div class="flex items-center text-cyan-200 text-sm font-medium tracking-wide opacity-0 group-hover:opacity-100 transition-opacity duration-500 delay-200">
              <span class="border-b border-transparent group-hover:border-cyan-200 transition-all duration-300">Open Gallery</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-3 group-hover:translate-x-2 transition-transform duration-300">
                <path d="M5 12h14"></path>
                <path d="m12 5 7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </router-link>
      </div>
    </main>
  </div>
</template>