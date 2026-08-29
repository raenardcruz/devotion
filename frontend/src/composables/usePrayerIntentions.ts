import { ref, computed } from 'vue';

export type DevotionKey = 'divine_mercy' | 'rosary' | 'st_michael';

export interface PrayerCategory {
  id: string;
  name: string;
  icon: string;
  description: string;
  color: string;
}

export interface PrayerIntention {
  id: string;
  title: string;
  description?: string;
  category: string;
  icon: string;
  createdAt: number;
  completed: boolean;
}

export const PRAYER_CATEGORIES: PrayerCategory[] = [
  {
    id: 'health',
    name: 'Healing & Health',
    icon: '🩺',
    description: 'For physical, emotional, and mental healing, recovery, and strength for the sick',
    color: '#30D158'
  },
  {
    id: 'family',
    name: 'Family & Loved Ones',
    icon: '👨‍👩‍👧',
    description: 'For unity, protection, blessings upon marriage, children, parents, and friends',
    color: '#0A84FF'
  },
  {
    id: 'peace',
    name: 'Peace & Reconciliation',
    icon: '🕊️',
    description: 'For peace in families, communities, our country, and war-torn regions',
    color: '#5E5CE6'
  },
  {
    id: 'souls',
    name: 'Holy Souls in Purgatory',
    icon: '🕯️',
    description: 'For departed loved ones, forgotten souls, and holy souls awaiting eternal rest',
    color: '#FF9F0A'
  },
  {
    id: 'faith',
    name: 'Faith & Conversion',
    icon: '🙏',
    description: 'For spiritual renewal, return of fallen-away Catholics, and perseverance in faith',
    color: '#BF5AF2'
  },
  {
    id: 'protection',
    name: 'Protection & Deliverance',
    icon: '🛡️',
    description: 'For deliverance from evil, spiritual warfare, safety from harm, and angelic guidance',
    color: '#72383D'
  },
  {
    id: 'thanksgiving',
    name: 'Thanksgiving & Praise',
    icon: '❤️',
    description: 'In gratitude for answered prayers, daily graces, and God\'s infinite mercy',
    color: '#FF375F'
  },
  {
    id: 'work',
    name: 'Work, Studies & Guidance',
    icon: '💼',
    description: 'For career discernment, employment, exams, financial needs, and wise decisions',
    color: '#64D2FF'
  },
  {
    id: 'personal',
    name: 'Personal Petitions',
    icon: '⭐',
    description: 'For special intentions, unspoken desires of the heart, and virtue growth',
    color: '#FFD60A'
  }
];

export const INTENTION_PRESETS = [
  { title: 'For the healing and recovery of a sick loved one', category: 'health', icon: '🩺', description: 'Asking Our Lord for complete bodily and spiritual healing and comfort.' },
  { title: 'For peace, harmony, and protection of my family', category: 'family', icon: '👨‍👩‍👧', description: 'Placing our family under the mantle of Our Lady and the Sacred Heart.' },
  { title: 'For the holy souls suffering in Purgatory', category: 'souls', icon: '🕯️', description: 'Offering this prayer for the relief and eternal rest of all departed souls.' },
  { title: 'For spiritual conversion and renewal of hearts', category: 'faith', icon: '🙏', description: 'Praying for loved ones to experience God\'s loving mercy and return to the sacraments.' },
  { title: 'In thanksgiving for blessings and answered prayers', category: 'thanksgiving', icon: '❤️', description: 'Thanking God for His faithfulness, graces, and providential care.' },
  { title: 'For angelic protection against spiritual attacks and harm', category: 'protection', icon: '🛡️', description: 'Invoking St. Michael and the guardian angels for defense and safe keeping.' }
];

export function usePrayerIntentions(devotionKey: DevotionKey) {
  const storageKey = `devotion_intentions_${devotionKey}`;
  const intentions = ref<PrayerIntention[]>(loadIntentions());

  function loadIntentions(): PrayerIntention[] {
    try {
      const stored = localStorage.getItem(storageKey);
      if (stored) {
        return JSON.parse(stored);
      }
    } catch (e) {
      console.error(`Failed to load intentions for ${devotionKey}:`, e);
    }
    return [];
  }

  function saveIntentions() {
    try {
      localStorage.setItem(storageKey, JSON.stringify(intentions.value));
    } catch (e) {
      console.error(`Failed to save intentions for ${devotionKey}:`, e);
    }
  }

  // Active intentions (can include all non-empty intentions; completed items remain accessible in list)
  const activeIntentions = computed(() => {
    return intentions.value.filter(item => !item.completed && item.title.trim().length > 0);
  });

  const totalCount = computed(() => intentions.value.length);
  const activeCount = computed(() => activeIntentions.value.length);
  const hasIntentions = computed(() => activeIntentions.value.length > 0);

  function addIntention(data: { title: string; description?: string; category?: string; icon?: string }) {
    if (!data.title.trim()) return;

    const matchedCategory = PRAYER_CATEGORIES.find(c => c.id === data.category) || PRAYER_CATEGORIES[0]!;
    
    const newIntention: PrayerIntention = {
      id: 'intent_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
      title: data.title.trim(),
      description: data.description?.trim() || '',
      category: matchedCategory.id,
      icon: data.icon || matchedCategory.icon,
      createdAt: Date.now(),
      completed: false
    };

    intentions.value.unshift(newIntention);
    saveIntentions();
    return newIntention;
  }

  function updateIntention(id: string, updates: Partial<Omit<PrayerIntention, 'id' | 'createdAt'>>) {
    const idx = intentions.value.findIndex(item => item.id === id);
    if (idx !== -1) {
      intentions.value[idx] = {
        ...intentions.value[idx]!,
        ...updates,
        title: updates.title !== undefined ? updates.title.trim() : intentions.value[idx]!.title,
        description: updates.description !== undefined ? updates.description.trim() : intentions.value[idx]!.description
      };
      saveIntentions();
    }
  }

  function deleteIntention(id: string) {
    intentions.value = intentions.value.filter(item => item.id !== id);
    saveIntentions();
  }

  function toggleComplete(id: string) {
    const item = intentions.value.find(i => i.id === id);
    if (item) {
      item.completed = !item.completed;
      saveIntentions();
    }
  }

  function clearAll() {
    intentions.value = [];
    saveIntentions();
  }

  function clearCompleted() {
    intentions.value = intentions.value.filter(item => !item.completed);
    saveIntentions();
  }

  function refresh() {
    intentions.value = loadIntentions();
  }

  return {
    intentions,
    activeIntentions,
    totalCount,
    activeCount,
    hasIntentions,
    addIntention,
    updateIntention,
    deleteIntention,
    toggleComplete,
    clearAll,
    clearCompleted,
    refresh
  };
}
