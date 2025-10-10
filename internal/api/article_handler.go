package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ArticleHandler struct{}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (ah *ArticleHandler) HandleGetArticleByID(w http.ResponseWriter, r *http.Request) {
	paramsArticleId := chi.URLParam(r, "id")
	if paramsArticleId == ""{
		http.NotFound(w, r)
		return
	}

	articleID, err := strconv.ParseInt(paramsArticleId, 10, 64)
	if err != nil{
		http.NotFound(w ,r)
		return
	}

	fmt.Fprintf(w, "this is the article id %d\n", articleID)
}

func (ah *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "created a article\n")
}