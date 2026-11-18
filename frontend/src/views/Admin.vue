<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import AppButton from '../components/common/AppButton.vue';
import { useDevotionApi } from '../composables/useDevotionApi';
import { useAccessControl } from '../composables/useAccessControl';

interface AccessCodeItem {
  id: number;
  code: string;
  label: string;
  is_active: boolean;
  daily_limit: number;
  used_today: number;
  remaining: number;
  created_at: string;
}

interface AccessRequestItem {
  id: number;
  name: string;
  email: string;
  reason: string;
  status: string;
  generated_code: string;
  created_at: string;
  updated_at: string;
}

interface MangaSummaryItem {
  id: number;
  date: string;
  style: string;
  title: string;
  moralLesson: string;
  imageModel: string;
  panelCount: number;
  characterCount: number;
  r2Prefix: string;
  r2Keys: string[];
  r2Stored: boolean;
  previewImage?: string;
  createdAt: string;
  updatedAt: string;
}

interface Settings {
  context_provider: string;
  fact_checker_provider: string;
  magisterium_llm_provider: string;
  gemini_api_key: string;
  magisterium_api_key: string;
  bible_api_key: string;
  gemini_model: string;
  gemini_live_model: string;
  gemini_image_model: string;
  ollama_model: string;
  context_model?: string;
  context_instruction: string;
  enable_fact_checker: boolean;
  r2_account_id: string;
  r2_endpoint: string;
  r2_bucket: string;
  r2_access_key_id: string;
  r2_secret_access_key: string;
  r2_public_url: string;
}

interface AIUsageRecordItem {
  id: number;
  provider: string;
  model: string;
  feature: string;
  operation: string;
  input_tokens: number;
  output_tokens: number;
  thoughts_tokens: number;
  cached_tokens: number;
  total_tokens: number;
  image_count: number;
  estimated_cost_usd: number;
  created_at: string;
}

interface AIUsageReport {
  period: 'day' | 'week' | 'month' | 'year';
  start_date: string;
  end_date: string;
  summary: {
    total_requests: number;
    total_input_tokens: number;
    total_output_tokens: number;
    total_thoughts_tokens: number;
    total_cached_tokens: number;
    total_tokens: number;
    total_images: number;
    total_cost_usd: number;
  };
  by_provider: {
    provider: string;
    model: string;
    requests: number;
    input_tokens: number;
    output_tokens: number;
    thoughts_tokens: number;
    cached_tokens: number;
    total_tokens: number;
    image_count: number;
    estimated_cost_usd: number;
  }[];
  by_feature: {
    feature: string;
    requests: number;
    total_tokens: number;
    image_count: number;
    estimated_cost_usd: number;
  }[];
  timeline: {
    bucket: string;
    requests: number;
    total_tokens: number;
    estimated_cost_usd: number;
  }[];
  recent_logs: AIUsageRecordItem[];
}

interface AIModelPricingItem {
  id?: number;
  provider: string;
  model: string;
  input_per_million: number;
  output_per_million: number;
  cached_per_million: number;
  image_cost_each: number;
}

type AdminTabId = 'ai' | 'integrations' | 'access' | 'usage';
type UsageSubTabId = 'overview' | 'activity' | 'pricing';

const tabs: { id: AdminTabId; label: string }[] = [
  { id: 'ai', label: 'AI Setup' },
  { id: 'integrations', label: 'Integrations' },
  { id: 'access', label: 'Access' },
  { id: 'usage', label: 'Usage & Cost' },
];

const { fetchWithAuth, getApiUrl, isServerUnreachable } = useDevotionApi();
const { syncAdminState } = useAccessControl();

// Navigation & UI State
const activeTab = ref<AdminTabId>('ai');
const showPromptEditor = ref(false);
const showR2Config = ref(false);
const usageSubTab = ref<UsageSubTabId>('overview');
const showLoginPassword = ref(false);
const showKeys = ref({
  gemini: false,
  magisterium: false,
  bible: false,
  r2Secret: false,
});
const promptTextareaRef = ref<HTMLTextAreaElement | null>(null);

// Auth State
const isAuthenticated = ref(false);
const loginUsername = ref('admin');
const loginPassword = ref('');
const authError = ref<string | null>(null);
const authLoading = ref(false);

// Settings State
const defaultSettings = (): Settings => ({
  context_provider: 'ollama',
  fact_checker_provider: 'gemini',
  magisterium_llm_provider: 'ollama',
  gemini_api_key: '',
  magisterium_api_key: '',
  bible_api_key: '',
  gemini_model: 'gemini-4-argon',
  gemini_live_model: 'gemini-3.8-live',
  gemini_image_model: 'gemini-nano-banana-2.1',
  ollama_model: 'gemma4:cloud',
  context_instruction: '',
  enable_fact_checker: true,
  r2_account_id: '',
  r2_endpoint: '',
  r2_bucket: '',
  r2_access_key_id: '',
  r2_secret_access_key: '',
  r2_public_url: '',
});

const settings = ref<Settings>(defaultSettings());
const savedSnapshot = ref<string>('');

const geminiModelsList = ref<string[]>([
  'gemini-4-argon',
  'gemini-3.8-flash',
  'gemini-3.5-flash-lite',
  'gemini-3.1-flash-lite',
]);
const geminiLiveModelsList = ref<string[]>([
  'gemini-3.8-live',
  'gemini-3.8-live-extended-thinking',
  'gemini-2.5-flash-native-audio',
  'gemini-2.0-flash-exp',
  'gemini-2.0-flash',
]);
const geminiImageModelsList = ref<string[]>([
  'gemini-nano-banana-2.1',
  'gemini-2.5-flash-image-preview',
  'gemini-2.0-flash-preview-image-generation',
  'imagen-3.0-generate-002',
]);
const ollamaModelsList = ref<string[]>([]);
const loadingOllamaModels = ref(false);
const loadingSettings = ref(false);
const saveLoading = ref(false);
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null);

// Access Codes & Guest Requests State
const accessCodes = ref<AccessCodeItem[]>([]);
const accessRequests = ref<AccessRequestItem[]>([]);
const loadingAccess = ref(false);
const newCodeLabel = ref('');
const newCustomCode = ref('');
const newDailyLimit = ref(10);
const creatingCode = ref(false);
const copiedKey = ref<string | null>(null);
const accessMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null);

// R2 Storage & Generated Manga Management State
const testingR2 = ref(false);
const r2TestMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null);
const generatedMangaList = ref<MangaSummaryItem[]>([]);
const loadingMangaList = ref(false);
const r2BackendConfigured = ref(false);
const regeneratingMangaIds = ref<Set<number>>(new Set());
const deletingMangaIds = ref<Set<number>>(new Set());
const confirmDeleteMangaId = ref<number | null>(null);
const mangaActionMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null);

const isR2FormConfigured = computed(() => {
  const s = settings.value;
  return Boolean(
    (s.r2_account_id?.trim() || s.r2_endpoint?.trim()) &&
      s.r2_bucket?.trim() &&
      s.r2_access_key_id?.trim() &&
      s.r2_secret_access_key?.trim()
  );
});

const resolveMangaPreviewUrl = (url?: string): string => {
  if (!url) return '';
  if (url.startsWith('data:') || url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  if (url.startsWith('/')) {
    return getApiUrl(url);
  }
  return url;
};

const hasUnsavedChanges = computed(() => {
  if (!savedSnapshot.value) return false;
  return JSON.stringify(settings.value) !== savedSnapshot.value;
});

const privacyModeLabel = computed(() => {
  const s = settings.value;
  const usesCloud =
    s.context_provider === 'gemini' ||
    s.magisterium_llm_provider === 'gemini' ||
    (s.enable_fact_checker && s.fact_checker_provider !== 'ollama');
  return usesCloud ? 'Hybrid Cloud + Local' : 'Local Privacy Mode';
});

const isLocalPrivacyOnly = computed(() => privacyModeLabel.value === 'Local Privacy Mode');

const configuredKeysCount = computed(() => {
  let count = 0;
  if (settings.value.gemini_api_key?.trim()) count++;
  if (settings.value.magisterium_api_key?.trim()) count++;
  if (settings.value.bible_api_key?.trim()) count++;
  return count;
});

const promptStats = computed(() => {
  const text = settings.value.context_instruction || '';
  const chars = text.length;
  const lines = text ? text.split('\n').length : 0;
  return { chars, lines };
});

const pendingRequestsCount = computed(() =>
  accessRequests.value.filter((r) => r.status === 'pending').length
);

const isSettingsTab = computed(() => activeTab.value === 'ai' || activeTab.value === 'integrations');

onMounted(() => {
  const token = localStorage.getItem('adminToken');
  if (token) {
    isAuthenticated.value = true;
    syncAdminState();
    loadSettings();
    loadAccessData();
    loadMangaList();
    loadAIUsage();
    loadAIPricing();
  }
});

// AI Usage & Cost State
const usagePeriod = ref<'day' | 'week' | 'month' | 'year'>('month');
const usageAnchorDate = ref<string>(new Date().toISOString().slice(0, 10));
const usageReport = ref<AIUsageReport | null>(null);
const loadingUsage = ref(false);
const usageError = ref<string | null>(null);
const cleanupBeforeDate = ref<string>(
  new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().slice(0, 10)
);
const cleaningUp = ref(false);
const cleanupMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null);

const loadAIUsage = async () => {
  loadingUsage.value = true;
  usageError.value = null;
  try {
    const params = new URLSearchParams({ period: usagePeriod.value });
    if (usageAnchorDate.value) {
      params.set('date', usageAnchorDate.value);
    }
    const res = await fetchWithAuth<AIUsageReport>(`/api/admin/ai-usage?${params.toString()}`);
    usageReport.value = res;
  } catch (err: any) {
    usageError.value = err.message || 'Failed to load AI usage report.';
  } finally {
    loadingUsage.value = false;
  }
};

const setUsagePeriod = (period: 'day' | 'week' | 'month' | 'year') => {
  usagePeriod.value = period;
  loadAIUsage();
};

const purgeAIUsage = async () => {
  if (!cleanupBeforeDate.value) {
    cleanupMessage.value = { type: 'error', text: 'Please select a cutoff date.' };
    return;
  }
  if (!window.confirm(`Purge all AI usage records older than ${cleanupBeforeDate.value}? This cannot be undone.`)) {
    return;
  }
  cleaningUp.value = true;
  cleanupMessage.value = null;
  try {
    const res = await fetchWithAuth<{ success: boolean; deleted_count: number; before_date: string }>(
      '/api/admin/ai-usage/cleanup',
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ before_date: cleanupBeforeDate.value }),
      }
    );
    cleanupMessage.value = {
      type: 'success',
      text: `Purged ${res.deleted_count} record${res.deleted_count === 1 ? '' : 's'} older than ${res.before_date}.`,
    };
    await loadAIUsage();
  } catch (err: any) {
    cleanupMessage.value = { type: 'error', text: err.message || 'Failed to purge old records.' };
  } finally {
    cleaningUp.value = false;
  }
};

const formatUSD = (val: number) => {
  if (!val) return '$0.0000';
  if (val < 0.01) return `$${val.toFixed(4)}`;
  return `$${val.toFixed(2)}`;
};

const formatNumber = (val: number) => {
  return (val || 0).toLocaleString();
};

const maxTimelineCost = computed(() => {
  const pts = usageReport.value?.timeline || [];
  if (!pts.length) return 0;
  return Math.max(...pts.map((p) => p.estimated_cost_usd), 0);
});

// Model Pricing Settings ($ / 1M tokens) & Frontend Cost Computation
const modelPricingList = ref<AIModelPricingItem[]>([]);
const loadingPricing = ref(false);
const savingPricing = ref(false);
const pricingMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null);
const newPricingRule = ref<AIModelPricingItem>({
  provider: 'gemini',
  model: '',
  input_per_million: 0.15,
  output_per_million: 0.60,
  cached_per_million: 0.0375,
  image_cost_each: 0,
});

const loadAIPricing = async () => {
  loadingPricing.value = true;
  try {
    const res = await fetchWithAuth<{ pricing: AIModelPricingItem[] }>('/api/admin/ai-pricing');
    modelPricingList.value = res?.pricing || [];
  } catch (err: any) {
    console.error('Failed to load AI model pricing:', err);
  } finally {
    loadingPricing.value = false;
  }
};

const saveAIPricing = async () => {
  savingPricing.value = true;
  pricingMessage.value = null;
  try {
    const res = await fetchWithAuth<{ success: boolean; pricing: AIModelPricingItem[] }>('/api/admin/ai-pricing', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ pricing: modelPricingList.value }),
    });
    if (res?.pricing) {
      modelPricingList.value = res.pricing;
    }
    pricingMessage.value = { type: 'success', text: 'Model pricing saved. Frontend costs updated.' };
  } catch (err: any) {
    pricingMessage.value = { type: 'error', text: err.message || 'Failed to save model pricing.' };
  } finally {
    savingPricing.value = false;
  }
};

