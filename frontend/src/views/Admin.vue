<script setup lang="ts">
import { ref, onMounted } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import AppButton from '../components/common/AppButton.vue';
import { useDevotionApi } from '../composables/useDevotionApi';

interface Settings {
  context_provider: string;
  gemini_api_key: string;
  bible_api_key: string;
  gemini_model: string;
  ollama_model: string;
  context_model?: string;
  context_instruction: string;
}

const { fetchWithAuth, isServerUnreachable } = useDevotionApi();

// Auth State
const isAuthenticated = ref(false);
const loginUsername = ref('admin');
const loginPassword = ref('');
const authError = ref<string | null>(null);
const authLoading = ref(false);

// Settings State
const settings = ref<Settings>({
  context_provider: 'ollama',
  gemini_api_key: '',
  bible_api_key: '',
  gemini_model: 'gemini-3.1-flash-lite',
  ollama_model: 'gemma4:cloud',
  context_instruction: '',
});

const geminiModelsList = ref<string[]>([
  'gemini-3.1-flash-lite',
  'gemini-2.5-flash',
  'gemini-2.5-pro',
  'gemini-1.5-flash',
  'gemini-1.5-pro',
]);
const ollamaModelsList = ref<string[]>([]);
const loadingOllamaModels = ref(false);
const saveLoading = ref(false);
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null);

