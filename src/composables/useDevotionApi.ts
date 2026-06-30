import { ref } from 'vue';

export function useDevotionApi() {
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchWithAuth = async <T>(endpoint: string, options: RequestInit = {}): Promise<T> => {
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'https://devotionapi.raenardcruz.com';
    const apiToken = import.meta.env.VITE_API_TOKEN;

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
    fetchWithAuth,
    getDevotion,
  };
}
