import React, { createContext, useState, useContext, useEffect } from 'react';
import api from '../services/api';

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (token) {
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      api.get('/me')
        .then(response => {
          setUser({ id: response.data.user_id, login: response.data.login }); // login может отсутствовать, но сохраним id
        })
        .catch(() => {
          logout();
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    const handleUnauthorized = () => logout();
    window.addEventListener('unauthorized', handleUnauthorized);
    return () => window.removeEventListener('unauthorized', handleUnauthorized);
  }, []);

  const login = async (login) => {
    try {
      const response = await api.post('/auth/login', { username: login });
      const token = response.data.Token;
      localStorage.setItem('token', token);
      setToken(token);
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      const meResponse = await api.get('/me');
      setUser({ id: meResponse.data.user_id, login });
      return { success: true };
    } catch (error) {
      return { success: false, error: error.response?.data?.message || 'Login failed' };
    }
};

  const register = async (login, name) => {
    try {
      const response = await api.post('/auth/register', { login, name });
      const { token } = response.data;
      localStorage.setItem('token', token);
      setToken(token);
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      const meResponse = await api.get('/me');
      setUser({ id: meResponse.data.user_id, login, name });
      return { success: true };
    } catch (error) {
      return { success: false, error: error.response?.data?.message || 'Registration failed' };
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    delete api.defaults.headers.common['Authorization'];
    setToken(null);
    setUser(null);
  };

  const value = {
    user,
    token,
    login,
    register,
    logout,
    loading
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};