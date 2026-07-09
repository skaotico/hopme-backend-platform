import { fetchWithAuth, API_AUTH_URL, TOKEN_KEY } from './api';

const AsyncStorage = {
  getItem: async (key: string) => localStorage.getItem(key),
  setItem: async (key: string, value: string) => localStorage.setItem(key, value),
  removeItem: async (key: string) => localStorage.removeItem(key),
};
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
    const url = `${API_AUTH_URL}/auth/login`;
    const body = JSON.stringify({ email, password });
    
    const curlCommand = `curl -X POST ${url} \\\n  -H "Content-Type: application/json" \\\n  -d '${body}'`;
    console.log('--- cURL DE LOGIN ---');
    console.log(curlCommand);

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body,
      });
      
      const responseText = await response.text();
      let result;
      try {
        result = JSON.parse(responseText);
      } catch (e) {
        result = { error: 'No se pudo parsear el JSON de respuesta.' };
      }
      
      if (result.success && result.data?.token) {
        await AsyncStorage.setItem(TOKEN_KEY, result.data.token);
      } else {
        result.error = `${result.error || 'Credenciales incorrectas'}\n\n[RQ Headers]: Content-Type: application/json\n[RQ Body]: ${body}\n\n[RS Status]: ${response.status}\n[RS Body]: ${responseText.substring(0, 150)}\n\ncURL generado:\n${curlCommand}`;
      }
      
      return result;
    } catch (e: any) {
      throw new Error(`${e.message}\n\n[RQ Headers]: Content-Type: application/json\n[RQ Body]: ${body}\n\n[RS Status]: NO HUBO RESPUESTA (Falló antes de llegar al servidor)\n\ncURL generado:\n${curlCommand}`);
    }
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
