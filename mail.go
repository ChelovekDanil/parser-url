package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	urlPtr := flag.String("url", "urls.txt", "путь к файлу с url")
	dirPtr := flag.String("dir", "htmls", "путь к папку с html страницами")

	flag.Parse()

	urls := readFile(*urlPtr)
	htmls := parseUrl(urls)
	saveHtmls(htmls, *dirPtr)
}

func saveHtmls(htmls []string, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Folder not exitst\nCreating new folder")

		// 0755 это права доспута
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	for i, v := range htmls {
		pathFile := path + "/" + strconv.Itoa(i+1) + ".html"

		file, err := os.Create(pathFile)
		if err != nil {
			panic(err)
		}

		_, err = file.WriteString(v)
		if err != nil {
			fmt.Println("Error while writing file")
			continue
		}

		file.Close()
	}

	fmt.Println("Files created")
}

func parseUrl(urls []string) []string {
	validUrls := []string{}

	for _, url := range urls {
		if url == "" {
			fmt.Println("bad address:", url)
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("bad address:", url)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		validUrls = append(validUrls, string(body))
	}

	return validUrls
}

func readFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("cant read file with urls")
		panic(err)
	}

	urls := strings.Split(string(data), "\n")

	return urls
}
