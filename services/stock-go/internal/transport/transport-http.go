package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GetCategoriesHandler struct {
	s Service
}

func (h GetCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.s.GetCategories(context.Background())
	ok(err)
	json, err := json.Marshal(v)
	ok(err)
	fmt.Fprint(w, string(json))
}

type httpTransport struct {
	s Service
}

func NewHTTPTransport(s Service) Transport {
	return httpTransport{
		s: s,
	}
}

func (t httpTransport) Serve() {
	gc := &GetCategoriesHandler{
		s: t.s,
	}
	http.Handle("/categories", gc)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
