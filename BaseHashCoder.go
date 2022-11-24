package main

import (
	"encoding/hex"
	"fmt"
	"hash"
	"sync"
)

type BaseHashCoder struct {
	datas   chan []byte
	Encoder hash.Hash
}

func (_self *BaseHashCoder) Name() string {
	return "-"
}

func (_self *BaseHashCoder) Create() {

}

func (_self *BaseHashCoder) ReadFromString(input string) {
	_self.ReadFromBytes([]byte(input))
}

func (_self *BaseHashCoder) ReadFromBytes(input []byte) {
	_self.Encoder.Write(input)
}

func (_self *BaseHashCoder) ReadFromChan(wg *sync.WaitGroup) {

	wg.Add(1)
	go func() {
		for data := range _self.datas {
			_self.Encoder.Write(data)
		}
		wg.Done()
	}()
}

func (_self *BaseHashCoder) WriteToChan(input []byte) {
	_self.datas <- input
}

func (_self *BaseHashCoder) CloseChan() {
	close(_self.datas)
}

func (_self *BaseHashCoder) GenHashHex() string {

	if _self.Encoder == nil {
		fmt.Errorf("no encoder")
		return ""
	}

	hash := hex.EncodeToString(_self.Encoder.Sum(nil))
	return hash
}
