package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aasenknut/note"
)

type NoteService struct {
	db *DB
}

func NewNoteService(db *DB) *NoteService {
	return &NoteService{db: db}
}

func (n *NoteService) GetAllNotes(ctx context.Context) ([]*note.Note, error) {
	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	notes, err := getAllNotes(ctx, tx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func getAllNotes(ctx context.Context, tx *Tx) ([]*note.Note, error) {
	q := `SELECT * FROM note`
	rows, err := tx.QueryContext(ctx, q)
	notes := make([]*note.Note, 0)
	for rows.Next() {
		var note note.Note
		if err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Text,
			&note.Created,
		); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}
	if rows.Err() != nil {
		return nil, err
	}
	if tx.Commit(); err != nil {
		return nil, fmt.Errorf("get tx: %v", err)
	}
	return notes, nil
}

func (n *NoteService) GetNoteByID(ctx context.Context, id int) (*note.Note, error) {
	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	note, err := getNote(ctx, id, tx)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func getNote(ctx context.Context, id int, tx *Tx) (*note.Note, error) {
	q := `SELECT *
	FROM note
	WHERE id=?`
	nt := note.Note{}
	err := tx.QueryRowContext(ctx, q, id).Scan(
		&nt.ID,
		&nt.Title,
		&nt.Text,
		&nt.Created,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("note not found: %v", err)
		default:
			return nil, err
		}
	}
	if tx.Commit(); err != nil {
		return nil, fmt.Errorf("get tx: %v", err)
	}
	return &nt, nil
}

func (n *NoteService) CreateNote(ctx context.Context, nt *note.Note) (*note.Note, error) {
	stmt := `INSERT INTO note (title, text) VALUES (?, ?) RETURNING *;`
	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create note: %v", err)
	}
	var insrt note.Note
	args := []any{nt.Title, nt.Text}
	err = tx.QueryRowContext(ctx, stmt, args...).Scan(
		&insrt.ID,
		&insrt.Title,
		&insrt.Text,
		&insrt.Created,
	)
	if err != nil {
		return nil, fmt.Errorf("create note: %v", err)
	}
	if tx.Commit(); err != nil {
		return nil, fmt.Errorf("create tx: %v", err)
	}
	return &insrt, nil
}
func (n *NoteService) DeleteNote(ctx context.Context, id int) error {
	return fmt.Errorf("_ not implemented _")
}
