package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ruleTemplates = map[string]string{
	"go":           "https://raw.githubusercontent.com/PatrickJS/awesome-cursorrules/main/rules/go/cursorrules",
	"vue-frontend": "https://raw.githubusercontent.com/PatrickJS/awesome-cursorrules/main/rules/vue/cursorrules",
}

var specialCommands = map[string]string{
	"supabase": "npx shadcn@latest add https://supabase.com/ui/r/ai-editor-rules.json",
}

func main() {
	var (
		projectType = flag.String("type", "go", "Тип проекта (go, vue, supabase)")
		customURL   = flag.String("url", "", "Кастомный URL для .cursorrules файла")
		directory   = flag.String("dir", ".", "Директория для установки .cursorrules")
		list        = flag.Bool("list", false, "Показать доступные типы проектов")
	)
	flag.Parse()

	if *list {
		fmt.Println("Доступные типы проектов:")
		for projectType := range ruleTemplates {
			fmt.Printf("  - %s\n", projectType)
		}
		fmt.Println("Специальные команды:")
		for projectType := range specialCommands {
			fmt.Printf("  - %s (выполняет команду)\n", projectType)
		}
		return
	}

	// Проверяем, что директория существует
	if info, err := os.Stat(*directory); err != nil || !info.IsDir() {
		log.Fatalf("Директория %s не существует", *directory)
	}

	// Проверяем, есть ли специальная команда для этого типа проекта
	if command, exists := specialCommands[strings.ToLower(*projectType)]; exists {
		fmt.Printf("🚀 Выполняем специальную команду для %s...\n", *projectType)
		if err := executeCommand(command, *directory); err != nil {
			log.Fatalf("Ошибка выполнения команды: %v", err)
		}
		fmt.Printf("✅ Команда успешно выполнена для проекта %s\n", *projectType)
		return
	}

	// Определяем URL для скачивания
	var rulesURL string
	if *customURL != "" {
		rulesURL = *customURL
	} else {
		var exists bool
		rulesURL, exists = ruleTemplates[strings.ToLower(*projectType)]
		if !exists {
			log.Fatalf("Неизвестный тип проекта: %s. Используйте -list для просмотра доступных типов", *projectType)
		}
	}

	// Путь к файлу .cursorrules
	rulesPath := filepath.Join(*directory, ".cursorrules")

	// Проверяем, существует ли уже файл
	if _, err := os.Stat(rulesPath); err == nil {
		fmt.Printf("Файл .cursorrules уже существует в %s\n", *directory)
		fmt.Print("Перезаписать? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Установка отменена")
			return
		}
	}

	// Скачиваем и устанавливаем .cursorrules
	if err := downloadAndInstallRules(rulesURL, rulesPath); err != nil {
		log.Fatalf("Ошибка установки cursor rules: %v", err)
	}

	fmt.Printf("✅ Cursor rules успешно установлены в %s\n", rulesPath)
	fmt.Printf("📁 Тип проекта: %s\n", *projectType)
	if *customURL != "" {
		fmt.Printf("🔗 URL: %s\n", *customURL)
	}
}

func downloadAndInstallRules(url, filePath string) error {
	// Скачиваем файл
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка скачивания: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка HTTP: %s", resp.Status)
	}

	// Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	// Копируем содержимое
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}

func executeCommand(command, directory string) error {
	// Разбиваем команду на части
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("пустая команда")
	}

	// Создаем команду
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Выполняем команду
	return cmd.Run()
}
