---
name: apple-design
description: Complete Apple-inspired UI & motion design system. Combines HIG principles, fluid motion physics (springs, velocity handoff, momentum projection, interruptibility), glassmorphism, dark mode, typography (optical sizing, clamp responsive scales), 8px spacing grid, and reusable UI component standards (cards, navbars, hero sections, buttons). Use when designing, building, or auditing Apple-style interfaces, responsive layouts, animations, and components.
---

# Apple Design System & Fluid Interface Guide

Combine Apple's WWDC design philosophy (*Designing Fluid Interfaces*, *Principles of Great Design*, *UI Typography*) with a ready-to-use CSS design system, component patterns, and accessibility guidelines.

---

## 1. Core Principles & Philosophy

### The Four Human Needs
Apple frames interface design as serving four core human needs:
- **Safety & Predictability**: Clear state feedback, non-destructive defaults, and reliable navigation.
- **Understanding**: Clear layout hierarchy, familiar metaphors, and spatial consistency.
- **Achievement**: Direct manipulation, low latency, and rapid task completion.
- **Joy**: Fluid spring motion, refined typography, translucent depth, and tactile polish.

### The Eight Design Principles (WWDC 2026)
1. **Purpose**: Make with intention; decide what *not* to build. Every feature asks for the user's time and attention.
2. **Agency**: Keep users in control with clear paths, simple undo, and minimal confirmation friction.
3. **Responsibility**: Protect user interest (privacy, safety, transparent disclaimers).
4. **Familiarity**: Build on known physical and UI metaphors; maintain spatial consistency (elements enter/dismiss along symmetric paths).
5. **Flexibility**: Adapt gracefully across screen dimensions, input modes (touch vs mouse), and accessibility settings.
6. **Simplicity**: Strip unnecessary clutter; prioritize clarity and content hierarchy over superficial minimalism.
7. **Craft**: Uncompromising attention to spacing, typography, alignment, contrast, and smooth transitions.
8. **Delight**: The natural result of mastering the other seven principles.

---

## 2. Fluid Motion & Physics System

### Key Motion Rules
- **Instant Response**: Respond on `pointerdown` / press, not release. Eliminate input latency.
- **1:1 Direct Manipulation**: Dragged elements track pointer location 1:1, respecting the initial grab offset.
- **Interruptibility**: All animations must be interruptible. Read presentation values on interrupt instead of target values.
- **Velocity Handoff**: Pass gesture release velocity into spring animations to avoid sudden speed discontinuities.
- **Momentum Projection**: Project endpoints based on velocity (`current + (v / 1000) * d / (1 - d)` with `d ≈ 0.998`).
- **Rubber-banding**: Apply progressive resistance at drag boundaries instead of hard stops.

### Spring Mechanics & Presets
Apple frames spring physics using **Damping Ratio** (overshoot control) and **Response Time** (settle speed in seconds).

| Interaction | Damping Ratio | Response Time | Bounce | Usage |
| :--- | :---: | :---: | :---: | :--- |
| **Default UI Transition** | `1.0` (Critically Damped) | `0.3s - 0.4s` | `0` | Standard cards, menus, page fades |
| **Reposition / Drag Release** | `0.8 - 1.0` | `0.4s` | Minimal | Sheet drags, floating panels |
| **Momentum / Flick** | `0.8` | `0.3s` | Subtle | Flicked cards, pull-to-refresh |

```js
// Web Spring Example (Motion / Framer Motion)
import { animate } from 'motion';

// Critically damped default (no overshoot)
animate(element, { y: 0 }, { type: 'spring', bounce: 0, duration: 0.4 });

// Momentum flick (subtle bounce from velocity handoff)
animate(element, { y: target }, { type: 'spring', bounce: 0.2, duration: 0.4, velocity: releaseVelocity });
```

---

## 3. Color System & Dark Mode

### CSS Design Tokens
```css
/* Light Mode */
:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f7;
  --text-primary: #1d1d1f;
  --text-secondary: #86868b;
  --accent-blue: #0071e3;
  --accent-green: #30d158;
  --border-color: rgba(0, 0, 0, 0.1);
  
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.04), 0 1px 4px rgba(0, 0, 0, 0.04);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.08), 0 16px 32px rgba(0, 0, 0, 0.08);
  --shadow-xl: 0 24px 48px rgba(0, 0, 0, 0.12);
}

/* Dark Mode Overrides */
@media (prefers-color-scheme: dark) {
  :root {
    --bg-primary: #000000;
    --bg-secondary: #1d1d1f;
    --text-primary: #f5f5f7;
    --text-secondary: #a1a1a6;
    --accent-blue: #0a84ff;
    --border-color: rgba(255, 255, 255, 0.1);
    
    --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
    --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.4);
    --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.5);
  }
}
```

---

## 4. Typography & Optical Sizing

- **System Font Stack**:
  ```css
  font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Segoe UI', sans-serif;
  ```
