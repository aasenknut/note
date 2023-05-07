package note

import (
	"context"
	"time"
)

type Note struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
}

type NoteService interface {
	GetAllNotes(ctx context.Context) ([]*Note, error)
	GetNoteByID(ctx context.Context, id int) (*Note, error)
	CreateNote(ctx context.Context, note *Note) (*Note, error)
	DeleteNote(ctx context.Context, id int) error
}
