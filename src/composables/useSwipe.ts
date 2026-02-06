import { onMounted, onUnmounted, type Ref } from 'vue';

interface SwipeOptions {
    onSwipeLeft?: () => void;
    onSwipeRight?: () => void;
    onSwipeUp?: () => void;
    onSwipeDown?: () => void;
    threshold?: number;
}

export function useSwipe(target: Ref<HTMLElement | null | undefined>, options: SwipeOptions) {
    const {
        onSwipeLeft,
        onSwipeRight,
        onSwipeUp,
        onSwipeDown,
        threshold = 50,
    } = options;

    let touchStartX = 0;
    let touchStartY = 0;
    let touchEndX = 0;
    let touchEndY = 0;

    const handleTouchStart = (e: TouchEvent) => {
        const touch = e.changedTouches[0];
        if (touch) {
            touchStartX = touch.screenX;
            touchStartY = touch.screenY;
        }
    };

    const handleTouchEnd = (e: TouchEvent) => {
        const touch = e.changedTouches[0];
        if (touch) {
            touchEndX = touch.screenX;
            touchEndY = touch.screenY;
            handleGesture();
        }
    };

    const handleGesture = () => {
        const deltaX = touchEndX - touchStartX;
        const deltaY = touchEndY - touchStartY;

        if (Math.abs(deltaX) > Math.abs(deltaY)) {
            // Horizontal Swipe
            if (Math.abs(deltaX) > threshold) {
                if (deltaX > 0) {
                    onSwipeRight?.();
                } else {
                    onSwipeLeft?.();
                }
            }
        } else {
            // Vertical Swipe
            if (Math.abs(deltaY) > threshold) {
                if (deltaY > 0) {
                    onSwipeDown?.();
                } else {
                    onSwipeUp?.();
                }
            }
        }
    };

    onMounted(() => {
        if (target.value) {
            target.value.addEventListener('touchstart', handleTouchStart, { passive: true });
            target.value.addEventListener('touchend', handleTouchEnd, { passive: true });
        }
    });

    onUnmounted(() => {
        if (target.value) {
            target.value.removeEventListener('touchstart', handleTouchStart);
            target.value.removeEventListener('touchend', handleTouchEnd);
        }
    });
}
