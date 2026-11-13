<template>
  <svg
    viewBox="0 0 220 175"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none pointer-events-none"
    aria-hidden="true"
  >
    <defs>
      <filter id="lectionary-shadow" x="-15%" y="-10%" width="130%" height="135%">
        <feDropShadow dx="0" dy="10" stdDeviation="8" flood-color="#1E0A0C" flood-opacity="0.42" />
      </filter>
      <linearGradient id="cover-grad" x1="15" y1="30" x2="205" y2="155" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#4A2125" />
        <stop offset="50%" stop-color="#361619" />
        <stop offset="100%" stop-color="#270E11" />
      </linearGradient>
      <linearGradient id="left-page" x1="25" y1="40" x2="110" y2="130" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#FAF4EB" />
        <stop offset="80%" stop-color="#F2E7D8" />
        <stop offset="100%" stop-color="#DECDBA" />
      </linearGradient>
      <linearGradient id="right-page" x1="110" y1="40" x2="195" y2="130" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#E1D0BE" />
        <stop offset="18%" stop-color="#F4EADB" />
        <stop offset="100%" stop-color="#FBF7F1" />
      </linearGradient>
      <linearGradient id="ribbon-grad" x1="108" y1="120" x2="108" y2="168" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#7B2D35" />
        <stop offset="100%" stop-color="#C89B54" />
      </linearGradient>
    </defs>

    <g filter="url(#lectionary-shadow)">
      <!-- Outer leather cover -->
      <path
        d="M22 44 C22 40, 26 37, 30 38 L104 52 C108 53, 112 53, 116 52 L190 38 C194 37, 198 40, 198 44 L204 134 C204 138, 200 142, 196 141 L118 135 C112 135, 108 135, 102 135 L24 141 C20 142, 16 138, 16 134 Z"
        fill="url(#cover-grad)"
        stroke="#C89B54"
        stroke-opacity="0.35"
        stroke-width="1.2"
      />

      <!-- Page thickness block (bottom edge) -->
      <path
        d="M25 131 L104 127 C108 127, 112 127, 116 127 L195 131 L192 137 L116 133 C112 133, 108 133, 104 133 L28 137 Z"
        fill="#D6C4B0"
      />
      <path
        d="M27 127 L104 123 C108 123, 112 123, 116 123 L193 127 L193 131 L116 127 C112 127, 108 127, 104 127 L27 131 Z"
        fill="#E7D8C7"
      />

      <!-- Hanging liturgical bookmark ribbon -->
      <path
        d="M107 124 L113 124 L114 164 L109.5 158 L105 164 Z"
        fill="url(#ribbon-grad)"
      />
      <path d="M107 124 L109 124 L109 158 L105 164 Z" fill="#5B1F26" fill-opacity="0.35" />

      <!-- Left page -->
      <path
        d="M32 34 C56 28, 86 30, 110 44 L110 126 C86 116, 56 118, 28 126 Z"
        fill="url(#left-page)"
      />
      <!-- Right page -->
      <path
        d="M188 34 C164 28, 134 30, 110 44 L110 126 C134 116, 164 118, 192 126 Z"
        fill="url(#right-page)"
      />

      <!-- Subtle page curl lines -->
      <path
        d="M35 37 C58 31, 86 33, 108 46"
        stroke="#FFFFFF"
        stroke-opacity="0.7"
        stroke-width="1"
      />
      <path
        d="M185 37 C162 31, 134 33, 112 46"
        stroke="#FFFFFF"
        stroke-opacity="0.7"
        stroke-width="1"
      />

      <!-- Center spine crease -->
      <line x1="110" y1="44" x2="110" y2="126" stroke="#B79E84" stroke-width="1.4" stroke-opacity="0.65" />

      <!-- Left page rubric header & text lines -->
      <path d="M44 48 C62 45, 82 47, 98 54" stroke="#8B2E36" stroke-width="2.2" stroke-linecap="round" stroke-opacity="0.85" />
      <g stroke="#9F8B78" stroke-width="1.6" stroke-linecap="round" stroke-opacity="0.55">
        <path d="M42 58 C62 55, 82 57, 98 63" />
        <path d="M41 66 C61 63, 82 65, 98 71" />
        <path d="M40 74 C60 71, 82 73, 98 79" />
        <path d="M39 82 C59 79, 82 81, 98 87" />
        <path d="M39 90 C58 87, 81 89, 98 95" />
        <path d="M38 98 C57 95, 76 97, 90 101" />
        <path d="M38 106 C56 103, 78 105, 96 110" />
      </g>

      <!-- Right page rubric header & text lines -->
      <path d="M122 54 C138 47, 158 45, 176 48" stroke="#8B2E36" stroke-width="2.2" stroke-linecap="round" stroke-opacity="0.85" />
      <g stroke="#9F8B78" stroke-width="1.6" stroke-linecap="round" stroke-opacity="0.55">
        <path d="M122 63 C138 57, 158 55, 178 58" />
        <path d="M122 71 C138 65, 159 63, 179 66" />
        <path d="M122 79 C138 73, 160 71, 180 74" />
        <path d="M122 87 C138 81, 161 79, 181 82" />
        <path d="M122 95 C139 89, 162 87, 181 90" />
        <path d="M122 103 C140 97, 160 95, 176 97" />
        <path d="M124 110 C142 105, 163 103, 182 106" />
      </g>
    </g>
  </svg>
</template>
