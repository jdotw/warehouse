package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GetCategoriesHandler struct {
	r Repository
}

func (h GetCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.r.GetCategories(context.Background())
	if err != nil {
		log.Panicln(err)
	}
	json, err := json.Marshal(v)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Fprint(w, string(json))
}

type httpTransport struct {
	r Repository
}

func NewHTTPTransport(r Repository) Transport {
	return httpTransport{
		r: r,
	}
}

func (t httpTransport) Serve() {
	gc := &GetCategoriesHandler{
		r: t.r,
	}
	http.Handle("/categories", gc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
