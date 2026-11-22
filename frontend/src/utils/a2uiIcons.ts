export type A2UIIconName =
  | 'scripture'
  | 'cross'
  | 'magisterium'
  | 'catechism'
  | 'scroll'
  | 'dove'
  | 'chalice'
  | 'lamp'
  | 'compass'
  | 'crown'
  | 'bridge'
  | 'quill'
  | 'pray'
  | 'heart'
  | 'scales'
  | 'quote'
  | 'bookmark'
  | 'checklist'
  | 'info'
  | 'globe'
  | 'shield'
  | 'hourglass'
  | 'olive'
  | 'star';

export type A2UIIconTone = 'primary' | 'gold' | 'neutral';

export interface A2UIIconDefinition {
  name: A2UIIconName;
  label: string;
  defaultTone: A2UIIconTone;
  /** Optional subtle translucent fill shapes rendered only in default (non-Kindle) theme */
  fills?: string[];
  /** Crisp 24x24 fine-line vector stroke paths */
  strokes: string[];
}

export const A2UI_ICON_PACK: Record<A2UIIconName, A2UIIconDefinition> = {
  scripture: {
    name: 'scripture',
    label: 'Sacred Scripture',
    defaultTone: 'primary',
    fills: [
      'M3 6.5C3 5.7 3.7 5 4.5 5H9c1.7 0 3 1.1 3 2.5V19c0-1.1-1.1-2-2.5-2H4.5C3.7 17 3 16.3 3 15.5V6.5Z',
      'M21 6.5C21 5.7 20.3 5 19.5 5H15c-1.7 0-3 1.1-3 2.5V19c0-1.1 1.1-2 2.5-2h5c.8 0 1.5-.7 1.5-1.5V6.5Z',
    ],
    strokes: [
      'M12 7.5C12 6.1 10.7 5 9 5H4.5C3.7 5 3 5.7 3 6.5v9C3 16.3 3.7 17 4.5 17H9.5c1.4 0 2.5.9 2.5 2V7.5Z',
      'M12 7.5C12 6.1 13.3 5 15 5h4.5c.8 0 1.5.7 1.5 1.5v9c0 .8-.7 1.5-1.5 1.5H14.5c-1.4 0-2.5.9-2.5 2V7.5Z',
      'M15.5 9.5h2.5',
      'M15.5 12.5h2.5',
      'M6 9.5h2.5',
      'M12 19v2.5',
    ],
  },
  cross: {
    name: 'cross',
    label: 'Sacred Cross',
    defaultTone: 'primary',
    fills: ['M12 4.5a3.5 3.5 0 1 0 0 7a3.5 3.5 0 1 0 0-7Z'],
    strokes: [
      'M12 3v18',
      'M7 8.5h10',
      'M10 3h4',
      'M9.5 21h5',
    ],
  },
  magisterium: {
    name: 'magisterium',
    label: 'Magisterium & Church Teaching',
    defaultTone: 'primary',
    fills: ['M5 10.5C5 7.4 8.1 5 12 5s7 2.4 7 5.5H5Z'],
    strokes: [
      'M12 2.5v2.5',
      'M10.7 3.7h2.6',
      'M5.5 10.5C5.5 7.5 8.4 5 12 5s6.5 2.5 6.5 5.5',
      'M3.5 10.5h17',
      'M6 10.5v8',
      'M10 10.5v8',
      'M14 10.5v8',
      'M18 10.5v8',
      'M3 18.5h18',
      'M4.5 21h15',
    ],
  },
  catechism: {
    name: 'catechism',
    label: 'Catechism',
    defaultTone: 'primary',
    fills: ['M5.5 4.5A2 2 0 0 1 7.5 2.5H19v16H7.5a2 2 0 0 0-2 2v-16Z'],
    strokes: [
      'M5.5 19.5A2 2 0 0 1 7.5 17.5H19',
      'M7.5 2.5H19v19H7.5a2 2 0 0 1-2-2v-15a2 2 0 0 1 2-2Z',
      'M12.5 7v5',
      'M10.5 9h4',
    ],
  },
  scroll: {
    name: 'scroll',
    label: 'Historical Context & Scroll',
    defaultTone: 'gold',
    fills: ['M7 5.5h10v13H7z'],
    strokes: [
      'M8 4h9a2 2 0 0 1 2 2v1.5a1.5 1.5 0 0 1-1.5 1.5H17v8.5a2.5 2.5 0 0 1-2.5 2.5H6.5A2 2 0 0 1 4.5 18v-1a1.5 1.5 0 0 1 1.5-1.5h1V6a2 2 0 0 1 2-2Z',
      'M10 9h4',
      'M10 12.5h4',
      'M10 16h2.5',
    ],
  },
  dove: {
    name: 'dove',
    label: 'Holy Spirit & Grace',
    defaultTone: 'gold',
    fills: ['M4.5 15.5c2.5-1.5 4.5-4.5 5.5-8.5c1.5 2 3.5 3.5 6.5 3.5c1 0 2 .3 2.8.8c-1.3 1.4-2.8 2.2-4.8 2.2c-1 2.8-3.8 4.8-7 4.8c-1.2 0-2.2-.3-3-.8Z'],
    strokes: [
      'M4.5 16c2.6-1.3 4.6-4.2 5.7-8.5c1.3 2.3 3.4 3.5 6.3 3.5c1.2 0 2.2.3 3 1c-1.1 1.3-2.6 2-4.5 2c-1.1 2.9-3.9 5-7.3 5c-1.2 0-2.3-.3-3.2-.8',
      'M13 11c-.5-2.8-2.2-5-4.5-6.2',
      'M16.8 11.4h.01',
    ],
  },
  chalice: {
    name: 'chalice',
    label: 'Eucharist & Sacraments',
    defaultTone: 'gold',
    fills: ['M7.5 9.5h9a4.5 4.5 0 0 1-9 0Z'],
    strokes: [
      'M12 4.5a2.2 2.2 0 1 0 0.01 0Z',
      'M7.2 9.5h9.6c0 3.2-2.1 5.5-4.8 5.5s-4.8-2.3-4.8-5.5Z',
      'M12 15v5',
      'M8.5 20h7',
    ],
  },
  lamp: {
    name: 'lamp',
    label: 'Spiritual Insight & Wisdom',
    defaultTone: 'gold',
    fills: ['M12 5a5 5 0 0 0-3.2 8.8c.7.6 1.2 1.4 1.2 2.2h4c0-.8.5-1.6 1.2-2.2A5 5 0 0 0 12 5Z'],
    strokes: [
      'M9 17.5h6',
      'M10 20.5h4',
      'M12 4.5a5.2 5.2 0 0 0-3.4 9.1c.7.6 1.2 1.4 1.4 2.4h4c.2-1 .7-1.8 1.4-2.4A5.2 5.2 0 0 0 12 4.5Z',
      'M12 2v1',
      'M4.5 9.5H3.5',
      'M20.5 9.5h-1',
    ],
  },
  compass: {
    name: 'compass',
    label: 'Moral Sense & Guidance',
    defaultTone: 'primary',
    fills: ['M12 3.5a8.5 8.5 0 1 0 0 17a8.5 8.5 0 1 0 0-17Z'],
    strokes: [
      'M12 3.5a8.5 8.5 0 1 0 0.01 0Z',
      'M15.4 8.6l-2.1 4.7l-4.7 2.1l2.1-4.7l4.7-2.1Z',
      'M12 12h.01',
    ],
  },
  crown: {
    name: 'crown',
    label: 'Anagogical Sense & Eternal Kingdom',
    defaultTone: 'gold',
    fills: ['M4 17l1.5-8.5L9.5 12L12 7l2.5 5l4-3.5L20 17H4Z'],
    strokes: [
      'M4 17l1.5-8.5L9.5 12L12 6.5L14.5 12l4-3.5L20 17H4Z',
      'M4.5 20h15',
      'M12 4v1',
    ],
  },
  bridge: {
    name: 'bridge',
    label: 'Typology & Covenant Fulfillment',
    defaultTone: 'primary',
    fills: ['M4 15c1.5-3.5 4.4-5.5 8-5.5s6.5 2 8 5.5H4Z'],
    strokes: [
      'M2.5 10.5h19',
      'M3.5 17.5v-7',
      'M20.5 17.5v-7',
      'M6 17.5c1-3.3 3.2-5 6-5s5 1.7 6 5',
      'M8.5 10.5v2.8',
      'M15.5 10.5v2.8',
      'M2.5 17.5h19',
    ],
  },
  quill: {
    name: 'quill',
    label: 'Church Fathers & Patristic Witness',
    defaultTone: 'gold',
    fills: ['M19.5 4.5c-4.5.3-8.7 2.8-11 7.2L7 15.5l3.8-1.5c4.4-2.3 6.9-6.5 8.7-9.5Z'],
    strokes: [
      'M20 4c-5.2.4-9.5 3.2-11.8 7.8L6.5 16l4.2-1.7C15.3 12 18.1 7.8 20 4Z',
      'M4 20l4.5-4.5',
      'M11 10.5h4.5',
      'M13.5 7.5h3.5',
    ],
  },
  pray: {
    name: 'pray',
    label: 'Prayer & Practical Meditation',
    defaultTone: 'primary',
    fills: ['M12 5.5L8.2 14.2L5.5 17l3.2 2.2L12 16l3.3 3.2L18.5 17l-2.7-2.8L12 5.5Z'],
    strokes: [
      'M12 5.2L8.4 13.8L5.5 16.5L8.4 19l3.6-3.2',
      'M12 5.2l3.6 8.6l2.9 2.7L15.6 19L12 15.8',
      'M12 5.2v10.6',
      'M12 2.5v1',
      'M7.5 4.5l.8.8',
      'M16.5 4.5l-.8.8',
    ],
  },
  heart: {
    name: 'heart',
    label: 'Charity, Mercy & Human Dignity',
    defaultTone: 'primary',
    fills: ['M12 20s-7.5-4.6-7.5-10.2A4.3 4.3 0 0 1 12 6.8a4.3 4.3 0 0 1 7.5 3C19.5 15.4 12 20 12 20Z'],
    strokes: [
      'M12 20s-7.5-4.6-7.5-10.2A4.3 4.3 0 0 1 12 6.8a4.3 4.3 0 0 1 7.5 3C19.5 15.4 12 20 12 20Z',
    ],
  },
  scales: {
    name: 'scales',
    label: 'Justice, Truth & Comparison',
    defaultTone: 'gold',
    fills: ['M4.5 14h5a2.5 2.5 0 0 1-5 0Z', 'M14.5 14h5a2.5 2.5 0 0 1-5 0Z'],
    strokes: [
      'M12 4v16',
      'M8.5 20h7',
      'M5 7.5h14',
      'M7 7.5L4.5 14h5L7 7.5Z',
      'M17 7.5L14.5 14h5L17 7.5Z',
      'M4.5 14a2.5 2.5 0 0 0 5 0',
      'M14.5 14a2.5 2.5 0 0 0 5 0',
    ],
  },
  quote: {
    name: 'quote',
    label: 'Magisterial Quotation',
    defaultTone: 'primary',
    fills: [
      'M5 11.5h4.5v5.5H5v-5.5Z',
      'M14.5 11.5H19v5.5h-4.5v-5.5Z',
    ],
    strokes: [
      'M9.5 11.5H5.5v5.5h4.5v-5.5c0-2.5 1-4 3-4.8',
      'M19 11.5h-4v5.5h4.5v-5.5c0-2.5 1-4 2.5-4.8',
      'M5.5 11.5C5.5 8.5 7 6.8 9.5 6',
      'M15 11.5c0-3 1.5-4.7 4-5.5',
    ],
  },
  bookmark: {
    name: 'bookmark',
    label: 'Key Takeaway',
    defaultTone: 'primary',
    fills: ['M6.5 4.5A1.5 1.5 0 0 1 8 3h8a1.5 1.5 0 0 1 1.5 1.5V20l-5.5-3.6L6.5 20V4.5Z'],
    strokes: [
      'M6.5 4.5A1.5 1.5 0 0 1 8 3h8a1.5 1.5 0 0 1 1.5 1.5V20l-5.5-3.6L6.5 20V4.5Z',
    ],
  },
  checklist: {
    name: 'checklist',
    label: 'Structured Points & Works',
    defaultTone: 'primary',
    strokes: [
      'M4.5 7.5l1.5 1.5l3-3',
      'M4.5 13.5l1.5 1.5l3-3',
      'M4.5 19.5l1.5 1.5l3-3',
      'M12 7.5h7.5',
      'M12 13.5h7.5',
      'M12 19.5h7.5',
    ],
  },
  info: {
    name: 'info',
    label: 'Clarification & Note',
    defaultTone: 'neutral',
    fills: ['M12 3.5a8.5 8.5 0 1 0 0 17a8.5 8.5 0 1 0 0-17Z'],
    strokes: [
      'M12 3.5a8.5 8.5 0 1 0 0.01 0Z',
      'M12 11.2v5',
      'M12 8.2h.01',
    ],
  },
  globe: {
    name: 'globe',
    label: 'Geography & Universal Church',
    defaultTone: 'neutral',
    fills: ['M12 3.5a8.5 8.5 0 1 0 0 17a8.5 8.5 0 1 0 0-17Z'],
    strokes: [
      'M12 3.5a8.5 8.5 0 1 0 0.01 0Z',
      'M3.8 12h16.4',
      'M12 3.5c2.4 2.4 3.6 5.2 3.6 8.5s-1.2 6.1-3.6 8.5c-2.4-2.4-3.6-5.2-3.6-8.5s1.2-6.1 3.6-8.5Z',
    ],
  },
  shield: {
    name: 'shield',
    label: 'Faith & Apologetics',
    defaultTone: 'primary',
    fills: ['M12 3.5l7 2.8v5.4c0 4.5-3 7.8-7 9.3c-4-1.5-7-4.8-7-9.3V6.3l7-2.8Z'],
    strokes: [
      'M12 3.5l7 2.8v5.4c0 4.5-3 7.8-7 9.3c-4-1.5-7-4.8-7-9.3V6.3l7-2.8Z',
      'M12 8v7',
      'M9.2 10.8h5.6',
    ],
  },
  hourglass: {
    name: 'hourglass',
    label: 'Chronology & Liturgical Time',
    defaultTone: 'gold',
    fills: ['M7 4h10c0 3.5-2.5 5.8-5 8c2.5 2.2 5 4.5 5 8H7c0-3.5 2.5-5.8 5-8c-2.5-2.2-5-4.5-5-8Z'],
    strokes: [
      'M6 3.5h12',
      'M6 20.5h12',
      'M7.5 3.5v2.2c0 2.7 1.8 4.9 4.5 6.3c2.7-1.4 4.5-3.6 4.5-6.3V3.5',
      'M7.5 20.5v-2.2c0-2.7 1.8-4.9 4.5-6.3c2.7 1.4 4.5 3.6 4.5 6.3v2.2',
    ],
  },
  olive: {
    name: 'olive',
    label: 'Hope & Reconciliation',
    defaultTone: 'gold',
    fills: ['M12 13.5c-2-2.5-1.5-6 2-8c2.5 2 2 5.5-2 8Z'],
    strokes: [
      'M5 19.5C8.5 16 12.5 12 18.5 5.5',
      'M10.5 14c-2.8-.6-4.3-2.5-4.3-5.2c2.7 0 4.5 1.5 4.8 4.4',
      'M13.5 11c.5-2.8 2.2-4.6 5-5c-.2 2.8-1.7 4.6-4.5 5',
      'M12 15.5c2.6.2 4.3 1.4 5.2 3.5c-2.5.4-4.4-.5-5.5-2.7',
    ],
  },
  star: {
    name: 'star',
    label: 'Guiding Truth & Virtue',
    defaultTone: 'gold',
    fills: ['M12 3.5l2.2 5.7l5.8 2.8l-5.8 2.8L12 20.5l-2.2-5.7L4 12l5.8-2.8L12 3.5Z'],
    strokes: [
      'M12 3.5l2.1 5.6L19.8 12l-5.7 2.9L12 20.5l-2.1-5.6L4.2 12l5.7-2.9L12 3.5Z',
    ],
  },
};

