package service

import (
	"context"

	"github.com/brightside-dev/dev-blog/internal/repository"
)

type ArticleService interface {
	Create(ctx context.Context, data []map[string]any) error
	GetAll(ctx context.Context, offset int) ([]map[string]any, error)
	GetById(ctx context.Context, id int) ([]map[string]any, error)
}

type articleService struct {
}

func NewArticleService(articeRepository repository.ArticleRepository) ArticleService {
	return &articleService{}
}

func (s *articleService) Create(ctx context.Context, data []map[string]any) error {
	return nil
}

func (s *articleService) GetAll(ctx context.Context, offset int) ([]map[string]any, error) {
	return []map[string]any{}, nil
}

func (s *articleService) GetById(ctx context.Context, id int) ([]map[string]any, error) {
	return []map[string]any{}, nil
}
