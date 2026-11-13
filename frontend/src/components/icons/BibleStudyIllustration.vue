<template>
  <svg
    viewBox="0 0 150 150"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none pointer-events-none"
    aria-hidden="true"
  >
    <defs>
      <filter id="bible-shadow" x="-18%" y="-12%" width="136%" height="136%">
        <feDropShadow dx="2" dy="10" stdDeviation="7" flood-color="#2E1308" flood-opacity="0.48" />
      </filter>
      <linearGradient id="bible-cover" x1="30" y1="20" x2="115" y2="120" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#3E2A25" />
        <stop offset="100%" stop-color="#241713" />
      </linearGradient>
      <linearGradient id="bible-spine" x1="25" y1="40" x2="55" y2="130" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#2C1C18" />
        <stop offset="100%" stop-color="#190F0D" />
      </linearGradient>
    </defs>

    <g filter="url(#bible-shadow)">
      <!-- Bookmark Ribbon -->
      <path d="M68 112 L75 110 L78 139 L72 133 L66 138 Z" fill="#E8C684" />

      <!-- Back Cover & Page Edges (3D isometric bottom-right thickness) -->
      <path
        d="M52 126 L124 98 L124 111 L52 138 Z"
        fill="#241713"
      />
      <!-- Cream Page Block -->
      <path
        d="M51 116 L121 89 L121 101 L51 128 Z"
        fill="#F4E9D8"
      />
      <line x1="53" y1="120" x2="120" y2="94" stroke="#D5C2AA" stroke-width="1" />
      <line x1="53" y1="124" x2="120" y2="98" stroke="#D5C2AA" stroke-width="1" />

      <!-- Spine (left thickness) -->
      <path
        d="M27 52 C24 54, 24 58, 25 62 L39 124 C40 129, 45 132, 52 129 L52 116 L36 47 Z"
        fill="url(#bible-spine)"
      />

      <!-- Front Leather Cover -->
      <path
        d="M35 45 L103 20 C107 18, 111 20, 112 24 L126 88 C127 92, 125 96, 121 97 L53 123 C49 124, 46 122, 45 118 L31 54 C30 50, 32 46, 35 45 Z"
        fill="url(#bible-cover)"
        stroke="#C89B54"
        stroke-opacity="0.35"
        stroke-width="1"
      />

      <!-- Spine Ribs -->
      <path d="M31 62 L37 59" stroke="#C89B54" stroke-opacity="0.5" stroke-width="1.5" stroke-linecap="round" />
      <path d="M35 80 L41 77" stroke="#C89B54" stroke-opacity="0.5" stroke-width="1.5" stroke-linecap="round" />
      <path d="M39 98 L45 95" stroke="#C89B54" stroke-opacity="0.5" stroke-width="1.5" stroke-linecap="round" />

      <!-- Embossed Gold Cross on Front Cover -->
      <g stroke="#DFB878" stroke-width="3.4" stroke-linecap="round">
        <line x1="72" y1="45" x2="84" y2="91" />
        <line x1="63" y1="63" x2="86" y2="55" />
      </g>
      <g stroke="#FFF1D6" stroke-opacity="0.45" stroke-width="1" stroke-linecap="round">
        <line x1="72" y1="46" x2="83.5" y2="90" />
        <line x1="64" y1="62.5" x2="85" y2="55" />
      </g>
    </g>
  </svg>
</template>
