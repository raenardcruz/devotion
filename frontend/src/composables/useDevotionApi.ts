import { ref } from 'vue';

const isServerUnreachable = ref(false);

export function useDevotionApi() {
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchWithAuth = async <T>(endpoint: string, options: RequestInit = {}): Promise<T> => {
    let baseUrl = import.meta.env.VITE_API_BASE_URL || 'https://devotionapi.raenardcruz.com';
    const apiToken = import.meta.env.VITE_API_TOKEN;

    // Validate and fall back if VITE_API_BASE_URL is invalid (e.g. set to the API token value)
    if (!baseUrl || (!baseUrl.startsWith('http://') && !baseUrl.startsWith('https://') && !baseUrl.startsWith('/'))) {
      console.warn(`Invalid VITE_API_BASE_URL "${baseUrl}" provided. Falling back to default.`);
      baseUrl = 'https://devotionapi.raenardcruz.com';
    }

    // Build absolute URL
    const url = endpoint.startsWith('http') ? endpoint : `${baseUrl.replace(/\/$/, '')}/${endpoint.replace(/^\//, '')}`;

    const headers = new Headers(options.headers);
    if (apiToken) {
      headers.set('Authorization', `Bearer ${apiToken}`);
    }
    headers.set('Accept', 'application/json');

    loading.value = true;
    error.value = null;

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      isServerUnreachable.value = false;

      if (response.status === 401) {
        throw new Error('Unauthorized access. Please check your API Token configuration.');
      }

      if (!response.ok) {
        let errorMessage = `Request failed: ${response.status} ${response.statusText}`;
        try {
          const data = await response.json();
          if (data && data.error) {
            errorMessage = data.error;
          }
        } catch {
          // If response isn't JSON, use status text
        }
        throw new Error(errorMessage);
      }

      return await response.json() as T;
    } catch (err: any) {
      if (
        err.name === 'TypeError' ||
        (err.message &&
          (err.message.includes('fetch') ||
            err.message.includes('NetworkError') ||
            err.message.includes('Failed to fetch') ||
            err.message.includes('Load failed')))
      ) {
        isServerUnreachable.value = true;
        const msg = 'Server cannot be reached. Please check if the API server is running and accessible.';
        error.value = msg;
        throw new Error(msg);
      }
      error.value = err.message || 'An error occurred during the API request.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const getDevotion = async <T>(date?: string): Promise<T> => {
    const query = date ? `?date=${date}` : '';
    return fetchWithAuth<T>(`/devotion${query}`);
  };

  return {
    loading,
    error,
    isServerUnreachable,
    fetchWithAuth,
    getDevotion,
  };
}
