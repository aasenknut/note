package http_test

import (
	"testing"

	"github.com/aasenknut/note/http"
	"github.com/aasenknut/note/mock"
)

type Server struct {
	*http.Server

	NoteService mock.NoteService
	UserService mock.UserService
	AuthService mock.AuthService
}

func OpenServer(tb testing.TB) *Server {
	tb.Helper()

	s := &Server{Server: http.NewServer()}
	s.Server.NoteService = &s.NoteService
	s.Server.UserService = &s.UserService
	s.Server.AuthService = &s.AuthService

	if err := s.Open(); err != nil {
		tb.Fatal(err)
	}

	return s
}

func CloseServer(tb testing.TB, s *Server) {
	tb.Helper()
	if err := s.Close(); err != nil {
		tb.Fatal(err)
	}
}
