# GitHub Workflows

## [Deploy to Coolify](./simple_deploy_to_coolify.yaml)

Этот workflow выполняет деплой в `Coolify`. Он запускается при каждом пуше в ветку ```main``` или может быть запущен вручную. Workflow просто делает вызов webhook URL Coolify для инициирования процесса деплоя.

В самом Coolify должен быть подключен репозиторий github из которого будет выгружаться ```Dockerfile``` или ```docker-compose.yaml```

Для работы требуется:
- `COOLIFY_WEBHOOK_URL` - URL для вызова webhook Coolify
- `COOLIFY_WEBHOOK_SECRET` - секретный ключ для авторизации в Coolify


## [Deploy to Coolify с миграциями базы данных](./deploy_to_coolify_with_migrations.yaml)

Этот workflow выполняет деплой в `Coolify` с предварительным запуском миграций базы данных. Он запускается при каждом пуше в ветку ```main``` или может быть запущен вручную.

Последовательность действий:
1. Выгружает код из репозитория
2. Устанавливает Go
3. Устанавливает инструмент миграций Goose
4. Запускает миграции базы данных из директории `./migrations` с использованием переменной окружения `DATABASE_URL`
5. Вызывает webhook URL Coolify для инициирования процесса деплоя

Для работы этого workflow требуется настроить следующие секреты в GitHub Actions:
- `DATABASE_URL` - строка подключения к базе данных PostgreSQL
- `COOLIFY_WEBHOOK_URL` - URL для вызова webhook Coolify
- `COOLIFY_WEBHOOK_SECRET` - секретный ключ для авторизации в Coolify

