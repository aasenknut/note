package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aasenknut/note"
)

const homeRoute = "/home"
const noteRoute = "/note/"
const createRoute = "/note/create"

func (s *Server) registerNoteRoutes() {
	s.router.HandleFunc("/", s.handleHome)
	s.router.HandleFunc(homeRoute, s.handleHome)
	s.router.HandleFunc(noteRoute, s.handleNoteView)
	s.router.HandleFunc(createRoute, s.validateSession(s.handleNoteCreate))
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	notes, err := s.NoteService.GetAllNotes(r.Context())
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		LogError(r, err)
		return
	}
	data := pageData{Notes: notes}
	err = s.render(w, "home.tmpl", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(notes)
		LogError(r, err)
	}
}

func (s *Server) handleNoteView(w http.ResponseWriter, r *http.Request) {
	rawID := readID(r.URL.Path)
	id, err := convertID(rawID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	note, err := s.NoteService.GetNoteByID(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	data := pageData{Note: note}
	err = s.render(w, "note.tmpl", data)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}

func (s *Server) handleNoteCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form: %v", err)
		}
		nt := &note.Note{
			Title: r.FormValue("title"),
			Text:  r.FormValue("text"),
		}
		log.Printf("attempting to insert note: %+v", nt)
		insertedNote, err := s.NoteService.CreateNote(r.Context(), nt)
		log.Printf("insirted note: %+v", insertedNote)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			LogError(r, err)
			return
		}
		http.Redirect(w, r, homeRoute, http.StatusPermanentRedirect)
	case http.MethodGet:
		data := pageData{
			Note: &note.Note{
				Title: "",
				Text:  "",
			},
		}
		err := s.render(w, "create.tmpl", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			LogError(r, err)
			return
		}
		status := http.StatusOK
		w.WriteHeader(status)
		w.Write([]byte(http.StatusText(status)))
	}
}

func readID(u string) string {
	urlStr := strings.Split(u, "/")
	id := urlStr[len(urlStr)-1]
	return id
}

func convertID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("can not convert ID to int")
	}
	return id, nil
}
