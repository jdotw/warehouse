package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jdotw/warehouse/services/stock-go/internal/service"
	"github.com/jdotw/warehouse/services/stock-go/internal/util"
)

type GetCategoriesHandler struct {
	s service.Service
}

func (h GetCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.s.GetCategories(context.Background())
	util.Ok(err)
	json, err := json.Marshal(v)
	util.Ok(err)
	fmt.Fprint(w, string(json))
}

type httpTransport struct {
	s service.Service
}

func NewHTTPTransport(s service.Service) Transport {
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
