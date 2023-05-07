package http

import "net/http"

func (s *Server) registerHealthRoute() {
	s.router.HandleFunc("/health", s.healthHandler)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}
