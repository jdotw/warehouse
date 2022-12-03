package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginTransport struct {
	r Repository
}

func NewGinTransport(r Repository) Transport {
	return ginTransport{
		r: r,
	}
}

func (t ginTransport) getCategories(c *gin.Context) {
	v, err := t.r.GetCategories(context.Background())
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, v)
}

func (t ginTransport) Serve() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/categories", t.getCategories)
	r.Run(":8080")
}
