package transport

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jdotw/warehouse/services/stock-go/internal/service"
	"github.com/jdotw/warehouse/services/stock-go/internal/util"
)

type ginTransport struct {
	s service.Service
}

func NewGinTransport(s service.Service) Transport {
	return ginTransport{
		s: s,
	}
}

func (t ginTransport) getCategories(c *gin.Context) {
	v, err := t.s.GetCategories(context.Background())
	util.Ok(err)
	c.JSON(http.StatusOK, v)
}

func (t ginTransport) Serve() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/categories", t.getCategories)
	r.Run(":8080")
}
