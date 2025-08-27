package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func generateMarkdownForDirectory(alerts []*AlertInfo, directory string, customTitle string) string {
	var sb strings.Builder

	// Если кастомный заголовок не указан, формируем автоматически
	title := customTitle
	if title == "" {
		group := extractGroupFromDirectory(directory)
		title = fmt.Sprintf("%s Алерты", group)
	}

	sb.WriteString(fmt.Sprintf("# %s\n\n", title))

	for i, alert := range alerts {
		sb.WriteString(fmt.Sprintf("### %d. %s\n", i+1, alert.Title))
		sb.WriteString(fmt.Sprintf("**Описание:** %s  \n", alert.Description))
		sb.WriteString(fmt.Sprintf("**Файл:** [`%s`](%s)\n\n", alert.RelativePath, alert.RelativePath))
	}

	return sb.String()
}

func extractGroupFromDirectory(directory string) string {
	// Получаем базовое имя директории
	baseName := filepath.Base(directory)

	// Разделяем по дефисам и берем последнее слово
	parts := strings.Split(baseName, "-")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		// Делаем первую букву заглавной
		return strings.Title(lastPart)
	}

	return "Unknown"
}
