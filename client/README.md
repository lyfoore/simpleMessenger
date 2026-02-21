# Messenger Frontend (React)

Современный React-фронтенд для мессенджера на Go.

## Технологии

- **React 18** с TypeScript
- **Vite** - быстрый сборщик
- **Tailwind CSS** - стилизация
- **React Router** - навигация
- **WebSocket** - realtime-сообщения

## Структура проекта

```
client2/
├── src/
│   ├── components/       # UI-компоненты
│   │   ├── ChatList.tsx
│   │   ├── ChatWindow.tsx
│   │   └── NewChatModal.tsx
│   ├── pages/           # Страницы
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   └── Chats.tsx
│   ├── contexts/        # React Context
│   │   └── AuthContext.tsx
│   ├── hooks/           # Custom hooks
│   ├── api.ts           # API клиент
│   ├── types.ts         # TypeScript типы
│   ├── main.tsx         # Точка входа
│   └── index.css        # Глобальные стили
├── Dockerfile
├── nginx.conf
└── package.json
```

## Запуск через Docker Compose

1. Скопируйте `.env.example` в `.env`:
   ```bash
   cp .env.example .env
   ```

2. Запустите все сервисы:
   ```bash
   docker-compose up --build
   ```

3. Откройте браузер:
   - Фронтенд: http://localhost:3000
   - API: http://localhost:8080

## Функционал

- ✅ Регистрация/авторизация по username
- ✅ Создание чатов с другими пользователями
- ✅ Отправка/получение сообщений (REST + WebSocket)
- ✅ Удаление сообщений и чатов
- ✅ Realtime-обновления через WebSocket
- ✅ Адаптивный дизайн
- ✅ Тёмная тема

## API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/auth/register` | Регистрация |
| POST | `/api/auth/login` | Авторизация |
| GET | `/api/chats` | Список чатов |
| POST | `/api/chats` | Создать чат |
| DELETE | `/api/chats/:chatId` | Удалить чат |
| GET | `/api/chats/:chatId/messages` | Получить сообщения |
| POST | `/api/chats/:chatId/messages` | Отправить сообщение |
| DELETE | `/api/messages/:messageId` | Удалить сообщение |
| GET | `/api/ws` | WebSocket подключение |
