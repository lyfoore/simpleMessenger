import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';
import ChatList from '../components/ChatList';

const ChatsPage = () => {
  const [chats, setChats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [newChatName, setNewChatName] = useState('');
  const { logout } = useAuth();
  const navigate = useNavigate();

  const fetchChats = async () => {
    try {
      const response = await api.get('/chats?limit=50');
      setChats(response.data.chats);
    } catch (err) {
      setError('Failed to load chats');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchChats();
  }, []);

  const handleCreateChat = async (e) => {
    e.preventDefault();
    if (!newChatName.trim()) return;
    try {
      await api.post('/chats', { name: newChatName });
      setNewChatName('');
      fetchChats();
    } catch (err) {
      setError('Failed to create chat');
    }
  };

  const handleDeleteChat = async (chatId) => {
    if (!window.confirm('Are you sure you want to delete this chat?')) return;
    try {
      await api.delete(`/chats/${chatId}`);
      fetchChats();
    } catch (err) {
      setError('Failed to delete chat');
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  if (loading) return <div>Loading chats...</div>;

  return (
    <div style={{ padding: '20px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between' }}>
        <h1>Your Chats</h1>
        <button onClick={handleLogout}>Logout</button>
      </div>
      {error && <div style={{ color: 'red' }}>{error}</div>}
      <form onSubmit={handleCreateChat} style={{ marginBottom: '20px' }}>
        <input
          type="text"
          placeholder="New chat name"
          value={newChatName}
          onChange={(e) => setNewChatName(e.target.value)}
          style={{ padding: '8px', width: '300px' }}
        />
        <button type="submit" style={{ padding: '8px 16px', marginLeft: '10px' }}>Create</button>
      </form>
      <ChatList chats={chats} onDelete={handleDeleteChat} />
    </div>
  );
};

export default ChatsPage;