const addPricingRule = async () => {
  const modelName = newPricingRule.value.model.trim();
  if (!modelName) {
    pricingMessage.value = { type: 'error', text: 'Please enter a model name.' };
    return;
  }
  savingPricing.value = true;
  pricingMessage.value = null;
  try {
    const res = await fetchWithAuth<{ success: boolean; pricing: AIModelPricingItem[] }>('/api/admin/ai-pricing', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ rule: { ...newPricingRule.value, model: modelName } }),
    });
    if (res?.pricing) {
      modelPricingList.value = res.pricing;
    }
    newPricingRule.value.model = '';
    pricingMessage.value = { type: 'success', text: `Added pricing rule for ${modelName}.` };
  } catch (err: any) {
    pricingMessage.value = { type: 'error', text: err.message || 'Failed to add pricing rule.' };
  } finally {
    savingPricing.value = false;
  }
};

const deletePricingRule = async (item: AIModelPricingItem) => {
  try {
    const params = new URLSearchParams({ provider: item.provider, model: item.model });
    const res = await fetchWithAuth<{ success: boolean; pricing: AIModelPricingItem[] }>(
      `/api/admin/ai-pricing?${params.toString()}`,
      { method: 'DELETE' }
    );
    if (res?.pricing) {
      modelPricingList.value = res.pricing;
    }
  } catch (err: any) {
    pricingMessage.value = { type: 'error', text: err.message || 'Failed to delete pricing rule.' };
  }
};

const findPricingRule = (provider: string, model: string): AIModelPricingItem | undefined => {
  const p = (provider || '').trim().toLowerCase();
  const m = (model || '').trim().toLowerCase();
  return modelPricingList.value.find(
    (r) => r.provider.trim().toLowerCase() === p && r.model.trim().toLowerCase() === m
  );
};

const computeCostFromRates = (
  provider: string,
  model: string,
  inputTokens: number,
  outputTokens: number,
  thoughtsTokens: number,
  cachedTokens: number,
  imageCount: number,
  fallbackCost: number
): number => {
  const rule = findPricingRule(provider, model);
  if (!rule) {
    return fallbackCost || 0;
  }
  const uncachedInput = Math.max(0, (inputTokens || 0) - (cachedTokens || 0));
  const cachedRate = rule.cached_per_million > 0 ? rule.cached_per_million : rule.input_per_million * 0.25;
  const inputCost = (uncachedInput / 1_000_000) * (rule.input_per_million || 0);
  const cachedCost = ((cachedTokens || 0) / 1_000_000) * cachedRate;
  const outputCost = (((outputTokens || 0) + (thoughtsTokens || 0)) / 1_000_000) * (rule.output_per_million || 0);
  const imgCost = (imageCount || 0) * (rule.image_cost_each || 0);
  return inputCost + cachedCost + outputCost + imgCost;
};

const computedByProvider = computed(() => {
  const rows = usageReport.value?.by_provider || [];
  return rows.map((row) => ({
    ...row,
    computed_cost_usd: computeCostFromRates(
      row.provider,
      row.model,
      row.input_tokens,
      row.output_tokens,
      row.thoughts_tokens || 0,
      row.cached_tokens || 0,
      row.image_count,
      row.estimated_cost_usd
    ),
  }));
});

const computedTotalCostUSD = computed(() => {
  if (!usageReport.value) return 0;
  return computedByProvider.value.reduce((sum, row) => sum + row.computed_cost_usd, 0);
});

const handleLogin = async () => {
  authError.value = null;
  authLoading.value = true;
  try {
    const res = await fetchWithAuth<{ success: boolean; token?: string; message?: string }>('/api/admin/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: loginUsername.value,
        password: loginPassword.value,
      }),
    });

    if (res.success && res.token) {
      localStorage.setItem('adminToken', res.token);
      syncAdminState();
      isAuthenticated.value = true;
      loadSettings();
      loadAccessData();
      loadMangaList();
      loadAIUsage();
      loadAIPricing();
    } else {
      authError.value = res.message || 'Invalid admin credentials.';
    }
  } catch (err: any) {
    authError.value = err.message || 'Login failed. Please check credentials.';
  } finally {
    authLoading.value = false;
  }
};

const handleLogout = () => {
  localStorage.removeItem('adminToken');
  syncAdminState();
  isAuthenticated.value = false;
  loginPassword.value = '';
  message.value = null;
  savedSnapshot.value = '';
};

const loadAccessData = async () => {
  loadingAccess.value = true;
  try {
    const [codesRes, reqsRes] = await Promise.all([
      fetchWithAuth<{ codes: AccessCodeItem[] }>('/api/admin/access-codes'),
      fetchWithAuth<{ requests: AccessRequestItem[] }>('/api/admin/access-requests'),
    ]);
    accessCodes.value = codesRes?.codes || [];
    accessRequests.value = reqsRes?.requests || [];
  } catch (err: any) {
    console.error('Failed to load access management data:', err);
  } finally {
    loadingAccess.value = false;
  }
};

const generateCode = async () => {
  creatingCode.value = true;
  accessMessage.value = null;
  try {
    const res = await fetchWithAuth<{ success: boolean; code: AccessCodeItem }>('/api/admin/access-codes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        label: newCodeLabel.value.trim(),
        code: newCustomCode.value.trim(),
        daily_limit: Number(newDailyLimit.value) || 10,
      }),
    });
    if (res.code) {
      accessCodes.value.unshift(res.code);
      newCodeLabel.value = '';
      newCustomCode.value = '';
      newDailyLimit.value = 10;
      accessMessage.value = { type: 'success', text: `Generated access code: ${res.code.code}` };
    }
  } catch (err: any) {
    accessMessage.value = { type: 'error', text: err.message || 'Failed to generate access code.' };
  } finally {
    creatingCode.value = false;
  }
};

const toggleCodeActive = async (item: AccessCodeItem) => {
  try {
    await fetchWithAuth('/api/admin/access-codes', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: item.id, is_active: !item.is_active }),
    });
    item.is_active = !item.is_active;
  } catch (err: any) {
    accessMessage.value = { type: 'error', text: err.message || 'Failed to update access code.' };
  }
};

const deleteCode = async (id: number) => {
  try {
    await fetchWithAuth(`/api/admin/access-codes?id=${id}`, { method: 'DELETE' });
    accessCodes.value = accessCodes.value.filter(c => c.id !== id);
  } catch (err: any) {
    accessMessage.value = { type: 'error', text: err.message || 'Failed to delete access code.' };
  }
};

const approveRequest = async (req: AccessRequestItem) => {
  accessMessage.value = null;
  try {
    const res = await fetchWithAuth<{ success: boolean; request: AccessRequestItem; code: AccessCodeItem }>(
      '/api/admin/access-requests/approve',
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: req.id, daily_limit: 10 }),
      }
    );
    if (res.request) {
      const idx = accessRequests.value.findIndex(r => r.id === req.id);
      if (idx !== -1) accessRequests.value[idx] = res.request;
    }
    if (res.code) {
      accessCodes.value.unshift(res.code);
      accessMessage.value = {
        type: 'success',
        text: `Approved ${req.name} — code ${res.code.code} created.`,
      };
    }
  } catch (err: any) {
    accessMessage.value = { type: 'error', text: err.message || 'Failed to approve request.' };
  }
};

const rejectRequest = async (req: AccessRequestItem) => {
  try {
    await fetchWithAuth('/api/admin/access-requests/reject', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: req.id }),
    });
    req.status = 'rejected';
  } catch (err: any) {
    accessMessage.value = { type: 'error', text: err.message || 'Failed to reject request.' };
  }
};

const deleteRequest = async (id: number) => {
  try {
    await fetchWithAuth(`/api/admin/access-requests?id=${id}`, { method: 'DELETE' });
    accessRequests.value = accessRequests.value.filter(r => r.id !== id);
  } catch (err: any) {
    accessMessage.value = { type: 'error', text: err.message || 'Failed to delete request.' };
  }
};

const getInviteLink = (code: string) => {
  if (typeof window === 'undefined') return '';
  return `${window.location.origin}${window.location.pathname}#/?access_code=${encodeURIComponent(code)}`;
};

const copyText = async (text: string, key: string) => {
  try {
    await navigator.clipboard.writeText(text);
    copiedKey.value = key;
    setTimeout(() => {
      if (copiedKey.value === key) copiedKey.value = null;
    }, 2000);
  } catch {
    // ignore clipboard errors
  }
};

const loadSettings = async () => {
  loadingSettings.value = true;
  try {
    const data = await fetchWithAuth<Settings>('/api/admin/settings');
    if (data) {
      settings.value = { ...defaultSettings(), ...data };
      if (!settings.value.gemini_image_model) {
        settings.value.gemini_image_model = 'gemini-nano-banana-2.1';
      }
      savedSnapshot.value = JSON.stringify(settings.value);
      showR2Config.value = !isR2FormConfigured.value;
    }
    fetchGeminiModels();
    fetchGeminiLiveModels();
    fetchGeminiImageModels();
    fetchOllamaModels();
  } catch (err: any) {
    console.error('Failed to load settings:', err);
  } finally {
    loadingSettings.value = false;
  }
};

const discardChanges = () => {
  if (!savedSnapshot.value) return;
  settings.value = JSON.parse(savedSnapshot.value);
  message.value = null;
};

const fetchGeminiModels = async () => {
  try {
    const res = await fetchWithAuth<{ models: string[] }>('/api/admin/gemini-models');
    if (res && res.models && res.models.length > 0) {
      geminiModelsList.value = res.models;
    }
  } catch (err) {
    console.error('Failed to fetch gemini models:', err);
  }
};

const fetchGeminiLiveModels = async () => {
  try {
    const res = await fetchWithAuth<{ models: string[] }>('/api/admin/gemini-live-models');
    if (res && res.models && res.models.length > 0) {
      geminiLiveModelsList.value = res.models;
    }
  } catch (err) {
    console.error('Failed to fetch gemini live models:', err);
  }
};

const fetchGeminiImageModels = async () => {
  try {
    const res = await fetchWithAuth<{ models: string[] }>('/api/admin/gemini-image-models');
    if (res && res.models && res.models.length > 0) {
      geminiImageModelsList.value = res.models;
    }
  } catch (err) {
    console.error('Failed to fetch gemini image models:', err);
  }
};

const fetchOllamaModels = async () => {
  loadingOllamaModels.value = true;
  try {
    const res = await fetchWithAuth<{ models: string[] }>('/api/admin/ollama-models');
    if (res && res.models && res.models.length > 0) {
      ollamaModelsList.value = res.models;
    } else {
      ollamaModelsList.value = ['gemma4:cloud', 'qwen2.5:4b', 'llama3.2:3b'];
    }
  } catch (err) {
    console.error('Failed to fetch ollama models:', err);
    ollamaModelsList.value = ['gemma4:cloud', 'qwen2.5:4b', 'llama3.2:3b'];
  } finally {
    loadingOllamaModels.value = false;
  }
};

const saveSettings = async () => {
  saveLoading.value = true;
  message.value = null;
  try {
    const res = await fetchWithAuth<{ success?: boolean; settings?: Settings }>('/api/admin/settings', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(settings.value),
    });
    if (res?.settings) {
      settings.value = { ...defaultSettings(), ...res.settings };
    }
    savedSnapshot.value = JSON.stringify(settings.value);
    r2BackendConfigured.value = isR2FormConfigured.value;
    message.value = { type: 'success', text: 'Settings successfully saved to PostgreSQL database.' };
  } catch (err: any) {
    message.value = { type: 'error', text: err.message || 'Failed to save settings to backend.' };
  } finally {
    saveLoading.value = false;
  }
};

const loadMangaList = async () => {
  loadingMangaList.value = true;
  try {
    const res = await fetchWithAuth<{
      success: boolean;
      r2Configured: boolean;
      stories: MangaSummaryItem[];
    }>('/api/admin/manga');
    generatedMangaList.value = res?.stories || [];
    r2BackendConfigured.value = Boolean(res?.r2Configured);
  } catch (err: any) {
    console.error('Failed to load generated manga list:', err);
  } finally {
    loadingMangaList.value = false;
  }
};

const testR2Connection = async () => {
  testingR2.value = true;
  r2TestMessage.value = null;
  try {
    const res = await fetchWithAuth<{
      success: boolean;
      message?: string;
    }>('/api/admin/r2/test', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        r2_account_id: settings.value.r2_account_id,
        r2_endpoint: settings.value.r2_endpoint,
        r2_bucket: settings.value.r2_bucket,
        r2_access_key_id: settings.value.r2_access_key_id,
        r2_secret_access_key: settings.value.r2_secret_access_key,
        r2_public_url: settings.value.r2_public_url,
      }),
    });
    r2TestMessage.value = {
      type: 'success',
      text: res?.message || 'R2 connection verified successfully.',
    };
  } catch (err: any) {
    r2TestMessage.value = {
      type: 'error',
      text: err?.message || 'Failed to connect to R2 bucket.',
    };
  } finally {
    testingR2.value = false;
  }
};

