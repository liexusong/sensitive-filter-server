package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	dict = NewDict()

	host   = flag.String("host", "", "Server host")
	port   = flag.Int("port", 8080, "Server port")
	source = flag.String("source", "./keywords", "Keywords source file")
)

func RespOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": data,
		"resp": "OK",
	})
}

func RespErr(ctx *gin.Context, code int, resp string) {
	ctx.JSON(200, gin.H{
		"code": code,
		"data": nil,
		"resp": resp,
	})
}

func MatchFirst(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespErr(ctx, 201, "Parameters invalid")
		return
	}

	values := dict.MatchAll(body, 1)
	if len(values) > 0 {
		RespOK(ctx, values[0])
		return
	}

	RespErr(ctx, 202, "Cannot found keywords")
}

func MatchAll(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespErr(ctx, 201, "Parameters invalid")
		return
	}

	values := dict.MatchAll(body, 0)
	if len(values) > 0 {
		RespOK(ctx, values)
		return
	}

	RespErr(ctx, 202, "Cannot found keywords")
}

func AddWord(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	if len(keyword) == 0 {
		RespErr(ctx, 201, "Parameters invalid")
		return
	}

	if dict.AddWord(keyword) {
		RespOK(ctx, nil)
		return
	}

	RespErr(ctx, 202, "Cannot add keyword")
}

func DelWord(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	if len(keyword) == 0 {
		RespErr(ctx, 201, "Parameters invalid")
		return
	}

	if dict.DelWord(keyword) {
		RespOK(ctx, nil)
		return
	}

	RespErr(ctx, 202, "Cannot delete keyword")
}

func main() {
	flag.Parse()

	if err := dict.LoadWordsFile(*source); err != nil {
		log.Fatalf("Failed to load keywords source file error: %v\n", err)
		return
	}

	router := gin.New()

	router.POST("/api/match_first", MatchFirst)
	router.POST("/api/match_all", MatchAll)
	router.GET("/api/add_word", AddWord)
	router.GET("/api/del_word", DelWord)

	if err := router.Run(fmt.Sprintf("%s:%d", *host, *port)); err != nil {
		log.Fatalf("Failed to run HTTP server error: %v\n", err)
	}
}
