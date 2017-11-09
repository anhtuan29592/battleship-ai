package main

import (
	"fmt"
	"github.com/anhtuan29592/paladin/config"
	"github.com/anhtuan29592/paladin/context"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func main() {

	// Init dependency injection
	var g inject.Graph
	var gameContext context.GameContext
	var notifyContext context.NotifyContext

	err := g.Provide(
		&inject.Object{Value: &gameContext},
		&inject.Object{Value: &notifyContext},
		&inject.Object{Value: config.DefaultRedis},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Init logger
	logFile, err := os.Create("/tmp/paladin.log")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// Init RESTful API
	router := gin.Default()
	paladin := router.Group("/paladin")
	{
		paladin.POST("/invite", gameContext.Invite)
		paladin.POST("/start", gameContext.Start)
		paladin.POST("/turn", gameContext.Turn)
		paladin.POST("/notify", notifyContext.Result)
	}
	router.Run(":8080")
}
