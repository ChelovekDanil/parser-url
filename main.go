package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
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

	var wg sync.WaitGroup

	// сохрание данных html в файлы
	for indexHtml, html := range htmlData {
		wg.Add(1)

		go func(i int, html string) {
			defer wg.Done()

			pathFile := pathDir + "/" + strconv.Itoa(i+1) + ".html"

			file, err := os.Create(pathFile)
			if err != nil {
				fmt.Println("Неудалось создать файл", err)
				return
			}
			defer file.Close()

			_, err = file.WriteString(html)
			if err != nil {
				fmt.Println("Ошибка при записи в файл", err)
				return
			}
		}(indexHtml, html)
	}
	wg.Wait()

	fmt.Println("Файлы созданы")
}

// parseUrl - парсит url и возвращает срез html
func parseUrl(urls []string) []string {
	htmlData := []string{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			if url == "" {
				fmt.Println("Неверный адрес:", url)
				return
			}

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Неверный адрес:", url)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Не удалось считать ответ", err)
				return
			}

			mu.Lock()
			htmlData = append(htmlData, string(body))
			mu.Unlock()
		}(url)
	}
	wg.Wait()

	return htmlData
}

// getUrlsFromFile - возвращает url из файла
func getUrlsFromFile(path string) []string {
	fileData, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Нет файла с url-ами")
		os.Exit(1)
	}

	urls := strings.Split(string(fileData), "\n")

	return urls
}
