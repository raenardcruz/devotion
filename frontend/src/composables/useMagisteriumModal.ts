import { ref } from 'vue';
import { useAccessControl } from './useAccessControl';

const isMagisteriumModalOpen = ref(false);
const sharedSessionId = ref<string | null>(null);

export function useMagisteriumModal() {
  const openMagisteriumModal = async (sessionId?: string) => {
    const { ensureAccess, openAccessModal } = useAccessControl();
    if (sessionId) {
      sharedSessionId.value = sessionId;
    }
    const allowed = await ensureAccess();
    if (allowed) {
      isMagisteriumModalOpen.value = true;
    } else {
      openAccessModal(() => {
        if (sessionId) {
          sharedSessionId.value = sessionId;
        }
        isMagisteriumModalOpen.value = true;
      });
    }
  };

  const closeMagisteriumModal = () => {
    isMagisteriumModalOpen.value = false;
  };

  const clearSharedSessionId = () => {
    sharedSessionId.value = null;
  };

  return {
    isMagisteriumModalOpen,
    sharedSessionId,
    openMagisteriumModal,
    closeMagisteriumModal,
    clearSharedSessionId,
  };
}
