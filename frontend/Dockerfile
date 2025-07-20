# Стадия сборки
FROM node:20-alpine AS builder

WORKDIR /app

# Копируем package-файлы
COPY package.json package-lock.json ./
RUN npm ci

# Копируем весь код
COPY . .

# Удаляем кэш билда, если был
RUN rm -rf .next

# Сборка проекта
RUN npm run build

# Стадия запуска
FROM node:20-alpine AS runner

WORKDIR /app

# Устанавливаем только прод-зависимости
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# Копируем билд и нужные файлы
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/next.config.js ./next.config.js
COPY --from=builder /app/tsconfig.json ./tsconfig.json
COPY --from=builder /app/src ./src

EXPOSE 3000

CMD ["npm", "start"]
