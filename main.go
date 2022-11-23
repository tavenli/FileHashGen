package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

/*
	支持MD5、SHA-1、SHA-256、SHA-512算法等
*/
func main() {
	fmt.Println("Hi Taven.Li")

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

	fmt.Println("ok")

}
