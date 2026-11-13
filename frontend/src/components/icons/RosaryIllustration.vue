<template>
  <svg
    viewBox="0 0 150 170"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none pointer-events-none"
    aria-hidden="true"
  >
    <defs>
      <filter id="rosary-shadow" x="-20%" y="-15%" width="140%" height="140%">
        <feDropShadow dx="0" dy="6" stdDeviation="5" flood-color="#220A0D" flood-opacity="0.45" />
      </filter>
      <radialGradient id="bead-gold" cx="35%" cy="30%" r="65%">
        <stop offset="0%" stop-color="#FFF7E8" />
        <stop offset="55%" stop-color="#E7CCA0" />
        <stop offset="100%" stop-color="#B98D52" />
      </radialGradient>
      <radialGradient id="bead-pater" cx="35%" cy="30%" r="65%">
        <stop offset="0%" stop-color="#FFFDF8" />
        <stop offset="50%" stop-color="#F2DEC0" />
        <stop offset="100%" stop-color="#C99A5B" />
      </radialGradient>
      <linearGradient id="cross-gold" x1="65" y1="125" x2="95" y2="162" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#F6E3C3" />
        <stop offset="50%" stop-color="#DFB87B" />
        <stop offset="100%" stop-color="#AD7E3D" />
      </linearGradient>
    </defs>

    <g filter="url(#rosary-shadow)">
      <!-- Subtle chain connecting the loop -->
      <path
        d="M75 18 C42 18, 22 40, 26 68 C29 90, 52 104, 75 106 C98 104, 121 90, 124 68 C128 40, 108 18, 75 18 Z"
        stroke="#D8B37C"
        stroke-opacity="0.55"
        stroke-width="1.6"
        stroke-dasharray="2 2"
      />
      <!-- Hanging chain -->
      <line x1="75" y1="106" x2="75" y2="132" stroke="#D8B37C" stroke-width="1.8" stroke-opacity="0.7" />

      <!-- Loop Ave Beads (left & right symmetrical arrangement) -->
      <!-- Top crown Pater bead -->
      <circle cx="75" cy="18" r="5" fill="url(#bead-pater)" />

      <!-- Left arc beads -->
      <circle cx="62" cy="19.5" r="3.4" fill="url(#bead-gold)" />
      <circle cx="51" cy="23" r="3.4" fill="url(#bead-gold)" />
      <circle cx="41" cy="29" r="3.4" fill="url(#bead-gold)" />
      <circle cx="33" cy="38" r="3.4" fill="url(#bead-gold)" />
      <circle cx="28" cy="48" r="4.6" fill="url(#bead-pater)" />
      <circle cx="25.5" cy="59" r="3.4" fill="url(#bead-gold)" />
      <circle cx="26" cy="70" r="3.4" fill="url(#bead-gold)" />
      <circle cx="29.5" cy="80" r="3.4" fill="url(#bead-gold)" />
      <circle cx="36" cy="89" r="3.4" fill="url(#bead-gold)" />
      <circle cx="45" cy="96" r="4.6" fill="url(#bead-pater)" />
      <circle cx="56" cy="101" r="3.4" fill="url(#bead-gold)" />
      <circle cx="66" cy="104" r="3.4" fill="url(#bead-gold)" />

      <!-- Right arc beads -->
      <circle cx="88" cy="19.5" r="3.4" fill="url(#bead-gold)" />
      <circle cx="99" cy="23" r="3.4" fill="url(#bead-gold)" />
      <circle cx="109" cy="29" r="3.4" fill="url(#bead-gold)" />
      <circle cx="117" cy="38" r="3.4" fill="url(#bead-gold)" />
      <circle cx="122" cy="48" r="4.6" fill="url(#bead-pater)" />
      <circle cx="124.5" cy="59" r="3.4" fill="url(#bead-gold)" />
      <circle cx="124" cy="70" r="3.4" fill="url(#bead-gold)" />
      <circle cx="120.5" cy="80" r="3.4" fill="url(#bead-gold)" />
      <circle cx="114" cy="89" r="3.4" fill="url(#bead-gold)" />
      <circle cx="105" cy="96" r="4.6" fill="url(#bead-pater)" />
      <circle cx="94" cy="101" r="3.4" fill="url(#bead-gold)" />
      <circle cx="84" cy="104" r="3.4" fill="url(#bead-gold)" />

      <!-- Centerpiece Marian Medal -->
      <path
        d="M75 98 L81.5 106 L75 114 L68.5 106 Z"
        fill="url(#cross-gold)"
        stroke="#FFF6E5"
        stroke-width="0.8"
      />
      <circle cx="75" cy="106" r="2.2" fill="#5C282D" />

      <!-- Tail Pater & 3 Ave Beads -->
      <circle cx="75" cy="119" r="3.5" fill="url(#bead-gold)" />
      <circle cx="75" cy="126.5" r="3.2" fill="url(#bead-gold)" />

      <!-- Crucifix -->
      <g>
        <!-- Vertical beam -->
        <rect x="72.5" y="131" width="5" height="31" rx="2" fill="url(#cross-gold)" />
        <!-- Horizontal beam -->
        <rect x="63.5" y="139" width="23" height="4.8" rx="2" fill="url(#cross-gold)" />
        <!-- Inner highlight -->
        <line x1="75" y1="133" x2="75" y2="159" stroke="#FFF8EB" stroke-opacity="0.6" stroke-width="1" stroke-linecap="round" />
        <line x1="65.5" y1="141.4" x2="84.5" y2="141.4" stroke="#FFF8EB" stroke-opacity="0.6" stroke-width="1" stroke-linecap="round" />
      </g>
    </g>
  </svg>
</template>