const regenerateMangaItem = async (item: MangaSummaryItem) => {
  const next = new Set(regeneratingMangaIds.value);
  next.add(item.id);
  regeneratingMangaIds.value = next;
  mangaActionMessage.value = null;
  confirmDeleteMangaId.value = null;

  try {
    const res = await fetchWithAuth<{
      success: boolean;
      summary?: MangaSummaryItem;
    }>('/api/admin/manga/regenerate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        id: item.id,
        date: item.date,
        style: item.style,
        imageModel: item.imageModel,
      }),
    });
    if (res?.summary) {
      const idx = generatedMangaList.value.findIndex((s) => s.id === item.id);
      if (idx !== -1) {
        generatedMangaList.value[idx] = res.summary;
      } else {
        await loadMangaList();
      }
    } else {
      await loadMangaList();
    }
    mangaActionMessage.value = {
      type: 'success',
      text: `Regenerated "${item.title || item.date}" and updated database & R2.`,
    };
  } catch (err: any) {
    mangaActionMessage.value = {
      type: 'error',
      text: err?.message || `Failed to regenerate manga for ${item.date}.`,
    };
  } finally {
    const done = new Set(regeneratingMangaIds.value);
    done.delete(item.id);
    regeneratingMangaIds.value = done;
  }
};

const deleteMangaItem = async (item: MangaSummaryItem) => {
  const next = new Set(deletingMangaIds.value);
  next.add(item.id);
  deletingMangaIds.value = next;
  mangaActionMessage.value = null;
  confirmDeleteMangaId.value = null;

  try {
    await fetchWithAuth(`/api/admin/manga?id=${item.id}`, {
      method: 'DELETE',
    });
    generatedMangaList.value = generatedMangaList.value.filter((s) => s.id !== item.id);
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem(`devotion_manga_state_${item.date}`);
    }
    mangaActionMessage.value = {
      type: 'success',
      text: `Deleted "${item.title || item.date}" from database and R2.`,
    };
  } catch (err: any) {
    mangaActionMessage.value = {
      type: 'error',
      text: err?.message || `Failed to delete manga for ${item.date}.`,
    };
  } finally {
    const done = new Set(deletingMangaIds.value);
    done.delete(item.id);
    deletingMangaIds.value = done;
  }
};

const insertPlaceholder = (placeholder: string) => {
  const el = promptTextareaRef.value;
  if (!el) {
    settings.value.context_instruction += placeholder;
    return;
  }
  const start = el.selectionStart ?? settings.value.context_instruction.length;
  const end = el.selectionEnd ?? start;
  const current = settings.value.context_instruction;
  settings.value.context_instruction = current.slice(0, start) + placeholder + current.slice(end);
  requestAnimationFrame(() => {
    el.focus();
    const pos = start + placeholder.length;
    el.setSelectionRange(pos, pos);
  });
};

const formatProviderLabel = (provider: string) => {
  switch (provider) {
    case 'ollama':
      return 'Ollama (Local)';
    case 'gemini':
      return 'Google Gemini';
    case 'magisterium':
      return 'Magisterium AI';
    default:
      return provider;
  }
};
</script>

