import React from 'react';
import { Link } from 'react-router-dom';

const ChatList = ({ chats, onDelete }) => {
  return (
    <ul style={{ listStyle: 'none', padding: 0 }}>
      {chats.map(chat => (
        <li key={chat.ID} style={{ margin: '10px 0', padding: '10px', border: '1px solid #ccc', borderRadius: '4px' }}>
          <Link to={`/chats/${chat.ID}`} style={{ textDecoration: 'none', color: 'inherit', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div>
              <strong>{chat.name}</strong>
              <div style={{ fontSize: '0.8em', color: '#666' }}>
                Last message: {chat.LastMessageAt ? new Date(chat.LastMessageAt).toLocaleString() : 'No messages'}
              </div>
            </div>
          </Link>
          <button onClick={() => onDelete(chat.ID)} style={{ padding: '4px 8px', background: 'red', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}>
            Delete
          </button>
        </li>
      ))}
    </ul>
  );
};

export default ChatList;