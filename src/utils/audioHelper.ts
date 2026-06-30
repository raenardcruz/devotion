/**
 * Utility for resolving audio files from src/data/ and managing prayer audio playback.
 */

// A map of prayer IDs that have confirmed audio files
export const AVAILABLE_AUDIO_PRAYERS = new Set([
  'sign-of-the-cross',
  'our-father',
  'hail-mary',
  'glory-be',
  'apostles-creed',
  'hail-holy-queen',
  'memorare',
  'angel-of-god',
  'st-michael-the-archangel',
  'fatima-prayer',
  'dm-opening-1',
  'dm-opening-2',
  'eternal-father',
  'sorrowful-passion',
  'holy-god',
  'dm-closing',
  'rosary-closing-prayer',
  'jesus-i-trust-in-you',
  'act-of-contrition'
]);

/**
 * Returns the resolved asset URL of the audio file for the given prayer ID.
 * Uses Vite's dynamic URL import mechanism.
 */
export const getPrayerAudioUrl = (prayerId: string): string => {
  if (!AVAILABLE_AUDIO_PRAYERS.has(prayerId)) {
    return '';
  }
  try {
    // Resolve relative path to src/data/ from this utility file (src/utils/)
    return new URL(`../data/${prayerId}.wav`, import.meta.url).href;
  } catch (e) {
    console.error(`Failed to resolve audio for prayer: ${prayerId}`, e);
    return '';
  }
};
