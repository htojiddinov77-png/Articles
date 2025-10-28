package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/htojiddinov77-png/Articles/internal/store"
	"github.com/htojiddinov77-png/Articles/internal/utils"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userstore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userstore,
		logger:    logger,
	}
}

func (uh *UserHandler) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("Error reading user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}
	user, err := uh.userStore.GetUserById(userID)
	if err != nil {
		uh.logger.Printf("Error getting user by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uh.logger.Printf("Error while decoding: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	createdUser, err := uh.userStore.CreateUser(&user)
	if err != nil {
		uh.logger.Printf("Error creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": createdUser})
}

func (uh *UserHandler) HandelUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("Error reading user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	existingUser, err := uh.userStore.GetUserById(userID)
	if err != nil {
		uh.logger.Printf("Error getting user by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if existingUser == nil {
		http.NotFound(w, r)
		return
	}

	var updatedUserRequest struct {
		FirstName    *string `json:"firsname"`
		LastName     *string `json:"lastname"`
		Email        *string `json:"email"`
		PasswordHash *string `json:"password_hash"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedUserRequest)
	if err != nil {
		uh.logger.Printf("error while decoding user: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	if updatedUserRequest.FirstName != nil {
		existingUser.FirstName = *updatedUserRequest.FirstName
	}
	if updatedUserRequest.LastName != nil {
		existingUser.LastName = *updatedUserRequest.LastName
	}
	if updatedUserRequest.FirstName != nil {
		existingUser.FirstName = *updatedUserRequest.FirstName
	}
	if updatedUserRequest.Email != nil {
		existingUser.Email = *updatedUserRequest.Email
	}
	if updatedUserRequest.PasswordHash != nil {
		existingUser.PasswordHash = *updatedUserRequest.PasswordHash
	}

	err = uh.userStore.UpdateUser(existingUser)
	if err != nil {
		uh.logger.Printf("Error updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": existingUser})
}

func (uh *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("Error reading user Id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	err = uh.userStore.DeleteUser(userID)
	if err != nil {
		uh.logger.Printf("Error deleting user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
