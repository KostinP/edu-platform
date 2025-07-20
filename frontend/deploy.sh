# Скрипт автоматизирует сборку фронтенда на сервере
# chmod +x deploy.sh
# ./deploy.sh

#!/bin/bash

# На сервере пересобираем образ и перезапускаем контейнер
ssh root@213.139.208.67 << EOF
  cd /srv/edu-platform
  git pull
  docker compose build frontend
  docker compose up -d frontend
EOF
