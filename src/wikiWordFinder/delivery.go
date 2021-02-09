package main

import (
	"io/ioutil"
	"net/http"
)

type Transport interface {
	Get(path string) (string, error)
	count() int
}

type HTTP struct {
	_count int
}

func NewHTTP() *HTTP {
	return &HTTP{_count: 0}
}

func (t *HTTP) Get(path string) (string, error) {
	t._count++
	res, err := http.Get(path)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (t HTTP) count() int {
	return t._count
}
