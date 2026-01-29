package sqlite

import (
	"context"
	"database/sql"

	"github.com/eventuallyconsistentwrites/high-tide-server/internal/post"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Create(ctx context.Context, p *post.Post) error {
	return nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*post.Post, error) {
	return nil, nil
}

func (s *Store) GetAll(ctx context.Context) ([]*post.Post, error) {
	return nil, nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	return nil
}
