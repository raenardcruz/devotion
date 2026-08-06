import { ref, onMounted, onUnmounted, type Ref } from 'vue';

export function useFullscreen(targetRef?: Ref<HTMLElement | null>) {
  const isFullscreen = ref(false);

  const checkFullscreen = () => {
    isFullscreen.value = !!(
      document.fullscreenElement ||
      (document as any).webkitFullscreenElement ||
      (document as any).msFullscreenElement
    );
  };

  const toggleFullscreen = async (elementOverride?: HTMLElement | null) => {
    try {
      const target = elementOverride || targetRef?.value || document.documentElement;
      if (!isFullscreen.value) {
        if (target.requestFullscreen) {
          await target.requestFullscreen();
        } else if ((target as any).webkitRequestFullscreen) {
          await (target as any).webkitRequestFullscreen();
        } else if ((target as any).msRequestFullscreen) {
          await (target as any).msRequestFullscreen();
        }
      } else {
        if (document.exitFullscreen) {
          await document.exitFullscreen();
        } else if ((document as any).webkitExitFullscreen) {
          await (document as any).webkitExitFullscreen();
        } else if ((document as any).msExitFullscreen) {
          await (document as any).msExitFullscreen();
        }
      }
    } catch (err) {
      console.error('Error toggling fullscreen mode:', err);
    }
  };

  const enterFullscreen = async (elementOverride?: HTMLElement | null) => {
    if (isFullscreen.value) return;
    await toggleFullscreen(elementOverride);
  };

  const exitFullscreen = async () => {
    if (!isFullscreen.value) return;
    await toggleFullscreen();
  };

  onMounted(() => {
    document.addEventListener('fullscreenchange', checkFullscreen);
    document.addEventListener('webkitfullscreenchange', checkFullscreen);
    document.addEventListener('msfullscreenchange', checkFullscreen);
  });

  onUnmounted(() => {
    document.removeEventListener('fullscreenchange', checkFullscreen);
    document.removeEventListener('webkitfullscreenchange', checkFullscreen);
    document.removeEventListener('msfullscreenchange', checkFullscreen);
  });

  return {
    isFullscreen,
    toggleFullscreen,
    enterFullscreen,
    exitFullscreen
  };
}
