# Скрипт автоматизирует сборку фронтенда и передачу с клиента на сервер
# chmod +x deploy.sh
# ./deploy.sh

#!/bin/bash

# # Установка зависимостей
# npm install

# # Сборка фронтенда
# npm run build

# # Отправляем исходники на сервер
# rsync -avz --delete . root@213.139.208.67:/srv/edu-platform/frontend/

# На сервере пересобираем образ и перезапускаем контейнер
ssh root@213.139.208.67 << EOF
  cd /srv/edu-platform
  git pull
  docker-compose build frontend
  docker-compose up -d frontend
EOF
