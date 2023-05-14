package http

import "net/http"

func (s *Server) validateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("token")
		if tok != "" {
			userID, err := s.AuthService.GetUserID(r.Context(), tok)
			if err != nil {
				LogError(r, err)
			}
		}
		next.ServeHTTP(w, r)
	})
}
