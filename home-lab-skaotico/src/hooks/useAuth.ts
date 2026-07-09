import { useState, useEffect } from 'react';
import { AuthService } from '../services/auth.service';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { TOKEN_KEY } from '../services/api';

export function useAuth() {
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    // Intentar cargar sesión al iniciar
    checkSession();
  }, []);

  const checkSession = async () => {
    setLoading(true);
    try {
      const token = await AsyncStorage.getItem(TOKEN_KEY);
      if (token) {
        const response = await AuthService.getMe();
        if (response.success) {
          setUser(response.data);
        } else {
          // Si el token es inválido, limpiarlo
          await AsyncStorage.removeItem(TOKEN_KEY);
        }
      }
    } catch (e) {
      console.log('Error checking session', e);
    } finally {
      setLoading(false);
    }
  };

  const login = async (email: string, pass: string) => {
    setLoading(true);
    try {
      const result = await AuthService.login(email, pass);
      if (result.success) {
        await checkSession();
        return { success: true };
      }
      return { success: false, error: result.error || 'Login failed' };
    } catch (e: any) {
      return { success: false, error: e.message };
    } finally {
      setLoading(false);
    }
  };

  const register = async (username: string, email: string, pass: string) => {
    setLoading(true);
    try {
      const result = await AuthService.register(username, email, pass);
      if (result.success) {
        return { success: true };
      }
      return { success: false, error: result.error || 'Register failed' };
    } catch (e: any) {
      return { success: false, error: e.message };
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    await AuthService.logout();
    setUser(null);
  };

  return {
    user,
    loading,
    login,
    register,
    logout,
  };
}
