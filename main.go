package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sf/config"
	"sf/src"
	"sf/src/threading"
	"strconv"
)

func downloadBook(input any) {
	var (
		bookId   string
		bookList []string
	)
	switch input.(type) {
	case string:
		bookId = input.(string)
	case int:
		bookId = strconv.Itoa(input.(int))
	case []string:
		bookList = input.([]string)
	}
	if bookId != "" && bookList == nil {
		BookData := src.GetBookDetailed(bookId)
		fmt.Printf("开始下载:%s\n", BookData.NovelName)
		if err := ioutil.WriteFile(fmt.Sprintf("save/%v.txt", BookData.NovelName),
			[]byte(BookData.NovelName), 0777); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		config.NewFile(fmt.Sprintf("cache/%v.json", BookData.NovelName))
		src.GetCatalogue(BookData)
	} else if bookList != nil && len(bookList) > 0 {

		ThreadLocks := threading.NewGoLimit(5)
		for _, bookId = range bookList {
			ThreadLocks.Add()
			go func(bookId string, t *threading.GoLimit) {
				defer ThreadLocks.Done()
				BookData := src.GetBookDetailed(bookId)
				fmt.Printf("开始下载:%s\n", BookData.NovelName)
				if err := ioutil.WriteFile(fmt.Sprintf("save/%v.txt", BookData.NovelName),
					[]byte(BookData.NovelName), 0777); err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				config.NewFile(fmt.Sprintf("cache/%v.json", BookData.NovelName))
				src.GetCatalogue(BookData)
			}(bookId, ThreadLocks)
		}
		ThreadLocks.WaitZero()
	}

}

func main() {
	config.NewMyJsonPro()
	if len(os.Args) >= 2 {
		inputs := os.Args[1:]
		switch {
		case inputs[0] == "l", inputs[0] == "login":
			if len(inputs) >= 3 {
				src.LoginAccount(inputs[1], inputs[2])
			} else {
				fmt.Println("parameters are not enough")
			}
		case inputs[0] == "d", inputs[0] == "download":
			if len(inputs) >= 2 {
				downloadBook(inputs[1])
			} else {
				fmt.Println("parameters are not enough")
			}
		}
	} else {
		fmt.Println("please input parameters, like: sf login username password")
	}
}