- **Fluid Heading Scale**:
  ```css
  h1 {
    font-size: clamp(2.5rem, 5vw + 1rem, 4.5rem);
    font-weight: 700;
    line-height: 1.05;
    letter-spacing: -0.02em; /* Tight negative tracking for large titles */
  }
  h2 {
    font-size: clamp(1.75rem, 3vw + 1rem, 2.75rem);
    font-weight: 600;
    line-height: 1.15;
    letter-spacing: -0.015em;
  }
  body {
    font-size: 1rem;
    line-height: 1.6;
    letter-spacing: 0;
  }
  ```
- **Typography Guidelines**:
  - Tighten letter-spacing (`-0.02em`) on large display titles.
  - Keep body copy letter-spacing neutral (`0`).
  - Tighten line-height on headings (`1.05 - 1.15`), expand on body copy (`1.5 - 1.6`).

---

## 5. Spacing Grid & Radius Tokens

```css
:root {
  /* 8px Spacing Scale */
  --space-1: 0.25rem;  /* 4px */
  --space-2: 0.5rem;   /* 8px */
  --space-3: 0.75rem;  /* 12px */
  --space-4: 1rem;     /* 16px */
  --space-6: 1.5rem;   /* 24px */
  --space-8: 2rem;     /* 32px */
  --space-12: 3rem;    /* 48px */
  --space-16: 4rem;    /* 64px */

  /* Border Radius Scale */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;
  --radius-2xl: 24px;
  --radius-full: 9999px;
}
```

---

## 6. Materials, Depth & Glassmorphism

```css
/* Glassmorphic Card Surface */
.glass-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  padding: var(--space-8);
  box-shadow: var(--shadow-lg), inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

/* Translucent Header Chrome */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid var(--border-color);
}
```

---

## 7. Reusable Component Patterns

### Hero Section
```html
<section class="hero">
  <div class="hero-content">
    <h1 class="hero-title">
      Hi, I'm <span class="gradient-text">Your Name</span>
    </h1>
    <p class="hero-subtitle">Creating fluid, human-centric web experiences.</p>
    <div class="hero-actions">
      <button class="btn btn-primary">Explore Work</button>
      <button class="btn btn-secondary">Get in Touch</button>
    </div>
  </div>
</section>
```

### Interactive Project Card
```css
.project-card {
  border-radius: var(--radius-2xl);
  overflow: hidden;
  background: var(--bg-secondary);
  box-shadow: var(--shadow-md);
  transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1), box-shadow 300ms cubic-bezier(0.16, 1, 0.3, 1);
}

.project-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-xl);
}
```

### Pill Buttons
```css
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 44px; /* Touch target minimum */
  padding: 0.75rem 1.5rem;
  border-radius: var(--radius-full);
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: transform 150ms cubic-bezier(0.16, 1, 0.3, 1), background-color 200ms ease;
}

.btn:active {
  transform: scale(0.96);
}

.btn-primary {
  background: var(--accent-blue);
  color: #ffffff;
}

.btn-secondary {
  background: transparent;
  color: var(--accent-blue);
  border: 1px solid var(--accent-blue);
}
```

---

## 8. Performance & Accessibility Guidelines

### Performance Rules
- ✅ Animate exclusively GPU-accelerated properties (`transform`, `opacity`).
- ❌ Avoid animating layout properties (`width`, `height`, `margin`, `top`, `left`).
- ✅ Use `will-change` only on active/imminent interactions.

### Accessibility Standards
- **Touch Targets**: Minimum **44px x 44px** hit target for all buttons and interactive controls.
- **Contrast**: Maintain WCAG AA minimum 4.5:1 ratio for body text against backgrounds.
- **Reduced Motion**:
  ```css
  @media (prefers-reduced-motion: reduce) {
    *, ::before, ::after {
      animation-duration: 0.01ms !important;
      transition-duration: 0.01ms !important;
      scroll-behavior: auto !important;
    }
    .sheet, .card {
      transform: none !important;
    }
  }
  ```
- **Reduced Transparency**:
  ```css
  @media (prefers-reduced-transparency: reduce) {
    .glass-card, .navbar {
      backdrop-filter: none !important;
      background: var(--bg-primary) !important;
    }
  }
  ```

---

## 9. AI Assistant Checklist for Apple Design

When designing, building, or auditing Apple-inspired UIs:
1. **Semantic HTML**: Use proper tags (`<main>`, `<section>`, `<nav>`, `<article>`, `<button>`).
2. **Fluid Scaling**: Use `clamp()` for headings and responsive grid layouts.
3. **Glass & Depth**: Implement semi-transparent materials with backdrop filters and layered ambient shadows.
4. **Touch & Focus**: Ensure min 44px targets and explicit visible focus indicators (`outline`).
5. **Spring Motion**: Use 300ms-400ms spring/ease curves and eliminate input lag.
6. **Dark Mode & A11y**: Always provide dark mode tokens and support `prefers-reduced-motion`.
