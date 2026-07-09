import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { AuthService } from '../services/auth.service';
import { TOKEN_KEY } from '../services/api';

interface AuthContextData {
  user: any;
  loading: boolean;
  login: (email: string, pass: string) => Promise<{ success: boolean; error?: string }>;
  register: (username: string, email: string, pass: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextData | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    checkSession();
  }, []);

  const checkSession = async () => {
    setLoading(true);
    try {
      const token = localStorage.getItem(TOKEN_KEY);
      if (token) {
        const response = await AuthService.getMe();
        if (response.success) {
          setUser(response.data);
        } else {
          localStorage.removeItem(TOKEN_KEY);
        }
      }
    } catch (e) {
      console.error('Error checking session', e);
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

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
