# Этап сборки
FROM node:18-alpine AS builder

# Установка рабочей директории
WORKDIR /app

# Копирование package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

COPY . .

# Vite сборка
RUN npm run build

# Проверка папки dist
RUN ls -la /app/dist

# Этап продакшн
FROM nginx:alpine

# Устанавливаем рабочую директорию для Nginx
WORKDIR /usr/share/nginx/html

# Копируем собранные файлы из этапа сборки в директорию Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 5173

CMD ["nginx", "-g", "daemon off;"]