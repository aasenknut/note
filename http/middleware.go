package http

import "net/http"

func (s *Server) registerMiddleware() {
	s.server.Handler = s.validateSession(s.server.Handler)
}

func (s *Server) validateSession(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := ""
		for _, v := range r.Cookies() {
			if v.Name == "session" {
				tok = v.Value
				break
			}
		}
		if tok == "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		_, err := s.AuthService.GetUserID(r.Context(), tok)
		if err != nil {
			LogError(r, err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

		next.ServeHTTP(w, r)
	})
}