onMounted(() => {
  const token = localStorage.getItem('adminToken');
  if (token) {
    isAuthenticated.value = true;
    loadSettings();
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
      isAuthenticated.value = true;
      loadSettings();
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
  isAuthenticated.value = false;
  loginPassword.value = '';
  message.value = null;
};

const loadSettings = async () => {
  try {
    const data = await fetchWithAuth<Settings>('/api/admin/settings');
    if (data) {
      settings.value = data;
    }
    fetchGeminiModels();
    fetchOllamaModels();
  } catch (err: any) {
    console.error('Failed to load settings:', err);
  }
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
    message.value = { type: 'success', text: 'Settings successfully saved to PostgreSQL database!' };
  } catch (err: any) {
    message.value = { type: 'error', text: err.message || 'Failed to save settings to backend.' };
  } finally {
    saveLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen bg-amber-50/40 text-stone-800 pb-24">
    <TopNav title="Admin Settings" />

    <main class="max-w-4xl mx-auto px-4 py-8">
      <!-- Server Unreachable Banner -->
      <div v-if="isServerUnreachable" class="mb-6 p-4 bg-red-100 border-l-4 border-red-600 rounded-r text-red-900 shadow-sm flex items-center space-x-3">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-red-600 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <h3 class="font-bold text-sm">Server Cannot Be Reached</h3>
          <p class="text-xs text-red-700 mt-0.5">Unable to connect to the backend service. Please check if the API server is running and accessible.</p>
        </div>
      </div>
      <!-- Login View -->
      <div v-if="!isAuthenticated" class="max-w-md mx-auto my-12">
        <ParchmentCard class="p-8 border border-amber-200 shadow-md">
          <div class="text-center mb-6">
            <h1 class="font-serif text-2xl font-bold text-amber-900">Admin Portal Login</h1>
            <p class="text-sm text-stone-600 mt-1">Authentication required to configure system settings</p>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-4">
            <div>
              <label class="block text-xs font-semibold uppercase tracking-wider text-amber-900 mb-1">Username</label>
              <input
                v-model="loginUsername"
                type="text"
                required
                class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800"
                placeholder="admin"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold uppercase tracking-wider text-amber-900 mb-1">Password</label>
              <input
                v-model="loginPassword"
                type="password"
                required
                class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800"
                placeholder="admin"
              />
            </div>

            <div v-if="authError" class="p-3 bg-red-100 border border-red-300 text-red-700 text-sm rounded">
              {{ authError }}
            </div>

            <div class="pt-2">
              <AppButton type="submit" variant="primary" class="w-full justify-center" :disabled="authLoading">
                {{ authLoading ? 'Authenticating...' : 'Login' }}
              </AppButton>
            </div>
          </form>
        </ParchmentCard>
      </div>

      <!-- Settings Dashboard -->
      <div v-else class="space-y-6">
        <div class="flex justify-between items-center pb-4 border-b border-amber-200">
          <div>
            <h1 class="font-serif text-2xl font-bold text-amber-950">AI & System Settings</h1>
            <p class="text-xs text-stone-600">Configure database-backed AI generation and API settings</p>
          </div>
          <AppButton variant="secondary" @click="handleLogout">
            Logout
          </AppButton>
        </div>

        <div v-if="message" class="p-4 rounded border text-sm font-medium" :class="message.type === 'success' ? 'bg-emerald-50 border-emerald-300 text-emerald-800' : 'bg-red-50 border-red-300 text-red-800'">
          {{ message.text }}
        </div>

        <ParchmentCard class="p-6 space-y-6 border border-amber-200 shadow-sm">
          <!-- Provider Selection -->
          <div>
            <label class="block font-serif font-bold text-amber-950 text-base mb-2">
              AI Generation Provider for Bible Context
            </label>
            <p class="text-xs text-stone-600 mb-3">Choose which AI provider generates scripture analysis & context.</p>
            <div class="grid grid-cols-2 gap-4">
              <label
                class="flex items-center p-3 border rounded cursor-pointer transition-colors"
                :class="settings.context_provider === 'ollama' ? 'border-amber-600 bg-amber-100/50 font-semibold text-amber-950' : 'border-amber-200 bg-white/50 text-stone-700'"
              >
                <input
                  type="radio"
                  value="ollama"
                  v-model="settings.context_provider"
                  class="mr-3 text-amber-600 focus:ring-amber-500"
                />
                <div>
                  <span class="block text-sm">Ollama (Local / Self-hosted)</span>
                  <span class="block text-xs text-stone-500">Uses local Ollama model</span>
                </div>
              </label>

              <label
                class="flex items-center p-3 border rounded cursor-pointer transition-colors"
                :class="settings.context_provider === 'gemini' ? 'border-amber-600 bg-amber-100/50 font-semibold text-amber-950' : 'border-amber-200 bg-white/50 text-stone-700'"
              >
                <input
                  type="radio"
                  value="gemini"
                  v-model="settings.context_provider"
                  class="mr-3 text-amber-600 focus:ring-amber-500"
                />
                <div>
                  <span class="block text-sm">Google Gemini</span>
                  <span class="block text-xs text-stone-500">Uses Gemini model & cloud API</span>
                </div>
              </label>
            </div>
          </div>

          <!-- Gemini API Key -->
          <div class="pt-2">
            <label class="block font-serif font-semibold text-amber-950 text-sm mb-1">
              Gemini API Key
            </label>
            <p class="text-xs text-stone-600 mb-2">Stored securely in PostgreSQL database. Required for Gemini Agent & Gemini Context.</p>
            <input
              v-model="settings.gemini_api_key"
              type="password"
              class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800 font-mono text-sm"
              placeholder="AIzaSy..."
            />
          </div>

          <!-- Bible API Key -->
          <div class="pt-2">
            <label class="block font-serif font-semibold text-amber-950 text-sm mb-1">
              API.Bible Key
            </label>
            <p class="text-xs text-stone-600 mb-2">Stored securely in PostgreSQL database. Required for word-for-word scripture fetching.</p>
            <input
              v-model="settings.bible_api_key"
              type="password"
              class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800 font-mono text-sm"
              placeholder="api-key-here..."
            />
          </div>

          <!-- Always Available Gemini Model Selection -->
          <div class="pt-2">
            <label class="block font-serif font-semibold text-amber-950 text-sm mb-1">
              Gemini Model Selection
            </label>
            <p class="text-xs text-stone-600 mb-2">
              Select the Gemini model used for the Catholic Assistant agent (main.go) and Gemini context generation.
            </p>
            <select
              v-model="settings.gemini_model"
              class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800 text-sm font-medium"
            >
              <option v-for="m in geminiModelsList" :key="m" :value="m">
                {{ m }}
              </option>
            </select>
          </div>

          <!-- Ollama Model Selection -->
          <div v-if="settings.context_provider === 'ollama'" class="pt-2">
            <label class="block font-serif font-semibold text-amber-950 text-sm mb-1">
              Ollama Model Selection
            </label>
            <p class="text-xs text-stone-600 mb-2">
              Select the local Ollama model used for Bible context generation.
            </p>
            <div class="relative">
              <select
                v-model="settings.ollama_model"
                :disabled="loadingOllamaModels"
                class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800 text-sm"
              >
                <option v-for="m in ollamaModelsList" :key="m" :value="m">
                  {{ m }}
                </option>
              </select>
              <span v-if="loadingOllamaModels" class="text-xs text-amber-700 mt-1 block">Loading local Ollama models...</span>
            </div>
          </div>

          <!-- Context Instruction Prompt -->
          <div class="pt-2">
            <label class="block font-serif font-semibold text-amber-950 text-sm mb-1">
              Bible Context Instruction Prompt
            </label>
            <p class="text-xs text-stone-600 mb-2">
              Custom prompt instructions given to the AI. Use <code v-pre>{{citation}}</code> and <code v-pre>{{passage_text}}</code> as dynamic placeholders.
            </p>
            <textarea
              v-model="settings.context_instruction"
              rows="10"
              class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800 font-mono text-xs leading-relaxed"
            ></textarea>
          </div>

          <!-- Save Button -->
          <div class="pt-4 border-t border-amber-200 flex justify-end">
            <AppButton variant="primary" @click="saveSettings" :disabled="saveLoading">
              {{ saveLoading ? 'Saving...' : 'Save Settings to Database' }}
            </AppButton>
          </div>
        </ParchmentCard>
      </div>
    </main>

    <BottomNav />
  </div>
</template>
