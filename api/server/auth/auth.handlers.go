package authhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"patoune-api/internal/auth"
)

type AuthHandlers struct {
	authService auth.AuthService
}

func NewAuthHandlers(authService auth.AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	var input auth.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid query body", http.StatusBadRequest)
		return
	}

	response, err := h.authService.Register(r.Context(), input)
	if err != nil {
		log.Printf("Error Register: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    response.Token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	var input auth.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid query body", http.StatusBadRequest)
		return
	}

	response, err := h.authService.Login(r.Context(), input)
	if err != nil {
		log.Printf("Error Login: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    response.Token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout success",
	})
}

func (h *AuthHandlers) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(auth.UserResponse{
		ID:    claims.UserID,
		Email: claims.Email,
	})
}
