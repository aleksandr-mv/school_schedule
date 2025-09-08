#!/bin/bash

# Директория со скриптом
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$SCRIPT_DIR"
COMPOSE_DIR="$SCRIPT_DIR/../compose"

# Используем переданную переменную ENV_SUBST или системный envsubst
if [ -z "$ENV_SUBST" ]; then
  if ! command -v envsubst &> /dev/null; then
    echo "❌ Ошибка: envsubst не найден в системе и не передан через ENV_SUBST!"
    echo "Запустите скрипт через task env:generate"
    exit 1
  fi
  ENV_SUBST=envsubst
fi

# Загружаем переменные из модульной структуры
if [ -f "$SCRIPT_DIR/.env.template" ]; then
  echo "🔄 Загружаем переменные из модульной структуры..."
  # Переходим в директорию env для корректной работы relative путей
  cd "$SCRIPT_DIR"
  # Экспортируем все переменные из модульных файлов
  set -a
  # Загружаем все services/*.env файлы
  for env_file in services/*.env; do
    if [ -f "$env_file" ]; then
      source "$env_file"
    fi
  done
  set +a
  cd - > /dev/null
  echo "✅ Переменные загружены из модульных конфигов"
elif [ -f "$SCRIPT_DIR/.env" ]; then
  echo "🔄 Загружаем переменные из .env файла..."
  set -a
  source "$SCRIPT_DIR/.env"
  set +a
else
  echo "❌ Ошибка: Ни .env.template, ни .env файл не найдены!"
  exit 1
fi

# Функция для обработки шаблона и создания .env файла
process_template() {
  local service=$1
  local template="$TEMPLATE_DIR/${service}.env.template"
  local output="$COMPOSE_DIR/${service}/.env"
  
  echo "Обработка шаблона для сервиса $service..."
  
  if [ ! -f "$template" ]; then
    echo "⚠️ Шаблон $template не найден, пропускаем..."
    return 0
  fi
  
  # Создаем директорию, если она еще не существует
  mkdir -p "$(dirname "$output")"
  
  # Используем envsubst для замены переменных в шаблоне
  $ENV_SUBST < "$template" > "$output"
  
  echo "✅ Создан файл $output"
}

# Определяем список сервисов из переменной окружения
if [ -z "$SERVICES" ]; then
  echo "⚠️ Переменная SERVICES не задана. Нет сервисов для обработки."
  exit 0
fi

# Разделяем список сервисов по запятой
IFS=',' read -ra services <<< "$SERVICES"
echo "🔍 Обрабатываем сервисы: ${services[*]}"

# Обрабатываем шаблоны для всех указанных сервисов
success_count=0
skip_count=0
for service in "${services[@]}"; do
  process_template "$service"
  if [ -f "$TEMPLATE_DIR/${service}.env.template" ]; then
    ((success_count++))
  else
    ((skip_count++))
  fi
done

if [ $success_count -eq 0 ]; then
  echo "⚠️ Ни один .env файл не создан. Проверьте список сервисов и наличие шаблонов."
else
  echo "🎉 Генерация завершена: $success_count файлов создано, $skip_count шаблонов пропущено"
fi 