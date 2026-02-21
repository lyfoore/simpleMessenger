import React from 'react';
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { apiClient } from '../api';
import { User } from '../types';

interface AuthContextType {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  login: (username: string) => Promise<void>;
  register: (username: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  validateToken: () => Promise<boolean>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    apiClient.setToken(null);
  };

  // Validate token on mount
  useEffect(() => {
    const validateAuth = async () => {
      const savedToken = localStorage.getItem('token');
      const savedUser = localStorage.getItem('user');

      if (!savedToken || !savedUser) {
        setIsLoading(false);
        return;
      }

      try {
        setUser(JSON.parse(savedUser));
        apiClient.setToken(savedToken);

        // Validate token by calling /api/me
        await apiClient.getMe();
        setToken(savedToken);
      } catch (e) {
        console.error('Token validation failed:', e);
        logout();
      } finally {
        setIsLoading(false);
      }
    };

    validateAuth();
  }, []);

  // Set up unauthorized callback
  useEffect(() => {
    apiClient.setOnUnauthorized(() => {
      console.log('Token expired, logging out');
      logout();
    });
  }, []);

  const login = async (username: string) => {
    const response = await apiClient.login(username);
    const token = response.Token || (response as any).token;
    const user = response.User || (response as any).user;
    setToken(token);
    setUser(user);
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    apiClient.setToken(token);
  };

  const register = async (username: string) => {
    const response = await apiClient.register(username);
    const token = response.Token || (response as any).token;
    const user = response.User || (response as any).user;
    setToken(token);
    setUser(user);
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    apiClient.setToken(token);
  };

  const validateToken = async (): Promise<boolean> => {
    if (!token) return false;
    try {
      await apiClient.getMe();
      return true;
    } catch {
      return false;
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isLoading,
        login,
        register,
        logout,
        isAuthenticated: !!token && !!user,
        validateToken,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
