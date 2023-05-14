package http

import "net/http"

func (s *Server) registerAdminRoutes() {
	s.router.HandleFunc("/admin/remove", s.basicAuth(s.adminRemoveNotes))
}

func (s *Server) adminRemoveNotes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pw, ok := r.BasicAuth()
		if !ok || user != s.AdminUser || pw != s.AdminPassword {
			status := http.StatusUnauthorized
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(status)
			w.Write([]byte(http.StatusText(status)))
			return
		}
		next.ServeHTTP(w, r)
	}
}
