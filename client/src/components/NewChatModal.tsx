import React from 'react';
import { useState, useEffect } from 'react';
import { apiClient } from '../api';
import { User } from '../types';

interface NewChatModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreateChat: (companionId: number) => void;
}

export default function NewChatModal({
  isOpen,
  onClose,
  onCreateChat,
}: NewChatModalProps) {
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState<User[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchError, setSearchError] = useState('');

  // Search users with debounce
  useEffect(() => {
    if (!searchTerm.trim()) {
      setSearchResults([]);
      setSearchError('');
      return;
    }

    const timer = setTimeout(async () => {
      setIsSearching(true);
      setSearchError('');
      try {
        const response = await apiClient.searchUsers(searchTerm, 20);
        setSearchResults(response.users || []);
      } catch (error) {
        if (error instanceof Error && error.message.includes('404')) {
          setSearchResults([]);
        } else {
          setSearchError('Failed to search users');
          console.error('Search error:', error);
        }
      } finally {
        setIsSearching(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [searchTerm]);

  const handleCreate = () => {
    if (selectedUserId !== null) {
      onCreateChat(selectedUserId);
      onClose();
      setSelectedUserId(null);
      setSearchTerm('');
      setSearchResults([]);
    }
  };

  const handleClose = () => {
    onClose();
    setSelectedUserId(null);
    setSearchTerm('');
    setSearchResults([]);
    setSearchError('');
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-dark-800 rounded-2xl w-full max-w-md shadow-2xl border border-dark-700 animate-fade-in">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-dark-700">
          <h3 className="text-lg font-semibold text-white">New Chat</h3>
          <button
            onClick={handleClose}
            className="p-2 rounded-lg text-dark-400 hover:text-white hover:bg-dark-700 transition-all"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Search Input */}
        <div className="p-4 border-b border-dark-700">
          <label className="block text-sm font-medium text-dark-300 mb-2">
            Search by username
          </label>
          <div className="relative">
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Enter username..."
              className="w-full px-4 py-2.5 pl-10 bg-dark-900/50 border border-dark-600 rounded-xl text-white placeholder-dark-500 text-sm focus:outline-none focus:border-primary-500 transition-all"
              autoFocus
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
            {isSearching && (
              <div className="absolute right-3 top-1/2 -translate-y-1/2">
                <svg className="animate-spin h-4 w-4 text-primary-500" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
              </div>
            )}
          </div>
          <p className="text-dark-500 text-xs mt-2">
            Search for users to start a conversation
          </p>
        </div>

        {/* Search Results */}
        <div className="max-h-64 overflow-y-auto p-2">
          {searchError && (
            <div className="text-center py-4 text-red-400">
              <p>{searchError}</p>
            </div>
          )}
          {isSearching && searchResults.length === 0 && !searchError && (
            <div className="text-center py-4 text-dark-400">
              <p>Searching...</p>
            </div>
          )}
          {!isSearching && searchTerm.trim() && searchResults.length === 0 && !searchError && (
            <div className="text-center py-4 text-dark-400">
              <p>No users found</p>
            </div>
          )}
          {!searchTerm.trim() && (
            <div className="text-center py-4 text-dark-400">
              <p>Start typing to search for users</p>
            </div>
          )}
          {searchResults.length > 0 && (
            <div className="space-y-1">
              {searchResults.map((user) => {
                const userLogin = user.Login || user.login || 'Unknown';
                const userId = user.ID;
                return (
                  <div
                    key={userId}
                    onClick={() => setSelectedUserId(userId)}
                    className={`flex items-center gap-3 p-3 rounded-xl cursor-pointer transition-all ${
                      selectedUserId === userId
                        ? 'bg-primary-500/20 border border-primary-500/50'
                        : 'hover:bg-dark-700/50 border border-transparent'
                    }`}
                  >
                    <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center text-white font-semibold">
                      {userLogin.charAt(0).toUpperCase()}
                    </div>
                    <span className="text-white font-medium">{userLogin}</span>
                    {selectedUserId === userId && (
                      <svg
                        className="w-5 h-5 text-primary-400 ml-auto"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {/* Actions */}
        <div className="flex gap-3 p-4 border-t border-dark-700">
          <button
            onClick={handleClose}
            className="flex-1 px-4 py-2.5 rounded-xl border border-dark-600 text-dark-300 hover:bg-dark-700 transition-all font-medium"
          >
            Cancel
          </button>
          <button
            onClick={handleCreate}
            disabled={selectedUserId === null}
            className="flex-1 btn-primary py-2.5 rounded-xl text-white font-medium disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
          >
            Create Chat
          </button>
        </div>
      </div>
    </div>
  );
}
