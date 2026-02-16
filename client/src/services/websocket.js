class WebSocketService {
  constructor(chatId, token, onMessage) {
    this.chatId = chatId;
    this.token = token;
    this.onMessage = onMessage;
    this.socket = null;
  }

  connect() {
    console.log('WebSocket token:', this.token);
    const baseUrl = window.__ENV_WS_URL__ || process.env.REACT_APP_WS_URL || 'ws://localhost:8080';
    const wsUrl = `${baseUrl}/api/ws?token=${encodeURIComponent(this.token)}`;
    console.log('Connecting to:', wsUrl);
    this.socket = new WebSocket(wsUrl);

    this.socket.onopen = () => {
      console.log('WebSocket connected');
    };

    this.socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (message.chatId === parseInt(this.chatId)) {
          this.onMessage(message);
        }
      } catch (e) {
        console.error('Failed to parse message', e);
      }
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error', error);
    };

    this.socket.onclose = () => {
      console.log('WebSocket disconnected');
    };
  }

  sendMessage(chatId, text) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      const message = {
        chatId: parseInt(chatId),
        text: text
      };
      this.socket.send(JSON.stringify(message));
    } else {
      console.error('WebSocket not connected');
    }
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
    }
  }
}

export default WebSocketService;