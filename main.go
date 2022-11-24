package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	help       = flag.Bool("h", false, "显示帮助信息")
	targetFile = flag.String("f", "", "指定文件名 或 目录路径")
	useCoders  = flag.String("c", "MD5,SHA-1,SHA-256,SHA-512", "指定文件名 或 目录路径")

	hashCoders  []HashCoder
	exeFileName = "FileHashGen"
)

/*
	支持MD5、SHA-1、SHA-256、SHA-512算法等
*/
func main() {

	flag.Parse()
	//flag.PrintDefaults()

	fmt.Println(`
★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★
FileHashGen v1.0
电子数据指纹生成工具

项目地址：
https://gitee.com/tavenli/FileHashGen

自动生成目录下所有文件的指纹信息，不带任何参数，直接执行（推荐）：
FileHashCode.exe

指定单个文件：
FileHashCode.exe -f "d:\检材目录\检材1.docx"

指定文件夹下所有文件：
FileHashCode.exe -f "d:\检材目录\视频文件\"

指定算法种类（默认是所有算法都计算）：
FileHashCode.exe -c "MD5,SHA-256"

FileHashCode.exe -c "MD5,SHA-256" -f "d:\检材目录\检材1.docx"

★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★
	`)

	coders := strings.Split(*useCoders, ",")
	for _, c := range coders {
		hashCoders = append(hashCoders, GetCodersByName(c))
	}

	exePath, _ := os.Executable()
	_, exeFileName = filepath.Split(exePath)
	fmt.Println("exeFileName", exeFileName)

	var fRespes []*FileResp

	if len(*targetFile) == 0 {
		workPath, _ := os.Getwd()
		fmt.Println(workPath)

		//默认从当前目录开始遍历
		fRespes = walkDirectory(workPath)

	} else {
		//有指定参数
		fInfo, err := os.Stat(*targetFile)
		if err != nil {
			fmt.Errorf("文件 %s 不存在\n", *targetFile)
			return
		}

		if fInfo.IsDir() {
			//指定了文件夹
			fRespes = walkDirectory(*targetFile)
			fmt.Println("results count：", len(fRespes))
		} else {
			//指定了单个文件
			fileResp := new(FileResp)
			fileResp.FileName = fInfo.Name()
			fileResp.FullPath = *targetFile

			genHashFile := &GenHashFile{FullPath: *targetFile}
			//genHashFile.LoadAllCoders()
			genHashFile.LoadCoders(hashCoders...)
			results := genHashFile.Generate()

			fileResp.Hashes = results

			fRespes = append(fRespes, fileResp)
		}

	}

	//输出结果
	var txtStr strings.Builder

	for _, f := range fRespes {
		fmt.Println("\n文件名：", f.FullPath)
		txtStr.WriteString("\n\n\n-----------------------------------------------------")
		txtStr.WriteString(fmt.Sprint("\n\n文件名：", f.FullPath))
		for _, h := range f.Hashes {
			fmt.Println(fmt.Sprint(h.CodeName, "：", h.HashVal))
			txtStr.WriteString(fmt.Sprint("\n", h.CodeName, "：", h.HashVal))
		}
	}

	fmt.Println("文件总数：", len(fRespes))
	txtStr.WriteString(fmt.Sprint("\n\n文件总数：", len(fRespes), "\n"))

	reportOutput := fmt.Sprint("FileHash-", time.Now().Format("20060102150405"), ".txt")
	_ = ioutil.WriteFile(reportOutput, []byte(txtStr.String()), 0600)
	fmt.Println("\n\n生成指纹报告文件：" + reportOutput)

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

//  判断文件/或文件夹是否存在
func FileIsExist(name string) bool {
	var exist = true
	if _, err := os.Stat(name); os.IsNotExist(err) {
		exist = false
	}
	return exist
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

			if IgnoreFile(fileInfo.Name()) {
				//忽略的文件
				return nil
			}

			fileResp := new(FileResp)
			fileResp.FileName = fileInfo.Name()
			fileResp.FullPath = curFullPath

			genHashFile := &GenHashFile{FullPath: curFullPath}
			//genHashFile.LoadAllCoders()
			genHashFile.LoadCoders(hashCoders...)
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
	fmt.Println("IgnoreFile check：", fileName)
	if fileName == exeFileName {
		return true
	}

	match, _ := regexp.MatchString("^FileHash-.+\\.txt$", fileName)
	if match {
		return true
	}

	return false
}

func GetAllCoders() []HashCoder {

	return []HashCoder{new(Md5Coder), new(SHA1Coder), new(SHA256Coder), new(SHA512Coder)}
}

func GetCodersByName(name string) HashCoder {
	//MD5,SHA-1,SHA-256,SHA-512
	switch name {
	case "MD5":
		return new(Md5Coder)
	case "SHA-1":
		return new(SHA1Coder)
	case "SHA-256":
		return new(SHA256Coder)
	case "SHA-512":
		return new(SHA512Coder)
	default:
		return new(Md5Coder)
	}

}
