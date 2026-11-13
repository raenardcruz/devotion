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

interface Settings {
  context_provider: string;
  fact_checker_provider: string;
  magisterium_llm_provider: string;
  gemini_api_key: string;
  magisterium_api_key: string;
  bible_api_key: string;
  gemini_model: string;
  gemini_live_model: string;
  ollama_model: string;
  context_model?: string;
  context_instruction: string;
  enable_fact_checker: boolean;
}

type AdminTabId = 'pipeline' | 'models' | 'keys' | 'prompts' | 'access';

const tabs: { id: AdminTabId; label: string; description: string }[] = [
  { id: 'pipeline', label: 'AI Pipeline', description: 'Providers & verification' },
  { id: 'models', label: 'Models', description: 'Gemini & local Ollama' },
  { id: 'keys', label: 'API Keys', description: 'Database credentials' },
  { id: 'prompts', label: 'Prompts', description: 'Context instructions' },
  { id: 'access', label: 'Access Codes', description: 'Guest requests & invites' },
];

const { fetchWithAuth, isServerUnreachable } = useDevotionApi();
const { syncAdminState } = useAccessControl();

// Navigation & UI State
const activeTab = ref<AdminTabId>('pipeline');
const showLoginPassword = ref(false);
const showKeys = ref({
  gemini: false,
  magisterium: false,
  bible: false,
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
  ollama_model: 'gemma4:cloud',
  context_instruction: '',
  enable_fact_checker: true,
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

onMounted(() => {
  const token = localStorage.getItem('adminToken');
  if (token) {
    isAuthenticated.value = true;
    syncAdminState();
    loadSettings();
    loadAccessData();
  }
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
      savedSnapshot.value = JSON.stringify(settings.value);
    }
    fetchGeminiModels();
    fetchGeminiLiveModels();
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
    await fetchWithAuth('/api/admin/settings', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(settings.value),
    });
    savedSnapshot.value = JSON.stringify(settings.value);
    message.value = { type: 'success', text: 'Settings successfully saved to PostgreSQL database.' };
  } catch (err: any) {
    message.value = { type: 'error', text: err.message || 'Failed to save settings to backend.' };
  } finally {
    saveLoading.value = false;
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
      <div v-else class="space-y-6">
        <!-- Dashboard Header -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-5 border-b border-[#D1C7BD]/70">
          <div>
            <div class="flex flex-wrap items-center gap-2.5">
              <h1 class="font-serif text-2xl sm:text-3xl font-bold text-[#322D29] tracking-tight">
                AI & System Settings
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
              Configure database-backed scripture generation, verification, models, and API credentials
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

        <!-- Tab Navigation Bar -->
        <div
          role="tablist"
          aria-label="Admin Settings Sections"
          class="flex items-center gap-2 overflow-x-auto no-scrollbar py-1 px-0.5 scroll-smooth"
        >
          <button
            v-for="tab in tabs"
            :key="tab.id"
            role="tab"
            :aria-selected="activeTab === tab.id"
            :aria-controls="`panel-${tab.id}`"
            @click="activeTab = tab.id"
            class="group flex items-center gap-2.5 px-4 py-2.5 rounded-2xl text-left transition-all duration-200 border shrink-0 min-h-[44px] cursor-pointer active:scale-[0.98]"
            :class="
              activeTab === tab.id
                ? 'bg-gradient-to-r from-[#72383D] to-[#322D29] text-white border-transparent shadow-sm shadow-[#72383D]/20'
                : 'bg-white/70 text-[#322D29]/80 border-[#D1C7BD]/80 hover:bg-white hover:text-[#72383D] hover:border-[#72383D]/40'
            "
          >
            <!-- Tab Icons -->
            <span class="flex-shrink-0">
              <svg v-if="tab.id === 'pipeline'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h10m-10 6h16" />
              </svg>
              <svg v-else-if="tab.id === 'models'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
              <svg v-else-if="tab.id === 'keys'" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.9">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </span>

            <div class="leading-tight">
              <div class="text-xs sm:text-sm font-bold tracking-wide flex items-center gap-1.5">
                {{ tab.label }}
                <span
                  v-if="tab.id === 'keys'"
                  class="text-[10px] px-1.5 py-0.2 rounded-full font-semibold"
                  :class="activeTab === tab.id ? 'bg-white/20 text-white' : 'bg-[#EFE9E1] text-[#72383D]'"
                >
                  {{ configuredKeysCount }}/3
                </span>
              </div>
              <div
                class="text-[10px] hidden sm:block"
                :class="activeTab === tab.id ? 'text-white/75' : 'text-[#322D29]/50'"
              >
                {{ tab.description }}
              </div>
            </div>
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
              class="text-current opacity-60 hover:opacity-100 p-1 rounded-lg cursor-pointer"
              aria-label="Dismiss notification"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </Transition>

        <!-- Main Tab Content Card -->
        <ParchmentCard class="p-5 sm:p-8">
          <Transition name="tab-fade" mode="out-in">
            <!-- TAB 1: AI PIPELINE -->
            <section
              v-if="activeTab === 'pipeline'"
              key="pipeline"
              id="panel-pipeline"
              role="tabpanel"
              class="space-y-6"
            >
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
                <div>
                  <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                    AI Generation & Fact-Checking Pipeline
                  </h2>
                  <p class="text-xs sm:text-sm text-[#322D29]/65 mt-0.5">
                    Route each stage of scripture context generation and doctrinal verification.
                  </p>
                </div>
              </div>

              <!-- Visual Pipeline Summary -->
              <div class="p-3.5 rounded-2xl bg-[#FAF7EE]/90 border border-[#D1C7BD]/70 flex flex-wrap items-center gap-2 text-xs">
                <span class="font-bold uppercase tracking-wider text-[10px] text-[#72383D] mr-1">Active Flow:</span>
                <span class="px-2.5 py-1 rounded-full bg-white border border-[#D1C7BD]/80 font-semibold text-[#322D29]">
                  1. {{ formatProviderLabel(settings.context_provider) }}
                </span>
                <span class="text-[#AC9C8D]">→</span>
                <span
                  class="px-2.5 py-1 rounded-full border font-semibold transition-opacity"
                  :class="
                    settings.enable_fact_checker
                      ? 'bg-white border-[#D1C7BD]/80 text-[#322D29]'
                      : 'bg-[#EFE9E1]/60 border-[#D1C7BD]/50 text-[#322D29]/45 line-through'
                  "
                >
                  2. {{ formatProviderLabel(settings.fact_checker_provider) }}
                </span>
                <span class="text-[#AC9C8D]">→</span>
                <span class="px-2.5 py-1 rounded-full bg-white border border-[#D1C7BD]/80 font-semibold text-[#322D29]">
                  3. {{ formatProviderLabel(settings.magisterium_llm_provider) }}
                </span>
              </div>

              <!-- 3-Stage Provider Grid -->
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <!-- 1. Context Provider -->
                <div class="p-4 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 flex flex-col justify-between gap-3">
                  <div>
                    <div class="flex items-center justify-between mb-1">
                      <span class="text-[10px] font-bold uppercase tracking-[0.12em] text-[#72383D]">Stage 1</span>
                      <span class="w-2 h-2 rounded-full bg-[#72383D]"></span>
                    </div>
                    <label for="context-provider" class="block text-sm font-bold text-[#322D29]">
                      Context Provider
                    </label>
                    <p class="text-xs text-[#322D29]/65 mt-0.5">
                      Generates initial scripture context and reflections.
                    </p>
                  </div>
                  <select
                    id="context-provider"
                    v-model="settings.context_provider"
                    class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-xs font-semibold text-[#322D29] min-h-[40px]"
                  >
                    <option value="ollama">Ollama (Local Model)</option>
                    <option value="gemini">Google Gemini (Cloud AI)</option>
                  </select>
                </div>

                <!-- 2. Fact Checker Provider -->
                <div
                  class="p-4 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 flex flex-col justify-between gap-3 transition-opacity"
                  :class="{ 'opacity-60': !settings.enable_fact_checker }"
                >
                  <div>
                    <div class="flex items-center justify-between mb-1">
                      <span class="text-[10px] font-bold uppercase tracking-[0.12em] text-[#72383D]">Stage 2</span>
                      <span
                        class="text-[10px] font-semibold px-2 py-0.5 rounded-full"
                        :class="settings.enable_fact_checker ? 'bg-emerald-100 text-emerald-800' : 'bg-[#D1C7BD]/40 text-[#322D29]/60'"
                      >
                        {{ settings.enable_fact_checker ? 'Active' : 'Bypassed' }}
                      </span>
                    </div>
                    <label for="fact-checker-provider" class="block text-sm font-bold text-[#322D29]">
                      Fact Checker Provider
                    </label>
                    <p class="text-xs text-[#322D29]/65 mt-0.5">
                      Verifies citations and alignment with Church teaching.
                    </p>
                  </div>
                  <select
                    id="fact-checker-provider"
                    v-model="settings.fact_checker_provider"
                    :disabled="!settings.enable_fact_checker"
                    class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-xs font-semibold text-[#322D29] min-h-[40px] disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <option value="gemini">Google Gemini (Cloud AI)</option>
                    <option value="ollama">Ollama (Local Model)</option>
                    <option value="magisterium">Magisterium AI (Doctrine Search & Verify)</option>
                  </select>
                </div>

                <!-- 3. Magisterium Summary LLM -->
                <div class="p-4 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 flex flex-col justify-between gap-3">
                  <div>
                    <div class="flex items-center justify-between mb-1">
                      <span class="text-[10px] font-bold uppercase tracking-[0.12em] text-[#72383D]">Stage 3</span>
                      <span class="w-2 h-2 rounded-full bg-[#AC9C8D]"></span>
                    </div>
                    <label for="magisterium-llm-provider" class="block text-sm font-bold text-[#322D29]">
                      Magisterium Summary LLM
                    </label>
                    <p class="text-xs text-[#322D29]/65 mt-0.5">
                      Summarizes Magisterium citations for responses.
                    </p>
                  </div>
                  <select
                    id="magisterium-llm-provider"
                    v-model="settings.magisterium_llm_provider"
                    class="w-full px-3 py-2 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-xs font-semibold text-[#322D29] min-h-[40px]"
                  >
                    <option value="ollama">Ollama (Local Model)</option>
                    <option value="gemini">Google Gemini (Cloud AI)</option>
                  </select>
                </div>
              </div>

              <!-- Fact Checker Toggle -->
              <div class="flex items-center justify-between p-4 sm:p-5 bg-[#FAF7EE]/85 border border-[#D1C7BD]/80 rounded-2xl">
                <div class="pr-4">
                  <label for="fact-checker-toggle" class="font-serif font-bold text-[#322D29] text-sm sm:text-base cursor-pointer">
                    Enable Scripture Fact Checker
                  </label>
                  <p class="text-xs text-[#322D29]/65 mt-0.5">
                    When enabled, the fact checker verifies and corrects generated context against Church doctrine before delivery.
                  </p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer flex-shrink-0">
                  <input
                    id="fact-checker-toggle"
                    type="checkbox"
                    v-model="settings.enable_fact_checker"
                    class="sr-only peer"
                  />
                  <div class="w-12 h-7 bg-[#D1C7BD] peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-[#72383D]/30 rounded-full peer peer-checked:after:translate-x-5 peer-checked:after:border-white after:content-[''] after:absolute after:top-[3px] after:left-[3px] after:bg-white after:border-[#D1C7BD] after:border after:rounded-full after:h-[22px] after:w-[22px] after:transition-all peer-checked:bg-[#72383D]"></div>
                </label>
              </div>

              <!-- Magnifica Humanitas Ethical AI Banner -->
              <div class="bg-[#FAF7EE]/70 border border-[#D1C7BD]/70 rounded-2xl p-4 sm:p-5 text-[#322D29] flex items-start gap-3.5 text-xs">
                <div class="w-8 h-8 rounded-xl bg-[#72383D]/10 text-[#72383D] flex items-center justify-center flex-shrink-0 mt-0.5">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                  </svg>
                </div>
                <div class="space-y-1.5">
                  <div class="font-bold text-[#72383D] uppercase tracking-[0.1em] text-[11px]">
                    Ethical AI & Data Privacy (Magnifica Humanitas)
                  </div>
                  <p class="text-[#322D29]/80 font-serif text-xs sm:text-sm leading-relaxed">
                    In accordance with Pope Leo XIV's encyclical <em>Magnifica Humanitas</em>, AI tools must serve human dignity, preserve truth in Sacred Scripture &amp; Magisterium, and remain transparent.
                  </p>
                  <div class="flex flex-wrap gap-x-4 gap-y-1 text-[11px] text-[#72383D] font-semibold pt-0.5">
                    <span>• Local Ollama mode keeps data on-device</span>
                    <span>• Generated reflections are explicitly tagged in UI</span>
                  </div>
                </div>
              </div>
            </section>

            <!-- TAB 2: MODELS -->
            <section
              v-else-if="activeTab === 'models'"
              key="models"
              id="panel-models"
              role="tabpanel"
              class="space-y-5"
            >
              <div>
                <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                  Model Selection
                </h2>
                <p class="text-xs sm:text-sm text-[#322D29]/65 mt-0.5">
                  Choose the cloud Gemini and local Ollama models used across text, voice, and scripture context.
                </p>
              </div>

              <div class="space-y-4">
                <!-- Gemini Text Model -->
                <div class="p-4 sm:p-5 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 space-y-2.5">
                  <div class="flex flex-wrap items-center justify-between gap-2">
                    <label for="gemini-model-select" class="font-serif font-bold text-[#322D29] text-sm sm:text-base">
                      Gemini Text & Context Model
                    </label>
                    <span
                      v-if="settings.context_provider === 'gemini' || settings.magisterium_llm_provider === 'gemini'"
                      class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-[#72383D]/10 text-[#72383D]"
                    >
                      Active in Pipeline
                    </span>
                  </div>
                  <p class="text-xs text-[#322D29]/65">
                    Used for the Catholic Assistant agent (<code class="font-mono bg-[#EFE9E1] px-1.5 py-0.5 rounded text-[#72383D]">main.go</code>) and Gemini context generation.
                  </p>
                  <select
                    id="gemini-model-select"
                    v-model="settings.gemini_model"
                    class="w-full px-3.5 py-2.5 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-sm font-medium min-h-[44px]"
                  >
                    <option v-for="m in geminiModelsList" :key="m" :value="m">
                      {{ m }}
                    </option>
                  </select>
                </div>

                <!-- Gemini Live Model -->
                <div class="p-4 sm:p-5 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 space-y-2.5">
                  <div class="flex flex-wrap items-center justify-between gap-2">
                    <label for="gemini-live-model-select" class="font-serif font-bold text-[#322D29] text-sm sm:text-base">
                      Gemini Live Model (Voice & Real-Time RAG)
                    </label>
                    <span class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-[#EFE9E1] text-[#322D29]/70 border border-[#D1C7BD]">
                      Voice Sanctuary
                    </span>
                  </div>
                  <p class="text-xs text-[#322D29]/65">
                    Dedicated exclusively to real-time voice sessions in Magisterium AI Sanctuary. Defaults to <code class="font-mono bg-[#EFE9E1] px-1.5 py-0.5 rounded text-[#72383D]">gemini-3.8-live</code>.
                  </p>
                  <select
                    id="gemini-live-model-select"
                    v-model="settings.gemini_live_model"
                    class="w-full px-3.5 py-2.5 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-sm font-medium min-h-[44px]"
                  >
                    <option v-for="m in geminiLiveModelsList" :key="m" :value="m">
                      {{ m }}
                    </option>
                  </select>
                </div>

                <!-- Ollama Local Model -->
                <div class="p-4 sm:p-5 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 space-y-2.5">
                  <div class="flex flex-wrap items-center justify-between gap-2">
                    <label for="ollama-model-select" class="font-serif font-bold text-[#322D29] text-sm sm:text-base">
                      Ollama Local Model
                    </label>
                    <div class="flex items-center gap-2">
                      <span
                        v-if="settings.context_provider === 'ollama' || settings.magisterium_llm_provider === 'ollama'"
                        class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-emerald-100 text-emerald-800"
                      >
                        Active Local
                      </span>
                      <button
                        type="button"
                        @click="fetchOllamaModels"
                        :disabled="loadingOllamaModels"
                        class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold border border-[#D1C7BD] bg-white hover:border-[#72383D] hover:text-[#72383D] transition-colors cursor-pointer disabled:opacity-50"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="w-3.5 h-3.5"
                          :class="{ 'animate-spin': loadingOllamaModels }"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                        {{ loadingOllamaModels ? 'Refreshing...' : 'Refresh' }}
                      </button>
                    </div>
                  </div>
                  <p class="text-xs text-[#322D29]/65">
                    Local Ollama model used for private on-device Bible context generation and summaries.
                  </p>
                  <select
                    id="ollama-model-select"
                    v-model="settings.ollama_model"
                    :disabled="loadingOllamaModels"
                    class="w-full px-3.5 py-2.5 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] text-sm font-medium min-h-[44px] disabled:opacity-60"
                  >
                    <option v-for="m in ollamaModelsList" :key="m" :value="m">
                      {{ m }}
                    </option>
                  </select>
                </div>
              </div>
            </section>

            <!-- TAB 3: API KEYS -->
            <section
              v-else-if="activeTab === 'keys'"
              key="keys"
              id="panel-keys"
              role="tabpanel"
              class="space-y-5"
            >
              <div>
                <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                  API Keys & Credentials
                </h2>
                <p class="text-xs sm:text-sm text-[#322D29]/65 mt-0.5">
                  Stored in PostgreSQL. Existing keys are masked on load—leave masked dots untouched to preserve the current key.
                </p>
              </div>

              <div class="space-y-4">
                <!-- Gemini API Key -->
                <div class="p-4 sm:p-5 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 space-y-2.5">
                  <div class="flex items-center justify-between gap-2">
                    <label for="gemini-api-key" class="font-serif font-bold text-[#322D29] text-sm sm:text-base">
                      Gemini API Key
                    </label>
                    <span
                      class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full"
                      :class="settings.gemini_api_key?.trim() ? 'bg-emerald-100 text-emerald-800' : 'bg-[#EFE9E1] text-[#322D29]/60 border border-[#D1C7BD]'"
                    >
                      {{ settings.gemini_api_key?.trim() ? 'Configured' : 'Not Set' }}
                    </span>
                  </div>
                  <p class="text-xs text-[#322D29]/65">
                    Required for Gemini Agent, Live Voice Sanctuary, and Gemini Context generation.
                  </p>
                  <div class="relative flex items-center">
                    <input
                      id="gemini-api-key"
                      v-model="settings.gemini_api_key"
                      :type="showKeys.gemini ? 'text' : 'password'"
                      class="w-full px-3.5 py-2.5 pr-20 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-sm min-h-[44px]"
                      placeholder="AIzaSy..."
                    />
                    <button
                      type="button"
                      @click="showKeys.gemini = !showKeys.gemini"
                      class="absolute right-2 px-2.5 py-1 text-xs font-semibold text-[#72383D] hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer"
                    >
                      {{ showKeys.gemini ? 'Hide' : 'Show' }}
                    </button>
                  </div>
                </div>

                <!-- Magisterium AI API Key -->
                <div class="p-4 sm:p-5 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 space-y-2.5">
                  <div class="flex items-center justify-between gap-2">
                    <label for="magisterium-api-key" class="font-serif font-bold text-[#322D29] text-sm sm:text-base">
                      Magisterium AI API Key
                    </label>
                    <span
                      class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full"
                      :class="settings.magisterium_api_key?.trim() ? 'bg-emerald-100 text-emerald-800' : 'bg-[#EFE9E1] text-[#322D29]/60 border border-[#D1C7BD]'"
                    >
                      {{ settings.magisterium_api_key?.trim() ? 'Configured' : 'Not Set' }}
                    </span>
                  </div>
                  <p class="text-xs text-[#322D29]/65">
                    Required for Catholic doctrine search and chat queries via the Magisterium AI API.
                  </p>
                  <div class="relative flex items-center">
                    <input
                      id="magisterium-api-key"
                      v-model="settings.magisterium_api_key"
                      :type="showKeys.magisterium ? 'text' : 'password'"
                      class="w-full px-3.5 py-2.5 pr-20 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-sm min-h-[44px]"
                      placeholder="magisterium-key-..."
                    />
                    <button
                      type="button"
                      @click="showKeys.magisterium = !showKeys.magisterium"
                      class="absolute right-2 px-2.5 py-1 text-xs font-semibold text-[#72383D] hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer"
                    >
                      {{ showKeys.magisterium ? 'Hide' : 'Show' }}
                    </button>
                  </div>
                </div>

                <!-- API.Bible Key -->
                <div class="p-4 sm:p-5 rounded-2xl border border-[#D1C7BD]/80 bg-white/60 space-y-2.5">
                  <div class="flex items-center justify-between gap-2">
                    <label for="bible-api-key" class="font-serif font-bold text-[#322D29] text-sm sm:text-base">
                      API.Bible Key
                    </label>
                    <span
                      class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full"
                      :class="settings.bible_api_key?.trim() ? 'bg-emerald-100 text-emerald-800' : 'bg-[#EFE9E1] text-[#322D29]/60 border border-[#D1C7BD]'"
                    >
                      {{ settings.bible_api_key?.trim() ? 'Configured' : 'Not Set' }}
                    </span>
                  </div>
                  <p class="text-xs text-[#322D29]/65">
                    Required for word-for-word scripture passage fetching.
                  </p>
                  <div class="relative flex items-center">
                    <input
                      id="bible-api-key"
                      v-model="settings.bible_api_key"
                      :type="showKeys.bible ? 'text' : 'password'"
                      class="w-full px-3.5 py-2.5 pr-20 border border-[#D1C7BD] rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-sm min-h-[44px]"
                      placeholder="api-key-here..."
                    />
                    <button
                      type="button"
                      @click="showKeys.bible = !showKeys.bible"
                      class="absolute right-2 px-2.5 py-1 text-xs font-semibold text-[#72383D] hover:bg-[#EFE9E1] rounded-lg transition-colors cursor-pointer"
                    >
                      {{ showKeys.bible ? 'Hide' : 'Show' }}
                    </button>
                  </div>
                </div>
              </div>
            </section>

            <!-- TAB 4: PROMPTS -->
            <section
              v-else-if="activeTab === 'prompts'"
              key="prompts"
              id="panel-prompts"
              role="tabpanel"
              class="space-y-4"
            >
              <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-2">
                <div>
                  <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                    Bible Context Instruction Prompt
                  </h2>
                  <p class="text-xs sm:text-sm text-[#322D29]/65 mt-0.5">
                    Custom system instructions given to the AI when generating scripture context.
                  </p>
                </div>
                <span class="text-[11px] font-mono text-[#322D29]/55">
                  {{ promptStats.lines }} lines · {{ promptStats.chars }} chars
                </span>
              </div>

              <!-- Clickable Placeholder Chips -->
              <div class="flex flex-wrap items-center gap-2 text-xs">
                <span class="text-[#322D29]/65 font-medium">Insert placeholder:</span>
                <button
                  type="button"
                  @click="insertPlaceholder('{{citation}}')"
                  class="px-2.5 py-1 rounded-full bg-[#FAF7EE] border border-[#D1C7BD] hover:border-[#72383D] text-[#72383D] font-mono text-xs transition-colors cursor-pointer"
                >
                  <span v-pre>+ {{citation}}</span>
                </button>
                <button
                  type="button"
                  @click="insertPlaceholder('{{passage_text}}')"
                  class="px-2.5 py-1 rounded-full bg-[#FAF7EE] border border-[#D1C7BD] hover:border-[#72383D] text-[#72383D] font-mono text-xs transition-colors cursor-pointer"
                >
                  <span v-pre>+ {{passage_text}}</span>
                </button>
              </div>

              <textarea
                ref="promptTextareaRef"
                v-model="settings.context_instruction"
                rows="13"
                placeholder="Enter custom instructions for Bible context generation..."
                class="w-full px-4 py-3 border border-[#D1C7BD] rounded-2xl bg-white/90 focus:outline-none focus:ring-2 focus:ring-[#72383D]/25 focus:border-[#72383D] text-[#322D29] font-mono text-xs sm:text-sm leading-relaxed"
              ></textarea>
            </section>

            <!-- TAB 5: ACCESS CODES & GUEST REQUESTS -->
            <section
              v-else-if="activeTab === 'access'"
              key="access"
              id="panel-access"
              role="tabpanel"
              class="space-y-6"
            >
              <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-2">
                <div>
                  <h2 class="font-serif text-lg sm:text-xl font-bold text-[#322D29]">
                    Guest Access Codes &amp; Invites
                  </h2>
                  <p class="text-xs sm:text-sm text-[#322D29]/65 mt-0.5">
                    Generate codes, copy direct invite links, and approve or revoke guest access for Bible Study &amp; Magisterium AI.
                  </p>
                </div>
                <button
                  type="button"
                  @click="loadAccessData"
                  :disabled="loadingAccess"
                  class="px-3 py-1.5 rounded-xl bg-[#EFE9E1] hover:bg-[#E2D8CC] text-xs font-semibold text-[#5C282D] transition-colors cursor-pointer self-start sm:self-auto"
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

              <!-- Generate New Access Code -->
              <div class="p-4 sm:p-5 rounded-2xl bg-[#FAF7EE] border border-[#D1C7BD]/80 space-y-4">
                <h3 class="font-serif text-sm font-bold text-[#4A2525]">Generate New Access Code</h3>
                <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label class="block text-[11px] font-bold uppercase tracking-[0.1em] text-[#72383D] mb-1">
                      Label / Guest Name
                    </label>
                    <input
                      v-model="newCodeLabel"
                      type="text"
                      placeholder="e.g. Parish Group"
                      class="w-full px-3 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/25"
                    />
                  </div>
                  <div>
                    <label class="block text-[11px] font-bold uppercase tracking-[0.1em] text-[#72383D] mb-1">
                      Custom Code <span class="font-normal text-[#7B6E65] normal-case">(optional)</span>
                    </label>
                    <input
                      v-model="newCustomCode"
                      type="text"
                      placeholder="Auto-generated if blank"
                      class="w-full px-3 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm font-mono uppercase text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/25"
                    />
                  </div>
                  <div>
                    <label class="block text-[11px] font-bold uppercase tracking-[0.1em] text-[#72383D] mb-1">
                      Daily Request Limit
                    </label>
                    <input
                      v-model.number="newDailyLimit"
                      type="number"
                      min="1"
                      max="500"
                      class="w-full px-3 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/25"
                    />
                  </div>
                </div>
                <div class="flex justify-end">
                  <AppButton
                    variant="primary"
                    customClass="!py-2 !px-4 text-xs"
                    :disabled="creatingCode"
                    @click="generateCode"
                  >
                    {{ creatingCode ? 'Generating...' : 'Generate Access Code' }}
                  </AppButton>
                </div>
              </div>

              <!-- Guest Access Requests -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <h3 class="font-serif text-base font-bold text-[#322D29]">
                    Guest Access Requests
                  </h3>
                  <span class="text-xs text-[#7B6E65]">{{ accessRequests.length }} total</span>
                </div>

                <div v-if="accessRequests.length === 0" class="p-4 rounded-xl bg-white/60 border border-[#D1C7BD]/60 text-xs text-[#7B6E65] text-center">
                  No guest access requests yet.
                </div>

                <div v-else class="space-y-2.5">
                  <div
                    v-for="req in accessRequests"
                    :key="req.id"
                    class="p-4 rounded-2xl bg-white/85 border border-[#D1C7BD]/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
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
                      <div v-if="req.generated_code" class="flex flex-wrap items-center gap-2 pt-1">
                        <span class="font-mono text-xs font-bold px-2 py-0.5 rounded bg-[#FAF7EE] border border-[#D1C7BD] text-[#72383D]">
                          {{ req.generated_code }}
                        </span>
                        <button
                          type="button"
                          @click="copyText(getInviteLink(req.generated_code), `req-link-${req.id}`)"
                          class="text-[11px] font-semibold text-[#72383D] hover:underline cursor-pointer"
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
                        class="px-3 py-1.5 rounded-xl bg-[#72383D] hover:bg-[#5C282D] text-[#FAF7F2] text-xs font-semibold transition-colors cursor-pointer"
                      >
                        Approve &amp; Generate
                      </button>
                      <button
                        v-if="req.status === 'pending'"
                        type="button"
                        @click="rejectRequest(req)"
                        class="px-3 py-1.5 rounded-xl bg-[#EFE9E1] hover:bg-[#E2D8CC] text-[#5C282D] text-xs font-semibold transition-colors cursor-pointer"
                      >
                        Reject
                      </button>
                      <button
                        type="button"
                        @click="deleteRequest(req.id)"
                        class="px-2.5 py-1.5 rounded-xl text-xs text-red-700 hover:bg-red-50 transition-colors cursor-pointer"
                        title="Delete request"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Active & Revoked Access Codes -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <h3 class="font-serif text-base font-bold text-[#322D29]">
                    Access Codes
                  </h3>
                  <span class="text-xs text-[#7B6E65]">{{ accessCodes.length }} total</span>
                </div>

                <div v-if="accessCodes.length === 0" class="p-4 rounded-xl bg-white/60 border border-[#D1C7BD]/60 text-xs text-[#7B6E65] text-center">
                  No access codes created yet. Generate one above.
                </div>

                <div v-else class="space-y-2.5">
                  <div
                    v-for="item in accessCodes"
                    :key="item.id"
                    class="p-4 rounded-2xl bg-white/85 border border-[#D1C7BD]/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
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
                          Today: {{ item.used_today }} / {{ item.daily_limit }}
                        </span>
                      </div>
                      <p v-if="item.label" class="text-xs text-[#6B5F56]">{{ item.label }}</p>
                    </div>

                    <div class="flex flex-wrap items-center gap-2 shrink-0">
                      <button
                        type="button"
                        @click="copyText(item.code, `code-${item.id}`)"
                        class="px-2.5 py-1.5 rounded-xl bg-[#FAF7EE] hover:bg-[#EFE9E1] border border-[#D1C7BD] text-xs font-semibold text-[#5C282D] transition-colors cursor-pointer"
                      >
                        {{ copiedKey === `code-${item.id}` ? 'Copied!' : 'Copy Code' }}
                      </button>
                      <button
                        type="button"
                        @click="copyText(getInviteLink(item.code), `link-${item.id}`)"
                        class="px-2.5 py-1.5 rounded-xl bg-[#72383D]/10 hover:bg-[#72383D]/20 border border-[#72383D]/25 text-xs font-semibold text-[#72383D] transition-colors cursor-pointer"
                      >
                        {{ copiedKey === `link-${item.id}` ? 'Link Copied!' : 'Copy Invite Link' }}
                      </button>
                      <button
                        type="button"
                        @click="toggleCodeActive(item)"
                        class="px-2.5 py-1.5 rounded-xl bg-[#EFE9E1] hover:bg-[#E2D8CC] text-xs font-semibold text-[#322D29] transition-colors cursor-pointer"
                      >
                        {{ item.is_active ? 'Revoke' : 'Activate' }}
                      </button>
                      <button
                        type="button"
                        @click="deleteCode(item.id)"
                        class="px-2.5 py-1.5 rounded-xl text-xs text-red-700 hover:bg-red-50 transition-colors cursor-pointer"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </section>
          </Transition>

          <!-- Persistent Action Footer inside Card -->
          <div class="mt-7 pt-5 border-t border-[#D1C7BD]/70 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
            <div class="flex items-center gap-2 text-xs">
              <span
                v-if="loadingSettings"
                class="text-[#322D29]/60"
              >
                Loading settings...
              </span>
              <template v-else-if="hasUnsavedChanges">
                <span class="w-2 h-2 rounded-full bg-amber-500 animate-pulse"></span>
                <span class="font-semibold text-[#72383D]">You have unsaved changes</span>
              </template>
              <template v-else>
                <span class="w-2 h-2 rounded-full bg-emerald-600/70"></span>
                <span class="text-[#322D29]/60">All changes synced with database</span>
              </template>
            </div>

            <div class="flex items-center justify-end gap-2.5">
              <AppButton
                v-if="hasUnsavedChanges"
                variant="outlined"
                customClass="!px-4 !py-2 text-xs min-h-[42px]"
                @click="discardChanges"
              >
                Discard
              </AppButton>
              <AppButton
                variant="primary"
                customClass="min-h-[42px]"
                @click="saveSettings"
                :disabled="saveLoading"
              >
                {{ saveLoading ? 'Saving...' : 'Save Settings' }}
              </AppButton>
            </div>
          </div>
        </ParchmentCard>
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
</style>
