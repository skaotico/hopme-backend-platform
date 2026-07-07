import { fetchWithAuth, API_AUTH_URL, TOKEN_KEY } from './api';
import AsyncStorage from '@react-native-async-storage/async-storage';

export const AuthService = {
  /**
   * Registra un nuevo usuario.
   */
  async register(username: string, email: string, password: string) {
    const response = await fetch(`${API_AUTH_URL}/auth/register`, {
      method: 'POST',
      body: JSON.stringify({ username, email, password }),
    });
    return response.json();
  },

  /**
   * Inicia sesión con el email y contraseña.
   * Si es exitoso, guarda el token devuelto.
   */
  async login(email: string, password: string) {
    const response = await fetch(`${API_AUTH_URL}/auth/login`, {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    const result = await response.json();
    
    if (result.success && result.data?.token) {
      await AsyncStorage.setItem(TOKEN_KEY, result.data.token);
    }
    
    return result;
  },

  /**
   * Obtiene la información del usuario autenticado actual.
   */
  async getMe() {
    const response = await fetchWithAuth(`${API_AUTH_URL}/auth/me`, {
      method: 'GET',
    });
    return response.json();
  },

  /**
   * Cierra la sesión (revoca token en servidor) y limpia AsyncStorage.
   */
  async logout() {
    // Intentar logout en el servidor (si expira, no importa, eliminamos local)
    try {
      await fetchWithAuth(`${API_AUTH_URL}/auth/logout`, {
        method: 'POST',
      });
    } catch (e) {
      console.warn('Error al hacer logout en servidor', e);
    }

    await AsyncStorage.removeItem(TOKEN_KEY);
  }
};
