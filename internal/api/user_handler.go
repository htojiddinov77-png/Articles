package api

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/htojiddinov77-png/Articles/internal/store"
	"github.com/htojiddinov77-png/Articles/internal/tokens"
	"github.com/htojiddinov77-png/Articles/internal/utils"
)

type ChangePasswordRequest struct {
	CurrentPassword    string `json:"current_password"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore  store.UserStore
	tokenStore store.TokenStore
	logger     *log.Logger
}

func NewUserHandler(userstore store.UserStore, tokenStore store.TokenStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore:  userstore,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

func (r *ChangePasswordRequest) Validate() error {
	if r.CurrentPassword == "" {
		return errors.New("current password is required")
	}

	if r.NewPassword == "" {
		return errors.New("new password is required")
	}

	if r.NewPassword != r.NewPasswordConfirm {
		return errors.New("new passwords don't match")
	}

	if r.NewPassword == r.CurrentPassword {
		return errors.New("new password cannot be the same as the current password")
	}

	return nil
}

func (r *registerUserRequest) validateRegisterRequest() error {
	if r.Username == "" {
		return errors.New("username is required")
	}

	if len(r.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (uh *UserHandler) HandlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	email := req.Email

	if email == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "email is required"})
		return
	}

	user, err := uh.userStore.GetUserByEmail(email)
	if err != nil {
		uh.logger.Printf("Error getting user by email: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if user == nil {
		// need feedback
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "User not found"})
		return
	}

	token, err := uh.tokenStore.CreateNewToken(user.ID, 10*time.Minute, tokens.ScopePasswordReset)
	if err != nil {
		uh.logger.Printf("Error creating  reset token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message": "password reset token generated",
		"token":   token,
	})
}

func (uh *UserHandler) HandlePasswordReset(w http.ResponseWriter, r *http.Request) {
	// Extract token from URL
	plainTextToken := chi.URLParam(r, "token")
	if plainTextToken == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "token is required"})
		return
	}

	hash := sha256.Sum256([]byte(plainTextToken))
	tokenHash := hash[:]

	// Lookup token in db
	token, err := uh.tokenStore.GetTokenByHash(tokenHash)
	if err != nil {
		uh.logger.Printf("Error getting token by hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if token == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid or expired token"})
		return
	}

	if token.Scope != tokens.ScopePasswordReset {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid token scope"})
		return
	}

	if time.Now().After(token.Expiry) {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Token expired"})
		return
	}

	var req struct {
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if req.NewPassword == "" || req.NewPassword != req.ConfirmPassword {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "passwords do not match"})
		return
	}

	user, err := uh.userStore.GetUserById(int64(token.UserID))
	if err != nil {
		uh.logger.Printf("Error getting user by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if user == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "User not found"})
		return
	}

	err = user.PasswordHash.Set(req.NewPassword)
	if err != nil {
		uh.logger.Printf("Error hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	err = uh.userStore.UpdateUser(user)
	if err != nil {
		uh.logger.Printf("Error updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if err := uh.tokenStore.DeleteAllTokensForUser(user.ID, token.Scope); err != nil {
		uh.logger.Printf("Error deleting token: %v", err)
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "password updated successfully"})

}

func (uh *UserHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Printf("Error reading user Id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	var req ChangePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	err = req.Validate()
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	oldUserPassword, err := uh.userStore.GetUserById(userId)
	if err != nil {
		uh.logger.Printf("Error getting user by ID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if oldUserPassword == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "User not found"})
		return
	}
	

	passwordsDomatch, err := oldUserPassword.PasswordHash.Matches(req.CurrentPassword)
	if err != nil {
		uh.logger.Printf("Error matching passwords: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if !passwordsDomatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Current password is incorrect"})
		return
	}

	err = oldUserPassword.PasswordHash.Set(req.NewPassword)
	if err != nil {
		uh.logger.Printf("Error hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	err = uh.userStore.UpdateUser(oldUserPassword)
	if err != nil {
		uh.logger.Printf("Error updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "password updated successfully"})
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = req.validateRegisterRequest()
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	existingUser, _ := uh.userStore.GetUserByUsername(req.Username)
	if existingUser != nil {
		uh.logger.Printf("ERROR: duplicaste entry: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Username is already exists"})
		return
	}

	existingEmailUser, err := uh.userStore.GetUserByEmail(req.Email)
	if err != nil {
		uh.logger.Printf("ERROR: getting user by email %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingEmailUser != nil {
		uh.logger.Printf("ERROR: duplicate entry: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Email is already exists"})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
		Bio:      req.Bio,
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: while creating user %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
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

func (uh *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
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
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Bio      *string `json:"bio"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedUserRequest)
	if err != nil {
		uh.logger.Printf("Error decoding update request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}



	if updatedUserRequest.Username != nil {

		if strings.TrimSpace(*updatedUserRequest.Username) == "" {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Username cannot be empty"})
			return
		}

		userByUsername, err := uh.userStore.GetUserByUsername(*updatedUserRequest.Username)
		if err != nil {
			uh.logger.Println("Error checking username:", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
			return
		}
		
		if userByUsername != nil && userByUsername.ID != existingUser.ID {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Username is already exists"})
			return
		}
		
		existingUser.Username = *updatedUserRequest.Username
	}
	if updatedUserRequest.Email != nil {
		if strings.TrimSpace(*updatedUserRequest.Email) == ""{
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
			return
		}

		userByEmail, err := uh.userStore.GetUserByEmail(*updatedUserRequest.Email)
		if err != nil {
			uh.logger.Println("Error getting username: ", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
			return
		}

		if userByEmail != nil && userByEmail.ID != existingUser.ID {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Email is already exists"})
			return
		}
		
		existingUser.Email = *updatedUserRequest.Email
	}

	if updatedUserRequest.Bio != nil {
		existingUser.Bio = *updatedUserRequest.Bio
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
