package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	yamlPattern              = "*.yml"
	baseDirectoryDocFileName = "README.md"
	defaultZipFileName       = "alerts.zip"
)

type Alert struct {
	Title       string                 `yaml:"title"`
	Annotations map[string]interface{} `yaml:"annotations"`
}

type AlertInfo struct {
	Title        string
	Description  string
	FilePath     string
	RelativePath string
}

type DirectoryInfo struct {
	Name   string
	Path   string
	Alerts []*AlertInfo
}

func main() {
	var (
		directory = flag.String("dir", ".", "Корневая директория с YAML файлами алертов")
		dirs      = flag.String("dirs", "", "Список директорий через запятую (например: feeds-team-http,feeds-team-kafka)")
		recursive = flag.Bool("r", false, "Рекурсивно обходить поддиректории")
		zipOutput = flag.Bool("zip", false, "Создать ZIP архив с алертами и документацией")
	)
	flag.Parse()

	if info, err := os.Stat(*directory); err != nil || !info.IsDir() {
		log.Fatalf("Ошибка: Директория %s не существует", *directory)
	}

	fmt.Printf("Сканирование директории: %s\n", *directory)

	var (
		directories []*DirectoryInfo
		err         error
	)
	if *dirs != "" {
		directories, err = specificScan(directory, dirs)
	} else if *recursive {
		directories, err = recursiveScan(directory)
	} else {
		directories, err = soloRepoScan(directory)
	}

	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	if *zipOutput {
		if err := createZipArchive(defaultZipFileName, directories); err != nil {
			log.Fatalf("Ошибка создания ZIP архива: %v", err)
		}
		fmt.Printf("ZIP архив создан: %s\n", defaultZipFileName)
	}
}

func specificScan(directory, dirs *string) ([]*DirectoryInfo, error) {
	targetDirs := strings.Split(*dirs, ",")
	fmt.Printf("Обработка директорий: %v\n", targetDirs)

	directories, err := scanSpecificDirectories(*directory, targetDirs)
	if err != nil {
		log.Fatalf("Ошибка сканирования указанных директорий: %v", err)
	}

	if len(directories) == 0 {
		log.Fatal("Не найдено ни одного алерта в указанных директориях")
	}

	totalAlerts := 0
	for _, dir := range directories {
		totalAlerts += len(dir.Alerts)
	}
	fmt.Printf("Найдено директорий с алертами: %d (всего алертов: %d)\n", len(directories), totalAlerts)

	// Генерируем документацию для каждой директории отдельно
	for _, dir := range directories {
		markdown := generateMarkdownForDirectory(dir.Alerts, dir.Path, "")
		outputFile := filepath.Join(dir.Path, baseDirectoryDocFileName)

		if err := ioutil.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
			log.Printf("Ошибка записи файла %s: %v", outputFile, err)
		} else {
			fmt.Printf("Документация сохранена: %s\n", outputFile)
		}
	}

	return directories, err
}

func recursiveScan(directory *string) ([]*DirectoryInfo, error) {
	directories, err := scanDirectoriesRecursively(*directory)
	if err != nil {
		log.Fatalf("Ошибка рекурсивного сканирования: %v", err)
	}

	if len(directories) == 0 {
		log.Fatal("Не найдено ни одного алерта в поддиректориях")
	}

	totalAlerts := 0
	for _, dir := range directories {
		totalAlerts += len(dir.Alerts)
	}
	fmt.Printf("Найдено директорий с алертами: %d (всего алертов: %d)\n", len(directories), totalAlerts)

	// Генерируем документацию для каждой директории отдельно
	for _, dir := range directories {
		markdown := generateMarkdownForDirectory(dir.Alerts, dir.Path, "")
		outputFile := filepath.Join(dir.Path, baseDirectoryDocFileName)

		if err := ioutil.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
			log.Printf("Ошибка записи файла %s: %v", outputFile, err)
		} else {
			fmt.Printf("Документация сохранена: %s\n", outputFile)
		}
	}

	return directories, err
}

func soloRepoScan(directory *string) ([]*DirectoryInfo, error) {
	// Обычное сканирование одной директории
	alerts, err := scanDirectory(*directory)
	if err != nil {
		log.Fatalf("Ошибка сканирования: %v", err)
	}

	if len(alerts) == 0 {
		log.Fatal("Не найдено ни одного алерта")
	}

	fmt.Printf("Найдено алертов: %d\n", len(alerts))

	markdown := generateMarkdownForDirectory(alerts, *directory, "")

	outputFile := filepath.Join(*directory, baseDirectoryDocFileName)
	if err := ioutil.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
		log.Fatalf("Ошибка записи файла %s: %v", outputFile, err)
	}

	fmt.Printf("Документация сохранена: %s\n", outputFile)

	return []*DirectoryInfo{
		{
			Name:   filepath.Base(*directory),
			Path:   *directory,
			Alerts: alerts,
		},
	}, err
}
