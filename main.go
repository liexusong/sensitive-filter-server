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

	host    = flag.String("host", "", "The server host")
	port    = flag.Int("port", 8080, "The server port")
	srcFile = flag.String("source", "./keywords", "The keywords source file")
	debug   = flag.Bool("debug", false, "Open debug mode")

	respErrText = map[int]string{
		200: "OK",
		201: "Arguments Invalid",
		202: "Not Found",
	}
)

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
		"resp": respErrText[code],
	})
}

func MatchFirst(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespJSON(ctx, 201)
		return
	}

	values := dict.MatchAll(body, 1)
	if len(values) > 0 {
		RespJSON(ctx, 200, values[0])
		return
	}

	RespJSON(ctx, 202)
}

func MatchAll(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespJSON(ctx, 201)
		return
	}

	values := dict.MatchAll(body, 0)
	if len(values) > 0 {
		RespJSON(ctx, 200, values)
		return
	}

	RespJSON(ctx, 202)
}

func Exists(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		RespJSON(ctx, 201)
		return
	}

	RespJSON(ctx, 200, dict.Exists(body))
}

func AddKeyword(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	if len(keyword) == 0 {
		RespJSON(ctx, 201)
		return
	}

	if dict.AddKeyword(keyword) {
		RespJSON(ctx, 200)
		return
	}

	RespJSON(ctx, 202)
}

func DelKeyword(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	if len(keyword) == 0 {
		RespJSON(ctx, 201)
		return
	}

	if dict.DelKeyword(keyword) {
		RespJSON(ctx, 200)
		return
	}

	RespJSON(ctx, 202)
}

func main() {
	flag.Parse()

	if err := dict.LoadWordsFromFile(*srcFile); err != nil {
		log.Fatalf("Failed to load keywords source file error: %v\n", err)
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

	if err := router.Run(fmt.Sprintf("%s:%d", *host, *port)); err != nil {
		log.Fatalf("Failed to run HTTP server error: %v\n", err)
	}
}
