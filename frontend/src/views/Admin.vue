<script setup lang="ts">
import { ref, onMounted } from 'vue';
import TopNav from '../components/common/TopNav.vue';
import BottomNav from '../components/common/BottomNav.vue';
import ParchmentCard from '../components/common/ParchmentCard.vue';
import AppButton from '../components/common/AppButton.vue';
import { useDevotionApi } from '../composables/useDevotionApi';

interface Settings {
  context_provider: string;
  fact_checker_provider: string;
  magisterium_llm_provider: string;
  gemini_api_key: string;
  magisterium_api_key: string;
  bible_api_key: string;
  gemini_model: string;
  ollama_model: string;
  context_model?: string;
  context_instruction: string;
  enable_fact_checker: boolean;
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
  fact_checker_provider: 'gemini',
  magisterium_llm_provider: 'ollama',
  gemini_api_key: '',
  magisterium_api_key: '',
  bible_api_key: '',
  gemini_model: 'gemini-3.1-flash-lite',
  ollama_model: 'gemma4:cloud',
  context_instruction: '',
  enable_fact_checker: true,
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

        <!-- Magnifica Humanitas Ethical AI Banner -->
        <div class="bg-amber-100/60 border border-amber-300/70 rounded-xl p-4 text-amber-950 flex items-start space-x-3 text-xs shadow-xs">
          <span class="text-base flex-shrink-0">📜</span>
          <div class="space-y-1">
            <div class="font-bold text-amber-900 flex items-center gap-1.5 uppercase tracking-wider text-[11px]">
              Ethical AI & Data Privacy Principles (Magnifica Humanitas)
            </div>
            <p class="text-amber-900/80 font-serif leading-relaxed">
              In accordance with Pope Leo XIV's encyclical <em>Magnifica Humanitas</em>, AI tools must serve human dignity, preserve truth in Sacred Scripture & Magisterium, and remain transparent.
            </p>
            <div class="flex flex-wrap gap-x-4 gap-y-1 text-[10px] text-amber-800 font-semibold pt-1">
              <span>🛡️ Local Ollama mode preserves data privacy on-device.</span>
              <span>🏷️ All generated reflections are explicitly tagged in the UI.</span>
            </div>
          </div>
        </div>

        <div v-if="message" class="p-4 rounded border text-sm font-medium" :class="message.type === 'success' ? 'bg-emerald-50 border-emerald-300 text-emerald-800' : 'bg-red-50 border-red-300 text-red-800'">
          {{ message.text }}
        </div>

        <ParchmentCard class="p-6 space-y-6 border border-amber-200 shadow-sm">
          <!-- Pipeline Architecture Information & Provider Selection -->
          <div>
            <label class="block font-serif font-bold text-amber-950 text-base mb-2">
              AI Generation & Fact-Checking Pipeline
            </label>
            <p class="text-xs text-stone-600 mb-3">System architecture and AI model providers for scripture context generation & verification:</p>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
              <!-- Context Provider Picker -->
              <div class="p-3 border border-amber-300 bg-amber-100/50 rounded text-amber-950 space-y-2">
                <div>
                  <span class="block text-sm font-semibold">1. Context Provider</span>
                  <span class="block text-xs text-stone-600">Select AI provider to generate scripture context.</span>
                </div>
                <select
                  v-model="settings.context_provider"
                  class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white focus:outline-none focus:ring-2 focus:ring-amber-500 text-xs font-semibold text-stone-800"
                >
                  <option value="ollama">Ollama (Local Model)</option>
                  <option value="gemini">Google Gemini (Cloud AI)</option>
                </select>
              </div>

              <!-- Fact Checker Provider Picker -->
              <div class="p-3 border border-amber-300 bg-amber-100/50 rounded text-amber-950 space-y-2">
                <div>
                  <span class="block text-sm font-semibold">2. Fact Checker Provider</span>
                  <span class="block text-xs text-stone-600">Select AI provider to verify citations.</span>
                </div>
                <select
                  v-model="settings.fact_checker_provider"
                  :disabled="!settings.enable_fact_checker"
                  class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white focus:outline-none focus:ring-2 focus:ring-amber-500 text-xs font-semibold text-stone-800 disabled:opacity-50"
                >
                  <option value="gemini">Google Gemini (Cloud AI)</option>
                  <option value="ollama">Ollama (Local Model)</option>
                  <option value="magisterium">Magisterium AI (Doctrine Search & Verify)</option>
                </select>
              </div>

              <!-- Magisterium LLM Summary Provider Picker -->
              <div class="p-3 border border-amber-300 bg-amber-100/50 rounded text-amber-950 space-y-2">
                <div>
                  <span class="block text-sm font-semibold">3. Magisterium Summary LLM</span>
                  <span class="block text-xs text-stone-600">Select LLM to summarize citations.</span>
                </div>
                <select
                  v-model="settings.magisterium_llm_provider"
                  class="w-full px-2.5 py-1.5 border border-amber-300 rounded bg-white focus:outline-none focus:ring-2 focus:ring-amber-500 text-xs font-semibold text-stone-800"
                >
                  <option value="ollama">Ollama (Local Model)</option>
                  <option value="gemini">Google Gemini (Cloud AI)</option>
                </select>
              </div>
            </div>

            <!-- Fact Checker Toggle Switch -->
            <div class="flex items-center justify-between p-4 bg-amber-100/40 border border-amber-300 rounded-lg">
              <div>
                <label for="fact-checker-toggle" class="font-serif font-semibold text-amber-950 text-sm cursor-pointer">
                  Enable Scripture Fact Checker
                </label>
                <p class="text-xs text-stone-600 mt-0.5">
                  When enabled, the fact checker verifies and corrects generated context against Church doctrine.
                </p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer ml-4">
                <input
                  id="fact-checker-toggle"
                  type="checkbox"
                  v-model="settings.enable_fact_checker"
                  class="sr-only peer"
                />
                <div class="w-11 h-6 bg-stone-300 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-amber-500 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-stone-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-amber-600"></div>
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

          <!-- Magisterium AI Key -->
          <div class="pt-2">
            <label class="block font-serif font-semibold text-amber-950 text-sm mb-1">
              Magisterium AI API Key
            </label>
            <p class="text-xs text-stone-600 mb-2">Stored securely in PostgreSQL database. Required for Catholic doctrine chat queries via Magisterium AI API.</p>
            <input
              v-model="settings.magisterium_api_key"
              type="password"
              class="w-full px-3 py-2 border border-amber-300 rounded bg-white/80 focus:outline-none focus:ring-2 focus:ring-amber-500 text-stone-800 font-mono text-sm"
              placeholder="magisterium-key-..."
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
          <div class="pt-2">
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
