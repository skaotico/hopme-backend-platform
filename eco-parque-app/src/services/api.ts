import AsyncStorage from '@react-native-async-storage/async-storage';

export const API_AUTH_URL = 'http://192.168.1.27:9090/api/v1';
export const API_PARQUE_URL = 'http://192.168.1.27:9091/api/v1';
export const API_ZONA_URL = 'http://192.168.1.27:9092/api/v1';
export const API_FLORA_URL = 'http://192.168.1.27:9093/api/v1';

export const TOKEN_KEY = '@eco_parque_access_token';

/**
 * Realiza una petición `fetch` inyectando automáticamente el token JWT
 * si existe en AsyncStorage.
 */
export async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const token = await AsyncStorage.getItem(TOKEN_KEY);
  
  const headers = new Headers(options.headers || {});
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }
  // Por defecto enviamos JSON si el método tiene body
  if (options.method && options.method !== 'GET' && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  return response;
}
