# skillsrock-tasks-service

## Запуск

```shell
USER_ID="$(id -u)" docker compose up
```

Запустятся следующие контейнеры:

- postgres:5432 (admin/admin)
- redis:6379
- prometheus:9090
- grafana:3000 (admin/admin)
- app:8080

Дашборд - <http://localhost:3000/d/ypFZFgvmz>

## Дополнительно

- [Примеры API-запросов](/api/api-samples.md)
- [Миграции базы данных](/db//migrations/)
- [Пример файла JSON для импорта/экспорта задач](/api/tasks-export.json)
- [Swagger-документацию для API](/api/openapi.yml)
- [Dev окружение](/flake.nix)
