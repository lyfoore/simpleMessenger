import React from 'react';
import { useState, useRef, useEffect } from 'react';
import { Message, User } from '../types';

interface ChatWindowProps {
  chatName: string;
  messages: Message[];
  currentUser: User | null;
  onSendMessage: (text: string) => void;
  onDeleteMessage: (messageId: number) => void;
  isMessagesLoading?: boolean;
}

export default function ChatWindow({
  chatName,
  messages,
  currentUser,
  onSendMessage,
  onDeleteMessage,
  isMessagesLoading = false,
}: ChatWindowProps) {
  const [newMessage, setNewMessage] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const prevMessageCount = useRef<number>(0);
  const isInitialLoad = useRef<boolean>(true);
  const prevChatName = useRef<string>(chatName);

  // Get current user ID (handle both PascalCase and lowercase)
  const currentUserId = currentUser?.ID || (currentUser as any)?.id;

  // Reset initial load flag when chat changes
  if (chatName !== prevChatName.current) {
    prevChatName.current = chatName;
    isInitialLoad.current = true;
    prevMessageCount.current = 0;
  }

  // Check if user is at bottom of chat
  const isAtBottom = () => {
    const container = messagesContainerRef.current;
    if (!container) return true;
    const threshold = 100; // pixels from bottom
    return container.scrollHeight - container.scrollTop - container.clientHeight < threshold;
  };

  const scrollToBottom = (smooth = true) => {
    messagesEndRef.current?.scrollIntoView({ behavior: smooth ? 'smooth' : 'auto' });
  };

  // Scroll on new messages
  useEffect(() => {
    if (isInitialLoad.current) {
      // First load - scroll to bottom without animation
      isInitialLoad.current = false;
      prevMessageCount.current = messages.length;
      scrollToBottom(false);
      return;
    }

    const messageDiff = messages.length - prevMessageCount.current;

    if (messageDiff > 0) {
      // New messages received - always scroll to bottom smoothly
      scrollToBottom(true);
    }

    prevMessageCount.current = messages.length;
  }, [messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newMessage.trim()) {
      onSendMessage(newMessage.trim());
      setNewMessage('');
    }
  };

  const formatTime = (dateString?: string) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) return 'Today';
    if (days === 1) return 'Yesterday';
    if (days < 7) return date.toLocaleDateString([], { weekday: 'long' });
    return date.toLocaleDateString();
  };

  const shouldShowDate = (index: number): boolean => {
    if (index === 0) return true;
    const prevMsg = messages[index - 1];
    const currMsg = messages[index];
    const prevDate = (prevMsg as Message).CreatedAt || (prevMsg as any).created_at;
    const currDate = (currMsg as Message).CreatedAt || (currMsg as any).created_at;
    return formatDate(prevDate) !== formatDate(currDate);
  };

  return (
    <div className="h-full flex flex-col bg-dark-900/20">
      {/* Header */}
      <div className="px-6 py-4 border-b border-dark-700 bg-dark-900/30 backdrop-blur-xl">
        <div className="flex items-center gap-4">
          <div className="relative">
            <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center text-white font-semibold text-lg">
              {chatName.charAt(0).toUpperCase()}
            </div>
            <div className="online-indicator absolute bottom-0 right-0"></div>
          </div>
          <div>
            <h3 className="text-white font-semibold text-lg">{chatName || 'Chat'}</h3>
            <p className="text-dark-400 text-sm">Online</p>
          </div>
        </div>
      </div>

      {/* Messages */}
      <div ref={messagesContainerRef} className="flex-1 overflow-y-auto p-6 space-y-4">
        {isMessagesLoading ? (
          <div className="flex items-center justify-center h-full">
            <div className="flex gap-2">
              <div className="typing-dot"></div>
              <div className="typing-dot"></div>
              <div className="typing-dot"></div>
            </div>
          </div>
        ) : messages.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <div className="w-20 h-20 rounded-full bg-dark-800 flex items-center justify-center mb-4">
              <svg className="w-10 h-10 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
            </div>
            <p className="text-white font-medium mb-1">No messages yet</p>
            <p className="text-dark-400 text-sm">Say hello to start the conversation!</p>
          </div>
        ) : (
          <>
            {messages.map((message, index) => (
              <React.Fragment key={message.ID}>
                {shouldShowDate(index) && (
                  <div className="flex justify-center">
                    <span className="px-3 py-1 bg-dark-800 rounded-full text-dark-400 text-xs">
                      {formatDate(message.CreatedAt || (message as any).created_at)}
                    </span>
                  </div>
                )}
                <div
                  className={`flex ${(message as Message).UserID === currentUserId || (message as any).user_id === currentUserId || (message as any).userId === currentUserId ? 'justify-end' : 'justify-start'} animate-message-in`}
                >
                  <div className="max-w-[70%]">
                    <div
                      className={`px-4 py-2.5 ${(message as Message).UserID === currentUserId || (message as any).user_id === currentUserId || (message as any).userId === currentUserId ? 'message-own' : 'message-other'
                      }`}
                    >
                      <p className="text-sm whitespace-pre-wrap break-words">{(message as Message).Text || (message as any).text}</p>
                    </div>
                    <div
                      className={`flex items-center gap-2 mt-1 ${(message as Message).UserID === currentUserId || (message as any).user_id === currentUserId || (message as any).userId === currentUserId ? 'justify-end' : 'justify-start'
                      }`}
                    >
                      <span className="text-dark-500 text-xs">{formatTime((message as Message).CreatedAt || (message as any).created_at)}</span>
                      {((message as Message).UserID === currentUserId || (message as any).user_id === currentUserId || (message as any).userId === currentUserId) && (
                        <button
                          onClick={() => onDeleteMessage(message.ID)}
                          className="text-dark-500 hover:text-red-400 transition-colors"
                          title="Delete message"
                        >
                          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M6 18L18 6M6 6l12 12"
                            />
                          </svg>
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              </React.Fragment>
            ))}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      {/* Input */}
      <div className="p-4 border-t border-dark-700 bg-dark-900/30 backdrop-blur-xl">
        <form onSubmit={handleSubmit} className="flex items-center gap-3">
          <input
            ref={inputRef}
            type="text"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type a message..."
            className="flex-1 px-4 py-3 bg-dark-800/50 border border-dark-600 rounded-xl text-white placeholder-dark-500 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/20 transition-all input-focus"
          />
          <button
            type="submit"
            disabled={!newMessage.trim()}
            className="p-3 rounded-xl btn-primary disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
          >
            <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
              />
            </svg>
          </button>
        </form>
      </div>
    </div>
  );
}
