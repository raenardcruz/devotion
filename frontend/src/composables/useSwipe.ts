import { onMounted, onUnmounted, type Ref } from 'vue';

interface SwipeOptions {
    onSwipeLeft?: (velocity?: number) => void;
    onSwipeRight?: (velocity?: number) => void;
    onSwipeUp?: (velocity?: number) => void;
    onSwipeDown?: (velocity?: number) => void;
    threshold?: number;
}

/**
 * Apple Design §6: Exponential momentum projection function
 */
export function projectMomentum(initialVelocity: number, decelerationRate = 0.998): number {
    return (initialVelocity / 1000) * decelerationRate / (1 - decelerationRate);
}

export function useSwipe(target: Ref<HTMLElement | null | undefined>, options: SwipeOptions) {
    const {
        onSwipeLeft,
        onSwipeRight,
        onSwipeUp,
        onSwipeDown,
        threshold = 40,
    } = options;

    let touchStartX = 0;
    let touchStartY = 0;
    let touchStartTime = 0;
    let touchEndX = 0;
    let touchEndY = 0;
    let touchEndTime = 0;

    const handleTouchStart = (e: TouchEvent) => {
        const touch = e.changedTouches[0];
        if (touch) {
            touchStartX = touch.screenX;
            touchStartY = touch.screenY;
            touchStartTime = performance.now();
        }
    };

    const handleTouchEnd = (e: TouchEvent) => {
        const touch = e.changedTouches[0];
        if (touch) {
            touchEndX = touch.screenX;
            touchEndY = touch.screenY;
            touchEndTime = performance.now();
            handleGesture();
        }
    };

    const handleGesture = () => {
        const deltaX = touchEndX - touchStartX;
        const deltaY = touchEndY - touchStartY;
        const deltaTime = Math.max(1, touchEndTime - touchStartTime); // in ms

        // Calculate gesture release velocity in pixels per second (Apple Design §5)
        const velocityX = (deltaX / deltaTime) * 1000;
        const velocityY = (deltaY / deltaTime) * 1000;

        if (Math.abs(deltaX) > Math.abs(deltaY)) {
            // Horizontal Swipe
            if (Math.abs(deltaX) > threshold) {
                if (deltaX > 0) {
                    onSwipeRight?.(velocityX);
                } else {
                    onSwipeLeft?.(velocityX);
                }
            }
        } else {
            // Vertical Swipe
            if (Math.abs(deltaY) > threshold) {
                if (deltaY > 0) {
                    onSwipeDown?.(velocityY);
                } else {
                    onSwipeUp?.(velocityY);
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
