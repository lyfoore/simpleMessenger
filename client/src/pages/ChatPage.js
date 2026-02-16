import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';
import MessageList from '../components/MessageList';
import MessageInput from '../components/MessageInput';
import WebSocketService from '../services/websocket';

const ChatPage = () => {
  const { chatId } = useParams();
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { token, user } = useAuth();
  const navigate = useNavigate();
  const ws = useRef(null);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await api.get(`/chats/${chatId}/messages?limit=50`);
        setMessages(response.data.messages);
      } catch (err) {
        setError('Failed to load messages');
      } finally {
        setLoading(false);
      }
    };
    fetchMessages();

    if (token) {
      ws.current = new WebSocketService(chatId, token, (newMessage) => {
        setMessages(prev => [...prev, newMessage]);
      });
      ws.current.connect();
    }

    return () => {
      if (ws.current) {
        ws.current.disconnect();
      }
    };
  }, [chatId, token]);

  const handleSendMessage = async (text) => {
    if (!text.trim()) return;
    try {
      ws.current.sendMessage(chatId, text);
    } catch (err) {
      setError('Failed to send message');
    }
  };

  const handleDeleteMessage = async (messageId) => {
    if (!window.confirm('Delete this message?')) return;
    try {
      await api.delete(`/messages/${messageId}`);
      setMessages(prev => prev.filter(m => m.ID !== messageId));
    } catch (err) {
      setError('Failed to delete message');
    }
  };

  if (loading) return <div>Loading messages...</div>;

  return (
    <div style={{ padding: '20px', height: '100vh', display: 'flex', flexDirection: 'column' }}>
      <div style={{ marginBottom: '10px' }}>
        <button onClick={() => navigate('/chats')}>Back to Chats</button>
      </div>
      {error && <div style={{ color: 'red' }}>{error}</div>}
      <MessageList messages={messages} currentUserId={user?.id} onDelete={handleDeleteMessage} />
      <MessageInput onSend={handleSendMessage} />
    </div>
  );
};

export default ChatPage;