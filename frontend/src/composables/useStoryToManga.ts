import { ref } from 'vue';
import { useDevotionApi } from './useDevotionApi';
import { useAccessControl } from './useAccessControl';

export type ComicStyle = 'manga' | 'comic';

export interface MangaCharacter {
  name: string;
  physicalDescription: string;
  personality: string;
  role: string;
}

export interface MangaSetting {
  timePeriod: string;
  location: string;
  mood: string;
}

export interface MangaStoryAnalysis {
  title: string;
  story?: string;
  moralLesson?: string;
  characters: MangaCharacter[];
  setting: MangaSetting;
}

export interface MangaCharacterReference {
  name: string;
  image: string; // base64 data URL or asset/R2 URL
  r2Key?: string;
  description: string;
}

export interface MangaUploadedReference {
  id?: string;
  name: string;
  image: string;
  fileName?: string;
}

export interface MangaPanelPlan {
  panelNumber: number;
  characters: string[];
  sceneDescription: string;
  dialogue?: string;
  narrative?: string;
  cameraAngle: string;
  visualMood: string;
  layoutRow?: number;
  colSpan?: number;
  aspectRatio?: string;
}

export interface MangaStoryBreakdown {
  panels: MangaPanelPlan[];
}

export interface MangaGeneratedPanel {
  panelNumber: number;
  image: string; // base64 data URL or asset/R2 URL
  r2Key?: string;
}

export interface MangaStoryVersionInfo {
  id: number;
  version: number;
  style: ComicStyle;
  title: string;
  imageModel: string;
  createdAt: string;
  updatedAt: string;
}

export interface MangaStoryRecord {
  id: number;
  date: string;
  style: ComicStyle;
  version?: number;
  title: string;
  moralLesson: string;
  story: string;
  imageModel: string;
  storyAnalysis: MangaStoryAnalysis;
  characterReferences: MangaCharacterReference[];
  storyBreakdown: MangaStoryBreakdown;
  generatedPanels: MangaGeneratedPanel[];
  r2Prefix?: string;
  r2Keys?: string[];
  r2Stored?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface ReadingsPayload {
  first_reading?: { citation: string; text: string; context?: string };
  responsorial_psalm?: { citation: string; text: string; context?: string };
  second_reading?: { citation: string; text: string; context?: string };
  gospel?: { citation: string; text: string; context?: string };
  pope_quote?: string;
}

const DB_NAME = 'MangaGeneratorDB';
const DB_VERSION = 1;
const IMAGE_STORE = 'images';
const IMAGE_MODEL_STORAGE_KEY = 'devotion_manga_image_model';
const STYLE_STORAGE_KEY = 'devotion_manga_style';

export const DEFAULT_IMAGE_MODEL = 'gemini-nano-banana-2.1';

export const AVAILABLE_IMAGE_MODELS = [
  'gemini-nano-banana-2.1',
  'gemini-2.5-flash-image-preview',
  'gemini-2.0-flash-preview-image-generation',
  'imagen-3.0-generate-002',
];

class MangaImageStorage {
  private db: IDBDatabase | null = null;

  async init(): Promise<void> {
    if (typeof window === 'undefined' || typeof indexedDB === 'undefined') {
      return Promise.resolve();
    }
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION);
      request.onerror = () => reject(request.error);
      request.onsuccess = () => {
        this.db = request.result;
        resolve();
      };
      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result;
        if (!db.objectStoreNames.contains(IMAGE_STORE)) {
          db.createObjectStore(IMAGE_STORE, { keyPath: 'id' });
        }
      };
    });
  }

  async storeImage(id: string, imageData: string): Promise<void> {
    try {
      if (!this.db) await this.init();
      if (!this.db) return;
      return new Promise((resolve, reject) => {
        const tx = this.db!.transaction([IMAGE_STORE], 'readwrite');
        const store = tx.objectStore(IMAGE_STORE);
        const req = store.put({ id, imageData, timestamp: Date.now() });
        req.onerror = () => reject(req.error);
        req.onsuccess = () => resolve();
      });
    } catch {
      // Ignore IndexedDB quota/private mode errors
    }
  }

  async getImage(id: string): Promise<string | null> {
    try {
      if (!this.db) await this.init();
      if (!this.db) return null;
      return new Promise((resolve, reject) => {
        const tx = this.db!.transaction([IMAGE_STORE], 'readonly');
        const store = tx.objectStore(IMAGE_STORE);
        const req = store.get(id);
        req.onerror = () => reject(req.error);
        req.onsuccess = () => {
          resolve(req.result ? req.result.imageData : null);
        };
      });
    } catch {
      return null;
    }
  }
}

