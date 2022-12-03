package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinTransport struct {
	r Repository
}

func NewGinTransport(r Repository) GinTransport {
	return GinTransport{
		r: r,
	}
}

func (t GinTransport) getCategories(c *gin.Context) {
	v, err := t.r.GetCategories(context.Background())
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, v)
}

func (t GinTransport) Serve() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/categories", t.getCategories)
	r.Run(":8080")
}