const ICON_ALIASES: Record<string, A2UIIconName> = {
  bible: 'scripture',
  'book-open': 'scripture',
  gospel: 'scripture',
  verse: 'scripture',
  lectionary: 'scripture',
  crucifix: 'cross',
  christ: 'cross',
  redemption: 'cross',
  church: 'magisterium',
  vatican: 'magisterium',
  council: 'magisterium',
  doctrine: 'magisterium',
  pillar: 'magisterium',
  papal: 'magisterium',
  ccc: 'catechism',
  keys: 'catechism',
  book: 'catechism',
  teaching: 'catechism',
  history: 'scroll',
  historical: 'scroll',
  context: 'scroll',
  tradition: 'scroll',
  manuscript: 'scroll',
  literal: 'scroll',
  spirit: 'dove',
  grace: 'dove',
  'holy-spirit': 'dove',
  peace: 'dove',
  eucharist: 'chalice',
  sacrament: 'chalice',
  sacraments: 'chalice',
  mass: 'chalice',
  liturgy: 'chalice',
  light: 'lamp',
  wisdom: 'lamp',
  insight: 'lamp',
  sparkle: 'lamp',
  allegory: 'lamp',
  allegorical: 'lamp',
  moral: 'compass',
  ethics: 'compass',
  guidance: 'compass',
  virtue: 'compass',
  path: 'compass',
  kingdom: 'crown',
  heaven: 'crown',
  anagogical: 'crown',
  glory: 'crown',
  eternity: 'crown',
  typology: 'bridge',
  fulfillment: 'bridge',
  covenant: 'bridge',
  link: 'bridge',
  fathers: 'quill',
  patristic: 'quill',
  saint: 'quill',
  saints: 'quill',
  doctor: 'quill',
  author: 'quill',
  feather: 'quill',
  encyclical: 'quill',
  prayer: 'pray',
  meditation: 'pray',
  devotion: 'pray',
  reflection: 'pray',
  contemplation: 'pray',
  mercy: 'heart',
  charity: 'heart',
  love: 'heart',
  dignity: 'heart',
  social: 'heart',
  humanitas: 'heart',
  justice: 'scales',
  compare: 'scales',
  comparison: 'scales',
  balance: 'scales',
  truth: 'scales',
  distinction: 'scales',
  quotation: 'quote',
  excerpt: 'quote',
  citation: 'quote',
  takeaway: 'bookmark',
  summary: 'bookmark',
  highlight: 'bookmark',
  pin: 'bookmark',
  key: 'bookmark',
  list: 'checklist',
  steps: 'checklist',
  works: 'checklist',
  check: 'checklist',
  note: 'info',
  callout: 'info',
  clarification: 'info',
  alert: 'info',
  geography: 'globe',
  culture: 'globe',
  world: 'globe',
  catholic: 'globe',
  defense: 'shield',
  apologetics: 'shield',
  faith: 'shield',
  guard: 'shield',
  timeline: 'hourglass',
  era: 'hourglass',
  time: 'hourglass',
  chronology: 'hourglass',
  hope: 'olive',
  reconciliation: 'olive',
  forgiveness: 'olive',
  branch: 'olive',
  mary: 'star',
  marian: 'star',
  beacon: 'star',
};

export function resolveA2UIIconName(raw?: string | null): A2UIIconName | null {
  if (!raw || typeof raw !== 'string') return null;
  const normalized = raw.trim().toLowerCase().replace(/[_\s]+/g, '-');
  if (!normalized) return null;
  if (normalized in A2UI_ICON_PACK) {
    return normalized as A2UIIconName;
  }
  if (normalized in ICON_ALIASES) {
    return ICON_ALIASES[normalized] ?? null;
  }
  return null;
}

export function getA2UIIconDefinition(raw?: string | null): A2UIIconDefinition | null {
  const name = resolveA2UIIconName(raw);
  return name ? A2UI_ICON_PACK[name] : null;
}
