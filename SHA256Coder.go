package main

import (
	"crypto/sha256"
)

type SHA256Coder struct {
	BaseHashCoder
}

func (_self *SHA256Coder) Create() {
	_self.Encoder = sha256.New()
	_self.datas = make(chan []byte, 10)
}
