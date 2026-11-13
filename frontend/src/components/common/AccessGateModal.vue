<template>
  <Teleport to="body">
    <transition name="fade">
      <div
        v-if="isAccessModalOpen"
        class="fixed inset-0 z-[80] flex items-center justify-center p-4 bg-[#201A17]/65 backdrop-blur-sm"
        role="dialog"
        aria-modal="true"
        aria-labelledby="access-gate-title"
        @click.self="closeAccessModal"
      >
        <div
          class="w-full max-w-md rounded-2xl bg-[#FAF7F2] border border-[#D1C7BD] shadow-2xl overflow-hidden text-[#322D29]"
        >
          <!-- Header -->
          <div class="px-6 pt-6 pb-4 bg-gradient-to-b from-[#F3ECE3] to-[#FAF7F2] border-b border-[#E5DDD3] flex items-start justify-between gap-4">
            <div>
              <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-[0.14em] bg-[#72383D]/10 text-[#72383D] border border-[#72383D]/20 mb-2">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                Guest Access Required
              </span>
              <h2 id="access-gate-title" class="font-serif text-xl font-bold text-[#322D29]">
                {{ mode === 'code' ? 'Enter Access Code' : 'Request Guest Access' }}
              </h2>
              <p class="text-xs text-[#6B5F56] mt-1 leading-relaxed">
                {{
                  mode === 'code'
                    ? 'Bible Study and Magisterium AI Chat require an access code (up to 10 requests per day).'
                    : 'Submit a request to the administrator to receive a guest access code.'
                }}
              </p>
            </div>
            <button
              type="button"
              @click="closeAccessModal"
              aria-label="Close modal"
              class="p-1.5 rounded-lg text-[#7B6E65] hover:text-[#322D29] hover:bg-[#EAE2D8]/60 transition-colors"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Body: Enter Code -->
          <div v-if="mode === 'code'" class="p-6 space-y-4">
            <form @submit.prevent="handleUnlock" class="space-y-4">
              <div>
                <label for="guest-access-code" class="block text-[11px] font-bold uppercase tracking-[0.12em] text-[#72383D] mb-1.5">
                  Access Code
                </label>
                <input
                  id="guest-access-code"
                  ref="codeInputRef"
                  v-model="codeInput"
                  type="text"
                  autocomplete="off"
                  placeholder="e.g. DEV-XXXX-XXXX"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-white border border-[#D1C7BD] text-[#322D29] font-mono text-sm tracking-wider uppercase placeholder:normal-case placeholder:tracking-normal placeholder:text-[#A0948B] focus:outline-none focus:ring-2 focus:ring-[#72383D]/30 focus:border-[#72383D]"
                />
              </div>

              <div
                v-if="accessError"
                class="p-3 rounded-xl bg-[#9B2C2C]/10 border border-[#9B2C2C]/25 text-xs text-[#822020] leading-relaxed"
              >
                {{ accessError }}
              </div>

              <div class="flex flex-col sm:flex-row gap-2.5 pt-1">
                <button
                  type="submit"
                  :disabled="verifying || !codeInput.trim()"
                  class="flex-1 px-4 py-2.5 rounded-xl bg-[#72383D] hover:bg-[#5C282D] disabled:opacity-50 disabled:cursor-not-allowed text-[#FAF7F2] text-sm font-semibold transition-colors shadow-sm cursor-pointer"
                >
                  {{ verifying ? 'Verifying...' : 'Unlock Access' }}
                </button>
                <button
                  type="button"
                  @click="switchToRequest"
                  class="px-4 py-2.5 rounded-xl bg-[#EAE2D8] hover:bg-[#DFD4C6] text-[#4A2525] text-sm font-semibold transition-colors cursor-pointer"
                >
                  Request Access
                </button>
              </div>
            </form>

            <div class="pt-3 border-t border-[#E5DDD3] flex items-center justify-between text-xs text-[#7B6E65]">
              <span>Your code is saved locally for future visits.</span>
              <router-link
                to="/admin"
                @click="closeAccessModal"
                class="text-[#72383D] font-semibold hover:underline whitespace-nowrap ml-2"
              >
                Admin Login
              </router-link>
            </div>
          </div>

          <!-- Body: Request Access -->
          <div v-else class="p-6 space-y-4">
            <div v-if="requestSubmitted" class="space-y-4 text-center py-2">
              <div class="w-12 h-12 rounded-full bg-emerald-100 text-emerald-700 mx-auto flex items-center justify-center">
                <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <div>
                <h3 class="font-serif text-lg font-bold text-[#322D29]">Request Submitted</h3>
                <p class="text-xs text-[#6B5F56] mt-1.5 leading-relaxed">
                  {{ requestSuccessMessage }}
                </p>
              </div>
              <div class="flex gap-2.5 pt-2">
                <button
                  type="button"
                  @click="mode = 'code'"
                  class="flex-1 px-4 py-2.5 rounded-xl bg-[#72383D] hover:bg-[#5C282D] text-[#FAF7F2] text-sm font-semibold transition-colors cursor-pointer"
                >
                  Enter Access Code
                </button>
                <button
                  type="button"
                  @click="closeAccessModal"
                  class="px-4 py-2.5 rounded-xl bg-[#EAE2D8] hover:bg-[#DFD4C6] text-[#4A2525] text-sm font-semibold transition-colors cursor-pointer"
                >
                  Close
                </button>
              </div>
            </div>

            <form v-else @submit.prevent="handleRequestSubmit" class="space-y-3.5">
              <div>
                <label for="req-name" class="block text-[11px] font-bold uppercase tracking-[0.12em] text-[#72383D] mb-1">
                  Your Name <span class="text-[#9B2C2C]">*</span>
                </label>
                <input
                  id="req-name"
                  v-model="reqName"
                  type="text"
                  required
                  placeholder="e.g. Maria Santos"
                  class="w-full px-3.5 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/30 focus:border-[#72383D]"
                />
              </div>

              <div>
                <label for="req-email" class="block text-[11px] font-bold uppercase tracking-[0.12em] text-[#72383D] mb-1">
                  Email <span class="text-[#7B6E65] font-normal normal-case">(optional, to send your code)</span>
                </label>
                <input
                  id="req-email"
                  v-model="reqEmail"
                  type="email"
                  placeholder="you@example.com"
                  class="w-full px-3.5 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/30 focus:border-[#72383D]"
                />
              </div>

              <div>
                <label for="req-reason" class="block text-[11px] font-bold uppercase tracking-[0.12em] text-[#72383D] mb-1">
                  Note / Reason <span class="text-[#7B6E65] font-normal normal-case">(optional)</span>
                </label>
                <textarea
                  id="req-reason"
                  v-model="reqReason"
                  rows="2"
                  placeholder="How do you plan to use Bible Study or Magisterium AI?"
                  class="w-full px-3.5 py-2 rounded-xl bg-white border border-[#D1C7BD] text-sm text-[#322D29] focus:outline-none focus:ring-2 focus:ring-[#72383D]/30 focus:border-[#72383D] resize-none"
                ></textarea>
              </div>

              <div
                v-if="reqError"
                class="p-3 rounded-xl bg-[#9B2C2C]/10 border border-[#9B2C2C]/25 text-xs text-[#822020]"
              >
                {{ reqError }}
              </div>

              <div class="flex flex-col sm:flex-row gap-2.5 pt-1">
                <button
                  type="submit"
                  :disabled="reqLoading || !reqName.trim()"
                  class="flex-1 px-4 py-2.5 rounded-xl bg-[#72383D] hover:bg-[#5C282D] disabled:opacity-50 disabled:cursor-not-allowed text-[#FAF7F2] text-sm font-semibold transition-colors shadow-sm cursor-pointer"
                >
                  {{ reqLoading ? 'Submitting...' : 'Submit Request' }}
                </button>
                <button
                  type="button"
                  @click="mode = 'code'"
                  class="px-4 py-2.5 rounded-xl bg-[#EAE2D8] hover:bg-[#DFD4C6] text-[#4A2525] text-sm font-semibold transition-colors cursor-pointer"
                >
                  Back
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue';
import { useAccessControl } from '../../composables/useAccessControl';

const {
  isAccessModalOpen,
  accessError,
  verifying,
  closeAccessModal,
  unlockWithCode,
  submitAccessRequest,
} = useAccessControl();

const mode = ref<'code' | 'request'>('code');
const codeInput = ref('');
const codeInputRef = ref<HTMLInputElement | null>(null);

const reqName = ref('');
const reqEmail = ref('');
const reqReason = ref('');
const reqLoading = ref(false);
const reqError = ref<string | null>(null);
const requestSubmitted = ref(false);
const requestSuccessMessage = ref('');

watch(isAccessModalOpen, async (open) => {
  if (open) {
    mode.value = 'code';
    requestSubmitted.value = false;
    reqError.value = null;
    await nextTick();
    codeInputRef.value?.focus();
  }
});

const switchToRequest = () => {
  reqError.value = null;
  requestSubmitted.value = false;
  mode.value = 'request';
};

const handleUnlock = async () => {
  if (!codeInput.value.trim()) return;
  const ok = await unlockWithCode(codeInput.value);
  if (ok) {
    codeInput.value = '';
  }
};

const handleRequestSubmit = async () => {
  if (!reqName.value.trim()) return;
  reqLoading.value = true;
  reqError.value = null;
  try {
    const res = await submitAccessRequest({
      name: reqName.value.trim(),
      email: reqEmail.value.trim(),
      reason: reqReason.value.trim(),
    });
    requestSubmitted.value = true;
    requestSuccessMessage.value =
      res.message || 'Your access request has been submitted. Once approved, enter your access code here.';
    reqName.value = '';
    reqEmail.value = '';
    reqReason.value = '';
  } catch (err: any) {
    reqError.value = err?.message || 'Failed to submit access request. Please try again.';
  } finally {
    reqLoading.value = false;
  }
};

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isAccessModalOpen.value) {
    closeAccessModal();
  }
};

onMounted(() => {
  window.addEventListener('keydown', onKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown);
});
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
