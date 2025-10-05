# Cursor Rules Installer

CLI инструмент для автоматической установки `.cursorrules` файлов в новые проекты.

## Что делает

Скачивает и устанавливает `.cursorrules` файлы из популярных шаблонов для разных типов проектов. Cursor Rules помогают настроить поведение AI-ассистента Cursor для конкретного типа проекта.

## Установка

```bash
# Сборка
cd scripts/cursor-rules-installer
go build -o cursor-rules-installer

# Установка в систему (опционально)
sudo mv cursor-rules-installer /usr/local/bin/
```

## Использование

### Основные команды

```bash
# Установить правила для Go проекта (по умолчанию)
cursor-rules-installer

# Установить правила для конкретного типа проекта
cursor-rules-installer -type javascript
cursor-rules-installer -type python
cursor-rules-installer -type react

# Установить в конкретную директорию
cursor-rules-installer -type typescript -dir /path/to/project

# Использовать кастомный URL
cursor-rules-installer -url https://raw.githubusercontent.com/user/repo/main/.cursorrules

# Показать доступные типы проектов
cursor-rules-installer -list
```

### Параметры

- `-type` - Тип проекта (go, javascript, typescript, python, react, vue, rust)
- `-dir` - Директория для установки .cursorrules (по умолчанию: текущая)
- `-url` - Кастомный URL для скачивания правил
- `-list` - Показать доступные типы проектов

## Примеры использования

```bash
# Инициализация нового Go проекта
mkdir my-go-project && cd my-go-project
go mod init my-go-project
cursor-rules-installer -type go

# Настройка React проекта
npx create-react-app my-app && cd my-app
cursor-rules-installer -type react

# Установка правил Nuxt из cursor.directory
npx nuxi@latest init my-nuxt-app && cd my-nuxt-app
cursor-rules-installer -type nuxt

# Настройка Supabase проекта (выполняет специальную команду)
cursor-rules-installer -type supabase

# Использование кастомных правил
cursor-rules-installer -url https://gist.githubusercontent.com/user/id/raw/custom-rules
```

## Что получаем на выходе

### Обычные правила (.cursorrules)
В указанной директории создается файл `.cursorrules` с настройками для AI-ассистента Cursor, оптимизированными для выбранного типа проекта.

### Правила из cursor.directory (.cursor/rules/*.mdc)
Создается директория `.cursor/rules/` с файлом правил в формате `.mdc`, который автоматически подключается в Cursor.

### Специальные команды
Выполняются соответствующие команды для настройки проекта (например, установка компонентов через shadcn).

Файлы содержат инструкции для:
- Стиля кодирования
- Лучших практик
- Специфичных для языка/фреймворка рекомендаций
- Структуры проекта
