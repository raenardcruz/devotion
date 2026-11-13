<template>
  <svg
    viewBox="0 0 160 130"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none pointer-events-none"
    aria-hidden="true"
  >
    <defs>
      <filter id="missal-shadow" x="-15%" y="-15%" width="130%" height="135%">
        <feDropShadow dx="0" dy="8" stdDeviation="6" flood-color="#3E281C" flood-opacity="0.28" />
      </filter>
      <linearGradient id="missal-cover" x1="15" y1="25" x2="145" y2="115" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#7B333A" />
        <stop offset="100%" stop-color="#4A1B20" />
      </linearGradient>
      <linearGradient id="missal-page-l" x1="25" y1="28" x2="80" y2="100" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#FBF6EE" />
        <stop offset="100%" stop-color="#E7D9C8" />
      </linearGradient>
      <linearGradient id="missal-page-r" x1="80" y1="28" x2="135" y2="100" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#E6D7C5" />
        <stop offset="25%" stop-color="#FAF3E8" />
        <stop offset="100%" stop-color="#FFFDF9" />
      </linearGradient>
    </defs>

    <g filter="url(#missal-shadow)">
      <!-- Burgundy Leather Cover -->
      <path
        d="M16 34 C16 30, 19 28, 23 29 L76 39 C79 40, 81 40, 84 39 L137 29 C141 28, 144 30, 144 34 L149 98 C149 102, 146 105, 142 104 L84 99 C81 99, 79 99, 76 99 L18 104 C14 105, 11 102, 11 98 Z"
        fill="url(#missal-cover)"
        stroke="#C89B54"
        stroke-opacity="0.45"
        stroke-width="1.2"
      />

      <!-- Gold Bookmark Ribbon -->
      <path d="M77 95 L83 95 L84 120 L80 115 L76 120 Z" fill="#C99951" />

      <!-- Left Page -->
      <path
        d="M23 24 C42 20, 62 22, 80 33 L80 94 C62 86, 42 88, 21 94 Z"
        fill="url(#missal-page-l)"
      />
      <!-- Right Page -->
      <path
        d="M137 24 C118 20, 98 22, 80 33 L80 94 C98 86, 118 88, 139 94 Z"
        fill="url(#missal-page-r)"
      />

      <!-- Center Spine -->
      <line x1="80" y1="33" x2="80" y2="94" stroke="#B69B80" stroke-width="1.2" stroke-opacity="0.7" />

      <!-- Left Page Lines -->
      <g stroke="#8F7A68" stroke-width="1.8" stroke-linecap="round" stroke-opacity="0.6">
        <path d="M31 39 C46 37, 58 38, 70 43" stroke="#7B333A" stroke-opacity="0.85" />
        <path d="M30 48 C46 46, 58 47, 70 52" />
        <path d="M29 56 C45 54, 58 55, 70 60" />
        <path d="M29 64 C45 62, 58 63, 70 68" />
        <path d="M28 72 C44 70, 56 71, 68 75" />
      </g>

      <!-- Right Page Liturgical Cross -->
      <g stroke="#7B333A" stroke-width="3.2" stroke-linecap="round">
        <line x1="109" y1="39" x2="109" y2="75" />
        <line x1="97" y1="52" x2="121" y2="52" />
      </g>
      <g stroke="#C99951" stroke-width="1.2" stroke-linecap="round" stroke-opacity="0.8">
        <line x1="109" y1="41" x2="109" y2="73" />
        <line x1="99" y1="52" x2="119" y2="52" />
      </g>
    </g>
  </svg>
</template>