<template>
  <div class="min-h-screen bg-[#EFE9E1] text-[#322D29] pb-32 sm:pb-24">
    <TopNav />

    <main class="max-w-5xl mx-auto px-4 sm:px-6 py-6 sm:py-10">
      <!-- Server Unreachable Banner -->
      <div
        v-if="isServerUnreachable"
        class="mb-6 p-4 bg-red-50/90 border border-red-300/80 rounded-2xl text-red-900 shadow-sm flex items-start gap-3"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <h3 class="font-semibold text-sm">Server Cannot Be Reached</h3>
          <p class="text-xs text-red-700 mt-0.5">Unable to connect to the backend service. Please verify the API server is running and accessible.</p>
        </div>
      </div>

      <!-- Login View -->
      <div v-if="!isAuthenticated" class="max-w-md mx-auto my-8 sm:my-14">
        <ParchmentCard class="p-7 sm:p-9">
          <div class="text-center mb-7">
            <div class="w-12 h-12 mx-auto mb-4 rounded-2xl bg-[#72383D]/10 border border-[#72383D]/20 flex items-center justify-center text-[#72383D]">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <h1 class="font-serif text-2xl sm:text-3xl font-bold text-[#322D29] tracking-tight">Admin Portal</h1>
            <p class="text-sm text-[#322D29]/65 mt-1.5">Sign in to manage AI pipelines, models, and credentials</p>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-4">
            <div>
              <label for="admin-username" class="block text-[11px] font-bold uppercase tracking-[0.12em] text-[#72383D] mb-1.5">
                Username
              </label>
              <input
                id="admin-username"
                v-model="loginUsername"
                type="text"
                required
                autocomplete="username"
                class="w-full px-3.5 py-2.5 border border-[#D1C7BD] rounded-xl bg-white/85 focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-sm transition-colors"
                placeholder="admin"
              />
            </div>

            <div>
              <label for="admin-password" class="block text-[11px] font-bold uppercase tracking-[0.12em] text-[#72383D] mb-1.5">
                Password
              </label>
              <div class="relative">
                <input
                  id="admin-password"
                  v-model="loginPassword"
                  :type="showLoginPassword ? 'text' : 'password'"
                  required
                  autocomplete="current-password"
                  class="w-full px-3.5 py-2.5 pr-11 border border-[#D1C7BD] rounded-xl bg-white/85 focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-sm transition-colors"
                  placeholder="••••••••"
                />
                <button
                  type="button"
                  @click="showLoginPassword = !showLoginPassword"
                  class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-[#322D29]/50 hover:text-[#72383D] rounded-lg transition-colors cursor-pointer"
                  :aria-label="showLoginPassword ? 'Hide password' : 'Show password'"
                >
                  <svg v-if="!showLoginPassword" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                </button>
              </div>
            </div>

            <div v-if="authError" class="p-3 bg-red-50 border border-red-200 text-red-700 text-xs sm:text-sm rounded-xl">
              {{ authError }}
            </div>

            <div class="pt-2">
              <AppButton type="submit" variant="primary" class="w-full justify-center min-h-[44px]" :disabled="authLoading">
                {{ authLoading ? 'Authenticating...' : 'Sign In' }}
              </AppButton>
            </div>
          </form>
        </ParchmentCard>
      </div>

      <!-- Authenticated Settings Dashboard -->
      <div v-else class="space-y-5">
        <!-- Dashboard Header -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-4 border-b border-[#D1C7BD]/70">
          <div>
            <div class="flex flex-wrap items-center gap-2.5">
              <h1 class="font-serif text-2xl sm:text-3xl font-bold text-[#322D29] tracking-tight">
                Admin Settings
              </h1>
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-semibold border"
                :class="
                  isLocalPrivacyOnly
                    ? 'bg-emerald-50/90 text-emerald-800 border-emerald-200'
                    : 'bg-white/75 text-[#72383D] border-[#D1C7BD]'
                "
              >
                <span
                  class="w-1.5 h-1.5 rounded-full"
                  :class="isLocalPrivacyOnly ? 'bg-emerald-600' : 'bg-[#72383D]'"
                ></span>
                {{ privacyModeLabel }}
              </span>
            </div>
            <p class="text-xs sm:text-sm text-[#322D29]/65 mt-1">
              Manage AI routing, models, integrations, guest access, and usage.
            </p>
          </div>

          <div class="flex items-center gap-2.5 self-start sm:self-auto">
            <AppButton variant="outlined" customClass="!px-4 !py-2 text-xs min-h-[40px]" @click="handleLogout">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Logout
            </AppButton>
          </div>
        </div>

        <!-- Clean 4-Tab Apple-Style Segmented Control (No Horizontal Scroll) -->
        <div
          role="tablist"
          aria-label="Admin Settings Sections"
          class="grid grid-cols-2 sm:grid-cols-4 gap-1.5 p-1.5 rounded-2xl bg-white/65 backdrop-blur-md border border-[#D1C7BD]/70 shadow-2xs"
        >
          <button
            v-for="tab in tabs"
            :key="tab.id"
            type="button"
            role="tab"
            :aria-selected="activeTab === tab.id"
            :aria-controls="`panel-${tab.id}`"
            @click="activeTab = tab.id"
            class="admin-seg-btn flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-xs sm:text-sm font-semibold transition-all duration-200 min-h-[44px] cursor-pointer border"
            :class="
              activeTab === tab.id
                ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white border-transparent hover:border-transparent shadow-sm shadow-[#72383D]/20'
                : 'bg-transparent text-[#322D29]/75 border-transparent hover:border-transparent hover:bg-[#EFE9E1]/80 hover:text-[#72383D]'
            "
          >
            <!-- Tab Icons -->
            <span class="flex-shrink-0">
              <svg v-if="tab.id === 'ai'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
              <svg v-else-if="tab.id === 'integrations'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
              <svg v-else-if="tab.id === 'access'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </span>

            <span class="truncate">{{ tab.label }}</span>

            <!-- Contextual Badges -->
            <span
              v-if="tab.id === 'integrations'"
              class="text-[10px] px-1.5 py-0.5 rounded-full font-bold leading-none"
              :class="activeTab === tab.id ? 'bg-white/20 text-white' : 'bg-[#EFE9E1] text-[#72383D]'"
            >
              {{ configuredKeysCount }}/3
            </span>
            <span
              v-else-if="tab.id === 'access' && pendingRequestsCount > 0"
              class="text-[10px] px-1.5 py-0.5 rounded-full font-bold leading-none"
              :class="activeTab === tab.id ? 'bg-amber-300 text-[#322D29]' : 'bg-amber-100 text-amber-900'"
            >
              {{ pendingRequestsCount }}
            </span>
          </button>
        </div>

        <!-- Global Save Feedback Banner -->
        <Transition name="tab-fade">
          <div
            v-if="message"
            class="p-3.5 px-4 rounded-2xl border text-xs sm:text-sm font-medium flex items-center justify-between gap-3 shadow-xs"
            :class="
              message.type === 'success'
                ? 'bg-emerald-50/90 border-emerald-300/80 text-emerald-900'
                : 'bg-red-50/90 border-red-300/80 text-red-900'
            "
          >
            <div class="flex items-center gap-2.5">
              <svg v-if="message.type === 'success'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-emerald-600 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-red-600 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{{ message.text }}</span>
            </div>
            <button
              type="button"
              @click="message = null"
              class="admin-inline-btn bg-transparent border-transparent text-current opacity-60 hover:opacity-100 hover:bg-black/5 p-1 rounded-lg cursor-pointer"
              aria-label="Dismiss notification"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </Transition>

        <!-- Main Tab Content Surface -->
        <div class="bg-gradient-to-br from-white/85 via-[#EFE9E1]/75 to-[#D9D9D9]/30 border border-[#D1C7BD]/60 backdrop-blur-md rounded-3xl p-5 sm:p-7 shadow-sm">
          <Transition name="tab-fade" mode="out-in">
            <!-- TAB 1: AI SETUP (Pipeline + Models + Collapsible Prompt) -->
            <section
              v-if="activeTab === 'ai'"
              key="ai"
              id="panel-ai"
              role="tabpanel"
              class="space-y-6"
            >
              <!-- 1. Pipeline Routing -->
              <div class="space-y-3.5">
                <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-2">
                  <div>
                    <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                      AI Pipeline
                    </h2>
                    <p class="text-xs text-[#322D29]/65 mt-0.5">
                      Choose providers for scripture generation, verification, and Magisterium summaries.
                    </p>
                  </div>

                  <!-- Compact Active Flow Indicator -->
                  <div class="inline-flex flex-wrap items-center gap-1.5 px-3 py-1.5 rounded-full bg-[#FAF7EE] border border-[#D1C7BD]/70 text-[11px] self-start md:self-auto">
                    <span class="font-semibold text-[#322D29]">{{ formatProviderLabel(settings.context_provider) }}</span>
                    <span class="text-[#AC9C8D]">→</span>
                    <span
                      class="font-semibold"
                      :class="settings.enable_fact_checker ? 'text-[#322D29]' : 'text-[#322D29]/40 line-through'"
                    >
                      {{ formatProviderLabel(settings.fact_checker_provider) }}
                    </span>
                    <span class="text-[#AC9C8D]">→</span>
                    <span class="font-semibold text-[#322D29]">{{ formatProviderLabel(settings.magisterium_llm_provider) }}</span>
                  </div>
                </div>

                <!-- 3-Stage Provider Grid (Fact Checker Toggle Integrated in Stage 2) -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                  <!-- Stage 1: Context Provider -->
                  <div class="p-4 rounded-2xl border border-[#D1C7BD]/80 bg-white/70 flex flex-col justify-between gap-2.5">
                    <div>
                      <span class="text-[10px] font-bold uppercase tracking-[0.12em] text-[#72383D]">Stage 1</span>
                      <label for="context-provider" class="block text-sm font-bold text-[#322D29] mt-0.5">
                        Context Provider
                      </label>
                      <p class="text-[11px] text-[#322D29]/60 mt-0.5">
                        Initial scripture context &amp; reflections.
                      </p>
                    </div>
                    <select
                      id="context-provider"
                      v-model="settings.context_provider"
                      class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-xs font-semibold text-[#322D29] min-h-[42px]"
                    >
                      <option value="ollama">Ollama (Local)</option>
                      <option value="gemini">Google Gemini (Cloud)</option>
                    </select>
                  </div>

                  <!-- Stage 2: Fact Checker Provider + Inline Enable Toggle -->
                  <div
                    class="p-4 rounded-2xl border border-[#D1C7BD]/80 bg-white/70 flex flex-col justify-between gap-2.5 transition-opacity"
                    :class="{ 'opacity-65': !settings.enable_fact_checker }"
                  >
                    <div>
                      <div class="flex items-center justify-between gap-2">
                        <span class="text-[10px] font-bold uppercase tracking-[0.12em] text-[#72383D]">Stage 2</span>
                        <label class="inline-flex items-center gap-1.5 cursor-pointer select-none">
                          <span class="text-[10px] font-semibold text-[#322D29]/65">
                            {{ settings.enable_fact_checker ? 'On' : 'Off' }}
                          </span>
                          <span class="relative inline-flex items-center">
                            <input
                              id="fact-checker-toggle"
                              type="checkbox"
                              v-model="settings.enable_fact_checker"
                              class="sr-only peer"
                            />
                            <span class="w-9 h-5 bg-[#D1C7BD] peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-[#72383D]/30 rounded-full peer peer-checked:after:translate-x-4 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#72383D] block"></span>
                          </span>
                        </label>
                      </div>
                      <label for="fact-checker-provider" class="block text-sm font-bold text-[#322D29] mt-0.5">
                        Fact Checker
                      </label>
                      <p class="text-[11px] text-[#322D29]/60 mt-0.5">
                        Verifies citations &amp; doctrinal fidelity.
                      </p>
                    </div>
                    <select
                      id="fact-checker-provider"
                      v-model="settings.fact_checker_provider"
                      :disabled="!settings.enable_fact_checker"
                      class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-xs font-semibold text-[#322D29] min-h-[42px] disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <option value="gemini">Google Gemini (Cloud)</option>
                      <option value="ollama">Ollama (Local)</option>
                      <option value="magisterium">Magisterium AI (Doctrine)</option>
                    </select>
                  </div>

                  <!-- Stage 3: Magisterium Summary LLM -->
                  <div class="p-4 rounded-2xl border border-[#D1C7BD]/80 bg-white/70 flex flex-col justify-between gap-2.5">
                    <div>
                      <span class="text-[10px] font-bold uppercase tracking-[0.12em] text-[#72383D]">Stage 3</span>
                      <label for="magisterium-llm-provider" class="block text-sm font-bold text-[#322D29] mt-0.5">
                        Magisterium Summary
                      </label>
                      <p class="text-[11px] text-[#322D29]/60 mt-0.5">
                        Summarizes Magisterium citations.
                      </p>
                    </div>
                    <select
                      id="magisterium-llm-provider"
                      v-model="settings.magisterium_llm_provider"
                      class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-xs font-semibold text-[#322D29] min-h-[42px]"
                    >
                      <option value="ollama">Ollama (Local)</option>
                      <option value="gemini">Google Gemini (Cloud)</option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- 2. Models (Compact 2x2 Grid) -->
              <div class="space-y-3 pt-5 border-t border-[#D1C7BD]/60">
                <div class="flex items-center justify-between gap-2">
                  <div>
                    <h3 class="font-serif text-base sm:text-lg font-bold text-[#322D29]">Models</h3>
                    <p class="text-xs text-[#322D29]/60">Cloud Gemini and local Ollama models across text, voice, and manga.</p>
                  </div>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <!-- Gemini Text -->
                  <div class="p-3.5 rounded-2xl border border-[#D1C7BD]/75 bg-white/70 space-y-2">
                    <div class="flex items-center justify-between gap-2">
                      <label for="gemini-model-select" class="text-xs font-bold text-[#322D29]">
                        Gemini Text &amp; Context
                      </label>
                      <span
                        v-if="settings.context_provider === 'gemini' || settings.magisterium_llm_provider === 'gemini'"
                        class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-[#72383D]/10 text-[#72383D]"
                      >
                        Active
                      </span>
                    </div>
                    <select
                      id="gemini-model-select"
                      v-model="settings.gemini_model"
                      class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-xs sm:text-sm font-medium min-h-[42px]"
                    >
                      <option v-for="m in geminiModelsList" :key="m" :value="m">{{ m }}</option>
                    </select>
                  </div>

                  <!-- Gemini Live -->
                  <div class="p-3.5 rounded-2xl border border-[#D1C7BD]/75 bg-white/70 space-y-2">
                    <div class="flex items-center justify-between gap-2">
                      <label for="gemini-live-model-select" class="text-xs font-bold text-[#322D29]">
                        Gemini Live (Voice Sanctuary)
                      </label>
                      <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-[#EFE9E1] text-[#322D29]/70">
                        Voice
                      </span>
                    </div>
                    <select
                      id="gemini-live-model-select"
                      v-model="settings.gemini_live_model"
                      class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-xs sm:text-sm font-medium min-h-[42px]"
                    >
                      <option v-for="m in geminiLiveModelsList" :key="m" :value="m">{{ m }}</option>
                    </select>
                  </div>

                  <!-- Gemini Image -->
                  <div class="p-3.5 rounded-2xl border border-[#D1C7BD]/75 bg-white/70 space-y-2">
                    <div class="flex items-center justify-between gap-2">
                      <label for="gemini-image-model-select" class="text-xs font-bold text-[#322D29]">
                        Gemini Image (Story-to-Manga)
                      </label>
                      <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-[#72383D]/10 text-[#72383D]">
                        Image
                      </span>
                    </div>
                    <select
                      id="gemini-image-model-select"
                      v-model="settings.gemini_image_model"
                      class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-xs sm:text-sm font-medium min-h-[42px]"
                    >
                      <option v-for="m in geminiImageModelsList" :key="m" :value="m">{{ m }}</option>
                    </select>
                  </div>

                  <!-- Ollama Local -->
                  <div class="p-3.5 rounded-2xl border border-[#D1C7BD]/75 bg-white/70 space-y-2">
                    <div class="flex items-center justify-between gap-2">
                      <label for="ollama-model-select" class="text-xs font-bold text-[#322D29]">
                        Ollama Local Model
                      </label>
                      <div class="flex items-center gap-1.5">
                        <span
                          v-if="settings.context_provider === 'ollama' || settings.magisterium_llm_provider === 'ollama'"
                          class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800"
                        >
                          Active
                        </span>
                        <button
                          type="button"
                          @click="fetchOllamaModels"
                          :disabled="loadingOllamaModels"
                          class="admin-inline-btn inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold border border-[#D1C7BD] bg-white hover:border-[#72383D] hover:text-[#72383D] transition-colors cursor-pointer disabled:opacity-50"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="w-3 h-3"
                            :class="{ 'animate-spin': loadingOllamaModels }"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                          </svg>
                          {{ loadingOllamaModels ? 'Syncing' : 'Refresh' }}
                        </button>
                      </div>
                    </div>
                    <select
                      id="ollama-model-select"
                      v-model="settings.ollama_model"
                      :disabled="loadingOllamaModels"
                      class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-xs sm:text-sm font-medium min-h-[42px] disabled:opacity-60"
                    >
                      <option v-for="m in ollamaModelsList" :key="m" :value="m">{{ m }}</option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- 3. Collapsible Bible Context Instruction Prompt -->
              <div class="pt-5 border-t border-[#D1C7BD]/60">
                <div class="rounded-2xl border border-[#D1C7BD]/80 bg-white/65 overflow-hidden">
                  <button
                    type="button"
                    @click="showPromptEditor = !showPromptEditor"
                    :aria-expanded="showPromptEditor"
                    class="admin-inline-btn w-full rounded-none border-0 bg-transparent flex items-center justify-between gap-3 px-4 py-3.5 text-left hover:bg-[#FAF7EE]/80 transition-colors cursor-pointer"
                  >
                    <div class="min-w-0">
                      <div class="flex flex-wrap items-center gap-2">
                        <span class="font-serif text-sm sm:text-base font-bold text-[#322D29]">
                          Bible Context Prompt
                        </span>
                        <span class="text-[11px] font-mono px-2 py-0.5 rounded-full bg-[#EFE9E1] text-[#322D29]/70">
                          {{ promptStats.lines }} lines · {{ promptStats.chars }} chars
                        </span>
                      </div>
                      <p class="text-xs text-[#322D29]/60 mt-0.5">
                        Custom system instructions for scripture context generation.
                      </p>
                    </div>
                    <span class="flex items-center gap-1.5 text-xs font-semibold text-[#72383D] shrink-0">
                      {{ showPromptEditor ? 'Hide' : 'Edit' }}
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="w-4 h-4 transition-transform duration-200"
                        :class="{ 'rotate-180': showPromptEditor }"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                      </svg>
                    </span>
                  </button>

                  <div v-if="showPromptEditor" class="px-4 pb-4 pt-2 border-t border-[#D1C7BD]/50 space-y-3 bg-[#FAF7EE]/40">
                    <div class="flex flex-wrap items-center gap-2 text-xs">
                      <span class="text-[#322D29]/65 font-medium">Insert placeholder:</span>
                      <button
                        type="button"
                        @click="insertPlaceholder('{{citation}}')"
                        class="admin-inline-btn px-2.5 py-1 rounded-full bg-white border border-[#D1C7BD] hover:border-[#72383D] text-[#72383D] font-mono text-xs transition-colors cursor-pointer"
                      >
                        <span v-pre>+ {{citation}}</span>
                      </button>
                      <button
                        type="button"
                        @click="insertPlaceholder('{{passage_text}}')"
                        class="admin-inline-btn px-2.5 py-1 rounded-full bg-white border border-[#D1C7BD] hover:border-[#72383D] text-[#72383D] font-mono text-xs transition-colors cursor-pointer"
                      >
                        <span v-pre>+ {{passage_text}}</span>
                      </button>
                    </div>

                    <textarea
                      ref="promptTextareaRef"
                      v-model="settings.context_instruction"
                      rows="10"
                      placeholder="Enter custom instructions for Bible context generation..."
                      class="w-full px-3.5 py-3 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm leading-relaxed scrollbar-thin"
                    ></textarea>
                  </div>
                </div>
              </div>

              <!-- Compact Ethical AI Caption -->
              <div class="flex items-start gap-2.5 px-3.5 py-2.5 rounded-xl bg-[#FAF7EE]/75 border border-[#D1C7BD]/60 text-[11px] text-[#322D29]/75">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-[#72383D] flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
                <span>
                  <strong class="text-[#72383D]">Magnifica Humanitas:</strong> Local Ollama keeps reflections on-device; AI-generated content is transparently labeled in the UI.
                </span>
              </div>
            </section>

            <!-- TAB 2: INTEGRATIONS (API Keys + R2 Storage + Manga Library) -->
            <section
              v-else-if="activeTab === 'integrations'"
              key="integrations"
              id="panel-integrations"
              role="tabpanel"
              class="space-y-6"
            >
              <!-- 1. API Keys (Apple-style grouped inset list) -->
              <div class="space-y-3">
                <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1">
                  <div>
                    <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                      API Keys
                    </h2>
                    <p class="text-xs text-[#322D29]/65 mt-0.5">
                      Stored in PostgreSQL. Masked dots preserve existing keys unless replaced.
                    </p>
                  </div>
                  <span class="text-xs font-semibold text-[#72383D] bg-[#FAF7EE] border border-[#D1C7BD]/80 px-2.5 py-1 rounded-full self-start sm:self-auto">
                    {{ configuredKeysCount }} / 3 configured
                  </span>
                </div>

                <div class="rounded-2xl border border-[#D1C7BD]/80 bg-white/70 divide-y divide-[#D1C7BD]/55 overflow-hidden">
                  <!-- Gemini API Key -->
                  <div class="p-4 flex flex-col sm:flex-row sm:items-center gap-3">
                    <div class="sm:w-56 shrink-0">
                      <div class="flex items-center gap-2">
                        <label for="gemini-api-key" class="text-sm font-bold text-[#322D29]">
                          Gemini API Key
                        </label>
                        <span
                          class="text-[10px] font-bold px-2 py-0.5 rounded-full"
                          :class="settings.gemini_api_key?.trim() ? 'bg-emerald-100 text-emerald-800' : 'bg-[#EFE9E1] text-[#322D29]/60'"
                        >
                          {{ settings.gemini_api_key?.trim() ? 'Set' : 'Missing' }}
                        </span>
                      </div>
                      <p class="text-[11px] text-[#322D29]/55 mt-0.5">Agent, Live Voice &amp; Context</p>
                    </div>
                    <div class="relative flex-1 flex items-center">
                      <input
                        id="gemini-api-key"
                        v-model="settings.gemini_api_key"
                        :type="showKeys.gemini ? 'text' : 'password'"
                        class="w-full px-3.5 py-2 pr-16 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                        placeholder="AIzaSy..."
                      />
                      <button
                        type="button"
                        @click="showKeys.gemini = !showKeys.gemini"
                        class="admin-inline-btn absolute right-1.5 px-2.5 py-1 text-xs font-semibold text-[#72383D] bg-transparent border-transparent hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer"
                      >
                        {{ showKeys.gemini ? 'Hide' : 'Show' }}
                      </button>
                    </div>
                  </div>

                  <!-- Magisterium AI API Key -->
                  <div class="p-4 flex flex-col sm:flex-row sm:items-center gap-3">
                    <div class="sm:w-56 shrink-0">
                      <div class="flex items-center gap-2">
                        <label for="magisterium-api-key" class="text-sm font-bold text-[#322D29]">
                          Magisterium AI
                        </label>
                        <span
                          class="text-[10px] font-bold px-2 py-0.5 rounded-full"
                          :class="settings.magisterium_api_key?.trim() ? 'bg-emerald-100 text-emerald-800' : 'bg-[#EFE9E1] text-[#322D29]/60'"
                        >
                          {{ settings.magisterium_api_key?.trim() ? 'Set' : 'Missing' }}
                        </span>
                      </div>
                      <p class="text-[11px] text-[#322D29]/55 mt-0.5">Doctrine search &amp; citations</p>
                    </div>
                    <div class="relative flex-1 flex items-center">
                      <input
                        id="magisterium-api-key"
                        v-model="settings.magisterium_api_key"
                        :type="showKeys.magisterium ? 'text' : 'password'"
                        class="w-full px-3.5 py-2 pr-16 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                        placeholder="magisterium-key-..."
                      />
                      <button
                        type="button"
                        @click="showKeys.magisterium = !showKeys.magisterium"
                        class="admin-inline-btn absolute right-1.5 px-2.5 py-1 text-xs font-semibold text-[#72383D] bg-transparent border-transparent hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer"
                      >
                        {{ showKeys.magisterium ? 'Hide' : 'Show' }}
                      </button>
                    </div>
                  </div>

                  <!-- API.Bible Key -->
                  <div class="p-4 flex flex-col sm:flex-row sm:items-center gap-3">
                    <div class="sm:w-56 shrink-0">
                      <div class="flex items-center gap-2">
                        <label for="bible-api-key" class="text-sm font-bold text-[#322D29]">
                          API.Bible Key
                        </label>
                        <span
                          class="text-[10px] font-bold px-2 py-0.5 rounded-full"
                          :class="settings.bible_api_key?.trim() ? 'bg-emerald-100 text-emerald-800' : 'bg-[#EFE9E1] text-[#322D29]/60'"
                        >
                          {{ settings.bible_api_key?.trim() ? 'Set' : 'Missing' }}
                        </span>
                      </div>
                      <p class="text-[11px] text-[#322D29]/55 mt-0.5">Scripture passage fetching</p>
                    </div>
                    <div class="relative flex-1 flex items-center">
                      <input
                        id="bible-api-key"
                        v-model="settings.bible_api_key"
                        :type="showKeys.bible ? 'text' : 'password'"
                        class="w-full px-3.5 py-2 pr-16 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                        placeholder="api-key-here..."
                      />
                      <button
                        type="button"
                        @click="showKeys.bible = !showKeys.bible"
                        class="admin-inline-btn absolute right-1.5 px-2.5 py-1 text-xs font-semibold text-[#72383D] bg-transparent border-transparent hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer"
                      >
                        {{ showKeys.bible ? 'Hide' : 'Show' }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 2. Cloudflare R2 Storage (Collapsible when already configured) -->
              <div class="pt-5 border-t border-[#D1C7BD]/60">
                <div class="rounded-2xl border border-[#D1C7BD]/80 bg-white/70 overflow-hidden">
                  <div class="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                    <div class="min-w-0">
                      <div class="flex flex-wrap items-center gap-2">
                        <h3 class="font-serif text-base font-bold text-[#322D29]">
                          Cloudflare R2 Storage
                        </h3>
                        <span
                          class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full"
                          :class="isR2FormConfigured ? 'bg-emerald-100 text-emerald-800' : 'bg-[#EFE9E1] text-[#322D29]/65 border border-[#D1C7BD]'"
                        >
                          {{ isR2FormConfigured ? 'Configured' : 'Not Set' }}
                        </span>
                      </div>
                      <p class="text-xs text-[#322D29]/60 mt-0.5">
                        Stores generated manga sheets, panels, and manifests.
                      </p>
                    </div>

                    <div class="flex items-center gap-2 shrink-0">
                      <button
                        type="button"
                        @click="testR2Connection"
                        :disabled="testingR2"
                        class="admin-inline-btn px-3 py-1.5 rounded-xl bg-[#FAF7EE] hover:bg-[#EFE9E1] border border-[#D1C7BD] text-xs font-semibold text-[#72383D] transition-colors disabled:opacity-50 cursor-pointer"
                      >
                        {{ testingR2 ? 'Testing...' : 'Test Connection' }}
                      </button>
                      <button
                        type="button"
                        @click="showR2Config = !showR2Config"
                        :aria-expanded="showR2Config"
                        class="admin-inline-btn inline-flex items-center gap-1 px-3 py-1.5 rounded-xl bg-white hover:bg-[#FAF7EE] border border-[#D1C7BD] text-xs font-semibold text-[#322D29] transition-colors cursor-pointer"
                      >
                        {{ showR2Config ? 'Hide Config' : 'Configure' }}
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="w-3.5 h-3.5 transition-transform duration-200"
                          :class="{ 'rotate-180': showR2Config }"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                    </div>
                  </div>

                  <div
                    v-if="r2TestMessage"
                    class="mx-4 mb-3 p-3 rounded-xl text-xs font-medium border"
                    :class="r2TestMessage.type === 'success' ? 'bg-emerald-50 border-emerald-200 text-emerald-800' : 'bg-red-50 border-red-200 text-red-800'"
                  >
                    {{ r2TestMessage.text }}
                  </div>

                  <div
                    v-if="showR2Config"
                    class="p-4 pt-3 border-t border-[#D1C7BD]/55 bg-[#FAF7EE]/45 grid grid-cols-1 sm:grid-cols-2 gap-3.5"
                  >
                    <div>
                      <label for="r2-account-id" class="block text-[11px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                        Cloudflare Account ID
                      </label>
                      <input
                        id="r2-account-id"
                        v-model="settings.r2_account_id"
                        type="text"
                        placeholder="e.g. 9f1c2d3e4a5b6c7d8e9f0a1b2c3d4e5f"
                        class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                      />
                    </div>

                    <div>
                      <label for="r2-bucket" class="block text-[11px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                        R2 Bucket Name
                      </label>
                      <input
                        id="r2-bucket"
                        v-model="settings.r2_bucket"
                        type="text"
                        placeholder="e.g. devotion-manga"
                        class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                      />
                    </div>

                    <div>
                      <label for="r2-access-key-id" class="block text-[11px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                        R2 Access Key ID
                      </label>
                      <input
                        id="r2-access-key-id"
                        v-model="settings.r2_access_key_id"
                        type="text"
                        placeholder="R2 S3 Access Key ID"
                        class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                      />
                    </div>

                    <div>
                      <label for="r2-secret-access-key" class="block text-[11px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                        R2 Secret Access Key
                      </label>
                      <div class="relative flex items-center">
                        <input
                          id="r2-secret-access-key"
                          v-model="settings.r2_secret_access_key"
                          :type="showKeys.r2Secret ? 'text' : 'password'"
                          placeholder="R2 S3 Secret Access Key"
                          class="w-full px-3 py-2 pr-16 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                        />
                        <button
                          type="button"
                          @click="showKeys.r2Secret = !showKeys.r2Secret"
                          class="admin-inline-btn absolute right-1.5 px-2 py-1 text-xs font-semibold text-[#72383D] bg-transparent border-transparent hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer"
                        >
                          {{ showKeys.r2Secret ? 'Hide' : 'Show' }}
                        </button>
                      </div>
                    </div>

                    <div>
                      <label for="r2-endpoint" class="block text-[11px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                        Custom S3 Endpoint <span class="font-normal text-[#322D29]/50 normal-case">(optional)</span>
                      </label>
                      <input
                        id="r2-endpoint"
                        v-model="settings.r2_endpoint"
                        type="text"
                        placeholder="https://<account_id>.r2.cloudflarestorage.com"
                        class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                      />
                    </div>

                    <div>
                      <label for="r2-public-url" class="block text-[11px] font-bold uppercase tracking-wider text-[#72383D] mb-1">
                        Public Base URL <span class="font-normal text-[#322D29]/50 normal-case">(optional)</span>
                      </label>
                      <input
                        id="r2-public-url"
                        v-model="settings.r2_public_url"
                        type="text"
                        placeholder="https://pub-xxx.r2.dev or leave blank to proxy"
                        class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm min-h-[42px]"
                      />
                    </div>
                  </div>
                </div>
              </div>

              <!-- 3. Generated Manga Library -->
              <div class="space-y-3 pt-5 border-t border-[#D1C7BD]/60">
                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                  <div>
                    <div class="flex items-center gap-2">
                      <h3 class="font-serif text-base font-bold text-[#322D29]">
                        Generated Manga Library
                      </h3>
                      <span class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-[#EFE9E1] text-[#72383D]">
                        {{ generatedMangaList.length }}
                      </span>
                    </div>
                    <p class="text-xs text-[#322D29]/60 mt-0.5">
                      Regenerate or delete stored daily stories from PostgreSQL and R2.
                    </p>
                  </div>
                  <button
                    type="button"
                    @click="loadMangaList"
                    :disabled="loadingMangaList"
                    class="admin-inline-btn inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-[#FAF7EE] hover:bg-[#EFE9E1] border border-[#D1C7BD] text-xs font-semibold text-[#5C282D] transition-colors cursor-pointer disabled:opacity-50 self-start sm:self-auto"
                  >
                    {{ loadingMangaList ? 'Refreshing...' : 'Refresh' }}
                  </button>
                </div>

                <div
                  v-if="mangaActionMessage"
                  class="p-3 rounded-xl text-xs font-medium border"
                  :class="mangaActionMessage.type === 'success' ? 'bg-emerald-50 border-emerald-200 text-emerald-800' : 'bg-red-50 border-red-200 text-red-800'"
                >
                  {{ mangaActionMessage.text }}
                </div>

                <div
                  v-if="loadingMangaList && generatedMangaList.length === 0"
                  class="p-5 rounded-2xl bg-white/60 border border-[#D1C7BD]/60 text-xs text-[#7B6E65] text-center"
                >
                  Loading generated manga stories...
                </div>

                <div
                  v-else-if="generatedMangaList.length === 0"
                  class="p-5 rounded-2xl bg-white/60 border border-[#D1C7BD]/60 text-xs text-[#7B6E65] text-center"
                >
                  No generated manga stories saved yet.
                </div>

                <div v-else class="space-y-2.5 max-h-[420px] overflow-y-auto scrollbar-thin pr-1">
                  <div
                    v-for="item in generatedMangaList"
                    :key="item.id"
                    class="p-3.5 rounded-2xl bg-white/85 border border-[#D1C7BD]/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
                  >
                    <div class="flex items-start gap-3 min-w-0 flex-1">
                      <div class="w-14 h-14 rounded-xl bg-[#EFE9E1] border border-[#D1C7BD] overflow-hidden shrink-0 flex items-center justify-center">
                        <img
                          v-if="item.previewImage"
                          :src="resolveMangaPreviewUrl(item.previewImage)"
                          :alt="item.title"
                          class="w-full h-full object-cover"
                        />
                        <span v-else class="text-lg">🎨</span>
                      </div>

                      <div class="space-y-1 min-w-0 flex-1">
                        <div class="flex flex-wrap items-center gap-1.5">
                          <span class="font-serif font-bold text-sm text-[#322D29] truncate">
                            {{ item.title || `Daily Story (${item.date})` }}
                          </span>
                          <span class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-[#EFE9E1] text-[#72383D]">
                            {{ item.style }}
                          </span>
                          <span
                            class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider"
                            :class="item.r2Stored ? 'bg-emerald-100 text-emerald-800' : 'bg-stone-200 text-stone-700'"
                          >
                            {{ item.r2Stored ? `R2 (${item.r2Keys?.length || 0})` : 'DB Only' }}
                          </span>
                        </div>

                        <div class="flex flex-wrap items-center gap-x-2.5 gap-y-0.5 text-xs text-[#5C5149]">
                          <span class="font-mono font-semibold text-[#72383D]">{{ item.date }}</span>
                          <span>·</span>
                          <span>{{ item.panelCount }} panels · {{ item.characterCount }} chars</span>
                          <span v-if="item.imageModel" class="font-mono text-[11px] text-[#322D29]/60">{{ item.imageModel }}</span>
                        </div>

                        <p v-if="item.moralLesson" class="text-xs text-[#6B5F56] line-clamp-1 font-serif italic">
                          {{ item.moralLesson }}
                        </p>
                      </div>
                    </div>

                    <!-- Actions -->
                    <div class="flex flex-wrap items-center gap-2 shrink-0">
                      <button
                        type="button"
                        @click="regenerateMangaItem(item)"
                        :disabled="regeneratingMangaIds.has(item.id) || deletingMangaIds.has(item.id)"
                        class="admin-inline-btn px-3 py-1.5 rounded-xl bg-[#72383D] hover:bg-[#5C282D] text-[#FAF7F2] text-xs font-semibold transition-colors disabled:opacity-50 cursor-pointer"
                      >
                        {{ regeneratingMangaIds.has(item.id) ? 'Regenerating...' : 'Regenerate' }}
                      </button>

                      <template v-if="confirmDeleteMangaId === item.id">
                        <button
                          type="button"
                          @click="deleteMangaItem(item)"
                          :disabled="deletingMangaIds.has(item.id)"
                          class="admin-inline-btn px-3 py-1.5 rounded-xl bg-red-700 hover:bg-red-800 text-white text-xs font-bold transition-colors disabled:opacity-50 cursor-pointer"
                        >
                          {{ deletingMangaIds.has(item.id) ? 'Deleting...' : 'Confirm' }}
                        </button>
                        <button
                          type="button"
                          @click="confirmDeleteMangaId = null"
                          class="admin-inline-btn px-2.5 py-1.5 rounded-xl bg-[#EFE9E1] hover:bg-[#E2D8CC] text-xs font-semibold text-[#322D29] transition-colors cursor-pointer"
                        >
                          Cancel
                        </button>
                      </template>
                      <button
                        v-else
                        type="button"
                        @click="confirmDeleteMangaId = item.id"
                        :disabled="regeneratingMangaIds.has(item.id) || deletingMangaIds.has(item.id)"
                        class="admin-inline-btn px-3 py-1.5 rounded-xl border border-red-200 text-red-700 hover:bg-red-50 text-xs font-semibold transition-colors disabled:opacity-50 cursor-pointer"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <!-- TAB 3: ACCESS CODES & GUEST REQUESTS -->
            <section
              v-else-if="activeTab === 'access'"
              key="access"
              id="panel-access"
              role="tabpanel"
              class="space-y-6"
            >
              <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                <div>
                  <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                    Guest Access &amp; Invites
                  </h2>
                  <p class="text-xs text-[#322D29]/65 mt-0.5">
                    Create invite codes and review guest access requests for Bible Study &amp; Magisterium AI.
                  </p>
                </div>
                <button
                  type="button"
                  @click="loadAccessData"
                  :disabled="loadingAccess"
                  class="admin-inline-btn px-3 py-1.5 rounded-xl bg-[#FAF7EE] hover:bg-[#EFE9E1] border border-[#D1C7BD] text-xs font-semibold text-[#5C282D] transition-colors cursor-pointer self-start sm:self-auto"
                >
                  {{ loadingAccess ? 'Refreshing...' : 'Refresh' }}
                </button>
              </div>

              <div
                v-if="accessMessage"
                class="p-3 rounded-xl text-xs font-medium border"
                :class="accessMessage.type === 'success'
                  ? 'bg-emerald-50 border-emerald-200 text-emerald-800'
                  : 'bg-[#9B2C2C]/10 border-[#9B2C2C]/25 text-[#822020]'"
              >
                {{ accessMessage.text }}
              </div>

              <!-- Compact Inline Generate New Access Code Bar -->
              <div class="p-4 rounded-2xl bg-[#FAF7EE]/85 border border-[#D1C7BD]/80 space-y-3">
                <h3 class="font-serif text-sm font-bold text-[#4A2525]">Create Access Code</h3>
                <div class="grid grid-cols-1 sm:grid-cols-12 gap-2.5 items-end">
                  <div class="sm:col-span-4">
                    <label class="block text-[10px] font-bold uppercase tracking-[0.1em] text-[#72383D] mb-1">
                      Label / Guest Name
                    </label>
                    <input
                      v-model="newCodeLabel"
                      type="text"
                      placeholder="e.g. Parish Group"
                      class="w-full px-3 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm text-[#322D29] min-h-[42px] focus:outline-none focus:ring-2 focus:ring-[#72383D]/25"
                    />
                  </div>
                  <div class="sm:col-span-4">
                    <label class="block text-[10px] font-bold uppercase tracking-[0.1em] text-[#72383D] mb-1">
                      Custom Code <span class="font-normal text-[#7B6E65] normal-case">(optional)</span>
                    </label>
                    <input
                      v-model="newCustomCode"
                      type="text"
                      placeholder="Auto if blank"
                      class="w-full px-3 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm font-mono uppercase text-[#322D29] min-h-[42px] focus:outline-none focus:ring-2 focus:ring-[#72383D]/25"
                    />
                  </div>
                  <div class="sm:col-span-2">
                    <label class="block text-[10px] font-bold uppercase tracking-[0.1em] text-[#72383D] mb-1">
                      Daily Limit
                    </label>
                    <input
                      v-model.number="newDailyLimit"
                      type="number"
                      min="1"
                      max="500"
                      class="w-full px-3 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm text-[#322D29] min-h-[42px] focus:outline-none focus:ring-2 focus:ring-[#72383D]/25"
                    />
                  </div>
                  <div class="sm:col-span-2">
                    <AppButton
                      variant="primary"
                      customClass="w-full !py-2 !px-3 text-xs min-h-[42px]"
                      :disabled="creatingCode"
                      @click="generateCode"
                    >
                      {{ creatingCode ? '...' : 'Generate' }}
                    </AppButton>
                  </div>
                </div>
              </div>

              <!-- Guest Access Requests -->
              <div class="space-y-2.5">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <h3 class="font-serif text-base font-bold text-[#322D29]">
                      Guest Requests
                    </h3>
                    <span
                      v-if="pendingRequestsCount > 0"
                      class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-amber-100 text-amber-900"
                    >
                      {{ pendingRequestsCount }} pending
                    </span>
                  </div>
                  <span class="text-xs text-[#7B6E65]">{{ accessRequests.length }} total</span>
                </div>

                <div v-if="accessRequests.length === 0" class="p-4 rounded-2xl bg-white/60 border border-[#D1C7BD]/60 text-xs text-[#7B6E65] text-center">
                  No guest access requests yet.
                </div>

                <div v-else class="space-y-2 max-h-[340px] overflow-y-auto scrollbar-thin pr-1">
                  <div
                    v-for="req in accessRequests"
                    :key="req.id"
                    class="p-3.5 rounded-2xl bg-white/85 border border-[#D1C7BD]/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
                  >
                    <div class="space-y-1 min-w-0">
                      <div class="flex flex-wrap items-center gap-2">
                        <span class="font-semibold text-sm text-[#322D29]">{{ req.name }}</span>
                        <span v-if="req.email" class="text-xs text-[#7B6E65]">({{ req.email }})</span>
                        <span
                          class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider"
                          :class="{
                            'bg-amber-100 text-amber-800': req.status === 'pending',
                            'bg-emerald-100 text-emerald-800': req.status === 'approved',
                            'bg-red-100 text-red-800': req.status === 'rejected',
                          }"
                        >
                          {{ req.status }}
                        </span>
                      </div>
                      <p v-if="req.reason" class="text-xs text-[#5C5149] leading-relaxed">{{ req.reason }}</p>
                      <div v-if="req.generated_code" class="flex flex-wrap items-center gap-2 pt-0.5">
                        <span class="font-mono text-xs font-bold px-2 py-0.5 rounded bg-[#FAF7EE] border border-[#D1C7BD] text-[#72383D]">
                          {{ req.generated_code }}
                        </span>
                        <button
                          type="button"
                          @click="copyText(getInviteLink(req.generated_code), `req-link-${req.id}`)"
                          class="admin-inline-btn bg-transparent border-transparent px-1 py-0.5 text-[11px] font-semibold text-[#72383D] hover:underline hover:bg-transparent cursor-pointer"
                        >
                          {{ copiedKey === `req-link-${req.id}` ? 'Invite Link Copied!' : 'Copy Invite Link' }}
                        </button>
                      </div>
                    </div>

                    <div class="flex items-center gap-2 shrink-0">
                      <button
                        v-if="req.status === 'pending'"
                        type="button"
                        @click="approveRequest(req)"
                        class="admin-inline-btn px-3 py-1.5 rounded-xl bg-[#72383D] hover:bg-[#5C282D] text-[#FAF7F2] text-xs font-semibold transition-colors cursor-pointer"
                      >
                        Approve
                      </button>
                      <button
                        v-if="req.status === 'pending'"
                        type="button"
                        @click="rejectRequest(req)"
                        class="admin-inline-btn px-3 py-1.5 rounded-xl bg-[#EFE9E1] hover:bg-[#E2D8CC] text-[#5C282D] text-xs font-semibold transition-colors cursor-pointer"
                      >
                        Reject
                      </button>
                      <button
                        type="button"
                        @click="deleteRequest(req.id)"
                        class="admin-inline-btn bg-transparent border-transparent px-2.5 py-1.5 rounded-xl text-xs text-red-700 hover:bg-red-50 transition-colors cursor-pointer"
                        title="Delete request"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Active & Revoked Access Codes -->
              <div class="space-y-2.5 pt-4 border-t border-[#D1C7BD]/60">
                <div class="flex items-center justify-between">
                  <h3 class="font-serif text-base font-bold text-[#322D29]">
                    Access Codes
                  </h3>
                  <span class="text-xs text-[#7B6E65]">{{ accessCodes.length }} total</span>
                </div>

                <div v-if="accessCodes.length === 0" class="p-4 rounded-2xl bg-white/60 border border-[#D1C7BD]/60 text-xs text-[#7B6E65] text-center">
                  No access codes created yet. Generate one above.
                </div>

                <div v-else class="space-y-2 max-h-[380px] overflow-y-auto scrollbar-thin pr-1">
                  <div
                    v-for="item in accessCodes"
                    :key="item.id"
                    class="p-3.5 rounded-2xl bg-white/85 border border-[#D1C7BD]/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
                  >
                    <div class="space-y-1 min-w-0">
                      <div class="flex flex-wrap items-center gap-2">
                        <span class="font-mono text-sm font-bold text-[#72383D] tracking-wide">
                          {{ item.code }}
                        </span>
                        <span
                          class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider"
                          :class="item.is_active ? 'bg-emerald-100 text-emerald-800' : 'bg-stone-200 text-stone-600'"
                        >
                          {{ item.is_active ? 'Active' : 'Revoked' }}
                        </span>
                        <span class="text-xs font-medium text-[#5C5149] bg-[#FAF7EE] px-2 py-0.5 rounded-full border border-[#E5DDD3]">
                          {{ item.used_today }} / {{ item.daily_limit }} today
                        </span>
                      </div>
                      <p v-if="item.label" class="text-xs text-[#6B5F56]">{{ item.label }}</p>
                    </div>

                    <div class="flex flex-wrap items-center gap-1.5 shrink-0">
                      <button
                        type="button"
                        @click="copyText(item.code, `code-${item.id}`)"
                        class="admin-inline-btn px-2.5 py-1.5 rounded-xl bg-[#FAF7EE] hover:bg-[#EFE9E1] border border-[#D1C7BD] text-xs font-semibold text-[#5C282D] transition-colors cursor-pointer"
                      >
                        {{ copiedKey === `code-${item.id}` ? 'Copied!' : 'Copy' }}
                      </button>
                      <button
                        type="button"
                        @click="copyText(getInviteLink(item.code), `link-${item.id}`)"
                        class="admin-inline-btn px-2.5 py-1.5 rounded-xl bg-[#72383D]/10 hover:bg-[#72383D]/20 border border-[#72383D]/25 text-xs font-semibold text-[#72383D] transition-colors cursor-pointer"
                      >
                        {{ copiedKey === `link-${item.id}` ? 'Link Copied!' : 'Invite Link' }}
                      </button>
                      <button
                        type="button"
                        @click="toggleCodeActive(item)"
                        class="admin-inline-btn px-2.5 py-1.5 rounded-xl bg-[#EFE9E1] hover:bg-[#E2D8CC] text-xs font-semibold text-[#322D29] transition-colors cursor-pointer"
                      >
                        {{ item.is_active ? 'Revoke' : 'Activate' }}
                      </button>
                      <button
                        type="button"
                        @click="deleteCode(item.id)"
                        class="admin-inline-btn bg-transparent border-transparent px-2.5 py-1.5 rounded-xl text-xs text-red-700 hover:bg-red-50 transition-colors cursor-pointer"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <!-- TAB 4: AI USAGE & COST -->
            <section
              v-else-if="activeTab === 'usage'"
              key="usage"
              id="panel-usage"
              role="tabpanel"
              class="space-y-5"
            >
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2.5">
                <div>
                  <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">AI Usage &amp; Cost</h2>
                  <p class="text-xs text-[#322D29]/60 mt-0.5">
                    Estimated from Gemini token metadata &amp; per-model rates. Local Ollama calls are $0.00.
                  </p>
                </div>
                <button
                  type="button"
                  @click="loadAIUsage"
                  :disabled="loadingUsage"
                  class="admin-inline-btn px-3 py-1.5 rounded-xl bg-[#FAF7EE] hover:bg-[#EFE9E1] border border-[#D1C7BD] text-xs font-semibold text-[#5C282D] transition-colors cursor-pointer disabled:opacity-50 self-start sm:self-auto"
                >
                  {{ loadingUsage ? 'Refreshing...' : 'Refresh' }}
                </button>
              </div>

              <!-- Compact Period Filter & Reference Date -->
              <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-3 rounded-2xl bg-[#FAF8F5] border border-[#D1C7BD]/70">
                <div class="inline-flex items-center gap-1 bg-white/80 p-1 rounded-xl border border-[#D1C7BD]/70 self-start">
                  <button
                    v-for="p in (['day', 'week', 'month', 'year'] as const)"
                    :key="p"
                    type="button"
                    @click="setUsagePeriod(p)"
                    class="admin-seg-btn min-h-[34px] px-3 py-1 rounded-lg text-xs font-semibold capitalize transition-all cursor-pointer border-0"
                    :class="
                      usagePeriod === p
                        ? 'bg-[#72383D] hover:bg-[#72383D] text-white shadow-2xs'
                        : 'bg-transparent text-[#322D29]/75 hover:bg-[#EFE9E1]/80 hover:text-[#72383D]'
                    "
                  >
                    {{ p }}
                  </button>
                </div>

                <div class="flex flex-wrap items-center gap-2.5 text-xs">
                  <span v-if="usageReport" class="text-[11px] font-mono text-[#322D29]/60">
                    {{ usageReport.start_date }} → {{ usageReport.end_date }}
                  </span>
                  <input
                    type="date"
                    v-model="usageAnchorDate"
                    @change="loadAIUsage"
                    aria-label="Reference date"
                    class="px-3 py-1.5 rounded-xl border border-[#D1C7BD] bg-white text-xs text-[#322D29] focus:outline-none focus:border-[#72383D]"
                  />
                </div>
              </div>

              <div
                v-if="usageError"
                class="p-3.5 rounded-xl bg-red-50 border border-red-200 text-xs text-red-700"
              >
                {{ usageError }}
              </div>

              <!-- Summary KPI Cards -->
              <div v-if="usageReport" class="grid grid-cols-2 lg:grid-cols-4 gap-3">
                <div class="p-3.5 rounded-2xl bg-white/75 border border-[#D1C7BD]/75">
                  <div class="text-[10px] font-bold uppercase tracking-wider text-[#322D29]/55">Estimated Cost</div>
                  <div class="mt-1 font-serif text-xl sm:text-2xl font-bold text-[#72383D]">
                    {{ formatUSD(computedTotalCostUSD) }}
                  </div>
                  <div class="mt-0.5 text-[11px] text-[#322D29]/50">Live from rates</div>
                </div>

                <div class="p-3.5 rounded-2xl bg-white/75 border border-[#D1C7BD]/75">
                  <div class="text-[10px] font-bold uppercase tracking-wider text-[#322D29]/55">AI Requests</div>
                  <div class="mt-1 font-serif text-xl sm:text-2xl font-bold text-[#322D29]">
                    {{ formatNumber(usageReport.summary.total_requests) }}
                  </div>
                  <div class="mt-0.5 text-[11px] text-[#322D29]/50">Logged calls</div>
                </div>

                <div class="p-3.5 rounded-2xl bg-white/75 border border-[#D1C7BD]/75">
                  <div class="text-[10px] font-bold uppercase tracking-wider text-[#322D29]/55">Total Tokens</div>
                  <div class="mt-1 font-serif text-xl sm:text-2xl font-bold text-[#322D29]">
                    {{ formatNumber(usageReport.summary.total_tokens) }}
                  </div>
                  <div class="mt-0.5 text-[11px] text-[#322D29]/50 truncate">
                    {{ formatNumber(usageReport.summary.total_input_tokens) }} in / {{ formatNumber(usageReport.summary.total_output_tokens) }} out
                  </div>
                </div>

                <div class="p-3.5 rounded-2xl bg-white/75 border border-[#D1C7BD]/75">
                  <div class="text-[10px] font-bold uppercase tracking-wider text-[#322D29]/55">Images</div>
                  <div class="mt-1 font-serif text-xl sm:text-2xl font-bold text-[#322D29]">
                    {{ formatNumber(usageReport.summary.total_images) }}
                  </div>
                  <div class="mt-0.5 text-[11px] text-[#322D29]/50">Manga &amp; magazine</div>
                </div>
              </div>

              <!-- Inner Usage Sub-Views (keeps screen clean & compact) -->
              <div class="inline-flex items-center gap-1 p-1 rounded-xl bg-[#FAF8F5] border border-[#D1C7BD]/75">
                <button
                  type="button"
                  @click="usageSubTab = 'overview'"
                  class="admin-seg-btn min-h-[34px] px-3 py-1 rounded-lg text-xs font-semibold transition-all cursor-pointer border-0"
                  :class="usageSubTab === 'overview' ? 'bg-white hover:bg-white text-[#72383D] shadow-2xs' : 'bg-transparent text-[#322D29]/65 hover:bg-white/60 hover:text-[#322D29]'"
                >
                  Overview
                </button>
                <button
                  type="button"
                  @click="usageSubTab = 'activity'"
                  class="admin-seg-btn min-h-[34px] px-3 py-1 rounded-lg text-xs font-semibold transition-all cursor-pointer border-0 flex items-center gap-1.5"
                  :class="usageSubTab === 'activity' ? 'bg-white hover:bg-white text-[#72383D] shadow-2xs' : 'bg-transparent text-[#322D29]/65 hover:bg-white/60 hover:text-[#322D29]'"
                >
                  Activity
                  <span
                    v-if="usageReport?.recent_logs?.length"
                    class="text-[10px] px-1.5 py-0.5 rounded-full bg-[#EFE9E1] text-[#72383D] leading-none"
                  >
                    {{ usageReport.recent_logs.length }}
                  </span>
                </button>
                <button
                  type="button"
                  @click="usageSubTab = 'pricing'"
                  class="admin-seg-btn min-h-[34px] px-3 py-1 rounded-lg text-xs font-semibold transition-all cursor-pointer border-0"
                  :class="usageSubTab === 'pricing' ? 'bg-white hover:bg-white text-[#72383D] shadow-2xs' : 'bg-transparent text-[#322D29]/65 hover:bg-white/60 hover:text-[#322D29]'"
                >
                  Pricing &amp; Cleanup
                </button>
              </div>

              <!-- Sub-View 1: Overview (Timeline + Breakdown) -->
              <div v-if="usageSubTab === 'overview'" class="space-y-4">
                <!-- Timeline Bar Chart -->
                <div v-if="usageReport && usageReport.timeline.length > 0" class="p-4 rounded-2xl bg-[#FAF8F5] border border-[#D1C7BD]/70 space-y-3">
                  <h3 class="text-xs font-bold uppercase tracking-wider text-[#72383D]">Cost Timeline</h3>
                  <div class="space-y-2 max-h-[260px] overflow-y-auto scrollbar-thin pr-1">
                    <div
                      v-for="pt in usageReport.timeline"
                      :key="pt.bucket"
                      class="flex items-center gap-3 text-xs"
                    >
                      <span class="w-24 font-mono text-[11px] text-[#322D29]/75 shrink-0">{{ pt.bucket }}</span>
                      <div class="flex-1 h-4 rounded-lg bg-[#EFE9E1] overflow-hidden">
                        <div
                          class="h-full bg-[#72383D] rounded-lg transition-all"
                          :style="{
                            width: maxTimelineCost > 0 ? `${Math.max((pt.estimated_cost_usd / maxTimelineCost) * 100, pt.requests > 0 ? 4 : 0)}%` : '4%'
                          }"
                        ></div>
                      </div>
                      <span class="w-20 text-right font-mono text-[11px] font-semibold text-[#322D29]">
                        {{ formatUSD(pt.estimated_cost_usd) }}
                      </span>
                      <span class="w-14 text-right text-[11px] text-[#322D29]/55">
                        {{ pt.requests }} req
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Breakdown by Provider / Model & Feature -->
                <div v-if="usageReport" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                  <div class="p-4 rounded-2xl bg-[#FAF8F5] border border-[#D1C7BD]/70 space-y-2.5">
                    <h3 class="text-xs font-bold uppercase tracking-wider text-[#72383D]">By Provider &amp; Model</h3>
                    <div v-if="computedByProvider.length === 0" class="text-xs text-[#322D29]/55 py-4 text-center">
                      No usage recorded for this period.
                    </div>
                    <div v-else class="space-y-2">
                      <div
                        v-for="row in computedByProvider"
                        :key="`${row.provider}-${row.model}`"
                        class="flex items-center justify-between p-2.5 rounded-xl bg-white border border-[#D1C7BD]/60 text-xs"
                      >
                        <div class="min-w-0">
                          <div class="font-semibold text-[#322D29] truncate">
                            <span class="uppercase text-[10px] px-1.5 py-0.5 rounded bg-[#EFE9E1] text-[#72383D] mr-1.5">{{ row.provider }}</span>
                            {{ row.model || 'default' }}
                          </div>
                          <div class="text-[11px] text-[#322D29]/60 mt-0.5">
                            {{ row.requests }} req · {{ formatNumber(row.total_tokens) }} tok
                            <span v-if="row.image_count > 0"> · {{ row.image_count }} img</span>
                          </div>
                        </div>
                        <div class="font-mono font-semibold text-[#72383D] shrink-0 ml-2">
                          {{ formatUSD(row.computed_cost_usd) }}
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="p-4 rounded-2xl bg-[#FAF8F5] border border-[#D1C7BD]/70 space-y-2.5">
                    <h3 class="text-xs font-bold uppercase tracking-wider text-[#72383D]">By Feature</h3>
                    <div v-if="usageReport.by_feature.length === 0" class="text-xs text-[#322D29]/55 py-4 text-center">
                      No usage recorded for this period.
                    </div>
                    <div v-else class="space-y-2">
                      <div
                        v-for="row in usageReport.by_feature"
                        :key="row.feature"
                        class="flex items-center justify-between p-2.5 rounded-xl bg-white border border-[#D1C7BD]/60 text-xs"
                      >
                        <div>
                          <div class="font-semibold text-[#322D29]">{{ row.feature }}</div>
                          <div class="text-[11px] text-[#322D29]/60 mt-0.5">
                            {{ row.requests }} req · {{ formatNumber(row.total_tokens) }} tok
                            <span v-if="row.image_count > 0"> · {{ row.image_count }} img</span>
                          </div>
                        </div>
                        <div class="font-mono font-semibold text-[#72383D]">
                          {{ formatUSD(row.estimated_cost_usd) }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Sub-View 2: Recent Activity (Latest 100) -->
              <div v-else-if="usageSubTab === 'activity'" class="p-4 rounded-2xl bg-[#FAF8F5] border border-[#D1C7BD]/70 space-y-3">
                <h3 class="text-xs font-bold uppercase tracking-wider text-[#72383D]">Recent Activity (Latest 100)</h3>
                <div v-if="!usageReport || usageReport.recent_logs.length === 0" class="text-xs text-[#322D29]/55 py-6 text-center">
                  No AI usage logs in this period yet.
                </div>
                <div v-else class="overflow-x-auto max-h-[440px] overflow-y-auto scrollbar-thin">
                  <table class="w-full text-left border-collapse text-xs">
                    <thead class="sticky top-0 bg-[#FAF8F5] z-10">
                      <tr class="border-b border-[#D1C7BD] text-[11px] text-[#322D29]/65">
                        <th class="py-2 pr-3 font-semibold">Time</th>
                        <th class="py-2 pr-3 font-semibold">Provider / Model</th>
                        <th class="py-2 pr-3 font-semibold">Feature</th>
                        <th class="py-2 pr-3 font-semibold text-right">In / Out</th>
                        <th class="py-2 pr-3 font-semibold text-right">Img</th>
                        <th class="py-2 font-semibold text-right">Est. Cost</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="log in usageReport.recent_logs"
                        :key="log.id"
                        class="border-b border-[#D1C7BD]/40 last:border-0"
                      >
                        <td class="py-2 pr-3 font-mono text-[11px] text-[#322D29]/70 whitespace-nowrap">
                          {{ new Date(log.created_at).toLocaleString() }}
                        </td>
                        <td class="py-2 pr-3 whitespace-nowrap">
                          <span class="uppercase text-[10px] px-1.5 py-0.5 rounded bg-[#EFE9E1] text-[#72383D] mr-1">{{ log.provider }}</span>
                          <span class="text-[#322D29]">{{ log.model }}</span>
                        </td>
                        <td class="py-2 pr-3 text-[#322D29]/80">
                          {{ log.feature }}<span v-if="log.operation" class="text-[#322D29]/50"> ({{ log.operation }})</span>
                        </td>
                        <td class="py-2 pr-3 text-right font-mono text-[11px] text-[#322D29]/80 whitespace-nowrap">
                          {{ formatNumber(log.input_tokens) }} / {{ formatNumber(log.output_tokens) }}
                        </td>
                        <td class="py-2 pr-3 text-right font-mono text-[11px]">
                          {{ log.image_count }}
                        </td>
                        <td class="py-2 text-right font-mono font-semibold text-[#72383D] whitespace-nowrap">
                          {{
                            formatUSD(
                              computeCostFromRates(
                                log.provider,
                                log.model,
                                log.input_tokens,
                                log.output_tokens,
                                log.thoughts_tokens,
                                log.cached_tokens,
                                log.image_count,
                                log.estimated_cost_usd
                              )
                            )
                          }}
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Sub-View 3: Model Pricing & Database Cleanup -->
              <div v-else class="space-y-4">
                <!-- Model Pricing Settings ($ / 1M Tokens) -->
                <div class="p-4 rounded-2xl bg-[#FAF8F5] border border-[#D1C7BD]/70 space-y-4">
                  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
                    <div>
                      <h3 class="text-sm font-bold text-[#72383D]">Model Pricing ($ per 1M Tokens)</h3>
                      <p class="text-xs text-[#322D29]/65 mt-0.5">
                        Set USD rates per model. Displayed costs update live from these rates.
                      </p>
                    </div>
                    <button
                      type="button"
                      @click="saveAIPricing"
                      :disabled="savingPricing || loadingPricing"
                      class="admin-inline-btn px-4 py-2 rounded-xl bg-[#72383D] hover:bg-[#5D2C30] text-white text-xs font-semibold transition-colors cursor-pointer disabled:opacity-50 self-start sm:self-auto"
                    >
                      {{ savingPricing ? 'Saving...' : 'Save Pricing' }}
                    </button>
                  </div>

                  <div v-if="modelPricingList.length > 0" class="overflow-x-auto scrollbar-thin">
                    <table class="w-full text-left border-collapse text-xs">
                      <thead>
                        <tr class="border-b border-[#D1C7BD] text-[11px] text-[#322D29]/65">
                          <th class="py-2 pr-2 font-semibold">Provider</th>
                          <th class="py-2 pr-2 font-semibold">Model</th>
                          <th class="py-2 pr-2 font-semibold">Input ($/1M)</th>
                          <th class="py-2 pr-2 font-semibold">Output ($/1M)</th>
                          <th class="py-2 pr-2 font-semibold">Cached ($/1M)</th>
                          <th class="py-2 pr-2 font-semibold">Image ($/img)</th>
                          <th class="py-2 font-semibold text-right">Action</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr
                          v-for="rule in modelPricingList"
                          :key="`${rule.provider}-${rule.model}`"
                          class="border-b border-[#D1C7BD]/40 last:border-0"
                        >
                          <td class="py-2 pr-2">
                            <span class="uppercase text-[10px] px-1.5 py-0.5 rounded bg-[#EFE9E1] text-[#72383D]">{{ rule.provider }}</span>
                          </td>
                          <td class="py-2 pr-2 font-mono text-[11px] text-[#322D29] whitespace-nowrap">{{ rule.model }}</td>
                          <td class="py-2 pr-2">
                            <input
                              type="number"
                              step="0.001"
                              min="0"
                              v-model.number="rule.input_per_million"
                              class="w-24 px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white font-mono text-[11px] text-[#322D29] focus:outline-none focus:border-[#72383D]"
                            />
                          </td>
                          <td class="py-2 pr-2">
                            <input
                              type="number"
                              step="0.001"
                              min="0"
                              v-model.number="rule.output_per_million"
                              class="w-24 px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white font-mono text-[11px] text-[#322D29] focus:outline-none focus:border-[#72383D]"
                            />
                          </td>
                          <td class="py-2 pr-2">
                            <input
                              type="number"
                              step="0.0001"
                              min="0"
                              v-model.number="rule.cached_per_million"
                              class="w-24 px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white font-mono text-[11px] text-[#322D29] focus:outline-none focus:border-[#72383D]"
                            />
                          </td>
                          <td class="py-2 pr-2">
                            <input
                              type="number"
                              step="0.001"
                              min="0"
                              v-model.number="rule.image_cost_each"
                              class="w-20 px-2 py-1 rounded-lg border border-[#D1C7BD] bg-white font-mono text-[11px] text-[#322D29] focus:outline-none focus:border-[#72383D]"
                            />
                          </td>
                          <td class="py-2 text-right">
                            <button
                              type="button"
                              @click="deletePricingRule(rule)"
                              class="admin-inline-btn bg-transparent border-transparent px-2 py-1 rounded-lg text-[11px] text-red-700 hover:bg-red-50 transition-colors cursor-pointer"
                            >
                              Remove
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>

                  <!-- Add New Model Pricing Row -->
                  <div class="pt-3 border-t border-[#D1C7BD]/60 flex flex-wrap items-end gap-2 text-xs">
                    <div>
                      <label class="block text-[10px] font-semibold text-[#322D29]/65 mb-1">Provider</label>
                      <select
                        v-model="newPricingRule.provider"
                        class="px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#322D29]"
                      >
                        <option value="gemini">gemini</option>
                        <option value="ollama">ollama</option>
                        <option value="magisterium">magisterium</option>
                      </select>
                    </div>
                    <div class="flex-1 min-w-[160px]">
                      <label class="block text-[10px] font-semibold text-[#322D29]/65 mb-1">Model Name</label>
                      <input
                        type="text"
                        v-model="newPricingRule.model"
                        placeholder="e.g. gemini-3.1-flash-lite"
                        class="w-full px-2.5 py-1.5 rounded-lg border border-[#D1C7BD] bg-white text-xs text-[#322D29]"
                      />
                    </div>
                    <div>
                      <label class="block text-[10px] font-semibold text-[#322D29]/65 mb-1">In ($/1M)</label>
                      <input
                        type="number"
                        step="0.001"
                        min="0"
                        v-model.number="newPricingRule.input_per_million"
                        class="w-20 px-2 py-1.5 rounded-lg border border-[#D1C7BD] bg-white font-mono text-xs text-[#322D29]"
                      />
                    </div>
                    <div>
                      <label class="block text-[10px] font-semibold text-[#322D29]/65 mb-1">Out ($/1M)</label>
                      <input
                        type="number"
                        step="0.001"
                        min="0"
                        v-model.number="newPricingRule.output_per_million"
                        class="w-20 px-2 py-1.5 rounded-lg border border-[#D1C7BD] bg-white font-mono text-xs text-[#322D29]"
                      />
                    </div>
                    <div>
                      <label class="block text-[10px] font-semibold text-[#322D29]/65 mb-1">Img ($)</label>
                      <input
                        type="number"
                        step="0.001"
                        min="0"
                        v-model.number="newPricingRule.image_cost_each"
                        class="w-16 px-2 py-1.5 rounded-lg border border-[#D1C7BD] bg-white font-mono text-xs text-[#322D29]"
                      />
                    </div>
                    <button
                      type="button"
                      @click="addPricingRule"
                      :disabled="savingPricing"
                      class="admin-inline-btn px-3.5 py-1.5 rounded-lg bg-[#EFE9E1] hover:bg-[#E2D8CC] text-xs font-semibold text-[#322D29] transition-colors cursor-pointer"
                    >
                      + Add Model
                    </button>
                  </div>

                  <div
                    v-if="pricingMessage"
                    class="p-2.5 rounded-xl text-xs"
                    :class="pricingMessage.type === 'success' ? 'bg-emerald-50 text-emerald-800 border border-emerald-200' : 'bg-red-50 text-red-700 border border-red-200'"
                  >
                    {{ pricingMessage.text }}
                  </div>
                </div>

                <!-- Database Cleanup / Purge Section -->
                <div class="p-4 rounded-2xl bg-[#FAF8F5] border border-red-200/80 space-y-3">
                  <div>
                    <h3 class="text-sm font-bold text-[#72383D]">Database Cleanup</h3>
                    <p class="text-xs text-[#322D29]/65 mt-0.5">
                      Purge all AI usage records strictly older than the selected cutoff date.
                    </p>
                  </div>
                  <div class="flex flex-col sm:flex-row sm:items-center gap-3">
                    <input
                      type="date"
                      v-model="cleanupBeforeDate"
                      class="px-3.5 py-2 rounded-xl border border-[#D1C7BD] bg-white text-xs text-[#322D29] focus:outline-none focus:border-[#72383D]"
                    />
                    <button
                      type="button"
                      @click="purgeAIUsage"
                      :disabled="cleaningUp || !cleanupBeforeDate"
                      class="admin-inline-btn px-4 py-2 rounded-xl bg-red-700 hover:bg-red-800 text-white text-xs font-semibold transition-colors cursor-pointer disabled:opacity-50"
                    >
                      {{ cleaningUp ? 'Purging...' : 'Purge Older Records' }}
                    </button>
                  </div>
                  <div
                    v-if="cleanupMessage"
                    class="p-2.5 rounded-xl text-xs"
                    :class="cleanupMessage.type === 'success' ? 'bg-emerald-50 text-emerald-800 border border-emerald-200' : 'bg-red-50 text-red-700 border border-red-200'"
                  >
                    {{ cleanupMessage.text }}
                  </div>
                </div>
              </div>
            </section>
          </Transition>

          <!-- Contextual Settings Action Footer (only shown on settings-backed tabs) -->
          <div
            v-if="isSettingsTab"
            class="mt-6 pt-4 border-t border-[#D1C7BD]/70 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3"
          >
            <div class="flex items-center gap-2 text-xs">
              <span v-if="loadingSettings" class="text-[#322D29]/60">
                Loading settings...
              </span>
              <template v-else-if="hasUnsavedChanges">
                <span class="w-2 h-2 rounded-full bg-amber-500 animate-pulse"></span>
                <span class="font-semibold text-[#72383D]">You have unsaved changes</span>
              </template>
              <template v-else>
                <span class="w-2 h-2 rounded-full bg-emerald-600/70"></span>
                <span class="text-[#322D29]/60">All settings synced with database</span>
              </template>
            </div>

            <div class="flex items-center justify-end gap-2.5">
              <AppButton
                v-if="hasUnsavedChanges"
                variant="outlined"
                customClass="!px-4 !py-2 text-xs min-h-[40px]"
                @click="discardChanges"
              >
                Discard
              </AppButton>
              <AppButton
                variant="primary"
                customClass="min-h-[40px] !py-2 text-xs sm:text-sm"
                @click="saveSettings"
                :disabled="saveLoading"
              >
                {{ saveLoading ? 'Saving...' : 'Save Settings' }}
              </AppButton>
            </div>
          </div>
        </div>

        <!-- Floating Unsaved Changes Pill (accessible from any tab or scroll position) -->
        <Transition name="save-bar">
          <div
            v-if="hasUnsavedChanges"
            class="fixed bottom-20 lg:bottom-6 left-1/2 -translate-x-1/2 z-40 w-[calc(100%-2rem)] max-w-lg px-4 py-2.5 rounded-2xl bg-[#322D29]/95 text-[#FAF7EE] backdrop-blur-xl border border-white/15 shadow-xl flex items-center justify-between gap-3"
          >
            <div class="flex items-center gap-2 text-xs sm:text-sm font-medium min-w-0">
              <span class="w-2 h-2 rounded-full bg-amber-400 animate-pulse shrink-0"></span>
              <span class="truncate">Unsaved settings changes</span>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <button
                type="button"
                @click="discardChanges"
                class="admin-inline-btn px-3 py-1.5 rounded-xl bg-white/10 hover:bg-white/20 text-xs font-semibold text-white/90 transition-colors cursor-pointer"
              >
                Discard
              </button>
              <button
                type="button"
                @click="saveSettings"
                :disabled="saveLoading"
                class="admin-inline-btn px-3.5 py-1.5 rounded-xl bg-[#72383D] hover:bg-[#8B464C] text-xs font-bold text-white transition-colors cursor-pointer disabled:opacity-50"
              >
                {{ saveLoading ? 'Saving...' : 'Save' }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </main>

    <BottomNav />
  </div>
</template>

<style scoped>
.tab-fade-enter-active,
.tab-fade-leave-active {
  transition: opacity 0.18s cubic-bezier(0.16, 1, 0.3, 1), transform 0.18s cubic-bezier(0.16, 1, 0.3, 1);
}

.tab-fade-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.tab-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

.save-bar-enter-active,
.save-bar-leave-active {
  transition: opacity 0.22s cubic-bezier(0.16, 1, 0.3, 1), transform 0.24s cubic-bezier(0.16, 1, 0.3, 1);
}

.save-bar-enter-from,
.save-bar-leave-to {
  opacity: 0;
  transform: translate(-50%, 14px) scale(0.97);
}

/* Keep admin segmented tabs and compact inline controls stable on hover */
.admin-seg-btn:hover,
.admin-inline-btn:hover {
  transform: none;
}

.admin-inline-btn {
  min-height: 0;
}
</style>
