# School Schedule - Kubernetes Deployment

Продакшен-ready Kubernetes манифесты для School Schedule приложения.

## Архитектура

### Namespaces
- `school-schedule-system` - системные компоненты
- `school-schedule-infrastructure` - инфраструктура (PostgreSQL, Redis, Kafka)
- `school-schedule-services` - микросервисы (IAM, RBAC)
- `school-schedule-monitoring` - мониторинг (OpenTelemetry)
- `school-schedule-ingress` - Ingress и сетевые политики

### Компоненты

#### Инфраструктура
- **PostgreSQL** - StatefulSet с персистентными томами
- **Redis Cluster** - 6 нод с автоматической инициализацией кластера
- **Kafka** - с Zookeeper для координации

#### Сервисы
- **IAM Service** - 3 реплики с health checks
- **RBAC Service** - 3 реплики с health checks

#### Мониторинг
- **OpenTelemetry Collector** - 2 реплики для высокой доступности
- **Jaeger** - трейсинг
- **Prometheus** - метрики
- **Elasticsearch** - логи

#### Безопасность
- **Network Policies** - изоляция трафика между namespace
- **Secrets** - безопасное хранение паролей и ключей
- **ConfigMaps** - конфигурация без секретов

## Быстрый старт

### Предварительные требования
- Kubernetes кластер (minikube, kind, или production)
- kubectl настроен для доступа к кластеру
- Docker images должны быть собраны и доступны

### Деплой
```bash
# Деплой всего приложения
./deploy.sh

# Или пошагово
kubectl apply -f namespaces/namespaces.yaml
kubectl apply -f infrastructure/
kubectl apply -f monitoring/
kubectl apply -f services/
kubectl apply -f ingress/
```

### Проверка статуса
```bash
# Все поды
kubectl get pods -A

# Логи сервисов
kubectl logs -f deployment/iam-service -n school-schedule-services
kubectl logs -f deployment/rbac-service -n school-schedule-services

# Порты для локального тестирования
kubectl port-forward service/iam-service 8080:8080 -n school-schedule-services
kubectl port-forward service/rbac-service 8080:8080 -n school-schedule-services
```

### Очистка
```bash
./cleanup.sh
```

## Продакшен особенности

### Масштабируемость
- **Horizontal Pod Autoscaler** - автоматическое масштабирование по CPU/Memory
- **Vertical Pod Autoscaler** - оптимизация ресурсов
- **Cluster Autoscaler** - автоматическое добавление нод

### Надежность
- **Health Checks** - liveness и readiness probes
- **Resource Limits** - предотвращение OOM
- **Anti-affinity** - распределение подов по нодам
- **Pod Disruption Budgets** - гарантия доступности

### Безопасность
- **Network Policies** - микросегментация сети
- **RBAC** - контроль доступа
- **Secrets Management** - безопасное хранение секретов
- **Image Security** - сканирование уязвимостей

### Мониторинг
- **OpenTelemetry** - единая телеметрия
- **Distributed Tracing** - трейсинг запросов
- **Metrics Collection** - сбор метрик
- **Log Aggregation** - централизованные логи

## Конфигурация

### Environment Variables
- `CONFIG_PATH` - путь к конфигурационному файлу
- `LOG_LEVEL` - уровень логирования
- `OTEL_ENDPOINT` - endpoint для OpenTelemetry

### Secrets
- `postgres-secret` - пароли для PostgreSQL
- `redis-secret` - пароли для Redis
- `app-secrets` - JWT и encryption ключи
- `otel-secrets` - токены для мониторинга

## Troubleshooting

### Проверка статуса
```bash
kubectl get all -n school-schedule-services
kubectl describe pod <pod-name> -n school-schedule-services
```

### Логи
```bash
kubectl logs <pod-name> -n school-schedule-services
kubectl logs <pod-name> -n school-schedule-services --previous
```

### События
```bash
kubectl get events -n school-schedule-services --sort-by='.lastTimestamp'
```

### Отладка сети
```bash
kubectl exec -it <pod-name> -n school-schedule-services -- nslookup <service-name>
kubectl exec -it <pod-name> -n school-schedule-services -- curl <service-url>
```
