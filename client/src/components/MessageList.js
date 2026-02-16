import React, { useRef, useEffect } from 'react';

const MessageList = ({ messages, currentUserId, onDelete }) => {
  const messagesEndRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  return (
    <div style={{ flex: 1, overflowY: 'auto', padding: '10px', border: '1px solid #ccc', marginBottom: '10px' }}>
      {messages.map(message => (
        <div key={message.ID} style={{ marginBottom: '10px', textAlign: message.userId === currentUserId ? 'right' : 'left' }}>
          <div style={{ display: 'inline-block', background: message.userId === currentUserId ? '#dcf8c6' : '#fff', padding: '8px 12px', borderRadius: '8px', maxWidth: '70%', boxShadow: '0 1px 1px rgba(0,0,0,0.1)' }}>
            <div><strong>User {message.userId}</strong></div>
            <div>{message.text}</div>
            <div style={{ fontSize: '0.7em', color: '#666', marginTop: '4px' }}>
              {new Date(message.CreatedAt).toLocaleTimeString()}
              {message.userId === currentUserId && (
                <button onClick={() => onDelete(message.ID)} style={{ marginLeft: '10px', background: 'none', border: 'none', color: 'red', cursor: 'pointer' }}>🗑️</button>
              )}
            </div>
          </div>
        </div>
      ))}
      <div ref={messagesEndRef} />
    </div>
  );
};

export default MessageList;