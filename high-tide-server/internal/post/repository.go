package post

import "context"

type Repository interface {
	Create(ctx context.Context, p *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
	GetAll(ctx context.Context) ([]*Post, error)
	Delete(ctx context.Context, id int64) error
}
