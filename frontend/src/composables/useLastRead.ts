import { ref } from 'vue';

export interface LastReadPage {
  title: string;
  path: string;
}

const LAST_READ_STORAGE_KEY = 'last_read_page';

function readStoredLastRead(): LastReadPage | null {
  try {
    const raw = localStorage.getItem(LAST_READ_STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as Partial<LastReadPage>;
    if (
      parsed &&
      typeof parsed.title === 'string' &&
      parsed.title.length > 0 &&
      typeof parsed.path === 'string' &&
      parsed.path.startsWith('/')
    ) {
      return { title: parsed.title, path: parsed.path };
    }
  } catch {
    // Ignore malformed storage entries or restricted storage access.
  }
  return null;
}

const storedLastRead = ref<LastReadPage | null>(readStoredLastRead());

/**
 * Composable for persisting and reading the last visited page in localStorage.
 */
export function useLastRead() {
  const saveLastRead = (page: LastReadPage) => {
    storedLastRead.value = page;
    try {
      localStorage.setItem(LAST_READ_STORAGE_KEY, JSON.stringify(page));
    } catch {
      // Ignore storage quota or restricted storage errors.
    }
  };

  return {
    storedLastRead,
    saveLastRead,
  };
}
