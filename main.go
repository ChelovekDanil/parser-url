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

	urlFilePtr, dirHtmlPtr, err := addFlags()
	if err != nil {
		fmt.Println("Ошибка при добавления флагов:", err)
		return
	}

	urls, err := getUrlsFromFile(*urlFilePtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	htmls := getHtmlData(urls)

	err = saveHtmlsInDir(htmls, *dirHtmlPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeFinish := time.Since(start)
	fmt.Printf("Время завершение программы: %s\n", fmt.Sprintf("%d.%dms", timeFinish.Milliseconds(), timeFinish.Microseconds()/10000))
}

// addFlags - добавляет флаги
func addFlags() (*string, *string, error) {
	defaultUrlFlag := "urls.txt"
	defaultDirFlag := "htmls"

	urlPtr := flag.String("url", "", "путь к файлу с url")
	dirPtr := flag.String("dir", "", "путь к папку с html страницами")

	flag.Parse()
	flag.PrintDefaults()

	if *urlPtr == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при чтении корневого каталога: %s", err)
		}
		urlPtr = &defaultUrlFlag
		fmt.Printf("Должен быть установлен флаг --url, который отвечает за путь к файлу с url адресами.\nПуть по умолчанию: %s/%s\n\n", currentDir, defaultUrlFlag)
	}

	if *dirPtr == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при чтении корневого каталога: %s", err)
		}
		dirPtr = &defaultDirFlag
		fmt.Printf("Должен быть установлен флаг --dir, который отвечает за путь в каталогу куда будут загружины html файл.\nДиректория по умолчанию: %s\n\n", currentDir+"/"+defaultDirFlag)
	}

	return urlPtr, dirPtr, nil
}

// getUrlsFromFile - возвращает url из файла
func getUrlsFromFile(pathFileUrl string) ([]string, error) {
	fileData, err := os.ReadFile(pathFileUrl)
	if err != nil {
		return nil, fmt.Errorf("нет файла с url-ами: %s", err)
	}

	urls := strings.Split(string(fileData), "\n")

	return urls, nil
}

// getHtmlData - возвращает срез html полученных из url-ов
func getHtmlData(urls []string) []string {
	htmlDataSlice := []string{}

	for _, url := range urls {
		htmlData, err := parseUrl(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Успешный ответ: %s\n", url)
		htmlDataSlice = append(htmlDataSlice, htmlData)
	}

	return htmlDataSlice
}

// parseUrl - парсит url и возвращает html
func parseUrl(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("неверный адрес: %s", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка в запросе: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("не удалось считать ответ: %s", err)
	}

	return string(body), nil
}

// saveHtmls - сохраняет html файлы в директори
func saveHtmlsInDir(htmlData []string, pathDir string) error {
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("ошибка при чтении корневого каталога: %s", err)
		}

		fmt.Printf("Директории не существует\nБудет создана новая директория по пути: %s\n", currentDir)

		// 0777 это права доступа
		err = os.MkdirAll(pathDir, 0777)
		if err != nil {
			return fmt.Errorf("неудалось создать директорию: %s", err)
		}
	}

	// сохрание данных html в файлы
	for indexHtml, html := range htmlData {
		err := createHtmlFile(indexHtml, pathDir, html)
		if err != nil {
			fmt.Printf("неудалось создать файл: %s", err)
		}
	}

	fmt.Println("Файлы созданы")
	return nil
}

// createHtmlFile - создает файлы в отпредленной директории
func createHtmlFile(indexHtml int, pathDir string, html string) error {
	pathFile := fmt.Sprintf("%s/%s.html", pathDir, strconv.Itoa(indexHtml+1))

	file, err := os.Create(pathFile)
	if err != nil {
		return fmt.Errorf("неудалось создать файл: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		return fmt.Errorf("ошибка при записи в файл: %s", err)
	}

	return nil
}
