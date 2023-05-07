package http

import (
	"net/http"

	"github.com/aasenknut/note/ui"
)

func (s *Server) registerAssets() {
	fileServer := http.FileServer(http.FS(ui.Files))
	s.router.Handle("/assets/", fileServer)
}
