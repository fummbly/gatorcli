-- +goose Up
CREATE TABLE feed_follows (
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id INTEGER NOT NULL,
  feed_id INTEGER NOT NULL,
  UNIQUE(user_id, feed_id),
  FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,
  FOREIGN KEY(feed_id) REFERENCES feeds(id)
    ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feed_follows;
