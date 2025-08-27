package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func createZipArchive(outputZip string, directories []*DirectoryInfo) error {
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return fmt.Errorf("ошибка создания zip файла: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, dir := range directories {
		for _, alert := range dir.Alerts {
			if err := addFileToZip(zipWriter, alert.FilePath, filepath.Join(dir.Name, alert.RelativePath)); err != nil {
				log.Printf("Предупреждение: не удалось добавить файл %s: %v", alert.FilePath, err)
			}
		}

		alertsMarkdownPath := filepath.Join(dir.Path, baseDirectoryDocFileName)
		if _, err := os.Stat(alertsMarkdownPath); err == nil {
			if err := addFileToZip(zipWriter, alertsMarkdownPath, filepath.Join(dir.Name, baseDirectoryDocFileName)); err != nil {
				log.Printf("Предупреждение: не удалось добавить файл %s: %v", alertsMarkdownPath, err)
			}
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(zipPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
