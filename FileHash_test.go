package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestTxtContent(t *testing.T) {
	md5HashCoder := new(Md5Coder)
	md5HashCoder.Create()
	md5HashCoder.ReadFromString("123456")
	fmt.Println(md5HashCoder.GenHashHex())

	t.Log("ok")
}

func TestDirFiles(t *testing.T) {
	//
	path, _ := os.Executable()
	fmt.Println(path)

	_, exec := filepath.Split(path)
	fmt.Println(exec)

	workPath, _ := os.Getwd()
	fmt.Println(workPath)

	//默认从当前目录开始遍历
	files := _walkDirectory(workPath)

	fmt.Println("results count：", len(files))
	for _, f := range files {
		fmt.Println(f.FullPath)
		for _, h := range f.Hashes {
			fmt.Println(fmt.Sprint(h.CodeName, "：", h.HashVal))
		}
	}

	t.Log("ok")
}

func _walkDirectory(toWalk string) []*FileResp {
	var files []*FileResp

	walkErr := filepath.WalkDir(toWalk, func(curFullPath string, curFile os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		//fmt.Println(curFullPath)
		if curFile.IsDir() {
			//fmt.Println(curFile.Name())
		} else {
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

func TestSingleFile(t *testing.T) {
	defer TimeCost(time.Now())

	//filePath := "C:\\Users\\HJKL\\Desktop\\cn_sql_server_2012_enterprise_edition_x86_x64_dvd_813295.iso"
	filePath := "C:\\Users\\HJKL\\Desktop\\payload.txt"
	genHashFile := &GenHashFile{FullPath: filePath}
	genHashFile.LoadAllCoders()
	results := genHashFile.Generate()
	for _, h := range results {
		fmt.Println(h)
	}

	t.Log("ok")
}

func TestBigFile(t *testing.T) {
	defer TimeCost(time.Now())
	filePath := "C:\\Users\\HJKL\\Desktop\\cn_sql_server_2012_enterprise_edition_x86_x64_dvd_813295.iso"
	//filePath := "C:\\Users\\HJKL\\Desktop\\payload.txt"
	/*
		file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			panic(err.Error())
		}
		fi, _ := file.Stat()
		fsize := fi.Size()

	*/

	//主要是考虑大文件，要并行计算

	wg := new(sync.WaitGroup)

	md5HashCoder := new(Md5Coder)
	md5HashCoder.Create()
	md5HashCoder.ReadFromChan(wg)

	sha1HashCoder := new(SHA1Coder)
	sha1HashCoder.Create()
	sha1HashCoder.ReadFromChan(wg)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	data := make([]byte, 4_194_304)
	for {
		n, err := file.Read(data)
		if err != nil && err != io.EOF {
			panic(err.Error())
			break
		}

		// 需要复制一份，否则goroutines 共同对 data 对象的操作可能出现问题
		tmp := make([]byte, len(data))
		copy(tmp, data)

		md5HashCoder.WriteToChan(tmp[:n])
		sha1HashCoder.WriteToChan(tmp[:n])

		if err == io.EOF {
			break
		}

	}

	md5HashCoder.CloseChan()
	sha1HashCoder.CloseChan()

	//需要先关闭通道，然后等待写入全部完成
	wg.Wait()

	hash := md5HashCoder.GenHashHex()
	fmt.Println(hash)

	fmt.Println(sha1HashCoder.GenHashHex())

	t.Log("ok")
}
