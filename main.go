package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	defaultUrlFlag := "urls.txt"
	defaultDirFlag := "htmls"

	urlPtr := flag.String("url", defaultUrlFlag, "путь к файлу с url")
	dirPtr := flag.String("dir", defaultDirFlag, "путь к папку с html страницами")

	flag.Parse()

	if *urlPtr == defaultUrlFlag {
		fmt.Println("Флаг не найден, будет использоваться значение по умолчанию:", defaultUrlFlag)
	}

	if *dirPtr == defaultDirFlag {
		fmt.Println("Флаг не найден, будет использоваться значение по умолчанию:", defaultDirFlag)
	}

	urls := getUrlsFromFile(*urlPtr)
	htmls := parseUrl(urls)
	saveHtmls(htmls, *dirPtr)

	fmt.Println("Время завершение программы", time.Since(start))
}

// saveHtmls - сохраняет html файлы в директори
func saveHtmls(htmlData []string, pathDir string) {
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		fmt.Println("Папки не существует\nБудет создана новая папка")

		// 0777 это права доступа
		err := os.MkdirAll(pathDir, 0777)
		if err != nil {
			fmt.Println("Неудалось создать директорию")
			os.Exit(1)
		}
	}

	// сохрание данных html в файлы
	for indexHtml, html := range htmlData {
		pathFile := pathDir + "/" + strconv.Itoa(indexHtml+1) + ".html"

		file, err := os.Create(pathFile)
		if err != nil {
			fmt.Println("Неудалось создать файл", err)
			continue
		}

		_, err = file.WriteString(html)
		if err != nil {
			fmt.Println("Ошибка при записи в файл", err)
			continue
		}

		file.Close()
	}

	fmt.Println("Файлы созданы")
}

// parseUrl - парсит url и возвращает срез html
func parseUrl(urls []string) []string {
	htmlData := []string{}

	for _, url := range urls {
		if url == "" {
			fmt.Println("Неверный адрес:", url)
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Неверный адрес:", url)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Не удалось считать ответ", err)
			continue
		}

		htmlData = append(htmlData, string(body))
	}

	return htmlData
}

// getUrlsFromFile - Возвращает url из файла
func getUrlsFromFile(path string) []string {
	fileData, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Нет файла с url-ами")
		os.Exit(1)
	}

	urls := strings.Split(string(fileData), "\n")

	return urls
}
