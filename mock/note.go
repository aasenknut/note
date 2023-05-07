package mock

import (
	"context"

	"github.com/aasenknut/note"
)

type NoteService struct {
	GetAllNotesFn func(ctx context.Context) ([]*note.Note, error)
	GetNoteByIDFn func(ctx context.Context, id int) (*note.Note, error)
	CreateNoteFn  func(ctx context.Context, note *note.Note) error
	DeleteNoteFn  func(ctx context.Context, id int) error
}

func (n *NoteService) GetAllNotes(ctx context.Context) ([]*note.Note, error) {
	return n.GetAllNotesFn(ctx)
}

func (n *NoteService) GetNoteByID(ctx context.Context, id int) (*note.Note, error) {
	return n.GetNoteByIDFn(ctx, id)
}
func (n *NoteService) CreateNote(ctx context.Context, note *note.Note) error {
	return n.CreateNoteFn(ctx, note)
}
func (n *NoteService) DeleteNote(ctx context.Context, id int) error {
	return n.DeleteNoteFn(ctx, id)
}
