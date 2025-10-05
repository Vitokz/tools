#!/bin/bash

# Vitokz Tools Installer
echo "🔧 Устанавливаем Vitokz Tools..."

# Установка repo-init
echo "📦 Устанавливаем repo-init..."
go install github.com/Vitokz/tools/scripts/repo-init@latest

# Установка export-grafan-based-alerts
echo "📊 Устанавливаем export-grafan-based-alerts..."
go install github.com/Vitokz/tools/scripts/export-grafan-based-alerts@latest

echo "✅ Установка завершена!"
echo "Инструменты доступны в: $(go env GOPATH)/bin"
echo ""
echo "Использование:"
echo "  repo-init -type vue-frontend"
echo "  export-grafan-based-alerts -help"
