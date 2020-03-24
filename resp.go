package main

import (
	"github.com/gin-gonic/gin"
)

const (
	ErrOK         = 200
	ErrArgInvalid = 201
	ErrNotFound   = 202
)

var rspErrTxt = map[int]string{
	ErrOK:         "OK",
	ErrArgInvalid: "Arguments Invalid",
	ErrNotFound:   "Not Found",
}

func RespJSON(ctx *gin.Context, code int, data... interface{}) {
	var (
		value interface{}
	)

	if len(data) > 0 {
		value = data[0]
	}

	ctx.JSON(200, gin.H{
		"code": code,
		"data": value,
		"resp": rspErrTxt[code],
	})
}
