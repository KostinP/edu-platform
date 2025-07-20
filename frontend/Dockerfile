#############################
# простая сборка на сервере #
#############################
# FROM node:20 AS builder
# WORKDIR /app
# COPY . .
# RUN npm install && npm run build

# FROM node:20
# WORKDIR /app
# COPY --from=builder /app .
# EXPOSE 3000
# CMD ["npm", "start"]

######################################
# оптимизированная сборка на сервере #
######################################
# FROM node:20 AS deps
# WORKDIR /app
# COPY package*.json ./
# RUN npm install

# FROM node:20 AS builder
# WORKDIR /app
# COPY --from=deps /app/node_modules ./node_modules
# COPY . .
# RUN npm run build

# FROM node:20
# WORKDIR /app
# COPY --from=builder /app ./
# EXPOSE 3000
# CMD ["npm", "start"]

###########################
# сборка на сервере с кешированием #
###########################
# Сборка 
# FROM node:20 AS deps
# WORKDIR /app
# COPY package.json package-lock.json ./
# RUN npm install

# FROM node:20 AS builder
# WORKDIR /app
# COPY --from=deps /app/node_modules ./node_modules
# COPY . .
# RUN npm run build

# FROM node:20 AS runner
# WORKDIR /app
# ENV NODE_ENV=production

# COPY --from=builder /app/public ./public
# COPY --from=builder /app/.next ./.next
# COPY --from=builder /app/node_modules ./node_modules
# COPY --from=builder /app/package.json ./package.json

# EXPOSE 3000
# CMD ["npm", "start"]

# Сборка на клиенте с отправкой на сервер
# Вариант 1:
# rsync -avz --exclude='.git' --exclude='node_modules/.cache' ./ root@your-vps-ip:/srv/edu-platform/frontend
# Вариант 2:
# scp -r .next public node_modules package.json Dockerfile root@your-vps-ip:/srv/edu-platform/frontend/

FROM node:20
WORKDIR /app

COPY . .

ENV NODE_ENV=production

EXPOSE 3000
CMD ["npm", "start"]
