package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

/*
	支持MD5、SHA-1、SHA-256、SHA-512算法等
*/
func main() {

	fmt.Println("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")
	fmt.Println("FileHashGen v1.0")
	fmt.Println("电子数据指纹生成工具")
	fmt.Println("")
	fmt.Println("项目地址：")
	fmt.Println("https://gitee.com/tavenli/FileHashGen")
	fmt.Println("")
	fmt.Println("★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★")

	workPath, _ := os.Getwd()
	fmt.Println(workPath)

	//默认从当前目录开始遍历
	files := walkDirectory(workPath)

	fmt.Println("results count：", len(files))
	for _, f := range files {
		fmt.Println("\n文件名：", f.FullPath)
		for _, h := range f.Hashes {
			fmt.Println(fmt.Sprint(h.CodeName, "：", h.HashVal))
		}
	}

}

//	go不支持三元表达式，可以使用自定义的函数实现
//	例如：max := If(x > y, x, y).(int)
func If(condition bool, trueVal, falseVal interface{}) interface{} {

	if condition {
		return trueVal
	}
	return falseVal
}

// 一行代码计算代码执行时间
// defer TimeCost(time.Now())
func TimeCost(start time.Time) {
	terminal := time.Since(start)
	fmt.Println("TimeCost：", terminal)
}

func walkDirectory(toWalk string) []*FileResp {
	var files []*FileResp

	walkErr := filepath.WalkDir(toWalk, func(curFullPath string, curFile os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		//fmt.Println(curFullPath)
		if !curFile.IsDir() {
			//是文件，则生成hash
			fileInfo, _ := curFile.Info()
			//fmt.Println(fileInfo.Name())

			fileResp := new(FileResp)
			fileResp.FileName = fileInfo.Name()
			fileResp.FullPath = curFullPath

			genHashFile := &GenHashFile{FullPath: curFullPath}
			genHashFile.LoadAllCoders()
			results := genHashFile.Generate()

			fileResp.Hashes = results

			files = append(files, fileResp)
		}

		return nil
	})

	if walkErr != nil {
		fmt.Println(walkErr.Error())
	}

	return files
}

func IgnoreFile(fileName string) bool {

	return true
}
