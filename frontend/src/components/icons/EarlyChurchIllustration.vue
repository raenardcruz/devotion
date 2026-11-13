<template>
  <svg
    viewBox="0 0 150 150"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none pointer-events-none"
    aria-hidden="true"
  >
    <defs>
      <filter id="scroll-shadow" x="-18%" y="-12%" width="136%" height="136%">
        <feDropShadow dx="1" dy="8" stdDeviation="6" flood-color="#1D2619" flood-opacity="0.42" />
      </filter>
      <linearGradient id="parchment-body" x1="35" y1="25" x2="115" y2="120" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#FCF6EB" />
        <stop offset="55%" stop-color="#F2E4CE" />
        <stop offset="100%" stop-color="#DEC9AB" />
      </linearGradient>
      <linearGradient id="scroll-roll" x1="30" y1="0" x2="120" y2="0" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#D7BE9D" />
        <stop offset="30%" stop-color="#F8EEDC" />
        <stop offset="75%" stop-color="#E7D3B7" />
        <stop offset="100%" stop-color="#C8AC87" />
      </linearGradient>
      <linearGradient id="wood-finial" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="#DFB878" />
        <stop offset="100%" stop-color="#9A7138" />
      </linearGradient>
      <radialGradient id="wax-seal" cx="38%" cy="35%" r="65%">
        <stop offset="0%" stop-color="#9B4149" />
        <stop offset="65%" stop-color="#6B252C" />
        <stop offset="100%" stop-color="#48151A" />
      </radialGradient>
    </defs>

    <g filter="url(#scroll-shadow)" transform="rotate(-6 75 75)">
      <!-- Wooden Roller Rod Ends (Top) -->
      <rect x="26" y="25" width="98" height="6" rx="3" fill="url(#wood-finial)" />
      <circle cx="25" cy="28" r="4.5" fill="#C99A5B" />
      <circle cx="125" cy="28" r="4.5" fill="#B58546" />

      <!-- Wooden Roller Rod Ends (Bottom) -->
      <rect x="26" y="113" width="98" height="6" rx="3" fill="url(#wood-finial)" />
      <circle cx="25" cy="116" r="4.5" fill="#C99A5B" />
      <circle cx="125" cy="116" r="4.5" fill="#B58546" />

      <!-- Main Unrolled Parchment Sheet -->
      <path
        d="M36 28 C38 56, 35 88, 37 116 L113 116 C111 88, 114 56, 112 28 Z"
        fill="url(#parchment-body)"
        stroke="#C9AF8B"
        stroke-width="0.9"
      />

      <!-- Subtle Fold Shadows under Top & Bottom Rolls -->
      <path d="M36 35 L112 35 L112 40 C86 38, 62 38, 36 40 Z" fill="#B89B75" fill-opacity="0.28" />
      <path d="M37 104 C62 106, 86 106, 113 104 L113 110 L37 110 Z" fill="#B89B75" fill-opacity="0.24" />

      <!-- Top Rolled Parchment Cylinder -->
      <rect
        x="33"
        y="20"
        width="84"
        height="16"
        rx="7.5"
        fill="url(#scroll-roll)"
        stroke="#C2A57F"
        stroke-width="1"
      />
      <ellipse cx="36.5" cy="28" rx="3" ry="7.2" fill="#B99971" />

      <!-- Bottom Rolled Parchment Cylinder -->
      <rect
        x="33"
        y="108"
        width="84"
        height="16"
        rx="7.5"
        fill="url(#scroll-roll)"
        stroke="#C2A57F"
        stroke-width="1"
      />
      <ellipse cx="36.5" cy="116" rx="3" ry="7.2" fill="#B99971" />

      <!-- Manuscript Rubric Header & Script Lines -->
      <path d="M47 47 C60 45.5, 72 45.5, 84 47" stroke="#7B333A" stroke-width="2.2" stroke-linecap="round" stroke-opacity="0.88" />
      <g stroke="#7C6956" stroke-width="1.7" stroke-linecap="round" stroke-opacity="0.62">
        <path d="M47 56 C66 54.5, 82 54.5, 101 56" />
        <path d="M47 64 C66 62.5, 82 62.5, 101 64" />
        <path d="M47 72 C64 70.8, 78 70.8, 95 72" />
        <path d="M47 80 C62 79, 73 79, 85 80" />
        <path d="M47 88 C59 87.2, 68 87.2, 77 88" />
      </g>

      <!-- Crimson Wax Seal with Gold Cross -->
      <circle cx="96" cy="91" r="11.5" fill="url(#wax-seal)" stroke="#DFB878" stroke-width="1.1" stroke-opacity="0.7" />
      <circle cx="96" cy="91" r="8.2" stroke="#F4E3C1" stroke-width="0.7" stroke-opacity="0.35" />
      <g stroke="#F5E5C8" stroke-width="1.6" stroke-linecap="round">
        <line x1="96" y1="86" x2="96" y2="96" />
        <line x1="92" y1="90" x2="100" y2="90" />
      </g>
    </g>
  </svg>
</template>
