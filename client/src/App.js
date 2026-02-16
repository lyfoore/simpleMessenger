import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import PrivateRoute from './components/PrivateRoute';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ChatsPage from './pages/ChatsPage';
import ChatPage from './pages/ChatPage';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/chats" element={
            <PrivateRoute>
              <ChatsPage />
            </PrivateRoute>
          } />
          <Route path="/chats/:chatId" element={
            <PrivateRoute>
              <ChatPage />
            </PrivateRoute>
          } />
          <Route path="/" element={<Navigate to="/chats" />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;