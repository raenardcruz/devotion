import { ref, computed } from 'vue';
import { useDevotionApi } from './useDevotionApi';

const STORAGE_KEY = 'guestAccessCode';
const ADMIN_TOKEN_KEY = 'adminToken';
const ADMIN_TOKEN_VALUE = 'admin-authenticated-token';

interface VerifyAccessResponse {
  valid: boolean;
  is_admin: boolean;
  code?: string;
  label?: string;
  daily_limit: number;
  used_today: number;
  remaining: number;
}

interface RequestAccessPayload {
  name: string;
  email: string;
  reason: string;
}

const adminTokenRef = ref<string | null>(
  typeof localStorage !== 'undefined' ? localStorage.getItem(ADMIN_TOKEN_KEY) : null
);
const storedAccessCode = ref<string>(
  typeof localStorage !== 'undefined' ? localStorage.getItem(STORAGE_KEY) || '' : ''
);
const isCodeVerified = ref<boolean>(false);
const dailyLimit = ref<number>(10);
const usedToday = ref<number>(0);
const isAccessModalOpen = ref<boolean>(false);
const accessError = ref<string | null>(null);
const verifying = ref<boolean>(false);
let pendingCallback: (() => void) | null = null;
let verifyPromise: Promise<boolean> | null = null;

export function useAccessControl() {
  const { fetchWithAuth } = useDevotionApi();

  const syncAdminState = () => {
    if (typeof localStorage !== 'undefined') {
      adminTokenRef.value = localStorage.getItem(ADMIN_TOKEN_KEY);
    }
  };

  const isAdmin = computed(() => adminTokenRef.value === ADMIN_TOKEN_VALUE);

  const remainingToday = computed(() => {
    if (isAdmin.value) return Infinity;
    return Math.max(0, dailyLimit.value - usedToday.value);
  });

  const hasAccess = computed(() => {
    if (isAdmin.value) return true;
    return isCodeVerified.value && !!storedAccessCode.value;
  });

  const verifyCode = async (rawCode: string, persist = true): Promise<boolean> => {
    const normalized = rawCode.trim().toUpperCase();
    if (!normalized) {
      accessError.value = 'Please enter an access code.';
      return false;
    }

    verifying.value = true;
    accessError.value = null;

    try {
      const res = await fetchWithAuth<VerifyAccessResponse>('/api/access/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code: normalized }),
      });

      if (res && res.valid) {
        if (res.is_admin) {
          isCodeVerified.value = true;
          return true;
        }
        const finalCode = (res.code || normalized).toUpperCase();
        storedAccessCode.value = finalCode;
        isCodeVerified.value = true;
        dailyLimit.value = res.daily_limit ?? 10;
        usedToday.value = res.used_today ?? 0;

        if (persist && typeof localStorage !== 'undefined') {
          localStorage.setItem(STORAGE_KEY, finalCode);
        }
        return true;
      }

      accessError.value = 'Invalid or inactive access code.';
      return false;
    } catch (err: any) {
      accessError.value = err?.message || 'Invalid or inactive access code.';
      if (persist && typeof localStorage !== 'undefined' && normalized === storedAccessCode.value) {
        localStorage.removeItem(STORAGE_KEY);
        storedAccessCode.value = '';
        isCodeVerified.value = false;
      }
      return false;
    } finally {
      verifying.value = false;
    }
  };

  const ensureAccess = async (): Promise<boolean> => {
    syncAdminState();
    if (isAdmin.value) {
      return true;
    }

    if (typeof localStorage !== 'undefined') {
      const fromStorage = localStorage.getItem(STORAGE_KEY) || '';
      if (fromStorage && fromStorage !== storedAccessCode.value) {
        storedAccessCode.value = fromStorage;
      }
    }

    if (!storedAccessCode.value) {
      isCodeVerified.value = false;
      return false;
    }

    if (verifyPromise) {
      return verifyPromise;
    }

    verifyPromise = verifyCode(storedAccessCode.value, true).finally(() => {
      verifyPromise = null;
    });
    return verifyPromise;
  };

  const acceptInviteCode = async (rawCode: string): Promise<boolean> => {
    const normalized = rawCode.trim().toUpperCase();
    if (!normalized) return false;
    // Save immediately in localStorage so subsequent checks/requests have it, then verify against DB
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(STORAGE_KEY, normalized);
    }
    storedAccessCode.value = normalized;
    return await verifyCode(normalized, true);
  };

  const refreshUsage = async (): Promise<void> => {
    if (isAdmin.value || !storedAccessCode.value) return;
    await verifyCode(storedAccessCode.value, false);
  };

  const clearAccessCode = () => {
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem(STORAGE_KEY);
    }
    storedAccessCode.value = '';
    isCodeVerified.value = false;
    usedToday.value = 0;
    accessError.value = null;
  };

  const openAccessModal = (onSuccess?: () => void) => {
    accessError.value = null;
    pendingCallback = onSuccess || null;
    isAccessModalOpen.value = true;
  };

  const closeAccessModal = () => {
    isAccessModalOpen.value = false;
    pendingCallback = null;
    accessError.value = null;
  };

  const unlockWithCode = async (code: string): Promise<boolean> => {
    const ok = await verifyCode(code, true);
    if (ok) {
      const cb = pendingCallback;
      pendingCallback = null;
      isAccessModalOpen.value = false;
      if (cb) cb();
    }
    return ok;
  };

  const submitAccessRequest = async (payload: RequestAccessPayload): Promise<{ success: boolean; message: string }> => {
    return await fetchWithAuth<{ success: boolean; message: string }>('/api/access/request', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
  };

  return {
    accessCode: storedAccessCode,
    isAdmin,
    isCodeVerified,
    hasAccess,
    dailyLimit,
    usedToday,
    remainingToday,
    isAccessModalOpen,
    accessError,
    verifying,
    syncAdminState,
    ensureAccess,
    acceptInviteCode,
    refreshUsage,
    clearAccessCode,
    openAccessModal,
    closeAccessModal,
    unlockWithCode,
    submitAccessRequest,
  };
}
