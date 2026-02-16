# Stage 1: build
FROM node:18-alpine AS builder

WORKDIR /app

# Копируем package.json и package-lock.json (если есть)
COPY package*.json ./

# Устанавливаем зависимости (используем npm install, так как lock-файл может отсутствовать)
RUN npm install

# Копируем исходный код и собираем приложение
COPY . .
RUN npm run build

# Stage 2: serve with nginx
FROM nginx:alpine

# Копируем собранные файлы
COPY --from=builder /app/build /usr/share/nginx/html

# Копируем скрипт для подстановки переменных окружения
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80

# ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]