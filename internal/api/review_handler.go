package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/htojiddinov77-png/Articles/internal/store"
	"github.com/htojiddinov77-png/Articles/internal/utils"
)

type ReviewHandler struct {
	reviewStore store.ReviewStore
	articleStore store.ArticleStore
	logger      *log.Logger
}

func NewReviewHandler(reviewstore store.ReviewStore, articlestore store.ArticleStore, logger *log.Logger) *ReviewHandler {
	return &ReviewHandler{
		reviewStore: reviewstore,
		articleStore: articlestore,
		logger:      logger,
	}
}

func (rh *ReviewHandler) HandleCreateReview(w http.ResponseWriter, r *http.Request) {
	var review store.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		rh.logger.Printf("ERROR: decodingreview %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	if review.Rating < 1 || review.Rating > 5 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "rating must be between 1 and 5"})
		return
	}
	existingArticle, err := rh.articleStore.GetArticleById(int64(review.ArticleId))
	if err != nil {
		rh.logger.Printf("ERROR: getting article: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to verify article"})
		return
	}
	if existingArticle == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "article not found"})
		return
	}

	createdReview, err := rh.reviewStore.CreateReview(&review)
	if err != nil {
		rh.logger.Printf("ERROR: createReview: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create review"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"review": createdReview})
}

func (rh *ReviewHandler) HandleGetReviewByid(w http.ResponseWriter, r *http.Request) {
	reviewId, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid review id"})
		return
	}

	review, err := rh.reviewStore.GetReviewById(reviewId)
	if err != nil {
		rh.logger.Printf("ERROR: getReviewByid: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if review == nil {
		http.NotFound(w, r)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"review": review})
}

func (rh *ReviewHandler) HandleUpdateReviewById(w http.ResponseWriter, r *http.Request) {
	reviewId, err := utils.ReadIDParam(r)
	if err != nil{
		rh.logger.Printf("Error reading review ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid review ID"})
		return
	}

	existingReview, err := rh.reviewStore.GetReviewById(reviewId)
	if err != nil {
		rh.logger.Printf("Error getting review by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if existingReview == nil {
		http.NotFound(w, r)
		return
	}


	type UpdatedReviewRequest struct {
		ReviewText *string `json:"review_text"`
		Rating     *int    `json:"rating"`
	}

	validate := func(req *UpdatedReviewRequest) error {
		if req.Rating != nil {
			if *req.Rating < 1 || *req.Rating > 5 {
				return fmt.Errorf("rating must be between 1 and 5")
			}
		}
		return nil
	}
	var req UpdatedReviewRequest 
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // Decode is checking rating integer
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	err = validate(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	
	if req.ReviewText != nil {
		existingReview.ReviewText = *req.ReviewText
	}
	
	if req.Rating != nil {
		existingReview.Rating = *req.Rating
	}

	

	err = rh.reviewStore.UpdateReview(existingReview)
	if err != nil {
		rh.logger.Printf("Error updating review: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"review": existingReview})
}

func (rh *ReviewHandler) HandleDeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Printf("Error reading review Id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid review ID"})
		return
	}

	err = rh.reviewStore.DeleteReview(reviewID)
	if err != nil {
		if errors.Is(err, store.ErrReviewNotfound) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Review not found"})
			return
		}

		rh.logger.Printf("Error deleting review: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return	
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
