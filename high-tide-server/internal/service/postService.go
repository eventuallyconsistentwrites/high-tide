package service

import (
	"context"
	"time"

	"github.com/eventuallyconsistentwrites/high-tide-server/internal/domain"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/repository"
)

type PostService struct {
	repo *repository.PostSQLRepository
}

func NewPostService(repo *repository.PostSQLRepository) *PostService {
	return &PostService{repo: repo}
}

// Service methods that perform relevant DB manipulations using PostSQLRepository

func (s *PostService) GetPost(ctx context.Context, id int64) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]domain.Post, error) {
	return s.repo.GetAllPosts(ctx)
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.Post) error {
	post.Timestamp = time.Now()
	return s.repo.Create(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, post *domain.Post) error {
	return s.repo.DeletePost(ctx, post)
}
