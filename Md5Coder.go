package main

import (
	"crypto/md5"
)

type Md5Coder struct {
	BaseHashCoder
}

func (_self *Md5Coder) Create() {
	_self.Encoder = md5.New()
	_self.datas = make(chan []byte, 10)
}
