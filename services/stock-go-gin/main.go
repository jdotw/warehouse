package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewCategory struct {
	Name string `json:"name"`
}

func getCategories(c *gin.Context) {
	var Catgegories = []Category{
		{ID: "1", Name: "Blue Train"},
		{ID: "2", Name: "Jeru"},
		{ID: "3", Name: "Sarah Vaughan and Clifford Brown"},
	}
	c.JSON(http.StatusOK, Catgegories)
}

func createCategory(c *gin.Context) {
	var newCategory NewCategory
	if err := c.BindJSON(&newCategory); err != nil {
		fmt.Printf("error: %v+\n", err)
		return
	}
	c.JSON(http.StatusAccepted, newCategory)
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/categories", getCategories)
	r.POST("/categories", createCategory)
	r.Run("localhost:8080")
}
