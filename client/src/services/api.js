import axios from 'axios';

const api = axios.create({
    baseURL: window.__ENV_API_URL__ || process.env.REACT_APP_API_URL || 'http://localhost:8080/api',  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      window.dispatchEvent(new Event('unauthorized'));
    }
    return Promise.reject(error);
  }
);

export default api;