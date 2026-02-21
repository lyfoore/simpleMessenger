export interface User {
  ID: number;
  Login?: string;
  login?: string;
  Name?: string;
  name?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
}

export interface Chat {
  ID: number;
  Name?: string;
  name?: string;
  LastMessageAt?: string;
  last_message_at?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
}

export interface Message {
  ID: number;
  Text?: string;
  text?: string;
  ChatID?: number;
  chat_id?: number;
  chatId?: number;
  UserID?: number;
  user_id?: number;
  userId?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
  DeletedAt?: string | null;
  User?: User;
}

export interface AuthResponse {
  User?: User;
  user?: User;
  Token?: string;
  token?: string;
  ExpiresIn?: number;
  expiresIn?: number;
}

export interface LoginRequest {
  Username: string;
}

export interface RegisterRequest {
  Username: string;
}

export interface CreateChatRequest {
  CompanionID: number;
}

export interface SendMessageRequest {
  Text: string;
}

export interface GetChatsResponse {
  chats?: Chat[];
  Chats?: Chat[];
}

export interface GetMessagesResponse {
  messages?: Message[];
  Messages?: Message[];
}
