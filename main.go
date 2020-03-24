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
	sHost   = flag.String("host", "", "The server host")
	sPort   = flag.Int("port", 8080, "The server port")
	srcFile = flag.String("source", "./keywords", "The keywords source file")
	debug   = flag.Bool("debug", false, "Open debug mode")

	matcher *Dict
)

func MatchFirst(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespJSON(ctx, ErrArgInvalid)
		return
	}

	values := matcher.MatchAll(body, 1)
	if len(values) > 0 {
		RespJSON(ctx, ErrOK, values[0])
		return
	}

	RespJSON(ctx, ErrNotFound)
}

func MatchAll(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespJSON(ctx, ErrArgInvalid)
		return
	}

	values := matcher.MatchAll(body, 0)
	if len(values) > 0 {
		RespJSON(ctx, ErrOK, values)
		return
	}

	RespJSON(ctx, ErrNotFound)
}

func Exists(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespJSON(ctx, ErrArgInvalid)
		return
	}

	RespJSON(ctx, ErrOK, matcher.Exists(body))
}

func AddKeyword(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	if len(keyword) == 0 {
		RespJSON(ctx, ErrArgInvalid)
		return
	}

	if matcher.AddKeyword(keyword) {
		RespJSON(ctx, ErrOK)
		return
	}

	RespJSON(ctx, ErrNotFound)
}

func DelKeyword(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	if len(keyword) == 0 {
		RespJSON(ctx, ErrArgInvalid)
		return
	}

	if matcher.DelKeyword(keyword) {
		RespJSON(ctx, ErrOK)
		return
	}

	RespJSON(ctx, ErrNotFound)
}

func main() {
	flag.Parse()

	matcher = NewDict()

	if err := matcher.LoadWordsFromFile(*srcFile); err != nil {
		log.Fatalf("Failed to load keywords from file error: %v\n", err)
		return
	}

	if *debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	group := router.Group("/api")
	{
		group.POST("/match_first", MatchFirst)
		group.POST("/match_all", MatchAll)
		group.POST("/exists", Exists)
		group.GET("/add_keyword", AddKeyword)
		group.GET("/del_keyword", DelKeyword)
	}

	if err := router.Run(fmt.Sprintf("%s:%d", *sHost, *sPort)); err != nil {
		log.Fatalf("Failed to run HTTP server error: %v\n", err)
	}
}
