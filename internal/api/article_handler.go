package api

import (
	"encoding/json"

	"log"
	"net/http"

	"github.com/htojiddinov77-png/Articles/internal/store"
	"github.com/htojiddinov77-png/Articles/internal/utils"
)

type ArticleHandler struct {
	articleStore store.ArticleStore
	logger       *log.Logger
}

func NewArticleHandler(articleStore store.ArticleStore, logger *log.Logger) *ArticleHandler {
	return &ArticleHandler{
		articleStore: articleStore,
		logger:       logger,
	}
}

func (ah *ArticleHandler) HandlerGetArticleById(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid article id"})
		return
	}

	article, err := ah.articleStore.GetArticleById(articleID)
	if err != nil {
		ah.logger.Printf("ERROR: getArticleByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"article": article})
}

func (ah *ArticleHandler) HandlerCreateArticle(w http.ResponseWriter, r *http.Request) {
	var article store.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		ah.logger.Printf("ERROR: decodingCreateArticle: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdArticle, err := ah.articleStore.CreateArticle(&article)
	if err != nil {
		ah.logger.Printf("ERROR: createArticle: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create article"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"article": createdArticle})
}

func (ah *ArticleHandler) HandleUpdateArticleById(w http.ResponseWriter, r *http.Request) {

	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid article update id"})
		return
	}

	existingArticle, err := ah.articleStore.GetArticleById(articleID)
	if err != nil {
		ah.logger.Printf("ERROR: getArticleById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingArticle == nil {
		http.NotFound(w, r)
		return
	}

	var UpdateArticleRequest struct {
		Title       *string           `json:"title"`
		Description *string           `json:"description"`
		Image       *string           `json:"image"`
		AuthorID    *int              `json:"author_id"`
		Paragraphs  []store.Paragraph `json:"paragraphs"`
	}

	err = json.NewDecoder(r.Body).Decode(&UpdateArticleRequest)
	if err != nil {
		ah.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if UpdateArticleRequest.Title != nil {
		existingArticle.Title = *UpdateArticleRequest.Title
	}

	if UpdateArticleRequest.Description != nil {
		existingArticle.Description = *UpdateArticleRequest.Description
	}

	if UpdateArticleRequest.Image != nil {
		existingArticle.Image = *UpdateArticleRequest.Image
	}

	if UpdateArticleRequest.AuthorID != nil {
		existingArticle.AuthorId = *UpdateArticleRequest.AuthorID
	}

	if UpdateArticleRequest.Paragraphs != nil {
		existingArticle.Paragraphs = *&UpdateArticleRequest.Paragraphs
	}

	err = ah.articleStore.UpdateArticle(existingArticle)
	if err != nil {
		ah.logger.Printf("ERROR: updatingArticle: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"article": existingArticle})
}

func (ah *ArticleHandler) HandleDeleteWorkoutbyId(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid article id"})
		return
	}

	existingArticle, err := ah.articleStore.GetArticleById(articleID)
	if err != nil {
		ah.logger.Printf("ERROR: getArticleByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingArticle == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	err = ah.articleStore.DeleteArticle(articleID)
	if err != nil {
		ah.logger.Printf("ERROR: deleteArticle: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to delete article"})
		return
	}
}
