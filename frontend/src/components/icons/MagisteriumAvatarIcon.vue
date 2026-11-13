<template>
  <svg
    viewBox="0 0 150 150"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    class="select-none pointer-events-none"
    aria-hidden="true"
  >
    <defs>
      <filter id="magisterium-avatar-shadow" x="-15%" y="-15%" width="130%" height="135%">
        <feDropShadow dx="0" dy="8" stdDeviation="6" flood-color="#4A3410" flood-opacity="0.34" />
      </filter>
      <linearGradient id="magisterium-medallion-bg" x1="25" y1="25" x2="125" y2="125" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#6B2D33" />
        <stop offset="100%" stop-color="#43181C" />
      </linearGradient>
      <linearGradient id="magisterium-gold-grad" x1="45" y1="24" x2="105" y2="115" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#F8ECD7" />
        <stop offset="50%" stop-color="#DFB878" />
        <stop offset="100%" stop-color="#C89B54" />
      </linearGradient>
      <linearGradient id="magisterium-robe-grad" x1="75" y1="80" x2="75" y2="122" gradientUnits="userSpaceOnUse">
        <stop offset="0%" stop-color="#F8ECD7" />
        <stop offset="100%" stop-color="#E8D2B1" />
      </linearGradient>
    </defs>

    <g filter="url(#magisterium-avatar-shadow)">
      <!-- Outer Medallion Ring (matches Home page CatechismIllustration) -->
      <circle
        cx="75"
        cy="75"
        r="56"
        fill="url(#magisterium-medallion-bg)"
        stroke="#F4E3C1"
        stroke-width="2.75"
        stroke-opacity="0.82"
      />

      <!-- Slowly Rotating Inner Filigree Ring -->
      <circle
        cx="75"
        cy="75"
        r="49.5"
        stroke="#DFB878"
        stroke-width="1.2"
        stroke-opacity="0.55"
        stroke-dasharray="3.5 3.5"
        class="avatar-ring-spin"
        style="transform-origin: 75px 75px;"
      />

      <!-- Radiant Halo Behind Head -->
      <circle
        cx="75"
        cy="60"
        r="24"
        fill="#DFB878"
        fill-opacity="0.14"
        stroke="#DFB878"
        stroke-width="1.5"
        stroke-opacity="0.7"
        class="avatar-halo-pulse"
        style="transform-origin: 75px 60px;"
      />

      <!-- Sacred Cross-Star at Apex -->
      <g stroke="#F8ECD7" stroke-width="2.4" stroke-linecap="round" class="avatar-star-twinkle" style="transform-origin: 75px 33px;">
        <line x1="75" y1="27" x2="75" y2="39" />
        <line x1="69" y1="33" x2="81" y2="33" />
      </g>

      <!-- Twinkling Starlight Motes -->
      <circle cx="43" cy="50" r="1.8" fill="#F8ECD7" class="avatar-mote" />
      <circle cx="107" cy="50" r="1.8" fill="#F8ECD7" class="avatar-mote-delayed" />
      <circle cx="48" cy="84" r="1.4" fill="#DFB878" class="avatar-mote-delayed" />
      <circle cx="102" cy="84" r="1.4" fill="#DFB878" class="avatar-mote" />

      <!-- Scholarly Cowl / Mantle Shoulders -->
      <path
        d="M42 112 C46 92, 58 84, 75 84 C92 84, 104 92, 108 112 C99 120, 88 124, 75 124 C62 124, 51 120, 42 112 Z"
        fill="url(#magisterium-robe-grad)"
        stroke="#DFB878"
        stroke-width="1.4"
      />
      <!-- Burgundy Stole / Scapular Accent -->
      <path
        d="M64 86 L69 123 L81 123 L86 86 C82 84.5, 68 84.5, 64 86 Z"
        fill="#7D3B41"
      />
      <path
        d="M68 85.5 L71 123 M82 85.5 L79 123"
        stroke="#DFB878"
        stroke-width="1.1"
        stroke-opacity="0.75"
      />

      <!-- Illuminated Open Tome at Heart -->
      <path
        d="M59 101 C65 98, 71 99, 75 102 C79 99, 85 98, 91 101 L93 114 C86 111, 80 112, 75 115 C70 112, 64 111, 57 114 Z"
        fill="#FAF4EB"
        stroke="#C89B54"
        stroke-width="1.3"
      />
      <line x1="75" y1="102" x2="75" y2="115" stroke="#B79E84" stroke-width="1.1" />
      <!-- Rubric Lines on Open Tome -->
      <path d="M62 105 C66 104, 69 104.5, 72 106" stroke="#8B2E36" stroke-width="1.3" stroke-linecap="round" />
      <path d="M78 106 C81 104.5, 84 104, 88 105" stroke="#9F8B78" stroke-width="1.2" stroke-linecap="round" />
      <path d="M61.5 109 C65.5 108, 69 108.5, 72 110" stroke="#9F8B78" stroke-width="1.2" stroke-linecap="round" />
      <path d="M78 110 C81 108.5, 84.5 108, 88.5 109" stroke="#9F8B78" stroke-width="1.2" stroke-linecap="round" />

      <!-- Serene Hooded Avatar Head -->
      <path
        d="M75 40 C60 40, 51 50, 51 64 C51 74, 56 82, 64 85.5 C67.5 87, 71 87.8, 75 87.8 C79 87.8, 82.5 87, 86 85.5 C94 82, 99 74, 99 64 C99 50, 90 40, 75 40 Z"
        fill="#E8D2B1"
        stroke="#DFB878"
        stroke-width="1.5"
      />
      <!-- Inner Face / Luminous Visage -->
      <path
        d="M75 45 C64 45, 57 53, 57 64 C57 73, 62 80, 75 82 C88 80, 93 73, 93 64 C93 53, 86 45, 75 45 Z"
        fill="#F8ECD7"
      />

      <!-- Expressive Gently Blinking Eyes -->
      <g class="avatar-blink" style="transform-origin: 75px 62px;">
        <ellipse cx="67.5" cy="62" rx="2.6" ry="3.2" fill="#5C282D" />
        <ellipse cx="82.5" cy="62" rx="2.6" ry="3.2" fill="#5C282D" />
        <circle cx="68.3" cy="60.9" r="0.9" fill="#FAF4EB" />
        <circle cx="83.3" cy="60.9" r="0.9" fill="#FAF4EB" />
      </g>

      <!-- Gentle Serene Smile -->
      <path
        d="M69.5 70.5 C71.2 72.5, 73 73.2, 75 73.2 C77 73.2, 78.8 72.5, 80.5 70.5"
        stroke="#7D3B41"
        stroke-width="1.8"
        stroke-linecap="round"
      />
    </g>
  </svg>
