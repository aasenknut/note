package http

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const createUser string = "/user/create"

const tokenLifetime time.Duration = 24 * time.Hour

func (s *Server) registerUserRoutes() {
	s.router.HandleFunc(createUser, s.signup)
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("signup decode: %v", err)
			status := http.StatusBadRequest
			http.Error(w, http.StatusText(status), status)
		}
		user := r.FormValue("username")
		pw := r.FormValue("password")
		hashedPW, err := bcrypt.GenerateFromPassword([]byte(pw), 16)
		if err != nil {
			log.Printf("generating hashed password: %v", err)
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		userID, err := s.UserService.Create(r.Context(), user, string(hashedPW))
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			LogError(r, err)
			return
		}
		token, err := generateToken()
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			LogError(r, err)
			return
		}
		ttl, err := s.AuthService.SetAuth(r.Context(), token, userID)
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			LogError(r, err)
			return
		}
		cookie := &http.Cookie{
			Name:    "session",
			Value:   token,
			Expires: time.Now().Add(ttl).UTC(),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	case http.MethodGet:
		if err := s.render(w, "signup.tmpl", pageData{}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			LogError(r, err)
			return
		}
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
