package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/htojiddinov77-png/Articles/internal/store"
	"github.com/htojiddinov77-png/Articles/internal/utils"
)

type ReviewHandler struct {
	reviewStore store.ReviewStore
	logger      *log.Logger
}

func NewReviewHandler(reviewstore store.ReviewStore, logger *log.Logger) *ReviewHandler {
	return &ReviewHandler{
		reviewStore: reviewstore,
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

	if review.ArticleId <=0 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid article id"})
	}

	if review.Rating < 1 || review.Rating > 5 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "rating must be between 1 and 5"})
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


	var UpdatedReviewRequest struct {
		ReviewText *string `json:"review_text"`
		Rating     *int    `json:"rating"`
	}

	err = json.NewDecoder(r.Body).Decode(&UpdatedReviewRequest)
	if err != nil {
		rh.logger.Printf("Error decoding update request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	if UpdatedReviewRequest.ReviewText != nil {
		existingReview.ReviewText = *UpdatedReviewRequest.ReviewText
	}

	if UpdatedReviewRequest.Rating != nil {
    if *UpdatedReviewRequest.Rating < 1 || *UpdatedReviewRequest.Rating > 5 {
        utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "rating must be between 1 and 5"})
        return
    }
    existingReview.Rating = *UpdatedReviewRequest.Rating
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
		rh.logger.Printf("Error deleting review: %v", err)
		utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{"error": "error deleting review"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