</template>

<style scoped>
@keyframes ringSpin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.avatar-ring-spin {
  animation: ringSpin 22s linear infinite;
}

@keyframes haloPulse {
  0%, 100% {
    opacity: 0.7;
    transform: scale(0.97);
  }
  50% {
    opacity: 1;
    transform: scale(1.04);
  }
}

.avatar-halo-pulse {
  animation: haloPulse 3.2s ease-in-out infinite;
}

@keyframes starTwinkle {
  0%, 100% {
    opacity: 0.8;
    transform: scale(0.95);
  }
  50% {
    opacity: 1;
    transform: scale(1.08);
  }
}

.avatar-star-twinkle {
  animation: starTwinkle 2.8s ease-in-out infinite;
}

@keyframes moteTwinkle {
  0%, 100% {
    opacity: 0.35;
  }
  50% {
    opacity: 1;
  }
}

.avatar-mote {
  animation: moteTwinkle 2.4s ease-in-out infinite;
}

.avatar-mote-delayed {
  animation: moteTwinkle 2.4s ease-in-out 1.2s infinite;
}

@keyframes avatarBlink {
  0%, 44%, 48%, 92%, 96%, 100% {
    transform: scaleY(1);
  }
  46%, 94% {
    transform: scaleY(0.12);
  }
}

.avatar-blink {
  animation: avatarBlink 5.2s ease-in-out infinite;
}
</style>
