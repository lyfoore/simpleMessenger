#!/bin/sh

# Заменяем плейсхолдеры вида __ENV_XXX__ на значения переменных окружения
# Например, в собранных JS-файлах ищем строку "__ENV_API_URL__" и заменяем на реальный URL
find /usr/share/nginx/html -type f -name "*.js" | while read file; do
    if [ -n "$REACT_APP_API_URL" ]; then
        sed -i "s|__ENV_API_URL__|$REACT_APP_API_URL|g" "$file"
    fi
    if [ -n "$REACT_APP_WS_URL" ]; then
        sed -i "s|__ENV_WS_URL__|$REACT_APP_WS_URL|g" "$file"
    fi
done

# Запускаем nginx
exec "$@"