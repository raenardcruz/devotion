import { ref } from 'vue';
import { useDevotionApi } from './useDevotionApi';
import { useAccessControl } from './useAccessControl';

export type MagazineReadingKey =
  | 'cover'
  | 'first_reading'
  | 'responsorial_psalm'
  | 'second_reading'
  | 'gospel';

export interface MagazineReadingImage {
  date: string;
  readingKey: MagazineReadingKey;
  image: string;
  r2Key?: string;
  cached?: boolean;
}

export interface MagazineReadingsInput {
  first_reading?: { citation: string; text: string; context?: string };
  responsorial_psalm?: { citation: string; text: string; context?: string };
  second_reading?: { citation: string; text: string; context?: string };
  gospel?: { citation: string; text: string; context?: string };
  pope_quote?: string;
}

const STORAGE_PREFIX = 'devotion_magazine_images_';

export function useMagazineImages() {
  const { fetchWithAuth, getApiUrl } = useDevotionApi();
  const { isAdmin, syncAdminState } = useAccessControl();

  const images = ref<Partial<Record<MagazineReadingKey, MagazineReadingImage>>>({});
  const generatingKeys = ref<Set<MagazineReadingKey>>(new Set());
  const isGeneratingAll = ref(false);
  const imageError = ref<string | null>(null);

  const readLocalCache = (date: string): Partial<Record<MagazineReadingKey, MagazineReadingImage>> => {
    if (typeof localStorage === 'undefined') return {};
    try {
      const raw = localStorage.getItem(`${STORAGE_PREFIX}${date}`);
      if (!raw) return {};
      return JSON.parse(raw) as Partial<Record<MagazineReadingKey, MagazineReadingImage>>;
    } catch {
      return {};
    }
  };

  const writeLocalCache = (
    date: string,
    data: Partial<Record<MagazineReadingKey, MagazineReadingImage>>
  ) => {
    if (typeof localStorage === 'undefined') return;
    try {
      localStorage.setItem(`${STORAGE_PREFIX}${date}`, JSON.stringify(data));
    } catch {
      // Ignore quota exceeded errors in localStorage
    }
  };

  const resolveImageSrc = (img?: MagazineReadingImage): string => {
    if (!img) return '';
    if (img.image && (img.image.startsWith('data:') || img.image.startsWith('http'))) {
      return img.image;
    }
    if (img.image && img.image.startsWith('/')) {
      return getApiUrl(img.image);
    }
    if (img.r2Key) {
      return getApiUrl(`/api/manga/asset?key=${encodeURIComponent(img.r2Key)}`);
    }
    return img.image || '';
  };

  const loadImagesForDate = async (date: string) => {
    imageError.value = null;
    const local = readLocalCache(date);
    images.value = { ...local };

    try {
      const res = await fetchWithAuth<{
        date: string;
        images?: Record<string, MagazineReadingImage>;
      }>(`/api/magazine/images?date=${encodeURIComponent(date)}`);

      if (res?.images) {
        const merged: Partial<Record<MagazineReadingKey, MagazineReadingImage>> = {
          ...local,
        };
        for (const [key, val] of Object.entries(res.images)) {
          if (val && val.image) {
            merged[key as MagazineReadingKey] = val;
          }
        }
        images.value = merged;
        writeLocalCache(date, merged);
      }
    } catch {
      // Keep local cache if backend cache lookup fails
    }
  };

  const generateReadingImage = async (
    date: string,
    readingKey: MagazineReadingKey,
    citation: string,
    text?: string,
    context?: string,
    forceRegenerate = false
  ): Promise<MagazineReadingImage | null> => {
    syncAdminState();
    if (!isAdmin.value) {
      imageError.value = 'Magazine image generation is restricted to administrators.';
      return null;
    }

    imageError.value = null;
    const nextSet = new Set(generatingKeys.value);
    nextSet.add(readingKey);
    generatingKeys.value = nextSet;

    try {
      const res = await fetchWithAuth<MagazineReadingImage>('/api/magazine/reading-image', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          date,
          readingKey,
          citation,
          text: text || '',
          context: context || '',
          forceRegenerate,
        }),
      });

      if (res?.image) {
        const updated = {
          ...images.value,
          [readingKey]: res,
        };
        images.value = updated;
        writeLocalCache(date, updated);
        return res;
      }
      return null;
    } catch (err: any) {
      imageError.value = err?.message || `Failed to generate illustration for ${readingKey}.`;
      return null;
    } finally {
      const afterSet = new Set(generatingKeys.value);
      afterSet.delete(readingKey);
      generatingKeys.value = afterSet;
    }
  };

  const generateMissingImages = async (date: string, readings: MagazineReadingsInput) => {
    syncAdminState();
    if (!isAdmin.value) {
      imageError.value = 'Magazine image generation is restricted to administrators.';
      return;
    }

    isGeneratingAll.value = true;
    imageError.value = null;

    try {
      const tasks: Array<{
        key: MagazineReadingKey;
        citation: string;
        text?: string;
        context?: string;
      }> = [];

      if (readings.gospel) {
        tasks.push({
          key: 'cover',
          citation: readings.gospel.citation,
          text: readings.gospel.text,
          context: readings.gospel.context,
        });
      }
      if (readings.first_reading) {
        tasks.push({
          key: 'first_reading',
          citation: readings.first_reading.citation,
          text: readings.first_reading.text,
          context: readings.first_reading.context,
        });
      }
      if (readings.responsorial_psalm) {
        tasks.push({
          key: 'responsorial_psalm',
          citation: readings.responsorial_psalm.citation,
          text: readings.responsorial_psalm.text,
          context: readings.responsorial_psalm.context,
        });
      }
      if (readings.second_reading) {
        tasks.push({
          key: 'second_reading',
          citation: readings.second_reading.citation,
          text: readings.second_reading.text,
          context: readings.second_reading.context,
        });
      }
      if (readings.gospel) {
        tasks.push({
          key: 'gospel',
          citation: readings.gospel.citation,
          text: readings.gospel.text,
          context: readings.gospel.context,
        });
      }

      for (const task of tasks) {
        if (!images.value[task.key]?.image) {
          await generateReadingImage(date, task.key, task.citation, task.text, task.context, false);
        }
      }
    } finally {
      isGeneratingAll.value = false;
    }
  };

  return {
    images,
    generatingKeys,
    isGeneratingAll,
    imageError,
    resolveImageSrc,
    loadImagesForDate,
    generateReadingImage,
    generateMissingImages,
  };
}
