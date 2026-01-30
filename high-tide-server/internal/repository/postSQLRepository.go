package repository

import (
	"context"

	"github.com/eventuallyconsistentwrites/high-tide-server/internal/domain"
	"gorm.io/gorm"
)

type PostSQLRepository struct {
	db *gorm.DB
}

// Creates *PostSQLRepository instance
func NewPostSQLRepository(db *gorm.DB) *PostSQLRepository {
	return &PostSQLRepository{db: db}
}

// Define operations on the table corresponding to the "Post" struct

func (s *PostSQLRepository) Create(ctx context.Context, post *domain.Post) error {
	err := s.db.WithContext(ctx).Create(post).Error
	return err
}

func (s *PostSQLRepository) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	var post domain.Post
	err := s.db.WithContext(ctx).First(&post, id).Error
	return &post, err
}

func (s *PostSQLRepository) GetAllPosts(ctx context.Context) ([]domain.Post, error) {
	var posts []domain.Post
	result := s.db.WithContext(ctx).Find(&posts)
	return posts, result.Error
}

func (s *PostSQLRepository) DeletePost(ctx context.Context, post *domain.Post) error {
	return s.db.WithContext(ctx).Delete(post).Error
}
