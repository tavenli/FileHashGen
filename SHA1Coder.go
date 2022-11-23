package main

import (
	"crypto/sha1"
)

type SHA1Coder struct {
	BaseHashCoder
}

func (_self *SHA1Coder) Create() {
	_self.Encoder = sha1.New()
	_self.datas = make(chan []byte, 10)
}
