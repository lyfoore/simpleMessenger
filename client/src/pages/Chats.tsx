import React from 'react';
import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import { apiClient } from '../api';
import { Chat, Message, User } from '../types';
import ChatList from '../components/ChatList';
import ChatWindow from '../components/ChatWindow';
import NewChatModal from '../components/NewChatModal';

export default function Chats() {
  const { user, token, logout } = useAuth();
  const navigate = useNavigate();

  const [chats, setChats] = useState<Chat[]>([]);
  const [messages, setMessages] = useState<Message[]>([]);
  const [selectedChat, setSelectedChat] = useState<Chat | null>(null);
  const [isMessagesLoading, setIsMessagesLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);

  // Load chats
  const loadChats = useCallback(async () => {
    try {
      console.log('Fetching chats...');
      const response = await apiClient.getChats();
      console.log('Chats response:', response);
      const chats = response.Chats || response.chats || [];
      setChats(chats);
    } catch (error) {
      console.error('Failed to load chats:', error);
    }
  }, []);

  // Load messages for selected chat
  const loadMessages = useCallback(async (chatId: number) => {
    setIsMessagesLoading(true);
    try {
      const response = await apiClient.getMessages(chatId);
      setMessages(response.Messages || response.messages || []);
    } catch (error) {
      console.error('Failed to load messages:', error);
    } finally {
      setIsMessagesLoading(false);
    }
  }, []);

  useEffect(() => {
    console.log('Chats: token=', token, 'user=', user);
    if (!token) {
      console.log('No token, navigating to login');
      navigate('/login');
      return;
    }
    console.log('Loading chats...');
    loadChats();
    setIsMessagesLoading(false);
  }, [token, navigate, loadChats]);

  useEffect(() => {
    if (selectedChat) {
      loadMessages(selectedChat.ID);
    } else {
      setMessages([]);
    }
  }, [selectedChat, loadMessages]);

  const handleSelectChat = (chat: Chat) => {
    setSelectedChat(chat);
  };

  const handleCreateChat = async (companionId: number) => {
    try {
      await apiClient.createChat(companionId);
      await loadChats();
    } catch (error) {
      console.error('Failed to create chat:', error);
      alert('Failed to create chat. User might not exist or chat already exists.');
    }
  };

  const handleDeleteChat = async (chatId: number) => {
    if (!confirm('Are you sure you want to delete this chat?')) return;
    try {
      await apiClient.deleteChat(chatId);
      if (selectedChat?.ID === chatId) {
        setSelectedChat(null);
      }
      await loadChats();
    } catch (error) {
      console.error('Failed to delete chat:', error);
    }
  };

  const handleSendMessage = async (text: string) => {
    if (!selectedChat || !user) return;

    // Send via WebSocket if connected
    if (ws && ws.readyState === WebSocket.OPEN) {
      const wsMessage = {
        text: text,
        chatId: selectedChat.ID,
        userId: user.ID,
      };
      ws.send(JSON.stringify(wsMessage));
    } else {
      // Fallback to REST API if WebSocket not connected
      // Message will be added when received via WebSocket
      try {
        await apiClient.sendMessage(selectedChat.ID, text);
      } catch (error) {
        console.error('Failed to send message:', error);
      }
    }
  };

  const handleDeleteMessage = async (messageId: number) => {
    if (!confirm('Delete this message?')) return;
    try {
      await apiClient.deleteMessage(messageId);
      setMessages((prev) => prev.filter((m) => m.ID !== messageId));
    } catch (error) {
      console.error('Failed to delete message:', error);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const getUserLogin = () => {
    if (!user) return '';
    return (user as any).Login || (user as any).login || user.Name || (user as any).name || '';
  };

  // WebSocket for real-time messages
  useEffect(() => {
    if (!selectedChat || !token) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host || 'localhost:3000';
    const wsUrl = `${protocol}//${host}/api/ws?token=${token}`;

    const wsConnection = new WebSocket(wsUrl);

    wsConnection.onopen = () => {
      console.log('WebSocket connected');
      setWs(wsConnection);
    };

    wsConnection.onclose = () => {
      console.log('WebSocket disconnected');
      setWs(null);
    };

    wsConnection.onmessage = (event) => {
      console.log('WebSocket message received:', event.data);
      try {
        const data = JSON.parse(event.data);
        console.log('Parsed WebSocket data:', data);
        
        // Backend sends message directly without type wrapper
        // Check if it's a message object (has ID and text)
        const newMessage: Message = data;
        
        if (newMessage.ID && (newMessage.text || newMessage.Text)) {
          const chatId = (newMessage as Message).ChatID || (newMessage as any).chat_id || (newMessage as any).chatId;
          console.log('New message for chat:', chatId, 'current selectedChat:', selectedChat);
          
          // Use functional update to get latest state
          setMessages((prev) => {
            console.log('Updating messages, prev length:', prev.length);
            
            // Check if message with this ID already exists
            const exists = prev.some((m) => m.ID === newMessage.ID);
            if (exists) {
              console.log('Message already exists, skipping');
              return prev;
            }
            
            console.log('Adding new message from WebSocket');
            return [...prev, newMessage];
          });
        } else {
          console.log('Not a message object, skipping', data);
        }
      } catch (e) {
        console.error('Failed to parse WebSocket message:', e);
      }
    };

    wsConnection.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => {
      if (wsConnection) {
        wsConnection.close();
      }
      setWs(null);
    };
  }, [selectedChat, token]);

  return (
    <div className="h-screen flex bg-dark-900">
      {/* Sidebar */}
      <div className="w-80 lg:w-96 flex-shrink-0">
        <ChatList
          chats={chats}
          selectedChatId={selectedChat?.ID || null}
          onSelectChat={handleSelectChat}
          onCreateChat={() => setIsModalOpen(true)}
          onDeleteChat={handleDeleteChat}
          currentUser={user ? { ID: user.ID, Login: getUserLogin() || '' } : null}
        />
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col min-w-0 relative">
        {/* Logout button - always visible */}
        <button
          onClick={handleLogout}
          className="absolute top-4 right-4 z-10 px-4 py-2 rounded-xl bg-dark-800 text-dark-300 hover:text-white hover:bg-dark-700 transition-all text-sm font-medium"
        >
          Logout
        </button>
        
        {selectedChat ? (
          <>
            <ChatWindow
              chatName={selectedChat.Name || selectedChat.name || 'Chat'}
              messages={messages}
              currentUser={user}
              onSendMessage={handleSendMessage}
              onDeleteMessage={handleDeleteMessage}
              isMessagesLoading={isMessagesLoading}
            />
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <div className="w-24 h-24 rounded-full bg-dark-800 flex items-center justify-center mx-auto mb-6">
                <svg
                  className="w-12 h-12 text-dark-500"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                  />
                </svg>
              </div>
              <h2 className="text-2xl font-bold text-white mb-2">Welcome to Messenger</h2>
              <p className="text-dark-400">Select a conversation to start messaging</p>
            </div>
          </div>
        )}
      </div>

      {/* New Chat Modal */}
      <NewChatModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onCreateChat={handleCreateChat}
      />
    </div>
  );
}