const imageStorage = new MangaImageStorage();

export function useStoryToManga() {
  const { fetchWithAuth, getApiUrl } = useDevotionApi();
  const { isAdmin, syncAdminState } = useAccessControl();

  const style = ref<ComicStyle>(
    (typeof localStorage !== 'undefined' && (localStorage.getItem(STYLE_STORAGE_KEY) as ComicStyle)) || 'manga'
  );
  const imageModel = ref<string>(
    (typeof localStorage !== 'undefined' && localStorage.getItem(IMAGE_MODEL_STORAGE_KEY)) || DEFAULT_IMAGE_MODEL
  );
  const availableImageModels = ref<string[]>([...AVAILABLE_IMAGE_MODELS]);

  const generatedStory = ref<string>('');
  const moralLesson = ref<string>('');
  const storyAnalysis = ref<MangaStoryAnalysis | null>(null);
  const characterReferences = ref<MangaCharacterReference[]>([]);
  const storyBreakdown = ref<MangaStoryBreakdown | null>(null);
  const generatedPanels = ref<MangaGeneratedPanel[]>([]);
  const r2Stored = ref<boolean>(false);
  const currentStoryId = ref<number | null>(null);
  const currentVersion = ref<number>(1);
  const availableVersions = ref<MangaStoryVersionInfo[]>([]);

  const isGenerating = ref(false);
  const isAnalyzing = ref(false);
  const isGeneratingCharacters = ref(false);
  const isChunking = ref(false);
  const isGeneratingPanels = ref(false);
  const currentStepText = ref('');
  const error = ref<string | null>(null);
  const failedStep = ref<'analysis' | 'characters' | 'layout' | 'panels' | null>(null);

  const regeneratingCharacters = ref<Set<number>>(new Set());
  const regeneratingPanels = ref<Set<number>>(new Set());

  const loadConfig = async () => {
    try {
      const res = await fetchWithAuth<{ imageModel?: string; models?: string[] }>('/api/manga/config');
      if (res?.models && res.models.length > 0) {
        availableImageModels.value = res.models;
      }
      if (res?.imageModel) {
        imageModel.value = res.imageModel;
        if (typeof localStorage !== 'undefined') {
          localStorage.setItem(IMAGE_MODEL_STORAGE_KEY, res.imageModel);
        }
      }
    } catch {
      // Fallback to localStorage / default gemini-nano-banana-2.1
    }
  };

  const updateImageModel = async (newModel: string) => {
    syncAdminState();
    if (!isAdmin.value) {
      error.value = 'Story & manga generation settings are restricted to administrators.';
      return;
    }
    const trimmed = newModel.trim() || DEFAULT_IMAGE_MODEL;
    imageModel.value = trimmed;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(IMAGE_MODEL_STORAGE_KEY, trimmed);
    }
    try {
      await fetchWithAuth('/api/manga/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ imageModel: trimmed }),
      });
    } catch {
      // Keep local setting if server config endpoint fails
    }
  };

  const updateStyle = (newStyle: ComicStyle) => {
    style.value = newStyle;
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(STYLE_STORAGE_KEY, newStyle);
    }
  };

  const checkStoryExists = async (date: string, anyStyle = false): Promise<boolean> => {
    try {
      const styleQuery = anyStyle ? '' : `&style=${encodeURIComponent(style.value)}`;
      const res = await fetchWithAuth<{
        success: boolean;
        found: boolean;
        story?: MangaStoryRecord;
      }>(`/api/manga/story?date=${encodeURIComponent(date)}${styleQuery}`);
      return Boolean(res?.found && res.story);
    } catch {
      return false;
    }
  };

  const resolveAssetImageUrl = (img: string): string => {
    const trimmed = (img || '').trim();
    if (!trimmed) return '';
    if (trimmed.startsWith('data:') || trimmed.startsWith('http://') || trimmed.startsWith('https://')) {
      return trimmed;
    }
    if (trimmed.startsWith('/')) {
      return getApiUrl(trimmed);
    }
    return trimmed;
  };

  const getProxyAssetUrl = (r2Key?: string): string => {
    if (!r2Key) return '';
    return getApiUrl(`/api/manga/asset?key=${encodeURIComponent(r2Key)}`);
  };

  const applyServerStoryRecord = (
    rec: MangaStoryRecord,
    versions?: MangaStoryVersionInfo[]
  ) => {
    currentStoryId.value = rec.id || null;
    currentVersion.value = rec.version && rec.version > 0 ? rec.version : 1;
    if (Array.isArray(versions)) {
      availableVersions.value = versions;
    }
    if (rec.style === 'manga' || rec.style === 'comic') {
      style.value = rec.style;
    }
    generatedStory.value = rec.story || rec.storyAnalysis?.story || '';
    moralLesson.value = rec.moralLesson || rec.storyAnalysis?.moralLesson || '';
    storyAnalysis.value = rec.storyAnalysis || null;
    storyBreakdown.value = rec.storyBreakdown || null;
    r2Stored.value = Boolean(rec.r2Stored);

    if (Array.isArray(rec.characterReferences)) {
      characterReferences.value = rec.characterReferences.map((c) => ({
        ...c,
        image: resolveAssetImageUrl(c.image),
      }));
    }
    if (Array.isArray(rec.generatedPanels)) {
      generatedPanels.value = rec.generatedPanels.map((p) => ({
        ...p,
        image: resolveAssetImageUrl(p.image),
      }));
    }
  };

  const loadStateFromServer = async (date: string, anyStyle = false): Promise<boolean> => {
    try {
      const styleQuery = anyStyle ? '' : `&style=${encodeURIComponent(style.value)}`;
      const res = await fetchWithAuth<{
        success: boolean;
        found: boolean;
        story?: MangaStoryRecord;
        versions?: MangaStoryVersionInfo[];
      }>(`/api/manga/story?date=${encodeURIComponent(date)}${styleQuery}`);

      if (res?.found && res.story) {
        applyServerStoryRecord(res.story, res.versions);
        return (
          Boolean(storyAnalysis.value) &&
          characterReferences.value.length > 0 &&
          generatedPanels.value.length > 0
        );
      }
      availableVersions.value = [];
      return false;
    } catch {
      return false;
    }
  };

  const loadVersionFromServer = async (id: number): Promise<boolean> => {
    try {
      const res = await fetchWithAuth<{
        success: boolean;
        found: boolean;
        story?: MangaStoryRecord;
        versions?: MangaStoryVersionInfo[];
      }>(`/api/manga/story?id=${encodeURIComponent(String(id))}`);

      if (res?.found && res.story) {
        applyServerStoryRecord(res.story, res.versions);
        return true;
      }
      return false;
    } catch {
      return false;
    }
  };

  const saveStateToServer = async (date: string, createNewVersion = false): Promise<boolean> => {
    if (!storyAnalysis.value) return false;
    try {
      const res = await fetchWithAuth<{
        success: boolean;
        story?: MangaStoryRecord;
        versions?: MangaStoryVersionInfo[];
      }>('/api/manga/story', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: createNewVersion ? undefined : currentStoryId.value || undefined,
          date,
          style: style.value,
          createNewVersion,
          imageModel: imageModel.value,
          generatedStory: generatedStory.value,
          moralLesson: moralLesson.value,
          storyAnalysis: storyAnalysis.value,
          characterReferences: characterReferences.value,
          storyBreakdown: storyBreakdown.value || { panels: [] },
          generatedPanels: generatedPanels.value,
        }),
      });
      if (res?.story) {
        applyServerStoryRecord(res.story, res.versions);
        return true;
      }
      return false;
    } catch {
      return false;
    }
  };

  const getCacheKey = (date: string) => `devotion_manga_state_${date}`;

  const saveStateForDate = async (date: string) => {
    if (typeof localStorage === 'undefined' || !storyAnalysis.value) return;
    try {
      const stateMeta = {
        date,
        style: style.value,
        imageModel: imageModel.value,
        generatedStory: generatedStory.value,
        moralLesson: moralLesson.value,
        storyAnalysis: storyAnalysis.value,
        characterReferences: characterReferences.value.map(({ image: _img, ...rest }) => rest),
        storyBreakdown: storyBreakdown.value,
        generatedPanels: generatedPanels.value.map(({ image: _img, ...rest }) => rest),
      };
      localStorage.setItem(getCacheKey(date), JSON.stringify(stateMeta));

      for (let i = 0; i < characterReferences.value.length; i++) {
        const char = characterReferences.value[i];
        if (char && char.image) {
          await imageStorage.storeImage(`${date}_char_${i}_${char.name}`, char.image);
        }
      }
      for (const panel of generatedPanels.value) {
        if (panel.image) {
          await imageStorage.storeImage(`${date}_panel_${panel.panelNumber}`, panel.image);
        }
      }
    } catch {
      // Ignore storage quota errors
    }
  };

  const loadStateForDate = async (date: string): Promise<boolean> => {
    if (typeof localStorage === 'undefined') return false;
    try {
      const raw = localStorage.getItem(getCacheKey(date));
      if (!raw) return false;
      const parsed = JSON.parse(raw);
      if (!parsed || !parsed.storyAnalysis) return false;

      generatedStory.value = parsed.generatedStory || parsed.storyAnalysis.story || '';
      moralLesson.value = parsed.moralLesson || parsed.storyAnalysis.moralLesson || '';
      storyAnalysis.value = parsed.storyAnalysis;
      storyBreakdown.value = parsed.storyBreakdown || null;

      if (Array.isArray(parsed.characterReferences)) {
        const hydratedChars: MangaCharacterReference[] = [];
        for (let i = 0; i < parsed.characterReferences.length; i++) {
          const c = parsed.characterReferences[i];
          const img = (await imageStorage.getImage(`${date}_char_${i}_${c.name}`)) || '';
          hydratedChars.push({ ...c, image: img });
        }
        characterReferences.value = hydratedChars;
      }

      if (Array.isArray(parsed.generatedPanels)) {
        const hydratedPanels: MangaGeneratedPanel[] = [];
        for (const p of parsed.generatedPanels) {
          const img = (await imageStorage.getImage(`${date}_panel_${p.panelNumber}`)) || '';
          if (img) {
            hydratedPanels.push({ panelNumber: p.panelNumber, image: img });
          }
        }
        generatedPanels.value = hydratedPanels;
      }

      return characterReferences.value.length > 0 && generatedPanels.value.length > 0;
    } catch {
      return false;
    }
  };

  const generateSingleCharacter = async (
    character: MangaCharacter,
    setting: MangaSetting,
    uploadedCharacterReferences?: MangaUploadedReference[]
  ): Promise<MangaCharacterReference> => {
    const payload: Record<string, unknown> = {
      characters: [character],
      setting,
      style: style.value,
      imageModel: imageModel.value,
    };
    if (uploadedCharacterReferences && uploadedCharacterReferences.length > 0) {
      payload.uploadedCharacterReferences = uploadedCharacterReferences;
    }

    const res = await fetchWithAuth<{
      success: boolean;
      characterReferences: MangaCharacterReference[];
    }>('/api/manga/generate-character-refs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    if (!res?.characterReferences || res.characterReferences.length === 0 || !res.characterReferences[0]) {
      throw new Error(`Failed to generate character sheet for ${character.name}`);
    }
    return res.characterReferences[0];
  };

  const generateSinglePanel = async (
    panel: MangaPanelPlan,
    chars: MangaCharacterReference[],
    setting: MangaSetting
  ): Promise<MangaGeneratedPanel> => {
    const res = await fetchWithAuth<{
      success: boolean;
      generatedPanel: MangaGeneratedPanel;
    }>('/api/manga/generate-panel', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        panel,
        characterReferences: chars,
        setting,
        style: style.value,
        imageModel: imageModel.value,
      }),
    });

    if (!res?.generatedPanel) {
      throw new Error(`Failed to generate panel ${panel.panelNumber}`);
    }
    return res.generatedPanel;
  };

  // One-click end-to-end generation pipeline (story from readings or user-supplied customStory/customCharacters -> character sheets -> panel layout -> panel images)
  const generateStoryAndManga = async (
    date: string,
    readings: ReadingsPayload,
    forceRegenerate = false,
    customStory?: string,
    customCharacters?: MangaCharacter[],
    uploadedCharacterReferences?: MangaUploadedReference[]
  ) => {
    syncAdminState();
    error.value = null;
    const trimmedCustomStory = (customStory || '').trim();
    const validCustomChars = (customCharacters || []).filter(
      (c) => c.name.trim() !== '' && c.physicalDescription.trim() !== ''
    );

    if (!forceRegenerate && !trimmedCustomStory && validCustomChars.length === 0) {
      const hasServer = await loadStateFromServer(date);
      if (hasServer) {
        return;
      }
      if (isAdmin.value) {
        const hasCached = await loadStateForDate(date);
        if (hasCached) {
          return;
        }
      }
    }

    if (!isAdmin.value) {
      generatedStory.value = '';
      moralLesson.value = '';
      storyAnalysis.value = null;
      characterReferences.value = [];
      storyBreakdown.value = null;
      generatedPanels.value = [];
      r2Stored.value = false;
      error.value = 'Story & manga generation is restricted to administrators. Log in as admin to generate new manga.';
      return;
    }

    isGenerating.value = true;
    error.value = null;
    failedStep.value = null;
    generatedStory.value = '';
    moralLesson.value = '';
    storyAnalysis.value = null;
    characterReferences.value = [];
    storyBreakdown.value = null;
    generatedPanels.value = [];
    r2Stored.value = false;

    try {
      // Step 1: Analyze user-provided story or generate custom story from the day's readings
      isAnalyzing.value = true;
      currentStepText.value = trimmedCustomStory
        ? 'Step 1/4: Analyzing your custom story, characters, and setting...'
        : 'Step 1/4: Crafting custom story from today’s readings & lessons (Problem → Listening to Readings → Applying Lesson → Resolution)...';

      const analyzePayload: Record<string, unknown> = {
        date,
        readings,
        style: style.value,
      };
      if (trimmedCustomStory) {
        analyzePayload.story = trimmedCustomStory;
      } else if (validCustomChars.length > 0) {
        const charSummary = validCustomChars
          .map((c) => `${c.name} (${c.role || 'Character'}: ${c.physicalDescription})`)
          .join('; ');
        analyzePayload.protagonistHint = `Feature these specific characters in the story (Problem -> Listening to Readings -> Applying Lesson -> Resolution): ${charSummary}`;
      }

      const analyzeRes = await fetchWithAuth<{
        success: boolean;
        story: string;
        moralLesson?: string;
        analysis: MangaStoryAnalysis;
      }>('/api/manga/analyze-story', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(analyzePayload),
      });

      if (!analyzeRes?.analysis) {
        throw new Error(
          trimmedCustomStory
            ? 'Failed to analyze your custom story'
            : 'Failed to generate and analyze story from readings'
        );
      }

      if (validCustomChars.length > 0) {
        analyzeRes.analysis.characters = validCustomChars;
      }

      generatedStory.value = analyzeRes.story || analyzeRes.analysis.story || '';
      moralLesson.value = analyzeRes.moralLesson || analyzeRes.analysis.moralLesson || '';
      storyAnalysis.value = analyzeRes.analysis;
      isAnalyzing.value = false;

      // Step 2: Generate Character Reference Sheets progressively
      isGeneratingCharacters.value = true;
      const totalChars = analyzeRes.analysis.characters.length;
      const charRefs: MangaCharacterReference[] = [];

      for (let i = 0; i < totalChars; i++) {
        const charData = analyzeRes.analysis.characters[i];
        if (!charData) continue;
        currentStepText.value = `Step 2/4: Generating character design ${i + 1}/${totalChars} (${charData.name}) with ${imageModel.value}...`;
        const refItem = await generateSingleCharacter(
          charData,
          analyzeRes.analysis.setting,
          uploadedCharacterReferences
        );
        charRefs.push(refItem);
        characterReferences.value = [...charRefs];
      }
      isGeneratingCharacters.value = false;

      // Step 3: Plan Comic Layout / Chunk Story into Sequential Panels + Narratives
      isChunking.value = true;
      currentStepText.value = 'Step 3/4: Planning manga/comic panels, dialogue & panel narratives...';

      const chunkRes = await fetchWithAuth<{
        success: boolean;
        storyBreakdown: MangaStoryBreakdown;
      }>('/api/manga/chunk-story', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          story: generatedStory.value,
          characters: analyzeRes.analysis.characters,
          setting: analyzeRes.analysis.setting,
          style: style.value,
        }),
      });

      if (!chunkRes?.storyBreakdown?.panels?.length) {
        throw new Error('Failed to plan comic panel layout');
      }

      storyBreakdown.value = chunkRes.storyBreakdown;
      isChunking.value = false;

      // Step 4: Generate Individual Manga/Comic Panels progressively
      isGeneratingPanels.value = true;
      const panels = chunkRes.storyBreakdown.panels;
      const renderedPanels: MangaGeneratedPanel[] = [];

      for (let i = 0; i < panels.length; i++) {
        const panelPlan = panels[i];
        if (!panelPlan) continue;
        currentStepText.value = `Step 4/4: Illustrating panel ${i + 1}/${panels.length} with ${imageModel.value}...`;
        const generated = await generateSinglePanel(
          panelPlan,
          characterReferences.value,
          analyzeRes.analysis.setting
        );
        renderedPanels.push(generated);
        generatedPanels.value = [...renderedPanels];
      }
      isGeneratingPanels.value = false;

      currentStepText.value = 'Story & Manga generation complete! Saving to storage...';
      const createNewVersion = forceRegenerate || Boolean(trimmedCustomStory) || validCustomChars.length > 0;
      const savedOnServer = await saveStateToServer(date, createNewVersion);
      if (!savedOnServer) {
        await saveStateForDate(date);
      }
      currentStepText.value = 'Story & Manga generation complete!';
    } catch (err: any) {
      if (isAnalyzing.value) failedStep.value = 'analysis';
      else if (isGeneratingCharacters.value) failedStep.value = 'characters';
      else if (isChunking.value) failedStep.value = 'layout';
      else if (isGeneratingPanels.value) failedStep.value = 'panels';

      isAnalyzing.value = false;
      isGeneratingCharacters.value = false;
      isChunking.value = false;
      isGeneratingPanels.value = false;
      currentStepText.value = '';
      error.value = err?.message || 'An error occurred while generating the story manga.';
    } finally {
      isGenerating.value = false;
    }
  };

  const regenerateCharacter = async (index: number, date?: string) => {
    syncAdminState();
    if (!isAdmin.value) {
      error.value = 'Story & manga generation is restricted to administrators.';
      return;
    }
    if (!storyAnalysis.value) return;
    const charData = storyAnalysis.value.characters[index];
    if (!charData) return;

    const nextSet = new Set(regeneratingCharacters.value);
    nextSet.add(index);
    regeneratingCharacters.value = nextSet;
    error.value = null;

    try {
      const updatedRef = await generateSingleCharacter(charData, storyAnalysis.value.setting);
      const updatedList = [...characterReferences.value];
      updatedList[index] = updatedRef;
      characterReferences.value = updatedList;
      if (date) {
        const savedOnServer = await saveStateToServer(date);
        if (!savedOnServer) {
          await saveStateForDate(date);
        }
      }
    } catch (err: any) {
      error.value = err?.message || `Failed to regenerate character ${charData.name}`;
    } finally {
      const finalSet = new Set(regeneratingCharacters.value);
      finalSet.delete(index);
      regeneratingCharacters.value = finalSet;
    }
  };

  const regeneratePanel = async (panelNumber: number, date?: string) => {
    syncAdminState();
    if (!isAdmin.value) {
      error.value = 'Story & manga generation is restricted to administrators.';
      return;
    }
    if (!storyAnalysis.value || !storyBreakdown.value) return;
    const panelPlan = storyBreakdown.value.panels.find((p) => p.panelNumber === panelNumber);
    if (!panelPlan) return;

    const nextSet = new Set(regeneratingPanels.value);
    nextSet.add(panelNumber);
    regeneratingPanels.value = nextSet;
    error.value = null;

    try {
      const updatedPanel = await generateSinglePanel(
        panelPlan,
        characterReferences.value,
        storyAnalysis.value.setting
      );
      const idx = generatedPanels.value.findIndex((p) => p.panelNumber === panelNumber);
      const updatedList = [...generatedPanels.value];
      if (idx >= 0) {
        updatedList[idx] = updatedPanel;
      } else {
        updatedList.push(updatedPanel);
      }
      generatedPanels.value = updatedList;
      if (date) {
        const savedOnServer = await saveStateToServer(date);
        if (!savedOnServer) {
          await saveStateForDate(date);
        }
      }
    } catch (err: any) {
      error.value = err?.message || `Failed to regenerate panel ${panelNumber}`;
    } finally {
      const finalSet = new Set(regeneratingPanels.value);
      finalSet.delete(panelNumber);
      regeneratingPanels.value = finalSet;
    }
  };

  return {
    style,
    imageModel,
    availableImageModels,
    generatedStory,
    moralLesson,
    storyAnalysis,
    characterReferences,
    storyBreakdown,
    generatedPanels,
    r2Stored,
    currentStoryId,
    currentVersion,
    availableVersions,
    isGenerating,
    isAnalyzing,
    isGeneratingCharacters,
    isChunking,
    isGeneratingPanels,
    currentStepText,
    error,
    failedStep,
    regeneratingCharacters,
    regeneratingPanels,
    loadConfig,
    updateImageModel,
    updateStyle,
    getProxyAssetUrl,
    checkStoryExists,
    loadStateFromServer,
    loadVersionFromServer,
    generateStoryAndManga,
    regenerateCharacter,
    regeneratePanel,
  };
}
