package server

import (
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	`aurora/initialize`
	`aurora/initialize/router/api`
	"aurora/middlewares"
)

func (s *Server) RegisterRouter() *gin.Engine {

	handler := NewHandle(
		checkProxy(),
	)

	router := gin.Default()
	//router.Use(middlewares.Cors)
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	router.GET("/3", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	router.POST("/upload", initialize.AddRag)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.OPTIONS("/v1/chat/models", optionsHandler)
	authGroup := router.Group("").Use(middlewares.Authorization)
	authGroup.POST("/v1/chat/completions", handler.duckduckgo)
	authGroup.GET("/v1/models", handler.engines)

	group := router.Group("/app/v1")
	group.POST("/question", api.QuestionController{}.Question)
	group.GET("/questionStream", api.QuestionController{}.QuestionStream)

	group.POST("/template", api.TemplateController{}.Template)
	group.POST("/imagetotext", api.ImageToTextController{}.ImageToText)
	group.POST("/embedding", api.EmbeddingController{}.Embedding)
	group.POST("/knowledge", api.KnowledgeController{}.Knowledge)
	group.POST("/queryKnowledge", api.QueryKnowledgeController{}.QueryKnowledge)

	rr := router.Group("/rss/v1")
	rr.GET("/feeds", s.FeedList)

	rr.POST("/feeds", s.FeedAdd)
	rr.POST("/opml/import", s.OPMLImport)
	//r.For("/opml/import", s.handleOPMLImport)
	//r.For("/opml/export", s.handleOPMLExport)
	return router
}
