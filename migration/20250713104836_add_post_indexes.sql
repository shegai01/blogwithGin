-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE INDEX idx_posts_author_id on posts(author_id);
CREATE INDEX idx_posts_created_at on posts(created_at);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX if EXISTS idx_posts_author_id;
DROP INDEX if EXISTS idx_posts_created_at;
-- +goose StatementEnd
