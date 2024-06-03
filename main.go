package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	targetFile = flag.String("f", "", "指定具体文件名 或 文件夹路径")
	useCoders  = flag.String("c", "MD5,SHA-1,SHA-256,SHA-512", "指定需要生成的算法种类，默认为所有算法都计算")
	reportFile = flag.String("o", "default", "指定生成Hash结果的文件名，默认生成随机文件名")
	outType    = flag.String("ot", "default", "生成报告文件的内容类型，支持参数 default、only-hash")
	helpTxt    = flag.String("helpTxt", "", `以 Windows 系统下为例

自动生成目录下所有文件的指纹信息，不带任何参数，直接执行（推荐）：
FileHashCode.exe

指定单个文件：
FileHashCode.exe -f "d:\检材目录\检材1.docx"

指定文件夹下所有文件：
FileHashCode.exe -f "d:\检材目录\视频文件\"

生成目录下所有文件，只使用两种算法：
FileHashCode.exe -c "MD5,SHA-256"

指定一个文件，只使用两种算法：
FileHashCode.exe -c "MD5,SHA-256" -f "d:\检材目录\检材1.docx"

注意： -c 参数不指定时，默认所有支持的算法都会生成

当前支持的算法有：
MD5,SHA-1,SHA-256,SHA-512

指定输出结果文件的文件名 report.txt，如果不指定则自动随机文件名：
FileHashGen.exe -f test.zip -c "SHA-256" -o report.txt

指定输出结果文件的文件名（内容仅含hash值）：
FileHashGen.exe -f test.zip -c "SHA-256" -ot "only-hash" -o test.zip.sha256


如果您有更进一步需求，请前往下面地址提交 issues
https://gitee.com/tavenli/FileHashGen
`)

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
FileHashGen v1.2
电子数据指纹生成工具

项目地址：
https://gitee.com/tavenli/FileHashGen

查看使用帮助：
FileHashCode -h

★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★
	`)

	coders := strings.Split(*useCoders, ",")
	for _, c := range coders {
		hashCoders = append(hashCoders, GetCodersByName(c))
	}

	exePath, _ := os.Executable()
	_, exeFileName = filepath.Split(exePath)
	//fmt.Println("exeFileName", exeFileName)

	//计算执行时间
	start := time.Now()

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
	reportContent := ""

	if *outType == "only-hash" {
		reportContent = outputResultForOnlyHash(fRespes)
	} else {
		reportContent = outputResultForDefault(fRespes)
	}

	reportOutput := fmt.Sprint("Hash-", time.Now().Format("20060102150405"), ".txt")
	if *reportFile != "default" {
		reportOutput = *reportFile
	}

	_ = os.WriteFile(reportOutput, []byte(reportContent), 0600)
	fmt.Println("\n\n生成指纹报告文件：" + reportOutput)

	terminal := time.Since(start)
	fmt.Println("\nTimeCost：", terminal)

	//var tpIn string
	//fmt.Println("\nOK，指纹信息已生成，按 [Enter] 键完成本次任务。")
	//fmt.Scanln(&tpIn)
}

func outputResultForDefault(fRespes []*FileResp) string {
	var txtStr strings.Builder

	for _, f := range fRespes {
		fmt.Println("\n文件名：", f.FullPath)

		txtStr.WriteString(fmt.Sprint("\n文件名：", f.FullPath))
		for _, h := range f.Hashes {
			fmt.Println(fmt.Sprint(h.CodeName, "：", h.HashVal))
			txtStr.WriteString(fmt.Sprint("\n", h.CodeName, "：", h.HashVal))
		}

		txtStr.WriteString("\n-----------------------------------------------------")

	}

	fmt.Println("\n文件总数：", len(fRespes))
	txtStr.WriteString(fmt.Sprint("\n\n文件总数：", len(fRespes), "\n"))

	return txtStr.String()
}

func outputResultForOnlyHash(fRespes []*FileResp) string {
	var txtStr strings.Builder

	for _, f := range fRespes {
		fmt.Println("\n文件名：", f.FullPath)

		for _, h := range f.Hashes {
			fmt.Println(fmt.Sprint(h.CodeName, "：", h.HashVal))
			txtStr.WriteString(fmt.Sprint(h.HashVal, "\n"))
		}

	}

	fmt.Println("\n文件总数：", len(fRespes))

	return txtStr.String()
}

// go不支持三元表达式，可以使用自定义的函数实现
// 例如：max := If(x > y, x, y).(int)
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

// 判断文件/或文件夹是否存在
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
	//fmt.Println("IgnoreFile check：", fileName)
	if fileName == exeFileName {
		return true
	}

	match, _ := regexp.MatchString("^Hash-.+\\.txt$", fileName)
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
