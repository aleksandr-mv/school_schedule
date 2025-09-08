# School Schedule - Микросервисная архитектура

Проект представляет собой микросервисную архитектуру для системы управления школьным расписанием с использованием Go, gRPC, Docker и современных инструментов наблюдаемости.

## 🏗️ Архитектура

### Микросервисы
- **IAM Service** - управление пользователями, аутентификация и авторизация
- **Event Bus** - Kafka для асинхронной коммуникации между сервисами
- **Cache Service** - Redis для кэширования данных

### Инфраструктура
- **PostgreSQL** - основная база данных
- **Redis** - кэширование и сессии
- **Kafka** - event streaming
- **Cassandra** - хранение трейсов

### Наблюдаемость
- **Elasticsearch + Kibana** - централизованное логирование
- **Prometheus + Grafana** - метрики и мониторинг
- **Jaeger** - распределенная трассировка
- **OpenTelemetry Collector** - сбор телеметрии

## 🚀 Быстрый старт

### Предварительные требования
- Docker & Docker Compose
- Go 1.21+
- Task (Taskfile)

### Установка и запуск

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/aleksandr-mv/school_schedule.git
   cd school_schedule
   ```

2. **Сгенерируйте переменные окружения:**
   ```bash
   task env:generate
   ```

3. **Запустите всю инфраструктуру:**
   ```bash
   task up-all
   ```

4. **Проверьте статус сервисов:**
   ```bash
   task status
   ```

## 🛠️ Разработка

### Основные команды

```bash
# Генерация protobuf файлов
task proto:generate

# Запуск линтеров
task lint

# Запуск тестов
task test

# Сборка IAM сервиса
task build:iam

# Запуск миграций
task migrate:up
```

### Структура проекта

```
├── iam/                    # IAM микросервис
│   ├── cmd/               # Точки входа
│   ├── internal/          # Внутренняя логика
│   ├── config/            # Конфигурации
│   └── db/               # Миграции БД
├── platform/             # Общие платформенные пакеты
├── shared/               # Общие компоненты
│   ├── proto/            # Protobuf определения
│   └── pkg/              # Сгенерированный код
├── deploy/               # Docker Compose конфигурации
│   ├── compose/          # Сервисы
│   └── env/              # Переменные окружения
└── workflows/            # GitHub Actions
```

## 🌐 Доступные сервисы

После запуска `task up-all` будут доступны:

- **Kibana**: http://localhost:5601
- **Grafana**: http://localhost:3000
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **Kafka UI**: http://localhost:8080
- **Elasticsearch**: http://localhost:9200

## 🔧 Конфигурация

### Переменные окружения
Все переменные настраиваются через систему шаблонов в `deploy/env/`:

- `deploy/env/services/` - модульные конфигурации
- `deploy/env/*.env.template` - шаблоны для генерации
- `deploy/env/generate-env.sh` - скрипт генерации

### Docker Compose
Сервисы организованы по проектам:
- `core` - сетевая инфраструктура
- `eventbus` - Kafka
- `tracing` - Jaeger + Cassandra
- `logs` - Elasticsearch + Kibana
- `metrics` - Prometheus + Grafana
- `iam` - PostgreSQL для IAM
- `cache` - Redis
- `otel` - OpenTelemetry Collector

## 🧪 Тестирование

```bash
# Запуск всех тестов
task test

# Тесты с покрытием
task test:coverage

# Интеграционные тесты
task test:integration
```

## 📊 Мониторинг

### Логи
- Централизованное логирование через Elasticsearch
- Kibana для анализа логов
- Структурированные логи в JSON формате

### Метрики
- Prometheus для сбора метрик
- Grafana дашборды для визуализации
- Бизнес-метрики и системные метрики

### Трассировка
- Jaeger для распределенной трассировки
- Cassandra для хранения трейсов
- OpenTelemetry для стандартизации

## 🔒 Безопасность

- JWT токены для аутентификации
- HTTPS/TLS для всех внешних соединений
- Валидация входных данных
- Логирование безопасности

## 📈 CI/CD

GitHub Actions workflows:
- `ci.yml` - основной CI pipeline
- `lint-reusable.yml` - проверка кода
- `test-reusable.yml` - запуск тестов
- `test-coverage-reusable.yml` - покрытие тестами
- `test-integration-reusable.yml` - интеграционные тесты

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия

Этот проект создан в образовательных целях.

## 👨‍💻 Автор

**Aleksandr Mv** - [GitHub](https://github.com/aleksandr-mv)

---

*Для получения дополнительной информации см. документацию в соответствующих директориях проекта.*
