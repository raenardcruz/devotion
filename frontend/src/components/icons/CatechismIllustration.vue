<template>
  <svg
    viewBox="0 0 150 150"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none pointer-events-none"
    aria-hidden="true"
  >
    <defs>
      <filter id="basilica-shadow" x="-15%" y="-15%" width="130%" height="135%">
        <feDropShadow dx="0" dy="8" stdDeviation="6" flood-color="#4A3410" flood-opacity="0.32" />
      </filter>
      <linearGradient id="medallion-bg" x1="25" y1="25" x2="125" y2="125" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#6B2D33" />
        <stop offset="100%" stop-color="#43181C" />
      </linearGradient>
    </defs>

    <g filter="url(#basilica-shadow)">
      <!-- Outer Medallion Ring -->
      <circle cx="75" cy="76" r="52" fill="url(#medallion-bg)" stroke="#F4E3C1" stroke-width="2.5" stroke-opacity="0.7" />
      <circle cx="75" cy="76" r="46" stroke="#DFB878" stroke-width="1" stroke-opacity="0.45" stroke-dasharray="3 3" />

      <!-- Cross at the Apex of the Dome -->
      <g stroke="#F8ECD7" stroke-width="2.6" stroke-linecap="round">
        <line x1="75" y1="29" x2="75" y2="41" />
        <line x1="70.5" y1="33.5" x2="79.5" y2="33.5" />
      </g>

      <!-- Cupola / Lantern -->
      <rect x="70.5" y="40" width="9" height="6" rx="1.5" fill="#F4E5CE" />

      <!-- Main Basilica Dome -->
      <path
        d="M55 64 C55 49, 64 44, 75 44 C86 44, 95 49, 95 64 Z"
        fill="#F4E5CE"
      />
      <!-- Dome Ribs -->
      <path d="M65 47 C63 53, 62 58, 62 64" stroke="#6B2D33" stroke-width="1.2" stroke-opacity="0.45" />
      <path d="M75 44 L75 64" stroke="#6B2D33" stroke-width="1.2" stroke-opacity="0.45" />
      <path d="M85 47 C87 53, 88 58, 88 64" stroke="#6B2D33" stroke-width="1.2" stroke-opacity="0.45" />

      <!-- Drum / Entablature -->
      <rect x="52" y="64" width="46" height="6" rx="1.5" fill="#DFB878" />

      <!-- Side Wings -->
      <path d="M36 84 L52 74 L52 104 L36 104 Z" fill="#E8D2B1" />
      <path d="M114 84 L98 74 L98 104 L114 104 Z" fill="#E8D2B1" />

      <!-- Central Facade & Pediment -->
      <polygon points="50,74 75,62 100,74" fill="#F8ECD7" />
      <rect x="52" y="74" width="46" height="30" fill="#F4E5CE" />

      <!-- Arched Doors / Columns -->
      <path d="M68 104 L68 86 C68 82, 71 79, 75 79 C79 79, 82 82, 82 86 L82 104 Z" fill="#5C282D" />
      <path d="M57 104 L57 89 C57 86.5, 58.8 85, 61 85 C63.2 85, 65 86.5, 65 89 L65 104 Z" fill="#7D3B41" />
      <path d="M85 104 L85 89 C85 86.5, 86.8 85, 89 85 C91.2 85, 93 86.5, 93 89 L93 104 Z" fill="#7D3B41" />

      <!-- Side Wing Windows -->
      <rect x="41" y="88" width="5" height="9" rx="2.5" fill="#5C282D" />
      <rect x="104" y="88" width="5" height="9" rx="2.5" fill="#5C282D" />

      <!-- Steps / Base -->
      <rect x="33" y="104" width="84" height="4" rx="2" fill="#DFB878" />
    </g>
  </svg>
</template>
