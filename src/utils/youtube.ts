/**
 * Utility for parsing YouTube URLs and returning embed-friendly links.
 */
export const getYouTubeEmbedUrl = (url: string | undefined): string => {
  if (!url) return '';
  let videoId = '';
  
  try {
    // Try youtu.be format first (e.g. https://youtu.be/DaWRPjNztTs?si=ISoMyr2dWhypBiQ_)
    if (url.includes('youtu.be/')) {
      const parts = url.split('youtu.be/');
      if (parts[1]) {
        videoId = parts[1].split('?')[0] || '';
      }
    } 
    // Try youtube.com/watch format (e.g. https://www.youtube.com/watch?v=DaWRPjNztTs)
    else if (url.includes('youtube.com/watch')) {
      const urlObj = new URL(url);
      videoId = urlObj.searchParams.get('v') || '';
    } 
    // Try youtube.com/embed format
    else if (url.includes('youtube.com/embed/')) {
      const parts = url.split('youtube.com/embed/');
      if (parts[1]) {
        videoId = parts[1].split('?')[0] || '';
      }
    }
  } catch (e) {
    console.error('Failed to parse YouTube URL:', url, e);
  }
  
  return videoId ? `https://www.youtube.com/embed/${videoId}?autoplay=1` : '';
};
