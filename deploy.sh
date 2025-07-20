# Скрипт автоматизирует сборку бэкенда и фронтенда на сервере
# chmod +x deploy.sh
# ./deploy.sh

#!/bin/bash

# На сервере пересобираем образ и перезапускаем контейнеры
ssh root@213.139.208.67 << EOF
  cd /srv/edu-platform
  git pull
  docker compose down
  docker compose build --no-cache
  docker compose up -d
EOF
