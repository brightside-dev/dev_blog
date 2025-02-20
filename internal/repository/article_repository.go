package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/brightside-dev/dev-blog/database/client"
	"github.com/brightside-dev/dev-blog/internal/model"
)

type ArticleRepository interface {
	Create(ctx context.Context, tx *sql.Tx, article *model.Article) error
	GetAll(ctx context.Context, offset int) ([]model.Article, error)
	GetByID(ctx context.Context, id int) (*model.Article, error)
}

type articleRepository struct {
	db client.DatabaseService
}

func NewArticleRepository(
	db client.DatabaseService,
) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, tx *sql.Tx, article *model.Article) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO users (admin_user_id, title, excerpt, content) VALUES (?, ?, ?, ?, ?, ?)",
		article.AdminUserID, article.Title, article.Excerpt, article.Content)
	if err != nil {
		return err
	}

	return nil
}

func (r *articleRepository) GetAll(ctx context.Context, offset int) ([]model.Article, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(
			&article.ID,
			&article.AdminUserID,
			&article.Title,
			&article.Excerpt,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}
	return articles, nil
}

func (r *articleRepository) GetByID(ctx context.Context, id int) (*model.Article, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM articles WHERE id = ?", id)

	var article model.Article
	if err := row.Scan(
		&article.ID,
		&article.AdminUserID,
		&article.Title,
		&article.Excerpt,
		&article.Content,
		&article.CreatedAt,
		&article.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("article not found: %w", err)
	}

	return &article, nil
}
