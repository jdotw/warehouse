package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginTransport struct {
	s Service
}

func NewGinTransport(s Service) Transport {
	return ginTransport{
		s: s,
	}
}

func (t ginTransport) getCategories(c *gin.Context) {
	v, err := t.s.GetCategories(context.Background())
	ok(err)
	c.JSON(http.StatusOK, v)
}

func (t ginTransport) Serve() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/categories", t.getCategories)
	r.Run(":8080")
}
