package model

type Article struct {
	ID          int    `json:"id"`
	AdminUserID int    `json:"admin_user_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
