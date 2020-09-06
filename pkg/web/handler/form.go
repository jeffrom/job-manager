package handler

import "github.com/go-playground/form"

var formDecoder = form.NewDecoder()

func init() {
	formDecoder.SetTagName("json")
	formDecoder.SetMaxArraySize(1000)
}
