import React from 'react';
import { useState } from 'react';
import { Chat } from '../types';

interface ChatListProps {
  chats: Chat[];
  selectedChatId: number | null;
  onSelectChat: (chat: Chat) => void;
  onCreateChat: () => void;
  onDeleteChat: (chatId: number) => void;
  currentUser: { ID: number; Login?: string } | null;
}

export default function ChatList({
  chats,
  selectedChatId,
  onSelectChat,
  onCreateChat,
  onDeleteChat,
  currentUser,
}: ChatListProps) {
  const [searchTerm, setSearchTerm] = useState('');

  const formatLastMessageTime = (dateString?: string) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const hours = Math.floor(diff / (1000 * 60 * 60));

    if (hours < 1) return 'Just now';
    if (hours < 24) return `${hours}h ago`;
    if (hours < 48) return 'Yesterday';
    return date.toLocaleDateString();
  };

  const filteredChats = chats.filter((chat) => {
    const chatName = (chat.Name || chat.name || '').toLowerCase();
    return chatName.includes(searchTerm.toLowerCase());
  });

  return (
    <div className="h-full flex flex-col bg-dark-900/30 backdrop-blur-xl border-r border-dark-700">
      {/* Header */}
      <div className="p-4 border-b border-dark-700">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-white">Messages</h2>
          <button
            onClick={onCreateChat}
            className="p-2 rounded-xl bg-primary-500/20 text-primary-400 hover:bg-primary-500/30 transition-all"
            title="New chat"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
          </button>
        </div>
        <div className="relative">
          <input
            type="text"
            placeholder="Search conversations..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full px-4 py-2.5 pl-10 bg-dark-800/50 border border-dark-600 rounded-xl text-white placeholder-dark-500 text-sm focus:outline-none focus:border-primary-500 transition-all"
          />
          <svg
            className="w-4 h-4 text-dark-500 absolute left-3 top-1/2 -translate-y-1/2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
      </div>

      {/* Chat List */}
      <div className="flex-1 overflow-y-auto">
        {filteredChats.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center p-6">
            <div className="w-16 h-16 rounded-full bg-dark-800 flex items-center justify-center mb-4">
              <svg className="w-8 h-8 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
            </div>
            {searchTerm ? (
              <>
                <p className="text-dark-400 font-medium mb-1">No matches found</p>
                <p className="text-dark-500 text-sm">Try a different search term</p>
              </>
            ) : (
              <>
                <p className="text-dark-400 font-medium mb-1">No conversations yet</p>
                <p className="text-dark-500 text-sm">Start a new chat to begin messaging</p>
              </>
            )}
          </div>
        ) : (
          <div className="p-2 space-y-1">
            {filteredChats.map((chat) => (
              <div
                key={chat.ID}
                onClick={() => onSelectChat(chat)}
                className={`chat-item p-3 rounded-xl cursor-pointer flex items-center gap-3 group ${
                  selectedChatId === chat.ID ? 'active' : ''
                }`}
              >
                <div className="relative flex-shrink-0">
                  <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center text-white font-semibold">
                    {((chat as Chat).Name || (chat as any).name || '?').charAt(0).toUpperCase()}
                  </div>
                  <div className="online-indicator absolute bottom-0 right-0"></div>
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <h3 className="text-white font-medium truncate">{(chat as Chat).Name || (chat as any).name || 'Unknown'}</h3>
                    <span className="text-dark-500 text-xs flex-shrink-0">
                      {formatLastMessageTime((chat as Chat).LastMessageAt || (chat as any).last_message_at)}
                    </span>
                  </div>
                  <p className="text-dark-400 text-sm truncate">Click to open chat</p>
                </div>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onDeleteChat(chat.ID);
                  }}
                  className="p-1.5 rounded-lg text-dark-500 hover:text-red-400 hover:bg-red-500/10 opacity-0 group-hover:opacity-100 transition-all"
                  title="Delete chat"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                  </svg>
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Current User */}
      {currentUser && (
        <div className="p-4 border-t border-dark-700">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center text-white font-semibold">
              {(currentUser.Login || '?').charAt(0).toUpperCase()}
            </div>
            <div className="flex-1">
              <p className="text-white font-medium">{currentUser.Login || 'Unknown'}</p>
              <p className="text-dark-500 text-xs">Online</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
