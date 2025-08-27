package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

func parseAlertFile(filePath string) (*AlertInfo, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %s: %v", filePath, err)
	}

	var alert Alert
	if err := yaml.Unmarshal(data, &alert); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML %s: %v", filePath, err)
	}

	// Извлекаем описание из annotations
	description := "Нет описания"
	if alert.Annotations != nil {
		if desc, ok := alert.Annotations["description"].(string); ok {
			description = desc
		} else if summary, ok := alert.Annotations["summary"].(string); ok {
			description = summary
		}
	}

	return &AlertInfo{
		Title:        alert.Title,
		Description:  description,
		FilePath:     filePath,
		RelativePath: filepath.Base(filePath),
	}, nil
}

func scanDirectory(directory string) ([]*AlertInfo, error) {
	var alerts []*AlertInfo

	files, err := filepath.Glob(filepath.Join(directory, yamlPattern))
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска файлов: %v", err)
	}

	for _, filePath := range files {
		alertInfo, err := parseAlertFile(filePath)
		if err != nil {
			log.Printf("Предупреждение: %v", err)
			continue
		}
		alerts = append(alerts, alertInfo)
	}

	// Сортируем по названию
	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].Title < alerts[j].Title
	})

	return alerts, nil
}

func scanSpecificDirectories(rootDir string, targetDirs []string) ([]*DirectoryInfo, error) {
	var directories []*DirectoryInfo

	for _, targetDir := range targetDirs {
		// Строим полный путь к директории
		fullPath := filepath.Join(rootDir, strings.TrimSpace(targetDir))

		// Проверяем существование директории
		if info, err := os.Stat(fullPath); err != nil || !info.IsDir() {
			log.Printf("Предупреждение: директория %s не найдена или не является директорией", fullPath)
			continue
		}

		// Сканируем алерты в директории
		alerts, err := scanDirectory(fullPath)
		if err != nil {
			log.Printf("Предупреждение при сканировании %s: %v", fullPath, err)
			continue
		}

		// Если есть алерты, добавляем директорию
		if len(alerts) > 0 {
			directories = append(directories, &DirectoryInfo{
				Name:   targetDir,
				Path:   fullPath,
				Alerts: alerts,
			})
		} else {
			log.Printf("В директории %s не найдено алертов", fullPath)
		}
	}

	return directories, nil
}

func scanDirectoriesRecursively(rootDir string) ([]*DirectoryInfo, error) {
	var directories []*DirectoryInfo

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем файлы, обрабатываем только директории
		if !info.IsDir() {
			return nil
		}

		// Сканируем алерты в текущей директории
		alerts, err := scanDirectory(path)
		if err != nil {
			log.Printf("Предупреждение при сканировании %s: %v", path, err)
			return nil
		}

		// Если есть алерты, добавляем директорию
		if len(alerts) > 0 {
			relPath, _ := filepath.Rel(rootDir, path)
			if relPath == "." {
				relPath = filepath.Base(rootDir)
			}

			directories = append(directories, &DirectoryInfo{
				Name:   relPath,
				Path:   path,
				Alerts: alerts,
			})
		}

		return nil
	})

	return directories, err
}
