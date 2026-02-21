import {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  GetChatsResponse,
  GetMessagesResponse,
  SendMessageRequest,
  User,
} from './types';

const API_BASE_URL = '';

type LogoutCallback = () => void;

class ApiClient {
  private token: string | null = null;
  private onUnauthorized: LogoutCallback | null = null;

  setToken(token: string | null) {
    this.token = token;
  }

  getToken(): string | null {
    return this.token;
  }

  setOnUnauthorized(callback: LogoutCallback | null) {
    this.onUnauthorized = callback;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = new Headers({
      'Content-Type': 'application/json',
    });

    if (this.token) {
      headers.append('Authorization', `Bearer ${this.token}`);
    }

    if (options.headers) {
      const extraHeaders = options.headers as Record<string, string>;
      Object.entries(extraHeaders).forEach(([key, value]) => {
        headers.append(key, value);
      });
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      // Handle 401 Unauthorized - token expired or invalid
      if (response.status === 401) {
        console.warn('Unauthorized request - token may be expired');
        if (this.onUnauthorized) {
          this.onUnauthorized();
        }
      }
      const error = await response.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // Auth endpoints
  async login(username: string): Promise<AuthResponse> {
    const data: LoginRequest = { Username: username };
    return this.request('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async register(username: string): Promise<AuthResponse> {
    const data: RegisterRequest = { Username: username };
    return this.request('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async getMe(): Promise<{ message: string; user_id: number }> {
    return this.request('/api/me');
  }

  // User endpoints
  async searchUsers(search: string, limit: number = 20): Promise<{ users: User[] }> {
    return this.request(`/api/users?search=${encodeURIComponent(search)}&limit=${limit}`);
  }

  async getUserByLogin(login: string): Promise<User> {
    return this.request(`/api/users?login=${encodeURIComponent(login)}`);
  }

  async getUserById(id: number): Promise<User> {
    return this.request(`/api/users?id=${id}`);
  }

  // Chat endpoints
  async getChats(limit: number = 50): Promise<GetChatsResponse> {
    return this.request(`/api/chats?limit=${limit}`);
  }

  async createChat(companionId: number): Promise<{ message: string }> {
    return this.request('/api/chats', {
      method: 'POST',
      body: JSON.stringify({ CompanionID: companionId }),
    });
  }

  async deleteChat(chatId: number): Promise<{ message: string }> {
    return this.request(`/api/chats/${chatId}`, {
      method: 'DELETE',
    });
  }

  // Message endpoints
  async getMessages(chatId: number, limit: number = 50): Promise<GetMessagesResponse> {
    return this.request(`/api/chats/${chatId}/messages?limit=${limit}`);
  }

  async sendMessage(chatId: number, text: string): Promise<{ message: string }> {
    const data: SendMessageRequest = { Text: text };
    return this.request(`/api/chats/${chatId}/messages`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async deleteMessage(messageId: number): Promise<{ message: string }> {
    return this.request(`/api/messages/${messageId}`, {
      method: 'DELETE',
    });
  }

  // WebSocket URL
  getWebSocketUrl(): string {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host || 'localhost:3000';
    const token = this.token ? `?token=${this.token}` : '';
    return `${protocol}//${host}/api/ws${token}`;
  }
}

export const apiClient = new ApiClient();
