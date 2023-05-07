package http

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/aasenknut/note"
	"github.com/aasenknut/note/ui"
)

// Server represents an HTTP server, and wraps all HTTP functionality.
type Server struct {
	ln     net.Listener
	server *http.Server
	router *http.ServeMux

	Addr        string
	NoteService note.NoteService
	TmplCache   map[string]*template.Template
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: http.NewServeMux(),
	}

	// Wrap router to do customisation which can't be done by middleware.
	s.server.Handler = http.HandlerFunc(s.serverHTTP)

	// Assets: .js, .css,...// Assets: .js, .css,...
	s.registerAssets()

	s.registerHealthRoute()
	s.registerNoteRoutes()
	return s
}

func (s *Server) Open() error {
	log.Printf("starting server on: %v", s.Addr)
	s.server.Addr = s.Addr

	var err error
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	go s.server.Serve(s.ln)
	return nil
}

func (s *Server) Close() error {
	if err := s.server.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Server) serverHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

const tmplIndex string = "html/index.tmpl"
const tmplPartials string = "html/partials/*.tmpl"
const pageDir string = "html/page/"

func (s *Server) SetTmplCache() error {
	cache := make(map[string]*template.Template, 0)

	paths, err := fs.Glob(ui.Files, pageDir+"*.tmpl")
	if err != nil {
		return fmt.Errorf("html from filesystem: %v", err)
	}

	for _, path := range paths {
		page := filepath.Base(path)
		patterns := []string{
			tmplIndex,
			tmplPartials,
			pageDir + page,
		}

		ts, err := template.New(page).ParseFS(ui.Files, patterns...)
		if err != nil {
			return fmt.Errorf("new template: %v", err)
		}
		cache[page] = ts
	}

	s.TmplCache = cache

	return nil
}