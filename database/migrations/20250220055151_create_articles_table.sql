-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,                
    admin_user_id BIGINT UNSIGNED NOT NULL,             
    title VARCHAR(255) NOT NULL,
    excerpt JSON NOT NULL,
    content JSON NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
    CONSTRAINT fk_articles_admin_user FOREIGN KEY (admin_user_id) REFERENCES admin_users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE articles;
-- +goose StatementEnd
