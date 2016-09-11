package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// searchPath  要找尋檔案的路徑
	searchPath = flag.String("input", "", "search file dir EX: /User/xxx/ddd/....")
	// searchFileName 要找尋的檔案名稱
	searchFileName = flag.String("file", "", "search file name")
	outputPath     = flag.String("output", "", "output file dir")
)

// 檔案路徑
var filelists []searchFile

var findFile *searchFile

// searchFile 搜尋的檔案結構
type searchFile struct {
	fileName string
	filePath string
	topPath  string
	fileType string
}

func main() {
	flag.Parse()
	if flag.NFlag() != 3 {
		flag.PrintDefaults()
	} else {
		body()
	}
}

func body() {

	// 列出所輸入的參數
	fmt.Println("Search Path : ", *searchPath)
	fmt.Println("Search File Name : ", *searchFileName)
	fmt.Println("Output File Path :", *outputPath)
	// 判斷路徑是否存在
	if *searchFileName == "null" {
		log.Fatal("Please Input the file name")
	}
	if err := pathExist(*searchPath); err != nil {
		log.Fatal("Error:", err)
	}

	// 開始進入目錄路徑

	if err := WalkDir(*searchPath); err != nil {
		fmt.Println("Not Find the Same File")
	} else {
		err := Copy(findFile)
		if err != nil {
			log.Println("Copy Failed")
		} else {
			log.Println("Copy Success")
		}
	}
}

// 拷貝資料到output之下
func Copy(fileinfo *searchFile) error {
	// 檢查Output目錄內是否存在相同目錄，如果存在則複製一份至該目錄
	// 否則建立一個新目錄
	outputDir := *outputPath + fileinfo.topPath
	if _, err := os.Stat(outputDir); err != nil {
		direrr := os.MkdirAll(outputDir, os.ModePerm)
		if direrr != nil {
			log.Fatal("Create Folder Failed ", direrr)
		} else {
			fmt.Println("Create Folder Success")
		}
	}
	// 複製檔案
	srcFile, err := os.Open(fileinfo.filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	outputFilePath := outputDir + fileinfo.fileName
	desFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// 查詢該目錄下所有的目錄與檔案
func WalkDir(path string) error {
	err := filepath.Walk(*searchPath, func(path string, f os.FileInfo, err error) error {

		// 如果取得檔案不是目錄，則紀錄
		if !f.IsDir() {
			// 取得上層路徑
			topPath := strings.TrimSuffix(strings.TrimPrefix(path, *searchPath), f.Name())
			tmp := searchFile{
				fileName: f.Name(),
				filePath: path,
				topPath:  topPath,
			}

			filelists = append(filelists, tmp)
			// 如果遍歷的檔案是所要則顯示出來並且退出
			if f.Name() == *searchFileName {
				findFile = &tmp
				return nil
			}
		}
		return nil
	})
	return err
}

// 列出該目錄下所有的檔案
func CheckFile(path string) []string {

	return nil
}

// pathExist 檢查路徑是否存在
func pathExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}
