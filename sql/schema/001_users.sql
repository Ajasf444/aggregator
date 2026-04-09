-- +goose Up
CREATE TABLE (
  id int primary key,
  created_at timestamp not null,
  updated_at timestamp not null,
  name string unique not null
);

-- +goose Down
DROP TABLE users;
