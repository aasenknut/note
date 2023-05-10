package http

import "net/http"

func (s *Server) validateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("token")
		s.AuthService.GetUserID(r.Context(), tok)
		next.ServeHTTP(w, r)
	})
}
