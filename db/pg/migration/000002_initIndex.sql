-- +migrate Up
CREATE INDEX IF NOT EXISTS idx_movies_deleted_at ON movies(deleted_at);
CREATE INDEX IF NOT EXISTS idx_cinemas_deleted_at ON cinemas(deleted_at);
CREATE INDEX IF NOT EXISTS idx_screens_deleted_at ON screens(deleted_at);
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON movies(deleted_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_movies_deleted_at;
DROP INDEX IF EXISTS idx_cinemas_deleted_at;
DROP INDEX IF EXISTS idx_screens_deleted_at;
DROP INDEX IF EXISTS idx_orders_deleted_at;