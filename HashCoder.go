package main

import "sync"

//Hash编码抽象接口
type HashCoder interface {
	Create()
	ReadFromString(input string)
	ReadFromBytes(input []byte)
	ReadFromChan(wg *sync.WaitGroup)
	WriteToChan(input []byte)
	CloseChan()
	GenHashHex() string
}
