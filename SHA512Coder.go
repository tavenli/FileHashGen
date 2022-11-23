package main

import (
	"crypto/sha512"
)

type SHA512Coder struct {
	BaseHashCoder
}

func (_self *SHA512Coder) Create() {
	_self.Encoder = sha512.New()
	_self.datas = make(chan []byte, 10)
}
