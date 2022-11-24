package main

import (
	"io"
	"os"
	"sync"
)

type GenHashFile struct {
	FullPath string
	WG       *sync.WaitGroup
	Coders   []HashCoder
}

type FileResp struct {
	FileName string
	FullPath string
	Hashes   []*HashResp
}

type HashResp struct {
	CodeName string
	HashVal  string
}

func (_self *GenHashFile) Generate() []*HashResp {
	file, err := os.Open(_self.FullPath)
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

		_self.WriteToChan(tmp[:n])

		if err == io.EOF {
			break
		}

	}

	_self.CloseChan()

	//需要先关闭通道，然后等待写入全部完成
	_self.WG.Wait()

	return _self.GetAllHash()
}

func (_self *GenHashFile) LoadAllCoders() {
	//_self.LoadCoders(new(Md5Coder), new(SHA1Coder), new(SHA256Coder), new(SHA512Coder))
	_self.LoadCoders(GetAllCoders()...)
}

func (_self *GenHashFile) LoadCoders(_coders ...HashCoder) {
	_self.Coders = _coders
	//
	_self.WG = new(sync.WaitGroup)

	for _, coder := range _self.Coders {
		coder.Create()
		//使用多协程同时读取
		coder.ReadFromChan(_self.WG)
	}
}

func (_self *GenHashFile) WriteToChan(input []byte) {
	for _, coder := range _self.Coders {
		coder.WriteToChan(input)
	}
}

func (_self *GenHashFile) CloseChan() {
	for _, coder := range _self.Coders {
		coder.CloseChan()
	}
}

func (_self *GenHashFile) GetAllHash() []*HashResp {
	var results []*HashResp
	for _, coder := range _self.Coders {

		results = append(results, &HashResp{CodeName: coder.Name(), HashVal: coder.GenHashHex()})
	}

	return results
}
