# 🛠️ Стадия сборки
FROM node:20-alpine AS builder

# Установим рабочую директорию
WORKDIR /app

# Установим зависимости (сначала только package.json — для кеша)
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

# Копируем весь исходный код
COPY frontend .

# Удаляем предыдущий кэш билда (если был)
RUN rm -rf .next

# Сборка приложения
RUN npm run build

# 🔒 Стадия запуска
FROM node:20-alpine AS runner

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем только продакшен-зависимости
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --omit=dev

# Копируем скомпилированные файлы из builder
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/next.config.js ./next.config.js
COPY --from=builder /app/tsconfig.json ./tsconfig.json
COPY --from=builder /app/src ./src

# Используем порт 3000
EXPOSE 3000

# Запуск Next.js в режиме продакшена
CMD ["npm", "start"]
