package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sarvesh-10/ReadEazeBackend/internal/domain"
	"github.com/Sarvesh-10/ReadEazeBackend/utility"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service *UserService
	logger  *utility.Logger
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Response struct {
	Message string `json:"message"`
	UserId  int64  `json:"userId"`
}

func NewUserHandler(service *UserService, logger *utility.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) error {
	h.service.logger.Info("Signing up user")
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.service.logger.Error("Error decoding user: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err

	}
	h.logger.Info("Decoded user successfully", user)
	h.service.logger.Info("generating token successfully")
	id, token, refreshToken, err := h.service.RegisterUser(user)
	h.service.logger.Info("User token successfully")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/refresh", // optional: scope to only refresh route
		HttpOnly: true,
		Secure:   false, // ✅ Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 7 days
	})
	resp := Response{Message: "User registered successfully", UserId: id}
	response, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	return nil
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	user, err := h.service.GetByEmail(loginRequest.Email)
	h.logger.Info("HERE IN GET USER BY EMAIL", user)

	if err != nil {
		fmt.Printf("HERE in teh get uesr fails")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	er := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	fmt.Printf("Here after compare and hash")

	if er != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return er
	}
	fmt.Printf("Error 1")
	token, err := utility.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	h.logger.Info("HERE IN ERROR 1")
	h.logger.Info("HERE IN ERROR 1")
	h.logger.Info("Generation successful")
	fmt.Println("Setting cookie: token=", token)

	refreshToken, err := utility.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token: "+err.Error(), http.StatusInternalServerError)
		return err
	}
	h.service.refreshRepo.Store(user.ID, utility.HashToken(refreshToken), time.Now().Add(time.Hour*24*7))

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Change to `true` in HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // ✅ Required for cross-origin cookies
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false, // Change to `true` in HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // ✅ Required for cross-origin cookies
	})

	h.logger.Info("Erro 3")
	resp := map[string]string{"message": "User logged in successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return nil
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire instantly
		HttpOnly: true,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

func (h *UserHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("Refreshing session", r)
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return err
	}
	refresh_token := cookie.Value
	if refresh_token == "" {
		http.Error(w, "Refresh token is empty", http.StatusUnauthorized)
		return err
	}
	accessToken, newRefreshToken, err := h.service.RefreshSession(refresh_token)
	if err != nil {
		http.Error(w, "Failed to refresh session: "+err.Error(), http.StatusInternalServerError)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false, // Change to `true` in HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // ✅ Required for cross-origin cookies
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   false, // Change to `true` in HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // ✅ Required for cross-origin cookies
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session refreshed successfully"))
	return nil

}
